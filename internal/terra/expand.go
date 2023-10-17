package terra

import (
	"fmt"
	"reflect"
)

// Expand converts Terraform-formatted data into the given object.
func (c *Converter) Expand(v, obj any) error {
	typ := reflect.TypeOf(obj)
	if typ.Kind() != reflect.Ptr {
		return fmt.Errorf("obj is not a pointer by %q", typ.String())
	}

	return c.expand(v, reflect.ValueOf(obj).Elem())
}

func (c *Converter) expand(a any, objVal reflect.Value) error {
	switch val := a.(type) {
	case []any:
		if len(val) == 0 {
			return nil
		}

		if objVal.Type().Kind() == reflect.Slice {
			return c.expandSlice(val, objVal)
		}

		m, ok := val[0].(map[string]any)
		if !ok {
			return fmt.Errorf("struct type requires data to be a map[string]any")
		}

		if objVal.Type().Kind() == reflect.Ptr && objVal.Type().Elem().Kind() != reflect.Struct && len(m) == 1 {
			// This is a pointer value struct, unwrap it.
			if v, ok := m["value"]; ok {
				return c.expandPrimitive(v, objVal)
			}
		}

		return c.expandStruct(m, objVal)
	case map[string]any:
		if len(val) == 0 {
			return nil
		}

		if objVal.Type().Kind() == reflect.Map {
			return c.expandMap(val, objVal)
		}

		// It could be we have a slice of struct.
		return c.expandStruct(val, objVal)
	default:
		return c.expandPrimitive(a, objVal)
	}
}

func (c *Converter) expandStruct(m map[string]any, objVal reflect.Value) error {
	if objVal.Type().Kind() == reflect.Ptr {
		if objVal.IsNil() {
			objVal.Set(reflect.New(objVal.Type().Elem()))
		}
		objVal = objVal.Elem()
	}

	t := objVal.Type()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %s", t.String())
	}

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		name := c.resolveName(sf)

		val, found := m[name]
		if !found {
			continue
		}

		if err := c.expand(val, objVal.Field(i)); err != nil {
			return err
		}
	}
	return nil
}

func (c *Converter) expandSlice(a []any, objVal reflect.Value) error {
	t := objVal.Type()
	if t.Kind() != reflect.Slice {
		return fmt.Errorf("expected slice, got %s", t.String())
	}

	l := len(a)
	if objVal.Len() < l {
		objVal.Set(reflect.MakeSlice(t, l, l))
	}

	for i, v := range a {
		if err := c.expand(v, objVal.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func (c *Converter) expandMap(m map[string]any, objVal reflect.Value) error {
	t := objVal.Type()
	if t.Kind() != reflect.Map {
		return fmt.Errorf("expected map, got %s", t.String())
	}

	if objVal.IsNil() {
		objVal.Set(reflect.MakeMap(t))
	}

	for k, v := range m {
		val := reflect.New(t.Elem()).Elem()
		if err := c.expand(v, val); err != nil {
			return err
		}

		objVal.SetMapIndex(reflect.ValueOf(k), val)
	}
	return nil
}

func (c *Converter) expandPrimitive(v any, objVal reflect.Value) error {
	if objVal.Type().Kind() == reflect.Ptr {
		if objVal.IsNil() {
			objVal.Set(reflect.New(objVal.Type().Elem()))
		}
		objVal = objVal.Elem()
	}

	objTyp := objVal.Type()

	if con, ok := c.conversions[objTyp]; ok {
		var err error
		v, err = con.expand(v)
		if err != nil {
			return err
		}
	}

	vVal := reflect.ValueOf(v)
	vTyp := vVal.Type()

	switch {
	case vTyp.AssignableTo(objTyp):
		objVal.Set(reflect.ValueOf(v))
		return nil
	case vTyp.ConvertibleTo(objTyp):
		objVal.Set(vVal.Convert(objTyp))
		return nil
	default:
		return fmt.Errorf("primitive of type %s not supported", objTyp.String())
	}
}
