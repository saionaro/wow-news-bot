package cacher

import (
	"crypto/md5"
	"encoding/hex"
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
	hasher   = md5.New()
	db       *bolt.DB
	useCache = true
)

func UnloadCache() {
	if db == nil {
		return
	}
	db.Close()
}

func checkCacheStatus() bool {
	disableCache := os.Getenv("DISABLE_CACHE")
	if disableCache != "" {
		useCache = false
		return false
	}
	return true
}

func LoadCache() {
	if !checkCacheStatus() {
		return
	}
	var err error
	db, err = bolt.Open(cacheDBPath, 0600, nil)
	helpers.Check(err)

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		helpers.Check(err)
		return nil
	})
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
	if db == nil {
		return
	}
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
