package main

import (
	"context"

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
	err = gorm.G[model.User](db).Create(ctx, &model.User{Username: "admin", PasswordHash: pwHash, Active: false, WeeklyWorkHours: 40})
}
