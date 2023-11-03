package date

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
	return NewYearMonth(ym.year, ym.month+time.Month(months))
}
