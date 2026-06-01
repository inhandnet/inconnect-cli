package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/inhandnet/ics-cli/internal/factory"
	"github.com/inhandnet/ics-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func NewCmdAPI(f *factory.Factory) *cobra.Command {
	var method string

	cmd := &cobra.Command{
		Use:   "api <path>",
		Short: "Make an authenticated API request",
		Long:  "Make a raw API request to the InConnect server. Example: ics api /api/invpn/networks/list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			path := args[0]
			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			var body []byte
			switch strings.ToUpper(method) {
			case "GET", "":
				body, err = client.Get(path, q)
			case "POST":
				body, err = client.Do("POST", path, q, nil)
			case "PUT":
				body, err = client.Do("PUT", path, q, nil)
			case "DELETE":
				body, err = client.Do("DELETE", path, q, nil)
			default:
				return fmt.Errorf("unsupported method: %s", method)
			}

			if err != nil {
				if body != nil {
					_ = iostreams.FormatOutput(body, f.IO, f.IO.Output)
				}
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVarP(&method, "method", "X", "GET", "HTTP method")

	return cmd
}
