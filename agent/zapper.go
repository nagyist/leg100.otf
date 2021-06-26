package agent

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func deleteBackendConfig(config []byte) (bool, []byte, error) {
	var deleted bool

	f, diags := hclwrite.ParseConfig([]byte(config), "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return false, nil, fmt.Errorf("unable to parse HCL: %s", diags.Error())
	}

OuterLoop:
	for _, block := range f.Body().Blocks() {
		if block.Type() == "terraform" {
			for _, b2 := range block.Body().Blocks() {
				if b2.Type() == "backend" {
					block.Body().RemoveBlock(b2)
					deleted = true
					break OuterLoop
				}
			}
		}
	}

	return deleted, f.Bytes(), nil
}
