package subscription

import "log"

type BaseNotifier interface {
	Notify(notiId string, message string)
}

type Notifier struct {
	// notiType -> Notifier
	notifiers map[string]BaseNotifier
}

func (n *Notifier) Notify(notiType, notiId string, message string) {
	log.Printf("Notifier.Notify: %s %s %s", notiType, notiId, message)
	notifier, ok := n.notifiers[notiType]
	if !ok {
		return
	}
	notifier.Notify(notiId, message)
}

type NotifierConfig struct {
	TelegramToken string
}

func NewNotifier(config *NotifierConfig) *Notifier {
	notifiers := make(map[string]BaseNotifier)
	notifiers["telegram"] = NewTelegramNotifier(config.TelegramToken)
	return &Notifier{
		notifiers: notifiers,
	}
}
