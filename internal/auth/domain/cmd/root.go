// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/ory/x/cmdx"
	"github.com/ory/x/dbal"
	"github.com/ory/x/jsonnetsecure"
	"github.com/ory/x/profilex"
	"my.com/secrets/internal/auth/domain/cmd/cleanup"
	"my.com/secrets/internal/auth/domain/cmd/courier"
	"my.com/secrets/internal/auth/domain/cmd/hashers"
	"my.com/secrets/internal/auth/domain/cmd/identities"
	"my.com/secrets/internal/auth/domain/cmd/jsonnet"
	"my.com/secrets/internal/auth/domain/cmd/migrate"
	"my.com/secrets/internal/auth/domain/cmd/remote"
	"my.com/secrets/internal/auth/domain/cmd/serve"
	"my.com/secrets/internal/auth/domain/driver"
	"my.com/secrets/internal/auth/domain/driver/config"
)

func NewRootCmd(driverOpts ...driver.RegistryOption) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use: "kratos",
	}
	cmdx.EnableUsageTemplating(cmd)

	courier.RegisterCommandRecursive(cmd, nil, driverOpts)
	cmd.AddCommand(identities.NewGetCmd())
	cmd.AddCommand(identities.NewDeleteCmd())
	cmd.AddCommand(jsonnet.NewFormatCmd())
	hashers.RegisterCommandRecursive(cmd)
	cmd.AddCommand(identities.NewImportCmd())
	cmd.AddCommand(jsonnet.NewLintCmd())
	cmd.AddCommand(identities.NewListCmd())
	migrate.RegisterCommandRecursive(cmd)
	serve.RegisterCommandRecursive(cmd, nil, driverOpts)
	cleanup.RegisterCommandRecursive(cmd)
	remote.RegisterCommandRecursive(cmd)
	cmd.AddCommand(identities.NewValidateCmd())
	cmd.AddCommand(cmdx.Version(&config.Version, &config.Commit, &config.Date))

	// Registers a hidden "jsonnet" subcommand for process-isolated Jsonnet VMs.
	cmd.AddCommand(jsonnetsecure.NewJsonnetCmd())

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() int {
	defer profilex.Profile().Stop()

	dbal.RegisterDriver(func() dbal.Driver {
		return driver.NewRegistryDefault()
	})

	jsonnetPool := jsonnetsecure.NewProcessPool(runtime.GOMAXPROCS(0))
	defer jsonnetPool.Close()

	c := NewRootCmd(driver.WithJsonnetPool(jsonnetPool))

	if err := c.Execute(); err != nil {
		if !errors.Is(err, cmdx.ErrNoPrintButFail) {
			_, _ = fmt.Fprintln(c.ErrOrStderr(), err)
		}
		return 1
	}
	return 0
}
