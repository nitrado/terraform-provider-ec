package schemagen_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nitrado/terraform-provider-ec/internal/terra/schemagen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_Struct(t *testing.T) {
	gen := schemagen.New(nil, nil, "")

	got, err := gen.Struct(&TestObject{})

	require.NoError(t, err)
	want := map[string]string{
		"float":    "{\nType: schema.TypeFloat,\n}",
		"int":      "{\nType: schema.TypeInt,\n}",
		"str":      "{\nType: schema.TypeString,\n}",
		"bool":     "{\nType: schema.TypeBool,\n}",
		"ptr_bool": "{\nType: schema.TypeList,\nMaxItems: 1,\nElem: &schema.Resource{\nSchema: map[string]*schema.Schema{\n\"value\": {\nType: schema.TypeBool,\nRequired: true,\n},\n},\n},\n}",
		"slice":    "{\nType: schema.TypeList,\nElem: &schema.Resource{\nSchema: map[string]*schema.Schema{\n\"a\": {\nType: schema.TypeString,\n},\n},\n},\n}",
		"map":      "{\nType: schema.TypeMap,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
		"struct":   "{\nType: schema.TypeList,\nMaxItems: 1,\nElem: &schema.Resource{\nSchema: map[string]*schema.Schema{\n\"a\": {\nType: schema.TypeString,\n},\n},\n},\n}",
	}
	assert.Equal(t, want, got)
}

func TestGenerator_StructUsesCustomFunction(t *testing.T) {
	c := func(_ any, _ *reflect.StructField, typ reflect.Type, _ *schema.Schema) (string, reflect.Kind, bool) {
		if typ == reflect.TypeOf(T{}) {
			return "myFunc", typ.Kind(), true
		}
		return "", typ.Kind(), true
	}

	gen := schemagen.New(nil, c, "")

	got, err := gen.Struct(S{})

	require.NoError(t, err)
	want := map[string]string{
		"struct": "{\nType: schema.TypeList,\nMaxItems: 1,\nElem: &schema.Resource{Schema:myFunc()},\n}",
	}
	assert.Equal(t, want, got)
}

func TestGenerator_StructUsesCustomKind(t *testing.T) {
	c := func(_ any, _ *reflect.StructField, typ reflect.Type, _ *schema.Schema) (string, reflect.Kind, bool) {
		if typ == reflect.TypeOf(T{}) {
			return "", reflect.String, true
		}
		return "", typ.Kind(), true
	}

	gen := schemagen.New(nil, c, "")

	got, err := gen.Struct(S{})

	require.NoError(t, err)
	want := map[string]string{
		"struct": "{\nType: schema.TypeString,\n}",
	}
	assert.Equal(t, want, got)
}

func TestGenerator_StructSkipsFields(t *testing.T) {
	c := func(_ any, _ *reflect.StructField, typ reflect.Type, _ *schema.Schema) (string, reflect.Kind, bool) {
		if typ == reflect.TypeOf(T{}) {
			return "", reflect.String, false
		}
		return "", typ.Kind(), true
	}

	gen := schemagen.New(nil, c, "")

	got, err := gen.Struct(S{})

	require.NoError(t, err)
	want := map[string]string{}
	assert.Equal(t, want, got)
}

type TestObject struct {
	Str     string         `json:"str"`
	Int     int            `json:"int"`
	Float   float64        `json:"float,omitempty"`
	Bool    bool           `json:"bool"`
	PtrBool *bool          `json:"ptrBool"`
	Slice   []T            `json:"slice"`
	Map     map[string]int `json:"map"`
	Struct  *T
}

type S struct {
	Struct T `json:"struct"`
}

type T struct {
	A string `json:"a"`
}
