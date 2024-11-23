FROM golang:1.23.1

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download


# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
#COPY /cmd/auth ./
#COPY /internal/auth ./
#COPY /internal/gateway/auth ./
#COPY /internal/gateway/utils ./
#COPY /internal/http/auth ./
#COPY /internal/logger ./
#COPY /internal/middleware ./
#COPY /internal/models ./
#COPY /internal/repository ./
COPY . .


# Build
RUN go build -o /auth_service_run ./cmd/auth/main.go

#
EXPOSE 50051

# Run
CMD ["/auth_service_run"]
#CMD ["ls -la"]