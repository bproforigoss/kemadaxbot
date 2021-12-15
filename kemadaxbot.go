package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"

	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

}
func convert(num int) string {
	var egyes = map[int]string{
		1:  "egy",
		2:  "kettő",
		3:  "három",
		4:  "négy",
		5:  "öt",
		6:  "hat",
		7:  "hét",
		8:  "nyolc",
		9:  "kilenc",
		10: "tíz",
	}
	value, ok := egyes[num]
	if !ok {
		return "This value is not currently found in our database to convert"
	}
	return value
}

func TestConversionUpToTen(t *testing.T) {
	conversionResult := convert(6)
	if conversionResult != "hat" {
		t.Errorf("Conversion was incorrect, got: %s, want: %s ", conversionResult, "hat")
	}
}

func main() {

	if len(os.Args) > 1 && os.Args[1] == "-v" {
		log.SetLevel(log.DebugLevel)
		log.Debug("Set loglevel to Debug")
	} else {
		log.SetLevel(log.WarnLevel)
	}
	purl := os.Getenv("PUBLIC_URL")
	token := os.Getenv("API_TOKEN")

	webHookURL := tgbotapi.NewWebhook(purl)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {

		log.WithError(err).Fatal("Something wrong with telegram token")

	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(webHookURL)
	if err != nil {
		log.WithError(err).Fatal("Something wrong with webhook")
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
	go func() { log.Panic(http.ListenAndServe(":8080", nil)) }()

	for update := range updates {

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if !update.Message.IsCommand() {
			log.Debug("Answering Hey by default")
			msg.Text = "Hey Buddy"
		}
		if update.Message.IsCommand() {

			switch update.Message.Command() {
			case "convert":
				log.Debug("Converting number to text")
				arg := update.Message.CommandArguments()

				num, err := strconv.Atoi(arg)
				if err != nil {
					log.Debug("/convert command parameter is not number")
					msg.Text = "Wrong parameter, only numbers as parameters are excepted"
				} else {
					convertedNum := convert(num)
					msg.Text = convertedNum
				}

			case "ping":
				log.Debug("Responding pong, to /ping command")
				msg.Text = "pong"
			default:
				log.Debug("Response to unknown command")
				msg.Text = "I don't know that command"
			}
		}

		if _, err := bot.Send(msg); err != nil {
			log.Debug("Bot sending response")
			log.Panic(err)
		}

	}

}
