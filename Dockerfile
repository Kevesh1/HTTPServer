# syntax=docker/dockerfile:1

FROM golang:1.21.3

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /httpserver

ENV PORT=8080

EXPOSE 8080

CMD ["/httpserver"]