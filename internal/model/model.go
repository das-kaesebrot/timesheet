package model

import (
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
