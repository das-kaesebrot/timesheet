FROM docker.io/library/golang:alpine@sha256:91eda9776261207ea25fd06b5b7fed8d397dd2c0a283e77f2ab6e91bfa71079d AS build

WORKDIR /usr/src/app

ENV CGO_ENABLED=1

# Add gcc and musl-dev
# https://wiki.alpinelinux.org/wiki/GCC
RUN apk add --no-cache \
    gcc \
    musl-dev

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/server/main.go

FROM scratch

COPY --from=build /usr/local/bin/app /usr/local/bin/timesheet
COPY contrib/passwd /etc/passwd

USER timesheet

CMD ["timesheet"]
