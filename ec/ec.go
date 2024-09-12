package ec

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ettle/strcase"
	"github.com/nitrado/tfconv"
	"gitlab.com/nitrado/b2b/ec/core/pkg/apiclient/clientset"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ProviderContext contains connection context information.
type ProviderContext struct {
	defaultCS clientset.Interface
	instances map[string]clientset.Interface
}

// NewProviderContext returns a provider context with the given default clientset and instances.
func NewProviderContext(defCS clientset.Interface, instances map[string]clientset.Interface) ProviderContext {
	if instances == nil {
		instances = map[string]clientset.Interface{}
	}
	return ProviderContext{
		defaultCS: defCS,
		instances: instances,
	}
}

// ResolveClientSet resolves the ClientSet from the given context.
func ResolveClientSet(m any, name string) (clientset.Interface, error) {
	connCtx, ok := m.(ProviderContext)
	if !ok {
		return nil, errors.New("invalid connection context")
	}

	if name == "" {
		if connCtx.defaultCS == nil {
			return nil, errors.New("no default clientset found")
		}
		return connCtx.defaultCS, nil
	}

	cs, ok := connCtx.instances[name]
	if !ok || cs == nil {
		return nil, fmt.Errorf("instance %q clientset not found", name)
	}
	return cs, nil
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
	c := tfconv.NewWithName(FieldName, "json")
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

// FieldName returns the terraform-styled field name from the given name.
func FieldName(name string) string {
	name = strcase.ToSnake(name)
	if strings.Contains(name, "cid_rs") {
		name = strings.ReplaceAll(name, "cid_rs", "cidrs")
	}
	return name
}
