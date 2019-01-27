package newsdeamon

import (
	"fmt"
	"regexp"
	"testing"
	"wow-news-bot/newsfactory"
	"wow-news-bot/types"
)

func TestFilter(t *testing.T) {
	var newsList []types.NewsItem
	var title = "Amazing news!"
	var titles = [2]string{title, "Wow news!"}

	for i := 0; i < len(titles); i++ {
		newsList = append(newsList, newsfactory.Create(
			titles[i],
			"some url",
			"",
		))
	}

	filteredList := filter(newsList, func(item types.NewsItem) bool {
		result, _ := regexp.MatchString("(?i)Wow", item.Title)
		if result {
			return false
		}
		return true
	})

	filteredLen := len(filteredList)

	if filteredLen != 1 {
		t.Error("Expected len 1, got ", filteredLen)
	}

	if filteredLen == 1 && filteredList[0].Title != title {
		t.Error(fmt.Sprintf("Expected title %s, got ", title), filteredList[0].Title)
	}

}

func TestFilterWowlessNews(t *testing.T) {
	var newsList []types.NewsItem

	for i := 0; i < len(wowlessKeywords); i++ {
		newsList = append(newsList, newsfactory.Create(
			wowlessKeywords[i],
			"some url",
			"",
		))
	}

	newList := filterWowlessNews(newsList)

	filteredLen := len(newList)

	if filteredLen != 0 {
		t.Error("Expected len 0, got ", filteredLen)
	}
}
