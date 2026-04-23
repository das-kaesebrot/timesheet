package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/das-kaesebrot/timesheet/internal/model"
	"github.com/das-kaesebrot/timesheet/internal/repository"
	"github.com/das-kaesebrot/timesheet/internal/template"
)

type Handler struct {
	repo     *repository.Repository
	renderer *template.Renderer
}

func New(repo *repository.Repository, renderer *template.Renderer) *Handler {
	return &Handler{repo: repo, renderer: renderer}
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/users", http.StatusFound)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.renderer.Render(w, "users/list", map[string]interface{}{"Users": users})
}

func (h *Handler) ShowUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	user, err := h.repo.GetUserByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.renderer.Render(w, "users/show", map[string]interface{}{"User": user})
}

func (h *Handler) NewUser(w http.ResponseWriter, r *http.Request) {
	h.renderer.Render(w, "users/new", nil)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.PostForm.Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	existing, _ := h.repo.GetUserByUsername(r.Context(), username)
	if existing != nil {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	weeklyWorkHours, _ := strconv.ParseUint(r.PostForm.Get("weekly_work_hours"), 10, 8)
	if weeklyWorkHours == 0 {
		weeklyWorkHours = 40
	}

	user := &model.User{
		Username:        username,
		Active:         r.PostForm.Get("active") == "on",
		WeeklyWorkHours: uint8(weeklyWorkHours),
	}

	if err := h.repo.CreateUser(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusFound)
}

func (h *Handler) EditUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	user, err := h.repo.GetUserByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.renderer.Render(w, "users/edit", map[string]interface{}{"User": user})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.PostForm.Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if username != user.Username {
		existing, _ := h.repo.GetUserByUsername(r.Context(), username)
		if existing != nil && existing.ID != user.ID {
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
		}
		user.Username = username
	}

	user.Active = r.PostForm.Get("active") == "on"
	if wh, err := strconv.ParseUint(r.PostForm.Get("weekly_work_hours"), 10, 8); err == nil && wh > 0 {
		user.WeeklyWorkHours = uint8(wh)
	}

	if err := h.repo.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d", id), http.StatusFound)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteUser(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusFound)
}

func (h *Handler) ListUserEntries(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	entries, err := h.repo.GetTimesheetEntriesByUserID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.renderer.Render(w, "entries/list", map[string]interface{}{"User": user, "Entries": entries})
}

func (h *Handler) NewUserEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.renderer.Render(w, "entries/new", map[string]interface{}{"User": user})
}

func (h *Handler) CreateUserEntry(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), uint(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	start, err := time.Parse("2006-01-02T15:04", r.PostForm.Get("start"))
	if err != nil {
		http.Error(w, "Invalid start time: "+err.Error(), http.StatusBadRequest)
		return
	}

	end, err := time.Parse("2006-01-02T15:04", r.PostForm.Get("end"))
	if err != nil {
		http.Error(w, "Invalid end time: "+err.Error(), http.StatusBadRequest)
		return
	}

	if end.Before(start) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	duration := end.Sub(start)
	if user.TimesheetGranularity != nil {
		minutes := int(duration.Minutes())
		granularityMinutes := int(user.TimesheetGranularity.Minutes())
		if granularityMinutes > 0 && minutes%granularityMinutes != 0 {
			http.Error(w, fmt.Sprintf("Duration must be divisible by %v", user.TimesheetGranularity), http.StatusBadRequest)
			return
		}
	}

	existingEntries, _ := h.repo.GetTimesheetEntriesByUserID(r.Context(), uint(userID))
	for _, e := range existingEntries {
		if !(end.Before(e.Start) || start.After(e.End)) {
			http.Error(w, "Time entry overlaps with existing entry", http.StatusBadRequest)
			return
		}
	}

	desc := r.PostForm.Get("description")
	entry := &model.TimesheetEntry{
		UserID:      uint(userID),
		Start:       start,
		End:         end,
		Description: &desc,
	}

	if err := h.repo.CreateTimesheetEntry(r.Context(), entry); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d/entries", userID), http.StatusFound)
}

func (h *Handler) EditEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	entry, err := h.repo.GetTimesheetEntryByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), entry.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.renderer.Render(w, "entries/edit", map[string]interface{}{"User": user, "Entry": entry})
}

func (h *Handler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	entry, err := h.repo.GetTimesheetEntryByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), entry.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	start, err := time.Parse("2006-01-02T15:04", r.PostForm.Get("start"))
	if err != nil {
		http.Error(w, "Invalid start time: "+err.Error(), http.StatusBadRequest)
		return
	}

	end, err := time.Parse("2006-01-02T15:04", r.PostForm.Get("end"))
	if err != nil {
		http.Error(w, "Invalid end time: "+err.Error(), http.StatusBadRequest)
		return
	}

	if end.Before(start) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	duration := end.Sub(start)
	if user.TimesheetGranularity != nil {
		minutes := int(duration.Minutes())
		granularityMinutes := int(user.TimesheetGranularity.Minutes())
		if granularityMinutes > 0 && minutes%granularityMinutes != 0 {
			http.Error(w, fmt.Sprintf("Duration must be divisible by %v", user.TimesheetGranularity), http.StatusBadRequest)
			return
		}
	}

	existingEntries, _ := h.repo.GetTimesheetEntriesByUserID(r.Context(), entry.UserID)
	for _, e := range existingEntries {
		if e.ID == entry.ID {
			continue
		}
		if !(end.Before(e.Start) || start.After(e.End)) {
			http.Error(w, "Time entry overlaps with existing entry", http.StatusBadRequest)
			return
		}
	}

	entry.Start = start
	entry.End = end
	desc := r.PostForm.Get("description")
	entry.Description = &desc

	if err := h.repo.UpdateTimesheetEntry(r.Context(), entry); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d/entries", entry.UserID), http.StatusFound)
}

func (h *Handler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	entry, err := h.repo.GetTimesheetEntryByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	userID := entry.UserID

	if err := h.repo.DeleteTimesheetEntry(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d/entries", userID), http.StatusFound)
}

type WeeklySummary struct {
	Week       time.Time
	StartDate  time.Time
	EndDate   time.Time
	TotalHours float64
	Delta     float64
}

func (h *Handler) OverviewUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	entries, err := h.repo.GetTimesheetEntriesByUserID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	weeklyMap := make(map[string]WeeklySummary)
	for _, e := range entries {
		year, week := e.Start.ISOWeek()
		key := fmt.Sprintf("%d-W%02d", year, week)
		if _, ok := weeklyMap[key]; !ok {
			monday := getMondayOfISOWeek(year, week)
			friday := monday.AddDate(0, 0, 4)
			weeklyMap[key] = WeeklySummary{
				Week:      time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
				StartDate: monday,
				EndDate:   friday,
			}
		}
		s := weeklyMap[key]
		s.TotalHours += e.End.Sub(e.Start).Hours()
		s.Delta = s.TotalHours - float64(user.WeeklyWorkHours)
		weeklyMap[key] = s
	}

	var summaries []WeeklySummary
	for _, s := range weeklyMap {
		summaries = append(summaries, s)
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].StartDate.Before(summaries[j].StartDate)
	})

	h.renderer.Render(w, "users/overview", map[string]interface{}{"User": user, "Summaries": summaries})
}

func getMondayOfISOWeek(year, week int) time.Time {
	jan1 := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	for jan1.Weekday() != time.Monday {
		jan1 = jan1.AddDate(0, 0, 1)
	}
	return jan1.AddDate(0, 0, (week-1)*7)
}

func (h *Handler) ExportUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var entries []model.TimesheetEntry
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	if start != "" && end != "" {
		startTime, err := time.Parse("2006-01-02", start)
		if err != nil {
			http.Error(w, "Invalid start date: "+err.Error(), http.StatusBadRequest)
			return
		}
		endTime, err := time.Parse("2006-01-02", end)
		if err != nil {
			http.Error(w, "Invalid end date: "+err.Error(), http.StatusBadRequest)
			return
		}
		entries, err = h.repo.GetTimesheetEntriesByUserIDInRange(r.Context(), uint(id), startTime, endTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		entries, err = h.repo.GetTimesheetEntriesByUserID(r.Context(), uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=timesheet_%d.csv", id))
	fmt.Fprintln(w, "user_id,username,start,end,description")

	for _, e := range entries {
		desc := ""
		if e.Description != nil {
			desc = *e.Description
		}
		fmt.Fprintf(w, "%d,%s,%s,%s,%s\n", user.ID, user.Username, e.Start.Format(time.RFC3339), e.End.Format(time.RFC3339), desc)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) JSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}