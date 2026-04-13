package ratelimit

import (
	"testing"
	"time"
)

func TestNew_InvalidMaxLines(t *testing.T) {
	_, err := New(Options{MaxLines: 0, Window: time.Second})
	if err == nil {
		t.Fatal("expected error for zero MaxLines")
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(Options{MaxLines: 10, Window: 0})
	if err == nil {
		t.Fatal("expected error for zero Window")
	}
}

func TestNew_WindowTooSmall(t *testing.T) {
	_, err := New(Options{MaxLines: 10, Window: 500 * time.Microsecond})
	if err == nil {
		t.Fatal("expected error for sub-millisecond window")
	}
}

func TestNew_ValidOptions(t *testing.T) {
	l, err := New(Options{MaxLines: 5, Window: time.Second})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil Limiter")
	}
}

func TestLimiter_Allow_UnderLimit(t *testing.T) {
	now := time.Now()
	l, _ := New(Options{MaxLines: 3, Window: time.Second, Now: func() time.Time { return now }})
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
}

func TestLimiter_Allow_ExceedsLimit(t *testing.T) {
	now := time.Now()
	l, _ := New(Options{MaxLines: 3, Window: time.Second, Now: func() time.Time { return now }})
	for i := 0; i < 3; i++ {
		l.Allow()
	}
	if l.Allow() {
		t.Fatal("expected Allow()=false after limit exceeded")
	}
}

func TestLimiter_Allow_ResetsAfterWindow(t *testing.T) {
	now := time.Now()
	l, _ := New(Options{MaxLines: 2, Window: time.Second, Now: func() time.Time { return now }})
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("expected false at limit")
	}
	// Advance past the window.
	now = now.Add(2 * time.Second)
	if !l.Allow() {
		t.Fatal("expected true after window reset")
	}
}

func TestLimiter_Apply_DropsExcess(t *testing.T) {
	now := time.Now()
	l, _ := New(Options{MaxLines: 3, Window: time.Second, Now: func() time.Time { return now }})
	input := []string{"a", "b", "c", "d", "e"}
	out := l.Apply(input)
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
}

func TestFormatSummary(t *testing.T) {
	s := FormatSummary(100, time.Second)
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}
