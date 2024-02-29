# Kobe's Migration Tool

A migration tool which uses Go and Make scripts

## Requirements
- [golang-migrate](https://github.com/golang-migrate/migrate) v4.16.2+
- Go v1.20.5+
- Make

## Setup
Create a `.env` file at the root directory with the following values:
```
POSTGRES_HOST=
POSTGRES_PORT=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DATABASE=
```

## Usage
### Create `up` and `down` migration files
```
make make-migration
```

This will generate two numbered `.sql` files. Write your migration code in the `.up.sql` file and the code needed to reverse the migration in the `.down.sql` file

### Run migrations
```
make migrate
```

### Revert migrations
```
make migrate-down
```
- Will revert one migration at a time
