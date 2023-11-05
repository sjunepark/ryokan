package yearmonth

import (
	"cloud.google.com/go/civil"
	"fmt"
	"time"
)

type YearMonth struct {
	year  int
	month time.Month
	date  civil.Date
}

func NewYearMonth(year int, month time.Month) YearMonth {
	date := civil.Date{Year: year, Month: month, Day: 1}
	return YearMonth{year: year, month: month, date: date}
}

func (ym YearMonth) String() string {
	return fmt.Sprintf("%04d-%02d", ym.year, ym.month)
}

func (ym YearMonth) Before(other YearMonth) bool {
	return ym.date.Before(other.date)
}

func (ym YearMonth) After(other YearMonth) bool {
	return ym.date.After(other.date)
}

func (ym YearMonth) AddMonths(months int) YearMonth {
	currentTime := time.Date(ym.year, ym.month, 1, 0, 0, 0, 0, time.UTC)
	addedTime := currentTime.AddDate(0, months, 0)
	return NewYearMonth(addedTime.Year(), addedTime.Month())
}
