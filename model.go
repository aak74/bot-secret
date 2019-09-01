package main

import (
	"log"
)

type Secret struct {
	id      uint64
	userID  int
	content string
}

func save(content string, userID int) (id int, err error) {
	newSecret := &Secret{
		userID:  userID,
		content: content,
	}
	r := &updater{
		s: newSecret,
	}
	r.save()
	store.Put(newSecret)
	log.Println("model save", content)
	return userID, nil
}
