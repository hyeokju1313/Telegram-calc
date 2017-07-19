package main

import (
	"os"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"net/http"
	"io/ioutil"
)

func main() {
	token := os.Args[1]
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			if update.Message.Command() == "calc" {
				a := update.Message.CommandArguments();
				a = strings.Replace(a, "+", "%2B", len(a))
				a = strings.Replace(a, "/", "%2F", len(a))
				r, err := http.Get("http://api.mathjs.org/v1/?expr=" + a)
				if err != nil {
					log.Panic(err)
				}
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					log.Panic(err)
				}
				m := tgbotapi.NewMessage(update.Message.Chat.ID, string(body))
				m.ReplyToMessageID = update.Message.MessageID
				bot.Send(m)
				log.Printf("[%s] %s", update.Message.Chat.UserName, a)
			}
		}
	}
}
