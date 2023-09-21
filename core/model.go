package core

type Sender struct {
	Id       int64
	UserName string
	History  []*Message
}

type Message struct {
	IsSend bool
	Date   int64
	Msg    string
}
