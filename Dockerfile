FROM golang:alpine AS builder
ENV GOOS=linux \
    GOARCH=amd64

ARG TOKEN
ARG PRODUCTION
ARG FIREBASE_CONFIG
ARG FIREBASE_PROJECT_ID
ARG DATADOG_API_KEY
ENV PRODUCTION=TRUE
ENV TOKEN=$TOKEN
ENV FIREBASE_CONFIG=$FIREBASE_CONFIG
ENV FIREBASE_PROJECT_ID=$FIREBASE_PROJECT_ID
ENV DATADOG_API_KEY=$DATADOG_API_KEY

WORKDIR /build

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o ./main.go

FROM alpine
WORKDIR /dist
COPY --from=builder /build/main /dist


ENTRYPOINT ./main