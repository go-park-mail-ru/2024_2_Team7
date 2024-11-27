#FROM golang:1.23.1
#WORKDIR /app
#COPY go.mod go.sum ./
#RUN go mod download
#COPY . .
#RUN go build -o /server_run ./cmd/server/main.go
#EXPOSE 8080
#CMD ["/server_run"]
FROM debian:bookworm-slim
WORKDIR /app
COPY ./bin/server_service /server_service_run
EXPOSE 8080
CMD ["/server_service_run"]