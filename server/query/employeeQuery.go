package query

type QueryResponse struct {
	Error  error
	Result interface{}
}

type EmployeeQuery interface {
	GetAll() <-chan QueryResponse
}
