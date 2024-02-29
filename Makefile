start-server:
	go run ./server

make-migration:
	migrate create -ext sql -dir ./migrations/sql -seq ${migration_name}
mm: make-migration

migrate:
	cd migrations; go run .
m: migrate

migrate-down:
	cd migrations; go run . -dir=down
