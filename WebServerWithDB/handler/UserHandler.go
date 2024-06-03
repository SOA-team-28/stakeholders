package handler

import (
	"database-example/repo"
	"database-example/service"
	"fmt"
	"log"

	saga "database-example/service/saga"
	events "database-example/service/saga/check_login"

	"gorm.io/gorm"
)

type UserHandler struct {
	UserService       *service.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
	tokenChannel      chan string
}

func NewUserHandler(db *gorm.DB, tokenRepo *repo.TokenVerificatonRepository, publisher saga.Publisher, subscriber saga.Subscriber, tokenChannel chan string) (*UserHandler, error) {
	userService := service.NewUserService(db, tokenRepo)
	u := &UserHandler{
		UserService:       userService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
		tokenChannel:      tokenChannel,
	}
	log.Println("subsrciber u handleru:", u.commandSubscriber)
	err := u.commandSubscriber.Subscribe(u.Handle)
	if err != nil {
		log.Println("Error subscribing to commands:", err)
		return nil, err
	}
	return u, nil

}
func (handler *UserHandler) Handle(reply *events.LoginReply) {
	fmt.Println("Usao u handle:")
	fmt.Printf("Reply primljen u handleru: %+v\n", reply)

	user, err := handler.UserService.FindUser(reply.Id)
	if err != nil {
		fmt.Println("Nije nadjen user!")
		handler.tokenChannel <- ""
		return
	}

	if reply.Type == events.CanLogin {
		token, err := handler.UserService.Login(user.Username, user.Password)
		if err != nil {
			fmt.Println("Ne moze se ulogovati, kaze odgovor!")
			handler.tokenChannel <- ""
			return
		}
		fmt.Println("Vracen reply da se moze logovati")
		handler.tokenChannel <- token
		return
	}

	handler.tokenChannel <- ""
}
