package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSchemaError(t *testing.T) {
	e := newError(TypesNotMatchError, "a")
	assert.Equal(t, "Error: not match one of types, Path: a", e.Error())
}
