package terra

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Flatten converts an object into a Terraform data structure using its schema.
func (c *Converter) Flatten(obj any, s map[string]*schema.Schema) (any, error) {
	v := reflect.ValueOf(obj)
	return c.flattenStruct(v, s)
}

func (c *Converter) flatten(v reflect.Value, s *schema.Schema) (any, error) {
	switch {
	case s.Type == schema.TypeInvalid:
		return nil, errors.New("invalid schema type")
	case s.Type == schema.TypeList && s.MaxItems == 1 && isResource(s.Elem):
		res := s.Elem.(*schema.Resource)

		if v.Type().Kind() == reflect.Ptr && v.Type().Elem().Kind() != reflect.Struct {
			val, err := c.flattenPrimitive(v, res.Schema["value"])
			if err != nil {
				return nil, err
			}
			return []any{map[string]any{"value": val}}, nil
		}

		return c.flattenStruct(v, res.Schema)
	case s.Type == schema.TypeList:
		return c.flattenSlice(v, s)
	case s.Type == schema.TypeMap:
		return c.flattenMap(v, s)
	default:
		return c.flattenPrimitive(v, s)
	}
}

func (c *Converter) flattenSlice(v reflect.Value, s *schema.Schema) (any, error) {
	t := v.Type()
	if t.Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected slice, got %s", t.String())
	}

	if v.IsNil() || v.Len() == 0 {
		return []any{}, nil
	}

	var elemS *schema.Schema
	switch {
	case isSchema(s.Elem):
		elemS = s.Elem.(*schema.Schema)
	case isResource(s.Elem):
		elemS = &schema.Schema{}
		*elemS = *s
		elemS.MaxItems = 1
	default:
		// Do nothing.
	}

	d := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		a, err := c.flatten(v.Index(i), elemS)
		if err != nil {
			return nil, err
		}
		if isResource(s.Elem) {
			// This is for the edge case where a slice of structs should
			// not return a []any[]any. Unwind the second slice.
			a = a.([]any)[0]
		}
		d[i] = a
	}
	return d, nil
}

func (c *Converter) flattenMap(v reflect.Value, s *schema.Schema) (any, error) {
	t := v.Type()
	if t.Kind() != reflect.Map {
		return nil, fmt.Errorf("expected map, got %s", t.String())
	}

	if v.Len() == 0 {
		return map[string]any{}, nil
	}

	elemS := s.Elem.(*schema.Schema)

	d := make(map[string]any, v.Len())
	for _, k := range v.MapKeys() {
		a, err := c.flatten(v.MapIndex(k), elemS)
		if err != nil {
			return nil, err
		}

		d[k.String()] = a
	}
	return d, nil
}

func (c *Converter) flattenStruct(v reflect.Value, s map[string]*schema.Schema) (any, error) {
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", t.String())
	}

	d := make(map[string]any, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		name := c.resolveName(sf)

		sch, found := s[name]
		if !found {
			continue
		}
		fv := v.Field(i)
		if fv.IsZero() {
			continue
		}

		val, err := c.flatten(fv, sch)
		if err != nil {
			return nil, err
		}
		d[name] = val
	}
	return []any{d}, nil
}

func (c *Converter) flattenPrimitive(v reflect.Value, s *schema.Schema) (any, error) {
	if v.Type().Kind() == reflect.Ptr {
		if v.IsNil() {
			//nolint:nilnil // This is expected.
			return nil, nil
		}
		v = v.Elem()
	}

	t := v.Type()

	if con, ok := c.conversions[t]; ok {
		obj, err := con.flatten(v.Interface())
		if err != nil {
			return nil, err
		}
		v = reflect.ValueOf(obj)
		t = v.Type()
	}

	switch s.Type {
	case schema.TypeString:
		if t.Kind() == reflect.String {
			return v.String(), nil
		}
	case schema.TypeInt:
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int(v.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int(v.Uint()), nil
		}
	case schema.TypeFloat:
		switch t.Kind() {
		case reflect.Float32, reflect.Float64:
			return v.Float(), nil
		}
	case schema.TypeBool:
		if t.Kind() == reflect.Bool {
			return v.Bool(), nil
		}
	}
	return nil, fmt.Errorf("unsupported primitve type %s for schema type %s", t.String(), s.Type.String())
}

func isResource(v any) bool {
	_, ok := v.(*schema.Resource)
	return ok
}

func isSchema(v any) bool {
	_, ok := v.(*schema.Schema)
	return ok
}
