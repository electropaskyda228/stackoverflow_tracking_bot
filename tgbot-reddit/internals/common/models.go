package common

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	ID   int64
	Name string
}

type Command struct {
	Name        string
	Description string
	Value       string
	ChatId      int64
	Bot         *tgbotapi.BotAPI
}

type Servicable interface {
	Service(command *Command)
}
