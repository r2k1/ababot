package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/telebot.v3/middleware"
	"log"
	"net/http"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

const CourtsCount = 12
const MinBookingDuration = 60 * time.Minute

func main() {
	pref := tele.Settings{
		Token:  "",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	store, err := NewStore("data.json")
	checkErr(err)

	b, err := tele.NewBot(pref)
	checkErr(err)
	b.Use(middleware.Logger())

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

type Booking struct {
	Court int       `json:"resourceId"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func fetchData() ([]Booking, error) {
	const layout = "2006-01-02"
	url := fmt.Sprintf("https://platform.aklbadminton.com/api/booking/feed?start=%s&end=%s", time.Now().Format(layout), time.Now().Add(time.Hour*24*7).Format(layout))
	log.Println("fetching", url)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected error code %d", resp.StatusCode)
	}
	var data []Booking
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	log.Printf("fetched %d bookings", len(data))
	return data, nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Calendar map[time.Time]uint

func NewCalendar(start, end time.Time) Calendar {
	c := make(Calendar)
	for t := start; t.Before(end); t = t.Add(MinBookingDuration) {
		c[t] = CourtsCount
	}
	return c
}

func (c Calendar) Book(b Booking) {
	start := b.Start
	start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), start.Minute(), 0, 0, start.Location())
	// some bookings start and end at :30 minutes mark.
	// for example 5:30-6:30, such interval are unavailable for us, so we need to reserve 5:00-7:00 slot in this case
	end := b.End
	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), end.Minute(), 0, 0, end.Location())
	if b.End.Sub(end) > 0 {
		end = end.Add(time.Hour)
	}
	for t := b.Start; t.Before(b.End); t = t.Add(MinBookingDuration) {
		c[t]--
	}
}

func (c Calendar) Available(start, end time.Time) bool {
	if end.After(start) {
		return false
	}
	for t := start; t.Before(end); t = t.Add(MinBookingDuration) {
		if c[t] == 0 {
			return false
		}
	}
	return true
}
