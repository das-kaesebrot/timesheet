FROM docker.io/library/golang:alpine@sha256:4cb7ac979db5fcc41cae44b2227ba5ab8a51e8807f40d9ba4dee20a0ad960b5b AS build

ARG VERSION="dev-docker"
WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# https://jerrynsh.com/3-easy-ways-to-add-version-flag-in-go/
RUN go build -v -ldflags "-X 'main.Version=${VERSION}'" -o /usr/local/bin/app ./cmd/server/main.go

FROM docker.io/library/alpine@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6

ARG APP_WORKDIR="/var/opt/timesheet"
ARG RUN_UID="10020"
ARG RUN_USER="timesheet"

ARG TIMESHEET_DATA_DIR="${APP_WORKDIR}/data"
ENV TIMESHEET_DB_FILE="${TIMESHEET_DATA_DIR}/timesheet.db"

RUN apk add --no-cache tzdata
RUN mkdir -pv "${APP_WORKDIR}/data"
RUN addgroup -g ${RUN_UID} ${RUN_USER} && \
    adduser -h ${APP_WORKDIR} -u ${RUN_UID} -G ${RUN_USER} -s /bin/false -D ${RUN_USER} && \
    chown -R ${RUN_USER}:${RUN_USER} "${APP_WORKDIR}"
WORKDIR ${APP_WORKDIR}

COPY --from=build /usr/local/bin/app /usr/local/bin/timesheet
WORKDIR "${APP_WORKDIR}/data"
USER ${RUN_USER}

CMD ["timesheet"]
