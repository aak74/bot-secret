package main

import (
	"log"
)

type repository interface {
	save(*secret) (id int, err error)
}

type updater struct {
	s *secret
}

func (u *updater) save() (id int, err error) {
	log.Println("updater save", u.s.content)
	return 1, nil
}
