package fetcher

import (
	"log"
	"net/http"
	"strings"
	"wow-news-bot/cacher"
	"wow-news-bot/types"

	"github.com/PuerkitoBio/goquery"
)

const newsSourceHost string = "https://www.noob-club.ru"

func FetchNews(channel chan []types.NewsItem) {
	res, err := http.Get(newsSourceHost)

	var newsList []types.NewsItem

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".entry.first .entry-header h1 a").Each(func(i int, s *goquery.Selection) {
		item := types.NewsItem{
			Title: strings.TrimSpace(s.Text()),
			Href:  newsSourceHost + s.AttrOr("href", ""),
			Hash:  "",
		}
		item.Hash = cacher.CalcHash(&item)
		newsList = append(newsList, item)
	})

	channel <- newsList
}
