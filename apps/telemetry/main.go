package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/redis/go-redis/v9"

	telemetryv1 "github.com/xuewentao/cheya/api/telemetry/v1"
	"github.com/xuewentao/cheya/apps/telemetry/consumer" // 引入我们刚才写的包
)

type TelemetryServer struct {
	telemetryv1.UnimplementedTelemetryServiceServer
}

func main() {
	//创建上下文用于控制生命周期
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//1.初始化 redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // 开发环境无密码
		DB:       0,  // 使用默认数据库
	})
	defer rdb.Close()
	// 测试 Redis 连接
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	log.Println("✅ Connected to Redis")

	//2.后台启动 kafka 消费者
	brokers := []string{"localhost:9092"}
	topic := "telemetry.raw"

	go consumer.StartTelemetryConsumer(ctx, brokers, topic, rdb)
	//3.启动 grpc server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	telemetryv1.RegisterTelemetryServiceServer(s, &TelemetryServer{})

	go func() {
		log.Println("📡 Telemetry Service is running on :50052")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve :%v", err)
		}
	}()

	//优雅退出
	quit := make(chan os.Signal, 1) //与 signal.Notify搭配使用
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down services...")
	cancel()
	s.GracefulStop() //grace 优雅退出 不要暴力 shut down 等所有的 io 操作完成再退出

}
