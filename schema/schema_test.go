package schema
import (
	"testing"
	"encoding/json"
	"strings"
	"code.yieldr.com/px/util/test/assert"
)

func TestNumericSchema(t*testing.T) {
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

func TestStringSchema(t*testing.T) {
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

func TestArraySchema(t *testing.T) {
	tests := []struct {
		jsonString              string

		expectedItems           interface{}
		expectedAdditionalItems interface{}
		expectedMax             int
		expectedMin             int
		expectedUnique          bool
	}{
		{
			jsonString : `
			{
				"maxItems": 5,
				"minItems": 3,
				"uniqueItems": true,
				"items": {
					"type": "integer"
				},
				"additionalItems": false
			}
			`,
			expectedItems: Schema{
				"type": "integer",
			},
			expectedAdditionalItems: false,
			expectedMax: 5,
			expectedMin: 3,
			expectedUnique: true,
		},
		{
			jsonString : `
			{
				"items": [
					{
						"type": "integer"
					},
					{
						"type": "number"
					}
				],
				"additionalItems": {
					"type": "integer"
				}
			}
			`,
			expectedItems: []Schema{
				Schema{
					"type": "integer",
				},
				Schema{
					"type": "number",
				},
			},
			expectedAdditionalItems: Schema{
				"type": "integer",
			},
			expectedUnique: false,
		},
	}

	for _, test := range tests {
		jsonSchema := make(Schema)
		decoder := json.NewDecoder(strings.NewReader(test.jsonString))
		decoder.UseNumber()

		err := decoder.Decode(&jsonSchema)
		assert.NoError(t, err)

		schema, schemaArray, exist := jsonSchema.Items()
		if exist {
			if schema != nil {
				assert.Equal(t, test.expectedItems, schema)
			} else {
				assert.Equal(t, test.expectedItems, schemaArray)
			}
		}

		schema, boolValue, exist := jsonSchema.AdditionalItems()
		if exist {
			if schema == nil {
				assert.Equal(t, test.expectedAdditionalItems, boolValue)
			} else {
				assert.Equal(t, test.expectedAdditionalItems, schema)
			}
		}

		if max, exist := jsonSchema.MaxItems(); exist {
			assert.Equal(t, test.expectedMax, max)
		}

		if min, exist := jsonSchema.MinItems(); exist {
			assert.Equal(t, test.expectedMin, min)
		}

		assert.Equal(t, test.expectedUnique, jsonSchema.UniqueItems())
	}
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

		actualType, actualTypes := jsonSchema.Type()
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