package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/nitrado/terraform-provider-ec/ec/provider"
)

func main() {
	var debugMode bool
	var pluginPath string

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.StringVar(&pluginPath, "registry", "registry.terraform.io/nitrado/ec", "specify path, useful for local debugging")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: provider.Provider,
		Debug:        debugMode,
		ProviderAddr: pluginPath,
	}

	plugin.Serve(opts)
}
