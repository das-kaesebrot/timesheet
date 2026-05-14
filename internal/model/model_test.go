package model

import (
	"testing"

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
