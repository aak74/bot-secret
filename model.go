package main

import (
	"log"
)

type Secret struct {
	id      string
	userId  int
	content string
}

func (s *Secret) save() (id int, err error) {
	log.Println("add", s.content)
	return 1, nil
}
