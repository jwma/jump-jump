# build stage
FROM golang:1.24-alpine AS builder

LABEL maintainer="m.mjw.ma@gmail.com"

ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0

WORKDIR /app_src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o apiserver ./cmd/apiserver \
    && go build -o landingserver ./cmd/landingserver \
    && go build -o createuser ./cmd/createuser

# final stage
FROM alpine:3.21

ENV CGO_ENABLED=0

WORKDIR /app
COPY --from=builder /app_src/apiserver /app_src/landingserver /app_src/createuser ./
COPY --from=builder /app_src/web ./web
EXPOSE 8080
ENTRYPOINT ["./apiserver"]
