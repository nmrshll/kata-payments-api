FROM golang:1.11.1-alpine3.7 AS build
RUN apk --no-cache add gcc g++ git libc-dev make ca-certificates
WORKDIR /go/src/gitlab.com/nmrshll/go-api-postgres
RUN go get -v github.com/rubenv/sql-migrate/...
RUN go get -u -t github.com/volatiletech/sqlboiler
RUN go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql

ENV GO111MODULE=on

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY . .
RUN go install ./... || true

#

FROM alpine:3.7 AS run-base
RUN apk add bash curl ca-certificates
RUN mkdir -p /scripts/ && curl -L https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh -o /scripts/wait-for-it.sh && chmod +x /scripts/wait-for-it.sh
WORKDIR /usr/bin
COPY ./migrations /migrations
COPY ./dbconfig.yml .
COPY ./sqlboiler.toml .
COPY --from=build /go/bin .