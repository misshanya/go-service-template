FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

ENV GOMODCACHE=/go/pkg/mod
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ENV GOCACHE=/root/.cache/go-build
ENV CGO_ENABLED=0
RUN --mount=type=cache,target="/root/.cache/go-build" \
    --mount=type=cache,target="/go/pkg/mod" \
    go build -ldflags="-s -w" -o service ./cmd/

FROM scratch AS runner

WORKDIR /app

COPY --from=builder /app/service ./service

USER 1000:1000

CMD ["./service"]