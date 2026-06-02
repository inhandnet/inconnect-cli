package config

import (
	"fmt"

	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCmdUseContext(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "use-context <name>",
		Short: "Switch the active context",
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

			cfg.CurrentContext = name
			if err := f.SaveConfig(); err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Switched to context %q\n", name)
			return nil
		},
	}
}
