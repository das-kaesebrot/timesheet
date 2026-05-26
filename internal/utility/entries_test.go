package utility

import (
	"testing"
	"time"

	"github.com/das-kaesebrot/timesheet/internal/model"
	"github.com/google/uuid"
)

func TestSumEntryDurationsEmpty(t *testing.T) {
	result := SumEntryDurations(nil)
	if result != 0 {
		t.Errorf("expected 0, got %v", result)
	}
}

func TestSumEntryDurations(t *testing.T) {
	start := time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC)
	entries := []*model.TimesheetEntry{
		{ID: uuid.UUID{}, Start: start, End: start.Add(8 * time.Hour)},
		{ID: uuid.UUID{}, Start: start.Add(24 * time.Hour), End: start.Add(24*time.Hour + 4*time.Hour)},
	}

	result := SumEntryDurations(entries)
	expected := 12 * time.Hour
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSumEntryDurationsZeroDuration(t *testing.T) {
	now := time.Now()
	entries := []*model.TimesheetEntry{
		{ID: uuid.UUID{}, Start: now, End: now},
	}

	result := SumEntryDurations(entries)
	if result != 0 {
		t.Errorf("expected 0, got %v", result)
	}
}

func TestSumEntryDurationsNegativeDuration(t *testing.T) {
	now := time.Now()
	entries := []*model.TimesheetEntry{
		{ID: uuid.UUID{}, Start: now.Add(1 * time.Hour), End: now},
	}

	result := SumEntryDurations(entries)
	if result >= 0 {
		t.Errorf("expected negative duration, got %v", result)
	}
}

func TestSumEntryDurationsMultipleEntries(t *testing.T) {
	start := time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)
	entries := []*model.TimesheetEntry{
		{ID: uuid.UUID{}, Start: start, End: start.Add(1 * time.Hour)},
		{ID: uuid.UUID{}, Start: start.Add(2 * time.Hour), End: start.Add(3 * time.Hour)},
		{ID: uuid.UUID{}, Start: start.Add(4 * time.Hour), End: start.Add(10 * time.Hour)},
	}

	result := SumEntryDurations(entries)
	expected := 8 * time.Hour
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
