package iostreams

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	gray   = "\033[90m"
)

func Green(s string) string  { return green + s + reset }
func Red(s string) string    { return red + s + reset }
func Yellow(s string) string { return yellow + s + reset }
func Bold(s string) string   { return bold + s + reset }
func Gray(s string) string   { return gray + s + reset }
