GRPC CRUD EXAMPLE USING POSTGRESS

- Run Server
```shell
DB_USER=postgres DB_PASSWORD=12345 DB_HOST=localhost DB_NAME=gp go run main.go
```

- Run Client

```shell
SERVER_HOST=localhost:8080 go run main.go
```
