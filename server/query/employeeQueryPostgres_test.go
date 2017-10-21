package query

import (
	"github.com/stretchr/testify/assert"
	"github.com/wuriyanto48/go-ddd-grpc/server/model"
	sqlMock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
	"time"
)

func TestEmployeeQuery(t *testing.T) {

	t.Run("TestGetAll", func(t *testing.T) {

		db, mock, _ := sqlMock.New()

		defer db.Close()

		rows := sqlMock.NewRows([]string{"id", "name", "age", "address", "salary", "created_at", "updated_at", "version"}).
			AddRow("d4515a17-f7c0-42d3-99e0-c999776e873c", "Wuriyanto", 16, "Jakarta", 68000000.0, time.Now(), time.Now(), 1)

		query := `SELECT "id", "name", "age", "address", "salary", "created_at", "updated_at", "version" FROM "employee"`

		mock.ExpectQuery(query).WillReturnRows(rows)

		q := &QueryEmployeePostgres{db}

		employeeResult := <-q.GetAll()

		assert.NoError(t, employeeResult.Error)

		employees, ok := employeeResult.Result.([]*model.Employee)

		assert.True(t, ok)

		assert.Len(t, employees, 1)

	})

}
