FROM golang:1.23.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /server_run ./cmd/server/main.go

EXPOSE 8080

CMD ["/server_run"]