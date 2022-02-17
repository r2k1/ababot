package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	tele "gopkg.in/telebot.v3"
)

type UserData struct {
	ID       string
	Times    []Weektime
	Notified []time.Time
}

type Weektime struct {
	Weekday time.Weekday
	// seconds from the beginning of the day
	Seconds int
}

const CourtsCount = 12
const MinBookingDuration = 30 * time.Minute

var testUser = UserData{
	ID: "167935153",
	Times: []Weektime{
		{
			Weekday: 7,
			Seconds: 6 * 60 * 60,
		},
		{
			Weekday: 7,
			Seconds: 7 * 60 * 60,
		},
	},
}

func main() {
	pref := tele.Settings{
		Token:  "",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	checkErr(err)

	b.Handle("/hello", func(c tele.Context) error {
		log.Print("hi")
		return c.Send("Hello!")
	})

	chat, err := b.ChatByUsername("167935153")
	checkErr(err)

	data, err := fetchData()
	checkErr(err)

	msg, err := b.Send(chat, fmt.Sprintf("%+v", data)[0:399])
	checkErr(err)
	log.Print(msg)

	b.Start()

}

type Booking struct {
	Court int       `json:"resourceId"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func fetchData() (map[time.Time]int, error) {
	const layout = "2006-01-02"
	url := fmt.Sprintf("https://platform.aklbadminton.com/api/booking/feed?start=%s&end=%s", time.Now().Format(layout), time.Now().Add(time.Hour*24*2).Format(layout))
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
	return toMap(data), nil
}

func toMap(data []Booking) map[time.Time]int {
	result := make(map[time.Time]int)
	for _, booking := range data {
		for t := booking.Start; t.Before(booking.End); t = t.Add(MinBookingDuration) {
			result[t]++
		}
	}
	return result
}

func availableSlots(bookings map[time.Time]int) Calendar {
	var minTime time.Time
	var maxTime time.Time
	for t := range bookings {
		if minTime.IsZero() || t.Before(minTime) {
			minTime = t
		}
		if maxTime.IsZero() || t.After(maxTime) {
			maxTime = t
		}
	}
	minTime = time.Date(minTime.Year(), minTime.Month(), minTime.Day(), 6, 0, 0, 0, minTime.Location())
	maxTime = time.Date(maxTime.Year(), maxTime.Month(), maxTime.Day(), 24, 0, 0, 0, maxTime.Location())

	result := make(Calendar)
	for t := minTime; t.Before(maxTime); t = t.Add(MinBookingDuration) {
		// closed
		if t.Hour() < 6 || t.Hour() > 18 {
			continue
		}
		if bookings[t] >= CourtsCount {
			continue
		}
		result[t] = nil
	}
	return result
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Calendar map[time.Time]interface{}

func (c Calendar) toSlice() []time.Time {
	result := make([]time.Time, 0, len(c))
	for t := range c {
		result = append(result, t)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Before(result[j])
	})
	return result
}
