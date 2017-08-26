package repository

import (
	"github.com/satori/go.uuid"
	"github.com/wuriyanto48/go-ddd-grpc/server/model"
)

type EmployeeResponse struct {
	Error    error
	Employee *model.Employee
}

type EmployeeRepository interface {
	Load(id uuid.UUID) <-chan EmployeeResponse
	Save(p *model.Employee) <-chan error
}
