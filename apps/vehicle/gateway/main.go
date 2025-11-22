package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	vehiclev1 "github.com/xuewentao/cheya/api/vehicle/v1"
)

func main() {
	//初始化 client  用网关来使用 http
	//生产环境一般使用服务发现
	conn ,err := grpc.NewClient("localhost: 50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil{
		log.Fatalf("❌ Failed to connect gRPC server %v",err)
	}
	defer conn.Close()
	//创建 grpc client 存根
	vehicleClient := vehiclev1.NewVehicleServiceClient(conn)	
	log.Println("✅ Connected to Vehicle Service(gRPC)")

	//2.初始化 Gin
	r := gin.Default()

	//定义路由 GET /api/vi/vehicles/:id
	r.GET("/api/v1/vehicles/:id",func(c *gin.Context){
		//获取 URL 参数
		vehicleID := c.Param("id")

		//设置超时上下文
		ctx,concel := context.WithTimeout(context.Background(),2 * time.Second)
		defer concel()

		//发起 gRPC 调用
		resp,err := vehicleClient.GetVehicle(ctx,&vehiclev1.GetVehicleRequest{
			VehicleId : vehicleID,
		})
		//错误处理
		if err != nil{
			log.Println("❌ gRPC called failed :%v",err)
			//返回 500 404
			c.JSON(http.StatusInternalServerError,gin.H{
				"error": err.Error(),
			})
			return
		}
		//成功响应
		c.JSON(http.StatusOK,gin.H{
			"code" : 200,
			"message" : "success",
			"data" : resp.Vehicle,
		})
	})
	//启动 HTTP 服务器
	log.Println("🚀 Gateway is running on :8080")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("❌ Failed to boost gateway :%v",err)
	}

}