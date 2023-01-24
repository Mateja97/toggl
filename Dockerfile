FROM golang:1.19-alpine

RUN apk add build-base

WORKDIR /app

COPY go.* ./

COPY . ./


RUN go build -o /toggl

CMD ["/toggl" ]