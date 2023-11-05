package scraper

import (
	"cloud.google.com/go/civil"
	"github.com/deckarep/golang-set/v2"
	"github.com/sjunepark/ryokan/internal/yearmonth"
)

type Scraper struct {
	from               civil.Date
	to                 civil.Date
	fromYM             yearmonth.YearMonth
	toYM               yearmonth.YearMonth
	yearMonthsToScrape mapset.Set[yearmonth.YearMonth]
	availableDates     mapset.Set[civil.Date]
	attempts           int
	maxAttempts        int
}

func NewScraper(from civil.Date, to civil.Date, maxAttempts int) *Scraper {
	fromYM := yearmonth.NewYearMonth(from.Year, from.Month)
	toYM := yearmonth.NewYearMonth(to.Year, to.Month)

	yearMonthsToScrape := mapset.NewSet[yearmonth.YearMonth]()
	for currentYearMonth := fromYM; currentYearMonth.Before(toYM); currentYearMonth = currentYearMonth.AddMonths(1) {
		yearMonthsToScrape.Add(currentYearMonth)
	}

	return &Scraper{
		from:               from,
		to:                 to,
		fromYM:             yearmonth.NewYearMonth(from.Year, from.Month),
		toYM:               yearmonth.NewYearMonth(to.Year, to.Month),
		yearMonthsToScrape: yearMonthsToScrape,
		availableDates:     mapset.NewSet[civil.Date](),
		attempts:           0,
		maxAttempts:        maxAttempts,
	}
}

func (s *Scraper) ShouldScrapeYearMonth(yearMonth yearmonth.YearMonth) bool {
	if !yearMonth.Before(s.toYM) {
		return false
	}

	return s.yearMonthsToScrape.Contains(yearMonth)
}

func (s *Scraper) AddScrapedYearMoth(yearMonth yearmonth.YearMonth) {
	s.yearMonthsToScrape.Remove(data.yearMonth)
	s.availableDates.Union(data.availableDates)
}

func (s *Scraper) AreAllMonthsScraped() bool {
	return s.yearMonthsToScrape.Cardinality() == 0
}

func (s *Scraper) KeepScraping() bool {
	return s.attempts < s.maxAttempts && !s.AreAllMonthsScraped()
}

type ScrapedYearMonthData struct {
	yearMonth      yearmonth.YearMonth
	availableDates mapset.Set[civil.Date]
}
