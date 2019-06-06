package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

var redirectURL string
var temporaryRedirectRegex *regexp.Regexp
var permanentRedirectRegex *regexp.Regexp

var nonPageRegex = regexp.MustCompile(`([^\s]+(\.(?i)(jpg|png|gif|bmp|json|js))$)`)
var pageRegex = regexp.MustCompile(`^([^\s]+(\.(?i)(html|php|asp))$)`)

var healthy bool
var ready bool

func healthEndpoint(w http.ResponseWriter, r *http.Request) {
	if healthy {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Service failing", http.StatusInternalServerError)
	}
}

func readinessEndpoint(w http.ResponseWriter, r *http.Request) {
	if ready {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Service failing", http.StatusInternalServerError)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	switch true {
	case temporaryRedirectRegex.MatchString(r.URL.RequestURI()):
		http.Redirect(w, r, redirectURL+r.URL.RequestURI(), http.StatusTemporaryRedirect)
	case permanentRedirectRegex.MatchString(r.URL.RequestURI()):
		http.Redirect(w, r, redirectURL+r.URL.RequestURI(), http.StatusPermanentRedirect)
	default:
		http.Redirect(w, r, redirectURL+r.URL.RequestURI(), http.StatusTemporaryRedirect)
	}

}

func cleanup() {
	ready = false
	log.Println("cleanup")
}

func main() {
	if os.Getenv("REDIRECT_URL") == "" {
		log.Fatal("REDIRECT_URL is needed")
	}
	redirectURL = os.Getenv("REDIRECT_URL")
	var err error

	if os.Getenv("TEMPORARY_REDIRECT_REGEX") != "" {
		temporaryRedirectRegex, err = regexp.Compile(os.Getenv("TEMPORARY_REDIRECT_REGEX"))
		if err != nil {
			log.Fatal("TEMPORARY_REDIRECT_REGEX is needed", err)
		}
	} else {
		temporaryRedirectRegex = pageRegex
	}

	if os.Getenv("PERMANENT_REDIRECT_REGEX") != "" {
		permanentRedirectRegex, err = regexp.Compile(os.Getenv("PERMANENT_REDIRECT_REGEX"))
		if err != nil {
			log.Fatal("PERMANENT_REDIRECT_REGEX is needed", err)
		}
	} else {
		permanentRedirectRegex = nonPageRegex
	}

	healthy = true
	ready = true
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}()

	redirector := http.NewServeMux()
	redirector.HandleFunc("/", redirect)
	go func() {
		err := http.ListenAndServe(":9090", redirector)
		if err != nil {
			healthy = false
			log.Fatal("ListenAndServe: ", err)
		}
	}()
	controller := http.NewServeMux()
	controller.HandleFunc("/healthy", healthEndpoint)
	controller.HandleFunc("/ready", readinessEndpoint)
	err = http.ListenAndServe(":9091", controller)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
