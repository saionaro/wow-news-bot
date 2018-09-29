package main

import (
	"log"
	"wow-news-bot/bot"
	"wow-news-bot/cacher"
	"wow-news-bot/newsdeamon"
	"wow-news-bot/types"
)

var (
	subscribeChannel   = make(chan int64)
	unsubscribeChannel = make(chan int64)
	freshNewsChannel   = make(chan []types.NewsItem)
)

func sendNews(freshNews []types.NewsItem) {
	log.Println("Starting news broadcast...")
	for i := len(freshNews) - 1; i >= 0; i-- {
		_, err := bot.SendMessage(freshNews[i].Title + "\n" + freshNews[i].Href)
		if err == nil {
			cacher.MarkSended(&freshNews[i])
		}
	}
}

func main() {
	cacher.LoadCache()
	go cacher.StartSyncDeamon()
	go newsdeamon.Start(freshNewsChannel)
	go bot.ObserveUpdates(subscribeChannel, unsubscribeChannel)
	for {
		select {
		case freshNews := <-freshNewsChannel:
			go sendNews(freshNews)
		case id := <-subscribeChannel:
			log.Printf("[%v] wants subscribe", id)
		case id := <-unsubscribeChannel:
			log.Printf("[%v] wants unsubscribe", id)
		}
	}
}
