package saga

type LoginCommandType int8

const (
	CheckLoginAvailability LoginCommandType = iota
	UnknownCommand
)

type LoginCommand struct {
	Id   int
	Type LoginCommandType
}

type LoginReplyType int8

const (
	CanLogin LoginReplyType = iota
	CannotLogin
	UnknownReply
)

type LoginReply struct {
	Type LoginReplyType
}
