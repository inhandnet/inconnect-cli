package router

import (
	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

type connectionEventsOptions struct {
	cmdutil.ListFlags
	EventType string
	After     string
	Before    string
}

func newCmdConnectionEvents(f *factory.Factory) *cobra.Command {
	opts := &connectionEventsOptions{}

	cmd := &cobra.Command{
		Use:     "connection-events <device-id>",
		Aliases: []string{"conn-events"},
		Short:   "List a device's MQTT online/offline events",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewCursorQuery(cmd, &opts.ListFlags)
			if opts.EventType != "" {
				q.Set("eventType", opts.EventType)
			}
			if t := cmdutil.ParseTimeFlag(opts.After); t != "" {
				q.Set("after", t)
			}
			if t := cmdutil.ParseTimeFlag(opts.Before); t != "" {
				q.Set("before", t)
			}

			body, err := client.Get("/api/devices/"+args[0]+"/connection-events", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.RegisterPagination(cmd)
	cmd.Flags().StringVar(&opts.EventType, "event-type", "", "Filter by event type (e.g. online, offline)")
	cmd.Flags().StringVar(&opts.After, "after", "", "Only events at/after this time (e.g. 2024-01-01T00:00:00Z)")
	cmd.Flags().StringVar(&opts.Before, "before", "", "Only events before this time")

	return cmd
}
