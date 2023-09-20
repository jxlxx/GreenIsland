package subjects

type Subject string

const (
	TickHour    Subject = "event.time.tick.hour"
	TickDay     Subject = "event.time.tick.day"
	TickQuarter Subject = "event.time.tick.quarter"
)

func (s Subject) String() string {
	return string(s)
}
