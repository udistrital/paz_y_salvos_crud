package utils

import (
	"time"

	"github.com/udistrital/utils_oas/v2/time_bogota"
)

func HoraBogota() time.Time {
	actual := timebogota.TiempoBogota()
	return time.Date(actual.Year(), actual.Month(), actual.Day(), actual.Hour(), actual.Minute(), actual.Second(), actual.Nanosecond(), time.UTC)
}
