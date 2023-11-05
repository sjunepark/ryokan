package yearmonth

import (
	"testing"
	"time"
)

func TestYearMonth_Before(t *testing.T) {
	tests := []struct {
		name       string
		thisYear   int
		thisMonth  time.Month
		otherYear  int
		otherMonth time.Month
		expected   bool
	}{
		{"before case", 2023, time.January, 2023, time.February, true},
		{"equal case", 2023, time.January, 2023, time.January, false},
		{"after case", 2023, time.February, 2023, time.January, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thisYM := NewYearMonth(tt.thisYear, tt.thisMonth)
			otherYM := NewYearMonth(tt.otherYear, tt.otherMonth)
			if got := thisYM.Before(otherYM); got != tt.expected {
				t.Errorf("YearMonth.Before() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestYearMonth_After(t *testing.T) {
	tests := []struct {
		name       string
		thisYear   int
		thisMonth  time.Month
		otherYear  int
		otherMonth time.Month
		expected   bool
	}{
		{"before case", 2023, time.January, 2023, time.February, false},
		{"equal case", 2023, time.January, 2023, time.January, false},
		{"after case", 2023, time.February, 2023, time.January, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thisYM := NewYearMonth(tt.thisYear, tt.thisMonth)
			otherYM := NewYearMonth(tt.otherYear, tt.otherMonth)
			if got := thisYM.After(otherYM); got != tt.expected {
				t.Errorf("YearMonth.After() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestYearMonth_AddMonths(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    time.Month
		add      int
		expected YearMonth
	}{
		{"add within year", 2023, time.January, 2, NewYearMonth(2023, time.March)},
		{"add cross year", 2023, time.December, 1, NewYearMonth(2024, time.January)},
		{"add negative", 2023, time.March, -2, NewYearMonth(2023, time.January)},
		{"add zero", 2023, time.January, 0, NewYearMonth(2023, time.January)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ym := NewYearMonth(tt.year, tt.month)
			result := ym.AddMonths(tt.add)
			if got := result; got != tt.expected {
				t.Errorf("YearMonth.AddMonths() = %v, want %v", got, tt.expected)
			}
		})
	}
}
