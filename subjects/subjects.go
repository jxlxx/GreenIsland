package subjects

import "fmt"

type Subject string

const (
	TickHour    Subject = "event.time.new.hour"
	TickDay     Subject = "event.time.new.day"
	TickQuarter Subject = "event.time.new.quarter"

	quarterlyCountryUpdate Subject = "news.country.%s.Q%d"
	quarterlyCompanyUpdate Subject = "news.company.%s.Q%d"
)

func (s Subject) String() string {
	return string(s)
}

func DailyTick(quarter, day, hour int) string {
	return fmt.Sprintf("event.time.Q%d.D%d.H%d", quarter, day, hour)
}

func QuarterlyCountryUpdate(code string, quarter int) string {
	return fmt.Sprintf(quarterlyCountryUpdate.String(), code, quarter)
}

func QuarterlyCompanyUpdate(code string, quarter int) string {
	return fmt.Sprintf(quarterlyCompanyUpdate.String(), code, quarter)
}
