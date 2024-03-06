FROM golang:1.22-alpine

# establish the working directory in the container
WORKDIR /app

# download air
RUN go install github.com/cosmtrek/air@latest

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY .air.toml .

# setup environment variables
ENV DB_HOST="host.docker.internal"
ENV DB_PORT=3432
ENV DB_USER="postgres"
ENV DB_PASSWORD="password"
ENV DB_NAME="go_test"
ENV SERVER_PORT=4321

# Entry point
CMD ["air", "."]
