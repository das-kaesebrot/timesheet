package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestUserBeforeCreateSetsUUID(t *testing.T) {
	user := &User{}
	err := user.BeforeCreate(&gorm.DB{})
	if err != nil {
		t.Fatalf("BeforeCreate returned an error: %v", err)
	}
	if user.ID == uuid.Nil {
		t.Error("BeforeCreate did not set a UUID for User")
	}
}

func TestUserBeforeCreateSetsV7UUID(t *testing.T) {
	user := &User{}
	err := user.BeforeCreate(&gorm.DB{})
	if err != nil {
		t.Fatalf("BeforeCreate returned an error: %v", err)
	}
	if user.ID.Version() != 7 {
		t.Errorf("expected UUID version 7, got version %d", user.ID.Version())
	}
}

func TestTimesheetEntryBeforeCreateSetsUUID(t *testing.T) {
	entry := &TimesheetEntry{}
	err := entry.BeforeCreate(&gorm.DB{})
	if err != nil {
		t.Fatalf("BeforeCreate returned an error: %v", err)
	}
	if entry.ID == uuid.Nil {
		t.Error("BeforeCreate did not set a UUID for TimesheetEntry")
	}
}

func TestTimesheetEntryBeforeCreateSetsV7UUID(t *testing.T) {
	entry := &TimesheetEntry{}
	err := entry.BeforeCreate(&gorm.DB{})
	if err != nil {
		t.Fatalf("BeforeCreate returned an error: %v", err)
	}
	if entry.ID.Version() != 7 {
		t.Errorf("expected UUID version 7, got version %d", entry.ID.Version())
	}
}

func TestUserBeforeCreateUniqueIDs(t *testing.T) {
	user1 := &User{}
	user2 := &User{}
	user1.BeforeCreate(&gorm.DB{})
	user2.BeforeCreate(&gorm.DB{})
	if user1.ID == user2.ID {
		t.Error("BeforeCreate generated duplicate UUIDs for User")
	}
}

func TestTimesheetEntryBeforeCreateUniqueIDs(t *testing.T) {
	entry1 := &TimesheetEntry{}
	entry2 := &TimesheetEntry{}
	entry1.BeforeCreate(&gorm.DB{})
	entry2.BeforeCreate(&gorm.DB{})
	if entry1.ID == entry2.ID {
		t.Error("BeforeCreate generated duplicate UUIDs for TimesheetEntry")
	}
}

var timesheetOverlapsTests = []struct {
	inExisting *TimesheetEntry
	inNew      *TimesheetEntry
	out        bool
}{
	{
		&TimesheetEntry{Start: time.Date(2026, 05, 11, 10, 0, 0, 0, time.UTC), End: time.Date(2026, 05, 11, 10, 15, 0, 0, time.UTC)},
		&TimesheetEntry{Start: time.Date(2026, 05, 11, 11, 0, 0, 0, time.UTC), End: time.Date(2026, 05, 11, 11, 15, 0, 0, time.UTC)},
		false,
	},
	{
		&TimesheetEntry{Start: time.Date(2026, 05, 11, 10, 0, 0, 0, time.UTC), End: time.Date(2026, 05, 11, 10, 15, 0, 0, time.UTC)},
		&TimesheetEntry{Start: time.Date(2026, 05, 11, 10, 5, 0, 0, time.UTC), End: time.Date(2026, 05, 11, 10, 20, 0, 0, time.UTC)},
		true,
	},
	{
		&TimesheetEntry{Start: time.Date(2026, 05, 11, 10, 0, 0, 0, time.UTC), End: time.Date(2026, 05, 11, 10, 15, 0, 0, time.UTC)},
		&TimesheetEntry{Start: time.Date(2026, 05, 11, 10, 15, 0, 0, time.UTC), End: time.Date(2026, 05, 11, 10, 30, 0, 0, time.UTC)},
		false,
	},
	{
		&TimesheetEntry{Start: time.Date(2026, 05, 11, 10, 0, 0, 0, time.UTC), End: time.Date(2026, 05, 11, 10, 15, 0, 0, time.UTC)},
		&TimesheetEntry{Start: time.Date(2026, 05, 11, 9, 15, 0, 0, time.UTC), End: time.Date(2026, 05, 11, 10, 5, 0, 0, time.UTC)},
		true,
	},
}

func TestTimesheetOverlaps(t *testing.T) {
	for _, tt := range timesheetOverlapsTests {
		result := tt.inNew.Overlaps(tt.inExisting)
		if result != tt.out {
			t.Errorf("result is not expected value! result=%v, expected=%v", result, tt.out)
		}
	}
}
