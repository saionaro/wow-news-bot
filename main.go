package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
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

func sendMessage(msg *types.Message) (bool, error) {
	if len(msg.Image) > 1 {
		return bot.SendImage(msg)
	}
	return bot.SendMessage(msg)
}

func downloadImage(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var data bytes.Buffer
	_, err = io.Copy(&data, res.Body)
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}

func sendNews(freshNews []types.NewsItem) {
	log.Println("Starting news broadcast...")
	for i := len(freshNews) - 1; i >= 0; i-- {
		message := types.Message{
			Text:  freshNews[i].Title + "\n" + freshNews[i].Href,
			Image: make([]byte, 1),
		}
		if freshNews[i].Image != "" {
			image, err := downloadImage(freshNews[i].Image)
			if err == nil {
				message.Image = image
			}
		}
		_, err := sendMessage(&message)
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
