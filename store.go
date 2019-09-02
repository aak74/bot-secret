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

func (s Store) GetList(userID int) ([]string, error) {
	s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("secrets"))

		// c := b.Cursor()

		sec := &Secret{}
		// for k, v := c.First(); k != nil; k, v = c.Next() {
		// 	json.Unmarshal(v, &sec)
		// 	fmt.Printf("key=%s, value=%v %v\n", k, sec, v)
		// }

		b.ForEach(func(k, v []byte) error {
			json.Unmarshal(v, &sec)
			fmt.Printf("key=%s, value=%v %v\n", k, sec, v)
			return nil
		})
		return nil
	})

	return []string{"qaz", "123"}, nil
}

func (s Store) Add(secret *Secret) error {
	log.Println("Put", secret)
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("secrets"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		// b := tx.Bucket([]byte("secrets"))
		log.Println("bucket", b, secret)
		id, _ := b.NextSequence()
		secret.id = id

		// Marshal user data into bytes.
		buf, err := json.Marshal(secret)
		if err != nil {
			return err
		}
		log.Printf("put to bucket %v %v\n", secret, buf)
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
