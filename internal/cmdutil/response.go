package cmdutil

import (
	"fmt"

	"github.com/inhandnet/ics-cli/internal/api"
	"github.com/inhandnet/ics-cli/internal/factory"
)

func WriteCreated(f *factory.Factory, resource string, body []byte) {
	id, name := api.ResultIDName(body)
	if name != "" {
		fmt.Fprintf(f.IO.ErrOut, "%s %q created (id: %s)\n", resource, name, id)
	} else if id != "" {
		fmt.Fprintf(f.IO.ErrOut, "%s created (id: %s)\n", resource, id)
	}
}

func WriteUpdated(f *factory.Factory, resource string, body []byte) {
	id, name := api.ResultIDName(body)
	if name != "" {
		fmt.Fprintf(f.IO.ErrOut, "%s %q updated (id: %s)\n", resource, name, id)
	} else if id != "" {
		fmt.Fprintf(f.IO.ErrOut, "%s updated (id: %s)\n", resource, id)
	}
}

func WriteDeleted(f *factory.Factory, resource, id string) {
	fmt.Fprintf(f.IO.ErrOut, "%s %s deleted\n", resource, id)
}
