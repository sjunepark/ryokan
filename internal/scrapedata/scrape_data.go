package scrapedata

import (
	"github.com/deckarep/golang-set/v2"
	"time"
)

type ScrapeData struct {
	from           time.Time
	to             time.Time
	firstMonth     time.Time
	lastMonth      time.Time
	scrapedMonths  mapset.Set[time.Time]
	availableDates mapset.Set[time.Time]
}

func NewScrapeData(from time.Time, to time.Time) *ScrapeData {
	return &ScrapeData{
		from:           from,
		to:             to,
		firstMonth:     time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC),
		lastMonth:      time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, time.UTC),
		scrapedMonths:  mapset.NewSet[time.Time](),
		availableDates: mapset.NewSet[time.Time](),
	}
}

func (sd *ScrapeData) IsDateInRange(date time.Time) bool {
	return (!date.Before(sd.from)) && (!date.After(sd.to))
}

func (sd *ScrapeData) IsMonthInRange(month time.Time) bool {
	return (!month.Before(sd.firstMonth)) && (!month.After(sd.lastMonth))
}

func (sd *ScrapeData) AreAllMonthsFuture(months []time.Time) bool {
	for _, month := range months {
		if !month.After(sd.lastMonth) {
			return false
		}
	}
	return true
}

func (sd *ScrapeData) AreAllMonthsScraped() bool {
	currentMonth := sd.firstMonth

	for currentMonth.Before(sd.lastMonth.AddDate(0, 1, 0)) {
		if !sd.scrapedMonths.Contains(currentMonth) {
			return false
		}
		currentMonth = currentMonth.AddDate(0, 1, 0)
	}
	return true
}
