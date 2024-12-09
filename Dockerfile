# Build stage
FROM golang:1.23-alpine3.20 AS builder

# Set SHELL option -o pipefail
SHELL ["/bin/ash", "-eo", "pipefail", "-c"]

WORKDIR /app
COPY . .

# Consolidate RUN instructions and use --no-cache
RUN apk add --no-cache curl=8.11.0-r2 && \
  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz && \
  go build -o main main.go

# Run stage
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
