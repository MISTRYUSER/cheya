package main

import (
	"context"
	"log"
	"net"

	vehiclev1 "github.com/xuewentao/cheya/api/vehicle/v1"
	"github.com/xuewentao/cheya/apps/vehicle/ent"
	"github.com/xuewentao/cheya/apps/vehicle/server"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

func main() {
	//1.链接数据库
	dns := "host=localhost port=5432 user=wentao_xue  dbname=cheya password=Woe89132 sslmode=disable"
	client, err := ent.Open("postgres", dns)
	if err != nil {
		log.Fatalf("❌ failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	//2.自动迁移
	//在 db 中自动创建 vehicles
	log.Println("📦 Migrating database schema...")
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("❌ failed creating schema resources: %v", err)
	}
	log.Println("✅ Schema migrated successfully!")

	//3.启动 grpc
	lis, err := net.Listen("tcp",":50051")
	if err != nil{
		log.Fatalf("❌ failed to listen : %v", err)
	}

	s := grpc.NewServer()
	//注入 client 到 server
	vehiclev1.RegisterVehicleServiceServer(s,server.NewVehicleServer(*client))
	log.Printf("🚀 Vehicle Service is running on :50051")

	//4.启动服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
