package schema
import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestObjectMaxMinRequired(t *testing.T) {
	tests := []struct {
		obj      map[string]interface{}
		schema   Schema
		expected []SchemaError
	}{
		{
			obj: map[string]interface{}{
				"a": json.Number("5"),
				"b": json.Number("6"),
			},
			schema: Schema{
				"maxProperties": json.Number("10"),
				"minProperties": json.Number("2"),
				"required": []interface{}{
					"a", "b",
				},
//				"properties": map[string]interface{}{
//					"a": map[string]interface{}{
//						"type": "integer",
//					},
//					"b": map[string]interface{}{
//						"type": "integer",
//					},
//				},
			},
			expected: nil,
		},
		{
			obj: map[string]interface{}{
				"a": json.Number("5"),
			},
			schema: Schema{
				"maxProperties": json.Number("10"),
				"minProperties": json.Number("2"),
				"required": []interface{}{
					"a", "b",
				},
//				"properties": map[string]interface{}{
//					"a": map[string]interface{}{
//						"type": "integer",
//					},
//					"b": map[string]interface{}{
//						"type": "integer",
//					},
//				},
			},
			expected: []SchemaError{
				newError(ObjectMinPropertiesError, "p"),
				newError(ObjectRequiredPropertiesError, "p"),
			},
		},
	}

	path := "p"
	for _, test :=range tests {
		c := NewObjectConstraint(test.schema)
		c.Validate(test.obj, path)

		assert.Equal(t, c.Errors(), test.expected)
	}
}


