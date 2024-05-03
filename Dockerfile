FROM golang:1.20-bookworm as builder

WORKDIR /app

ADD . /app

RUN go build -o bin/student-api .

FROM debian:bookworm-slim

COPY --from=builder /app/bin/student-api .

COPY ./migrations /app/migrations

ENV MIGRATIONS_PATH="/app/migrations"

EXPOSE 8080

CMD ["./student-api"]
