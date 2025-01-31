FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o gopher-tech-blog ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gopher-tech-blog .

# Copy static files and templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# Add a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Change ownership of the working directory to the non-root user
RUN chown -R appuser:appgroup /app

# Copy the .env file
COPY .env.sample .env

# Switch to the non-root user
USER appuser

EXPOSE 8080

# Command to run the application
CMD ["./gopher-tech-blog"]