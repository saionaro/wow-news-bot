package main

import (
	"log"
	"wow-news-bot/bot"
	"wow-news-bot/broadcaster"
	"wow-news-bot/cacher"
	"wow-news-bot/newsdeamon"
	"wow-news-bot/types"
)

var (
	subscribeChannel   = make(chan int64)
	unsubscribeChannel = make(chan int64)
	freshNewsChannel   = make(chan []types.NewsItem)
)

func main() {
	cacher.LoadCache()
	defer cacher.UnloadCache()
	go newsdeamon.Start(freshNewsChannel)
	go bot.ObserveUpdates(subscribeChannel, unsubscribeChannel)
	for {
		select {
		case freshNews := <-freshNewsChannel:
			go broadcaster.Broadcast(freshNews)
		case id := <-subscribeChannel:
			log.Printf("[%v] wants subscribe", id)
		case id := <-unsubscribeChannel:
			log.Printf("[%v] wants unsubscribe", id)
		}
	}
}
