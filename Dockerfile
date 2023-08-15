FROM harbor.iblog.pro/test/golang:main.golang.1.21.custom as builder
# golang:1.21.0-alpine3.18 as builder

WORKDIR /app
COPY . .


ENV GOLANG_VERSION="1.21.0"
ENV GOPROXY='https://nexus3.iblog.pro/repository/go-proxy/'
ENV GONOSUMDB="https://gitlab.iblog.pro/*"
ENV GONOPROXY="https://gitlab.iblog.pro/*"
#export GOSUMDB='sum.golang.org https://nexus.iblog.pro/repository/golang-sum/'


RUN go mod download
RUN go build -o /app/goshop ./cmd

FROM harbor.iblog.pro/test/alpine:main.scratch.3.18.stage.4
#FROM scratch

WORKDIR /app
COPY --from=builder /app/goshop /app/goshop
COPY ./config/config.sample.yaml ./config/config.yaml

EXPOSE 8888
EXPOSE 8081
ENTRYPOINT ["/app/goshop"]
