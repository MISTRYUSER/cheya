package main

import (
	"context"
	"log"
	"time"

	vehiclev1 "github.com/xuewentao/cheya/api/vehicle/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
func main() {
	//1.连接 grpc 服务器
	//这里是 client 所以 server 是监听 这里是 Dial被弃用现在是 NewClient
	// WithTransportCredentials(insecure...) 表示不使用 TLS (仅限内网/开发)
	conn,err := grpc.NewClient("localhost: 50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil{
		log.Fatalf("❌failed to create client: %v",err)
	}
	defer conn.Close()

	//创建 client
	c := vehiclev1.NewVehicleServiceClient(conn)

	//设置超时 1s
	ctx,cancel := context.WithTimeout(context.Background(),time.Second)
	defer cancel()

	//mock 一辆新车
	timestamp := time.Now().Format("150405")
	vin := "VIN-TEST-" + timestamp
	plate := "沪A-" + timestamp

	log.Printf("🛠️  Creating Vehicle with VIN: %s ...", vin)

	createResp,err := c.CreateVehicle(ctx,&vehiclev1.CreateVehicleRequest{
		Vin: vin,
		LicensePlate: plate,
	})

	if err != nil {
		log.Printf("⚠️  Create failed: %v", err)
	} else {
		log.Printf("✅ Vehicle Created! DB_ID: %s", createResp.VehicleId)
	}

	//查询
	log.Printf("🔍 Querying Vehicle with VIN: %s ...", vin)
	GetResp,err := c.GetVehicle(ctx,&vehiclev1.GetVehicleRequest{
		VehicleId: vin,
	})
	if err != nil {
		log.Fatalf("❌ Get failed: %v", err)
	}
	//打印结果
	v := GetResp.Vehicle
	log.Printf("🎉 Found Vehicle:")
	log.Printf("   - VIN:   %s", v.Vin)
	log.Printf("   - Plate: %s", v.LicensePlate)
	// 这里的 Status 是个枚举值 (0, 1, 2)，打印出来是数字或 String (取决于 Protobuf 生成配置)
	log.Printf("   - Status: %v", v.Status)

}