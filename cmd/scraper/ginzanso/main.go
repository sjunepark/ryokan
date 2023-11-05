package main

import (
	"github.com/sjunepark/ryokan/internal/scraper"
	"log/slog"
)

func main() {
	dates, err := scraper.ScrapeAvailableDates()
	if err != nil {
		slog.Error("Failed to scrape available dates", "error", err)
	}
	slog.Info("Available dates", "dates", dates)
}
