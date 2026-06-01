package firmware

import (
	"fmt"
	"net/url"

	"github.com/inhandnet/ics-cli/internal/cmdutil"
	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdDevices(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices",
		Short: "Manage devices in a firmware upgrade job",
	}

	cmd.AddCommand(
		newCmdDevicesList(f),
		newCmdDevicesAdd(f),
		newCmdDevicesRemove(f),
	)

	return cmd
}

type devicesListOptions struct {
	cmdutil.ListFlags
	Status       string
	SerialNumber string
}

func newCmdDevicesList(f *factory.Factory) *cobra.Command {
	opts := &devicesListOptions{}

	cmd := &cobra.Command{
		Use:     "list <firmware-id>",
		Aliases: []string{"ls"},
		Short:   "List devices in a firmware upgrade job",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, &opts.ListFlags)
			if opts.Status != "" {
				q.Set("status", opts.Status)
			}
			if opts.SerialNumber != "" {
				q.Set("serialNumber", opts.SerialNumber)
			}

			body, err := client.Get("/api/job/"+args[0]+"/devices", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.Status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&opts.SerialNumber, "serial", "", "Filter by serial number")
	return cmd
}

func newCmdDevicesAdd(f *factory.Factory) *cobra.Command {
	var groupIDs []string

	cmd := &cobra.Command{
		Use:   "add <firmware-id> [device-id...]",
		Short: "Add devices to a firmware upgrade job",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if len(args) > 1 {
				body["deviceIds"] = args[1:]
			}
			if len(groupIDs) > 0 {
				body["deviceGroupIds"] = groupIDs
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			respBody, err := client.Do("POST", "/api/firmware/"+args[0]+"/devices", q, body)
			if err != nil {
				if respBody != nil {
					_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				}
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Devices added to firmware job %s\n", args[0])
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&groupIDs, "group", nil, "Device group IDs")
	return cmd
}

func newCmdDevicesRemove(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <firmware-id> <device-id>",
		Short: "Remove a device from a firmware upgrade job",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			_, err = client.Delete("/api/job/" + args[0] + "/devices/" + args[1])
			if err != nil {
				return err
			}

			fmt.Fprintf(f.IO.ErrOut, "Device %s removed from firmware job %s\n", args[1], args[0])
			return nil
		},
	}
}
