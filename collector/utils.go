package collector

import "time"

func TimePeriod(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	hour := t.Hour()

	switch {
	case hour >= 0 && hour < 3:
		return "Late Night"
	case hour >= 3 && hour < 6:
		return "Dawn"
	case hour >= 6 && hour < 12:
		return "Morning"
	case hour == 12:
		return "Noon"
	case hour > 12 && hour < 18:
		return "Afternoon"
	case hour >= 18 && hour < 21:
		return "Evening"
	case hour >= 21 && hour < 24:
		return "Late Night"
	default:
		return "Unknown time period"
	}
}
