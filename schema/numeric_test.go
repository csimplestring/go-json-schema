package schema
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestNumericConstraint(t*testing.T) {

	tests := []struct{
		path string
		n interface{}
		schema Schema

		expected []SchemaError
	} {
		{
			path: "a",
			n: json.Number("100"),
			schema: Schema{
				"multipleOf" : json.Number("3"),
			},

			expected: []SchemaError {
				&schemaError{NumericMultipleOfError, "a"},
			},
		},
		{
			path: "a",
			n: json.Number("80"),
			schema: Schema{
				"multipleOf" : json.Number("3.1"),
				"minimum": json.Number("90"),
			},

			expected: []SchemaError {
				&schemaError{NumericMultipleOfError, "a"},
				&schemaError{NumericMinimumError, "a"},
			},
		},
		{
			path: "a",
			n: "1",
			schema: Schema{
				"maximum": json.Number("100.4"),
				"exclusiveMaximum": true,
			},

			expected: []SchemaError {
				&schemaError{NumericTypeMismatchError, "a"},
			},
		},
	}

	for _, test :=range tests {
		constraint := NewNumericConstraint(test.schema)
		constraint.Validate(test.n, test.path)
		assert.Equal(t, test.expected, constraint.Errors())
	}
}
