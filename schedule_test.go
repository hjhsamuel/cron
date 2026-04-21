package cron

import (
	"testing"
	"time"
)

func TestLastNDomSchedule(t *testing.T) {
	// Mock Feb 2024 (Leap year, 29 days)
	now := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	// L (Last day: 29th 5:00)
	sch, err := Parse("@monthly L 5:00")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	next := sch.Next(now)
	expected := time.Date(2024, 2, 29, 5, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("L: expected %v, got %v", expected, next)
	}

	// L-1 (Second to last day: 28th 5:00)
	sch, _ = Parse("@monthly L-1 5:00")
	next = sch.Next(now)
	expected = time.Date(2024, 2, 28, 5, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("L-1: expected %v, got %v", expected, next)
	}

	// Cross-month test
	// After Feb 29, 2024 6:00, it should be Mar 31, 2024 (March has 31 days, L is 31st)
	now = time.Date(2024, 2, 29, 6, 0, 0, 0, time.UTC)
	sch, _ = Parse("@monthly L 5:00")
	next = sch.Next(now)
	expected = time.Date(2024, 3, 31, 5, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("L across month: expected %v, got %v", expected, next)
	}

	// Test case where targetDay < 1 (e.g., requesting L-30 in February)
	sch, _ = Parse("@monthly L-30 5:00")
	now = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	next = sch.Next(now)
	expected = time.Date(2024, 2, 1, 5, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("L-30 (targetDay < 1): expected %v, got %v", expected, next)
	}

	// Test @monthly 15,L 23:00
	sch, _ = Parse("@monthly 15,L 23:00")
	now = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	// First execution should be on the 15th
	next = sch.Next(now)
	expected = time.Date(2024, 2, 15, 23, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("@monthly 15,L: expected %v, got %v", expected, next)
	}

	// After 15th, it should be on the 29th
	next = sch.Next(next)
	expected = time.Date(2024, 2, 29, 23, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("@monthly 15,L: expected %v, got %v", expected, next)
	}
}

func TestNewLastNDomSchedule_Validation(t *testing.T) {
	// Valid parameters
	_, err := NewLastNDomSchedule(0, 5, 0)
	if err != nil {
		t.Errorf("Expected valid schedule, got error: %v", err)
	}

	// Invalid n
	_, err = NewLastNDomSchedule(31, 5, 0)
	if err == nil {
		t.Error("Expected error for n=31, got nil")
	}

	// Invalid hour
	_, err = NewLastNDomSchedule(0, 24, 0)
	if err == nil {
		t.Error("Expected error for hour=24, got nil")
	}
}
