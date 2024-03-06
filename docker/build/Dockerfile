FROM golang:1.22-apline

# establish the working directory in the container
WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# copy files from current directory to the container's working directory
COPY . .

# setup environment variables
ENV DB_HOST="host.docker.internal"
ENV DB_PORT=5432
ENV DB_USER="postgres"
ENV DB_PASSWORD="password"
ENV DB_NAME="go_test"
ENV SERVER_PORT=4321

# build the app
RUN go build -v -o ./app .

# expose a port on the container
EXPOSE $PORT

# establish the entry point for the container
CMD ["./app"]
