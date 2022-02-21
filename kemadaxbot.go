package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	rand.Seed(time.Now().UnixNano())

}

type primePair struct {
	factor    int
	remainder int
}
type factorizedInt []primePair

func (f factorizedInt) factorsWithCommas() string {
	factors := ""
	for i := 0; i < len(f); i++ {
		factors += fmt.Sprint(f[i].factor) + ", "
	}
	return factors[:len(factors)-2]
}
func (f factorizedInt) factorTree() string {
	factorTree := ""
	offset := ""
	numDigits := CountDigits(f[0].remainder)

	for i := 0; i < len(f); i++ {
		if CountDigits(f[i].remainder) < numDigits {
			numDigits = CountDigits(f[i].remainder)
			offset += "  "
			factorTree += offset + fmt.Sprint(f[i].remainder) + "|" + fmt.Sprint(f[i].factor) + "\n"
		} else {
			factorTree += offset + fmt.Sprint(f[i].remainder) + "|" + fmt.Sprint(f[i].factor) + "\n"
		}

	}
	lastRemainder := f[len(f)-1].remainder / f[len(f)-1].factor
	if CountDigits(lastRemainder) < CountDigits(f[len(f)-1].remainder) {
		offset += "  "

	}
	factorTree += offset + fmt.Sprint(lastRemainder) + "|"

	return factorTree
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
func primeFactorization(num int) factorizedInt {
	f := factorizedInt{}
	for i := 2; i < num; i++ {

		if IsPrime(i) && num > 1 {
			for num%i == 0 {
				pair := primePair{i, num}
				f = append(f, pair)

				num = num / i
			}

		}

	}
	return f
}
func generateBigPrime() int {
	min := 100000000
	max := 1000000000
	randint := rand.Intn(max-min+1) + min
	for i := randint; i > 2; i-- {

		if IsPrime(i) {
			return i
		}

	}
	return 0

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

		if num < 30 {
			return egyes[num]
		}
		if num < 100 {
			return tizes[num/10] + egyes[num%10]
		}
		if num < 1000 {
			return egyes[num/100] + "száz" + convert(num%100)
		}
		if num < 1000000 {
			return convert(num/1000) + "ezer" + convert(num%1000)
		}
		if num < 1000000000 {
			return convert(num/1000000) + "millió" + convert(num%1000000)
		}

		return convert(num/1000000000) + "milliárd" + convert(num%1000000000)

	} else {

		if num < 30 {
			return egyes[num]
		}
		if num < 100 {
			return tizes[num/10] + egyes[num%10]
		}
		if num < 1000 {
			return egyes[num/100] + "száz-" + convert(num%100)
		}
		if num < 1000000 {
			return convert(num/1000) + "ezer-" + convert(num%1000)
		}
		if num < 1000000000 {
			return convert(num/1000000) + "millió-" + convert(num%1000000)
		}
		return convert(num/1000000000) + "milliárd-" + convert(num%1000000000)
	}
}

type InputsDeploy struct {
	ChatID string `json:"chatID"`
}
type RequestToGithubDeploy struct {
	Ref    string       `json:"ref"`
	Inputs InputsDeploy `json:"inputs"`
}

func deploy(URL string, PAT string, chatid string) error {
	encoded := base64.StdEncoding.EncodeToString([]byte(PAT))

	client := &http.Client{}

	reqBody := RequestToGithubDeploy{
		Ref:    "main",
		Inputs: InputsDeploy{ChatID: chatid},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WithError(err).Fatal("deploy func marshal JSON failed")
		return fmt.Errorf("deploy func marshal JSON failed: %v", err)
	}

	r, _ := http.NewRequest("POST", URL, bytes.NewBuffer(reqBytes))
	if err != nil {
		log.WithError(err).Fatal("deploy func making new request failed")
		return fmt.Errorf("deploy func making new request failed: %v", err)
	}
	r.Header.Set("Content-type", "application/vnd.github.v3+json")
	r.Header.Set("Authorization", fmt.Sprintf("Basic %s", encoded))

	resp, err := client.Do(r)
	if err != nil {
		log.WithError(err).Fatal("Something went wrong while deploy func  was sending request to GitHub API")
		return fmt.Errorf("Something went wrong while deploy func  was sending request to GitHub API: %v", err)
	}

	defer resp.Body.Close()

	log.Debug(" GitHub API's HTTP response StatusCode" + fmt.Sprint(resp.StatusCode))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil && body == nil {
		log.WithError(err).Fatal("deploy func could not read request body")
		return fmt.Errorf("deploy func could not read request body: %v", err)
	}

	log.Debug(" GitHub API's HTTP response body content" + string(json.RawMessage(body)))
	return nil
}

type InputsReplicaCount struct {
	ChatID       string `json:"chatID"`
	ReplicaCount string `json:"number_of_replicas"`
	CustomUrl    string `json:"customURL"`
}
type RequestToGithubReplicaCount struct {
	Ref    string             `json:"ref"`
	Inputs InputsReplicaCount `json:"inputs"`
}

func setReplicaCount(URL string, PAT string, chatid string, replicaCount string, customUrl string) error {
	encoded := base64.StdEncoding.EncodeToString([]byte(PAT))

	client := &http.Client{}

	reqBody := RequestToGithubReplicaCount{
		Ref:    "main",
		Inputs: InputsReplicaCount{ChatID: chatid, ReplicaCount: replicaCount, CustomUrl: customUrl},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WithError(err).Fatal("setReplicaCount func marshal JSON failed")
		return fmt.Errorf("setReplicaCountfunc marshal JSON failed: %v", err)
	}

	r, err := http.NewRequest("POST", URL, bytes.NewBuffer(reqBytes))
	if err != nil {
		log.WithError(err).Fatal("setReplicaCount func making new request failed")
		return fmt.Errorf("setReplicaCount func making new request failed: %v", err)
	}
	r.Header.Set("Content-type", "application/vnd.github.v3+json")
	r.Header.Set("Authorization", fmt.Sprintf("Basic %s", encoded))

	resp, err := client.Do(r)
	if err != nil {
		log.WithError(err).Fatal("Something went wrong while setReplicaCount func was sending request to GitHub API")
		return fmt.Errorf("Something went wrong while setReplicaCount func was sending request to GitHub API: %v", err)
	}

	defer resp.Body.Close()

	log.Debug(" GitHub API's HTTP response StatusCode" + fmt.Sprint(resp.StatusCode))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil && body == nil {
		log.WithError(err).Fatal("setReplicaCount func could not read request body")
		return fmt.Errorf("setReplicaCount func could not read request body: %v", err)
	}

	log.Debug(" GitHub API's HTTP response body content" + string(json.RawMessage(body)))
	return nil
}

type MessageFromGitHub struct {
	ChatID string `json:"chat_id"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var randomURL = make([]string, 0)

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
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
	pat := os.Getenv("PAT")

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
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		log.Debug("Request for saying Hello")
		_, err := fmt.Fprintf(w, "Hello %s", os.Getenv("USER"))
		if err != nil {
			log.WithError(err).Warn("Hello handler failed while writing response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	responseAPIHandler := func(w http.ResponseWriter, req *http.Request) {
		log.Debug(req.URL.Path[len(req.URL.Path)-10:])

		for i := 0; i < len(randomURL); i++ {
			if randomURL[i] == req.URL.Path[len(req.URL.Path)-10:] {
				update := MessageFromGitHub{}
				log.Debug("Request from GitHub to responseAPI")
				body, err := ioutil.ReadAll(req.Body)
				if err != nil {
					log.WithError(err).Warn("responseAPI could not read request body")
				}
				err = json.Unmarshal(body, &update)
				if err != nil {
					log.WithError(err).Warn("responseAPI could not Unmarshal request JSON")
				}
				chatid, _ := strconv.ParseInt(update.ChatID, 10, 64)
				msg := tgbotapi.NewMessage(chatid, "Deploy to kubernetes cluster completed")
				if _, err := bot.Send(msg); err != nil {
					log.WithError(err).Warn("responseAPI could not send a message to chat")
				}
				break
			}

		}

	}

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/responseAPI/", responseAPIHandler)

	updates := bot.ListenForWebhook("/")

	go func() { log.Panic(http.ListenAndServe(":8080", nil)) }()

	for update := range updates {

		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			if !update.Message.IsCommand() {
				log.Debug("Answering Hey by default")
				msg.Text = "Hey Buddy\nAvailable commands are th following:\n/Convert + (positive whole number as parameter, number< 999 999 999 999) Converting number into words. \n/PrimeFactorization + (positive whole number which is greater than 2, accepted as parameter)"
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
						log.Debug("/Convert command parameter is greater than 999999999999")
						msg.Text = "Wrong parameter, only number less than 999.999.999.999 is accepted"
					} else if num == 0 {
						msg.Text = "Nulla"
					} else if num < 0 {
						log.Debug("/Convert command parameter is less than 0")
						msg.Text = "Wrong parameter, parameter must be greater than 0"

					} else {
						convertedNum := convert(num)
						convertedNumFirst := convertedNum[:1]
						convertedNumRest := convertedNum[1:]
						convertedNumFirst = strings.Title(convertedNumFirst)
						convertedNum = convertedNumFirst + convertedNumRest

						msg.Text = convertedNum
					}

				case "PrimeFactorization":
					log.Debug("Prime factorization request")
					arg := update.Message.CommandArguments()

					num, err := strconv.Atoi(arg)
					if err != nil {
						log.Debug("/PrimeFactorization command parameter is not number")
						msg.Text = "Wrong parameter, only positive whole number which is greater than 2, accepted as parameter"

					} else if num < 2 {
						log.Debug("/PrimeFactorization command parameter is less than 2")
						msg.Text = "Wrong parameter, parameter must be greater than 2"

					} else if IsPrime(num) {
						log.Debug("/PrimeFactorization command parameter is a prime")
						msg.Text = fmt.Sprint(num) + " is a prime"

					} else {
						result := primeFactorization(num)
						msg.Text = "Prime factors: " + result.factorsWithCommas() + "\n" + "Factor tree:" + "\n" + result.factorTree()
					}

				case "GenerateBigPrime":
					log.Debug("GenerateBigPrime request")
					msg.Text = fmt.Sprint(generateBigPrime())

				case "Deploy":
					err := deploy("https://api.github.com/repos/bproforigoss/kemadaxbot/actions/workflows/chatbot_deploy.yaml/dispatches", pat, fmt.Sprint(update.Message.Chat.ID))
					if err != nil {
						log.Debug("/Deploy failed sending error to chat")
						msg.Text = fmt.Sprint(err)
					}
					msg.Text = "Deployment has been started"

				case "Deploy_debug":
					err := deploy("https://api.github.com/repos/bproforigoss/kemadaxbot/actions/workflows/chatbot_deploy_debug.yaml/dispatches", pat, fmt.Sprint(update.Message.Chat.ID))
					if err != nil {
						log.Debug("/Deploy_debug failed, sending error to chat")
						msg.Text = fmt.Sprint(err)
					}
					msg.Text = "Deployment has been started"

				case "SetReplicaCount":
					arg := update.Message.CommandArguments()
					num, err := strconv.Atoi(arg)
					if err != nil {
						log.Debug("/SetReplicaCount command parameter is not positive whole number")
						msg.Text = "Wrong parameter, parameter is not number"

					} else if num > 50 {
						log.Debug("/SetReplicaCount command parameter is greater than 50")
						msg.Text = "Wrong parameter, parameter is greater 50"

					} else if num < 1 {
						log.Debug("/SetReplicaCount command parameter is less than 50")
						msg.Text = "Wrong parameter, parameter is less than 1"

					} else {
						url := RandStringBytes(10)
						randomURL = append(randomURL, url)
						log.Debug(url)
						err := setReplicaCount("https://api.github.com/repos/bproforigoss/kemadaxbot/actions/workflows/chatbot_set_replica.yaml/dispatches", pat, fmt.Sprint(update.Message.Chat.ID), fmt.Sprint(num), url)
						if err != nil {
							log.Debug("/SetReplicaCount failed sending error to chat")
							msg.Text = fmt.Sprint(err)
						}
						msg.Text = "Deployment has been started"
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
		} else {
			continue
		}

	}

}
