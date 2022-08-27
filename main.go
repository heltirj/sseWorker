package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	startMessage = "Каждый из нас - беспонтовый пирожок"
	apiAddress   = ":4000"
	timeLayout   = "02.01.06 15:04:05"
	retry        = 1000
)

type SSEWorker struct {
	Text string
	m    *sync.Mutex
}

type Request struct {
	Word string `json:"word"`
}

func main() {
	runServer()
}

func runServer() {
	log.Fatal(http.ListenAndServe(apiAddress, handlers()))
}

func handlers() http.Handler {
	r := http.NewServeMux()

	s := SSEWorker{
		Text: startMessage,
		m:    new(sync.Mutex),
	}

	r.Handle("/", http.FileServer(http.Dir("client")))
	r.HandleFunc("/say", s.say)
	r.HandleFunc("/listen", s.listen)

	return r
}

func (s *SSEWorker) listen(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")

	fmt.Fprintf(w, "retry: %d\n\ndata: %s %s\n\n", retry, time.Now().Format(timeLayout), s.Text)
}

func (s *SSEWorker) say(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Это POST-метод. идиот!"))
		return
	}

	var message Request
	err := json.NewDecoder(r.Body).Decode(&message)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.m.Lock()
	s.Text = message.Word
	s.m.Unlock()

	w.Write([]byte(message.Word))
}
