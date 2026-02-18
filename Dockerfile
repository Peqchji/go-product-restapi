FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/httpserver



FROM alpine:latest
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /server .

RUN apk --no-cache add curl

RUN chown -R appuser:appgroup /app
USER appuser

CMD ["./server"]
