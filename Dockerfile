# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o diffengine ./cmd/diffengine

# Final stage
FROM scratch

COPY --from=builder /app/diffengine /diffengine

ENTRYPOINT ["/diffengine"]
