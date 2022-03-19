package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/telebot.v3"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Store struct {
	sync.RWMutex
	Data Data
	File *os.File
}

type Data struct {
	Users map[string]*UserData `json:"users"`
}

type UserData struct {
	ID            string                 `json:"id"`
	Subscriptions []Subscription         `json:"subscriptions"`
	Notified      map[time.Time]struct{} `json:"notified"`
}

func (u *UserData) addToNotified(t time.Time) {
	if u.Notified == nil {
		u.Notified = make(map[time.Time]struct{})
	}
	u.Notified[t] = struct{}{}
}

func NewUserData(id string) *UserData {
	return &UserData{
		ID:            id,
		Subscriptions: make([]Subscription, 0),
		Notified:      make(map[time.Time]struct{}),
	}
}

type Subscription struct {
	Weekday time.Weekday `json:"weekday"`
	Time    Clock        `json:"time"`
	Hours   int          `json:"hours"`
}

func ParseTimeRange(input string) (Subscription, error) {
	data := strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")
	if len(data) < 2 || len(data) > 3 {
		return Subscription{}, errors.New(`incorrect time format, expected format: "Mon 15:00 2"`)
	}
	weekday, ok := weekdayMapping[data[0]]
	if !ok {
		return Subscription{}, fmt.Errorf("unknown day of the week: %s", weekday)
	}

	startClock, err := parseTime(data[1])
	if err != nil {
		return Subscription{}, err
	}
	hours := 1
	if len(data) >= 3 {
		hours, err = strconv.Atoi(data[2])
		if err != nil {
			return Subscription{}, err
		}
	}
	if hours < 1 {
		return Subscription{}, errors.New("hours must be equal or greater than 1")
	}

	return Subscription{
		Weekday: weekday,
		Time:    startClock,
		Hours:   hours,
	}, nil
}

func (r *Subscription) String() string {
	return fmt.Sprintf("%s %02d:%02d-%02d:%02d", r.Weekday.String()[:3], r.Time.Hour, r.Time.Minute, r.Time.Hour+r.Hours, r.Time.Minute)
}

type Clock struct {
	Hour   int `json:"hours"`
	Minute int `json:"minutes"`
}

func NewStore(path string) (*Store, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	var data Data
	err = json.NewDecoder(file).Decode(&data)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if len(data.Users) == 0 {
		data.Users = make(map[string]*UserData)
	}
	return &Store{
		File: file,
		Data: data,
	}, nil
}

func (s *Store) addTime(userID string, time Subscription) {
	user, ok := s.Data.Users[userID]
	if !ok {
		user = NewUserData(userID)
		s.Data.Users[userID] = user
	}
	// avoid duplicates
	for _, t := range user.Subscriptions {
		if t == time {
			return
		}
	}
	user.Subscriptions = append(user.Subscriptions, time)
}

func (s *Store) removeTime(userID string, time Subscription) error {
	user, ok := s.Data.Users[userID]
	if !ok {
		return errors.New("user not found")
	}
	for i, t := range user.Subscriptions {
		if t == time {
			user.Subscriptions = append(user.Subscriptions[:i], user.Subscriptions[i+1:]...)
			return nil
		}
	}
	return errors.New("time not found")
}

var weekdayMapping = map[string]time.Weekday{
	"mon":       time.Monday,
	"monday":    time.Monday,
	"tue":       time.Tuesday,
	"tuesday":   time.Tuesday,
	"wed":       time.Wednesday,
	"wednesday": time.Wednesday,
	"thu":       time.Thursday,
	"thursday":  time.Thursday,
	"fri":       time.Friday,
	"friday":    time.Friday,
	"sat":       time.Saturday,
	"saturday":  time.Saturday,
	"sun":       time.Sunday,
	"sunday":    time.Sunday,
}

// time format example "Mon 15:00"
// TODO: "Weekday 15:00"
// TODO: "Weekend 15:00"
// TODO: "Mon-Tue 15:00-17:00"
func (s *Store) Subscribe(userID string, input string) error {
	if len(userID) == 0 {
		return errors.New("userID can't be blank")
	}
	tr, err := ParseTimeRange(input)
	if err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	s.addTime(userID, tr)
	err = s.save()
	if err != nil {
		return fmt.Errorf("could not save data: %w", err)
	}
	return nil
}

func (s *Store) Unsubscribe(userID, input string) error {
	if len(userID) == 0 {
		return errors.New("userID can't be blank")
	}
	tr, err := ParseTimeRange(input)
	if err != nil {
		return fmt.Errorf("incorrect time format: %w", err)
	}
	s.Lock()
	defer s.Unlock()
	err = s.removeTime(userID, tr)
	if err != nil {
		return err
	}
	err = s.save()
	if err != nil {
		return fmt.Errorf("could not save data: %w", err)
	}
	return nil
}

func (s *Store) DeleteUser(userID string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.Data.Users, userID)
	return s.save()
}

func (s *Store) Subscriptions(userID string) string {
	s.RLock()
	defer s.RUnlock()
	user, ok := s.Data.Users[userID]
	if !ok {
		return ""
	}
	var subMsgs []string
	for _, tr := range user.Subscriptions {
		subMsgs = append(subMsgs, tr.String())
	}
	if len(subMsgs) == 0 {
		return "no subscriptions"
	}
	return strings.Join(subMsgs, "\n")
}

func (s *Store) NotifyAll(b *telebot.Bot, cal Calendar) {
	s.RLock()
	defer s.RUnlock()
	for _, user := range s.Data.Users {
		err := s.notifyUser(b, user, cal)
		if err != nil {
			log.Println("could not notify user", err)
			continue
		}
	}
}

func (s *Store) notifyUser(b *telebot.Bot, user *UserData, cal Calendar) error {
	userCal := cal.ForUserSubscriptions(user)
	if len(userCal) == 0 {
		return nil
	}
	id, err := strconv.Atoi(user.ID)
	if err != nil {
		return err
	}
	teleUser := &telebot.User{ID: int64(id)}
	msg := fmt.Sprintf("New booking available:\n%s", userCal.String())
	_, err = b.Send(teleUser, msg)
	if err != nil {
		return err
	}
	for t := range userCal {
		user.addToNotified(t)
	}
	return s.save()
}

func parseTime(timeS string) (Clock, error) {
	clock := strings.Split(timeS, ":")
	hour, err := strconv.ParseInt(clock[0], 10, 64)
	if err != nil {
		return Clock{}, fmt.Errorf("incorrect time format: \"%s\"", timeS)
	}
	if hour < 0 || hour > 23 {
		return Clock{}, errors.New("hour should be between 0 and 23")
	}
	minute, err := strconv.ParseInt(clock[1], 10, 64)
	if err != nil {
		return Clock{}, fmt.Errorf("incorrect time format: \"%s\"", timeS)
	}
	if minute < 0 || minute > 59 {
		return Clock{}, errors.New("minute should be between 0 and 59")
	}
	return Clock{
		Hour:   int(hour),
		Minute: int(minute),
	}, nil
}

func (s *Store) save() error {
	if err := s.File.Truncate(0); err != nil {
		return err
	}
	if _, err := s.File.Seek(0, 0); err != nil {
		return err
	}
	enc := json.NewEncoder(s.File)
	enc.SetIndent("", "  ")
	return enc.Encode(s.Data)
}

type Calendar map[time.Time]uint

func NewCalendar(start, end time.Time) Calendar {
	c := make(Calendar)
	start = time.Date(start.Year(), start.Month(), start.Day(), StartHour, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), EndHour, 0, 0, 0, end.Location())
	for t := start; t.Before(end); t = t.Add(MinBookingDuration) {
		if t.Hour() < StartHour || t.Hour() > EndHour {
			continue
		}
		c[t] = CourtsCount
	}
	return c
}

func (cal Calendar) Book(b Booking) {
	// some bookings start and end at :30 minutes mark.
	// for example 5:30-6:30, such interval are unavailable for us, so we need to reserve 5:00-7:00 slot in this case
	start := b.Start
	start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), 0, 0, 0, start.Location())
	end := b.End
	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), 0, 0, 0, end.Location())
	if b.End.Sub(end) > 0 {
		end = end.Add(time.Hour)
	}
	for t := start; t.Before(end); t = t.Add(MinBookingDuration) {
		if cal[t] > 0 {
			cal[t]--
		}
	}
}

func (cal Calendar) IsAvailable(start, end time.Time) bool {
	if end.After(start) {
		return false
	}
	for t := start; t.Before(end); t = t.Add(MinBookingDuration) {
		if cal[t] == 0 {
			return false
		}
	}
	return true
}

func (cal Calendar) NonZero() Calendar {
	c := make(Calendar)
	for t, v := range cal {
		if v > 0 {
			c[t] = v
		}
	}
	return c
}

func (cal Calendar) String() string {
	var buf bytes.Buffer
	for _, v := range cal.toSlice() {
		if v.Count > 0 {
			buf.WriteString(fmt.Sprintf("%s - %d\n", v.Time.Format("2006-01-02 15:04"), v.Count))
		}
	}
	return buf.String()
}

func (cal Calendar) ForSubscription(subscription Subscription) Calendar {
	result := make(Calendar)
	for t, slots := range cal {
		if t.Weekday() != subscription.Weekday {
			continue
		}
		if t.Hour() != subscription.Time.Hour {
			continue
		}
		for i := 1; i <= subscription.Hours; i++ {
			nextSlots, ok := cal[t.Add(time.Hour*time.Duration(i))]
			if !ok {
				continue
			}
			if nextSlots < slots {
				nextSlots = slots
			}
		}
		result[t] = slots
	}
	return result
}

func (cal Calendar) ForUserSubscriptions(user *UserData) Calendar {
	result := make(Calendar)
	for _, subscription := range user.Subscriptions {
		subCal := cal.ForSubscription(subscription)
		for k, v := range subCal {
			result[k] = v
		}
	}
	for t := range user.Notified {
		delete(result, t)
	}
	return result
}

type Entry struct {
	Time  time.Time
	Count uint
}

func (cal Calendar) toSlice() []Entry {
	var entries []Entry
	for t, v := range cal {
		entries = append(entries, Entry{t, v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Time.Before(entries[j].Time)
	})
	return entries
}

type Booking struct {
	Court int       `json:"resourceId"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func fetchData() ([]Booking, error) {
	const layout = "2006-01-02"
	url := fmt.Sprintf("https://platform.aklbadminton.com/api/booking/feed?start=%s&end=%s", time.Now().Format(layout), time.Now().Add(time.Hour*24*8).Format(layout))
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

func fetchCalendar() (Calendar, error) {
	data, err := fetchData()
	if err != nil {
		return nil, err
	}
	cal := NewCalendar(time.Now(), time.Now().Add(time.Hour*24*7))
	for _, b := range data {
		cal.Book(b)
	}
	return cal, nil
}
