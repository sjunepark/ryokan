package scraper

import (
	"testing"
	"time"
)

func TestDateInRange(t *testing.T) {
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)
	s := newScraper(from, to)

	testCases := []struct {
		date     time.Time
		expected bool
	}{
		{time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2022, 12, 31, 23, 59, 59, 0, time.UTC), false},
		{time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), true},
		{time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC), true},
		{time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC), true},
		{time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), false},
	}

	for _, tc := range testCases {
		if result := s.dateInRange(tc.date); result != tc.expected {
			t.Errorf("Expected %v for date %v, but got %v", tc.expected, tc.date, result)
		}
	}
}

func TestMonthInRange(t *testing.T) {
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)
	s := newScraper(from, to)

	testCases := []struct {
		month    time.Time
		expected bool
	}{
		{time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), true},
		{time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC), true},
		{time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC), false},
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), false},
	}

	for _, tc := range testCases {
		if result := s.monthInRange(tc.month); result != tc.expected {
			t.Errorf("Expected %v for month %v, but got %v", tc.expected, tc.month, result)
		}
	}
}
