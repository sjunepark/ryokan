package scraper

import (
	"cloud.google.com/go/civil"
	"sort"
)

func ScrapeAvailableDates() ([]civil.Date, error) {
	scrapeFrom := civil.Date{Year: 2023, Month: 12, Day: 22}
	scrapeTo := civil.Date{Year: 2024, Month: 4, Day: 30}

	s, err := NewBanScraper("ginzanso", scrapeFrom, scrapeTo, 12)
	if err != nil {
		return []civil.Date{}, err
	}

	availableDates := s.GetAvailableDates().ToSlice()

	sort.Slice(availableDates, func(i, j int) bool {
		return availableDates[i].Before(availableDates[j])
	})

	return availableDates, nil
}
