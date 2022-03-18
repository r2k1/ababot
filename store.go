package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
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
