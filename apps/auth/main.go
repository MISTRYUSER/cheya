package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/golang-jwt/jwt/v5"
	authv1 "github.com/xuewentao/cheya/api/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var jwtSecret = []byte("cheya-super-secret-key-2025")

type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
}

func (c *AuthServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	//暂时模拟校验
	if req.Username != "admin" || req.Password != "123456" {
		return nil, status.Errorf(codes.Unauthenticated, "用户名或密码错误")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  "u-001",
		"username": req.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "生成token失败: %v", err)
	}

	return &authv1.LoginResponse{
		AccessToken: tokenString,
		ExpiresIn:   86400, // 24小时 = 86400秒
		Userme:      req.Username,
	}, nil
}
func main() {
	lis,err := net.Listen("tcp",":50054")
	if err != nil { 
		log.Fatalf("failed to listen : %v",err)
	}	
	s := grpc.NewServer()
	authv1.RegisterAuthServiceServer(s,&AuthServer{})

	log.Println("Auth service is running on:50054")
	s.Serve(lis)
}