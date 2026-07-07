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
│       └── main.go              # Web server entry point
├── internal/
│   ├── handler/                 # HTTP handlers (net/http)
│   ├── httperror/               # Structured HTTP error type
│   ├── middleware/              # Logger + error recovery middleware
│   ├── model/                   # GORM models (do not modify)
│   ├── password/                # bcrypt password hashing (unused)
│   ├── repository/              # GORM database operations
│   ├── template/                # Go HTML template rendering engine
│   └── utility/                 # Timezone, week, entry helpers
├── web/
│   ├── static/
│   │   ├── js/
│   │   │   └── script.js        # Dark mode toggle + batch-delete
│   │   └── libs/
│   │       └── bootstrap@5.3.8/ # Vendored Bootstrap 5.3.8
│   └── template/
│       ├── layouts/
│       │   └── base.html        # Base layout with Bootstrap + HTMX
│       ├── partials/            # Reusable template fragments
│       └── *.html               # Page templates
├── Dockerfile
├── go.mod
└── go.sum
```

## Feature Implementation Plan

### Phase 1: Foundation
- [x] Create `cmd/server/main.go` with `http.ServeMux`
- [x] Reuse existing model and password packages
- [x] Create the database using existing GORM setup
- [x] Create repository layer (User CRUD, TimesheetEntry CRUD)
- [x] Create base template with Bootstrap + HTMX
- [x] Write HTML templates for all pages

### Phase 2: User Management Pages
- [x] List users page (`GET /users`)
- [x] Create user form (`GET /users/new`, `POST /users`)
- [x] Edit user form (`GET /users/:id/edit`)
- [x] Update user (`PATCH /users/:id`)
- [x] Delete user (`DELETE /users/:id`)

### Phase 3: Timesheet Entry Submission
- [x] Submission form (`GET /users/:id/entries/new`)
- [x] Create entry (`POST /users/:id/entries`)
- [x] Client-side validation for TimesheetGranularity
- [x] Server-side validation for TimesheetGranularity

### Phase 4: Timesheet Entry Management
- [x] List entries page (`GET /users/:id/entries`)
- [x] Edit entry form (`GET /entries/:id/edit`)
- [x] Update entry (`POST /entries/:id/edit`)
- [x] Delete entry (`POST /entries/:id/delete`)
- [x] Batch delete entries with per-week checkboxes and select-all (`POST /users/:id/entries/batch-delete`)

### Phase 5: Timesheet Overview
- [x] Overview page with weekly summary (`GET /users/:id/overview`)
- [x] Calculate weekly sums (Monday-Friday grouped by ISO week)
- [x] Calculate delta: `sum hours - user.WeeklyWorkHours`
- [x] Pagination with configurable weeks per page (1, 5, 10) via `?page=&per_page=`

### Phase 6: Export Functionality
- [x] CSV export endpoint (`GET /users/:id/entries/export`)
- [x] Optional date range filtering (`?start=YYYY-MM-DD&end=YYYY-MM-DD`)
- [x] Handle nil granularity (allow any value) for exported entries

### Phase 7: Error Handling
- [x] Validate form input server-side (times, granularity, overlap, timezone, week day)
- [x] Return meaningful error messages via `httperror` + `error.html`
- [x] Handle database errors gracefully via `InternalServerError` + error middleware

### Phase 8: Import Functionality
- [x] CSV import page (`GET /users/:id/entries/import`)
- [x] CSV import with MIME type validation, overlap detection, RFC3339 parsing

## API Routes

| Method | Path | Description |
|--------|------|-------------|
| Any | `/` | Redirect to `/users` |
| GET | `/favicon.ico` | Watch emoji favicon |
| GET | `/users` | List all users |
| GET | `/users/new` | Create user form |
| POST | `/users` | Create user action |
| GET | `/users/{id}` | User overview with weekly summaries |
| GET | `/users/{id}/edit` | Edit user form |
| POST | `/users/{id}` | Update user action |
| POST | `/users/{id}/delete` | Delete user (soft-delete) |
| POST | `/users/{id}/delete-entries` | Delete all entries for user |
| GET | `/users/{id}/entries` | New entry form |
| GET | `/users/{id}/entries/quick` | Natural language entry form (unfinished) |
| POST | `/users/{id}/entries` | Create entry action |
| GET | `/users/{id}/entries/export` | Export timesheet to CSV |
| GET | `/users/{id}/entries/import` | CSV import form |
| POST | `/users/{id}/entries/import` | Import CSV entries |
| GET | `/entries/{id}/edit` | Edit entry form |
| POST | `/entries/{id}` | Update entry action |
| POST | `/users/{id}/entries/delete` | Batch delete entries |

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
start,end,is_paidtimeoff,description
2024-01-01T09:00:00Z,2024-01-01T17:00:00Z,false,Working on feature X
```

### CSV Import Format
```csv
start,end,is_paidtimeoff,description
2024-01-01T09:00:00Z,2024-01-01T17:00:00Z,false,Working on feature X
```

## Notes
- No JavaScript compilation steps
- Server-side HTML rendering with Go templates
- Bootstrap + HTMX via CDN
- No external router - use `http.ServeMux` or manual routing
