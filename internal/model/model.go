package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username         string
	PasswordHash     string
	Mail             *string // can be nullable
	Active           bool
	WeeklyWorkHours  uint8
	TimesheetEntries []TimesheetEntry
}

type TimesheetEntry struct {
	gorm.Model
	Start       time.Time
	End         time.Time
	UserID      uint
	Description *string
}
