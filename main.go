package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3/middleware"
	"log"
	"os"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

const CourtsCount = 12
const MinBookingDuration = 60 * time.Minute
const StartHour = 6
const EndHour = 22

func main() {
	err := godotenv.Load(".env")
	pref := tele.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	store, err := NewStore("data.json")
	checkErr(err)

	b, err := tele.NewBot(pref)
	checkErr(err)
	b.Use(middleware.Logger())

	b.Handle("/available", func(c tele.Context) error {
		data, err := fetchData()
		if err != nil {
			return c.Send(err.Error())
		}
		cal := NewCalendar(time.Now(), time.Now().Add(time.Hour*24*7))
		for _, b := range data {
			cal.Book(b)
		}
		return c.Send(limitString(fmt.Sprintf("%v", cal.String()), 4096))
	})

	b.Handle("/subscriptions", func(c tele.Context) error {
		id := strconv.FormatInt(c.Sender().ID, 10)
		msg := fmt.Sprintf("Current subscriptions:\n%s", store.Subscriptions(id))
		return c.Send(msg)
	})

	b.Handle("/subscribe", func(c tele.Context) error {
		id := strconv.FormatInt(c.Sender().ID, 10)
		err := store.Subscribe(id, c.Data())
		if err != nil {
			return c.Send(err.Error())
		}
		msg := fmt.Sprintf("You are subscribed to\n%s", store.Subscriptions(id))
		return c.Send(msg)
	})

	b.Handle("/unsubscribe", func(c tele.Context) error {
		id := strconv.FormatInt(c.Sender().ID, 10)
		err := store.Unsubscribe(id, c.Data())
		if err != nil {
			return c.Send(err.Error())
		}
		msg := fmt.Sprintf("Unsubscribed. Current subscriptions:\n%s", store.Subscriptions(id))
		return c.Send(msg)
	})

	b.Handle("/clean", func(c tele.Context) error {
		id := strconv.FormatInt(c.Sender().ID, 10)
		err := store.DeleteUser(id)
		if err != nil {
			return c.Send(err.Error())
		}
		return c.Send("Unsubscribed from everything")
	})

	b.OnError = func(err error, c tele.Context) {
		log.Println(err)
	}

	go func() {
		log.Println("starting refresher")
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		check := func() {
			log.Println("refreshing data")
		}

		check()
		for range ticker.C {
			check()
		}
	}()

	log.Println("listening to messages")
	b.Start()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func limitString(in string, l int) string {
	if len(in) > l {
		return in[:l]
	}
	return in
}
