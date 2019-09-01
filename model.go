package main

import (
	"log"
)

type secret struct {
	id      string
	userID  int
	content string
}

func save(content string, userID int) (id int, err error) {
	newSecret := &secret{
		userID:  userID,
		content: content,
	}
	r := &updater{
		s: newSecret,
	}
	r.save()
	log.Println("model save", content)
	return userID, nil
}
