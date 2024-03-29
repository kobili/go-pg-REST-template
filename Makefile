build:
	docker build -t go-rest -f ./docker/dev/Dockerfile ./server

start-server:
	docker compose up
