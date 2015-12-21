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

func TestArrayConstraintItems(t *testing.T) {
	path := "a"

	listTests := []struct {
		itemSchema     Schema
		value          []interface{}
		expectedErrors []SchemaError
	}{
		{
			itemSchema: Schema{
				"items": map[string]interface{} {
					"type": "integer",
				},
			},
			value: []interface{} {json.Number("1")},
			expectedErrors: nil,
		},
		{
			itemSchema: Schema{
				"items": map[string]interface{} {
					"type": "integer",
				},
			},
			value: []interface{} {json.Number("1.1")},
			expectedErrors: []SchemaError{
				newError(TypeNotMatchError, "a[0]"),
			},
		},
	}

	for _, test := range listTests {
		c :=  NewArrayConstraint(test.itemSchema)
		c.validateItems(test.value, path)

		assert.Equal(t, test.expectedErrors, c.Errors())
	}

	tupleTests := []struct{
		itemSchema     Schema
		value          []interface{}
		expectedErrors []SchemaError
	}{
		{
			itemSchema: Schema{
				"items": []interface{} {
					map[string]interface{} {
						"type": "integer",
					},
				},
			},
			value: []interface{} {json.Number("1")},
			expectedErrors: nil,
		},
		{
			itemSchema: Schema{
				"items": []interface{} {
					map[string]interface{} {
						"type": "integer",
					},
				},
			},
			value: []interface{} {json.Number("1.1")},
			expectedErrors: []SchemaError{
				newError(TypeNotMatchError, "a[0]"),
			},
		},
		{
			itemSchema: Schema{
				"items": []interface{} {
					map[string]interface{} {
						"type": "integer",
					},
					map[string]interface{} {
						"type": "string",
					},
				},
			},
			value: []interface{} {
				json.Number("1"),
				json.Number("1"),
			},
			expectedErrors: []SchemaError{
				newError(TypeNotMatchError, "a[1]"),
			},
		},
		{
			itemSchema: Schema{
				"items": []interface{} {
					map[string]interface{} {
						"type": "integer",
					},
				},
				"additionalItems": true,
			},
			value: []interface{} {
				json.Number("1"),
				"str",
			},
			expectedErrors: nil,
		},
		{
			itemSchema: Schema{
				"items": []interface{} {
					map[string]interface{} {
						"type": "integer",
					},
				},
				"additionalItems": false,
			},
			value: []interface{} {
				json.Number("1"),
				"str",
			},
			expectedErrors: []SchemaError{
				newError(ArrayAdditionalItemError, "a[1]"),
			},
		},
		{
			itemSchema: Schema{
				"items": []interface{} {
					map[string]interface{} {
						"type": "integer",
					},
				},
				"additionalItems": map[string]interface{} {
					"type": "string",
				},
			},
			value: []interface{} {
				json.Number("1"),
				"str",
				json.Number("2"),
			},
			expectedErrors: []SchemaError{
				newError(TypeNotMatchError, "a[2]"),
			},
		},
	}

	for _, test := range tupleTests {
		c :=  NewArrayConstraint(test.itemSchema)
		c.validateItems(test.value, path)

		assert.Equal(t, test.expectedErrors, c.Errors())
	}
}