package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/das-kaesebrot/timesheet/internal/model"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *Repository {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.TimesheetEntry{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return New(db)
}

func setupUser(t *testing.T, repo *Repository, username string) *model.User {
	t.Helper()
	user := &model.User{
		Name:                 username,
		Active:               true,
		WeeklyWorkTime:       40 * time.Hour,
		TimesheetGranularity: 15 * time.Minute,
		StartOfWeek:          time.Monday,
		DefaultTimezone:      "UTC",
	}
	if err := repo.CreateUser(context.Background(), user); err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return user
}

func TestCreateUser(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "testuser")

	if user.ID == uuid.Nil {
		t.Error("expected user ID to be set")
	}
}

func TestGetUserByID(t *testing.T) {
	repo := setupTestDB(t)
	created := setupUser(t, repo, "testuser")

	got, err := repo.GetUserByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}

	if got.Name != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", got.Name)
	}
	if got.ID != created.ID {
		t.Errorf("expected ID %v, got %v", created.ID, got.ID)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	repo := setupTestDB(t)

	_, err := repo.GetUserByID(context.Background(), uuid.Nil)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound, got %v", err)
	}
}

func TestListUsers(t *testing.T) {
	repo := setupTestDB(t)

	users, err := repo.ListUsers(context.Background())
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}

	setupUser(t, repo, "user1")
	setupUser(t, repo, "user2")
	setupUser(t, repo, "user3")

	users, err = repo.ListUsers(context.Background())
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

func TestUpdateUser(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "beforeupdate")

	user.Name = "afterupdate"
	if err := repo.UpdateUser(context.Background(), user); err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	got, err := repo.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GetUserByID after update failed: %v", err)
	}

	if got.Name != "afterupdate" {
		t.Errorf("expected username 'afterupdate', got '%s'", got.Name)
	}
}

func TestDeleteUser(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "todelete")

	if err := repo.DeleteUser(context.Background(), user.ID); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	_, err := repo.GetUserByID(context.Background(), user.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound after delete, got %v", err)
	}
}

func TestCreateTimesheetEntry(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "entryuser")

	desc := "test entry"
	now := time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC)
	entry := &model.TimesheetEntry{
		UserID:      user.ID,
		Start:       now,
		End:         now.Add(8 * time.Hour),
		Description: desc,
	}

	if err := repo.CreateTimesheetEntry(context.Background(), entry); err != nil {
		t.Fatalf("CreateTimesheetEntry failed: %v", err)
	}

	if entry.ID == uuid.Nil {
		t.Error("expected entry ID to be set")
	}
}

func TestGetTimesheetEntryByID(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "entryuser2")

	desc := "find me"
	now := time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC)
	entry := &model.TimesheetEntry{
		UserID:      user.ID,
		Start:       now,
		End:         now.Add(4 * time.Hour),
		Description: desc,
	}
	if err := repo.CreateTimesheetEntry(context.Background(), entry); err != nil {
		t.Fatalf("CreateTimesheetEntry failed: %v", err)
	}

	got, err := repo.GetTimesheetEntryByID(context.Background(), entry.ID)
	if err != nil {
		t.Fatalf("GetTimesheetEntryByID failed: %v", err)
	}

	if got.ID != entry.ID {
		t.Errorf("expected ID %v, got %v", entry.ID, got.ID)
	}
}

func TestGetTimesheetEntryByIDNotFound(t *testing.T) {
	repo := setupTestDB(t)

	_, err := repo.GetTimesheetEntryByID(context.Background(), uuid.Nil)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound, got %v", err)
	}
}

func TestGetTimesheetEntriesByUserID(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "listentries")

	now := time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC)
	for i := range 3 {
		desc := "entry"
		entry := &model.TimesheetEntry{
			UserID:      user.ID,
			Start:       now.Add(time.Duration(i) * 24 * time.Hour),
			End:         now.Add(time.Duration(i)*24*time.Hour + time.Hour),
			Description: desc,
		}
		if err := repo.CreateTimesheetEntry(context.Background(), entry); err != nil {
			t.Fatalf("CreateTimesheetEntry failed: %v", err)
		}
	}

	entries, err := repo.GetTimesheetEntriesByUserID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GetTimesheetEntriesByUserID failed: %v", err)
	}

	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}

	for i := 1; i < len(entries); i++ {
		if entries[i].Start.After(entries[i-1].Start) {
			t.Error("entries not ordered by start DESC")
		}
	}
}

func TestGetTimesheetEntriesByUserIDEmpty(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "noentries")

	entries, err := repo.GetTimesheetEntriesByUserID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GetTimesheetEntriesByUserID failed: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestGetEarliestTimesheetEntryByUserID(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "earlyentry")

	now := time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC)
	desc1 := "later"
	desc2 := "earlier"

	entry1 := &model.TimesheetEntry{UserID: user.ID, Start: now.Add(2 * time.Hour), End: now.Add(3 * time.Hour), Description: desc1}
	entry2 := &model.TimesheetEntry{UserID: user.ID, Start: now, End: now.Add(1 * time.Hour), Description: desc2}

	if err := repo.CreateTimesheetEntry(context.Background(), entry1); err != nil {
		t.Fatalf("CreateTimesheetEntry failed: %v", err)
	}
	if err := repo.CreateTimesheetEntry(context.Background(), entry2); err != nil {
		t.Fatalf("CreateTimesheetEntry failed: %v", err)
	}

	got, err := repo.GetEarliestTimesheetEntryByUserID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GetEarliestTimesheetEntryByUserID failed: %v", err)
	}

	if got.ID != entry2.ID {
		t.Errorf("expected earliest entry ID %v, got %v", entry2.ID, got.ID)
	}
}

func TestGetEarliestTimesheetEntryByUserIDNotFound(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "noearlyentry")

	_, err := repo.GetEarliestTimesheetEntryByUserID(context.Background(), user.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound, got %v", err)
	}
}

func TestGetTimesheetEntriesByUserIDInRange(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "rangeentry")

	now := time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)
	desc := "entry"

	entries := []*model.TimesheetEntry{
		{UserID: user.ID, Start: now, End: now.Add(1 * time.Hour), Description: desc},
		{UserID: user.ID, Start: now.Add(24 * time.Hour), End: now.Add(25 * time.Hour), Description: desc},
		{UserID: user.ID, Start: now.Add(48 * time.Hour), End: now.Add(49 * time.Hour), Description: desc},
		{UserID: user.ID, Start: now.Add(72 * time.Hour), End: now.Add(73 * time.Hour), Description: desc},
	}

	for _, e := range entries {
		if err := repo.CreateTimesheetEntry(context.Background(), e); err != nil {
			t.Fatalf("CreateTimesheetEntry failed: %v", err)
		}
	}

	got, err := repo.GetTimesheetEntriesByUserIDInRange(context.Background(), user.ID, now, now.Add(24*time.Hour))
	if err != nil {
		t.Fatalf("GetTimesheetEntriesByUserIDInRange failed: %v", err)
	}

	if len(got) != 2 {
		t.Errorf("expected 2 entries in range, got %d", len(got))
	}
}

func TestGetTimesheetEntriesByUserIDInRangeEmpty(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "emptyrange")

	now := time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)
	got, err := repo.GetTimesheetEntriesByUserIDInRange(context.Background(), user.ID, now, now.Add(24*time.Hour))
	if err != nil {
		t.Fatalf("GetTimesheetEntriesByUserIDInRange failed: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("expected 0 entries, got %d", len(got))
	}
}

func TestUpdateTimesheetEntry(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "updateentry")

	desc := "original"
	now := time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC)
	entry := &model.TimesheetEntry{
		UserID:      user.ID,
		Start:       now,
		End:         now.Add(8 * time.Hour),
		Description: desc,
	}
	if err := repo.CreateTimesheetEntry(context.Background(), entry); err != nil {
		t.Fatalf("CreateTimesheetEntry failed: %v", err)
	}

	newDesc := "updated"
	entry.Description = newDesc
	entry.End = now.Add(6 * time.Hour)
	if err := repo.UpdateTimesheetEntry(context.Background(), entry); err != nil {
		t.Fatalf("UpdateTimesheetEntry failed: %v", err)
	}

	got, err := repo.GetTimesheetEntryByID(context.Background(), entry.ID)
	if err != nil {
		t.Fatalf("GetTimesheetEntryByID after update failed: %v", err)
	}

	if got.Description != "updated" {
		t.Errorf("expected description 'updated', got '%s'", got.Description)
	}
	if got.End != now.Add(6*time.Hour) {
		t.Errorf("expected end %v, got %v", now.Add(6*time.Hour), got.End)
	}
}

func TestDeleteTimesheetEntry(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "deleteentry")

	desc := "delete me"
	now := time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC)
	entry := &model.TimesheetEntry{
		UserID:      user.ID,
		Start:       now,
		End:         now.Add(1 * time.Hour),
		Description: desc,
	}
	if err := repo.CreateTimesheetEntry(context.Background(), entry); err != nil {
		t.Fatalf("CreateTimesheetEntry failed: %v", err)
	}

	if err := repo.DeleteTimesheetEntries(context.Background(), []uuid.UUID{entry.ID}); err != nil {
		t.Fatalf("DeleteTimesheetEntries failed: %v", err)
	}

	_, err := repo.GetTimesheetEntryByID(context.Background(), entry.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected ErrRecordNotFound after delete, got %v", err)
	}
}

func TestPreloadTimesheetEntries(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "preloaduser")

	now := time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC)
	desc := "preloaded"
	for range 2 {
		entry := &model.TimesheetEntry{
			UserID:      user.ID,
			Start:       now,
			End:         now.Add(1 * time.Hour),
			Description: desc,
		}
		if err := repo.CreateTimesheetEntry(context.Background(), entry); err != nil {
			t.Fatalf("CreateTimesheetEntry failed: %v", err)
		}
	}

	got, err := repo.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}

	if len(got.TimesheetEntries) != 2 {
		t.Errorf("expected 2 preloaded entries, got %d", len(got.TimesheetEntries))
	}
}

func TestDeleteUserDoesNotDeleteEntries(t *testing.T) {
	repo := setupTestDB(t)
	user := setupUser(t, repo, "orphanentries")

	now := time.Date(2026, 5, 11, 9, 0, 0, 0, time.UTC)
	desc := "orphaned entry"
	entry := &model.TimesheetEntry{
		UserID:      user.ID,
		Start:       now,
		End:         now.Add(1 * time.Hour),
		Description: desc,
	}
	if err := repo.CreateTimesheetEntry(context.Background(), entry); err != nil {
		t.Fatalf("CreateTimesheetEntry failed: %v", err)
	}

	entriesBefore, _ := repo.GetTimesheetEntriesByUserID(context.Background(), user.ID)
	if len(entriesBefore) != 1 {
		t.Fatalf("expected 1 entry before delete, got %d", len(entriesBefore))
	}

	if err := repo.DeleteUser(context.Background(), user.ID); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	entriesAfter, _ := repo.GetTimesheetEntriesByUserID(context.Background(), user.ID)
	if len(entriesAfter) != 1 {
		t.Errorf("expected entries to remain after user delete (no cascade), got %d", len(entriesAfter))
	}
}
