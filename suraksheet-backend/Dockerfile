# Dockerfile for Go API Backend
FROM golang:1.22.5

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o bin/suraksheet cmd/main.go

COPY start.sh .

EXPOSE 8080

CMD ["./start.sh"]
