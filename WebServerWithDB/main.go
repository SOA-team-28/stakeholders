package main

import (
	"context"
	"database-example/db"
	"database-example/handler"
	"database-example/model"
	user_service "database-example/proto/user"
	"database-example/repo"
	"database-example/service"
	"sync"
	"time"

	"database-example/service/saga/nats"
	"fmt"
	"log"
	"net"

	saga "database-example/service/saga"
	events "database-example/service/saga/check_login"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

var (
	tokenChannel chan string
)

type Server struct {
	user_service.UnimplementedUserServiceServer
	UserService      *service.UserService
	commandPublisher *nats.Publisher
	replySubscriber  *nats.Subscriber
	UserHandler      *handler.UserHandler
	TokenService     *service.VerificationTokenService
}

var (
	globalToken string
	mu          sync.Mutex
)

func NewServer(db *gorm.DB, tokenRepo *repo.TokenVerificatonRepository, commandPublisher *nats.Publisher, replySubscriber *nats.Subscriber) (*Server, error) {
	userService := service.NewUserService(db, tokenRepo)
	tokenChannel = make(chan string)
	userHandler, err := handler.NewUserHandler(db, tokenRepo, commandPublisher, replySubscriber, tokenChannel)
	if err != nil {
		return nil, err
	}
	tokenService := service.NewVerificationTokenService(db)
	return &Server{
		UserService:      userService,
		commandPublisher: commandPublisher,
		replySubscriber:  replySubscriber,
		UserHandler:      userHandler, // Dodano
		TokenService:     tokenService,
	}, nil
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
	user, err := s.UserService.FindUserByUseranem(req.Username)
	command := &events.LoginCommand{
		Id:   user.Id,
		Type: events.CheckLoginAvailability,
	}

	err = s.commandPublisher.Publish(command)
	if err != nil {
		return nil, fmt.Errorf("failed to publish login command: %v", err)
	}

	time.Sleep(3 * time.Second)
	updatedUser, err := s.UserService.FindUser(user.Id)
	if updatedUser.CanLogin == true {
		token, err := s.TokenService.FindVerificationTokenByUser(updatedUser.Id)
		if err != nil {
			fmt.Println("cannot find token")
			return &user_service.LoginUserResponse{
				Token: token.TokenData,
			}, nil
		}
		return &user_service.LoginUserResponse{
			Token: token.TokenData,
		}, nil

	} else {
		return &user_service.LoginUserResponse{
			Token: "user cannot login bcs he has too many reports",
		}, nil
	}

}

func main() {
	host := "localhost"
	port := "4222"
	user := "user"
	password := "password"
	commandSubject := "LoginCommand"
	//replySubject := "LoginReply"
	queueGroup := "user-service"

	commandPublisher, err := nats.NewNATSPublisher(host, port, user, password, commandSubject)
	if err != nil {
		panic(err)
	}

	replySubscriber, err := nats.NewNATSSubscriber(host, port, user, password, "LoginReply", queueGroup)
	if err != nil {
		panic(err)
	}
	database := db.InitDB()
	if database == nil {
		log.Fatal("FAILED TO CONNECT TO DB")
	}

	// Eksplicitno konvertovanje tipova
	commandPublisherConverted := commandPublisher.(*nats.Publisher)
	replySubscriberConverted := replySubscriber.(*nats.Subscriber)
	tokenRepo := repo.NewTokenVerificatinRepository(database)

	server, err := NewServer(database, tokenRepo, commandPublisherConverted, replySubscriberConverted)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	loginOrchestrator := server.initLoginOrchestrator(commandPublisher, replySubscriber)
	if loginOrchestrator == nil {
		log.Fatal("failed to create login orchestrator")
	}

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

func (server *Server) initLoginOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) *service.LoginOrchestrator {
	orchestrator, err := service.NewLoginOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}
