FROM golang:1.22-bullseye as builder
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -o app -v ./cmd/app

FROM alpine:latest
RUN mkdir /app
WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/app .
COPY --from=builder /src/internal/migrations ./migrations

ENTRYPOINT ["./app"]
