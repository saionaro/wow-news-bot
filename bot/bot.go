package bot

import (
	"log"
	"os"
	"strconv"
	"wow-news-bot/helpers"

	"gopkg.in/telegram-bot-api.v4"
)

var (
	botInstance *tgbotapi.BotAPI
	channelID   int64 = -1
)

func createBot() {
	var err error
	botInstance, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	helpers.Check(err)
	log.Printf("Authorized on account %s", botInstance.Self.UserName)
}

func setupChannelID() {
	channelIDVar := os.Getenv("CHANNEL_ID")
	if channelIDVar == "" {
		log.Println("There is not output telegram channel")
	} else {
		channelIDVarInt, err := strconv.ParseInt(channelIDVar, 10, 64)
		helpers.Check(err)
		channelID = channelIDVarInt
		log.Printf("Output telegram channel setted as %d", channelIDVarInt)
	}
}

func GetUpdates() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return botInstance.GetUpdatesChan(u)
}

func ObserveUpdates(subscribeChannel, unsubscribeChannel chan int64) {
	createBot()
	setupChannelID()
	updates, err := GetUpdates()
	helpers.Check(err)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Text == "u" {
			unsubscribeChannel <- update.Message.Chat.ID
		} else {
			subscribeChannel <- update.Message.Chat.ID
		}
	}
}

func GetBot() *tgbotapi.BotAPI {
	return botInstance
}

func SendMessage(text string) (bool, error) {
	if channelID != -1 {
		msg := tgbotapi.NewMessage(channelID, text)
		_, err := botInstance.Send(msg)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
