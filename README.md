GRPC CRUD EXAMPLE USING POSTGRESS

# Creating SSL/TLS Certificates

 - server.key : a private RSA key to sign and authenticate the public key
	```shell
	openssl genrsa -out server.key 2048
	```

 - server.pem/ server.crt : self-signed x.509 public keys for distribution
	```shell
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

	```
 - server.csr: a certificate signing request to acces the CA(Certificate Authority)
	```shell
	openssl req -new -sha256 -key server.key -out server.csr

	```

	```shell
	openssl x509 -req -sha256 -in server.csr -signkey server.key -out server.crt -days 3650

	```
# App
- Run Server
```shell
DB_USER=postgres DB_PASSWORD=12345 DB_HOST=localhost DB_NAME=gp go run main.go
```

- Run Client

```shell
SERVER_HOST=localhost:8080 go run main.go
```
