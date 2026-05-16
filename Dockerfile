# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -o /out/goshop ./cmd/api

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /out/goshop /app/goshop
COPY config.sample.yaml /app/config.yaml

EXPOSE 8888
ENTRYPOINT ["/app/goshop"]
