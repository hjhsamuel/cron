package cron

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

type Schedule interface {
	cron.Schedule
}

type MultipleSchedule struct {
	sch []Schedule
}

func (s *MultipleSchedule) Next(t time.Time) time.Time {
	var next time.Time
	for _, item := range s.sch {
		tmp := item.Next(t)
		if tmp.IsZero() {
			continue
		}
		if next.IsZero() || tmp.Before(next) {
			next = tmp
		}
	}
	return next
}

func NewMultipleSchedule(cs ...Schedule) (*MultipleSchedule, error) {
	if len(cs) == 0 {
		return nil, fmt.Errorf("at least one cron schedule is required")
	}
	s := &MultipleSchedule{sch: cs}
	if s.Next(time.Now()).IsZero() {
		return nil, fmt.Errorf("invalid schedule")
	}
	return s, nil
}

type LastNDomSchedule struct {
	n      int
	hour   int
	minute int
}

func (s *LastNDomSchedule) Next(t time.Time) time.Time {
	year, month, _ := t.Date()
	loc := t.Location()

	for {
		lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, loc).Day()
		targetDay := lastDay - s.n
		if targetDay < 1 {
			// switch to first day
			targetDay = 1
		}
		next := time.Date(year, month, targetDay, s.hour, s.minute, 0, 0, loc)
		if next.After(t) {
			return next
		}

		month += 1
		if month > 12 {
			month = 1
			year += 1
		}
	}
}

func NewLastNDomSchedule(n, hour, minute int) (*LastNDomSchedule, error) {
	if n < 0 || n >= 31 {
		return nil, fmt.Errorf("invalid n: valid range is [0, 30], got: %d", n)
	}
	if hour < 0 || hour > 23 {
		return nil, fmt.Errorf("invalid hour: valid range is [0, 23], got: %d", hour)
	}
	if minute < 0 || minute > 59 {
		return nil, fmt.Errorf("invalid minute: valid range is [0, 59], got: %d", minute)
	}
	return &LastNDomSchedule{n: n, hour: hour, minute: minute}, nil
}
