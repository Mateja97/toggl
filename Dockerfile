FROM golang:1.19-alpine

ENV PORT ${PORT}

EXPOSE ${PORT} ${PORT}

RUN apk add build-base

WORKDIR /app

COPY go.* ./

COPY . ./


RUN go build -o /toggl

CMD ["/toggl" ]