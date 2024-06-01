package handler

import (
	"database-example/repo"
	"database-example/service"

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

	//OVDJE NISAM SIGURNA STA IDE
	/*
		id, err := primitive.ObjectIDFromHex(command.Order.Id)
		if err != nil {
			return
		}
		order := &domain.Order{Id: id}

		reply := events.CreateOrderReply{Order: command.Order}

		switch command.Type {
		case events.ApproveOrder:
			err := handler.orderService.Approve(order)
			if err != nil {
				return
			}
			reply.Type = events.OrderApproved
		case events.CancelOrder:
			err := handler.orderService.Cancel(order)
			if err != nil {
				return
			}
			reply.Type = events.OrderCancelled
		default:
			reply.Type = events.UnknownReply
		}

		if reply.Type != events.UnknownReply {
			_ = handler.replyPublisher.Publish(reply)
		}
	*/
}
