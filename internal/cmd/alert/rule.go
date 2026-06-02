package alert

import (
	"net/url"

	"github.com/inhandnet/inconnect-cli/internal/cmdutil"
	"github.com/inhandnet/inconnect-cli/internal/factory"
	"github.com/inhandnet/inconnect-cli/internal/iostreams"
	"github.com/spf13/cobra"
)

func newCmdRule(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "Manage alert rules",
	}

	cmd.AddCommand(
		newCmdRuleList(f),
		newCmdRuleGet(f),
		newCmdRuleCreate(f),
		newCmdRuleUpdate(f),
		newCmdRuleDelete(f),
	)

	return cmd
}

func newCmdRuleList(f *factory.Factory) *cobra.Command {
	opts := &cmdutil.ListFlags{}
	var routerName string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List alert rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := cmdutil.NewQuery(cmd, opts)
			if routerName != "" {
				q.Set("routerName", routerName)
			}

			body, err := client.Get("/api/invpn/alerts/rules", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&routerName, "router-name", "", "Filter by router name")
	return cmd
}

func newCmdRuleGet(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an alert rule by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body, err := client.Get("/api/invpn/alerts/rules/"+args[0], url.Values{})
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, f.IO.Output)
		},
	}
}

type ruleCreateOptions struct {
	RuleType       string
	RouterIDs      []string
	AlertType      string
	Retention      int64
	NotifyChannels []string
	NotifyUsers    []string
}

func newCmdRuleCreate(f *factory.Factory) *cobra.Command {
	opts := &ruleCreateOptions{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an alert rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			rule := map[string]any{
				"type": opts.AlertType,
			}
			if opts.Retention > 0 {
				rule["param"] = map[string]any{"retention": opts.Retention}
			}

			body := map[string]any{
				"ruleType": opts.RuleType,
				"rules":    []any{rule},
			}
			if len(opts.RouterIDs) > 0 {
				routers := make([]map[string]any, len(opts.RouterIDs))
				for i, id := range opts.RouterIDs {
					routers[i] = map[string]any{"_id": id}
				}
				body["routers"] = routers
			}

			notify := map[string]any{}
			if len(opts.NotifyChannels) > 0 {
				notify["channels"] = opts.NotifyChannels
			}
			if len(opts.NotifyUsers) > 0 {
				notify["users"] = opts.NotifyUsers
			}
			if len(notify) > 0 {
				body["notify"] = notify
			}

			q := url.Values{}
			if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
				q.Set("oid", oid)
			}

			respBody, err := client.Do("POST", "/api/invpn/alerts/rules", q, body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteCreated(f, "Alert rule", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&opts.RuleType, "rule-type", "all", "Rule type (all or router)")
	cmd.Flags().StringSliceVar(&opts.RouterIDs, "router-ids", nil, "Router IDs (for rule-type=router)")
	cmd.Flags().StringVar(&opts.AlertType, "alert-type", "", "Alert type: vpn_connected, vpn_disconnected (required)")
	cmd.Flags().Int64Var(&opts.Retention, "retention", 0, "Retention period in seconds")
	cmd.Flags().StringSliceVar(&opts.NotifyChannels, "notify-channels", nil, "Notification channels (email, sms)")
	cmd.Flags().StringSliceVar(&opts.NotifyUsers, "notify-users", nil, "User IDs to notify")
	_ = cmd.MarkFlagRequired("alert-type")

	return cmd
}

func newCmdRuleUpdate(f *factory.Factory) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an alert rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if name != "" {
				body["name"] = name
			}

			respBody, err := client.Put("/api/invpn/alerts/rules/"+args[0], body)
			if err != nil {
				_ = iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
				return err
			}

			cmdutil.WriteUpdated(f, "Alert rule", respBody)
			return iostreams.FormatOutput(respBody, f.IO, f.IO.Output)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Rule name")
	return cmd
}

func newCmdRuleDelete(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an alert rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			_, err = client.Delete("/api/invpn/alerts/rules/" + args[0])
			if err != nil {
				return err
			}

			cmdutil.WriteDeleted(f, "Alert rule", args[0])
			return nil
		},
	}
}
