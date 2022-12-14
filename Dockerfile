# syntax=docker/dockerfile:1

## Build
FROM golang:1.19 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./
COPY internal ./
COPY pkg ./

RUN go build -o /promqlinter ./cmd/promqlinter

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /promqlinter /promqlinter

USER nonroot:nonroot

ENTRYPOINT ["/promqlinter"]