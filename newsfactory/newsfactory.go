package newsfactory

import (
	"strings"
	"wow-news-bot/cacher"
	"wow-news-bot/types"
)

func Create(title, herf, imageSrc string) types.NewsItem {
	item := types.NewsItem{
		Title: strings.TrimSpace(title),
		Href:  herf,
		Image: imageSrc,
		Hash:  "",
	}
	item.Hash = cacher.CalcHash(&item)
	return item
}
