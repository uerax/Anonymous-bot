package core

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/uerax/Anonymous-bot/core"
	"github.com/uerax/goconf"
)

type Telegram struct {
	token   string
	chatId  int
	replyId int
	cmd     string
	bot     *tgbot.BotAPI
	senders map[int]*core.Sender
}

func NewTelegram() *Telegram {
	token, err := goconf.VarString("telegram", "token")
	if err != nil {
		log.Panic(err)
	}

	chatId, err := goconf.VarInt("telegram", "chatId")
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	tg := &Telegram{
		token:   token,
		chatId:  chatId,
		bot:     bot,
		replyId: chatId,
		senders: make(map[int]*core.Sender),
	}

	return tg
}

func (t *Telegram) Start() {
	bot, err := tgbot.NewBotAPI(t.token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if update.Message.IsCommand() {
				// 非持有者不允许调用命令
				if update.Message.MessageID != t.chatId {
					go t.SendMsg(update.Message.MessageID, "Do not call the command, send the message directly.")
					continue
				}
				// 命令列表
				switch update.Message.Command() {
				case "list":
					msg := "List:"
					for id, m := range t.senders {
						msg += fmt.Sprintf("\n%s(`%d`)", m.UserName, id)
					}
					go t.SendMsg(t.chatId, msg)
				// 设置回复Id
				case "reply":
					id, err := strconv.Atoi(update.Message.Text)
					if err != nil {
						go t.SendMsg(t.chatId, err.Error())
					} else {
						t.replyId = id
						go t.SendMsg(t.chatId, fmt.Sprintln("当前回复对象: ", t.replyId))
					}
				}
			} else {
				if update.Message.MessageID != t.chatId {
					if _, ok := t.senders[update.Message.MessageID]; !ok {
						t.senders[update.Message.MessageID] = new(core.Sender)
					}
					t.senders[update.Message.MessageID].History = append(t.senders[update.Message.MessageID].History, &core.Message{
						Date: time.Now().Unix(),
						Msg: update.Message.Text,
						IsSend: true,
					})
					// 自动将其设置成待回复的ID
					t.replyId = update.Message.MessageID
					// 拼接发送消息
					msg := fmt.Sprintf("*%s(%d) :*\n%s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)
					go t.SendMsg(t.chatId, msg)
				} else {
					if _, ok := t.senders[t.replyId]; !ok {
						t.senders[t.replyId] = new(core.Sender)
					}
					t.senders[t.replyId].History = append(t.senders[t.replyId].History, &core.Message{
						Date: time.Now().Unix(),
						Msg: update.Message.Text,
						IsSend: false,
					})
					go t.SendMsg(t.replyId, update.Message.Text)
				}
			}
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbot.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func (t *Telegram) SendMsg(id int, msg string) {
	mc := tgbot.NewMessage(int64(id), msg)
	mc.ParseMode = "Markdown"
	mc.DisableWebPagePreview = false
}
