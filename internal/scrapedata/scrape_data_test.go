package scrapedata

import (
	"cloud.google.com/go/civil"
	"github.com/sjunepark/ryokan/internal/date"
	"testing"
	"time"
)

func createScrapeData(t testing.TB, from string, to string) *ScrapeData {
	t.Helper()
	fromTime, err := time.Parse("2006-01-02", from)
	if err != nil {
		t.Fatalf("Error parsing from date: %v", err)
	}
	fromDate := civil.DateOf(fromTime)
	toTime, err := time.Parse("2006-01-02", to)
	if err != nil {
		t.Fatalf("Error parsing to date: %v", err)
	}
	toDate := civil.DateOf(toTime)

	sd := NewScrapeData(fromDate, toDate)
	return sd
}

func TestIsDateInRange(t *testing.T) {
	sd := createScrapeData(t, "2022-12-25", "2023-02-28")

	testCases := []struct {
		date     civil.Date
		expected bool
	}{
		{civil.DateOf(time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC)), false},
		{civil.DateOf(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC)), true},
		{civil.DateOf(time.Date(2022, 12, 31, 23, 59, 59, 0, time.UTC)), true},
		{civil.DateOf(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)), true},
		{civil.DateOf(time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)), true},
		{civil.DateOf(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)), false},
		{civil.DateOf(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)), false},
	}

	for _, tc := range testCases {
		if result := sd.IsDateInRange(tc.date); result != tc.expected {
			t.Errorf("Expected %v for date %v, but got %v", tc.expected, tc.date, result)
		}
	}
}

func TestIsMonthInRange(t *testing.T) {
	sd := createScrapeData(t, "2022-12-25", "2023-02-28")

	testCases := []struct {
		yearMonth date.YearMonth
		expected  bool
	}{
		{date.NewYearMonth(2022, 12), true},
		{date.NewYearMonth(2023, 1), true},
		{date.NewYearMonth(2023, 2), true},
		{date.NewYearMonth(2023, 3), false},
		{date.NewYearMonth(2023, 4), false},
		{date.NewYearMonth(2023, 12), false},
		{date.NewYearMonth(2024, 1), false},
	}

	for _, tc := range testCases {
		if result := sd.IsYearMonthInRange(tc.yearMonth); result != tc.expected {
			t.Errorf("Expected %v for month %v, but got %v", tc.expected, tc.yearMonth, result)
		}
	}
}

func TestAreAllMonthsFuture(t *testing.T) {
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)
	s := NewScrapeData(civil.DateOf(from), civil.DateOf(to))

	testCases := []struct {
		yearMonth []date.YearMonth
		expected  bool
	}{
		{
			[]date.YearMonth{
				date.NewYearMonth(2022, 11),
				date.NewYearMonth(2022, 12),
				date.NewYearMonth(2023, 1),
			},
			false,
		},
		{
			[]date.YearMonth{
				date.NewYearMonth(2023, 1),
				date.NewYearMonth(2023, 2),
				date.NewYearMonth(2023, 3),
			},
			false,
		},
		{
			[]date.YearMonth{
				date.NewYearMonth(2023, 3),
				date.NewYearMonth(2023, 4),
				date.NewYearMonth(2023, 5),
			},
			true,
		},
		{
			[]date.YearMonth{
				date.NewYearMonth(2024, 1),
				date.NewYearMonth(2024, 2),
				date.NewYearMonth(2024, 3),
			},
			true,
		},
	}

	for _, tc := range testCases {
		if result := s.AreAllMonthsFuture(tc.yearMonth); result != tc.expected {
			t.Errorf("Expected %v for months %v, but got %v", tc.expected, tc.yearMonth, result)
		}
	}
}

func TestAreAllMonthsScraped(t *testing.T) {
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)
	s := NewScrapeData(civil.DateOf(from), civil.DateOf(to))

	// Scenario 1: No months scraped
	if result := s.AreAllMonthsScraped(); result {
		t.Errorf("Expected false for AreAllMonthsScraped, but got %v", result)
	}

	// Scenario 2: Some months scraped
	s.scrapedYearMonths.Add(date.NewYearMonth(2023, 1))
	if result := s.AreAllMonthsScraped(); result {
		t.Errorf("Expected false for AreAllMonthsScraped, but got %v", result)
	}

	// Scenario 3: All months scraped
	s.scrapedYearMonths.Add(date.NewYearMonth(2023, 2))
	s.scrapedYearMonths.Add(date.NewYearMonth(2023, 3))

	if result := s.AreAllMonthsScraped(); !result {
		t.Errorf("Expected true for AreAllMonthsScraped, but got %v", result)
	}
}
