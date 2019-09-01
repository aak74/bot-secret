package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"

	bolt "github.com/coreos/bbolt"
)

type Store struct {
	db *bolt.DB
}

func NewBoltDB() *Store {
	db, err := bolt.Open("./bot.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Store{db}
}

func (s Store) Put(secret *Secret) error {
	log.Println("Put", secret)
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("secrets"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		// b := tx.Bucket([]byte("secrets"))
		log.Println("bucket", b)
		id, _ := b.NextSequence()
		secret.id = id

		// Marshal user data into bytes.
		buf, err := json.Marshal(secret)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(secret.id), buf)
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
