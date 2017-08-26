package repository

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"github.com/wuriyanto48/go-ddd-grpc/server/model"
	"os"
	"time"
)

type EmployeeRepositoryPostgres struct {
	Db *sql.DB
}

func NewEmployeeRepoPostgres(host string, username string, password string, dbName string) (*EmployeeRepositoryPostgres, error) {

	_, devOk := os.LookupEnv("NO_SSL")

	var connStr string
	if devOk {
		connStr = fmt.Sprintf("host='%s' port=5432 dbname='%s' user='%s' password='%s' sslmode=disable",
			host, dbName, username, password)
	} else {
		connStr = fmt.Sprintf("host='%s' port=5432 dbname='%s' user='%s' password='%s' sslmode=require",
			host, dbName, username, password)
	}
	db, err := sql.Open("postgres", connStr)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	if err != nil {
		return nil, err
	}

	return NewEmployeeRepoPostgresDB(db)
}

func NewEmployeeRepoPostgresDB(db *sql.DB) (*EmployeeRepositoryPostgres, error) {
	if db == nil {
		return nil, errors.New("Cannot assign nil db as order repository")
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return &EmployeeRepositoryPostgres{db}, nil

}

func (repo *EmployeeRepositoryPostgres) Load(id uuid.UUID) <-chan EmployeeResponse {
	result := make(chan EmployeeResponse)
	go func() {
		defer close(result)
		employee := new(model.Employee)
		query := "SELECT * FROM EMPLOYEE WHERE id=$1"
		prep, err := repo.Db.Prepare(query)
		defer prep.Close()
		if err != nil {
			result <- EmployeeResponse{Error: err}
			return
		}
		err = prep.QueryRow(id).Scan(&employee.Id, &employee.Name, &employee.Age, &employee.Address, &employee.Salary, &employee.CreatedAt, &employee.UpdatedAt, &employee.Version)
		if err != nil {
			result <- EmployeeResponse{Error: err}
			return
		}
		result <- EmployeeResponse{Error: nil, Employee: employee}
	}()
	return result
}

func (repo *EmployeeRepositoryPostgres) Save(e *model.Employee) <-chan error {
	result := make(chan error)
	go func() {

		defer close(result)
		if e == nil {
			result <- errors.New("Employee required")
			return
		} else if e.Name == "" {
			result <- errors.New("Employee required")
			return
		}

		tx, err := repo.Db.Begin()

		if err != nil {
			result <- err
			return
		}

		readstmt, err := tx.Prepare(`SELECT "version" FROM "employee" WHERE "id"=$1`)

		if err != nil {
			tx.Rollback()
			result <- err
			return
		}

		defer readstmt.Close()

		var version int
		err = readstmt.QueryRow(e.Id).Scan(&version)

		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			result <- err
			return
		}

		if version > e.Version {
			tx.Rollback()
			result <- errors.New("There's conflict during save")
			return
		}
		query := `INSERT INTO employee("id", "name", "age", "address", "salary", "created_at", "updated_at", "version")
    VALUES($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT("id") DO UPDATE SET "name" = $2, "age" = $3, "address" = $4, "salary" = $5,
		"created_at" = $6, "updated_at" = $7, "version" = $8`
		prep, err := repo.Db.Prepare(query)
		defer prep.Close()
		if err != nil {
			tx.Rollback()
			result <- err
			return
		}

		e.Version += 1
		e.UpdatedAt = time.Now()

		_, err = prep.Exec(e.Id, e.Name, e.Age, e.Address, e.Salary, e.CreatedAt, e.UpdatedAt, e.Version)
		if err != nil {
			tx.Rollback()
			result <- err
			return
		}

		tx.Commit()

		result <- err
	}()
	return result

}
