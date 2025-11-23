package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket" // ✅ 新增
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	vehiclev1 "github.com/xuewentao/cheya/api/vehicle/v1"
)

// 简易连接池
var (
	clients   = make(map[*websocket.Conn]bool) // WebSocket 客户端连接池
	broadcast = make(chan string)              // 广播消息通道
	mutex     sync.Mutex                       // 保护 clients map 的互斥锁
)

// WebSocket upgrader 配置
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源（生产环境需要更严格的检查）
	},
}

func main() {
	//初始化 client  用网关来使用 http
	//生产环境一般使用服务发现
	conn, err := grpc.NewClient("localhost: 50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("❌ Failed to connect gRPC server %v", err)
	}
	defer conn.Close()
	//创建 grpc client 存根
	vehicleClient := vehiclev1.NewVehicleServiceClient(conn)
	log.Println("✅ Connected to Vehicle Service(gRPC)")
	//1.Redis 订阅
	go func() {
		rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
		log.Println("👂 Gateway subscribing to Redis channel: vehicle:update")

		sub := rdb.Subscribe(context.Background(), "vehicle:update")
		ch := sub.Channel()

		for msg := range ch {
			log.Printf("📩 Received from Redis: %s", msg.Payload)
			broadcast <- msg.Payload
		}
	}()
	//2.WebSocket 广播协程
	go func() {
		for {
			msg := <-broadcast
			log.Printf("📡 Broadcasting to %d clients", len(clients))

			mutex.Lock()
			for client := range clients {
				err := client.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					log.Printf("❌ WS Error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
			mutex.Unlock()
		}
	}()
	//3.初始化 Gin
	r := gin.Default()

	//定义路由 GET /api/vi/vehicles/:id
	r.GET("/api/v1/vehicles/:id", func(c *gin.Context) {
		//获取 URL 参数
		vehicleID := c.Param("id")

		//设置超时上下文
		ctx, concel := context.WithTimeout(context.Background(), 2*time.Second)
		defer concel()

		//发起 gRPC 调用
		resp, err := vehicleClient.GetVehicle(ctx, &vehiclev1.GetVehicleRequest{
			VehicleId: vehicleID,
		})
		//错误处理
		if err != nil {
			log.Printf("❌ gRPC called failed: %v", err)
			//返回 500
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		//成功响应
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    resp.Vehicle,
		})
	})

	//WebSocket 结构
	//ws 指的是 WebSocket 连接对象
	r.GET("/ws", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Fatalf("❌ WS Upgrade failed: %v", err)
			return
		}
		mutex.Lock()
		clients[ws] = true
		mutex.Unlock()
		log.Println("🔌 New Browser Connected!")

		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				mutex.Lock()
				delete(clients, ws)
				mutex.Unlock()
				break

			}
		}
	})

	// 提供静态文件（test.html）
	r.StaticFile("/test.html", "./test.html")
	r.StaticFile("/", "./test.html") // 根路径也返回 test.html

	//启动 HTTP 服务器
	log.Println("🚀 Gateway is running on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("❌ Failed to boost gateway :%v", err)
	}

}
