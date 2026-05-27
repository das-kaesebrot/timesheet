FROM docker.io/library/golang:alpine@sha256:91eda9776261207ea25fd06b5b7fed8d397dd2c0a283e77f2ab6e91bfa71079d AS build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /usr/local/bin/app ./cmd/server/main.go

FROM scratch

COPY contrib/passwd /etc/passwd
COPY --from=build /usr/local/bin/app /timesheet

USER timesheet

CMD ["./timesheet"]
