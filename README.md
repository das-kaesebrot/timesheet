# timesheet

timesheet is a web-based time tracking application that allows you to manage users and their timesheet entries with support for timezone-aware weekly summaries, CSV export, and natural language time input.

## Features

- [x] multi-user with configurable granularity and weekly work time per user
- [x] timesheet entries with start/end times and descriptions
- [x] timezone-aware entries
- [x] weekly summaries with time logged and weekly diff
- [x] CSV export
- [ ] CSV import

## Build and run

Clone the repository and build the binary:

```bash
cd timesheet
go build -o timesheet ./cmd/server
./timesheet
```

The server will start on `[::]:8080` by default and create a `timesheet.db` SQLite database in the current directory.

## Configuration

The application can be configured via environment variables.

### Environment variables

| Variable | Description | Default | Required? |
|----------|-------------|---------|-----------|
| `TIMESHEET_HOST` | Host address to bind to | `[::]` | No |
| `TIMESHEET_PORT` | Port to listen on | `8080` | No |

## Open Source License Attribution

This application uses Open Source components. You can find the source code of their open source projects along with license information below. We acknowledge and are grateful to these developers for their contributions to open source.

### [GORM](https://github.com/go-gorm/gorm)

- Copyright (c) 2013-present [jinzhu](https://github.com/jinzhu) and contributors
- [MIT License](https://github.com/go-gorm/gorm/blob/master/LICENSE)

### [GORM SQLite Driver](https://github.com/go-gorm/sqlite)

- Copyright (c) 2013-present [jinzhu](https://github.com/jinzhu) and contributors
- [MIT License](https://github.com/go-gorm/sqlite/blob/master/License)

### [uuid](https://github.com/google/uuid)

- Copyright (c) 2009,2014 Google Inc
- [BSD33-Clause License](https://github.com/google/uuid/blob/master/LICENSE)

### [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)

- Copyright (c) 2009 The Go Authors
- [BSD33-Clause License](https://go.googlesource.com/crypto/+//master/LICENSE)

### [golang.org/x/text](https://pkg.go.dev/golang.org/x/text)

- Copyright (c) 2009 The Go Authors
- [BSD-3-Clause License](https://go.googlesource.com/text/+/master/LICENSE)
