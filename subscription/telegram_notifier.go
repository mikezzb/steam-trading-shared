package subscription

// implements BaseNotifier

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct {
	// telegram bot token
	Token string
	// telegram bot
	bot *tgbotapi.BotAPI
}

func NewTelegramNotifier(token string) *TelegramNotifier {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("NewTelegramNotifier: %v", err)
	}
	return &TelegramNotifier{
		Token: token,
		bot:   bot,
	}
}

func (t *TelegramNotifier) Notify(chatId, message string) {
	// send telegram message
	chatIdInt, err := strconv.ParseInt(chatId, 10, 64)
	if err != nil {
		log.Fatalf("Notify: %v", err)
	}

	msg := tgbotapi.NewMessage(
		chatIdInt,
		message,
	)

	_, err = t.bot.Send(msg)
	if err != nil {
		log.Fatalf("Notify: %v", err)
	}
}
