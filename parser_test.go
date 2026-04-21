package cron

import (
	"testing"
	"time"
)

func TestParse_Every(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		spec     string
		wantErr  bool
		expected time.Time
	}{
		{"@every 10s", false, now.Add(10 * time.Second)},
		{"@every 1m30s", false, now.Add(1*time.Minute + 30*time.Second)},
		{"@every 1h", false, now.Add(time.Hour)},
		{"@every 1d", true, time.Time{}}, // time.ParseDuration 不支持 d
		{"@every invalid", true, time.Time{}},
		{"@every", true, time.Time{}},
	}

	for _, tt := range tests {
		c, err := Parse(tt.spec)
		if err != nil {
			if !tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.spec, err, tt.wantErr)
			}
			continue
		}
		if !c.Next(now).Equal(tt.expected) {
			t.Errorf("Parse(%q) = %v, want %v", tt.spec, c.Next(now), tt.expected)
		}
	}
}

func TestParse_Daily(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		spec     string
		expected []time.Time
		wantErr  bool
	}{
		{
			spec: "@daily 5:00",
			expected: []time.Time{
				time.Date(2024, 1, 1, 5, 0, 0, 0, time.UTC),
			},
		},
		{
			spec: "@daily 0:00,12:00",
			expected: []time.Time{
				time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{"@daily 24:00", nil, true},
		{"@daily 5:60", nil, true},
		{"@daily", nil, true},
	}

	for _, tt := range tests {
		sch, err := Parse(tt.spec)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse(%q) error = %v, wantErr %v", tt.spec, err, tt.wantErr)
			continue
		}
		if tt.wantErr {
			continue
		}

		// 对于 MultipleSchedule 或单个 Schedule，验证 Next
		// 注意：Next(t) 返回的是 t 之后的时间
		checkTime := now
		for _, exp := range tt.expected {
			next := sch.Next(checkTime)
			if !next.Equal(exp) {
				t.Errorf("%s: expected %v, got %v", tt.spec, exp, next)
			}
			checkTime = next
		}
	}
}

func TestParse_Weekly(t *testing.T) {
	// 2024-01-01 是周一 (Monday = 1 in our parser's context of standardParser)
	// robfig/cron: Sunday = 0 or 7
	now := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC) // 2023-12-31 Sunday

	tests := []struct {
		spec     string
		expected time.Time
		wantErr  bool
	}{
		{"@weekly 1 5:00", time.Date(2024, 1, 1, 5, 0, 0, 0, time.UTC), false},
		{"@weekly 0 5:00", time.Date(2023, 12, 31, 5, 0, 0, 0, time.UTC), false}, // Sunday
		{"@weekly 1-5 9:00", time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), false},
		{"@weekly 1,3 5:00", time.Date(2024, 1, 1, 5, 0, 0, 0, time.UTC), false},
		{"@weekly invalid 5:00", time.Time{}, true},
		{"@weekly 1 invalid", time.Time{}, true},
	}

	for _, tt := range tests {
		sch, err := Parse(tt.spec)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse(%q) error = %v, wantErr %v", tt.spec, err, tt.wantErr)
			continue
		}
		if tt.wantErr {
			continue
		}
		next := sch.Next(now)
		if !next.Equal(tt.expected) {
			t.Errorf("%s: expected %v, got %v", tt.spec, tt.expected, next)
		}
	}
}

func TestParse_Monthly(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		spec     string
		expected time.Time
		wantErr  bool
	}{
		{"@monthly 1 5:00", time.Date(2024, 1, 1, 5, 0, 0, 0, time.UTC), false},
		{"@monthly L 5:00", time.Date(2024, 1, 31, 5, 0, 0, 0, time.UTC), false},
		{"@monthly L-1 5:00", time.Date(2024, 1, 30, 5, 0, 0, 0, time.UTC), false},
		{"@monthly 1-5 5:00", time.Date(2024, 1, 1, 5, 0, 0, 0, time.UTC), false},
		{"@monthly 1,15 5:00", time.Date(2024, 1, 1, 5, 0, 0, 0, time.UTC), false},
		{"@monthly 1 5:00,23:00", time.Date(2024, 1, 1, 5, 0, 0, 0, time.UTC), false},
		{"@monthly 32 5:00", time.Time{}, true},
		{"@monthly LX 5:00", time.Time{}, true},
		{"@monthly L-31 5:00", time.Time{}, true},
		{"@monthly L 0:00", time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC), false},
	}

	for _, tt := range tests {
		sch, err := Parse(tt.spec)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse(%q) error = %v, wantErr %v", tt.spec, err, tt.wantErr)
			continue
		}
		if tt.wantErr {
			continue
		}
		next := sch.Next(now)
		if !next.Equal(tt.expected) {
			t.Errorf("%s: expected %v, got %v", tt.spec, tt.expected, next)
		}
	}
}

func TestParse_Standard(t *testing.T) {
	spec := "0 0 5 * * *"
	sch, err := Parse(spec)
	if err != nil {
		t.Fatalf("Parse standard cron error: %v", err)
	}
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	next := sch.Next(now)
	expected := time.Date(2024, 1, 1, 5, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("Standard cron: expected %v, got %v", expected, next)
	}
}

func TestParse_EmptyAndInvalid(t *testing.T) {
	if _, err := Parse(""); err == nil {
		t.Error("Expected error for empty spec, got nil")
	}
	if _, err := Parse("   "); err == nil {
		t.Error("Expected error for whitespace spec, got nil")
	}
}
