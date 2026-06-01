package network

import (
	"fmt"
	"net/url"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdMembers(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "members <network-id>",
		Short: "Update network members (routers and accounts)",
		Args:  cobra.ExactArgs(1),
	}

	cmd.AddCommand(
		newCmdMembersUpdate(f),
	)

	return cmd
}

func newCmdMembersUpdate(f *factory.Factory) *cobra.Command {
	var addRouters, delRouters, addAccounts, delAccounts []string
	var transferTo string

	cmd := &cobra.Command{
		Use:   "update <network-id>",
		Short: "Add or remove routers and accounts from a network",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if len(addRouters) > 0 {
				body["addRouterIds"] = addRouters
			}
			if len(delRouters) > 0 {
				body["delRouterIds"] = delRouters
			}
			if len(addAccounts) > 0 {
				body["addAccountIds"] = addAccounts
			}
			if len(delAccounts) > 0 {
				body["delAccountIds"] = delAccounts
			}
			if transferTo != "" {
				body["delTransferToNetworkId"] = transferTo
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			_, err = client.Do("PUT", "/api/invpn/networks/"+args[0]+"/members", q, body)
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Network %s members updated\n", args[0])
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&addRouters, "add-routers", nil, "Router IDs to add")
	cmd.Flags().StringSliceVar(&delRouters, "del-routers", nil, "Router IDs to remove")
	cmd.Flags().StringSliceVar(&addAccounts, "add-accounts", nil, "Account IDs to add")
	cmd.Flags().StringSliceVar(&delAccounts, "del-accounts", nil, "Account IDs to remove")
	cmd.Flags().StringVar(&transferTo, "transfer-to", "", "Network ID to transfer removed members to")

	return cmd
}

func newCmdAccounts(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accounts",
		Short: "List accounts in a network",
	}

	cmd.AddCommand(
		newCmdAccountsList(f),
		newCmdDefaultAccountsList(f),
	)

	return cmd
}

func newCmdAccountsList(f *factory.Factory) *cobra.Command {
	opts := &cmdutil.ListFlags{}
	var email string

	cmd := &cobra.Command{
		Use:     "list <network-id>",
		Aliases: []string{"ls"},
		Short:   "List accounts in a network",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, opts)
			if email != "" {
				q.Set("email", email)
			}

			body, err := client.Get("/api/invpn/networks/"+args[0]+"/accounts", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&email, "email", "", "Filter by email")

	return cmd
}

func newCmdDefaultAccountsList(f *factory.Factory) *cobra.Command {
	opts := &cmdutil.ListFlags{}
	var email string

	cmd := &cobra.Command{
		Use:   "default",
		Short: "List accounts in the default network",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, opts)
			if email != "" {
				q.Set("email", email)
			}

			body, err := client.Get("/api/invpn/networks/default/accounts", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&email, "email", "", "Filter by email")

	return cmd
}

func newCmdRouters(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routers",
		Short: "List routers in a network",
	}

	cmd.AddCommand(
		newCmdRoutersList(f),
		newCmdDefaultRoutersList(f),
	)

	return cmd
}

func newCmdRoutersList(f *factory.Factory) *cobra.Command {
	opts := &cmdutil.ListFlags{}
	var serialNumber string

	cmd := &cobra.Command{
		Use:     "list <network-id>",
		Aliases: []string{"ls"},
		Short:   "List routers in a network",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, opts)
			if serialNumber != "" {
				q.Set("serialNumber", serialNumber)
			}

			body, err := client.Get("/api/invpn/networks/"+args[0]+"/routers", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&serialNumber, "serial", "", "Filter by serial number")

	return cmd
}

func newCmdDefaultRoutersList(f *factory.Factory) *cobra.Command {
	opts := &cmdutil.ListFlags{}
	var serialNumber, center string

	cmd := &cobra.Command{
		Use:   "default",
		Short: "List routers in the default network",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, opts)
			if serialNumber != "" {
				q.Set("serialNumber", serialNumber)
			}
			if center != "" {
				q.Set("center", center)
			}

			body, err := client.Get("/api/invpn/networks/default/routers", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&serialNumber, "serial", "", "Filter by serial number")
	cmd.Flags().StringVar(&center, "center", "", "Filter by center")

	return cmd
}

func newCmdEndpoints(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "endpoints",
		Short: "List endpoints in a network",
	}

	cmd.AddCommand(
		newCmdEndpointsList(f),
		newCmdDefaultEndpointsList(f),
	)

	return cmd
}

func newCmdEndpointsList(f *factory.Factory) *cobra.Command {
	opts := &cmdutil.ListFlags{}
	var vip, rip string

	cmd := &cobra.Command{
		Use:     "list <network-id>",
		Aliases: []string{"ls"},
		Short:   "List endpoints in a network",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, opts)
			if vip != "" {
				q.Set("vip", vip)
			}
			if rip != "" {
				q.Set("rip", rip)
			}

			body, err := client.Get("/api/invpn/networks/"+args[0]+"/endpoints", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&vip, "vip", "", "Filter by VIP")
	cmd.Flags().StringVar(&rip, "rip", "", "Filter by real IP")

	return cmd
}

func newCmdDefaultEndpointsList(f *factory.Factory) *cobra.Command {
	opts := &cmdutil.ListFlags{}
	var vip, rip string

	cmd := &cobra.Command{
		Use:   "default",
		Short: "List endpoints in the default network",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, opts)
			if vip != "" {
				q.Set("vip", vip)
			}
			if rip != "" {
				q.Set("rip", rip)
			}

			body, err := client.Get("/api/invpn/networks/default/endpoints", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&vip, "vip", "", "Filter by VIP")
	cmd.Flags().StringVar(&rip, "rip", "", "Filter by real IP")

	return cmd
}

func newCmdCenters(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "centers <network-id>",
		Short: "Get center routers of a network",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			body, err := client.Get("/api/invpn/networks/"+args[0]+"/centers", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}
