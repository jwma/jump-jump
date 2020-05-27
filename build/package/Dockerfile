# build stage
FROM golang:1.13-alpine as builder

LABEL maintainer="m.mjw.ma@gmail.com"

ENV GOPROXY=https://goproxy.io
ENV CGO_ENABLED=0

WORKDIR /app_src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o apiserver ./cmd/apiserver \
    && go build -o landingserver ./cmd/landingserver \
    && go build -o createuser ./cmd/createuser \
    && go build -o migraterequesthistory ./cmd/migraterequesthistory

# final stage
FROM alpine

ARG EXPOSE_PORT=8080
ENV CGO_ENABLED=0

EXPOSE $EXPOSE_PORT
WORKDIR /app
COPY --from=builder /app_src/apiserver /app_src/landingserver /app_src/createuser /app_src/migraterequesthistory ./
COPY --from=builder /app_src/web ./web
ENTRYPOINT ["./apiserver"]
