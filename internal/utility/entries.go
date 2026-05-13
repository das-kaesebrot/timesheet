package utility

import (
	"time"

	"github.com/das-kaesebrot/timesheet/internal/model"
)

func SumEntryDurations(entries []model.TimesheetEntry) (duration time.Duration) {
	for _, entry := range entries {
		duration += entry.End.Sub(entry.Start)
	}

	return duration
}
