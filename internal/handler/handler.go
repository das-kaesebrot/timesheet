package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/client9/nowandlater"
	"github.com/das-kaesebrot/timesheet/internal/model"
	"github.com/das-kaesebrot/timesheet/internal/repository"
	"github.com/das-kaesebrot/timesheet/internal/template"
	"github.com/das-kaesebrot/timesheet/internal/utility"
	"github.com/google/uuid"
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
	h.renderer.Render(w, "users_list", map[string]interface{}{"Users": users})
}

func (h *Handler) ShowUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	summaries, err := h.GetWeeklySummariesForUser(user, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.renderer.Render(w, "users_show", map[string]interface{}{"User": user, "Summaries": summaries})
}

func (h *Handler) NewUser(w http.ResponseWriter, r *http.Request) {
	availableTimezones, err := utility.GetAllTimezones(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.renderer.Render(w, "users_new", map[string]interface{}{"Timezones": availableTimezones})
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

	parsedDuration, err := time.ParseDuration(r.PostForm.Get("timesheet_granularity"))
	if err != nil {
		http.Error(w, "Bad value for timesheet granularity, has to be a time string!", http.StatusBadRequest)
		return
	}
	parsedWeeklyWorkTime, err := time.ParseDuration(r.PostForm.Get("weekly_work_time"))
	if err != nil {
		http.Error(w, "Bad value for weekly work time, has to be a time string!", http.StatusBadRequest)
		return
	}

	availableTimezones, err := utility.GetAllTimezones(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	timezone := r.PostForm.Get("default_timezone")
	if !slices.Contains(availableTimezones, timezone) {
		http.Error(w, "Given timezone is not a valid timezone!", http.StatusBadRequest)
		return
	}

	user := &model.User{
		Username:             username,
		Description:          r.PostForm.Get("description"),
		Active:               r.PostForm.Get("active") == "on",
		WeeklyWorkTime:       parsedWeeklyWorkTime.Abs(),
		TimesheetGranularity: parsedDuration.Abs(),
		DefaultTimezone:      timezone,
	}

	if err := h.repo.CreateUser(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusFound)
}

func (h *Handler) EditUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	user, err := h.repo.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	timezones, err := utility.GetAllTimezones(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.renderer.Render(w, "users_edit", map[string]interface{}{"User": user, "Timezones": timezones})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), id)
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
	user.Description = r.PostForm.Get("description")

	parsedDuration, err := time.ParseDuration(r.PostForm.Get("timesheet_granularity"))
	if err != nil {
		http.Error(w, "Bad value for timesheet granularity, has to be a time string!", http.StatusBadRequest)
		return
	}
	user.TimesheetGranularity = parsedDuration.Abs()

	parsedWeeklyWorkTime, err := time.ParseDuration(r.PostForm.Get("weekly_work_time"))
	if err != nil {
		http.Error(w, "Bad value for weekly work time, has to be a time string!", http.StatusBadRequest)
		return
	}
	user.WeeklyWorkTime = parsedWeeklyWorkTime.Abs()

	availableTimezones, err := utility.GetAllTimezones(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	timezone := r.PostForm.Get("default_timezone")
	if !slices.Contains(availableTimezones, timezone) {
		http.Error(w, "Given timezone is not a valid timezone!", http.StatusBadRequest)
		return
	}
	user.DefaultTimezone = timezone

	if err := h.repo.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d", id), http.StatusFound)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusFound)
}

func (h *Handler) NewUserEntry(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	timezones, err := utility.GetAllTimezones(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.renderer.Render(w, "entries_new", map[string]interface{}{"User": user, "Timezones": timezones})
}

func (h *Handler) NewUserEntryQuick(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	timezones, err := utility.GetAllTimezones(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.renderer.Render(w, "entries_new_quick", map[string]interface{}{"User": user, "Timezones": timezones})
}

func (h *Handler) CreateUserEntryQuick(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loc, err := time.LoadLocation(user.DefaultTimezone)
	if err != nil {
		http.Error(w, "error parsing timezone: "+err.Error(), http.StatusBadRequest)
		return
	}

	p := nowandlater.Parser{Location: loc}
	start, err := p.Parse(r.PostForm.Get("natural_language_time_start"))
	if err != nil {
		http.Error(w, "Error parsing start: "+err.Error(), http.StatusBadRequest)
		return
	}
	end, err := p.Parse(r.PostForm.Get("natural_language_time_end"))
	if err != nil {
		http.Error(w, "Error parsing end: "+err.Error(), http.StatusBadRequest)
		return
	}

	if end.Before(start) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	existingEntries, _ := h.repo.GetTimesheetEntriesByUserID(r.Context(), userID)
	for _, e := range existingEntries {
		if !(end.Before(e.Start) || start.After(e.End)) {
			http.Error(w, "Time entry overlaps with existing entry", http.StatusBadRequest)
			return
		}
	}

	desc := r.PostForm.Get("description")
	entry := &model.TimesheetEntry{
		UserID:      userID,
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
func (h *Handler) CreateUserEntry(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loc, err := time.LoadLocation(r.PostForm.Get("timezone"))
	if err != nil {
		http.Error(w, "error parsing timezone: "+err.Error(), http.StatusBadRequest)
		return
	}

	start, err := time.ParseInLocation("2006-01-02T15:04", r.PostForm.Get("start"), loc)
	if err != nil {
		http.Error(w, "Invalid start time: "+err.Error(), http.StatusBadRequest)
		return
	}

	end, err := time.ParseInLocation("2006-01-02T15:04", r.PostForm.Get("end"), loc)
	if err != nil {
		http.Error(w, "Invalid end time: "+err.Error(), http.StatusBadRequest)
		return
	}

	if end.Before(start) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	duration := end.Sub(start)
	minutes := int(duration.Minutes())
	granularityMinutes := int(user.TimesheetGranularity.Minutes())
	if granularityMinutes > 0 && minutes%granularityMinutes != 0 {
		http.Error(w, fmt.Sprintf("Duration must be divisible by %v", user.TimesheetGranularity), http.StatusBadRequest)
		return
	}

	existingEntries, _ := h.repo.GetTimesheetEntriesByUserID(r.Context(), userID)
	for _, e := range existingEntries {
		if !(end.Before(e.Start) || start.After(e.End)) {
			http.Error(w, "Time entry overlaps with existing entry", http.StatusBadRequest)
			return
		}
	}

	desc := r.PostForm.Get("description")
	entry := &model.TimesheetEntry{
		UserID:      userID,
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
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	entry, err := h.repo.GetTimesheetEntryByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), entry.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.renderer.Render(w, "entries_edit", map[string]interface{}{"User": user, "Entry": entry})
}

func (h *Handler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	entry, err := h.repo.GetTimesheetEntryByID(r.Context(), id)
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

	loc, err := time.LoadLocation(r.PostForm.Get("timezone"))
	if err != nil {
		http.Error(w, "error parsing timezone: "+err.Error(), http.StatusBadRequest)
		return
	}

	start, err := time.ParseInLocation("2006-01-02T15:04", r.PostForm.Get("start"), loc)
	if err != nil {
		http.Error(w, "Invalid start time: "+err.Error(), http.StatusBadRequest)
		return
	}

	end, err := time.ParseInLocation("2006-01-02T15:04", r.PostForm.Get("end"), loc)
	if err != nil {
		http.Error(w, "Invalid end time: "+err.Error(), http.StatusBadRequest)
		return
	}

	if end.Before(start) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}

	duration := end.Sub(start)
	minutes := int(duration.Minutes())
	granularityMinutes := int(user.TimesheetGranularity.Minutes())
	if granularityMinutes > 0 && minutes%granularityMinutes != 0 {
		http.Error(w, fmt.Sprintf("Duration must be divisible by %v", user.TimesheetGranularity), http.StatusBadRequest)
		return
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
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid entry ID", http.StatusBadRequest)
		return
	}

	entry, err := h.repo.GetTimesheetEntryByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	userID := entry.UserID

	if err := h.repo.DeleteTimesheetEntry(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%d/entries", userID), http.StatusFound)
}

func (h *Handler) ExportUser(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), id)
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
		entries, err = h.repo.GetTimesheetEntriesByUserIDInRange(r.Context(), id, startTime, endTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		entries, err = h.repo.GetTimesheetEntriesByUserID(r.Context(), id)
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

type WeeklySummary struct {
	WeekNumberISO int
	StartOfWeek   time.Time
	EndOfWeek     time.Time
	TimeLogged    time.Duration
	Entries       []model.TimesheetEntry
}

func (h *Handler) GetWeeklySummariesForUser(u *model.User, r *http.Request) ([]WeeklySummary, error) {
	firstEntry, err := h.repo.GetEarliestTimesheetEntryByUserID(r.Context(), u.ID)
	if err != nil {
		return nil, err
	}

	rangeStart := utility.GetPreviousWeekStartDate(firstEntry.Start, u.StartOfWeek)
	rangeEnd := utility.GetNextWeekStartDate(time.Now(), u.StartOfWeek)

	weekStarts, err := utility.GetWeekRangeInWindow(rangeStart, rangeEnd, u.StartOfWeek)
	if err != nil {
		return nil, err
	}

	summaries := make([]WeeklySummary, len(weekStarts))

	for i, startDate := range weekStarts {
		entriesInWeek, err := h.repo.GetTimesheetEntriesByUserIDInRange(r.Context(), u.ID, startDate, startDate.AddDate(0, 0, 7))
		if err != nil {
			return nil, err
		}

		timeLogged := utility.SumEntryDurations(entriesInWeek)
		_, week := startDate.ISOWeek()
		summaries[i] = WeeklySummary{WeekNumberISO: week, StartOfWeek: startDate, EndOfWeek: startDate.AddDate(0, 0, 7).Add(-time.Second), TimeLogged: timeLogged, Entries: entriesInWeek}
	}

	return summaries, nil
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) JSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
