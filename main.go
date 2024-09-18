package main

import (
	"context"
	"flag"
	"log"

	"github.com/conduktor/terraform-provider-conduktor/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// Generate terraform provider code from the provider_code_spec.json file
// https://developer.hashicorp.com/terraform/plugin/code-generation/framework-generator
// Json spec format https://developer.hashicorp.com/terraform/plugin/code-generation/specification
//go:generate go run github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework generate all --input  provider_code_spec.json --output internal/schema

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name conduktor

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"
	commit  string = "none"
	date    string = "unknown"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/conduktor/conduktor",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version, commit, date), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
