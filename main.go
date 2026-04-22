package main

import (
	"context"
	"fmt"
	"time"

	"github.com/das-kaesebrot/timesheet/internal/model"
	"github.com/das-kaesebrot/timesheet/internal/password"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	ctx := context.Background()

	// Migrate the schema
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.TimesheetEntry{})

	// Create
	pwHash, err := password.HashPassword("admin")
	if err != nil {
		panic("Failed to hash password!")
	}

	description := "getting some work done"
	err = gorm.G[model.User](db).Create(ctx, &model.User{Username: "admin", PasswordHash: pwHash, Active: true, WeeklyWorkHours: 40})
	user, err := gorm.G[model.User](db).Where("username = ?", "admin").First(ctx)
	err = gorm.G[model.TimesheetEntry](db).Create(ctx, &model.TimesheetEntry{UserID: user.ID, Start: time.Now(), End: time.Now().Add(time.Hour), Description: &description})
	err = gorm.G[model.TimesheetEntry](db).Create(ctx, &model.TimesheetEntry{UserID: user.ID, Start: time.Now().Add(time.Hour), End: time.Now().Add(time.Hour), Description: &description})
	err = gorm.G[model.TimesheetEntry](db).Create(ctx, &model.TimesheetEntry{UserID: user.ID, Start: time.Now(), End: time.Now().Add(time.Hour), Description: &description})

	fmt.Printf("Found user: %v", user)
	fmt.Printf("Found time entries: %v", user.TimesheetEntries)
}
