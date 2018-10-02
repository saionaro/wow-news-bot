package cacher

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"wow-news-bot/helpers"
	"wow-news-bot/types"

	"github.com/boltdb/bolt"
)

const (
	cacheFilePath       string = "./cache.json"
	cacheDBPath         string = "./cache.db"
	bucket              string = "cache"
	cacheSyncPeriodMins int    = 2
)

var (
	sended = make(map[string]bool)
	hasher = md5.New()
	db     *bolt.DB
)

func UnloadCache() {
	db.Close()
}

func LoadCache() {
	var err error
	db, err = bolt.Open(cacheDBPath, 0600, nil)
	helpers.Check(err)

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		helpers.Check(err)
		return nil
	})
}

func syncCache() {
	fmt.Println("Starting cache sync...")
	cacheFile, err := os.OpenFile(cacheFilePath, os.O_RDWR|os.O_CREATE, 0666)
	helpers.Check(err)
	defer cacheFile.Close()
	jsonString, parseErr := json.Marshal(sended)
	helpers.Check(parseErr)
	n2, writeErr := cacheFile.Write(jsonString)
	helpers.Check(writeErr)
	fmt.Printf("Wrote %d bytes\n", n2)
}

func CheckExistence(hash string) bool {
	var exist bool
	if db == nil {
		return false
	}
	db.View(func(tx *bolt.Tx) error {
		table := tx.Bucket([]byte(bucket))
		value := table.Get([]byte(hash))
		if len(value) > 0 {
			exist = true
		} else {
			exist = false
		}
		return nil
	})
	return exist
}

func MarkSended(item *types.NewsItem) {
	sended[item.Hash] = true
	db.Update(func(tx *bolt.Tx) error {
		table := tx.Bucket([]byte(bucket))
		err := table.Put([]byte(item.Hash), []byte("1"))
		return err
	})
}

func CalcHash(item *types.NewsItem) string {
	hasher.Reset()
	io.WriteString(hasher, item.Href)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
