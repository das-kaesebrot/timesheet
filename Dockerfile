FROM docker.io/library/golang:alpine@sha256:91eda9776261207ea25fd06b5b7fed8d397dd2c0a283e77f2ab6e91bfa71079d AS build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/server

FROM docker.io/library/alpine

ARG APP_WORKDIR="/var/opt/timesheet"
ARG RUN_UID="10020"
ARG RUN_USER="timesheet"

ARG TIMESHEET_WEB_DIR="${APP_WORKDIR}/web"
ENV TIMESHEET_WEB_DIR="${TIMESHEET_WEB_DIR}"
ARG TIMESHEET_DATA_DIR="${APP_WORKDIR}/data"
ENV TIMESHEET_DB_FILE="${TIMESHEET_DATA_DIR}/timesheet.db"

RUN mkdir -pv "${APP_WORKDIR}/data"
RUN addgroup -g ${RUN_UID} ${RUN_USER} && \
    adduser -h ${APP_WORKDIR} -u ${RUN_UID} -G ${RUN_USER} -s /bin/false -D ${RUN_USER} && \
    chown -R ${RUN_USER}:${RUN_USER} "${APP_WORKDIR}"
WORKDIR ${APP_WORKDIR}

COPY --from=build /usr/local/bin/app /usr/local/bin/timesheet
COPY web "${TIMESHEET_WEB_DIR}"
WORKDIR "${APP_WORKDIR}/data"
USER ${RUN_USER}

CMD ["timesheet"]
