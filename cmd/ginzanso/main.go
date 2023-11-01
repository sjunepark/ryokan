package main

import (
	scraper "github.com/sjunepark/ryokan/internal/ginzanso"
	"time"
)

func main() {
	from := time.Date(2023, 12, 23, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC)
	scraper.Scrape(from, to)
}
