FROM golang:1.24-alpine AS builder

COPY go.mod ./
RUN go mod download

COPY . .

RUN apk add --no-cache git ca-certificates

RUN git clone https://github.com/pressly/goose -b v3.25.0

WORKDIR ./goose

RUN go build -tags='no_clickhouse no_libsql no_mssql no_mysql no_sqlite3 no_vertica no_ydb' -o /go/bin/goose ./cmd/goose

WORKDIR ..

RUN GOOS=linux go build -o /main cmd/main.go

FROM alpine AS release-auth

COPY --from=builder /main /
COPY --from=builder /go/migrations /
COPY --from=builder /go/bin/goose /

CMD /goose up && /main