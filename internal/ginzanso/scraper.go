package scraper

import (
	"cloud.google.com/go/civil"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/sjunepark/ryokan/internal/scraper"
	"github.com/sjunepark/ryokan/internal/yearmonth"
	"strings"
	"time"
)

func Scrape(from civil.Date, to civil.Date) (availableDates []civil.Date, err error) {
	browser := rod.New().Trace(true).SlowMotion(time.Second).NoDefaultDevice().MustConnect()
	defer func(browser *rod.Browser) {
		err := browser.Close()
		if err != nil {
			panic(err)
		}
	}(browser)
	page := browser.MustPage("https://reserve.489ban.net/client/ginzanso/4/plan/availability/daily?#content").MustWaitLoad().MustWindowFullscreen()

	s := scraper.NewScraper(from, to, 12)

	for s.KeepScraping() {
		availableYearMonths := getAvailableMonths(page)
		if len(availableYearMonths) == 0 {
			return []civil.Date{}, fmt.Errorf("no available year months")
		}

		for _, currentYearMonth := range availableYearMonths {
			s.AddScrapedYearMoth(currentYearMonth)
			if !s.ShouldScrapeYearMonth(currentYearMonth) {
				continue
			}
			availableDates = append(availableDates, getAvailableDates(page)...)
		}
		// todo: Go to next
	}

	time.Sleep(time.Hour)
	return availableDates, nil
}

func getAvailableDates(page *rod.Page) []civil.Date {
	availableElements := getAvailableDateElements(page)

	var dates []civil.Date
	for _, availableElement := range availableElements {
		parent, err := getFirstPtagParent(availableElement)
		if err != nil {
			fmt.Println("Error getting parent p tag:", err)
			continue
		}
		parentClassesStr, err := parent.Attribute("class")
		if err != nil {
			fmt.Println("No class:", err)
			continue
		}
		parentClasses := parseClasses(*parentClassesStr)

		dates = append(dates, getDatesFromClasses(parentClasses)...)
	}

	return dates
}

func getAvailableDateElements(page *rod.Page) []*rod.Element {
	availableTriangles, triangleErr := page.Elements(".fa-exclamation-triangle")
	availableCircles, circleErr := page.Element(".fa-circle-o")
	if triangleErr != nil && circleErr != nil {
		return []*rod.Element{}
	}
	return append(availableTriangles, availableCircles)
}

func getFirstPtagParent(element *rod.Element) (*rod.Element, error) {
	current := element
	for current != nil {
		var parentErr error
		current, parentErr = current.Parent()
		if parentErr != nil {
			return nil, parentErr
		}

		tagName := current.MustEval(`() => this.tagName`)
		if tagName.Str() == "P" {
			return current, nil
		}
	}
	return nil, fmt.Errorf("no parent p tag found")
}

func parseClasses(classes string) []string {
	return strings.Fields(classes)
}

func getDatesFromClasses(classes []string) []civil.Date {
	var dates []civil.Date
	for _, class := range classes {
		date, _ := parseDate(class)
		dates = append(dates, date)
	}
	return dates
}

func parseDate(date string) (civil.Date, error) {
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return civil.Date{}, err
	}
	return civil.DateOf(dateTime), nil
}

func getAvailableMonths(page *rod.Page) []yearmonth.YearMonth {
	var errors []string
	var yearMonthElements []*rod.Element

	for i := 0; i < 3; i++ {
		selector := fmt.Sprintf("#yearMonth_%d", i)
		yearMonths, err := page.Elements(selector)
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}
		yearMonthElements = append(yearMonthElements, yearMonths...)
	}

	if len(yearMonthElements) == 0 {
		return []yearmonth.YearMonth{}
	}

	yearMonths := make([]yearmonth.YearMonth, 0, len(yearMonthElements))
	for _, yme := range yearMonthElements {
		yearMonthText, err := yme.Text()
		if err != nil {
			fmt.Println("Error getting yearMonth text:", err)
			continue
		}
		yearMonthTime, err := time.Parse("Jan 2006", yearMonthText)
		if err != nil {
			fmt.Println("Error parsing yearMonth:", err)
			continue
		}
		yearMonth := yearmonth.NewYearMonth(yearMonthTime.Year(), yearMonthTime.Month())
		yearMonths = append(yearMonths, yearMonth)
	}
	return yearMonths
}
