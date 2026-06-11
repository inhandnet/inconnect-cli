package cmd

import (
	"os"
	"strconv"

	"github.com/inhandnet/inconnect-cli/internal/build"
	"github.com/inhandnet/inconnect-cli/internal/debug"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/spf13/cobra"
)

func NewCmdRoot(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "inconnect",
		Short:         "InConnect CLI — manage VPN networks, servers, and routers",
		Version:       build.Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if v, _ := cmd.Flags().GetBool("debug"); v {
				debug.Enabled = true
			}
			if os.Getenv("INCONNECT_DEBUG") != "" {
				debug.Enabled = true
			}

			output, _ := cmd.Flags().GetString("output")
			jq, _ := cmd.Flags().GetString("jq")
			if jq != "" {
				f.IO.JQ = jq
				if !cmd.Flags().Changed("output") {
					output = "json"
				}
			}
			if output == "" {
				if f.IO.IsTTY {
					output = "table"
				} else {
					output = "json"
				}
			}
			f.IO.Output = output

			if cols, _ := cmd.Flags().GetStringSlice("columns"); len(cols) > 0 {
				f.IO.Columns = cols
			}

			if ctx, _ := cmd.Flags().GetString("context"); ctx != "" {
				os.Setenv("INCONNECT_CONTEXT", ctx)
			}

			// Resolve the effective org ID: an explicit --oid wins; otherwise fall
			// back to the active context's saved OrgID (set via 'auth switch-org').
			// We backfill both the flag and INCONNECT_OID so every read path sees it.
			oid, _ := cmd.Flags().GetString("oid")
			if oid == "" {
				if cfg, err := f.Config(); err == nil {
					if actx, err := cfg.ActiveContext(); err == nil && actx.OrgID != "" {
						oid = actx.OrgID
						_ = cmd.Flags().Set("oid", oid)
					}
				}
			}
			if oid != "" {
				os.Setenv("INCONNECT_OID", oid)
			}
			if cmd.Flags().Changed("verbose") {
				v, _ := cmd.Flags().GetInt("verbose")
				os.Setenv("INCONNECT_VERBOSE", strconv.Itoa(v))
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringP("output", "o", "", "Output format: json, table, yaml (default: table for TTY, json otherwise)")
	cmd.PersistentFlags().StringSliceP("columns", "c", nil, "Table columns to show (comma-separated dot-paths; prefix ! to exclude)")
	cmd.PersistentFlags().String("jq", "", "Filter JSON output using a jq expression")
	cmd.PersistentFlags().String("oid", "", "Organization ID")
	cmd.PersistentFlags().String("context", "", "Config context to use")
	cmd.PersistentFlags().Bool("debug", false, "Enable debug output")
	cmd.PersistentFlags().Int("verbose", 100, "API field verbosity for GET requests (1-100, higher = more fields; 0 to omit)")

	return cmd
}
