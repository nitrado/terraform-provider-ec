package resource

import (
	"errors"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SetData sets the given data into a schema ResourceData.
func SetData(d *schema.ResourceData, v any) error {
	if a, ok := v.([]any); ok {
		if len(a) == 0 {
			return errors.New("unexpected empty []any")
		}
		v = a[0]
	}
	m, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("expected map, got %T", v)
	}
	for k, val := range m {
		if err := d.Set(k, val); err != nil {
			return err
		}
	}
	return nil
}

// SetDataKeys sets only the given keys of the given data into a schema ResourceData.
func SetDataKeys(d *schema.ResourceData, v any, keys ...string) error {
	if a, ok := v.([]any); ok {
		if len(a) == 0 {
			return errors.New("unexpected empty []any")
		}
		v = a[0]
	}
	m, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("expected map, got %T", v)
	}
	for k, val := range m {
		if !slices.Contains(keys, k) {
			continue
		}
		if err := d.Set(k, val); err != nil {
			return err
		}
	}
	return nil
}
