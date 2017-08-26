package model

import (
	"github.com/satori/go.uuid"
	"time"
)

type Employee struct {
	Id        uuid.UUID
	Name      string
	Age       int
	Address   string
	Salary    float64
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}

func NewEmployee(name string, age int, address string, salary float64) *Employee {
	return &Employee{
		Id:        uuid.NewV4(),
		Name:      name,
		Age:       age,
		Salary:    salary,
		Address:   address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   0,
	}
}
