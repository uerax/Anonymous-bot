package core

import (
	"github.com/uerax/goconf"
)

type Server interface {
	Start()
	SendMsg(id int64, msg string)
}

func NewServer() Server {
	mode := goconf.VarIntOrDefault(1, "mode")
	switch mode {
	case 1:
		return NewTelegram()
	default:
		return NewTelegram()
	}
}