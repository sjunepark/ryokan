package scraper

import (
	"cloud.google.com/go/civil"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/sjunepark/ryokan/internal/yearmonth"
	"testing"
	"time"
)

func createScrapeData(t testing.TB, from string, to string) *ScrapedData {
	t.Helper()
	fromTime, err := time.Parse("2006-01-02", from)
	if err != nil {
		t.Fatalf("Error parsing from yearmonth: %v", err)
	}
	fromDate := civil.DateOf(fromTime)
	toTime, err := time.Parse("2006-01-02", to)
	if err != nil {
		t.Fatalf("Error parsing to yearmonth: %v", err)
	}
	toDate := civil.DateOf(toTime)

	sd := NewScraper(fromDate, toDate)
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
			t.Errorf("Expected %v for yearmonth %v, but got %v", tc.expected, tc.date, result)
		}
	}
}

func TestIsMonthInRange(t *testing.T) {
	sd := createScrapeData(t, "2022-12-25", "2023-02-28")

	testCases := []struct {
		yearMonth yearmonth.YearMonth
		expected  bool
	}{
		{yearmonth.NewYearMonth(2022, 12), true},
		{yearmonth.NewYearMonth(2023, 1), true},
		{yearmonth.NewYearMonth(2023, 2), true},
		{yearmonth.NewYearMonth(2023, 3), false},
		{yearmonth.NewYearMonth(2023, 4), false},
		{yearmonth.NewYearMonth(2023, 12), false},
		{yearmonth.NewYearMonth(2024, 1), false},
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
	s := NewScraper(civil.DateOf(from), civil.DateOf(to))

	testCases := []struct {
		yearMonth []yearmonth.YearMonth
		expected  bool
	}{
		{
			[]yearmonth.YearMonth{
				yearmonth.NewYearMonth(2022, 11),
				yearmonth.NewYearMonth(2022, 12),
				yearmonth.NewYearMonth(2023, 1),
			},
			false,
		},
		{
			[]yearmonth.YearMonth{
				yearmonth.NewYearMonth(2023, 1),
				yearmonth.NewYearMonth(2023, 2),
				yearmonth.NewYearMonth(2023, 3),
			},
			false,
		},
		{
			[]yearmonth.YearMonth{
				yearmonth.NewYearMonth(2023, 3),
				yearmonth.NewYearMonth(2023, 4),
				yearmonth.NewYearMonth(2023, 5),
			},
			true,
		},
		{
			[]yearmonth.YearMonth{
				yearmonth.NewYearMonth(2024, 1),
				yearmonth.NewYearMonth(2024, 2),
				yearmonth.NewYearMonth(2024, 3),
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
	s := NewScraper(civil.DateOf(from), civil.DateOf(to))

	// Scenario 1: No months scraped
	if result := s.AreAllMonthsScraped(); result {
		t.Errorf("Expected false for AreAllMonthsScraped, but got %v", result)
	}

	// Scenario 2: Some months scraped
	s.scrapedYearMonths.Add(yearmonth.NewYearMonth(2023, 1))
	if result := s.AreAllMonthsScraped(); result {
		t.Errorf("Expected false for AreAllMonthsScraped, but got %v", result)
	}

	// Scenario 3: All months scraped
	s.scrapedYearMonths.Add(yearmonth.NewYearMonth(2023, 2))
	s.scrapedYearMonths.Add(yearmonth.NewYearMonth(2023, 3))

	if result := s.AreAllMonthsScraped(); !result {
		t.Errorf("Expected true for AreAllMonthsScraped, but got %v", result)
	}
}
