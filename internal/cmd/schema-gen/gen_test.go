package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var update = flag.Bool("update", false, "update the golden files of this test")

func TestGenerator_Generate(t *testing.T) {
	gen := NewGenerator()

	got, err := gen.Generate(&TestObject{}, "testdata", "testObjectSchema")
	require.NoError(t, err)
	if *update {
		_ = os.WriteFile("testdata/schema_testobject.go", got, 0o644)
	}

	want, err := os.ReadFile("testdata/schema_testobject.go")
	require.NoError(t, err)
	assert.Equal(t, string(want), string(got))
}

type TestObject struct {
	metav1.TypeMeta   `json:""`
	metav1.ObjectMeta `json:"metadata"`

	Str    string            `json:"str"`
	Int    int               `json:"int"`
	Float  float64           `json:"float,omitempty"`
	Bool   bool              `json:"bool"`
	Slice  []T               `json:"slice"`
	Map    map[string]int    `json:"map"`
	Struct *T                `json:"struct"`
	CPU    resource.Quantity `json:"cpu"`
	Status string            `json:"status"`
}

func (o TestObject) Docs() map[string]string {
	return map[string]string{
		"str": "String docs",
	}
}

func (o TestObject) Attributes() map[string]string {
	return map[string]string{
		"int":  "readonly",
		"bool": "required",
	}
}

type T struct {
	A string `json:"a"`
}
