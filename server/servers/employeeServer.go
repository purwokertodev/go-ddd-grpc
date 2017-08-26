package servers

import (
	"fmt"
	"github.com/satori/go.uuid"
	pb "github.com/wuriyanto48/go-ddd-grpc/api"
	model "github.com/wuriyanto48/go-ddd-grpc/server/model"
	repo "github.com/wuriyanto48/go-ddd-grpc/server/repository"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
  "net"
)

type EmployeeServer struct {
	server       *grpc.Server
	employeeRepo repo.EmployeeRepository
}

func NewEmployeeServer(server *grpc.Server, repo repo.EmployeeRepository) *EmployeeServer {
	return &EmployeeServer{
		server:       server,
		employeeRepo: repo,
	}
}

func (s *EmployeeServer) Serve(port uint) error {

	address := fmt.Sprintf(":%d", port)

  log.Println("your server is running")

	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	s.server = grpc.NewServer()

	pb.RegisterEmployeeServiceServer(s.server, s)

	err = s.server.Serve(l)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *EmployeeServer) CreateEmployee(ctx context.Context, e *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	em := model.NewEmployee(e.Name, int(e.Age), e.Address, e.Salary)

	err := <-s.employeeRepo.Save(em)
	if err != nil {
		return nil, err
	}

	response := &pb.EmployeeResponse{Id: em.Id.String(), Success: true}
	return response, nil
}

func (s *EmployeeServer) GetEmployee(key *pb.EmployeeFilter, stream pb.EmployeeService_GetEmployeeServer) error {

	id, err := uuid.FromString(key.Key)

	if err != nil {
		fmt.Println("uuid error")
	}

	emRes := <-s.employeeRepo.Load(id)
	if emRes.Error != nil {
		return emRes.Error
	}

	em := emRes.Employee

	res := &pb.EmployeeRequest{
		Id:        em.Id.String(),
		Name:      em.Name,
		Age:       int32(em.Age),
		Address:   em.Address,
		Salary:    em.Salary,
		CreatedAt: em.CreatedAt.String(),
		UpdatedAt: em.UpdatedAt.String(),
		Version:   int32(em.Version),
	}

	if err = stream.Send(res); err != nil {
		return err
	}

  return nil
}
