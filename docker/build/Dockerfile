FROM golang:1.22-alpine

# establish the working directory in the container
WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# copy files from current directory to the container's working directory
COPY . .

# setup environment variables
ENV SERVER_PORT=4321

# build the app
RUN go build -v -o ./app .

# expose a port on the container
EXPOSE $PORT

# establish the entry point for the container
CMD ["./app"]
