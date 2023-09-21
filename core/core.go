package core

type Server interface {
	Start()
	SendMsg(id int, msg string)
}