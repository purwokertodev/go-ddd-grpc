package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/wuriyanto48/go-ddd-grpc/api"
	"google.golang.org/grpc"
)

func main() {

	serverHost, ok := os.LookupEnv("SERVER_HOST")

	if !ok {
		fmt.Println("SERVER HOST not set in environment variable or test script")
	}

	// dial without tls or something
	conn, err := grpc.Dial(serverHost, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewEmployeeServiceClient(conn)
	//find
	id := "9e10c3fc-f3dd-4dca-b7e9-8f1f1b038dcc"

	GetEmployee(client, id)

	// var salary float64 = 15000000.0
	// em := &pb.EmployeeRequest{
	// 	Name:    "Bimo",
	// 	Age:     30,
	// 	Address: "Jakarta Barat",
	// 	Salary:  salary,
	// }
	//
	// CreateEmployee(client, em)
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
