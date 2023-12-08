package main

import (
	"fmt"
	"io"
	"os"

	armadav1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/armada/v1"
	containerv1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/container/v1"
	corev1 "gitlab.com/nitrado/b2b/ec/armada/pkg/api/core/v1"
)

type objInfo struct {
	Pkg      string
	Obj      any
	Filename string
	FuncName string
}

var objs = []objInfo{
	{
		Pkg:      "armada",
		Obj:      corev1.EnvVar{},
		Filename: "ec/armada/schema_env.go",
		FuncName: "envSchema",
	},
	{
		Pkg:      "armada",
		Obj:      &armadav1.Armada{},
		Filename: "ec/armada/schema_armada.go",
		FuncName: "armadaSchema",
	},
	{
		Pkg:      "armada",
		Obj:      &armadav1.ArmadaSet{},
		Filename: "ec/armada/schema_armadaset.go",
		FuncName: "armadaSetSchema",
	},
	{
		Pkg:      "container",
		Obj:      &containerv1.Branch{},
		Filename: "ec/container/schema_branch.go",
		FuncName: "branchSchema",
	},
	{
		Pkg:      "core",
		Obj:      corev1.Resources{},
		Filename: "ec/core/schema_resources.go",
		FuncName: "resourcesSchema",
	},
	{
		Pkg:      "core",
		Obj:      corev1.EnvVar{},
		Filename: "ec/core/schema_env.go",
		FuncName: "envSchema",
	},
	{
		Pkg:      "core",
		Obj:      &corev1.Environment{},
		Filename: "ec/core/schema_environment.go",
		FuncName: "environmentSchema",
	},
	{
		Pkg:      "core",
		Obj:      &corev1.Region{},
		Filename: "ec/core/schema_region.go",
		FuncName: "regionSchema",
	},
	{
		Pkg:      "core",
		Obj:      &corev1.Site{},
		Filename: "ec/core/schema_site.go",
		FuncName: "siteSchema",
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

		b, err := gen.Generate(info.Obj, info.Pkg, info.FuncName)
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
