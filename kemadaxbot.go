package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {

	log.SetOutput(os.Stdout)
	
}

func main() {

	purl := os.Getenv("PUBLIC_URL")

	webHookURL := tgbotapi.NewWebhook(purl)

	bot, err := tgbotapi.NewBotAPI("2042481517:AAGd1WViLeY8fpNEdmkDF1C0qIjlr1i6p4g")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(webHookURL)
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/")

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		log.Debug("Request for saying Hello")
		_, err := fmt.Fprintf(w, "Hello %s", os.Getenv("USER"))
		if err != nil {
			log.WithError(err).Warn("Hello handler failed while writing response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	http.HandleFunc("/hello", helloHandler)
	go log.Panic(http.ListenAndServe(":8080", nil))

	for update := range updates {
		log.Printf("%+v\n", update)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello Buddy!")
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}

}
