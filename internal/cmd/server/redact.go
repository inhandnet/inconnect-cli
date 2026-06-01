package server

import (
	"regexp"

	"github.com/spf13/cobra"
)

var pemPrivateKeyRE = regexp.MustCompile(`(?s)-----BEGIN [A-Z0-9 ]*PRIVATE KEY-----.*?-----END [A-Z0-9 ]*PRIVATE KEY-----`)

// redactBody masks PEM private keys in server payloads (e.g. keyPair.key) unless
// the caller passed --show-secrets. It operates on the raw JSON bytes so it works
// regardless of nesting or whether the response is a single object or a list.
func redactBody(cmd *cobra.Command, body []byte) []byte {
	if show, _ := cmd.Flags().GetBool("show-secrets"); show {
		return body
	}
	return pemPrivateKeyRE.ReplaceAll(body, []byte("***REDACTED (use --show-secrets to reveal)***"))
}
