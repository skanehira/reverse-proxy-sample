FROM golang:alpine as build_stage

ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLE 0

COPY . /app
WORKDIR /app
RUN go build -o server

FROM alpine:3.10
COPY --from=build_stage /app/server /usr/local/bin

ENTRYPOINT server
