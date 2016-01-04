package constraint
import (
	"reflect"
	"fmt"
)

type Keyword interface {
	Name() string
	Validate(v interface{}, ctx ValidationContext) error
}

// Type defines keyword "type"
type Type struct {
	value string
	values []string
}

func (t *Type) Name() string {
	return "type"
}

func (t *Type) Validate(v interface{}, ctx ValidationContext) error {
	if reflect.TypeOf(v) == reflect.Bool {
		return nil
	}

	return fmt.Errorf("Type %s is required but got %s at %s",
	 t.value, "f", ctx.Path())
}