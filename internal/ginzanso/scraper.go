package scraper

import (
	"fmt"
	"github.com/go-rod/rod"
	"strings"
	"time"
)

func Scrape() {
	browser := rod.New().Trace(true).SlowMotion(time.Second).NoDefaultDevice().MustConnect()
	defer func(browser *rod.Browser) {
		err := browser.Close()
		if err != nil {
			panic(err)
		}
	}(browser)

	page := browser.MustPage("https://reserve.489ban.net/client/ginzanso/4/plan/availability/daily?#content").MustWaitLoad().MustWindowFullscreen()

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

func getAvailableMonths(page *rod.Page) ([]time.Time, error) {
	month0, month0err := page.Elements("#yearMonth_0")
	month1, month1err := page.Elements("#yearMonth_1")
	month2, month2err := page.Elements("#yearMonth_2")

	if month0err != nil && month1err != nil && month2err != nil {
		var errors []string
		errors = append(errors, month0err.Error())
		errors = append(errors, month1err.Error())
		errors = append(errors, month2err.Error())
		return nil, fmt.Errorf(strings.Join(errors, "; "))
	}

	monthElements := append(month0, month1...)
	monthElements = append(monthElements, month2...)

	var months []time.Time
	for _, monthElement := range monthElements {
		monthText, err := monthElement.Text()
		if err != nil {
			fmt.Println("Error getting month text:", err)
			continue
		}
		month, err := time.Parse("Jan 2006", monthText)
		if err != nil {
			fmt.Println("Error parsing month:", err)
			continue
		}
		months = append(months, month)
	}
	return months, nil
}

// availableMonth is the time.Time version of yearMonth, which is the first day of the month
func shouldStop(availableMonths []time.Time, targetDate time.Time) bool {
	for _, availableMonth := range availableMonths {
		if !availableMonth.After(targetDate) {
			return false
		}
	}
	return true
}
