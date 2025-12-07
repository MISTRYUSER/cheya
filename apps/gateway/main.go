package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	_"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket" // ✅ 新增
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/xuewentao/cheya/api/auth/v1"
	vehiclev1 "github.com/xuewentao/cheya/api/vehicle/v1"
)


//jwt-Secret
var jwtSecret = []byte("cheya-super-secret-key-2025")

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
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("❌ Failed to connect gRPC server %v", err)
	}
	defer conn.Close()
	//创建 grpc client 存根
	vehicleClient := vehiclev1.NewVehicleServiceClient(conn)
	log.Println("✅ Connected to Vehicle Service(gRPC)")

	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()

	//1.Redis 订阅
	go func() {
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

	// CORS 中间件 - 允许前端跨域访问
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

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

	//GET /api/v1/vehicles
	r.GET("/api/v1/vehicles", func(c *gin.Context) {
		// 从查询参数获取分页信息，设置默认值
		page := int32(1)
		pageSize := int32(100)

		// 解析 page 参数
		if pageParam := c.Query("page"); pageParam != "" {
			if p, err := strconv.ParseInt(pageParam, 10, 32); err == nil && p > 0 {
				page = int32(p)
			}
		}

		// 解析 pageSize 参数
		if pageSizeParam := c.Query("pageSize"); pageSizeParam != "" {
			if ps, err := strconv.ParseInt(pageSizeParam, 10, 32); err == nil && ps > 0 {
				pageSize = int32(ps)
			}
		}

		//构造 grpc 请求
		req := &vehiclev1.ListVehiclesRequest{
			Page:     page,
			PageSize: pageSize,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		//2.调用 grpc
		resp, err := vehicleClient.ListVehicles(ctx, req)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		//返回 json
		c.JSON(200, gin.H{
			"code": 200,
			"data": gin.H{
				"items": resp.Vehicles,
				"total": resp.TotalCount,
			},
		})
	})

	// 车辆控制接口
	r.POST("/api/v1/vehicles/:vin/control", func(c *gin.Context) {
		vin := c.Param("vin")

		var body struct {
			Action string `json:"action"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		// 验证动作类型
		if body.Action != "STOP" && body.Action != "START" {
			c.JSON(400, gin.H{"error": "Invalid action. Must be STOP or START"})
			return
		}

		// 构造命令并发布到 Redis
		cmd := body.Action + ":" + vin
		err := rdb.Publish(context.Background(), "vehicle:commands", cmd).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		log.Printf("📢 Command sent: %s for vehicle %s", body.Action, vin)
		c.JSON(200, gin.H{
			"code":    200,
			"message": "Command sent successfully",
			"data": gin.H{
				"vin":    vin,
				"action": body.Action,
			},
		})
	})

	//连接 auth service login
	authConn, _ := grpc.NewClient("localhost:50054", grpc.WithTransportCredentials(insecure.NewCredentials()))
	authClient := authv1.NewAuthServiceClient(authConn)
	//Login
	r.POST("/api/v1/auth/login", func(c *gin.Context) {
		var req authv1.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid body"})
			return
		}

		resp, err := authClient.Login(context.Background(), &req)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"access_token": resp.AccessToken,
			"expires_in":   resp.ExpiresIn,
			"username":     resp.Username,
		})
	})
	r.POST("/api/v1/auth/register",func(c *gin.Context) {
		var req authv1.RegisterRequest
		if err := c.ShouldBindJSON(&req);err != nil {
			c.JSON(400, gin.H{"error": "Invalid body"})
			return
		}

		resp,err := authClient.Register(context.Background(),&req) 
		if err != nil {
			c.JSON(401,gin.H{"error":err.Error()})
			return
		}
		c.JSON(201,gin.H {
			"code" :    resp.Code,
			"message": 	resp.Message,
			"user_id":  resp.UserId,
			"username": resp.Username,
			"created_at":resp.CreatedAt,
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
