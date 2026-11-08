package scalar

import (
	"fmt"
)

type GraphQLType interface {
	ImplementsGraphQLType(string) bool
	UnmarshalGraphQL(any) error
}

type Null[T any] struct {
	Value *T
	Set   bool
}

func (t Null[T]) ImplementsGraphQLType(name string) bool {
	switch typ := any(t.Value).(type) {
	case *string:
		return name == "String"
	case *bool:
		return name == "Boolean"
	case *int32:
		return name == "Int"
	case *float64:
		return name == "Float"
	case GraphQLType:
		return typ.ImplementsGraphQLType(name)
	default:
		return false
	}
}

func (t *Null[T]) UnmarshalGraphQL(input any) error {
	t.Set = true

	if input == nil {
		return nil
	}

	if v, ok := input.(T); ok {
		t.Value = &v
		return nil
	}

	switch v := any(t.Value).(type) {
	case *float64:
		switch i := input.(type) {
		case int:
			coerced := float64(i)
			t.Value = any(&coerced).(*T)
			return nil
		case int32:
			coerced := float64(i)
			t.Value = any(&coerced).(*T)
			return nil
		case int64:
			coerced := float64(i)
			t.Value = any(&coerced).(*T)
			return nil
		}
	case GraphQLType:
		return v.UnmarshalGraphQL(input)
	}

	var v T
	return fmt.Errorf("wrong type for %T: %T", v, input)
}

func (*Null[T]) Nullable() {}

func (t Null[T]) Data() (any, bool) {
	return t.Value, t.Set
}
