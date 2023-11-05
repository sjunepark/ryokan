package scraper

import (
	"cloud.google.com/go/civil"
	"github.com/sjunepark/ryokan/internal/scraper"
	"log/slog"
)

func Run() {
	scrapeFrom := civil.Date{Year: 2023, Month: 12, Day: 22}
	scrapeTo := civil.Date{Year: 2024, Month: 4, Day: 30}
	s, err := scraper.NewBanScraper("ginzanso", scrapeFrom, scrapeTo, 12)
	if err != nil {
		slog.Error("Error creating BanScraper", "error", err)
		return
	}

	availableDates := s.GetAvailableDates()
	slog.Info("Finished Scraping", "availableDates", availableDates)
}
