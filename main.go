package main

import (
	"fmt"
	eq "github.com/wuriyanto48/go-ddd-grpc/server/query"
	repo "github.com/wuriyanto48/go-ddd-grpc/server/repository"
	serv "github.com/wuriyanto48/go-ddd-grpc/server/servers"
	"google.golang.org/grpc"
	"os"
)

const (
	SERVER_CERT = "./cert/server.crt"
	SERVER_KEY  = "./cert/server.key"
	CA          = "./cert/server.crt"
)

func main() {

	dbHost, ok := os.LookupEnv("DB_HOST")

	if !ok {
		fmt.Println("DB_HOST not set in environment variable or test script")
	}

	dbName, ok := os.LookupEnv("DB_NAME")

	if !ok {
		dbName = "gp"
	}

	dbUser, ok := os.LookupEnv("DB_USER")

	if !ok {
		fmt.Println("DB_USER not set in environmet variable or test script")
	}

	dbPassword, ok := os.LookupEnv("DB_PASSWORD")

	if !ok {
		fmt.Println("DB_USER not set in environmet variable or test script")
	}

	repoEmployee, err := repo.NewEmployeeRepoPostgres(dbHost, dbUser, dbPassword, dbName)

	if err != nil {
		fmt.Println("error during create repository employee")
	}

	queryEmployee, err := eq.NewQueryEmployeePostgres(dbHost, dbUser, dbPassword, dbName)

	if err != nil {
		fmt.Println("error during create query employee")
	}

	grpcServer := grpc.NewServer()
	employeeServer := serv.NewEmployeeServer(grpcServer, repoEmployee, queryEmployee, SERVER_CERT, SERVER_KEY, CA)
	err = employeeServer.ServeMutualTLS(8080)
	if err != nil {
		fmt.Sprintf("error create employee grpc server : %s", err)
	}

}
