package scraper

import "cloud.google.com/go/civil"

type Scraper interface {
	NewScraper(from civil.Date, to civil.Date, attemptCount int) *Scraper
	GetAvailableDates() []civil.Date
}
