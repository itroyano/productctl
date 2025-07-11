package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/opdev/productctl/internal/cli"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/apply"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/archivecomponent"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/bridge"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/cleanup"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/create"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/deleteproductlisting"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/fetch"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/jsonschema"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/sanitize"
	"github.com/opdev/productctl/internal/cmd/productctl/cmd/version"
	libversion "github.com/opdev/productctl/internal/version"
)

// Execute runs the top-most command structure of the CLI.
func Execute() error {
	return RootCmd().ExecuteContext(context.Background())
}

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "productctl",
		Short:             "A CLI for managing Red Hat Partner Product Listings.",
		Long:              "A basic CLI useful for helping Red Hat Certification Partners define their Product Listings, and create and manage certification projects associated with those product listings.",
		Version:           libversion.Version.String(),
		PersistentPreRunE: configureCLIPreRunE,
	}

	cmd.AddCommand(version.Command())
	cmd.PersistentFlags().String(cli.FlagIDLogLevel, "info", "The verbosity of the tool itself. Ex. error, warn, info, debug")

	util := bridge.Command("util", "Utilities for the management of your Partner Connect account")
	util.PersistentFlags().String(cli.FlagIDEndpoint, "prod", "The catalog API environment to use. Choose from stage, prod")
	util.PersistentFlags().String(cli.FlagIDCustomEndpoint, "", "Define a custom API endpoint. Supersedes predefined environment values like \"prod\" if set")
	util.AddCommand(archivecomponent.Command())
	util.AddCommand(deleteproductlisting.Command())
	cmd.AddCommand(util)

	// Build the product management command tree.
	product := bridge.Command("product", "Manage your Product Listing")
	product.PersistentFlags().String(cli.FlagIDEndpoint, "prod", "The catalog API environment to use. Choose from stage, prod")
	product.PersistentFlags().String(cli.FlagIDCustomEndpoint, "", "Define a custom API endpoint. Supersedes predefined environment values like \"prod\" if set")
	product.AddCommand(create.Command())
	product.AddCommand(apply.Command())
	product.AddCommand(fetch.Command())
	product.AddCommand(sanitize.Command())
	product.AddCommand(cleanup.Command())
	product.AddCommand(jsonschema.Command())
	cmd.AddCommand(product)

	return cmd
}

var ErrConfiguringCLI = errors.New("failed to configure CLI")

func configureCLIPreRunE(cmd *cobra.Command, args []string) error {
	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Join(ErrConfiguringCLI, err)
	}

	loglevel, err := cmd.Flags().GetString(cli.FlagIDLogLevel)
	if err != nil {
		return errors.Join(ErrConfiguringCLI, err)
	}
	ctx, _, err := cli.ConfigureLogger(loglevel, os.Stderr)
	if err != nil {
		return errors.Join(ErrConfiguringCLI, err)
	}

	cmd.SetContext(ctx)

	return nil
}
