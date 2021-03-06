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
	JsonNumber  = JsonType("number")
	JsonString  = JsonType("string")
	JsonArray   = JsonType("array")
	JsonObject  = JsonType("object")
	JsonBoolean = JsonType("boolean")
	JsonNull    = JsonType("null")
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

type SchemaInterface interface {
	Validate(v interface{}, path string, interactive bool) error
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

func (s Schema) Type() (jsonType JsonType, jsonTypes []JsonType, exist bool) {
	exist = false
	v, ok := s["type"]

	if !ok {
		return
	}

	exist = true
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
	v, exist := s["enum"]
	if !exist {
		return
	}

	exist = true
	enums = v.([]interface{})
	return
}

func (s Schema) AllOf() (all []Schema, exist bool) {
	v, exist := s["allOf"]
	if !exist {
		return
	}

	// v must be an array of valid schema
	exist = true
	for _, one := range v.([]interface{}) {
		all = append(all, Schema(one.(map[string]interface{})))
	}
	return
}

func (s Schema) AnyOf() (any []Schema, exist bool) {
	v, exist := s["anyOf"]
	if !exist {
		return
	}

	// v must be an array of valid schema
	exist = true
	for _, one := range v.([]interface{}) {
		any = append(any, Schema(one.(map[string]interface{})))
	}
	return
}

func (s Schema) OneOf() (all []Schema, exist bool) {
	v, exist := s["oneOf"]
	if !exist {
		return
	}

	// v must be an array of valid schema
	exist = true
	for _, one := range v.([]interface{}) {
		all = append(all, Schema(one.(map[string]interface{})))
	}
	return
}

func (s Schema) Not() (not Schema, exist bool) {
	v, exist := s["not"]
	if !exist {
		return
	}

	// v must be an object
	exist = true
	not = Schema(v.(map[string]interface{}))
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
	v, exist := s["additionalItems"]
	if !exist {
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
	v, exist := s["items"]
	if !exist {
		return
	}

	// item can be an object or an array of objects
	switch v.(type) {
	case map[string]interface{}:
		exist = true
		schema = Schema(v.(map[string]interface{}))
		return
	case []interface{}:
		exist = true
		for _, item := range v.([]interface{}) {
			schemaArray = append(schemaArray, Schema(item.(map[string]interface{})))
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

// validation keywords for object

func (s Schema) MaxProperties() (maxProperties int, exist bool) {
	return s.getIntValue("maxProperties")
}

func (s Schema) MinProperties() (minProperties int, exist bool) {
	return s.getIntValue("minProperties")
}

func (s Schema) Required() (required []string, exist bool) {
	v, exist := s["required"]
	if !exist {
		return
	}

	for _, r := range v.([]interface{}) {
		required = append(required, r.(string))
	}

	return
}

// properties returns a map where the key is the property name and the value
// is a valid json schema for this property.
func (s Schema) Properties() (propertiesSchema map[string]Schema, exist bool) {
	v, exist := s["properties"]
	if !exist {
		return
	}

	propertiesSchema = make(map[string]Schema)
	for key, val := range v.(map[string]interface{}) {
		propertiesSchema[key] = Schema(val.(map[string]interface{}))
	}

	return
}

// The value of "additionalProperties" MUST be a boolean or an object.
// If it is an object, it MUST also be a valid JSON Schema.
func (s Schema) AdditionalProperties() (additionSchema Schema, boolValue bool, exist bool) {
	v, exist := s["additionalProperties"]
	if !exist {
		return
	}

	switch v.(type) {
	case bool:
		boolValue = v.(bool)
		return
	case map[string]interface{}:
		additionSchema = Schema(v.(map[string]interface{}))
		return
	default:
		return
	}
}

func (s Schema) PatternProperties() (patternSchema PatternProperties, exist bool) {
	v, exist := s["patternProperties"]
	if !exist {
		return
	}

	patternSchema = make(PatternProperties)
	for key, val := range v.(map[string]interface{}) {
		patternSchema[key] = Schema(val.(map[string]interface{}))
	}

	return
}

type PatternProperties map[string]Schema

func (p PatternProperties) match(prop string, val interface{}) bool {

}