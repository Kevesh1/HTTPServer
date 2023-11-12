# syntax=docker/dockerfile:1

FROM golang:1.21.3

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

# Go can't call C code
# OS: Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /httpserver

ENV PORT=8080

EXPOSE 8080

CMD ["/httpserver"]