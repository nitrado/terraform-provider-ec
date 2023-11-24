package ec

import (
	"errors"
	"strings"

	"github.com/nitrado/tfconv"
	"gitlab.com/nitrado/b2b/ec/armada/pkg/apiclient/clientset"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ResolveClientSet resolves the ClientSet from the given context.
func ResolveClientSet(m any) (clientset.Interface, error) {
	clientSet, ok := m.(clientset.Interface)
	if !ok {
		return nil, errors.New("invalid clientset")
	}
	return clientSet, nil
}

// ScopedName returns the encoded name of an object.
func ScopedName(env, name string) string {
	if env != "" {
		return env + "/" + name
	}
	return name
}

// SplitName decodes the key into its parts.
func SplitName(key string) (env, name string) {
	parts := strings.SplitN(key, "/", 2)
	switch len(parts) {
	case 1:
		return "", parts[0]
	default:
		return parts[0], parts[1]
	}
}

// Converter returns the configured converter.
func Converter() *tfconv.Converter {
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
