package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"src/github.com/bproforigoss/kemadaxbot/Structs/ChatbotStructs"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	rand.Seed(time.Now().UnixNano())

}

type TelegramBotRequest struct {
	UpdateID int                       `json:"update_id"`
	Message  TelegramBotRequestMessage `json:"message"`
}
type TelegramBotRequestMessage struct {
	MessageID int                               `json:"message_id"`
	Date      int64                             `json:"date"`
	Chat      TelegramBotRequestMessageChat     `json:"chat"`
	Entities  TelegramBotRequestMessageEntities `json:"entities"`
	Text      string                            `json:"text"`
}
type TelegramBotRequestMessageEntities []struct {
	Type   string `json:"type"`
	Length int    `json:"length"`
}
type TelegramBotRequestMessageChat struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

/*type RequestBody struct { //Ez veszélyes ---->külön package
	URL              string `json:"url"`
	RequestNumber    int    `json:"number"`
	RequestFrequency int    `json:"frequency"`
	RequestChatId    int    `json:"chat_id"`
}*/

/*{"update_id":1,"message":{"message_id":1,"date":1649352456,"chat":{"id":2006716105,"type":"private"},
"entities":[{"type":"bot_command","length":17}],"text":"/GenerateBigPrime"}*/
func load(URL string, ChatID int) error {
	client := &http.Client{}
	reqBody := TelegramBotRequest{
		UpdateID: 1,
		Message: TelegramBotRequestMessage{MessageID: 1, Date: time.Now().Unix(), Chat: TelegramBotRequestMessageChat{ID: ChatID, Type: "private"},
			Entities: TelegramBotRequestMessageEntities{{Type: "bot_command", Length: 17}}, Text: "/GenerateBigPrime",
		}}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.WithError(err).Fatal("load func marshal JSON failed")
		return fmt.Errorf("load func marshal JSON failed: %v", err)
	}

	r, _ := http.NewRequest("POST", URL, bytes.NewBuffer(reqBytes))
	if err != nil {
		log.WithError(err).Fatal("load func making new request failed")
		return fmt.Errorf("load func making new request failed: %v", err)
	}
	r.Header.Set("Content-type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		log.WithError(err).Fatal("Something went wrong while load func  was sending request to Chatbot")
		return fmt.Errorf("Something went wrong while  func  was sending request to GitHub API: %v", err)
	}

	defer resp.Body.Close()

	log.Debug(" GenerateBigPrime microservice API's HTTP response StatusCode" + fmt.Sprint(resp.StatusCode))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil && body == nil {
		log.WithError(err).Fatal("load func could not read response body")
		return fmt.Errorf("load func could not read response body: %v", err)
	}

	log.Debug(" GenerateBigPrime microservice API's HTTP response body content" + string(json.RawMessage(body)))
	return nil

}

func main() {

	if len(os.Args) > 1 && os.Args[1] == "-v" {
		log.SetLevel(log.DebugLevel)
		log.Debug("Set loglevel to Debug")
	} else {
		log.SetLevel(log.WarnLevel)
	}

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		log.Debug("Request for saying Hello")
		_, err := fmt.Fprintf(w, "Hello %s", os.Getenv("USER"))
		if err != nil {
			log.WithError(err).Warn("Hello handler failed while writing response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	loadTestingReqHandler := func(w http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var reqBody ChatbotStructs.RequestToLoad
		err := decoder.Decode(&reqBody)
		if err != nil {
			log.WithError(err).Warn("Unmarshal JSON failed at loadTestingHandler")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ticker := time.NewTicker(time.Duration(reqBody.RequestFrequency) * time.Second)
		defer ticker.Stop()
		reqNum := reqBody.RequestNumber
		i := 0
		for range ticker.C {
			err := load(reqBody.URL, reqBody.RequestChatId)
			if err != nil {
				log.WithError(err).Warn("loadTestingReqHandler calling load func() failed")
			}
			i++
			if i == reqNum {
				break

			}

		}

	}
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", loadTestingReqHandler)

	log.Panic(http.ListenAndServe(":8080", nil))

}
