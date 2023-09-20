package core

import (
	"log"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/uerax/goconf"
)

type Telegram struct {
	token  string
	chatId int64
	bot    *tgbot.BotAPI
}

func NewTelegram() *Telegram {
	token, err := goconf.VarString("telegram", "token")
	if err != nil {
		log.Panic(err)
	}

	chatId, err := goconf.VarInt64("telegram", "chatId")
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	tg := &Telegram{
		token:  token,
		chatId: chatId,
		bot:    bot,
	}

	return tg
}

