package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

func init() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	rand.Seed(time.Now().UnixNano())

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

func generateBigPrime() int {
	min := 10000000
	max := 100000000
	randint := rand.Intn(max-min+1) + min
	for i := randint; i > 2; i-- {

		if IsPrime(i) {
			return i
		}

	}
	return 0

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

	responseAPIHandler := func(w http.ResponseWriter, req *http.Request) {
		log.Debug(req.URL.Path)
		prime := generateBigPrime()
		strPrime := strconv.Itoa(prime)

		_, err := w.Write([]byte(strPrime))
		if err != nil {
			log.WithError(err).Warn("Hello handler failed while writing response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/", responseAPIHandler)

	log.Panic(http.ListenAndServe(":8080", nil))

}
