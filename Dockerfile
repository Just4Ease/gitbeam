FROM golang:1.21 AS builder
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . .


RUN GIT_TERMINAL_PROMPT=1 CGO_ENABLED=1 GOOS=linux go build -o app -a -ldflags '-linkmode external -extldflags "-static"' .

#FROM debian:latest
FROM alpine:3.13
# convert build-arg to env variables
#RUN apk add --no-cache tzdata
#RUN apk add --no-cache gcc musl-dev sqlite-dev
#ENV TZ Africa/Lagos
RUN mkdir -p /svc/
COPY --from=builder /go/src/app/app /svc/

WORKDIR /svc/

EXPOSE 8080

CMD ["./app"]
