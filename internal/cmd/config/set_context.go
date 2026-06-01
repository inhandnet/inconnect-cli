package config

import (
	"fmt"

	iconfig "github.com/inhandnet/ics-cli/internal/config"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCmdSetContext(f *factory.Factory) *cobra.Command {
	var host string

	cmd := &cobra.Command{
		Use:   "set-context <name>",
		Short: "Create or update a context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			name := args[0]
			ctx, ok := cfg.Contexts[name]
			if !ok {
				ctx = &iconfig.Context{}
			}
			if host != "" {
				ctx.Host = host
			}

			cfg.SetContext(name, ctx)
			if err := f.SaveConfig(); err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Context %q set\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&host, "host", "", "Host URL or region")
	_ = cmd.MarkFlagRequired("host")

	return cmd
}
