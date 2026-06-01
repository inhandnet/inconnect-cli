package config

import (
	"fmt"
	"sort"

	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/spf13/cobra"
)

func newCmdListContexts(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list-contexts",
		Aliases: []string{"ls"},
		Short:   "List all contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := f.Config()
			if err != nil {
				return err
			}

			current := cfg.ActiveContextName()
			var names []string
			for name := range cfg.Contexts {
				names = append(names, name)
			}
			sort.Strings(names)

			fmt.Fprintf(f.IO.Out, "%-8s %-16s %-40s %s\n", "CURRENT", "NAME", "HOST", "USER")
			for _, name := range names {
				ctx := cfg.Contexts[name]
				marker := ""
				if name == current {
					marker = "*"
				}
				fmt.Fprintf(f.IO.Out, "%-8s %-16s %-40s %s\n", marker, name, ctx.Host, ctx.User)
			}
			return nil
		},
	}
}
