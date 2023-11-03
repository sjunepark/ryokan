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

func Scrape(from civil.Date, to civil.Date) {
	browser := rod.New().Trace(true).SlowMotion(time.Second).NoDefaultDevice().MustConnect()
	defer func(browser *rod.Browser) {
		err := browser.Close()
		if err != nil {
			panic(err)
		}
	}(browser)

	page := browser.MustPage("https://reserve.489ban.net/client/ginzanso/4/plan/availability/daily?#content").MustWaitLoad().MustWindowFullscreen()

	sd := scraper.NewScraper(from, to)

	for shouldStop(sd.sc, to) {

	}

	availableTriangles, triangleErr := page.Elements(".fa-exclamation-triangle")
	availableCircles, circleErr := page.Element(".fa-circle-o")
	if triangleErr != nil && circleErr != nil {
		fmt.Println("No available rooms")
		return
	}
	availableElements := append(availableTriangles, availableCircles)

	var dates []time.Time
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

	time.Sleep(time.Hour)
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

func getDatesFromClasses(classes []string) []time.Time {
	var dates []time.Time
	for _, class := range classes {
		date, _ := parseDate(class)
		dates = append(dates, date)
	}
	return dates
}

func parseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

func getAvailableMonths(page *rod.Page) ([]yearmonth.YearMonth, error) {
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
		combinedError := fmt.Errorf(strings.Join(errors, "; "))
		fmt.Println(nil, combinedError)
		return nil, combinedError
	}

	var yearMonths []yearmonth.YearMonth
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
	return yearMonths, nil
}

// shouldStop stops when i) all months are after to or ii) all months are found
// availableMonths is the time.Time version of yearMonth, which is the first day of the month
func shouldStop(availableMonths []yearmonth.YearMonth, to civil.Date) bool {
	for _, availableMonth := range availableMonths {
		if !availableMonth.After(yearmonth.NewYearMonth(to.Year, to.Month)) {
			return false
		}
	}
	return true
}
