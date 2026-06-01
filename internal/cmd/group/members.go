package group

import (
	"fmt"
	"net/url"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdAccounts(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "Manage accounts in a group",
	}

	cmd.AddCommand(
		newCmdAccountsList(f),
		newCmdAccountsUpdate(f),
	)

	return cmd
}

func newCmdAccountsList(f *factory.Factory) *cobra.Command {
	opts := &cmdutil.ListFlags{}

	cmd := &cobra.Command{
		Use:     "list <group-id>",
		Aliases: []string{"ls"},
		Short:   "List accounts in a group (use 'none' for unassigned)",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, opts)
			body, err := client.Get("/api/invpn/groups/"+args[0]+"/accounts", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	return cmd
}

func newCmdAccountsUpdate(f *factory.Factory) *cobra.Command {
	var addIDs, removeIDs []string

	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Add or remove accounts from a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if len(addIDs) > 0 {
				body["addAccountIds"] = addIDs
			}
			if len(removeIDs) > 0 {
				body["delAccountIds"] = removeIDs
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			_, err = client.Do("PUT", "/api/invpn/groups/"+args[0]+"/accounts", q, body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Group %s accounts updated\n", args[0])
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&addIDs, "add", nil, "Account IDs to add")
	cmd.Flags().StringSliceVar(&removeIDs, "remove", nil, "Account IDs to remove")

	return cmd
}

func newCmdRouters(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routers",
		Short: "Manage routers in a group",
	}

	cmd.AddCommand(
		newCmdRoutersList(f),
		newCmdRoutersUpdate(f),
	)

	return cmd
}

func newCmdRoutersList(f *factory.Factory) *cobra.Command {
	opts := &cmdutil.ListFlags{}

	cmd := &cobra.Command{
		Use:     "list <group-id>",
		Aliases: []string{"ls"},
		Short:   "List routers in a group (use 'none' for unassigned)",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, opts)
			body, err := client.Get("/api/invpn/groups/"+args[0]+"/routers", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	return cmd
}

func newCmdRoutersUpdate(f *factory.Factory) *cobra.Command {
	var addIDs, removeIDs []string

	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Add or remove routers from a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if len(addIDs) > 0 {
				body["addRouterIds"] = addIDs
			}
			if len(removeIDs) > 0 {
				body["delRouterIds"] = removeIDs
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			_, err = client.Do("PUT", "/api/invpn/groups/"+args[0]+"/routers", q, body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Group %s routers updated\n", args[0])
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&addIDs, "add", nil, "Router IDs to add")
	cmd.Flags().StringSliceVar(&removeIDs, "remove", nil, "Router IDs to remove")

	return cmd
}
