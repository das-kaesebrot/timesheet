package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
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
	Start       time.Time
	End         time.Time
	UserID      uint
	Description *string
}
