package scrapedata

import (
	"cloud.google.com/go/civil"
	"github.com/deckarep/golang-set/v2"
	"github.com/sjunepark/ryokan/internal/yearmonth"
)

type ScrapeData struct {
	from              civil.Date
	to                civil.Date
	firstMonth        yearmonth.YearMonth
	lastMonth         yearmonth.YearMonth
	scrapedYearMonths mapset.Set[yearmonth.YearMonth]
	availableDates    mapset.Set[civil.Date]
}

func NewScrapeData(from civil.Date, to civil.Date) *ScrapeData {
	return &ScrapeData{
		from:              from,
		to:                to,
		firstMonth:        yearmonth.NewYearMonth(from.Year, from.Month),
		lastMonth:         yearmonth.NewYearMonth(to.Year, to.Month),
		scrapedYearMonths: mapset.NewSet[yearmonth.YearMonth](),
		availableDates:    mapset.NewSet[civil.Date](),
	}
}

func (sd *ScrapeData) IsDateInRange(date civil.Date) bool {
	return (!date.Before(sd.from)) && (!date.After(sd.to))
}

func (sd *ScrapeData) IsYearMonthInRange(yearMonth yearmonth.YearMonth) bool {
	return (!yearMonth.Before(sd.firstMonth)) && (!yearMonth.After(sd.lastMonth))
}

func (sd *ScrapeData) AreAllMonthsFuture(yearMonths []yearmonth.YearMonth) bool {
	for _, ym := range yearMonths {
		if !ym.After(sd.lastMonth) {
			return false
		}
	}
	return true
}

func (sd *ScrapeData) AreAllMonthsScraped() bool {
	currentYearMonth := sd.firstMonth

	for currentYearMonth.Before(sd.lastMonth.AddMonths(1)) {
		if !sd.scrapedYearMonths.Contains(currentYearMonth) {
			return false
		}
		currentYearMonth = currentYearMonth.AddMonths(1)

	}
	return true
}
