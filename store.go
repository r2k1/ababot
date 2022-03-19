package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	ID       string                 `json:"id"`
	Times    []TimeRange            `json:"times"`
	Notified map[time.Time]struct{} `json:"notified"`
}

func NewUserData(id string) *UserData {
	return &UserData{
		ID:       id,
		Times:    make([]TimeRange, 0),
		Notified: make(map[time.Time]struct{}),
	}
}

type TimeRange struct {
	Weekday time.Weekday `json:"weekday"`
	Start   Clock
	End     Clock
}

func ParseTimeRange(input string) (TimeRange, error) {
	data := strings.Split(strings.ToLower(strings.TrimSpace(input)), " ")
	if len(data) != 2 {
		return TimeRange{}, errors.New(`incorrect time format, expected format: "Mon 15:00-17:00"`)
	}
	weekday, ok := weekdayMapping[data[0]]
	if !ok {
		return TimeRange{}, fmt.Errorf("unknown day of the week: %s", weekday)
	}

	clocks := strings.Split(data[1], "-")
	if (len(clocks)) != 2 {
		return TimeRange{}, errors.New("time should be provided in 13:00-15:00 format")
	}
	startClock, err := parseTime(clocks[0])
	if err != nil {
		return TimeRange{}, err
	}
	endClock, err := parseTime(clocks[1])
	if err != nil {
		return TimeRange{}, err
	}

	return TimeRange{
		Weekday: weekday,
		Start:   startClock,
		End:     endClock,
	}, nil
}

func (r *TimeRange) String() string {
	return fmt.Sprintf("%s %02d:%02d-%02d:%02d", r.Weekday.String()[:3], r.Start.Hour, r.Start.Minute, r.End.Hour, r.End.Minute)
}

type Clock struct {
	Hour   uint `json:"hours"`
	Minute uint `json:"minutes"`
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

func (s *Store) addTime(userID string, time TimeRange) {
	user, ok := s.Data.Users[userID]
	if !ok {
		user = NewUserData(userID)
		s.Data.Users[userID] = user
	}
	// avoid duplicates
	for _, t := range user.Times {
		if t == time {
			return
		}
	}
	user.Times = append(user.Times, time)
}

func (s *Store) removeTime(userID string, time TimeRange) error {
	user, ok := s.Data.Users[userID]
	if !ok {
		return errors.New("user not found")
	}
	for i, t := range user.Times {
		if t == time {
			user.Times = append(user.Times[:i], user.Times[i+1:]...)
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
		return fmt.Errorf("incorrect time format: %w", err)
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
	var ranges []string
	for _, tr := range user.Times {
		ranges = append(ranges, tr.String())
	}
	return strings.Join(ranges, "\n")
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
		Hour:   uint(hour),
		Minute: uint(minute),
	}, nil
}

type Range struct {
	Start time.Time
	End   time.Time
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
