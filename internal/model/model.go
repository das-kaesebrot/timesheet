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
	WeeklyWorkHours      uint8
	TimesheetEntries     []TimesheetEntry
	TimesheetGranularity *time.Duration
}

type TimesheetEntry struct {
	gorm.Model
	Start       time.Time
	End         time.Time
	UserID      uint
	Description *string
}
