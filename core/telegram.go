package core

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/uerax/goconf"
)

type Telegram struct {
	token   string
	chatId  int64
	replyId int64
	cmd     string
	bot     *tgbot.BotAPI
	senders map[int64]*Sender
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
		token:   token,
		chatId:  chatId,
		bot:     bot,
		replyId: chatId,
		senders: make(map[int64]*Sender),
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
				if update.Message.From.ID != t.chatId {
					go t.SendMsg(update.Message.From.ID, "Do not call the command, send the message directly.")
					continue
				}
				// 命令列表
				switch update.Message.Command() {
				// 历史聊天记录
				case "history":
				// 列表
				case "list":
					msg := fmt.Sprintf("当前回复对象: `%d`\n联系人列表:", t.replyId)
					for id, m := range t.senders {
						msg += fmt.Sprintf("\n%s(`%d`)", m.UserName, id)
					}
					go t.SendMsg(t.chatId, msg)
				// 设置回复Id
				case "to":
					replyId := update.Message.Text[6:]
					id, err := strconv.ParseInt(strings.TrimSpace(replyId), 0, 64)
					if err != nil {
						go t.SendMsg(t.chatId, err.Error())
					} else {
						t.replyId = id
						go t.SendMsg(t.chatId, fmt.Sprintln("当前回复对象: ", t.replyId))
					}
				}
			} else {
				if update.Message.From.ID != t.chatId {
					if _, ok := t.senders[update.Message.From.ID]; !ok {
						t.senders[update.Message.From.ID] = &Sender{
							Id:       update.Message.From.ID,
							UserName: update.Message.From.FirstName + update.Message.From.LastName,
							History:  make([]*Message, 0),
						}
					}
					t.senders[update.Message.From.ID].History = append(t.senders[update.Message.From.ID].History, &Message{
						Date:   time.Now().Unix(),
						Msg:    update.Message.Text,
						IsSend: true,
					})
					// 自动将其设置成待回复的ID
					t.replyId = update.Message.From.ID
					// 拼接发送消息
					mc := tgbot.NewForward(t.chatId, update.Message.From.ID, update.Message.MessageID)

					go t.bot.Send(mc)
				} else {
					if _, ok := t.senders[t.replyId]; !ok {
						t.senders[update.Message.From.ID] = &Sender{
							Id:       update.Message.From.ID,
							UserName: update.Message.From.FirstName + update.Message.From.LastName,
							History:  make([]*Message, 0),
						}
					}

					if update.Message.Text != "" {
						t.senders[t.replyId].History = append(t.senders[t.replyId].History, &Message{
							Date:   time.Now().Unix(),
							Msg:    update.Message.Text,
							IsSend: false,
						})
					}
					mc := tgbot.NewCopyMessage(t.replyId, update.Message.From.ID, update.Message.MessageID)

					go t.bot.Send(mc)
				}
			}
		}
	}
}

func (t *Telegram) SendMsg(id int64, msg string) {
	mc := tgbot.NewMessage(id, msg)
	mc.ParseMode = "Markdown"
	mc.DisableWebPagePreview = false
	t.bot.Send(mc)
}
