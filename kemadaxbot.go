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
func CountDigits(i int) int {
	count := 0
	for i > 0 {
		i = i / 10
		count += 1
	}

	return count
}
func IsPrime(num int) bool {
	if num < 2 {
		return false
	}
	for i := 2; i < num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}
func primeFactors(num int) ([]string, string) {
	factors := make([]string, 0)
	var factorTree string
	offset := ""
	numDigits := CountDigits(num)
	for i := 2; i < num; i++ {

		if IsPrime(i) && num > 1 {
			for num%i == 0 {
				factors = append(factors, fmt.Sprint(i))
				if CountDigits(num) < numDigits {
					numDigits = CountDigits(num)
					offset += "  "
					factorTree += offset + fmt.Sprint(num) + "|" + fmt.Sprint(i) + "\n"
					num = num / i

				} else {
					factorTree += offset + fmt.Sprint(num) + "|" + fmt.Sprint(i) + "\n"
					num = num / i
				}
			}

		}

	}
	factorTree += offset + fmt.Sprint(num) + "|"
	return factors, factorTree
}

/*func generateBigPrime() int {
	min := 100000000000000000
	max := 1000000000000000000
	rand.Seed(time.Now().UnixNano())
	randint := rand.Intn(max-min+1) + min
	for i := randint; i > 2; i-- {

		if IsPrime(i) {
			return i
		}

	}
	return 0

}*/
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
			msg.Text = "Hey Buddy\nAvailable commands are th following:\n/Convert + (positive whole number as parameter, number< 999 999 999 999) Converting number into words. \nPrimeFactorization + (positive whole number which is greater than 2, accepted as parameter)"
		}
		if update.Message.IsCommand() {

			switch update.Message.Command() {
			case "Convert":
				log.Debug("Converting number to text")
				arg := update.Message.CommandArguments()

				num, err := strconv.Atoi(arg)
				if err != nil {
					log.Debug("/Convert command parameter is not number")
					msg.Text = "Wrong parameter, only positive whole number is accepted as parameter"
				} else if num > 999999999999 {
					msg.Text = "Wrong parameter, only number less than 999.999.999.999 is accepted"
				} else {
					if num == 0 {
						msg.Text = "Nulla"
					} else {
						convertedNum := convert(num)
						convertedNumFirst := convertedNum[:1]
						convertedNumRest := convertedNum[1:]
						convertedNumFirst = strings.Title(convertedNumFirst)
						convertedNum = convertedNumFirst + convertedNumRest

						msg.Text = convertedNum
					}
				}
			case "PrimeFactorization":
				log.Debug("Prime factorization request")
				arg := update.Message.CommandArguments()

				num, err := strconv.Atoi(arg)
				if err != nil || num < 2 {
					log.Debug("/PrimeFactorization command parameter is not number or parameter is less than 2")
					msg.Text = "Wrong parameter, only positive whole number which is greater than 2, accepted as parameter"
				} else {
					if IsPrime(num) {
						msg.Text = fmt.Sprint(num) + " is a prime"
					} else {
						factor, factorTree := primeFactors(num)
						factorJoin := strings.Join(factor, ", ")
						msg.Text = "Prime factors: " + factorJoin + "\n" + "Factor tree:" + "\n" + factorTree
					}
				}

			case "Ping":
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
