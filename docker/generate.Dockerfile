FROM golang:1.11.1-alpine3.7 AS build
RUN apk add gcc g++ git libc-dev make ca-certificates curl bash postgresql-client
RUN go get -v github.com/rubenv/sql-migrate/...
RUN go get -u -t github.com/volatiletech/sqlboiler
RUN go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql
RUN mkdir -p /scripts/ && curl -L https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh -o /scripts/wait-for-it.sh && chmod +x /scripts/wait-for-it.sh

ENV GO111MODULE=on

WORKDIR /usr/local/codegen
COPY ./migrations /migrations
COPY ./dbconfig.yml .
COPY ./sqlboiler.toml .