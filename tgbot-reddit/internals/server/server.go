package server

import (
	"log"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func checkConstants() {
	if os.Getenv("BOT_PORT") == "" {
		log.Println("Os's constants have not been found")
		os.Exit(1)
	}
}

func MakeServer(bot *tgbotapi.BotAPI) {
	checkConstants()
	http.HandleFunc("/bot/send", func(w http.ResponseWriter, req *http.Request) {
		body := "Method is not allowed"
		status := http.StatusBadRequest
		w.Header().Set("Content-Type", "text/plain")
		if req.Method != "POST" {
			w.WriteHeader(status)
			w.Write([]byte(body))
			return
		}

		err := req.ParseForm()
		if err != nil {
			body = "Failed to send message"
			w.WriteHeader(status)
			w.Write([]byte(body))
			return
		}

		params := req.Form
		chatIdString := params.Get("chat_id")
		questionId := params.Get("question")

		chatId, err := strconv.ParseInt(chatIdString, 10, 64)
		if err != nil {
			body = "Failed to send message"
			w.WriteHeader(status)
			w.Write([]byte(body))
			return
		}

		body = "Success of sending message"
		status = http.StatusOK
		w.WriteHeader(status)
		w.Write([]byte(body))

		log.Println("New answer on question", questionId)

		msg := tgbotapi.NewMessage(chatId, "Появились обновления в вопросе с номером "+questionId)
		bot.Send(msg)
	})

	http.ListenAndServe(os.Getenv("BOT_PORT"), nil)
}
