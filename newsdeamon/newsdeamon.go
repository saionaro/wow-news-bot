package newsdeamon

import (
	"fmt"
	"regexp"
	"time"
	"wow-news-bot/cacher"
	"wow-news-bot/fetcher"
	"wow-news-bot/types"
)

const newsCheckTimeoutMins = 5

var (
	newsChannel          = make(chan []types.NewsItem)
	notificationsChannel chan []types.NewsItem
	wowlessKeywords      = [11]string{"hearthstone", "overwatch", "heroes of the storm", "diablo", "swtor", "hurtstone", "warframe", "mtg arena", "guild wars", "doom", "herostorm"}
)

func filter(vs []types.NewsItem, f func(types.NewsItem) bool) []types.NewsItem {
	vsf := make([]types.NewsItem, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func filterWowlessNews(list []types.NewsItem) []types.NewsItem {
	return filter(list, func(item types.NewsItem) bool {
		for i := 0; i < len(wowlessKeywords); i++ {
			result, _ := regexp.MatchString("(?i)"+wowlessKeywords[i], item.Title)
			if result {
				return false
			}
		}
		return true
	})
}

func filterUnsendedNews(list []types.NewsItem) []types.NewsItem {
	return filter(list, func(item types.NewsItem) bool {
		return !cacher.CheckExistence(item.Hash)
	})
}

func checkNews() {
	fmt.Println("Starting check news...")
	newsList := filterUnsendedNews(filterWowlessNews(fetcher.FetchNews()))
	if len(newsList) > 0 {
		notificationsChannel <- newsList
	}
}

func Start(freshNewsChannel chan []types.NewsItem) {
	notificationsChannel = freshNewsChannel
	newsCheckTicker := time.NewTicker(time.Duration(newsCheckTimeoutMins) * time.Minute)
	defer newsCheckTicker.Stop()
	for {
		select {
		case <-newsCheckTicker.C:
			go checkNews()
		}
	}
}
