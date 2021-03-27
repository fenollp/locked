package locked

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// DecodeFile FIXME
func DecodeFile(fn string) (*T, error) {
	src, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Failed to read configuration",
				Detail:   fmt.Sprintf("Can't read %s: %s.", fn, err),
			},
		}
	}

	file, diags := hclsyntax.ParseConfig(src, fn, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags
	}

	var t T
	if diags = gohcl.DecodeBody(file.Body, nil, &t); diags.HasErrors() {
		return nil, diags
	}

	return &t, nil
}
