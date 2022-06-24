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

	"github.com/bproforigoss/kemadaxbot/chatboterrors"
	"github.com/bproforigoss/kemadaxbot/chatbotstructs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	rand.Seed(time.Now().UnixNano())
	prometheus.MustRegister(respDuration)
	prometheus.MustRegister(reqCounter)
	prometheus.MustRegister(respDurationAvg)
	prometheus.MustRegister(reqFrequencyCounter)
}

var (
	respDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "generatBigPrime_request_duration",
		Help:    "Durations till primeGenerator component responds with prime",
		Buckets: []float64{60, 80, 100, 120, 140, 160, 180, 200, 220, 240, 260},
	})
	reqCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "Command_request_counter",
		Help: "Number of primeGenerator component responds with prime",
	},
		[]string{"command"})
	respDurationAvg = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "generatBigPrime_request_duration_avg",
		Help: "Avarage durations primeGenerator component responds with prime",
	},
		[]string{"command"})
	reqFrequencyCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "generatBigPrime_request_frequency_counter",
		Help: "Number of primeGenerator component responds with prime",
	})

	randomURL = make([]string, 0)
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type factorizedInt []chatbotstructs.PrimePair

func (f factorizedInt) factorsWithCommas() string {
	factors := ""
	for i := 0; i < len(f); i++ {
		factors += fmt.Sprint(f[i].Factor) + "," + " "
	}
	return factors[:len(factors)-2]
}
func (f factorizedInt) factorTree() string {
	factorTree := ""
	offset := ""
	numDigits := CountDigits(f[0].Remainder)

	for i := 0; i < len(f); i++ {
		if CountDigits(f[i].Remainder) < numDigits {
			numDigits = CountDigits(f[i].Remainder)
			offset += "  "
			factorTree += offset + fmt.Sprint(f[i].Remainder) + "|" + fmt.Sprint(f[i].Factor) + "\n"
		} else {
			factorTree += offset + fmt.Sprint(f[i].Remainder) + "|" + fmt.Sprint(f[i].Factor) + "\n"
		}

	}
	lastRemainder := f[len(f)-1].Remainder / f[len(f)-1].Factor
	if CountDigits(lastRemainder) < CountDigits(f[len(f)-1].Remainder) {
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
				pair := chatbotstructs.PrimePair{i, num}
				f = append(f, pair)

				num = num / i
			}

		}

	}
	return f
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
	if (num < 2000) || (num > 2000 && num%1000 == 0) {

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

//done
func deploy(URL string, PAT string, chatid string) error {
	encoded := base64.StdEncoding.EncodeToString([]byte(PAT))

	client := &http.Client{}

	reqBody := chatbotstructs.RequestToGithubDeploy{
		Ref:    "main",
		Inputs: chatbotstructs.InputsDeploy{ChatID: chatid},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("deploy func marshal JSON failed: %v", err)
	}

	r, _ := http.NewRequest("POST", URL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return fmt.Errorf("deploy func making new request failed: %v", err)
	}
	r.Header.Set("Content-type", "application/vnd.github.v3+json")
	r.Header.Set("Authorization", fmt.Sprintf("Basic %s", encoded))

	resp, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("something went wrong while deploy func  was sending request to GitHub API: %v", err)
	}

	defer resp.Body.Close()

	log.Debug("GitHub API's HTTP response StatusCode" + fmt.Sprint(resp.StatusCode))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil && body == nil {
		return fmt.Errorf("deploy func could not read request body: %v", err)
	}

	log.Debug(" GitHub API's HTTP response body content" + string(json.RawMessage(body)))
	return nil
}

//done
func loadRequest(URL string, chatbotURL string, num int, chatid int64) error {
	client := &http.Client{}

	reqBody := chatbotstructs.RequestToLoad{
		Url:    chatbotURL,
		Number: num,
		ChatID: chatid,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("loadRequest func marshal JSON failed: %v", err)
	}

	r, _ := http.NewRequest("POST", URL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return fmt.Errorf("loadRequest func making new request failed: %v", err)
	}
	r.Header.Set("Content-type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("something went wrong while loadRequest func was sending request to loadTestingTool: %v", err)
	}

	defer resp.Body.Close()

	log.Debug("loadTestingTool's HTTP response StatusCode" + fmt.Sprint(resp.StatusCode))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil && body == nil {
		return fmt.Errorf("loadRequest func could not read response body: %v", err)
	}

	log.Debug("loadTestingTool's HTTP response body content" + string(json.RawMessage(body)))
	return nil
}

func setReplicaCount(URL string, PAT string, chatid string, replicaCount string, customUrl string) error {
	encoded := base64.StdEncoding.EncodeToString([]byte(PAT))

	client := &http.Client{}

	reqBody := chatbotstructs.RequestToGithubReplicaCount{
		Ref:    "main",
		Inputs: chatbotstructs.InputsReplicaCount{ChatID: chatid, ReplicaCount: replicaCount, CustomUrl: customUrl},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("setReplicaCountfunc marshal JSON failed: %v", err)
	}

	r, err := http.NewRequest("POST", URL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return fmt.Errorf("setReplicaCount func making new request failed: %v", err)
	}
	r.Header.Set("Content-type", "application/vnd.github.v3+json")
	r.Header.Set("Authorization", fmt.Sprintf("Basic %s", encoded))

	resp, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("something went wrong while setReplicaCount func was sending request to GitHub API: %v", err)
	}

	defer resp.Body.Close()

	log.Debug("GitHub API's HTTP response StatusCode:" + fmt.Sprint(resp.StatusCode) + "(setReplicaCount func)") //szebben?

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil && body == nil {
		return fmt.Errorf("setReplicaCount func could not read request body: %v", err)
	}

	log.Debug("GitHub API's HTTP response body content:" + string(json.RawMessage(body)) + "(setReplicaCount func)") //szebben?
	return nil
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func generatePrimeRequest() string {
	resp, err := http.Get("http://primegenerator-service")
	if err != nil {
		log.WithError(err).Warning("Something went wrong while generatePrimeRequest func was sending request to primegenerator-service")
		return fmt.Sprint("HTTP.GET failed with error: ", err)
	}
	defer resp.Body.Close()
	log.WithFields(log.Fields{
		"Status": resp.Status,
	}).Debug("Response recieved")

	prime, _ := ioutil.ReadAll(resp.Body)
	sprime := string(prime)
	return sprime

}
func checkLoadArgs(args string) error {
	log.Debug("Checking Load command args")
	var err error
	split := strings.Split(args, ",")
	num, err := strconv.Atoi(split[0])
	if err != nil {
		return chatboterrors.ErrParamaterIsNotNumber

	} else {

		if len(split) != 2 {
			return chatboterrors.ErrWrongInsufficientNumberOfParameter

		} else if strings.ContainsAny(split[1], " ") {
			return chatboterrors.ErrParamaterIsInvalidUrl
		} else if num > 500 {
			return chatboterrors.ErrParamaterIsTooLarge

		} else if num < 1 {
			return chatboterrors.ErrParamaterIsToosmall

		} else {
			return nil
		}

	}

}
func checkSetReplicaCountArg(arg string) error {
	log.Debug("Checking SetReplicaCount command arg")
	var err error
	num, err := strconv.Atoi(arg)
	if err != nil {
		return chatboterrors.ErrParamaterIsNotNumber
	} else if num > 50 {
		return chatboterrors.ErrParamaterIsTooLargeSetReplicaCount
	} else if num < 1 {
		return chatboterrors.ErrParamaterIsToosmallSetReplicaCount
	} else {
		return nil
	}

}
func checkPrimeFactorizationArg(arg string) error {
	log.Debug("Checking PrimeFactorization command arg")
	var err error
	num, err := strconv.Atoi(arg)
	if err != nil {
		return chatboterrors.ErrParamaterIsNotNumber
	} else if num < 2 {
		return chatboterrors.ErrParamaterIsToosmallPrimeFactorization
	} else if IsPrime(num) {
		return chatboterrors.ErrParamaterIsPrime
	} else {
		return nil
	}

}
func checkConvertArg(arg string) error {
	log.Debug("Checking Convert command arg")
	var err error
	num, err := strconv.Atoi(arg)
	if err != nil {
		return chatboterrors.ErrParamaterIsNotNumber
	} else if num < 1 {
		return chatboterrors.ErrParamaterIsToosmallConvert
	} else if num > 999999999999 {
		return chatboterrors.ErrParamaterIsTooLargeConvert
	} else {
		return nil
	}

}

func main() {

	if len(os.Args) > 1 && os.Args[1] == "-v" {
		log.SetLevel(log.DebugLevel)
		log.Debug("Set log level to Debug")
	} else {
		log.SetLevel(log.WarnLevel)
	}

	purl := os.Getenv("PUBLIC_URL")
	token := os.Getenv("API_TOKEN")
	pat := os.Getenv("PAT")

	webHookURL := tgbotapi.NewWebhook(purl)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.WithError(err).Fatal("Something went wrong while setting telegram token")

	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(webHookURL)
	if err != nil {
		log.WithError(err).Fatal("Something went wrong while setting webhook")
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
		//log.Debug(req.URL.Path[len(req.URL.Path)-10:])
		//log.Debug(randomURL)
		update := chatbotstructs.MessageFromGitHub{}
		log.Debug("Request from GitHub to responseAPI")
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.WithError(err).Warn("responseAPI could not read request body")
		}
		err = json.Unmarshal(body, &update)
		if err != nil {
			log.WithError(err).Warn("responseAPI could not Unmarshal request's JSON")
		}
		chatid, _ := strconv.ParseInt(update.ChatID, 10, 64)
		msg := tgbotapi.NewMessage(chatid, "Process is completed")
		if _, err := bot.Send(msg); err != nil {
			log.WithError(err).Warn("responseAPI could not send a message to chat")
		}

	}

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/responseAPI", responseAPIHandler)

	updates := bot.ListenForWebhook("/")

	go func() { log.Panic(http.ListenAndServe(":8080", nil)) }()

	for update := range updates {
		log.Debug(fmt.Print(update))
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			if !update.Message.IsCommand() {
				log.Debug("Answering Hey by default")
				msg.Text = "Hey Buddy\nAvailable commands are the following:\n/Convert + (positive whole number as parameter, number< 999 999 999 999) Converting number into words. \n/PrimeFactorization + (positive whole number which is greater than 2, accepted as parameter)"
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}

			if update.Message.IsCommand() {

				switch update.Message.Command() {
				case "Convert":
					reqCounter.WithLabelValues("Convert").Inc()
					log.Debug("Converting number to text")
					arg := update.Message.CommandArguments()
					if err := checkConvertArg(arg); err != nil {
						log.WithError(err).Fatal()
						msg.Text = fmt.Sprint(err)
						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}
					} else {
						num, _ := strconv.Atoi(arg)
						convertedNum := convert(num)
						convertedNumFirst := convertedNum[:1]
						convertedNumRest := convertedNum[1:]
						c := cases.Title(language.Hungarian)
						convertedNumFirst = c.String(convertedNumFirst)
						convertedNum = convertedNumFirst + convertedNumRest
						msg.Text = convertedNum
						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}
					}
				case "PrimeFactorization":
					reqCounter.WithLabelValues("PrimeFactorization").Inc()
					log.Debug("Prime factorization request")
					arg := update.Message.CommandArguments()
					if err := checkPrimeFactorizationArg(arg); err != nil {
						log.WithError(err).Fatal()
						msg.Text = fmt.Sprint(err)
						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}
					} else {
						num, _ := strconv.Atoi(arg)
						result := primeFactorization(num)
						msg.Text = "Prime factors: " + result.factorsWithCommas() + "\n" + "Factor tree:" + "\n" + result.factorTree()
						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}
					}
				case "GenerateBigPrime":
					go func() {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
						reqFrequencyCounter.Inc()
						log.Debug("reqFrequencyCounter.Inc()")
						log.Debug("GenerateBigPrime request")
						reqCounter.WithLabelValues("GenerateBigPrime").Inc()
						start := time.Now()
						log.Debug("Time when calling primegenerator component: ", start)
						msg.Text = generatePrimeRequest()
						duration := time.Since(start)
						log.Debug(" PrimeGenerator component response duration in seconds: ", duration.Seconds())
						respDuration.Observe(duration.Seconds())
						log.Debug("Response arrived from primeGenerator to chat")
						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}
					}()
				case "Load":
					reqCounter.WithLabelValues("Load").Inc()
					args := update.Message.CommandArguments()
					split := strings.Split(args, ",")
					if err := checkLoadArgs(args); err != nil {
						log.WithError(err).Fatal()
						msg.Text = fmt.Sprint(err)
						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}
					} else {
						num, _ := strconv.Atoi(split[0])
						URL := split[1]
						err := loadRequest("http://loadtestingtool-service", URL, num, update.Message.Chat.ID)
						if err != nil {
							log.WithError(err).Fatal()
							msg.Text = fmt.Sprint(err)
							if _, err := bot.Send(msg); err != nil {
								log.Panic(err)
							}
						}
					}
				case "Deploy_primeGenerator":
					reqCounter.WithLabelValues("Deploy_primeGenerator").Inc()
					err := deploy("https://api.github.com/repos/bproforigoss/kemadaxbot/actions/workflows/chatbot_primeGenerator_deploy.yaml/dispatches", pat, fmt.Sprint(update.Message.Chat.ID))
					if err != nil {
						log.WithError(err).Fatal()
						msg.Text = fmt.Sprint(err)
						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}
					}
				case "Deploy_loadTestingTool":
					reqCounter.WithLabelValues("Deploy_loadTestingTool").Inc()
					err := deploy("https://api.github.com/repos/bproforigoss/kemadaxbot/actions/workflows/chatbot_loadTestingTool_deploy.yaml/dispatches", pat, fmt.Sprint(update.Message.Chat.ID))
					if err != nil {
						log.WithError(err).Fatal()
						msg.Text = fmt.Sprint(err)
						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}
					}
				case "Deploy_primeGenerator_debug":
					reqCounter.WithLabelValues("Deploy_primeGenerator_debug").Inc()
					err := deploy("https://api.github.com/repos/bproforigoss/kemadaxbot/actions/workflows/chatbot_primeGenerator_deploy_debug.yaml/dispatches", pat, fmt.Sprint(update.Message.Chat.ID))
					if err != nil {
						log.WithError(err).Fatal()
						msg.Text = fmt.Sprint(err)
						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}
					}
				case "Deploy_loadTestingTool_debug":
					reqCounter.WithLabelValues("Deploy_loadTestingTool_debug").Inc()
					err := deploy("https://api.github.com/repos/bproforigoss/kemadaxbot/actions/workflows/chatbot_loadTestingTool_deploy_debug.yaml/dispatches", pat, fmt.Sprint(update.Message.Chat.ID))
					if err != nil {
						log.WithError(err).Fatal()
						msg.Text = fmt.Sprint(err)
						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}
					}
				case "SetReplicaCount":
					reqCounter.WithLabelValues("SetReplicaCount").Inc()
					arg := update.Message.CommandArguments()
					if err := checkSetReplicaCountArg(arg); err != nil {
						log.WithError(err).Fatal()
						msg.Text = fmt.Sprint(err)
						if _, err := bot.Send(msg); err != nil {
							log.Panic(err)
						}
					} else {
						num, _ := strconv.Atoi(arg)
						url := RandStringBytes(10)
						randomURL = append(randomURL, url)
						log.Debug(url)
						err := setReplicaCount("https://api.github.com/repos/bproforigoss/kemadaxbot/actions/workflows/chatbot_set_replica.yaml/dispatches", pat, fmt.Sprint(update.Message.Chat.ID), fmt.Sprint(num), url)
						if err != nil {
							log.WithError(err).Fatal()
							msg.Text = fmt.Sprint(err)
							if _, err := bot.Send(msg); err != nil {
								log.Panic(err)
							}
						}
					}
				case "Ping":
					log.Debug("Request for Ping command")
					reqCounter.WithLabelValues("Ping").Inc()
					msg.Text = "pong"
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
				default:
					log.Debug("Response to unknown command")
					msg.Text = "I don't know that command"
					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
				}
			}

		} else {
			continue
		}

	}

}
