# Timesheet Application Project Plan

## Overview
A lightweight Go-based web application for tracking timesheet entries using existing models, SQLite, Bootstrap, HTMX, and Go's standard templating.

## Technology Stack
- **Backend**: Go with GORM (existing)
- **Database**: SQLite (existing)
- **Frontend**: Go templates, Bootstrap 5 (CDN), HTMX (CDN)
- **Architecture**: Standard Go project layout (`/internal/`)

## Project Structure
```
timesheet/
├── cmd/
│   └── server/
│       └── main.go          # Web server entry point
├── internal/
│   ├── handler/             # HTTP handlers (net/http)
│   ├── repository/          # GORM database operations
│   ├── template/            # Go HTML templates
│   └── model/               # Existing models (do not modify)
├── templates/               # HTML templates
│   ├── layout.html          # Base layout with Bootstrap + HTMX
│   ├── users/
│   │   ├── list.html
│   │   ├── form.html
│   │   └── ...
│   └── entries/
│       └── ...
├── go.mod
└── main.go                  # Keep existing for GORM setup
```

## Feature Implementation Plan

### Phase 1: Foundation
- [ ] Create `cmd/server/main.go` with `http.ServeMux`
- [ ] Reuse existing model and password packages
- [ ] Create the database using existing GORM setup
- [ ] Create repository layer (User CRUD, TimesheetEntry CRUD)
- [ ] Create base template with Bootstrap + HTMX
- [ ] Write HTML templates for all pages

### Phase 2: User Management Pages
- [ ] List users page (`GET /users`)
- [ ] Create user form (`GET /users/new`, `POST /users`)
- [ ] Edit user form (`GET /users/:id/edit`)
- [ ] Update user (`PATCH /users/:id`)
- [ ] Delete user (`DELETE /users/:id`)

### Phase 3: Timesheet Entry Submission
- [ ] Submission form (`GET /users/:id/entries/new`)
- [ ] Create entry (`POST /users/:id/entries`)
- [ ] Client-side validation for TimesheetGranularity
- [ ] Server-side validation for TimesheetGranularity
- [ ] Handle nil TimesheetGranularity (allow any value)

### Phase 4: Timesheet Entry Management
- [ ] List entries page (`GET /users/:id/entries`)
- [ ] Edit entry form (`GET /entries/:id/edit`)
- [ ] Update entry (`PATCH /entries/:id`)
- [ ] Delete entry (`DELETE /entries/:id`)

### Phase 5: Timesheet Overview
- [ ] Overview page with weekly summary (`GET /users/:id/overview`)
- [ ] Calculate weekly sums (Monday-Friday grouped by ISO week)
- [ ] Calculate delta: `sum hours - user.WeeklyWorkHours`

### Phase 6: Export Functionality
- [ ] CSV export endpoint (`GET /users/:id/export`)
- [ ] Optional date range filtering (`?start=YYYY-MM-DD&end=YYYY-MM-DD`)
- [ ] Handle nil granularity (allow any value) for exported entries

### Phase 7: Error Handling
- [ ] Validate form input server-side
- [ ] Return meaningful error messages
- [ ] Handle database errors gracefully

## API Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Home/Dashboard |
| GET | `/users` | List all users |
| GET | `/users/new` | Create user form |
| POST | `/users` | Create user action |
| GET | `/users/:id` | User details |
| GET | `/users/:id/edit` | Edit user form |
| PATCH | `/users/:id` | Update user action |
| DELETE | `/users/:id` | Delete user action |
| GET | `/users/:id/entries/new` | New entry form |
| POST | `/users/:id/entries` | Create entry action |
| GET | `/users/:id/entries` | List entries for user |
| GET | `/users/:id/overview` | Weekly overview page |
| GET | `/entries/:id/edit` | Edit entry form |
| GET | `/entries/:id/edit` | Edit entry form |
| PATCH | `/entries/:id` | Update entry action |
| DELETE | `/entries/:id` | Delete entry action |
| GET | `/users/:id/export` | Export timesheet to CSV |
| GET | `/users/:id/export?start=&end=` | Export with time range filter |

## Key Implementation Details

### Timesheet Granularity Validation
- If `User.TimesheetGranularity` is nil, allow any duration
- Otherwise, entry duration must be divisible by granularity
- Duration calculation: `entry.End - entry.Start`
- Validate on both client (JS) and server (Go)

### Weekly Overview Calculation
- Group entries by ISO week (Monday-Friday)
- Sum duration per week
- Calculate delta: `sum - user.WeeklyWorkHours`

### Export Format (CSV)
```csv
user_id,username,start,end,description
1,admin,2024-01-01T09:00:00Z,2024-01-01T17:00:00Z,Working on feature X
```

## Technology Stack
- **Backend**: Go with standard `net/http` router
- **Database**: SQLite with GORM (existing)
- **Frontend**: Go templates, Bootstrap 5 (CDN), HTMX (CDN)
- **Architecture**: Standard Go project layout (`/internal/`)

## Notes
- No JavaScript compilation steps
- Server-side HTML rendering with Go templates
- Bootstrap + HTMX via CDN
- No external router - use `http.ServeMux` or manual routing