package command

import (
	"strconv"
	"strings"
	. "tgbot-reddit/internals/common"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var commandDescriptions = map[string]string{
	"/start":        GetStartCommand(nil, 0).Description,
	"/help":         "Список команд",
	"/list":         GetListCommand(nil, 0).Description,
	"/registration": GetRegistrationCommand("", nil, 0).Description,
	"/track":        GetTrackCommand("", nil, 0).Description,
	"/untrack":      GetUntrackCommand("", nil, 0).Description,
}

func ParseCommand(row string, bot *tgbotapi.BotAPI, chatId int64, userName string) {
	words := strings.Fields(row)
	var service Servicable
	switch words[0] {
	case "/start":
		service = StartService{}
		service.Service(GetStartCommand(bot, chatId))
	case "/help":
		service = HelpService{}
		service.Service(GetHelpCommand(bot, chatId))
	case "/list":
		service = ListService{}
		service.Service(GetListCommand(bot, chatId))
	case "/registration":
		service = RegistrationService{}
		service.Service(GetRegistrationCommand(userName, bot, chatId))
	case "/track":
		if len(words) <= 1 {
			ErrorEnter(bot, chatId)
			return
		}
		service = TrackService{}
		service.Service(GetTrackCommand(words[1], bot, chatId))
	case "/untrack":
		if len(words) <= 1 {
			ErrorEnter(bot, chatId)
			return
		}
		service = UntrackService{}
		service.Service(GetUntrackCommand(words[1], bot, chatId))
	}
}

func ErrorEnter(bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Неправильный ввод")
	bot.Send(msg)
}

func GetDescriptions() string {
	var answer strings.Builder

	for cmd, desc := range commandDescriptions {
		answer.WriteString(cmd)
		answer.WriteString(": ")
		answer.WriteString(desc)
		answer.WriteString("\n")
	}

	return answer.String()
}

func GetStartCommand(bot *tgbotapi.BotAPI, chatId int64) *Command {
	return &Command{
		Name:        "start",
		Description: "Старотовая команда",
		Value:       "Привет, это бот для отслеживаний вопросов с сайта stackOverFlow",
		ChatId:      chatId,
		Bot:         bot,
	}
}

func GetHelpCommand(bot *tgbotapi.BotAPI, chatId int64) *Command {
	return &Command{
		Name:        "help",
		Description: "Список команд",
		ChatId:      chatId,
		Bot:         bot,
		Value:       GetDescriptions(),
	}
}

func GetListCommand(bot *tgbotapi.BotAPI, chatId int64) *Command {
	return &Command{
		Name:        "list",
		Description: "Список отслеживаемых вопросов",
		Value:       strconv.FormatUint(uint64(chatId), 10),
		ChatId:      chatId,
		Bot:         bot,
	}
}

func GetRegistrationCommand(username string, bot *tgbotapi.BotAPI, chatId int64) *Command {
	return &Command{
		Name:        "registration",
		Description: "Регистрация",
		Value:       username,
		ChatId:      chatId,
		Bot:         bot,
	}
}

func GetTrackCommand(url string, bot *tgbotapi.BotAPI, chatId int64) *Command {
	return &Command{
		Name:        "track",
		Description: "Добавить ссылку на отслеживание",
		Value:       url,
		ChatId:      chatId,
		Bot:         bot,
	}
}

func GetUntrackCommand(url string, bot *tgbotapi.BotAPI, chatId int64) *Command {
	return &Command{
		Name:        "untrack",
		Description: "Отключить отслеживание ссылки",
		Value:       url,
		ChatId:      chatId,
		Bot:         bot,
	}
}
