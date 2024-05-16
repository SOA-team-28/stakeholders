package main

import (
	"context"
	"database-example/db"
	"database-example/model"
	user_service "database-example/proto/user"
	"database-example/service"
	"fmt"
	"log"
	"net"

	"database-example/repo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

type Server struct {
	user_service.UnimplementedUserServiceServer
	UserService *service.UserService
}

func NewServer(db *gorm.DB, tokenRepo *repo.TokenVerificatonRepository) *Server {
	userService := service.NewUserService(db, tokenRepo)
	return &Server{
		UserService: userService,
	}
}

func (s *Server) GetUser(ctx context.Context, req *user_service.GetUserRequest) (*user_service.GetUserResponse, error) {
	user, err := s.UserService.FindUser(int(req.GetId()))
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	return &user_service.GetUserResponse{
		User: &user_service.User{
			Id:                int32(user.Id),
			Username:          user.Username,
			Password:          user.Password,
			Role:              user_service.User_UserRole(user.Role),
			IsActive:          user.IsActive,
			Email:             user.Email,
			VerificationToken: user.VerificationToken,
			IsVerified:        user.IsVerified,
		},
	}, nil
}

func (s *Server) UpsertUser(ctx context.Context, req *user_service.UpsertUserRequest) (*user_service.UpsertUserResponse, error) {
	user := &model.User{
		Id:                int(req.User.GetId()),
		Username:          req.User.GetUsername(),
		Password:          req.User.GetPassword(),
		Role:              model.UserRole(req.User.GetRole()),
		IsActive:          req.User.IsActive,
		Email:             req.User.Email,
		VerificationToken: req.User.VerificationToken,
		IsVerified:        req.User.IsVerified,
	}

	err := s.UserService.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert user: %v", err)
	}

	return &user_service.UpsertUserResponse{
		User: &user_service.User{
			Id:                int32(user.Id),
			Username:          user.Username,
			Password:          user.Password,
			Role:              user_service.User_UserRole(user.Role),
			IsActive:          user.IsActive,
			Email:             user.Email,
			VerificationToken: user.VerificationToken,
			IsVerified:        user.IsVerified,
		},
	}, nil
}

func (s *Server) LoginUser(ctx context.Context, req *user_service.LoginUserRequest) (*user_service.LoginUserResponse, error) {
	token, err := s.UserService.Login(req.GetUsername(), req.GetPassword())
	fmt.Print(req.GetPassword())
	fmt.Print(req.GetUsername())
	if err != nil {
		return nil, fmt.Errorf("login unsuccessful: %v", err)
	}

	return &user_service.LoginUserResponse{
		Token: token,
	}, nil
}

func main() {
	database := db.InitDB()
	if database == nil {
		log.Fatal("FAILED TO CONNECT TO DB")
	}
	tokenRepo := repo.NewTokenVerificatinRepository(database)
	server := NewServer(database, tokenRepo)

	grpcServer := grpc.NewServer()
	user_service.RegisterUserServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("gRPC server listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
