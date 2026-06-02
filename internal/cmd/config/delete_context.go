package config

import (
	"fmt"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCmdDeleteContext(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-context <name>",
		Short: "Delete a context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			name := args[0]
			if _, ok := cfg.Contexts[name]; !ok {
				return fmt.Errorf("context %q not found", name)
			}

			cfg.DeleteContext(name)
			if err := f.SaveConfig(); err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Context %q deleted\n", name)
			return nil
		},
	}
}
