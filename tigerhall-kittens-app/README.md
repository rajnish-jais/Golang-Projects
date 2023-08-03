## [Tigerhall-Kittens-App]

Technology: Golang

### SetUP

- Postgres and RabbitMQ has been used for this project. 
- Add your configuration for Postgres and RabbitMq in config/local/server.yml file in order established the connection for dependency of this project.

### Run Migration
- goose -dir migrations postgres "postgres://postgres:postgres@localhost:5432/tiger_hall?sslmode=disable" up
- `postgres://postgres:postgres@localhost:5432/tiger_hall?sslmode=disable` this a connection string and need to be updated according to configuration of postgres connection.
### Start Server
- Go to cmd file and execute `go run main.go`
### Run Tests
- execute `go test file/folder name`

Note: You can find all the endpoint in postman collection in file [Tigerhall_Kittens.postman_collection.json](Tigerhall_Kittens.postman_collection.json). 