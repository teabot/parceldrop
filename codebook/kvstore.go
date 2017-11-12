package codebook

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

const (
	bucketName = "access_codes"
)

var db *bolt.DB
var open bool

func OpenStore() error {
	var err error
	dbfile := "codebook.db"
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(dbfile, 0600, config)
	if err != nil {
		log.Fatal(err)
	}
	err2 := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("CODEBOOK: create bucket: %s", err)
		}
		return err
	})
	if err2 != nil {
		log.Fatal(err2)
	}
	open = true
	return nil
}

func CloseStore() {
	open = false
	db.Close()
}

func (p *AccessCode) save() error {
	if !open {
		return fmt.Errorf("CODEBOOK: db must be opened before saving")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		codes, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("CODEBOOK: create bucket: %s", err)
		}
		enc, err := p.encode()
		if err != nil {
			return fmt.Errorf("CODEBOOK: could not encode AccessCode %s: %s", p.Digits, err)
		}
		err = codes.Put([]byte(p.Digits), enc)
		return err
	})
	return err
}

func (p *AccessCode) encode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func decode(data []byte) (*AccessCode, error) {
	var p *AccessCode
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetAccessCode(digits string) (*AccessCode, error) {
	if !open {
		return nil, fmt.Errorf("CODEBOOK: db must be opened before saving")
	}
	var p *AccessCode
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(bucketName))
		k := []byte(digits)
		record := b.Get(k)
		if len(record) <= 0 {
			return nil
		}
		p, err = decode(record)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("CODEBOOK: Could not get AccessCode %s", digits)
		return nil, err
	}
	return p, nil
}

func List() {
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucketName)).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("CODEBOOK: key=%s, value=%s\n", k, v)
		}
		return nil
	})
}
