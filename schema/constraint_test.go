package schema
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)


func TestTypeConstraint(t *testing.T) {
	tests := []struct{
		schema Schema
		value interface{}
		expected []SchemaError
	}{
		{
			schema: Schema{
				"type": "integer",
			},
			value: 1,
			expected: nil,
		},
		{
			schema: Schema{
				"type": []interface{}{"integer", "string"},
			},
			value: "str",
			expected: nil,
		},
		{
			schema: Schema{
				"type": []interface{}{"integer"},
			},
			value: "str",
			expected: []SchemaError{
				newError(TypesNotMatchError, "a"),
			},
		},
		{
			schema: Schema{
				"type": "object",
			},
			value: "str",
			expected: []SchemaError{
				newError(TypeNotMatchError, "a"),
			},
		},
	}

	for _, test := range tests {
		c := NewBaseConstraint(test.schema)
		c.validateType(test.value, "a")
		assert.Equal(t, test.expected, c.Errors())
	}
}

func TestEnumConstraint(t *testing.T) {
	tests := []struct{
		schema Schema
		value interface{}
		expected []SchemaError
	}{
		{
			schema: Schema{
				"enum": []interface{}{json.Number("1")},
			},
			value: json.Number("1"),
			expected: nil,
		},
		{
			schema: Schema{
				"enum": []interface{}{
					json.Number("1"),
					json.Number("1.2"),
					"str",
				},
			},
			value: json.Number("1.2"),
			expected: nil,
		},
		{
			schema: Schema{
				"enum": []interface{}{
					json.Number("1"),
					json.Number("1.2"),
					"str",
				},
			},
			value: json.Number("1.21"),
			expected: []SchemaError{
				newError(EnumError, "a"),
			},
		},
	}

	for _, test := range tests {
		c := NewBaseConstraint(test.schema)
		c.validateEnum(test.value, "a")
		assert.Equal(t, test.expected, c.Errors())
	}
}