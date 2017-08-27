package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/wuriyanto48/go-ddd-grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	SERVER_CERT = "../cert/server.crt"
	SERVER_KEY  = "../cert/server.key"
	CA          = "../cert/server.crt"
	SERVER_NAME = "localhost"
)

func main() {

	serverHost, ok := os.LookupEnv("SERVER_HOST")

	if !ok {
		fmt.Println("SERVER HOST not set in environment variable or test script")
	}

	client, err := clientWithMutualTLS(serverHost)
	if err != nil {
		fmt.Println(err)
	}

	//find all

	GetEmployeeAll(*client)

	//find
	// id := "1bd06dd0-fa00-40a8-8bfa-ddef1342aec5"
	//
	// GetEmployee(*client, id)

	//create

	// var salary float64 = 87000000.0
	// em := &pb.EmployeeRequest{
	// 	Name:    "Andree",
	// 	Age:     35,
	// 	Address: "Bogor",
	// 	Salary:  salary,
	// }
	//
	// CreateEmployee(*client, em)
}

func clientWithInsecure(serverHost string) (*pb.EmployeeServiceClient, error) {
	conn, err := grpc.Dial(serverHost, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	//defer conn.Close()

	client := pb.NewEmployeeServiceClient(conn)

	return &client, nil
}

func clientWithTLS(serverHost string) (*pb.EmployeeServiceClient, error) {

	//create client TLS
	creds, err := credentials.NewClientTLSFromFile(SERVER_CERT, "")
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(serverHost, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	//defer conn.Close()

	client := pb.NewEmployeeServiceClient(conn)

	return &client, nil
}

func clientWithMutualTLS(serverHost string) (*pb.EmployeeServiceClient, error) {

	//get from disk
	certificate, err := tls.LoadX509KeyPair(SERVER_CERT, SERVER_KEY)
	if err != nil {
		return nil, fmt.Errorf("cannot load server key pair : %s", err)
	}

	//create certificate pool from CA
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(CA)
	if err != nil {
		return nil, fmt.Errorf("cannot load certificate authority : %s", err)
	}

	//append the client certificate from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, errors.New("failed append client cert")
	}

	creds := credentials.NewTLS(&tls.Config{
		ServerName:   SERVER_NAME,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	conn, err := grpc.Dial(serverHost, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	//defer conn.Close()

	client := pb.NewEmployeeServiceClient(conn)

	return &client, nil
}

func CreateEmployee(client pb.EmployeeServiceClient, e *pb.EmployeeRequest) {

	resp, err := client.CreateEmployee(context.Background(), e)

	if err != nil {
		log.Fatalf("Could not create employee: %v", err)
	}
	if resp.Success {
		log.Printf("A new employee has been added with id: %d", resp.Id)
	}

}

func GetEmployee(client pb.EmployeeServiceClient, id string) {
	filter := &pb.EmployeeFilter{Key: id}
	resStream, err := client.GetEmployee(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	employee, err := resStream.Recv()

	if err == io.EOF {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatalf("%v.GetEmployee(_) = _, %v", client, err)
	}

	fmt.Println(employee)
}

func GetEmployeeAll(client pb.EmployeeServiceClient) {
	filter := &pb.EmployeeFilter{}
	resStream, err := client.GetAll(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	result, err := resStream.Recv()

	if err == io.EOF {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatalf("%v.GetEmployee(_) = _, %v", client, err)
	}

	for _, e := range result.Employees {
		fmt.Println(e)
	}
}
