package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeConstraint(t *testing.T) {
	tests := []struct {
		schema   Schema
		value    interface{}
		expected []ValidationError
	}{
		{
			schema:   Schema{},
			value:    1,
			expected: nil,
		},
		{
			schema: Schema{
				"type": "integer",
			},
			value:    1,
			expected: nil,
		},
		{
			schema: Schema{
				"type": []interface{}{"integer", "string"},
			},
			value:    "str",
			expected: nil,
		},
		{
			schema: Schema{
				"type": []interface{}{"integer"},
			},
			value: "str",
			expected: []ValidationError{
				newError(TypesNotMatchError, "a"),
			},
		},
		{
			schema: Schema{
				"type": "object",
			},
			value: "str",
			expected: []ValidationError{
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
	tests := []struct {
		schema   Schema
		value    interface{}
		expected []ValidationError
	}{
		{
			schema: Schema{
				"enum": []interface{}{json.Number("1")},
			},
			value:    json.Number("1"),
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
			value:    json.Number("1.2"),
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
			expected: []ValidationError{
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

func TestAllOfConstraint(t *testing.T) {
	tests := []struct {
		schema   Schema
		value    interface{}
		expected []ValidationError
	}{
		{
			schema: Schema{
				"allOf": []interface{}{
					map[string]interface{}{
						"type": "integer",
					},
					map[string]interface{}{
						"enum": []interface{}{
							json.Number("1"),
							json.Number("2"),
						},
					},
				},
			},
			value:    json.Number("1"),
			expected: nil,
		},
		{
			schema: Schema{
				"allOf": []interface{}{
					map[string]interface{}{
						"type": "integer",
					},
					map[string]interface{}{
						"enum": []interface{}{
							json.Number("1"),
							json.Number("2"),
						},
					},
				},
			},
			value: json.Number("3"),
			expected: []ValidationError{
				newError(AllOfError, "a"),
			},
		},
		{
			schema: Schema{
				"allOf": []interface{}{
					map[string]interface{}{
						"type": "integer",
					},
					map[string]interface{}{
						"enum": []interface{}{
							json.Number("1"),
							json.Number("2"),
						},
					},
				},
			},
			value: json.Number("1.3"),
			expected: []ValidationError{
				newError(AllOfError, "a"),
				newError(AllOfError, "a"),
			},
		},
	}

	for _, test := range tests {
		c := NewBaseConstraint(test.schema)
		c.validateAllOf(test.value, "a")
		assert.Equal(t, test.expected, c.Errors())
	}
}

func TestAnyOfConstraint(t *testing.T) {
	tests := []struct {
		schema   Schema
		value    interface{}
		expected []ValidationError
	}{
		{
			schema: Schema{
				"anyOf": []interface{}{
					map[string]interface{}{
						"type": "integer",
					},
					map[string]interface{}{
						"enum": []interface{}{
							json.Number("1.2"),
							json.Number("2.1"),
						},
					},
				},
			},
			value:    json.Number("1"),
			expected: nil,
		},
		{
			schema: Schema{
				"anyOf": []interface{}{
					map[string]interface{}{
						"type": "integer",
					},
					map[string]interface{}{
						"enum": []interface{}{
							json.Number("1.2"),
							json.Number("2.1"),
						},
					},
				},
			},
			value: json.Number("3.1"),
			expected: []ValidationError{
				newError(AnyOfError, "a"),
			},
		},
	}

	for _, test := range tests {
		c := NewBaseConstraint(test.schema)
		c.validateAnyOf(test.value, "a")
		assert.Equal(t, test.expected, c.Errors())
	}
}

func TestOneOfConstraint(t *testing.T) {
	tests := []struct {
		schema   Schema
		value    interface{}
		expected []ValidationError
	}{
		{
			schema: Schema{
				"oneOf": []interface{}{
					map[string]interface{}{
						"type": "integer",
					},
					map[string]interface{}{
						"enum": []interface{}{
							json.Number("1.2"),
							json.Number("2.1"),
						},
					},
				},
			},
			value:    json.Number("1"),
			expected: nil,
		},
		{
			schema: Schema{
				"oneOf": []interface{}{
					map[string]interface{}{
						"type": "number",
					},
					map[string]interface{}{
						"enum": []interface{}{
							json.Number("1.2"),
							json.Number("2.1"),
						},
					},
				},
			},
			value: json.Number("1.2"),
			expected: []ValidationError{
				newError(OneOfError, "a"),
			},
		},
	}

	for _, test := range tests {
		c := NewBaseConstraint(test.schema)
		c.validateOneOf(test.value, "a")
		assert.Equal(t, test.expected, c.Errors())
	}
}

func TestNotConstraint(t *testing.T) {
	tests := []struct {
		schema   Schema
		value    interface{}
		expected []ValidationError
	}{
		{
			schema: Schema{
				"not": map[string]interface{}{
					"type": "integer",
				},
			},
			value:    json.Number("1.1"),
			expected: nil,
		},
		{
			schema: Schema{
				"not": map[string]interface{}{
					"type": "integer",
				},
			},
			value: json.Number("1"),
			expected: []ValidationError{
				newError(NotError, "a"),
			},
		},
	}

	for _, test := range tests {
		c := NewBaseConstraint(test.schema)
		c.validateNot(test.value, "a")
		assert.Equal(t, test.expected, c.Errors())
	}
}

func TestBaseConstraint(t *testing.T) {
	tests := []struct {
		schema   Schema
		value    interface{}
		expected []ValidationError
	}{
		{
			schema: Schema{
				"type": "integer",
			},
			value:    json.Number("1"),
			expected: nil,
		},
	}

	for _, test := range tests {
		c := NewBaseConstraint(test.schema)
		c.Validate(test.value, "a")
		assert.Equal(t, test.expected, c.Errors())
	}
}