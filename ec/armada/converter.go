package armada

import (
	"github.com/nitrado/tfconv"
	"k8s.io/apimachinery/pkg/api/resource"
)

func converter() *tfconv.Converter {
	c := tfconv.New("json")
	c.Register(resource.Quantity{}, expandQuantity, flattenQuantity)
	return c
}

func expandQuantity(v any) (any, error) {
	return resource.ParseQuantity(v.(string))
}

func flattenQuantity(v any) (any, error) {
	q := v.(resource.Quantity)
	return (&q).String(), nil
}
