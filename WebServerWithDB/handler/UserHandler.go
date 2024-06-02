package handler

import (
	"database-example/repo"
	"database-example/service"
	"fmt"

	saga "database-example/service/saga"
	events "database-example/service/saga/check_login"

	"gorm.io/gorm"
)

type UserHandler struct {
	UserService       *service.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewUserHandler(db *gorm.DB, tokenRepo *repo.TokenVerificatonRepository, publisher saga.Publisher, subscriber saga.Subscriber) (*UserHandler, error) {
	userService := service.NewUserService(db, tokenRepo)
	u := &UserHandler{
		UserService:       userService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := u.commandSubscriber.Subscribe(u.handle)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (handler *UserHandler) handle(command *events.LoginCommand) {
	reply := &events.LoginReply{}

	switch command.Type {
	case events.CheckLoginAvailability:
		token, err := handler.UserService.Login(command.Username, command.Password)
		fmt.Printf("Usao u login komandu: ")
		if err != nil {
			reply.Type = events.CannotLogin
			fmt.Printf("Ne moze se ulogovat, kaze odgovor! ")
		} else {
			reply.Type = events.CanLogin
			reply.Token = token // Povratak tokena u slučaju uspešnog login-a
		}
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		err := handler.replyPublisher.Publish(reply)
		if err != nil {
			fmt.Printf("Failed to publish reply: %v\n", err)
		}
	}
}
