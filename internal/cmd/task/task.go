package task

import (
	"fmt"

	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

var statusMap = map[string]string{
	"running":   "1",
	"waiting":   "0,4,5",
	"failed":    "-1,2",
	"completed": "3",
}

func NewCmdTask(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Manage tasks",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdCancel(f),
		newCmdRestart(f),
	)

	return cmd
}

type listOptions struct {
	cmdutil.ListFlags
	Status   string
	ObjectID string
}

func newCmdList(f *factory.Factory) *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, &opts.ListFlags)
			if opts.Status != "" {
				if mapped, ok := statusMap[opts.Status]; ok {
					q.Set("states", mapped)
				} else {
					q.Set("states", opts.Status)
				}
			}
			if opts.ObjectID != "" {
				q.Set("object_id", opts.ObjectID)
			}

			body, err := client.Get("/api2/tasks", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.Status, "status", "", "Filter by status (running, waiting, failed, completed)")
	cmd.Flags().StringVar(&opts.ObjectID, "object-id", "", "Filter by object ID")

	return cmd
}

func newCmdCancel(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "cancel <id>",
		Short: "Cancel a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			_, err = client.Delete("/api2/tasks/" + args[0])
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Task %s cancelled\n", args[0])
			return nil
		},
	}
}

func newCmdRestart(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "restart <id>",
		Short: "Restart a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			_, err = client.Post("/api2/tasks/"+args[0]+"/restart", nil)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Task %s restarted\n", args[0])
			return nil
		},
	}
}

