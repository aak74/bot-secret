package main

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Secret struct {
	id      uint64
	userID  int
	content string
}

func save(msg *tgbotapi.Message) (id int, err error) {
	newSecret := &Secret{
		userID:  msg.From.ID,
		content: msg.CommandArguments(),
	}
	// r := &updater{
	// 	s: newSecret,
	// }
	store.Add(newSecret)
	log.Println("model save", msg.CommandArguments())
	return 1, nil
}

func list(msg *tgbotapi.Message) (result []string, err error) {
	return store.GetList(msg.From.ID)
}
