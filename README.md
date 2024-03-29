# GO + MongoDB + Docker REST API Template

A template of a REST API written in Go with a development environment run using Docker.

## Requirements:
- [Docker](https://www.docker.com/products/docker-desktop/)
    - Docker Compose
- [Go](https://go.dev/)
- [Make](https://www.gnu.org/software/make/)

## Local Dev Environment
### Running the server
```
make start-server
```
- This will use Docker Compose to spool up all the required Docker containers
- The Go server will be run using [air](https://github.com/cosmtrek/air) which automatically restarts the server when code changes are made
- Server environment variables are located in the `docker-compose.yaml` file.

# Helpful MongoDB Resources
- [Golang MongoDB Driver Quick Reference](https://www.mongodb.com/docs/drivers/go/current/quick-reference/#std-label-golang-quick-reference)
