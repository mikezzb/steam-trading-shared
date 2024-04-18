package subscription_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	subs "github.com/mikezzb/steam-trading-shared/subscription"
)

func TestTelegram(t *testing.T) {
	// load secrets from json
	secretsPath := "./secrets.json"
	secretsFile, err := os.ReadFile(secretsPath)
	if err != nil {
		log.Fatalf("TestTelegram: %v", err)
	}
	var secrets map[string]string
	// unmarshal json
	json.Unmarshal(secretsFile, &secrets)

	notifier := subs.NewNotifier(&subs.NotifierConfig{
		TelegramToken: secrets["telegramToken"],
	})

	t.Run("TestTelegramMessage", func(t *testing.T) {
		notifier.Notify("telegram", secrets["telegramTestChatId"], "Hello Test!")

		// wait for message to be sent
		time.Sleep(time.Second)
	})
}
