package debug

import (
	"fmt"
	"os"
)

var Enabled bool

func Log(format string, args ...any) {
	if Enabled {
		fmt.Fprintf(os.Stderr, "[debug] "+format+"\n", args...)
	}
}
