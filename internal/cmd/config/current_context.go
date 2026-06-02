package config

import (
	"fmt"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCmdCurrentContext(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "current-context",
		Short: "Show the current context name",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			name := cfg.ActiveContextName()
			if name == "" {
				fmt.Fprintln(f.IO.Out, "(none)")
			} else {
				fmt.Fprintln(f.IO.Out, name)
			}
			return nil
		},
	}
}
