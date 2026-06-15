package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/das-kaesebrot/timesheet/internal/httperror"
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

type SortOrder int

const (
	SortOrderDescending SortOrder = iota
	SortOrderAscending
)

func New(repo *repository.Repository, renderer *template.Renderer) *Handler {
	return &Handler{repo: repo, renderer: renderer}
}

// catchall route
func (h *Handler) Root(w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, "/users", http.StatusFound)
	return nil
}

func (h *Handler) GetFavicon(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "image/svg+xml")
	fmt.Fprintf(w, "<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 100 100\">\n<text y=\".9em\" font-size=\"90\">%s</text>\n</svg>", "⌚")
	return nil
}

func (h *Handler) GetUsersList(w http.ResponseWriter, r *http.Request) error {
	users, err := h.repo.ListUsers(r.Context())
	if err != nil {
		return httperror.InternalServerError(err)
	}
	h.renderer.Render(w, "users_list", map[string]interface{}{"Users": users})
	return nil
}

func (h *Handler) GetUserOverview(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return httperror.BadRequest("Invalid user ID")
	}

	user, err := h.repo.GetUserByID(r.Context(), id)
	if err != nil {
		return httperror.New(http.StatusNotFound, "User not found", err)
	}

	summaries, err := h.GetWeeklySummariesForUser(user, r, SortOrderDescending)
	if err != nil {
		return httperror.InternalServerError(err)
	}

	availablePageSizes := []int{1, 5, 10}
	page := 1
	perPage := 5
	if queryPage := r.URL.Query().Get("page"); queryPage != "" {
		if n, err := strconv.Atoi(queryPage); err == nil && n > 0 {
			page = n
		}
	}
	if queryPerPage := r.URL.Query().Get("per_page"); queryPerPage != "" {
		if n, err := strconv.Atoi(queryPerPage); err == nil {
			for item := range availablePageSizes {
				if item == n {
					perPage = n
				}
			}
		}
	}

	totalSummaries := len(summaries)
	totalPages := (totalSummaries + perPage - 1) / perPage
	if totalPages == 0 {
		totalPages = 1
	}

	// overflow protection
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * perPage
	if start > totalSummaries {
		start = totalSummaries
	}
	end := start + perPage
	if end > totalSummaries {
		end = totalSummaries
	}

	pageSummaries := summaries[start:end]

	h.renderer.Render(w, "users_show", map[string]interface{}{
		"User":               user,
		"Summaries":          pageSummaries,
		"Page":               page,
		"PerPage":            perPage,
		"AvailablePageSizes": availablePageSizes,
		"TotalPages":         totalPages,
		"TotalSummaries":     totalSummaries,
	})
	return nil
}

func (h *Handler) GetUserNew(w http.ResponseWriter, r *http.Request) error {
	h.renderer.Render(w, "users_new", map[string]interface{}{})
	return nil
}

func (h *Handler) PostUserNew(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid form data", err)
	}
	userUpdate, err := parseUserForm(r.PostForm)
	if err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid user data", err)
	}

	user := &model.User{}
	user.UpdateFromForm(userUpdate)

	if err := h.repo.CreateUser(r.Context(), user); err != nil {
		return httperror.InternalServerError(err)
	}

	http.Redirect(w, r, "/users", http.StatusFound)
	return nil
}

func (h *Handler) GetUserEdit(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return httperror.BadRequest("Invalid user ID")
	}
	user, err := h.repo.GetUserByID(r.Context(), id)
	if err != nil {
		return httperror.New(http.StatusNotFound, "User not found", err)
	}

	h.renderer.Render(w, "users_edit", map[string]interface{}{"User": user})
	return nil
}

func (h *Handler) PostUserUpdate(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return httperror.BadRequest("Invalid user ID")
	}
	user, err := h.repo.GetUserByID(r.Context(), id)
	if err != nil {
		return httperror.New(http.StatusNotFound, "User not found", err)
	}
	if err := r.ParseForm(); err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid form data", err)
	}
	userUpdate, err := parseUserForm(r.PostForm)
	if err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid user data", err)
	}
	user.UpdateFromForm(userUpdate)
	if err := h.repo.UpdateUser(r.Context(), user); err != nil {
		return httperror.InternalServerError(err)
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%s", id.String()), http.StatusFound)
	return nil
}

func (h *Handler) PostUserDelete(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return httperror.BadRequest("Invalid user ID")
	}

	if err := h.repo.DeleteUser(r.Context(), id); err != nil {
		return httperror.InternalServerError(err)
	}

	http.Redirect(w, r, "/users", http.StatusFound)
	return nil
}

func (h *Handler) GetEntryNew(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return httperror.BadRequest("Invalid user ID")
	}

	user, err := h.repo.GetUserByID(r.Context(), id)
	if err != nil {
		return httperror.New(http.StatusNotFound, "User not found", err)
	}

	h.renderer.Render(w, "entries_new", map[string]interface{}{"User": user})
	return nil
}

func (h *Handler) GetEntryNewQuick(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return httperror.BadRequest("Invalid user ID")
	}

	user, err := h.repo.GetUserByID(r.Context(), id)
	if err != nil {
		return httperror.New(http.StatusNotFound, "User not found", err)
	}

	h.renderer.Render(w, "entries_new_quick", map[string]interface{}{"User": user})
	return nil
}

func (h *Handler) PostEntryNew(w http.ResponseWriter, r *http.Request) error {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return httperror.BadRequest("Invalid user ID")
	}

	user, err := h.repo.GetUserByID(r.Context(), userID)
	if err != nil {
		return httperror.New(http.StatusNotFound, "User not found", err)
	}

	if err := r.ParseForm(); err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid form data", err)
	}

	loc, err := time.LoadLocation(r.PostForm.Get("timezone"))
	if err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid timezone", err)
	}

	newEntryStart, err := time.ParseInLocation("2006-01-02T15:04", r.PostForm.Get("start"), loc)
	if err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid start time", err)
	}

	newEntryEnd, err := time.ParseInLocation("2006-01-02T15:04", r.PostForm.Get("end"), loc)
	if err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid end time", err)
	}

	if newEntryEnd.Before(newEntryStart) {
		return httperror.BadRequest("End time must be after start time")
	}

	duration := newEntryEnd.Sub(newEntryStart)

	if duration == 0 {
		return httperror.BadRequest("Duration must be longer than 0")
	}

	minutes := int(duration.Minutes())
	granularityMinutes := int(user.TimesheetGranularity.Minutes())
	if granularityMinutes > 0 && minutes%granularityMinutes != 0 {
		return httperror.New(http.StatusBadRequest, fmt.Sprintf("Duration must be divisible by %v", user.TimesheetGranularity), nil)
	}

	desc := r.PostForm.Get("description")
	newEntry := &model.TimesheetEntry{
		UserID:      userID,
		Start:       newEntryStart,
		End:         newEntryEnd,
		Description: desc,
	}

	existingEntries, _ := h.repo.GetTimesheetEntriesByUserID(r.Context(), userID)
	for _, existingEntry := range existingEntries {
		if newEntry.Overlaps(existingEntry) {
			return httperror.BadRequest("Time entry overlaps with existing entry")
		}
	}

	if err := h.repo.CreateTimesheetEntry(r.Context(), newEntry); err != nil {
		return httperror.InternalServerError(err)
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%s", userID.String()), http.StatusFound)
	return nil
}

func (h *Handler) GetEntryEdit(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return httperror.BadRequest("Invalid entry ID")
	}

	entry, err := h.repo.GetTimesheetEntryByID(r.Context(), id)
	if err != nil {
		return httperror.New(http.StatusNotFound, "Entry not found", err)
	}

	user, err := h.repo.GetUserByID(r.Context(), entry.UserID)
	if err != nil {
		return httperror.New(http.StatusNotFound, "User not found", err)
	}

	h.renderer.Render(w, "entries_edit", map[string]interface{}{"User": user, "Entry": entry})
	return nil
}

func (h *Handler) PostEntryUpdate(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return httperror.BadRequest("Invalid entry ID")
	}

	entry, err := h.repo.GetTimesheetEntryByID(r.Context(), id)
	if err != nil {
		return httperror.New(http.StatusNotFound, "Entry not found", err)
	}

	user, err := h.repo.GetUserByID(r.Context(), entry.UserID)
	if err != nil {
		return httperror.New(http.StatusNotFound, "User not found", err)
	}

	if err := r.ParseForm(); err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid form data", err)
	}

	loc, err := time.LoadLocation(r.PostForm.Get("timezone"))
	if err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid timezone", err)
	}

	start, err := time.ParseInLocation("2006-01-02T15:04", r.PostForm.Get("start"), loc)
	if err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid start time", err)
	}

	end, err := time.ParseInLocation("2006-01-02T15:04", r.PostForm.Get("end"), loc)
	if err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid end time", err)
	}

	if end.Before(start) {
		return httperror.BadRequest("End time must be after start time")
	}

	duration := end.Sub(start)
	minutes := int(duration.Minutes())
	granularityMinutes := int(user.TimesheetGranularity.Minutes())
	if granularityMinutes > 0 && minutes%granularityMinutes != 0 {
		return httperror.New(http.StatusBadRequest, fmt.Sprintf("Duration must be divisible by %v", user.TimesheetGranularity), nil)
	}

	existingEntries, _ := h.repo.GetTimesheetEntriesByUserID(r.Context(), entry.UserID)
	for _, e := range existingEntries {
		if e.ID == entry.ID {
			continue
		}
		if !(end.Before(e.Start) || start.After(e.End)) {
			return httperror.BadRequest("Time entry overlaps with existing entry")
		}
	}

	entry.Start = start
	entry.End = end
	desc := r.PostForm.Get("description")
	entry.Description = desc

	if err := h.repo.UpdateTimesheetEntry(r.Context(), entry); err != nil {
		return httperror.InternalServerError(err)
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%s", entry.UserID.String()), http.StatusFound)
	return nil
}

func (h *Handler) PostEntryDelete(w http.ResponseWriter, r *http.Request) error {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return httperror.BadRequest("Invalid user ID")
	}

	if err := r.ParseForm(); err != nil {
		return httperror.New(http.StatusBadRequest, "Invalid form data", err)
	}

	entryIDs := r.PostForm["entry_ids"]
	if len(entryIDs) == 0 {
		return httperror.BadRequest("No entries selected")
	}

	ids := make([]uuid.UUID, 0, len(entryIDs))
	for _, idStr := range entryIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return httperror.BadRequest("Invalid entry ID")
		}
		ids = append(ids, id)
	}

	if err := h.repo.DeleteTimesheetEntries(r.Context(), ids); err != nil {
		return httperror.InternalServerError(err)
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%s", userID.String()), http.StatusFound)
	return nil
}

func (h *Handler) ExportUser(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return httperror.BadRequest("Invalid user ID")
	}

	user, err := h.repo.GetUserByID(r.Context(), id)
	if err != nil {
		return httperror.New(http.StatusNotFound, "User not found", err)
	}

	var entries []*model.TimesheetEntry
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	if start != "" && end != "" {
		startTime, err := time.Parse("2006-01-02", start)
		if err != nil {
			return httperror.New(http.StatusBadRequest, "Invalid start date", err)
		}
		endTime, err := time.Parse("2006-01-02", end)
		if err != nil {
			return httperror.New(http.StatusBadRequest, "Invalid end date", err)
		}
		entries, err = h.repo.GetTimesheetEntriesByUserIDInRange(r.Context(), id, startTime, endTime)
		if err != nil {
			return httperror.InternalServerError(err)
		}
	} else {
		entries, err = h.repo.GetTimesheetEntriesByUserID(r.Context(), id)
		if err != nil {
			return httperror.InternalServerError(err)
		}
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=timesheet_%d.csv", id))
	fmt.Fprintln(w, "user_id,username,start,end,description")

	for _, e := range entries {
		desc := e.Description
		fmt.Fprintf(w, "%d,%s,%s,%s,%s\n", user.ID, user.Name, e.Start.Format(time.RFC3339), e.End.Format(time.RFC3339), desc)
	}
	return nil
}

type WeeklySummary struct {
	WeekNumberISO int
	StartOfWeek   time.Time
	EndOfWeek     time.Time
	TimeLogged    time.Duration
	WeeklyDiff    time.Duration
	Entries       []*model.TimesheetEntry
}

func (h *Handler) GetWeeklySummariesForUser(u *model.User, r *http.Request, order SortOrder) ([]WeeklySummary, error) {
	entries, err := h.repo.CountTimesheetEntriesByUserID(r.Context(), u.ID)
	if err != nil {
		return nil, err
	}
	if entries <= 0 {
		return nil, nil
	}

	firstEntry, err := h.repo.GetEarliestTimesheetEntryByUserID(r.Context(), u.ID)
	if err != nil {
		return nil, err
	}

	rangeStart := utility.GetPreviousWeekStartDate(firstEntry.Start, u.StartOfWeek)
	rangeEnd := utility.GetNextWeekStartDate(time.Now(), u.StartOfWeek)

	if rangeEnd.Equal(rangeStart) {
		rangeEnd = rangeEnd.AddDate(0, 0, 1)
	}

	weekStarts, err := utility.GetWeekRangeInWindow(rangeStart, rangeEnd, u.StartOfWeek)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(weekStarts, utility.CompareTimes) // ascending order

	var orderDescending = true
	if order == SortOrderAscending {
		orderDescending = false
	} else {
		slices.Reverse(weekStarts)
	}

	summaries := make([]WeeklySummary, len(weekStarts))

	for i, startDate := range weekStarts {
		entriesInWeek, err := h.repo.GetTimesheetEntriesByUserIDInRangeInOrder(r.Context(), u.ID, startDate, startDate.AddDate(0, 0, 7), orderDescending)
		if err != nil {
			return nil, err
		}

		timeLogged := utility.SumEntryDurations(entriesInWeek)
		_, week := startDate.ISOWeek()
		summaries[i] = WeeklySummary{WeekNumberISO: week, StartOfWeek: startDate, EndOfWeek: startDate.AddDate(0, 0, 7).Add(-time.Second), TimeLogged: timeLogged, Entries: entriesInWeek, WeeklyDiff: (timeLogged - u.WeeklyWorkTime)}
	}

	return summaries, nil
}

func parseUserForm(form url.Values) (*model.UserUpdate, error) {
	var userUpdate = new(model.UserUpdate)

	userUpdate.Name = form.Get("name")

	parsedDuration, err := time.ParseDuration(form.Get("timesheet_granularity"))
	if err != nil {
		return nil, fmt.Errorf("Bad value for timesheet granularity, has to be a time string! %w", err)
	}
	parsedDuration = parsedDuration.Abs()
	userUpdate.TimesheetGranularity = &parsedDuration

	parsedWeeklyWorkTime, err := time.ParseDuration(form.Get("weekly_work_time"))
	if err != nil {
		return nil, fmt.Errorf("Bad value for weekly work time, has to be a time string! %w", err)
	}
	parsedWeeklyWorkTime = parsedWeeklyWorkTime.Abs()
	userUpdate.WeeklyWorkTime = &parsedWeeklyWorkTime

	n, err := strconv.Atoi(form.Get("week_start_day"))
	if err != nil {
		return nil, err
	}
	if n < 0 || n > int(time.Saturday) {
		return nil, fmt.Errorf("Week day has to be between 0 and 6! Given value: %d", n)
	}
	parsedWeekStartDay := time.Weekday(n)
	userUpdate.StartOfWeek = &parsedWeekStartDay

	availableTimezones, err := utility.GetAllTimezones(true)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	timezone := form.Get("default_timezone")
	if !slices.Contains(availableTimezones, timezone) {
		return nil, fmt.Errorf("Given timezone is not a valid timezone! %w", err)
	}
	userUpdate.DefaultTimezone = timezone

	userUpdate.Active = form.Get("active") == "on"

	return userUpdate, nil
}
