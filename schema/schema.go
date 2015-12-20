package schema
import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

//
//type Schema struct {
//	ID string
//	Schema string
//	Title string
//	Description string
//}

// firstly, assume the schema passed in is valid, after the schema validation is finished, you can validate the schema
// against it.

type JsonType string

const (
	JsonInteger = JsonType("integer")
	JsonNumber = JsonType("number")
	JsonString = JsonType("string")
	JsonArray = JsonType("array")
	JsonObject = JsonType("object")
	JsonBoolean = JsonType("boolean")
	JsonNull = JsonType("null")
)

func getJsonType(v interface{}) (JsonType, error) {
	switch v.(type) {
	case bool:
		return JsonBoolean, nil
	case int8, int16, int32, int, int64:
		return JsonInteger, nil
	case float32, float64:
		return JsonNumber, nil
	case string:
		return JsonString, nil
	case []interface{}:
		return JsonArray, nil
	case map[string]interface{}:
		return JsonObject, nil
	case nil:
		return JsonNull, nil
	case json.Number:
		if strings.Contains(v.(json.Number).String(), ".") {
			return JsonNumber, nil
		} else {
			return JsonInteger, nil
		}
	default:
		return JsonType(""), fmt.Errorf("Unsupported json type %s", reflect.TypeOf(v).String())
	}
}

type Schema map[string]interface{}

func (s Schema) getFloat64Value(key string) (value float64, exist bool) {
	value = 0
	exist = false

	if v, ok := s[key]; ok {
		value, _ = v.(json.Number).Float64()
		exist = true
	}

	return
}

func (s Schema) getBoolValue(key string) (value bool, exist bool) {
	value = false
	exist = false

	if v, ok := s[key]; ok {
		value, _ = v.(bool)
		exist = true
	}

	return
}

func (s Schema) getIntValue(key string) (value int, exist bool) {
	value = 0
	exist = false

	if v, ok := s[key]; ok {
		i64, _ := v.(json.Number).Int64()
		value = int(i64)
		exist = true
	}

	return
}

func (s Schema) getStringValue(key string) (value string, exist bool) {
	value = ""
	exist = false

	if v, ok := s[key]; ok {
		value, _ = v.(string)
		exist = true
	}

	return
}

// validation keywords for any instance

func (s Schema) Type() (jsonType JsonType, jsonTypes []JsonType) {
	v, ok := s["type"]

	if !ok {
		return
	}

	switch v.(type) {
	case string:
		jsonType = JsonType(v.(string))
		return
	case []interface{}:
		for _, t := range v.([]interface{}) {
			jsonTypes = append(jsonTypes, JsonType(t.(string)))
		}
		return
	}

	return
}

func (s Schema) Enum() (enums []interface{}, exist bool) {
	exist = false

	v, ok := s["enum"]
	if !ok {
		return
	}

	exist = true
	enums = v.([]interface{})
	return
}

// validation keywords for numeric

func (s Schema) MultipleOf() (divided float64, exist bool) {
	return s.getFloat64Value("multipleOf")
}

func (s Schema) Maximum() (max float64, exist bool) {
	return s.getFloat64Value("maximum")
}

func (s Schema) Minimum() (max float64, exist bool) {
	return s.getFloat64Value("minimum")
}

func (s Schema) ExclusiveMinimum() bool {
	if v, ok := s.getBoolValue("exclusiveMinimum"); ok {
		return v
	}
	return false
}

func (s Schema) ExclusiveMaximum() bool {
	if v, ok := s.getBoolValue("exclusiveMaximum"); ok {
		return v
	}
	return false
}

// validation keywords for string

func (s Schema) MaxLength() (maxLen int, exist bool) {
	return s.getIntValue("maxLength")
}

func (s Schema) MinLength() (minLen int, exist bool) {
	return s.getIntValue("minLength")
}

func (s Schema) Pattern() (pattern string, exist bool) {
	return s.getStringValue("pattern")
}

// validation keywords for array

func (s Schema) AdditionalItems() (schema Schema, boolValue bool, exist bool) {
	exist = false

	v, ok := s["additionalItems"]
	if !ok {
		return
	}

	switch v.(type) {
	case bool:
		exist = true
		boolValue = v.(bool)
		return
	case map[string]interface{}:
		exist = true
		schema = Schema(v.(map[string]interface{}))
		return
	default:
		return
	}
}

func (s Schema) Items() (schema Schema, schemaArray []Schema, exist bool) {
	schema, schemaArray, exist = nil, nil, false

	v, ok := s["items"]
	if !ok {
		return
	}

	switch v.(type) {
	case map[string]interface{}:
		exist = true
		schema = Schema(v.(map[string]interface{}))
		return
	case []map[string]interface{}:
		exist = true
		for _, itemSchema := range v.([]map[string]interface{}) {
			schemaArray = append(schemaArray, itemSchema)
		}
		return
	default:
		return
	}
}

func (s Schema) MaxItems() (maxItems int, exist bool) {
	return s.getIntValue("maxItems")
}

func (s Schema) MinItems() (minItems int, exist bool) {
	return s.getIntValue("minItems")
}

func (s Schema) UniqueItems() bool {
	if v, ok := s.getBoolValue("uniqueItems"); ok {
		return v
	}
	return false
}