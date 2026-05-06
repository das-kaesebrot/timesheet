package repository

import (
	"context"
	"time"

	"github.com/das-kaesebrot/timesheet/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("TimesheetEntries").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) ListUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *Repository) UpdateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *Repository) CreateTimesheetEntry(ctx context.Context, entry *model.TimesheetEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *Repository) GetTimesheetEntryByID(ctx context.Context, id uuid.UUID) (*model.TimesheetEntry, error) {
	var entry model.TimesheetEntry
	err := r.db.WithContext(ctx).First(&entry, id).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *Repository) GetTimesheetEntriesByUserID(ctx context.Context, userID uuid.UUID) ([]model.TimesheetEntry, error) {
	var entries []model.TimesheetEntry
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("start DESC").Find(&entries).Error
	return entries, err
}

func (r *Repository) GetTimesheetEntriesByUserIDInRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]model.TimesheetEntry, error) {
	var entries []model.TimesheetEntry
	end = end.AddDate(0, 0, 1)
	err := r.db.WithContext(ctx).Where("user_id = ? AND start >= ? AND end < ?", userID, start, end).Order("start DESC").Find(&entries).Error
	return entries, err
}

func (r *Repository) UpdateTimesheetEntry(ctx context.Context, entry *model.TimesheetEntry) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

func (r *Repository) DeleteTimesheetEntry(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.TimesheetEntry{}, id).Error
}
