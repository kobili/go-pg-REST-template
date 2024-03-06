build:
	docker build -t go-rest -f ./docker/dev/Dockerfile ./server

start-server:
	docker-compose up

make-migration:
	docker run -v $(shell pwd)/migrations/sql:/migrations migrate/migrate create -ext sql -dir ./migrations -seq ${migration_name}
mm: make-migration

migrate:
	docker run  --rm -v $(shell pwd)/migrations/sql:/migrations --network go-test-rest_go-server migrate/migrate -path=/migrations/ -database postgres://postgres:password@postgres:5432/go_test?sslmode=disable up
m: migrate

count ?= 1 # use -all to go all the way down
migrate-down:
	docker run --rm -v $(shell pwd)/migrations/sql:/migrations --network go-test-rest_go-server migrate/migrate -path=/migrations/ -database postgres://postgres:password@postgres:5432/go_test?sslmode=disable down ${count}
