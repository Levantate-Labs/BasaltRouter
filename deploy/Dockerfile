# Shared multi-stage Dockerfile for BasaltRouter services.
# Build with: docker build -f deploy/Dockerfile --build-arg SERVICE=gateway .

ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE=gateway
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/service ./cmd/${SERVICE}

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /out/service /app/service
COPY migrations /app/migrations

USER nobody
EXPOSE 8080

ENTRYPOINT ["/app/service"]
