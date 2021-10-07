package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type DogFacts struct {
	DogFacts []*Dogfact
}

type Dogfact struct {
	Dogfact string `json:"fact"`
}

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Type      string `json:"type"`
}
type From struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LanguageCode string `json:"language_code"`
}
type Message struct {
	MessageID int    `json:"message_id"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
	Chat      Chat
	From      From
}

type Update struct {
	UpdateID int `json:"update_id"`
	Message  Message
}
type sendMessageReqBody struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

func getFact() ([]byte, error) {

	resp, err := http.Get("https://dog-facts-api.herokuapp.com/api/v1/resources/dogs?number=1")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.WithFields(log.Fields{
		"Status": resp.Status,
	}).Debug("Response recieved")

	fact, _ := ioutil.ReadAll(resp.Body)
	log.Debug("Request sent to dogFact api")
	return fact, err

}
func init() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	var portNumber string = "8080"

	if len(os.Args) > 1 && os.Args[1] == "-v" {
		log.SetLevel(log.DebugLevel)
		//loglevel=>Debuglevel
		log.Debug("Set loglevel to Debug")
	} else {
		log.SetLevel(log.WarnLevel)
		//loglevel=>Warnlevel

	}

	//setwebhook
	purl := os.Getenv("PUBLIC_URL")
	url := "https://api.telegram.org/bot2042481517:AAGd1WViLeY8fpNEdmkDF1C0qIjlr1i6p4g/setWebhook?url=" + purl
	print(url)

	resp, err := http.Get(url)
	if err != nil {
		log.WithError(err).Warn("Bye handler failed while writing response")
	}
	defer resp.Body.Close()
	log.WithFields(log.Fields{
		"Status": resp.Status,
	}).Debug("Response recieved")

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		log.Debug("Request for saying Hello")
		_, err := fmt.Fprintf(w, "Hello %s", os.Getenv("USER"))
		if err != nil {
			log.WithError(err).Warn("Hello handler failed while writing response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	byeHandler := func(w http.ResponseWriter, req *http.Request) {
		log.Debug("Request for saying Bye")
		_, err := fmt.Fprintf(w, "Bye %s", os.Getenv("USER"))
		if err != nil {
			log.WithError(err).Warn("Bye handler failed while writing response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
	dogFactHandler := func(w http.ResponseWriter, req *http.Request) {
		log.Debug("Request for dogfact")
		fact := DogFacts{}
		resp, err := getFact()
		if err != nil {
			log.WithError(err).Warn("Something went wrong while trying to get a fact dogFact from api")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(resp, &fact.DogFacts)
		if err != nil {
			log.WithError(err).Warn("Unmarshal failed")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = io.WriteString(w, fact.DogFacts[0].Dogfact)
		if err != nil {
			log.WithError(err).Warn("Something wrong with struct")
		}

	}
	hiBuddyHandler := func(w http.ResponseWriter, req *http.Request) {
		update := Update{}
		log.Debug("Request from telegram")
		body, _ := ioutil.ReadAll(req.Body)
		err = json.Unmarshal(body, &update)
		if err != nil {
			log.WithError(err).Warn("Unmarshal failed")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		url := "https://api.telegram.org/bot2042481517:AAGd1WViLeY8fpNEdmkDF1C0qIjlr1i6p4g/sendMessage"

		reqBody := &sendMessageReqBody{
			ChatID: update.Message.Chat.ID,
			Text:   "Hi Buddy!",
		}
		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			log.WithError(err).Warn("Marshal failed")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err = http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
		if err != nil {
			log.WithError(err).Warn("Bye handler failed while writing response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Debug(resp.Status)
		defer resp.Body.Close()
	}

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/bye", byeHandler)
	http.HandleFunc("/dog", dogFactHandler)
	http.HandleFunc("/", hiBuddyHandler)

	log.WithFields(log.Fields{
		"portNumber": portNumber,
	}).Info("Server is starting")

	log.Panic(http.ListenAndServe(":8080", nil))

}
