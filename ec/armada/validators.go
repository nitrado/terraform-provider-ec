package armada

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

const maxNameLength = 63

var nameRegexp = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)

func validateName(value any, path cty.Path) (diags diag.Diagnostics) {
	v := value.(string)

	if !nameRegexp.MatchString(v) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       `"` + v + `" is not a valid name`,
			AttributePath: path,
		},
		)
	}
	if len(v) > maxNameLength {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("%q must be no more than %d characters", v, maxNameLength),
			AttributePath: path,
		},
		)
	}
	return diags
}
