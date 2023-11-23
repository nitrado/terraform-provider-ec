package meta

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

const (
	maxNameLength        = 63
	maxEnvironmentLength = 4
)

var (
	nameRegexp        = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)
	environmentRegexp = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)
)

func validateName(value any, path cty.Path) (diags diag.Diagnostics) {
	v := value.(string)

	if len(v) > maxNameLength {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("%q must be no more than %d characters", v, maxNameLength),
			AttributePath: path,
		})
	}
	if !nameRegexp.MatchString(v) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       `"` + v + `" is not a valid name`,
			AttributePath: path,
		})
	}
	return diags
}

func validateEnvironment(value any, path cty.Path) (diags diag.Diagnostics) {
	v := value.(string)

	if len(v) > maxEnvironmentLength {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("%q must be no more than %d characters", v, maxNameLength),
			AttributePath: path,
		})
	}
	if !environmentRegexp.MatchString(v) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       `"` + v + `" is not a valid environment`,
			AttributePath: path,
		})
	}
	return diags
}
