# Build stage
FROM golang:1.24-alpine AS builder

# Set SHELL option -o pipefail
SHELL ["/bin/ash", "-eo", "pipefail", "-c"]

WORKDIR /app
COPY . .

# Consolidate RUN instructions and use --no-cache
RUN go build -o main main.go

# Run stage
FROM alpine:3.22

WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
