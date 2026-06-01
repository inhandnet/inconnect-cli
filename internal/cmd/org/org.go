package org

import (
	"net/url"

	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func NewCmdOrg(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org",
		Short: "Manage organization settings",
	}

	cmd.AddCommand(
		newCmdList(f),
		newCmdGet(f),
		newCmdCreate(f),
		newCmdDelete(f),
		newCmdExport(f),
		newCmdSettings(f),
		newCmdUpdateSettings(f),
	)

	return cmd
}

func newCmdSettings(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "settings",
		Short: "Get current organization settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			body, err := client.Get("/api/invpn/org/settings", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}

type updateSettingsOptions struct {
	RealIPEnabled               bool
	DiffAccountAndDeviceIPPool  bool
	DataThresholdInBytes        int64
	NotificationEmails          []string
}

func newCmdUpdateSettings(f *factory.Factory) *cobra.Command {
	opts := &updateSettingsOptions{}

	cmd := &cobra.Command{
		Use:   "update-settings <org-id>",
		Short: "Update organization settings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if cmd.Flags().Changed("real-ip-enabled") {
				body["realIpEnabled"] = opts.RealIPEnabled
			}
			if cmd.Flags().Changed("diff-ip-pool") {
				body["diffAccountAndDeviceIpPool"] = opts.DiffAccountAndDeviceIPPool
			}
			if cmd.Flags().Changed("data-threshold") {
				body["dataThresholdInBytes"] = opts.DataThresholdInBytes
			}
			if cmd.Flags().Changed("notification-emails") {
				body["notificationEmails"] = opts.NotificationEmails
			}

			respBody, err := client.Put("/api/invpn/org/"+args[0]+"/settings", body)
			if err != nil {
				if respBody != nil {
					_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				}
				return err
			}

			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().BoolVar(&opts.RealIPEnabled, "real-ip-enabled", false, "Enable real IP for routers")
	cmd.Flags().BoolVar(&opts.DiffAccountAndDeviceIPPool, "diff-ip-pool", false, "Separate account and device IP pools")
	cmd.Flags().Int64Var(&opts.DataThresholdInBytes, "data-threshold", 0, "Data usage threshold in bytes for notification")
	cmd.Flags().StringSliceVar(&opts.NotificationEmails, "notification-emails", nil, "Notification emails (max 5)")

	return cmd
}
