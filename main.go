package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/telebot.v3/middleware"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

const CourtsCount = 12
const MinBookingDuration = 30 * time.Minute

var testUser = UserData{
	ID: "167935153",
	Times: []Weektime{
		{
			Weekday: 6,
			Hour:    6,
			Minute:  0,
		},
		{
			Weekday: 6,
			Hour:    7,
			Minute:  0,
		},
	},
}

func main() {
	pref := tele.Settings{
		Token:  "",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	store, err := NewStore("data.json")
	checkErr(err)

	//checkErr(store.Save())

	b, err := tele.NewBot(pref)
	checkErr(err)
	b.Use(middleware.Logger())

	b.Handle("/check", func(c tele.Context) error {
		data, err := fetchData()
		if err != nil {
			return err
		}
		availableSlots := availableSlots(data)
		return c.Send(availableSlots.String(), tele.ModeHTML)
	})
	b.Handle("/filter", func(c tele.Context) error {
		data, err := fetchData()
		if err != nil {
			return err
		}
		availableSlots := availableSlots(data).filterTimes(testUser.Times)
		if len(availableSlots) == 0 {
			return c.Send("no available times for selected times")
		}
		return c.Send(availableSlots.String(), tele.ModeHTML)

	})
	b.Handle("/subscribe", func(c tele.Context) error {
		id := strconv.FormatInt(c.Sender().ID, 10)
		err := store.AddTime2(id, c.Data())
		if err != nil {
			return c.Send(err.Error())
		}
		return c.Send("noted")
	})
	log.Println("listening to messages")
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

func availableSlots(bookings map[time.Time]int) Slots {
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

	result := make(Slots)
	for t := minTime; t.Before(maxTime); t = t.Add(MinBookingDuration) {
		// closed
		if t.Hour() < 6 || t.Hour() > 18 {
			continue
		}
		if bookings[t] >= CourtsCount {
			continue
		}
		result[t] = CourtsCount - bookings[t]
	}
	return result
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Slots map[time.Time]int

type Slot struct {
	Timestamp time.Time
	Courts    int
}

func (c Slots) toSlice() []Slot {
	result := make([]Slot, 0, len(c))
	for t, v := range c {
		result = append(result, Slot{Timestamp: t, Courts: v})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	return result
}

func (c Slots) String() string {
	var result string
	for _, slot := range c.toSlice() {
		result += fmt.Sprintf("<code>%s - %d</code>\n", slot.Timestamp.Format("Mon 2006-01-02 15:04"), slot.Courts)
	}
	return result
}

func (c Slots) filterTimes(times []Weektime) Slots {
	result := make(Slots)
	for k, v := range c {
		for _, t := range times {
			if t.Weekday == k.Weekday() && t.Hour == k.Hour() && t.Minute == k.Minute() {
				result[k] = v
			}
		}
	}
	return result
}
