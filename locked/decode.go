package locked

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const lckfl = "Lockfile"

func decodeFile(fn string) (*T, error) {
	src, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	t, err := parseFile(fn, src)
	if err != nil {
		return nil, err
	}

	dst := hclwrite.Format(src)
	if !bytes.Equal(src, dst) {
		if err := os.WriteFile(fn, dst, 0644); err != nil { //TODO: use existing rights, not 0644
			return nil, err
		}
	}

	return t, nil
}

func parseFile(fn string, src []byte) (*T, error) {
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

// Load looks for Lockfile in the current and parent directories and finally $XDG_CONFIG_HOME.
// If the file exists (and is authorized), it is loaded and merged into FIXME
func Load() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	// Fail early in case of rights problems
	gfname, err := xdg.ConfigFile(lckfl)
	if err != nil {
		return err
	}

	fn := filepath.Join(wd, lckfl)
	t0, err := decodeFile(fn)
	if err != nil {
		return err
	}
	fmt.Printf(">>> t0 %#v\n", t0)

	parent := func(p string) string { return filepath.Join(filepath.Dir(filepath.Dir(p)), lckfl) }
	for previous, fname := fn, parent(fn); fname != previous; previous, fname = fname, parent(fname) {
		fmt.Println(">>> trying", fname)
		ti, err := decodeFile(fname)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err // FIXME: break on 403, return otherwise
		}
		fmt.Printf(">>> ti %#v\n", ti)
	}

	tg, err := decodeFile(gfname)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err // FIXME: break on 403, return otherwise
	}
	fmt.Printf(">>> tg %#v\n", tg)

	fmt.Printf(">>> prelude %#v\n", prelude)
	return nil
}
