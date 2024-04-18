package subscription

// implements BaseNotifier

import (
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct {
	// telegram bot token
	Token string
	// telegram bot
	bot *tgbotapi.BotAPI
	// notification channel
	notiCh chan NotiReq
}

type NotiReq struct {
	ChatId  int64
	Message string
}

func NewTelegramNotifier(token string) *TelegramNotifier {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("NewTelegramNotifier: %v", err)
	}
	notifier := &TelegramNotifier{
		Token:  token,
		bot:    bot,
		notiCh: make(chan NotiReq, 100),
	}
	go notifier.processNotifications()
	return notifier
}

func (t *TelegramNotifier) processNotifications() {
	// 5 messages per second rate limit
	limiter := time.NewTicker(time.Second / 5)
	defer limiter.Stop()

	for range limiter.C {
		select {
		case req := <-t.notiCh:
			t.sendMessage(req.ChatId, req.Message)
		default:
		}
	}
}

func (t *TelegramNotifier) sendMessage(chatId int64, message string) {
	msg := tgbotapi.NewMessage(
		chatId,
		message,
	)
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Fatalf("sendMessage: %v", err)
	}
}

func (t *TelegramNotifier) Notify(chatId, message string) {
	// send telegram message
	chatIdInt, err := strconv.ParseInt(chatId, 10, 64)

	if err != nil {
		log.Fatalf("Notify: %v", err)
	}

	t.notiCh <- NotiReq{
		ChatId:  chatIdInt,
		Message: message,
	}
}

func (t *TelegramNotifier) Close() {
	close(t.notiCh)
}
