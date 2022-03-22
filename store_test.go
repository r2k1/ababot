package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestParseTimeRange(t *testing.T) {
	t.Run("without duration", func(t *testing.T) {
		tr, err := ParseTimeRange("Mon 16:00")
		require.NoError(t, err)
		assert.Equal(t, Subscription{
			Weekday: time.Monday,
			Time: Clock{
				Hour:   16,
				Minute: 0,
			},
			Hours: 1,
		}, tr)
		assert.Equal(t, "Mon 16:00-17:00", tr.String())
	})

	t.Run("with duration", func(t *testing.T) {
		tr, err := ParseTimeRange("Mon 16:00 2")
		require.NoError(t, err)
		assert.Equal(t, Subscription{
			Weekday: time.Monday,
			Time: Clock{
				Hour:   16,
				Minute: 0,
			},
			Hours: 2,
		}, tr)
		assert.Equal(t, "Mon 16:00-18:00", tr.String())
	})

}

func Test_NewCalendar(t *testing.T) {
	calendar := NewCalendar(time.Date(2020, 1, 1, 6, 0, 0, 0, time.Local), time.Date(2020, 1, 2, 18, 0, 0, 0, time.Local))
	fmt.Println(calendar)
}

func TestCalendar_Available(t *testing.T) {
	data, err := fetchData()
	require.NoError(t, err)
	cal := NewCalendar(time.Now(), time.Now().Add(time.Hour*24*7))
	for _, b := range data {
		cal.Book(b)
	}
	fmt.Println(cal.NonZero().String())
}

func TestCalendar_Book(t *testing.T) {
	cal := Calendar{
		time.Date(2020, 1, 1, 6, 0, 0, 0, time.Local): 1,
		time.Date(2020, 1, 1, 7, 0, 0, 0, time.Local): 1,
		time.Date(2020, 1, 1, 8, 0, 0, 0, time.Local): 1,
		time.Date(2020, 1, 1, 9, 0, 0, 0, time.Local): 1,
	}
	cal.Book(Booking{
		Start: time.Date(2020, 1, 1, 7, 30, 0, 0, time.Local),
		End:   time.Date(2020, 1, 1, 8, 30, 0, 0, time.Local),
	})
	assert.Equal(t, Calendar{
		time.Date(2020, 1, 1, 6, 0, 0, 0, time.Local): 1,
		time.Date(2020, 1, 1, 7, 0, 0, 0, time.Local): 0,
		time.Date(2020, 1, 1, 8, 0, 0, 0, time.Local): 0,
		time.Date(2020, 1, 1, 9, 0, 0, 0, time.Local): 1,
	}, cal)
}

func TestCalendar_ForSubscription(t *testing.T) {
	// 1-1-2020 is Wednesday
	cal := Calendar{
		time.Date(2020, 1, 1, 6, 0, 0, 0, time.Local): 1,
		time.Date(2020, 1, 1, 7, 0, 0, 0, time.Local): 2,
		time.Date(2020, 1, 1, 8, 0, 0, 0, time.Local): 3,
		time.Date(2020, 1, 1, 9, 0, 0, 0, time.Local): 4,
	}
	t.Run("no match", func(t *testing.T) {
		assert.Equal(t, Calendar{}, cal.ForSubscription(Subscription{
			Weekday: time.Monday,
			Time: Clock{
				Hour:   6,
				Minute: 0,
			},
			Hours: 1,
		}))
	})

	t.Run("match", func(t *testing.T) {
		assert.Equal(t, Calendar{
			time.Date(2020, 1, 1, 6, 0, 0, 0, time.Local): 1,
		}, cal.ForSubscription(Subscription{
			Weekday: time.Wednesday,
			Time: Clock{
				Hour:   6,
				Minute: 0,
			},
			Hours: 1,
		}))
	})

	t.Run("double booking match", func(t *testing.T) {
		assert.Equal(t, Calendar{
			time.Date(2020, 1, 1, 6, 0, 0, 0, time.Local): 1,
			time.Date(2020, 1, 1, 7, 0, 0, 0, time.Local): 2,
		}, cal.ForSubscription(Subscription{
			Weekday: time.Wednesday,
			Time: Clock{
				Hour:   6,
				Minute: 0,
			},
			Hours: 2,
		}))
	})

	t.Run("no match for a double booking", func(t *testing.T) {
		assert.Equal(t, Calendar{}, cal.ForSubscription(Subscription{
			Weekday: time.Wednesday,
			Time: Clock{
				Hour:   8,
				Minute: 0,
			},
			Hours: 2,
		}))
	})

}
