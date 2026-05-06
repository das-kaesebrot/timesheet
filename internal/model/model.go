package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                   uuid.UUID
	Username             string
	Description          string
	PasswordHash         string
	Active               bool
	WeeklyWorkTime       time.Duration
	TimesheetEntries     []TimesheetEntry
	TimesheetGranularity time.Duration
	DefaultTimezone      string
}

type TimesheetEntry struct {
	gorm.Model
	ID          uuid.UUID
	Start       time.Time
	End         time.Time
	UserID      uuid.UUID
	Description *string
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

func (timesheetEntry *TimesheetEntry) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	timesheetEntry.ID = newUuid
	return
}
