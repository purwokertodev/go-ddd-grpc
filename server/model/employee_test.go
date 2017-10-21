package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEmployee(t *testing.T) {

	e := NewEmployee("Wuriyanto", 16, "Jakarta", 60000000.0)

	assert.NotNil(t, e)

}
