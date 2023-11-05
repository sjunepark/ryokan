package scraper

import (
	"cloud.google.com/go/civil"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-rod/rod"
	"github.com/sjunepark/ryokan/internal/yearmonth"
	"log/slog"
	"strings"
	"time"
)

type PageData struct {
	yearMonths     mapset.Set[yearmonth.YearMonth]
	availableDates mapset.Set[civil.Date]
}

func GetPageData(page *rod.Page) (*PageData, error) {
	yearMonths, err := getYearMonths(page)
	if err != nil {
		return &PageData{}, nil
	}
	availableDates := getAvailalbleDates(page)
	slog.Info("Got page data", "yearMonths", yearMonths, "availableDates", availableDates)
	return &PageData{yearMonths: yearMonths, availableDates: availableDates}, nil
}

func getYearMonths(page *rod.Page) (mapset.Set[yearmonth.YearMonth], error) {
	var yearMonths = mapset.NewSet[yearmonth.YearMonth]()
	elements, err := page.Elements(".webc_cal_head p")
	if err != nil {
		return mapset.NewSet[yearmonth.YearMonth](), fmt.Errorf("error getting year months: %w", err)
	}

	for _, element := range elements {
		elementText, err := element.Text()
		if err != nil {
			continue
		}
		yearMonthTime, err := time.Parse("Jan 2006", elementText)
		if err != nil {
			continue
		}
		yearMonth := yearmonth.NewYearMonth(yearMonthTime.Year(), yearMonthTime.Month())
		yearMonths.Add(yearMonth)
	}
	return yearMonths, nil
}

func getAvailalbleDates(page *rod.Page) mapset.Set[civil.Date] {
	availableDates := mapset.NewSet[civil.Date]()

	availableElements := rod.Elements{}
	triangles, triangleErr := page.Elements(".fa-exclamation-triangle")
	circles, circleErr := page.Elements(".fa-circle-o")
	if (triangleErr != nil) && (circleErr != nil) {
		return mapset.NewSet[civil.Date]()
	}
	availableElements = append(availableElements, triangles...)
	availableElements = append(availableElements, circles...)

	for _, element := range availableElements {
		parent, err := element.Parent()
		if err != nil {
			continue
		}
		grandParent, err := parent.Parent()
		if err != nil {
			continue
		}

		classesString, err := grandParent.Attribute("class")
		if err != nil {
			continue
		}

		classes := strings.Fields(*classesString)
		for _, class := range classes {
			// class is type of 2023-02-12. Parse this using time package
			datetime, err := time.Parse("2006-01-02", class)
			if err != nil {
				continue
			}
			date := civil.Date{Year: datetime.Year(), Month: datetime.Month(), Day: datetime.Day()}
			availableDates.Add(date)
		}
	}
	return availableDates
}
