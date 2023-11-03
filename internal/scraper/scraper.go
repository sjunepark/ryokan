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
}

func NewScraper(from civil.Date, to civil.Date) *Scraper {
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
	}
}

func (sd *Scraper) ShouldScrapeYearMonth(yearMonth yearmonth.YearMonth) bool {
	if !yearMonth.Before(sd.toYM) {
		return false
	}

	return sd.yearMonthsToScrape.Contains(yearMonth)
}

func (sd *Scraper) AddScrapedYearMoth(data ScrapedYearMonthData) {
	sd.yearMonthsToScrape.Remove(data.yearMonth)
	sd.availableDates.Union(data.availableDates)
}

func (sd *Scraper) AreAllMonthsScraped() bool {
	currentYearMonth := sd.fromYM
	for currentYearMonth.Before(sd.toYM.AddMonths(1)) {
		if !sd.scrapedYearMonths.Contains(currentYearMonth) {
			return false
		}
		currentYearMonth = currentYearMonth.AddMonths(1)
	}
	return true
}

func (sd *Scraper) IsDateInRange(date civil.Date) bool {
	return (!date.Before(sd.from)) && (!date.After(sd.to))
}

func (sd *Scraper) IsYearMonthInRange(yearMonth yearmonth.YearMonth) bool {
	return (!yearMonth.Before(sd.fromYM)) && (!yearMonth.After(sd.toYM))
}

func (sd *Scraper) AreAllMonthsFuture(yearMonths []yearmonth.YearMonth) bool {
	for _, ym := range yearMonths {
		if !ym.After(sd.toYM) {
			return false
		}
	}
	return true
}

type ScrapedYearMonthData struct {
	yearMonth      yearmonth.YearMonth
	availableDates mapset.Set[civil.Date]
}
