package broadcaster

import (
	"log"
	"wow-news-bot/bot"
	"wow-news-bot/cacher"
	"wow-news-bot/fetcher"
	"wow-news-bot/types"
)

func sendMessage(msg *types.Message) (bool, error) {
	if len(msg.Image) > 1 {
		return bot.SendImage(msg)
	}
	return bot.SendMessage(msg)
}

func Broadcast(freshNews []types.NewsItem) {
	log.Println("Starting news broadcast...")
	for i := len(freshNews) - 1; i >= 0; i-- {
		message := types.Message{
			Text:  freshNews[i].Title + "\n" + freshNews[i].Href,
			Image: make([]byte, 1),
		}
		if freshNews[i].Image != "" {
			image, err := fetcher.FetchImage(freshNews[i].Image)
			if err == nil {
				message.Image = image
			}
		}
		_, err := sendMessage(&message)
		if err == nil {
			cacher.MarkSended(&freshNews[i])
		}
	}
	log.Println("News broadcast finished.")
}
