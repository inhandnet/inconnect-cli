package iostreams

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/itchyny/gojq"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"gopkg.in/yaml.v3"
)

func FormatOutput(body []byte, ios *IOStreams, output string) error {
	data := unwrapResult(body)

	if ios.JQ != "" {
		return applyJQ(data, ios.JQ, ios.Out, ios.IsTTY)
	}

	switch output {
	case "table":
		return formatTable(data, ios.Out)
	case "yaml":
		return formatYAML(data, ios.Out)
	default:
		return formatJSON(data, ios.Out, ios.IsTTY)
	}
}

func formatJSON(data []byte, w io.Writer, isTTY bool) error {
	var buf []byte
	if isTTY {
		buf = pretty.Color(pretty.Pretty(data), nil)
	} else {
		var compact json.RawMessage = data
		b, err := json.Marshal(compact)
		if err != nil {
			b = data
		}
		buf = b
	}
	_, err := fmt.Fprintln(w, string(buf))
	return err
}

func formatTable(data []byte, w io.Writer) error {
	parsed := gjson.ParseBytes(data)

	if parsed.IsArray() {
		return formatArrayTable(parsed, w)
	}
	return formatObjectTable(parsed, w)
}

func formatArrayTable(arr gjson.Result, w io.Writer) error {
	items := arr.Array()
	if len(items) == 0 {
		fmt.Fprintln(w, "(no results)")
		return nil
	}

	var keys []string
	items[0].ForEach(func(key, _ gjson.Result) bool {
		keys = append(keys, key.String())
		return true
	})

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	header := make([]string, len(keys))
	for i, k := range keys {
		header[i] = strings.ToUpper(k)
	}
	fmt.Fprintln(tw, strings.Join(header, "\t"))

	for _, item := range items {
		vals := make([]string, len(keys))
		for i, k := range keys {
			vals[i] = cell(item.Get(k).String())
		}
		fmt.Fprintln(tw, strings.Join(vals, "\t"))
	}
	return tw.Flush()
}

func formatObjectTable(obj gjson.Result, w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	obj.ForEach(func(key, value gjson.Result) bool {
		fmt.Fprintf(tw, "%s\t%s\n", key.String(), cell(value.String()))
		return true
	})
	return tw.Flush()
}

const maxCellWidth = 48

// cell flattens a value for single-row table display: embedded tabs/newlines
// would corrupt column alignment, so collapse them to spaces, and overly long
// values (nested JSON, certs) are truncated so one wide column can't blow out
// the whole table. Use -o json/yaml or --jq for full values.
func cell(s string) string {
	s = strings.NewReplacer("\t", " ", "\n", " ", "\r", " ").Replace(s)
	if r := []rune(s); len(r) > maxCellWidth {
		return string(r[:maxCellWidth-1]) + "…"
	}
	return s
}

func formatYAML(data []byte, w io.Writer) error {
	var obj any
	if err := json.Unmarshal(data, &obj); err != nil {
		_, err = w.Write(data)
		return err
	}
	return yaml.NewEncoder(w).Encode(obj)
}

func applyJQ(data []byte, expr string, w io.Writer, isTTY bool) error {
	query, err := gojq.Parse(expr)
	if err != nil {
		return fmt.Errorf("invalid jq expression: %w", err)
	}

	var input any
	if err := json.Unmarshal(data, &input); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	iter := query.Run(input)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, isErr := v.(error); isErr {
			return err
		}
		out, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if isTTY {
			out = pretty.Color(pretty.Pretty(out), nil)
		}
		fmt.Fprintln(w, string(out))
	}
	return nil
}

func unwrapResult(body []byte) []byte {
	parsed := gjson.ParseBytes(body)
	if !parsed.IsObject() {
		return body
	}

	var keys []string
	parsed.ForEach(func(key, _ gjson.Result) bool {
		keys = append(keys, key.String())
		return true
	})

	if len(keys) == 1 && keys[0] == "result" {
		return []byte(parsed.Get("result").Raw)
	}

	if len(keys) == 2 {
		hasResult := false
		hasTotal := false
		for _, k := range keys {
			if k == "result" {
				hasResult = true
			}
			if k == "total" || k == "count" {
				hasTotal = true
			}
		}
		if hasResult && hasTotal {
			return []byte(parsed.Get("result").Raw)
		}
	}

	return body
}
