package servers

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	pb "github.com/wuriyanto48/go-ddd-grpc/api"
	model "github.com/wuriyanto48/go-ddd-grpc/server/model"
	repo "github.com/wuriyanto48/go-ddd-grpc/server/repository"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
)

type EmployeeServer struct {
	server               *grpc.Server
	employeeRepo         repo.EmployeeRepository
	serverCert           string
	serverKey            string
	certificateAuthority string
}

func NewEmployeeServer(server *grpc.Server, repo repo.EmployeeRepository, serverCert, serverKey, ca string) *EmployeeServer {
	return &EmployeeServer{
		server:               server,
		employeeRepo:         repo,
		serverCert:           serverCert,
		serverKey:            serverKey,
		certificateAuthority: ca,
	}
}

func (s *EmployeeServer) ServeMutualTLS(port uint) error {

	address := fmt.Sprintf(":%d", port)

	log.Println("your server is running")

	//get from disk
	certificate, err := tls.LoadX509KeyPair(s.serverCert, s.serverKey)
	if err != nil {
		return fmt.Errorf("cannot load server key pair : %s", err)
	}

	//create certificate pool from CA
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(s.certificateAuthority)
	if err != nil {
		return fmt.Errorf("cannot load certificate authority : %s", err)
	}

	//append the client certificate from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return errors.New("failed append client cert")
	}

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	})

	s.server = grpc.NewServer(grpc.Creds(creds))

	pb.RegisterEmployeeServiceServer(s.server, s)

	err = s.server.Serve(l)

	if err != nil {
		return err
	}

	return nil
}

func (s *EmployeeServer) ServeTLS(port uint) error {

	address := fmt.Sprintf(":%d", port)

	log.Println("your server is running")

	creds, err := credentials.NewServerTLSFromFile(s.serverCert, s.serverKey)
	if err != nil {
		return fmt.Errorf("Cannot load TLS keys : %s", err)
	}

	s.server = grpc.NewServer(grpc.Creds(creds))

	pb.RegisterEmployeeServiceServer(s.server, s)

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	err = s.server.Serve(l)

	if err != nil {
		return err
	}

	return nil
}

// server insecure server/ no server side encryption
func (s *EmployeeServer) Serve(port uint) error {

	address := fmt.Sprintf(":%d", port)

	log.Println("your server is running")

	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.server = grpc.NewServer()

	pb.RegisterEmployeeServiceServer(s.server, s)

	err = s.server.Serve(l)

	if err != nil {
		return err
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
		return err
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
