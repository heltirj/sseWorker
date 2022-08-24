package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type SSEChannel struct {
	Word string
}

type Request struct {
	Word string `json:"word"`
}

func main() {

	s := SSEChannel{
		Word: "Привет",
	}

	http.Handle("/", http.FileServer(http.Dir("client")))
	http.HandleFunc("/say", s.say)
	http.HandleFunc("/listen", s.listen)

	log.Println("Запуск сервера на http://localhost:4000")
	err := http.ListenAndServe(":4000", nil)
	log.Fatal(err)

}

func (s *SSEChannel) listen(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	fmt.Fprintf(w, "retry: 1000\n\ndata: %s %s\n\n", s.Word, time.Now().Format("15:04:05"))
}

func (s *SSEChannel) say(w http.ResponseWriter, r *http.Request) {
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

	s.Word = message.Word

	w.Write([]byte(message.Word))
}
