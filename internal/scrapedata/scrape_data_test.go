package scrapedata

import (
	"testing"
	"time"
)

func TestIsDateInRange(t *testing.T) {
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)
	s := NewScrapeData(from, to)

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
		if result := s.IsDateInRange(tc.date); result != tc.expected {
			t.Errorf("Expected %v for date %v, but got %v", tc.expected, tc.date, result)
		}
	}
}

func TestIsMonthInRange(t *testing.T) {
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)
	s := NewScrapeData(from, to)

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
		if result := s.IsMonthInRange(tc.month); result != tc.expected {
			t.Errorf("Expected %v for month %v, but got %v", tc.expected, tc.month, result)
		}
	}
}

func TestAreAllMonthsFuture(t *testing.T) {
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)
	s := NewScrapeData(from, to)

	testCases := []struct {
		months   []time.Time
		expected bool
	}{
		{
			[]time.Time{
				time.Date(2022, 11, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			[]time.Time{
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			[]time.Time{
				time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC),
			},
			true,
		},
		{
			[]time.Time{
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			},
			true,
		},
	}

	for _, tc := range testCases {
		if result := s.AreAllMonthsFuture(tc.months); result != tc.expected {
			t.Errorf("Expected %v for months %v, but got %v", tc.expected, tc.months, result)
		}
	}
}

func TestAreAllMonthsScraped(t *testing.T) {
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)
	s := NewScrapeData(from, to)

	// Scenario 1: No months scraped
	if result := s.AreAllMonthsScraped(); result {
		t.Errorf("Expected false for AreAllMonthsScraped, but got %v", result)
	}

	// Scenario 2: Some months scraped
	s.scrapedMonths.Add(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	if result := s.AreAllMonthsScraped(); result {
		t.Errorf("Expected false for AreAllMonthsScraped, but got %v", result)
	}

	// Scenario 3: All months scraped
	s.scrapedMonths.Add(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC))
	s.scrapedMonths.Add(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC))
	if result := s.AreAllMonthsScraped(); !result {
		t.Errorf("Expected true for AreAllMonthsScraped, but got %v", result)
	}
}
