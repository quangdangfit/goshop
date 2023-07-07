FROM golang:1.20.5

MAINTAINER quangdp<quangdangfit@gmail.com>

WORKDIR /app
COPY . ./
RUN go mod download

COPY ./config/config.sample.yaml ./config/config.yaml
RUN go build -o /app/goshop

EXPOSE 8888
ENTRYPOINT ["/app/goshop"]
