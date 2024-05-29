package service

import (
	saga "github.com/SOA-team-28/common/common/saga/messaging"

	events "database-example/service"
)

type LoginOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewLoginOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*LoginOrchestrator, error) {
	o := &LoginOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (o *LoginOrchestrator) Start(id int) error {
	event := &events.LoginCommand{
		id:   id,
		Type: events.CheckLoginAvailability,
	}

	return o.commandPublisher.Publish(event)
}

func (o *LoginOrchestrator) handle(reply *events.LoginReply) {

}

//ako zatreba

/*

func (o *CreateOrderOrchestrator) nextCommandType(reply events.CreateOrderReplyType) events.CreateOrderCommandType {
	switch reply {
	case events.InventoryUpdated:
		return events.ShipOrder
	case events.InventoryNotUpdated:
		return events.CancelOrder
	case events.InventoryRolledBack:
		return events.CancelOrder
	case events.OrderShippingScheduled:
		return events.ApproveOrder
	case events.OrderShippingNotScheduled:
		return events.RollbackInventory
	default:
		return events.UnknownCommand
	}
}
*/
