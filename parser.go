package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var standardParser = cron.NewParser(
	cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
)

// Parse parses the spec string
func Parse(spec string) (Schedule, error) {
	if !strings.HasPrefix(spec, DescriptionPrefix) {
		return standardParser.Parse(spec)
	}

	parts := strings.Fields(spec)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty spec")
	}

	switch parts[0] {
	case Every:
		return parseEvery(parts)
	case Daily:
		return parseDaily(parts)
	case Weekly:
		return parseWeekly(parts)
	case Monthly:
		return parseMonthly(parts)
	default:
		return standardParser.Parse(spec)
	}
}

func parseEvery(parts []string) (Schedule, error) {
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid @every spec: %v", parts)
	}
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return nil, err
	}
	return cron.Every(duration), nil
}

func parseDaily(parts []string) (Schedule, error) {
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid @daily spec: %v", parts)
	}
	// Supports @daily 5:00 or @daily 0:00,12:00
	times := strings.Split(parts[1], AndChar)
	var schedules []Schedule
	for _, t := range times {
		hour, minute, err := parseTime(t)
		if err != nil {
			return nil, err
		}
		// Construct standard cron format: 0 minute hour * * *
		spec := fmt.Sprintf("0 %d %d * * *", minute, hour)
		sch, err := standardParser.Parse(spec)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, sch)
	}
	if len(schedules) == 1 {
		return schedules[0], nil
	}
	return NewMultipleSchedule(schedules...)
}

func parseWeekly(parts []string) (Schedule, error) {
	// @weekly 1 5:00 -> Monday 5:00
	// @weekly 1-5 5:00 -> Mon-Fri 5:00
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid @weekly spec: %v", parts)
	}
	days := parts[1]
	timeStr := parts[2]

	times := strings.Split(timeStr, AndChar)
	var schedules []Schedule
	for _, t := range times {
		hour, minute, err := parseTime(t)
		if err != nil {
			return nil, err
		}
		// Construct standard cron format: 0 minute hour * * days
		// Note: robfig/cron Sunday is 0 (or 7)
		spec := fmt.Sprintf("0 %d %d * * %s", minute, hour, days)
		sch, err := standardParser.Parse(spec)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, sch)
	}
	if len(schedules) == 1 {
		return schedules[0], nil
	}
	return NewMultipleSchedule(schedules...)
}

func parseMonthly(parts []string) (Schedule, error) {
	// @monthly 1 5:00
	// @monthly L 5:00
	// @monthly 15,L 23:00

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid @monthly spec fields: expect 3, got %d", len(parts))
	}
	doms := strings.Split(parts[1], AndChar)

	ts := make([][2]int, 0)
	times := strings.Split(parts[2], AndChar)
	for _, t := range times {
		hour, minute, err := parseTime(t)
		if err != nil {
			return nil, err
		}
		ts = append(ts, [2]int{hour, minute})
	}

	var schedules []Schedule
	for _, dom := range doms {
		if strings.HasPrefix(dom, LastNDomPrefix) {
			// L
			// L-1
			n := 0
			if len(dom) > 1 {
				// L-1
				if dom[1] != '-' {
					return nil, fmt.Errorf("invalid L format: %s", dom)
				}
				val, err := strconv.Atoi(dom[2:])
				if err != nil {
					return nil, fmt.Errorf("invalid L value: %s", dom)
				}
				n = val
			}
			for _, t := range ts {
				c, err := NewLastNDomSchedule(n, t[0], t[1])
				if err != nil {
					return nil, err
				}
				schedules = append(schedules, c)
			}
		} else {
			for _, t := range ts {
				spec := fmt.Sprintf("0 %d %d %s * *", t[1], t[0], dom)
				c, err := standardParser.Parse(spec)
				if err != nil {
					return nil, err
				}
				schedules = append(schedules, c)
			}
		}
	}

	if len(schedules) == 1 {
		return schedules[0], nil
	}
	return NewMultipleSchedule(schedules...)
}

func parseTime(t string) (hour, minute int, err error) {
	parts := strings.Split(t, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time format: %s", t)
	}
	hour, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	if hour < 0 || hour > 23 {
		return 0, 0, fmt.Errorf("invalid hour: %d", hour)
	}

	minute, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	if minute < 0 || minute > 59 {
		return 0, 0, fmt.Errorf("invalid minute: %d", minute)
	}

	return hour, minute, nil
}
