package main

import (
	"fmt"
	"io"
	"os"

	armadav1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/apis/armada/v1"
	containerv1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/apis/container/v1"
)

type objInfo struct {
	Obj      any
	Filename string
	FuncName string
}

var objs = []objInfo{
	{
		Obj:      armadav1.Resources{},
		Filename: "ec/armada/schema_resources.go",
		FuncName: "resourcesSchema",
	},
	{
		Obj:      armadav1.EnvVar{},
		Filename: "ec/armada/schema_env.go",
		FuncName: "envSchema",
	},
	{
		Obj:      &armadav1.Site{},
		Filename: "ec/armada/schema_site.go",
		FuncName: "siteSchema",
	},
	{
		Obj:      &armadav1.Region{},
		Filename: "ec/armada/schema_region.go",
		FuncName: "regionSchema",
	},
	{
		Obj:      &armadav1.Armada{},
		Filename: "ec/armada/schema_armada.go",
		FuncName: "armadaSchema",
	},
	{
		Obj:      &armadav1.ArmadaSet{},
		Filename: "ec/armada/schema_armadaset.go",
		FuncName: "armadaSetSchema",
	},
	{
		Obj:      &containerv1.Branch{},
		Filename: "ec/armada/schema_branch.go",
		FuncName: "branchSchema",
	},
}

func main() {
	os.Exit(realMain(os.Stderr))
}

func realMain(out io.Writer) int {
	gen := NewGenerator()

	for _, info := range objs {
		// Remove the file if it exists.
		if _, err := os.Stat(info.Filename); err != nil {
			_ = os.Remove(info.Filename)
		}

		b, err := gen.Generate(info.Obj, "armada", info.FuncName)
		if err != nil {
			_, _ = fmt.Fprintln(out, err.Error())
			continue
		}

		//nolint:gosec // The mask 0o644 is fine.
		if err = os.WriteFile(info.Filename, b, 0o644); err != nil {
			_, _ = fmt.Fprintf(out, "Could not write schema file %q: %s\n", info.Filename, err.Error())
			return 1
		}
	}
	return 0
}
