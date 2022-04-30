package utils

import (
	"fmt"
	"time"
)

func GetFormatRequestTime(time time.Time) string {
	return fmt.Sprintf("%d.%d", time.Unix(), time.Nanosecond()/1e3)
}

func GetRequestCost(start, end time.Time) float64 {
	return float64(end.Sub(start).Nanoseconds()/1e4) / 100.0
}
