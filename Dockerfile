FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY ./pkg/config/config.sample.yaml ./pkg/config/config.yaml
RUN CGO_ENABLED=0 go build -o goshop ./cmd/api

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/goshop .
COPY --from=builder /app/pkg/config/config.yaml ./pkg/config/config.yaml

EXPOSE 8888
ENTRYPOINT ["/app/goshop"]
