package schema
import (
	"testing"
	"encoding/json"
	"strings"
	"code.yieldr.com/px/util/test/assert"
)

func TestNumericSchema(t*testing.T)  {
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