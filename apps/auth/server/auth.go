package server

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	authv1 "github.com/xuewentao/cheya/api/auth/v1"
	"github.com/xuewentao/cheya/apps/auth/ent"
	"github.com/xuewentao/cheya/apps/auth/ent/user" // 用于 Where 条件
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var jwtSecret = []byte("cheya-super-secret-key-2025")

type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
	client ent.Client
}

func (s *AuthServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	//1.查库
	u, err := s.client.User.Query().
		Where(user.Username(req.Username)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, status.Error(codes.Unauthenticated, "用户不存在")
	}
	//2.校验密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password),[]byte(req.Password))
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "密码错误")
	}
	//3.签发 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      u.ID,
		"username": u.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)
	return &authv1.LoginResponse{
		AccessToken: tokenString,
		ExpiresIn:   86400,
		Username:    u.Username,
	}, nil
}
func (s *AuthServer) Register(ctx context.Context,req *authv1.RegisterRequest) (*authv1.RegisterResponse,error) {
	//1.校验代码一致性
	if req.Password != req.ConfirmPassword {
		return nil,status.Error(codes.InvalidArgument,"两次输入的密码不一致")
	}
	exists,err := s.client.User.Query().
								Where(user.Username(req.Username)).
								Exist(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "数据库错误: %v", err)
	}
	if exists {
		return nil ,status.Error(codes.AlreadyExists,"用户名已经存在")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password),bcrypt.DefaultCost)
	if err != nil {
		return nil,status.Error(codes.Unavailable,"密码加密失败")
	}	
	u, err := s.client.User.Create().
							SetUsername(req.Username).
							SetPassword(string(hashedPassword)).
							SetEmail(req.Email).
							Save(ctx)
	return &authv1.RegisterResponse{
		Code:  		201,
		Message:    "注册成功",
		UserId:		int64(u.ID),
		Username:   u.Username,
		CreatedAt:  u.CreatedAt.Format(time.RFC3339),
	},nil

}

func NewAuthServer(client *ent.Client) *AuthServer {
	return &AuthServer{
		client: *client,
	}
}