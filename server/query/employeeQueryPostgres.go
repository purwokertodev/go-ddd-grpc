package query

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/wuriyanto48/go-ddd-grpc/server/model"
	"os"
)

type QueryEmployeePostgres struct {
	db *sql.DB
}

func NewQueryEmployeePostgres(host string, username string, password string, dbName string) (*QueryEmployeePostgres, error) {
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
	if err != nil {
		return nil, err
	}

	return NewQueryEmployeePostgresDB(db)
}

func NewQueryEmployeePostgresDB(db *sql.DB) (*QueryEmployeePostgres, error) {
	if db == nil {
		return nil, errors.New("Cannot assign nil db as employee query")
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return &QueryEmployeePostgres{db}, nil
}

func (q *QueryEmployeePostgres) GetAll() <-chan QueryResponse {
	output := make(chan QueryResponse)

	go func() {
		defer close(output)

		query := `SELECT "id", "name", "age", "address", "salary", "created_at", "updated_at", "version" FROM "employee"`

		rows, err := q.db.Query(query)

		defer rows.Close()

		if err != nil {
			output <- QueryResponse{Error: err}
			return
		}

		var employees []*model.Employee

		for rows.Next() {
			var em model.Employee

			err = rows.Scan(&em.Id, &em.Name, &em.Age, &em.Address, &em.Salary, &em.CreatedAt, &em.UpdatedAt, &em.Version)

			if err != nil {
				output <- QueryResponse{Error: err}
				return
			}

			employees = append(employees, &em)
		}

		output <- QueryResponse{Result: employees}
	}()

	return output
}
