package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

}
func convert(num int) string {

	var egyes = map[int]string{
		1: "egy", 2: "kettő", 3: "három", 4: "négy", 5: "öt", 6: "hat", 7: "hét", 8: "nyolc", 9: "kilenc", 10: "tíz",
		11: "tizenegy", 12: "tizenkettő", 13: "tizenhárom", 14: "tizennégy", 15: "tizenöt", 16: "tizenhat", 17: "tizenhét", 18: "tizennyolc", 19: "tizenkilenc",
		20: "húsz", 21: "huszonegy", 22: "huszonkettő", 23: "huszonhárom", 24: "huszonégy", 25: "huszonöt", 26: "huszonhat", 27: "huszonhét", 28: "huszonnyolc", 29: "huszonkilnec",
	}
	var tizes = map[int]string{
		1: "", 2: "", 3: "harminc", 4: "negyven", 5: "ötven", 6: "hatvan", 7: "hetven", 8: "nyolcvan", 9: "kilencven",
	}
	if num < 2000 || num > 2000 && num%1000 == 0 {

		if num < 0 {
			return "mínusz " + convert(-num)
		}
		if num < 30 {
			return egyes[num]
		}
		if num < 100 {
			return tizes[num/10] + egyes[int(num%10)]
		}
		if num < 1000 {
			return egyes[num/100] + "száz" + convert(int(num%100))
		}
		if num < 1000000 {
			return convert(num/1000) + "ezer" + convert(int(num%1000))
		}
		if num < 1000000000 {
			return convert(num/1000000) + "millió" + convert(int(num%1000000))
		}

		return convert(num/1000000000) + "milliárd" + convert(int(num%1000000000))

	} else {

		if num < 0 {
			return "mínusz " + convert(-num)
		}
		if num < 30 {
			return egyes[num]
		}
		if num < 100 {
			return tizes[num/10] + egyes[int(num%10)]
		}
		if num < 1000 {
			return egyes[num/100] + "száz-" + convert(int(num%100))
		}
		if num < 1000000 {
			return convert(num/1000) + "ezer-" + convert(int(num%1000))
		}
		if num < 1000000000 {
			return convert(num/1000000) + "millió-" + convert(int(num%1000000))
		}
		return convert(num/1000000000) + "milliárd-" + convert(int(num%1000000000))
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
					if num == 0 {
						msg.Text = "Nulla"
					} else {
						converetedNum := convert(num)
						msg.Text = strings.Title(converetedNum)
					}
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
