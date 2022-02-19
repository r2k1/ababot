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
	ID       string      `json:"id"`
	Times    []Weektime  `json:"times"`
	Notified []time.Time `json:"notified"`
}

type Weektime struct {
	Weekday time.Weekday `json:"weekday"`
	Hour    int          `json:"hours"`
	Minute  int          `json:"minutes"`
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

func (s *Store) addTime(userID string, time Weektime) {
	user, ok := s.Data.Users[userID]
	if !ok {
		user = &UserData{
			ID: userID,
		}
		s.Data.Users[userID] = user
	}
	user.Times = append(user.Times, time)
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
func (s *Store) AddTime2(userID string, timeS string) error {
	if len(userID) == 0 {
		return errors.New("userID can't be blank")
	}
	s.Lock()
	defer s.Unlock()
	data := strings.Split(strings.ToLower(strings.TrimSpace(timeS)), " ")
	if len(data) != 2 {
		return errors.New(`incorrect time format, expected format: "Mon 15:00"`)
	}
	weekday, ok := weekdayMapping[data[0]]
	if !ok {
		return fmt.Errorf("unknown day of the week: %s", weekday)
	}
	clock := strings.Split(data[1], ":")
	hour, err := strconv.ParseInt(clock[0], 10, 64)
	if err != nil {
		return fmt.Errorf("incorrect time format: \"%s\"", data[1])
	}
	if hour < 0 || hour > 23 {
		return errors.New("hour should be between 0 and 23")
	}
	minute, err := strconv.ParseInt(clock[1], 10, 64)
	if err != nil {
		return fmt.Errorf("incorrect time format: \"%s\"", data[1])
	}
	if minute < 0 || minute > 59 {
		return errors.New("minute should be between 0 and 59")
	}
	s.addTime(userID, Weektime{
		Weekday: weekday,
		Hour:    int(hour),
		Minute:  int(minute),
	})
	return s.save()
}

func (s *Store) save() error {
	if err := s.File.Truncate(0); err != nil {
		return err
	}
	if _, err := s.File.Seek(0, 0); err != nil {
		return err
	}
	return json.NewEncoder(s.File).Encode(s.Data)
}
