package schema
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestStringConstraint(t*testing.T)  {
	tests := []struct{
		path string
		n interface{}
		schema Schema

		expected []SchemaError
	}{
		{
			path: "a",
			n: "f",
			schema: Schema{
				"minLength": json.Number("2"),
				"pattern": "foo(\\d+)",
			},

			expected: []SchemaError{
				&schemaError{StringMinLengthError, "a"},
				&schemaError{StringPatternError, "a"},
			},
		},
	}

	for _, test :=range tests {
		constraint := NewStringConstraint(test.schema)
		constraint.Validate(test.n, test.path)
		assert.Equal(t, test.expected, constraint.Errors())
	}
}