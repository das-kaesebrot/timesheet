package model

import (
	"cmp"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                   uuid.UUID `gorm:"type:uuid"`
	Name                 string
	PasswordHash         string
	Active               bool
	WeeklyWorkTime       time.Duration
	StartOfWeek          time.Weekday
	TimesheetEntries     []TimesheetEntry
	TimesheetGranularity time.Duration
	DefaultTimezone      string
}

type UserUpdate struct {
	Name                 string
	Active               bool
	WeeklyWorkTime       *time.Duration
	StartOfWeek          *time.Weekday
	TimesheetGranularity *time.Duration
	DefaultTimezone      string
}

type TimesheetEntry struct {
	gorm.Model
	ID          uuid.UUID
	Start       time.Time
	End         time.Time
	UserID      uuid.UUID
	Description string
}

type TimesheetEntryUpdate struct {
	Name        string
	Start       time.Time
	End         time.Time
	Description string
}

// Note: Gorm will fail if the function signature
//	does not include `*gorm.DB` and `error`

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	user.ID = newUuid
	return
}

func (user *User) UpdateFromForm(form *UserUpdate) {
	if user.Name != form.Name {
		user.Name = form.Name
	}
	if user.Active != form.Active {
		user.Active = form.Active
	}
	if user.WeeklyWorkTime != *form.WeeklyWorkTime {
		user.WeeklyWorkTime = *form.WeeklyWorkTime
	}
	if user.StartOfWeek != *form.StartOfWeek {
		user.StartOfWeek = *form.StartOfWeek
	}
	if user.TimesheetGranularity != *form.TimesheetGranularity {
		user.TimesheetGranularity = *form.TimesheetGranularity
	}
	if user.DefaultTimezone != form.DefaultTimezone {
		user.DefaultTimezone = form.DefaultTimezone
	}
}

func (timesheetEntry *TimesheetEntry) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	timesheetEntry.ID = newUuid
	return
}

func (t *TimesheetEntry) String() string {
	return fmt.Sprintf("TimesheetEntry{ID=%s, Start=%s, End=%s, UserID=%s, Description=%s, CreatedAt=%s, UpdatedAt=%s, DeletedAt=%s}", t.ID.String(), t.Start.Format(time.RFC3339), t.End.Format(time.RFC3339), t.UserID.String(), t.Description, t.CreatedAt.Format(time.RFC3339), t.UpdatedAt.Format(time.RFC3339), t.DeletedAt.Time.Format(time.RFC3339))
}

func (timesheetEntry *TimesheetEntry) Overlaps(other *TimesheetEntry) bool {
	return timesheetEntry.Start.Before(other.End) && other.Start.Before(timesheetEntry.End)
}

func comparyByDate(a, b *TimesheetEntry) int {
	return cmp.Compare(a.Start.Nanosecond(), b.Start.Nanosecond())
	/*
		if a == b {
			return 0
		}
		if a.Start.Before(b.Start) {
			return -1
		}
		return 1
	*/
}
