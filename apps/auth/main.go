package main

import (
	"context"
	"log"
	"net"

	authv1 "github.com/xuewentao/cheya/api/auth/v1"
	"github.com/xuewentao/cheya/apps/auth/ent"
	"github.com/xuewentao/cheya/apps/auth/server"
	"google.golang.org/grpc"
	_ "github.com/lib/pq"
)
func main() {
	//1.链接数据库
	dns := "host=localhost port=5432 user=wentao_xue dbname=cheya password=Woe89132 sslmode=disable"
	client, err := ent.Open("postgres",dns)
	if err != nil {
		log.Fatalf("❌ failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	//2.自动迁移
	//在 db中自动创建 user
	log.Println("📦 Migrating database schema...")
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("❌ failed creating schema resources: %v", err)
	}
	log.Println("✅ Schema migrated successfully!")

	//3.启动 grpc
	lis, err := net.Listen("tcp",":50054")
	if err != nil{
		log.Fatalf("❌ failed to listen : %v", err)
	}

	s := grpc.NewServer()
	//注入 client 到 server
	authv1.RegisterAuthServiceServer(s,server.NewAuthServer(client))
	log.Printf("🚀 Auth Service is running on :50054")
	//4.启动服务
	if err := s.Serve(lis);err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
