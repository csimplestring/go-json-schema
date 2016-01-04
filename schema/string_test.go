package schema

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringConstraint(t *testing.T) {
	tests := []struct {
		path   string
		n      interface{}
		schema Schema

		expected []ValidationError
	}{
		{
			path: "a",
			n:    "f",
			schema: Schema{
				"minLength": json.Number("2"),
				"pattern":   "foo(\\d+)",
			},

			expected: []ValidationError{
				&validationError{StringMinLengthError, "a"},
				&validationError{StringPatternError, "a"},
			},
		},
	}

	for _, test := range tests {
		constraint := NewStringConstraint(test.schema)
		constraint.Validate(test.n, test.path)
		assert.Equal(t, test.expected, constraint.Errors())
	}
}
