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
	fmt.Printf("Usao u handle: ")
	fmt.Printf("Reply primljen u handleru: ", reply)
	user, err := handler.UserService.FindUser(reply.Id)

	if err != nil {

		fmt.Printf("Nije nadjen user! ", user)

	}

	if reply.Type == events.CanLogin {
		token, err := handler.UserService.Login(user.Username, user.Password)
		fmt.Printf("Vracen reply da se moze logovati ")
		if err != nil {

			fmt.Printf("Ne moze se ulogovat, kaze odgovor! ")
			fmt.Printf("token! ", token)

		} else {
			fmt.Printf("updejtuj status na true! ", token)
			//reply rekao da moze, login metoda iz servisa bila uspjesna
			handler.UserService.UpdateStatus(user.Id, true)
		}
	} else {
		fmt.Printf("updejtuj status na false! ")
		//reply rekao da se ne moze logovati
		handler.UserService.UpdateStatus(user.Id, false)
	}

}
