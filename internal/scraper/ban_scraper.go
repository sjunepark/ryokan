package scraper

import (
	"cloud.google.com/go/civil"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/sjunepark/ryokan/internal/yearmonth"
	"log/slog"
)

type BanScraper struct {
	client             string
	scrapeFrom         civil.Date
	scrapeTo           civil.Date
	yearMonthsToScrape mapset.Set[yearmonth.YearMonth]
	availableDates     mapset.Set[civil.Date]
	navigateCount      int
	maxNavigateCount   int
}

func NewBanScraper(client string, scrapeFrom civil.Date, scrapeTo civil.Date, maxNavigateCount int) (*BanScraper, error) {
	if scrapeFrom.After(scrapeTo) {
		return &BanScraper{}, fmt.Errorf("scrapeFrom must be before scrapeTo")
	}
	yearMonthsToScrape := mapset.NewSet[yearmonth.YearMonth]()
	for current := yearmonth.NewYearMonth(scrapeFrom.Year, scrapeFrom.Month); !current.After(yearmonth.NewYearMonth(scrapeTo.Year, scrapeTo.Month)); current = current.AddMonths(1) {
		yearMonthsToScrape.Add(current)
	}
	slog.Info("Created", "yearMonthsToScrape", yearMonthsToScrape)
	return &BanScraper{
			client:             client,
			scrapeFrom:         scrapeFrom,
			scrapeTo:           scrapeTo,
			yearMonthsToScrape: yearMonthsToScrape,
			availableDates:     mapset.NewSet[civil.Date](),
			navigateCount:      0,
			maxNavigateCount:   maxNavigateCount,
		},
		nil
}

func (s *BanScraper) GetAvailableDates() mapset.Set[civil.Date] {
	page, err := createRodPage(s.client)
	if err != nil {
		slog.Error("Error creating rod page", "error", err)
		return mapset.NewSet[civil.Date]()
	}

	for s.keepNavigating() {
		s.navigateCount++
		pageData, err := GetPageData(page)
		if err != nil {
			slog.Error("Error getting page data", "error", err)
			return mapset.NewSet[civil.Date]()
		}

		s.processPageData(pageData)

		err = goNext(page)
		if err != nil {
			continue
		}
	}

	return s.availableDates
}

func (s *BanScraper) processPageData(pageData *PageData) {
	pageData.yearMonths.Each(func(elem yearmonth.YearMonth) bool {
		if s.yearMonthsToScrape.Contains(elem) {
			s.yearMonthsToScrape.Remove(elem)
			slog.Info("Removed", "yearMonth", elem, "remaining", s.yearMonthsToScrape)
		}
		return false
	})
	s.availableDates = s.availableDates.Union(pageData.availableDates)
	slog.Info("Processed pageData", "pageData.availableDates", pageData.availableDates, "total availableDates", s.availableDates)
}

func createRodPage(client string) (*rod.Page, error) {
	url := fmt.Sprintf("https://reserve.489ban.net/client/%s/4/plan/availability/daily?#content", client)
	var page *rod.Page
	err := rod.Try(func() {
		page = rod.New().MustConnect().MustPage(url).MustWaitLoad().MustWindowFullscreen()
	})
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (s *BanScraper) keepNavigating() bool {
	countWithinLimit := s.navigateCount < s.maxNavigateCount
	yearMonthsToScrapeNotEmpty := s.yearMonthsToScrape.Cardinality() > 0
	return countWithinLimit && yearMonthsToScrapeNotEmpty
}

func goNext(page *rod.Page) error {
	nextElement, err := page.Element("#next")
	if err != nil {
		return err
	}
	err = nextElement.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return err
	}
	return nil
}
