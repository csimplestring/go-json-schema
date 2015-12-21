package schema

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayConstraintUniqueItems(t *testing.T) {

	tests := []struct {
		input          []interface{}
		expectedErrors []SchemaError
	}{
		{
			input: []interface{}{
				json.Number("1"),
				json.Number("1"),
			},
			expectedErrors: nil,
		},
		{
			input: []interface{}{
				"a", "a",
			},
			expectedErrors: nil,
		},
		{
			input: []interface{}{
				map[string]interface{}{},
				map[string]interface{}{},
			},
			expectedErrors: nil,
		},
		{
			input: []interface{}{
				map[string]interface{}{
					"a": 1,
					"b": 2,
				},
				map[string]interface{}{
					"a": 1,
					"b": 2,
				},
			},
			expectedErrors: nil,
		},
		{
			input: []interface{}{
				map[string]interface{}{
					"a": 1,
					"b": 2,
				},
				map[string]interface{}{
					"a": 1,
					"b": 2,
					"c": 3,
				},
			},
			expectedErrors: []SchemaError{
				&schemaError{ArrayUniqueItemError, "a[1]"},
			},
		},
	}

	c := NewArrayConstraint(nil)
	path := "a"
	for _, test := range tests {
		c.validateUniqueItem(test.input, path)
		assert.Equal(t, test.expectedErrors, c.Errors())
	}
}
