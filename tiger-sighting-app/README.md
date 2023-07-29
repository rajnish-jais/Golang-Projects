goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/tiger_hall?sslmode=disable" up
connection string= "postgres://postgres:postgres@localhost:5432/tiger_hall?sslmode=disable"
goose create create_all_tables sql
docker run -d --name tiger-hall-postgis -v /Users/rajnishk/postgresql:/var/lib/postgresql/data -e POSTGRES_PASSWORD=postgres -p 5432:5432 timescale/timescaledb-postgis:latest-pg11