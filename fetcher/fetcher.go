package fetcher

import (
	"net/http"
	"strings"
	"wow-news-bot/cacher"
	"wow-news-bot/helpers"
	"wow-news-bot/types"

	"github.com/PuerkitoBio/goquery"
)

const newsSourceHost string = "https://www.noob-club.ru"

func FetchNews(channel chan []types.NewsItem) {
	res, err := http.Get(newsSourceHost)
	helpers.Check(err)
	defer res.Body.Close()
	var newsList []types.NewsItem
	doc, err := goquery.NewDocumentFromReader(res.Body)
	helpers.Check(err)
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
