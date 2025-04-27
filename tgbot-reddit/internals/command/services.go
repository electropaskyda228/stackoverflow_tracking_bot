package command

import (
	"errors"
	"log"
	"strconv"
	. "tgbot-reddit/internals/api"
	. "tgbot-reddit/internals/common"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartService struct{}

func (ss StartService) Service(command *Command) {
	msg := tgbotapi.NewMessage(command.ChatId, command.Value)
	command.Bot.Send(msg)
}

type HelpService struct{}

func (hs HelpService) Service(command *Command) {
	msg := tgbotapi.NewMessage(command.ChatId, command.Value)
	command.Bot.Send(msg)
}

type ListService struct{}

func (ls ListService) Service(command *Command) {
	trackings, err := Listing(command.Value)
	if err != nil {
		msg := tgbotapi.NewMessage(command.ChatId, "Ошибка на сервере. Попробуйте еще раз")
		command.Bot.Send(msg)
		return
	}
	var answer string = "Список вопросов:\n"
	for _, question := range trackings {
		answer += question + "\n"
	}
	msg := tgbotapi.NewMessage(command.ChatId, answer)
	command.Bot.Send(msg)
}

type RegistrationService struct{}

func (rs RegistrationService) Service(command *Command) {
	isExist, err := CheckUserExisting(strconv.FormatInt(command.ChatId, 10))
	if err != nil {
		msg := tgbotapi.NewMessage(command.ChatId, "Ошибка на сервере. Попробуйте еще раз")
		command.Bot.Send(msg)
		return
	}
	if isExist {
		msg := tgbotapi.NewMessage(command.ChatId, "Пользователь уже существует")
		command.Bot.Send(msg)
		return
	}
	err = MakeUser(strconv.FormatInt(command.ChatId, 10), command.Value)
	msg := tgbotapi.NewMessage(command.ChatId, "Пользователь добавлен")
	if err != nil {
		msg = tgbotapi.NewMessage(command.ChatId, "Ошибка. Пользователь не добавлен. Попробйте еще раз")
	} else {
		log.Println("Registration", command.Value)
	}
	command.Bot.Send(msg)
}

type TrackService struct{}

func (ts TrackService) Service(command *Command) {
	err := Tracking(strconv.FormatInt(command.ChatId, 10), command.Value)
	msg := tgbotapi.NewMessage(command.ChatId, "Отслеживание добавлено")
	if errors.Is(err, &FailedAddTracking{}) {
		msg = tgbotapi.NewMessage(command.ChatId, "Такой вопрос уже отслеживается")
	} else if err != nil {
		msg = tgbotapi.NewMessage(command.ChatId, "Отслеживание не добавлено. Ошибка на сервере. Попробуйте еще раз")
	} else {
		log.Println("Tracked", command.Value)
	}
	command.Bot.Send(msg)
}

type UntrackService struct{}

func (us UntrackService) Service(command *Command) {
	err := Untracking(strconv.FormatInt(command.ChatId, 10), command.Value)
	msg := tgbotapi.NewMessage(command.ChatId, "Отслеживание снято")
	if errors.Is(err, &FailedUntrack{}) {
		msg = tgbotapi.NewMessage(command.ChatId, "Такой вопрос не отслеживается")
	} else if err != nil {
		log.Println(err.Error())
		msg = tgbotapi.NewMessage(command.ChatId, "Отслеживание не добавлено. Ошибка на сервере. Попробуйте еще раз")
	} else {
		log.Println("Untracked", command.Value)
	}
	command.Bot.Send(msg)
}
