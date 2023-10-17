package armada

import (
	"github.com/nitrado/terraform-provider-ec/internal/terra"
	"k8s.io/apimachinery/pkg/api/resource"
)

func converter() *terra.Converter {
	c := terra.New()
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
