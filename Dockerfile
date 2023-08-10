FROM golang:1.21.0-alpine3.18 as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/goshop ./cmd

FROM scratch 

WORKDIR /app
COPY --from=builder /app/goshop /app/goshop
COPY ./config/config.sample.yaml ./config/config.yaml

EXPOSE 8888
ENTRYPOINT ["/app/goshop"]
