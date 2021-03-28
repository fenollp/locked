package locked

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// DecodeFile FIXME
func DecodeFile(fn string) (*T, error) {
	src, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Failed to read lockfile",
				Detail:   fmt.Sprintf("Can't read %s: %s.", fn, err),
			},
		}
	}

	t, err := decodeFile(fn, src)
	if err != nil {
		return nil, err
	}

	dst := hclwrite.Format(src)
	if !bytes.Equal(src, dst) {
		if err := ioutil.WriteFile(fn, dst, 0644); err != nil {
			return nil, hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Failed to format lockfile",
					Detail:   fmt.Sprintf("Can't write %s: %s.", fn, err),
				},
			}
		}
	}

	return t, nil
}

func decodeFile(fn string, src []byte) (*T, error) {
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
