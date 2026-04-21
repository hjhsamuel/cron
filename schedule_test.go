package cron

import (
	"testing"
	"time"
)

func TestLastNDomSchedule(t *testing.T) {
	// 模拟 2024年2月 (闰年，29天)
	now := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	// L (最后一天 29号 5:00)
	sch, err := Parse("@monthly L 5:00")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	next := sch.Next(now)
	expected := time.Date(2024, 2, 29, 5, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("L: expected %v, got %v", expected, next)
	}

	// L-1 (倒数第二天 28号 5:00)
	sch, _ = Parse("@monthly L-1 5:00")
	next = sch.Next(now)
	expected = time.Date(2024, 2, 28, 5, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("L-1: expected %v, got %v", expected, next)
	}

	// 跨月测试
	// 2024年2月29号 6:00 之后，应该是 2024年3月31号 (3月有31天，L 是 31号)
	now = time.Date(2024, 2, 29, 6, 0, 0, 0, time.UTC)
	sch, _ = Parse("@monthly L 5:00")
	next = sch.Next(now)
	expected = time.Date(2024, 3, 31, 5, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("L across month: expected %v, got %v", expected, next)
	}

	// 测试 targetDay < 1 的情况 (例如在 2月份请求 L-30)
	sch, _ = Parse("@monthly L-30 5:00")
	now = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	next = sch.Next(now)
	expected = time.Date(2024, 2, 1, 5, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("L-30 (targetDay < 1): expected %v, got %v", expected, next)
	}

	// 测试 @monthly 15,L 23:00
	sch, _ = Parse("@monthly 15,L 23:00")
	now = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	// 第一次应该在 15号
	next = sch.Next(now)
	expected = time.Date(2024, 2, 15, 23, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("@monthly 15,L: expected %v, got %v", expected, next)
	}

	// 15号之后应该在 29号
	next = sch.Next(next)
	expected = time.Date(2024, 2, 29, 23, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("@monthly 15,L: expected %v, got %v", expected, next)
	}
}

func TestNewLastNDomSchedule_Validation(t *testing.T) {
	// 有效的参数
	_, err := NewLastNDomSchedule(0, 5, 0)
	if err != nil {
		t.Errorf("Expected valid schedule, got error: %v", err)
	}

	// 无效的 n
	_, err = NewLastNDomSchedule(31, 5, 0)
	if err == nil {
		t.Error("Expected error for n=31, got nil")
	}

	// 无效的 hour
	_, err = NewLastNDomSchedule(0, 24, 0)
	if err == nil {
		t.Error("Expected error for hour=24, got nil")
	}
}
