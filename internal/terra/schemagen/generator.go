package schemagen

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/ettle/strcase"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/maps"
)

// DocsFunc is a function that returns the documentation for a struct field.
type DocsFunc = func(obj any, sf reflect.StructField) string

// CustomizerFunc is a function used to customize a schema.
// It returns an optional function name, the amended kind and a bool determining
// if the field should be skipped.
type CustomizerFunc = func(obj any, sf *reflect.StructField, typ reflect.Type, s *schema.Schema) (string, reflect.Kind, bool)

var errSkip = errors.New("skip")

// Generator generates Terraform schemas.
type Generator struct {
	docsFn       DocsFunc
	customizerFn CustomizerFunc
	tag          string
}

// New returns a schema generator.
func New(docsFn DocsFunc, customizerFn CustomizerFunc, tag string) *Generator {
	if tag == "" {
		tag = "json"
	}
	if docsFn == nil {
		docsFn = func(any, reflect.StructField) string { return "" }
	}
	if customizerFn == nil {
		customizerFn = func(_ any, _ *reflect.StructField, typ reflect.Type, _ *schema.Schema) (string, reflect.Kind, bool) {
			return "", typ.Kind(), true
		}
	}

	return &Generator{
		docsFn:       docsFn,
		customizerFn: customizerFn,
		tag:          tag,
	}
}

// Struct generates the schema from the given object.
func (g *Generator) Struct(obj any) (map[string]string, error) {
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("obj must be a struct or pointer to a struct")
	}

	fields := make(map[string]string, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)

		name := g.ResolveName(sf)
		code, err := g.genField(name, obj, &sf, sf.Type, false)
		if err != nil {
			switch {
			case errors.Is(err, errSkip):
				continue
			default:
				return nil, err
			}
		}

		fields[strcase.ToSnake(name)] = code
	}
	return fields, nil
}

//nolint:cyclop
func (g *Generator) genField(name string, obj any, sf *reflect.StructField, typ reflect.Type, isNested bool) (string, error) {
	typ, isPtr := derefType(typ)
	s := &schema.Schema{}

	fnName, kind, ok := g.customizerFn(obj, sf, typ, s)
	if !ok {
		return "", errSkip
	}

	if sf != nil {
		s.Description = g.docsFn(obj, *sf)
	}

	var setFunc string
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s.Type = schema.TypeInt
	case reflect.Float32, reflect.Float64:
		s.Type = schema.TypeFloat
	case reflect.String:
		s.Type = schema.TypeString
	case reflect.Bool:
		s.Type = schema.TypeBool
	case reflect.Slice:
		s.Type = schema.TypeList
		elem, err := g.genField("", obj, nil, typ.Elem(), true)
		if err != nil {
			return "", fmt.Errorf("generating elem for %q: %w", name, err)
		}
		s.Elem = elem
	case reflect.Map:
		s.Type = schema.TypeMap
		elem, err := g.genField("", obj, nil, typ.Elem(), true)
		if err != nil {
			return "", fmt.Errorf("generating elem for %q: %w", name, err)
		}
		s.Elem = elem
	case reflect.Struct:
		s.Type = schema.TypeList
		s.MaxItems = 1

		var elem string
		if fnName == "" {
			elem = "&schema.Resource{\nSchema: map[string]*schema.Schema{\n"

			newObj := reflect.New(typ).Elem().Interface()
			fields, err := g.Struct(newObj)
			if err != nil {
				return "", err
			}

			fieldNames := maps.Keys(fields)
			sort.Strings(fieldNames)
			for _, fieldName := range fieldNames {
				elem += fmt.Sprintf("%q: %s,\n", fieldName, fields[fieldName])
			}
			elem += "},\n}"
		} else {
			elem = "&schema.Resource{Schema:" + fnName + "()}"
		}

		if isNested {
			return elem, nil
		}
		s.Elem = elem
	default:
		return "", fmt.Errorf("unsupported type %q for %q", typ.String(), name)
	}

	if isPtr &&
		(s.Type == schema.TypeInt || s.Type == schema.TypeFloat || s.Type == schema.TypeString || s.Type == schema.TypeBool) {
		// Terraform always sets a default value for all properties not set. In the pointer case,
		// we can no longer tell if it is nil or set by the user to a default value. Instead a block should
		// be created in this case.
		ptrS := &schema.Schema{}
		*ptrS = *s

		// Ensure we override all values.
		s.Optional = false
		s.Computed = false
		s.Required = true
		elem, err := genSchemaField(s, "", false)
		if err != nil {
			return "", err
		}

		ptrS.Type = schema.TypeList
		ptrS.MaxItems = 1
		ptrS.Default = nil
		ptrS.Elem = "&schema.Resource{\nSchema: map[string]*schema.Schema{\n\"value\": " + elem + ",\n},\n}"
		s = ptrS
	}

	return genSchemaField(s, setFunc, isNested)
}

// ResolveName resolves the name of a field.
func (g *Generator) ResolveName(sf reflect.StructField) string {
	jsonName := sf.Tag.Get(g.tag)
	if name, _, ok := strings.Cut(jsonName, ","); ok {
		jsonName = name
	}
	if jsonName != "" && jsonName != "-" {
		return jsonName
	}
	return sf.Name
}

func genSchemaField(s *schema.Schema, setFunc string, isNested bool) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	err := schemaTmpl.Execute(buf, struct {
		Schema   *schema.Schema
		SetFunc  string
		IsNested bool
	}{
		Schema:   s,
		SetFunc:  setFunc,
		IsNested: isNested,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func derefType(t reflect.Type) (reflect.Type, bool) {
	if t.Kind() == reflect.Ptr {
		return t.Elem(), true
	}
	return t, false
}

var schemaTmpl = template.Must(template.New("schema").Parse(
	`{{if .IsNested}}&schema.Schema{{end}}{{"{"}}{{if not .IsNested}}
{{end}}Type: schema.{{.Schema.Type}},{{if ne .Schema.Description ""}}
Description: {{printf "%q" .Schema.Description}},{{end}}{{if .Schema.Required}}
Required: {{.Schema.Required}},{{end}}{{if .Schema.Optional}}
Optional: {{.Schema.Optional}},{{end}}{{if .Schema.ForceNew}}
ForceNew: {{.Schema.ForceNew}},{{end}}{{if .Schema.Computed}}
Computed: {{.Schema.Computed}},{{end}}{{if gt .Schema.MaxItems 0}}
MaxItems: {{.Schema.MaxItems}},{{end}}{{if .Schema.Elem}}
Elem: {{.Schema.Elem}},{{end}}{{if ne .SetFunc ""}}{{if not .IsNested}}
{{end}}Set: {{.SetFunc}},{{end}}{{if not .IsNested}}
{{end}}{{"}"}}`,
))
