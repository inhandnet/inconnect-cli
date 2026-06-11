package iostreams

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/itchyny/gojq"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"gopkg.in/yaml.v3"
)

func FormatOutput(body []byte, ios *IOStreams, output string) error {
	if ios.JQ != "" {
		return applyJQ(unwrapResult(body), ios.JQ, ios.Out, ios.IsTTY)
	}

	switch output {
	case "table":
		// table renders from the original body so it can read pagination
		// metadata (total/page) from the envelope.
		return formatTable(body, ios)
	case "yaml":
		return formatYAML(unwrapResult(body), ios.Out)
	default:
		return formatJSON(unwrapResult(body), ios.Out, ios.IsTTY)
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

// formatTable renders JSON as a table. It unwraps a "result" envelope when
// present, then renders an array of objects as headered rows, or a single
// object as KEY/VALUE pairs. ios.Columns filters and orders the columns.
func formatTable(data []byte, ios *IOStreams) error {
	parsed := gjson.ParseBytes(data)

	items := parsed.Get("result")
	if !items.Exists() {
		items = parsed
	}

	tp := NewTablePrinter(ios.Out, ios.IsTTY)

	switch {
	case items.IsArray():
		arr := items.Array()
		if ios.IsTTY {
			if header := paginationHeader(&parsed, len(arr)); header != "" {
				fmt.Fprintln(ios.Out, Gray(header))
			}
		}
		if len(arr) == 0 {
			fmt.Fprintln(ios.Out, "(no results)")
			return nil
		}
		renderArray(tp, arr, ios)
	case items.IsObject():
		renderObject(tp, &items, ios)
	default:
		fmt.Fprintln(ios.Out, formatResult(&items))
		return nil
	}
	return tp.Render()
}

// paginationHeader builds a "Showing X of Y results (Page M of N)" line from
// the envelope metadata. Returns "" when no pagination info is available.
func paginationHeader(raw *gjson.Result, count int) string {
	if !raw.IsObject() || count == 0 {
		return ""
	}
	total := raw.Get("total")
	totalPages := raw.Get("totalPages")
	page := raw.Get("page")

	var parts []string
	if total.Exists() && total.Int() > 0 {
		parts = append(parts, fmt.Sprintf("Showing %d of %d results", count, total.Int()))
	} else {
		parts = append(parts, fmt.Sprintf("Showing %d results", count))
	}
	if totalPages.Exists() && totalPages.Int() > 0 {
		parts = append(parts, fmt.Sprintf("(Page %d of %d)", page.Int()+1, totalPages.Int()))
	} else if page.Exists() {
		parts = append(parts, fmt.Sprintf("(Page %d)", page.Int()+1))
	}
	return strings.Join(parts, " ")
}

// resolveColumns computes the effective column list. Entries prefixed with "!"
// are exclusions. Explicit includes (if any) define order; otherwise all keys
// are shown. Exclusions are removed from whichever set is used.
func resolveColumns(columns, allKeys []string) []string {
	if len(columns) == 0 {
		return allKeys
	}
	var includes []string
	excluded := map[string]bool{}
	for _, c := range columns {
		if strings.HasPrefix(c, "!") {
			excluded[strings.TrimPrefix(c, "!")] = true
		} else {
			includes = append(includes, c)
		}
	}

	base := includes
	if len(base) == 0 {
		base = allKeys
	}
	result := make([]string, 0, len(base))
	for _, c := range base {
		if !excluded[c] {
			result = append(result, c)
		}
	}
	return result
}

// renderArray renders an array of objects as a table with a header row.
func renderArray(tp *TablePrinter, items []gjson.Result, ios *IOStreams) {
	first := items[0]
	if !first.IsObject() {
		for i := range items {
			tp.AddRow(formatResult(&items[i]))
		}
		return
	}

	cols := resolveColumns(ios.Columns, flattenKeys(&first))

	header := make([]string, len(cols))
	for i, col := range cols {
		h := strings.ToUpper(col)
		if ios.IsTTY {
			h = Bold(h)
		}
		header[i] = h
	}
	tp.AddRow(header...)

	for _, item := range items {
		if !item.IsObject() {
			continue
		}
		row := make([]string, len(cols))
		for j, col := range cols {
			v := item.Get(col)
			row[j] = cell(formatResult(&v))
		}
		tp.AddRow(row...)
	}
}

// renderObject renders a single object as KEY / VALUE pairs.
func renderObject(tp *TablePrinter, obj *gjson.Result, ios *IOStreams) {
	cols := resolveColumns(ios.Columns, flattenKeys(obj))
	for _, key := range cols {
		val := obj.Get(key)
		if !val.Exists() {
			continue
		}
		k := key
		if ios.IsTTY {
			k = Bold(key)
		}
		tp.AddRow(k, formatResult(&val))
	}
}

// formatResult converts a gjson.Result to a display string, localizing
// ISO-8601 timestamps and rendering nested JSON inline.
func formatResult(r *gjson.Result) string {
	switch r.Type {
	case gjson.Null:
		return ""
	case gjson.String:
		if s := formatLocalTime(r.Str); s != "" {
			return s
		}
		return r.Str
	case gjson.True:
		return "true"
	case gjson.False:
		return "false"
	case gjson.Number:
		if r.Num == float64(int64(r.Num)) {
			return fmt.Sprintf("%d", int64(r.Num))
		}
		return formatFloat(r.Num)
	case gjson.JSON:
		return r.Raw
	default:
		return r.String()
	}
}

// cell collapses embedded whitespace so a value can't break table alignment.
func cell(s string) string {
	return strings.NewReplacer("\t", " ", "\n", " ", "\r", " ").Replace(s)
}

// escapeGjsonKey escapes dots and wildcards so gjson treats a key segment as a
// literal key rather than a nested path.
func escapeGjsonKey(s string) string {
	if !strings.ContainsAny(s, ".?*\\") {
		return s
	}
	var b strings.Builder
	for _, c := range s {
		switch c {
		case '.', '?', '*', '\\':
			b.WriteByte('\\')
		}
		b.WriteRune(c)
	}
	return b.String()
}

// flattenKeys collects leaf-level gjson-escaped dot-paths from an object,
// expanding nested objects (arrays and scalars are leaves), sorted.
func flattenKeys(r *gjson.Result) []string {
	keys := flattenKeysWithPrefix(r, "")
	sort.Strings(keys)
	return keys
}

func flattenKeysWithPrefix(r *gjson.Result, prefix string) []string {
	var keys []string
	r.ForEach(func(key, value gjson.Result) bool {
		seg := escapeGjsonKey(key.Str)
		path := seg
		if prefix != "" {
			path = prefix + "." + seg
		}
		if value.IsObject() {
			keys = append(keys, flattenKeysWithPrefix(&value, path)...)
		} else {
			keys = append(keys, path)
		}
		return true
	})
	return keys
}

var timestampLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
}

// formatLocalTime parses s as an ISO-8601 timestamp and returns it in local
// time, or "" if s is not a recognized timestamp.
func formatLocalTime(s string) string {
	if len(s) < 19 || s[4] != '-' || s[7] != '-' || s[10] != 'T' {
		return ""
	}
	for _, layout := range timestampLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Local().Format("2006-01-02 15:04:05")
		}
	}
	return ""
}

// formatFloat formats a float64 for table display with up to 3 decimals,
// trimming trailing zeros, falling back to %g for tiny values.
func formatFloat(f float64) string {
	s := strconv.FormatFloat(f, 'f', 3, 64)
	if f != 0 && strings.TrimRight(strings.TrimRight(s, "0"), ".") == "0" {
		return strconv.FormatFloat(f, 'g', -1, 64)
	}
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
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

// unwrapResult strips the envelope when the JSON object has "result" as its
// only key, or "result" plus a single pagination key (total/count). Richer
// envelopes are left intact so json/yaml output retains pagination metadata.
func unwrapResult(body []byte) []byte {
	parsed := gjson.ParseBytes(body)
	if !parsed.IsObject() {
		return body
	}

	result := parsed.Get("result")
	if !result.Exists() {
		return body
	}

	var keys []string
	parsed.ForEach(func(key, _ gjson.Result) bool {
		keys = append(keys, key.String())
		return true
	})

	if len(keys) == 1 {
		return []byte(result.Raw)
	}
	if len(keys) == 2 {
		for _, k := range keys {
			if k == "total" || k == "count" {
				return []byte(result.Raw)
			}
		}
	}
	return body
}
