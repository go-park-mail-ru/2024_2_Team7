FROM golang:1.23.1

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Build
RUN go build -o /csat_run ./cmd/csat/main.go

#
EXPOSE 50055

# Run
CMD ["/csat_run"]