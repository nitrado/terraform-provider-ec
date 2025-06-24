package main

import (
	"fmt"
	"io"
	"os"

	armadav1 "github.com/gamefabric/gf-core/pkg/api/armada/v1"
	authenticationv1beta1 "github.com/gamefabric/gf-core/pkg/api/authentication/v1beta1"
	containerv1 "github.com/gamefabric/gf-core/pkg/api/container/v1"
	corev1 "github.com/gamefabric/gf-core/pkg/api/core/v1"
	formationv1 "github.com/gamefabric/gf-core/pkg/api/formation/v1"
	protectionv1 "github.com/gamefabric/gf-core/pkg/api/protection/v1"
	storagev1beta1 "github.com/gamefabric/gf-core/pkg/api/storage/v1beta1"
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
		Pkg:      "authentication",
		Obj:      &authenticationv1beta1.Provider{},
		Filename: "ec/authentication/schema_provider.go",
		FuncName: "providerSchema",
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
		Obj:      &corev1.Location{},
		Filename: "ec/core/schema_location.go",
		FuncName: "locationSchema",
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
	{
		Pkg:      "formation",
		Obj:      corev1.EnvVar{},
		Filename: "ec/formation/schema_env.go",
		FuncName: "envSchema",
	},
	{
		Pkg:      "formation",
		Obj:      &formationv1.Formation{},
		Filename: "ec/formation/schema_formation.go",
		FuncName: "formationSchema",
	},
	{
		Pkg:      "formation",
		Obj:      &formationv1.Vessel{},
		Filename: "ec/formation/schema_vessel.go",
		FuncName: "vesselSchema",
	},
	{
		Pkg:      "protection",
		Obj:      &protectionv1.GatewayPolicy{},
		Filename: "ec/protection/schema_gatewaypolicy.go",
		FuncName: "gatewayPolicySchema",
	},
	{
		Pkg:      "protection",
		Obj:      &protectionv1.Protocol{},
		Filename: "ec/protection/schema_protocol.go",
		FuncName: "protocolSchema",
	},
	{
		Pkg:      "protection",
		Obj:      &protectionv1.Mitigation{},
		Filename: "ec/protection/schema_migration.go",
		FuncName: "migrationSchema",
	},
	{
		Pkg:      "storage",
		Obj:      &storagev1beta1.VolumeStore{},
		Filename: "ec/storage/schema_volumestore.go",
		FuncName: "volumeStoreSchema",
	},
	{
		Pkg:      "storage",
		Obj:      &storagev1beta1.VolumeStoreRetentionPolicy{},
		Filename: "ec/storage/schema_volumestoreretentionpolicy.go",
		FuncName: "volumeStoreRetentionPolicySchema",
	},
	{
		Pkg:      "storage",
		Obj:      &storagev1beta1.Volume{},
		Filename: "ec/storage/schema_volume.go",
		FuncName: "volumeSchema",
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
