FROM golang:1.23-alpine

# Specify the environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Setting the working directory
WORKDIR /app

# Copy the mod ans sum files to app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Install goose for DB migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy project files to app
COPY . .

# Build app
WORKDIR /app/cmd
RUN go build -o ../expense-manager
WORKDIR /app
# Expose build app to a port
EXPOSE 8080

# Start the app
CMD ["./expense-manager"]