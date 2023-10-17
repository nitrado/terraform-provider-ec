package terra_test

import (
	"testing"

	"github.com/nitrado/terraform-provider-ec/internal/terra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestConverter_Flatten(t *testing.T) {
	c := terra.New()
	c.Register(resource.Quantity{}, func(v any) (any, error) {
		return resource.ParseQuantity(v.(string))
	}, func(v any) (any, error) {
		q := v.(resource.Quantity)
		return (&q).String(), nil
	})

	obj := TestObject{
		Str:   "test-str",
		Alias: StrAlias("test-alias"),
		Int:   1,
		Float: 2.3,
		Bool:  true,
		Slice: []T{
			{A: "test-t"},
			{A: "test-t-also"},
		},
		Map: map[string]int{
			"foo": 4,
		},
		Struct: &T{
			A: "test-ptr-t",
			B: newInt(16),
			C: nil,
		},
		Q: resource.MustParse("205m"),
	}

	got, err := c.Flatten(obj, testObjectSchema())

	require.NoError(t, err)
	want := []any{map[string]any{
		"str":   "test-str",
		"alias": "test-alias",
		"int":   1,
		"float": 2.3,
		"bool":  true,
		"slice": []any{map[string]any{
			"a": "test-t",
		}, map[string]any{
			"a": "test-t-also",
		}},
		"map": map[string]any{
			"foo": 4,
		},
		"struct": []any{map[string]any{
			"a": "test-ptr-t",
			"b": []any{map[string]any{
				"value": 16,
			}},
		}},
		"q": "205m",
	}}
	assert.Equal(t, want, got)
}
