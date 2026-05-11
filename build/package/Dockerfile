# frontend build stage
FROM node:22-alpine AS frontend-builder

WORKDIR /web-ui
COPY web-ui/package.json web-ui/package-lock.json ./
RUN npm ci
COPY web-ui/ ./
RUN npm run build

# backend build stage
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
    && go build -o createsuperuser ./cmd/createsuperuser

# final stage
FROM alpine:3.21

ENV CGO_ENABLED=0

WORKDIR /app
COPY --from=builder /app_src/apiserver /app_src/landingserver /app_src/createsuperuser ./
COPY --from=builder /app_src/docs ./docs
COPY --from=frontend-builder /web-ui/dist ./web-ui/dist
EXPOSE 8080
ENTRYPOINT ["./apiserver"]
