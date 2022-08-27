package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-rfc/sse/pkg/eventsource"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sync"
	"testing"
	"time"
	"unicode/utf8"
)

func TestListen(t *testing.T) {

	var tests = []struct {
		clientsCount int
	}{
		{1},
		{10},
		{100},
		{1000},
		{5000},
	}

	fmt.Println("Тестируем /listen")
	for _, test := range tests {
		fmt.Printf("Тестируем на %d клиентах\n", test.clientsCount)
		srv := httptest.NewServer(handlers())

		esources := make([]*eventsource.EventSource, 0, 0)
		for i := 0; i < test.clientsCount; i++ {
			es, err := eventsource.New(srv.URL + "/listen")
			if err != nil {
				fmt.Println("Ошибка создания eventsource", err)
				continue
			}
			esources = append(esources, es)
		}

		datetimes := make([]int64, len(esources), len(esources))

		for i := 0; i < 10; i++ {
			for j, es := range esources {
				event := <-es.MessageEvents()

				datetime, text, err := parseMessage(event.Data)

				if err != nil {
					t.Error("Ошибка парсинга даты и времени")
				}

				duration := retry

				if datetimes[j] != 0 {
					duration = int(datetime - datetimes[j])
				}

				if duration/retry > 3 {
					t.Error("Превышение разницы между двумя собщениями более чем в 3 раза")
				}

				if text != startMessage {
					t.Error("Стартовое сообщение неправильное")
				}

			}
		}

		for i := range esources {
			esources[i].Close()
		}

		srv.Close()

		fmt.Printf("Тест на %d клиентах успешно пройден!\n", test.clientsCount)
	}
}

func TestSay(t *testing.T) {
	var tests = []struct {
		clientsCount int
	}{
		{1},
		{10},
		{100},
		{1000},
		{5000},
	}

	fmt.Println("Тестируем /say")
	for _, test := range tests {
		fmt.Printf("Тестируем на %d клиентах\n", test.clientsCount)

		srv := httptest.NewServer(handlers())
		fmt.Println(srv.URL)

		esources := make([]*eventsource.EventSource, 0, 0)
		for i := 0; i < test.clientsCount; i++ {
			es, err := eventsource.New(srv.URL + "/listen")
			if err != nil {
				t.Error("Ошибка создания eventsource")
				continue
			}
			esources = append(esources, es)
		}

		sayInd := rand.Intn(10)
		newMessage := "Всё идёт по плану"

		for i := 0; i < 10; i++ {
			if i == sayInd {
				mu := new(sync.Mutex)
				mu.Lock()

				reqData := Request{
					Word: newMessage,
				}

				jsonData, err := json.Marshal(reqData)
				if err != nil {
					t.Error("Ошибка парсинга реквеста")
				}

				_, err = http.Post(srv.URL+"/say", "application/json", bytes.NewBuffer(jsonData))
				if err != nil {
					t.Error("Ошибка отправки json-запроса")
				}

				mu.Unlock()
			}

			for _, es := range esources {

				event := <-es.MessageEvents()

				_, text, err := parseMessage(event.Data)
				if err != nil {
					t.Error("Ошибка парсинга даты и времени")
				}

				if i >= sayInd && text != newMessage {
					t.Error("Не работает say")
					return

				}
			}
		}

		for i := range esources {
			esources[i].Close()
		}

		srv.Close()

		fmt.Printf("Тест эндпойнта /say на %d клиентах успешно пройден!\n", test.clientsCount)
	}
}

func parseMessage(message string) (int64, string, error) {
	r, err := regexp.Compile(fmt.Sprintf("^(?P<datetime>.{%d})\\s(?P<text>.+)", utf8.RuneCountInString(timeLayout)))
	if err != nil {
		return 0, "", err
	}

	m := r.FindStringSubmatch(message)
	if m == nil {
		panic("mo match")
	}

	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = m[i]
		}
	}

	datetime, err := time.Parse(timeLayout, result["datetime"])
	if err != nil {
		return 0, "", err
	}

	return datetime.UnixMilli(), result["text"], nil
}
