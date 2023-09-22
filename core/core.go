package core

import (
	"github.com/uerax/goconf"
)

type Server interface {
	Start()
	SendMsg(id int64, msg string)
}

func NewServer() Server {
	mode, err := goconf.VarInt("mode")
	if err != nil {
		mode = 1
	}
	switch mode {
	case 1:
		return NewTelegram()
	}
	return nil
}