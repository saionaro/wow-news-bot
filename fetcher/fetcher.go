package fetcher

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"wow-news-bot/helpers"
	"wow-news-bot/newsfactory"
	"wow-news-bot/types"

	"github.com/PuerkitoBio/goquery"
)

const newsSourceHost string = "https://www.noob-club.ru"

func FetchImage(url string) ([]byte, error) {
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

func FetchNews() []types.NewsItem {
	res, err := http.Get(newsSourceHost)
	helpers.Check(err)
	defer res.Body.Close()
	var newsList []types.NewsItem
	doc, err := goquery.NewDocumentFromReader(res.Body)
	helpers.Check(err)
	doc.Find(".entry.first").Each(func(i int, article *goquery.Selection) {
		header := article.Find(".entry-header h1 a")
		image := article.Find(".entry-content img")
		item := newsfactory.Create(
			header.Text(),
			newsSourceHost+header.AttrOr("href", ""),
			image.AttrOr("src", ""),
		)
		newsList = append(newsList, item)
	})
	return newsList
}
