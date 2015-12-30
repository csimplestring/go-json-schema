package schema

import (
	"encoding/json"
	"strings"
	"testing"

	"code.yieldr.com/px/util/test/assert"
)

func deserializeSchema(str string) (Schema, error) {
	jsonSchema := make(Schema)
	decoder := json.NewDecoder(strings.NewReader(str))
	decoder.UseNumber()
	err := decoder.Decode(&jsonSchema)

	return jsonSchema, err
}

func TestNumericSchema(t *testing.T) {
	jsonString := `
	{
		"multipleOf": 123,
		"maximum": 100,
		"exclusiveMaximum": true,
		"minimum": 0,
		"exclusiveMinimum": true
	}
	`
	jsonSchema := make(Schema)

	decoder := json.NewDecoder(strings.NewReader(jsonString))
	decoder.UseNumber()

	err := decoder.Decode(&jsonSchema)
	assert.NoError(t, err)

	divided, exist := jsonSchema.MultipleOf()
	assert.Equal(t, float64(123), divided)
	assert.True(t, exist)

	max, exist := jsonSchema.Maximum()
	assert.Equal(t, float64(100), max)
	assert.True(t, exist)

	min, exist := jsonSchema.Minimum()
	assert.Equal(t, float64(0), min)
	assert.True(t, exist)
}

func TestStringSchema(t *testing.T) {
	jsonString := `
	{
		"maxLength": 10,
		"minLength": 1,
		"pattern": "abc"
	}
	`
	jsonSchema := make(Schema)

	decoder := json.NewDecoder(strings.NewReader(jsonString))
	decoder.UseNumber()

	err := decoder.Decode(&jsonSchema)
	assert.NoError(t, err)

	maxLen, exist := jsonSchema.MaxLength()
	assert.Equal(t, 10, maxLen)
	assert.True(t, exist)

	minLen, exist := jsonSchema.MinLength()
	assert.Equal(t, 1, minLen)
	assert.True(t, exist)

	pattern, exist := jsonSchema.Pattern()
	assert.Equal(t, "abc", pattern)
	assert.True(t, exist)
}

func TestObjectSchema(t *testing.T) {
	// common
	jsonString := `
	{
		"maxProperties": 10,
		"minProperties": 2,
		"required": ["a", "b"],
		"properties": {
			"a": {
				"type": "integer"
			},
			"b": {
				"type": "integer"
			}
		},
		"patternProperties": {
			"a[0-9]b": {
				"type": "integer"
			}
		}
	}
	`

	schema, err := deserializeSchema(jsonString)
	assert.NoError(t, err)

	max, exist := schema.MaxProperties()
	assert.Equal(t, true, exist)
	assert.Equal(t, 10, max)

	min, exist := schema.MinProperties()
	assert.Equal(t, true, exist)
	assert.Equal(t, 2, min)

	required, exist := schema.Required()
	assert.Equal(t, true, exist)
	assert.Equal(t, []string{"a", "b"}, required)

	prop, exist := schema.Properties()
	assert.Equal(t, true, exist)
	assert.Equal(t, Properties(map[string]Schema{
		"a": Schema{"type": "integer"},
		"b": Schema{"type": "integer"},
	}), prop)

	// test for pattern properties

	pattern, exist := schema.PatternProperties()
	assert.Equal(t, true, exist)

	match, patternSchema := pattern.Match("a2b")
	assert.Equal(t, true, match)
	assert.Equal(t, Schema{
		"type": "integer",
	}, patternSchema)

	// test for additionalProperties
	tests := []struct {
		input            string
		expectedAddition *AdditionalProperties
	}{
		{
			input: `
			{
				"additionalProperties": false
			}
			`,
			expectedAddition: &AdditionalProperties{
				Schema:    nil,
				BoolValue: false,
				IsBool:    true,
				IsSchema:  false,
			},
		},
		{
			input: `
			{
				"additionalProperties": {
					"type": "integer"
				}
			}
			`,
			expectedAddition: &AdditionalProperties{
				Schema:    Schema{
					"type": "integer",
				},
				BoolValue: false,
				IsBool:    false,
				IsSchema:  true,
			},
		},
	}

	for _, test := range tests {
		s, err := deserializeSchema(test.input)
		assert.NoError(t, err)

		addition, _ := s.AdditionalProperties()
		assert.Equal(t, test.expectedAddition, addition)
	}
}

func TestArraySchema(t *testing.T) {

	// list
	jsonString := `
		{
			"items": {
				"type": "integer"
			}
		}
	`

	expectedItem := Schema{
		"type": "integer",
	}

	schema, err := deserializeSchema(jsonString)
	assert.NoError(t, err)

	actual, schemaArray, exist := schema.Items()
	assert.Nil(t, schemaArray)
	assert.True(t, exist)
	assert.Equal(t, expectedItem, actual)

	// tuple
	jsonString = `
			{
				"items": [
					{
						"type": "integer"
					},
					{
						"type": "number"
					}
				]
			}
			`
	expectedItems := []Schema{
		Schema{
			"type": "integer",
		},
		Schema{
			"type": "number",
		},
	}

	schema, err = deserializeSchema(jsonString)
	assert.NoError(t, err)

	single, schemaArray, exist := schema.Items()
	assert.Nil(t, single)
	assert.True(t, exist)
	assert.Equal(t, expectedItems, schemaArray)

	// additional item
	jsonString = `
			{
				"additionalItems": true
			}
			`

	schema, err = deserializeSchema(jsonString)
	assert.NoError(t, err)
	additionSchema, isAllowAddition, existAddition := schema.AdditionalItems()
	assert.True(t, existAddition)
	assert.Nil(t, additionSchema)
	assert.True(t, isAllowAddition)
}

func TestSchemaType(t *testing.T) {
	tests := []struct {
		jsonString   string
		expectedType interface{}
	}{
		{
			jsonString: `
			{
				"type": "boolean"
			}
			`,
			expectedType: JsonBoolean,
		},
		{
			jsonString: `
			{
				"type": ["integer", "boolean"]
			}
			`,
			expectedType: []JsonType{JsonInteger, JsonBoolean},
		},
	}

	for _, test := range tests {
		jsonSchema := make(Schema)
		decoder := json.NewDecoder(strings.NewReader(test.jsonString))
		decoder.UseNumber()

		err := decoder.Decode(&jsonSchema)
		assert.NoError(t, err)

		actualType, actualTypes, _ := jsonSchema.Type()
		if actualType != "" {
			assert.Equal(t, test.expectedType, actualType)
		} else if actualTypes != nil {
			assert.Equal(t, test.expectedType, actualTypes)
		} else {
			t.Fatalf("unknown json schema type %v", actualType)
		}
	}
}

func TestSchemaEnum(t *testing.T) {
	tests := []struct {
		jsonString   string
		expectedEnum interface{}
	}{
		{
			jsonString: `
			{
				"enum": [1, 2.3, "str", ["a", "b"], {"foo": "bar"}]
			}
			`,
			expectedEnum: []interface{}{
				json.Number("1"),
				json.Number("2.3"),
				"str",
				[]interface{}{"a", "b"},
				map[string]interface{}{
					"foo": "bar",
				},
			},
		},
	}

	for _, test := range tests {
		jsonSchema := make(Schema)
		decoder := json.NewDecoder(strings.NewReader(test.jsonString))
		decoder.UseNumber()

		err := decoder.Decode(&jsonSchema)
		assert.NoError(t, err)

		actual, _ := jsonSchema.Enum()
		assert.Equal(t, test.expectedEnum, actual)
	}
}

func TestSchemaAllOf(t *testing.T) {
	tests := []struct {
		jsonString    string
		expectedAllOf []Schema
	}{
		{
			jsonString: `
			{
				"allOf": [
					{
						"type": "integer"
					},
					{
						"type": "string"
					}
				]
			}
			`,
			expectedAllOf: []Schema{
				Schema{
					"type": "integer",
				},
				Schema{
					"type": "string",
				},
			},
		},
	}

	for _, test := range tests {
		jsonSchema := make(Schema)
		decoder := json.NewDecoder(strings.NewReader(test.jsonString))
		decoder.UseNumber()

		err := decoder.Decode(&jsonSchema)
		assert.NoError(t, err)

		actual, _ := jsonSchema.AllOf()
		assert.Equal(t, test.expectedAllOf, actual)
	}
}
