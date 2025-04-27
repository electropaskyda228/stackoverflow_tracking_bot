package main

import (
	"log"
	"os"

	. "tgbot-reddit/internals/command"
	. "tgbot-reddit/internals/server"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func checkConstants() {
	if os.Getenv("TGBOT_REDDIT") == "" {
		log.Println("Os's constants have not been found")
		os.Exit(1)
	}
	if os.Getenv("LOCAL_HOST") == "" || os.Getenv("STACKOVERFLOW_SERVER_PORT") == "" {
		log.Println("Os's constants have not been found")
		os.Exit(1)
	}
}

func main() {
	checkConstants()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGBOT_REDDIT"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	go MakeServer(bot)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			go ParseCommand(update.Message.Text, bot, update.Message.Chat.ID, update.Message.Chat.UserName)
		}
	}
}
