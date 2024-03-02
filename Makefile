build:
	docker build -t go-rest ./server

start-server:
	docker run -p 4321:4321 --name go-rest-server go-rest

make-migration:
	migrate create -ext sql -dir ./migrations/sql -seq ${migration_name}
mm: make-migration

migrate:
	cd migrations; go run .
m: migrate

migrate-down:
	cd migrations; go run . -dir=down
