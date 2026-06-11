package iostreams

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/tidwall/gjson"
)

func TestResolveColumns(t *testing.T) {
	all := []string{"id", "name", "vip", "oid"}
	tests := []struct {
		name    string
		columns []string
		want    []string
	}{
		{"empty returns all", nil, all},
		{"includes define order", []string{"name", "id"}, []string{"name", "id"}},
		{"exclude removes from all", []string{"!oid", "!vip"}, []string{"id", "name"}},
		{"include plus exclude", []string{"name", "id", "!id"}, []string{"name"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveColumns(tt.columns, all)
			if strings.Join(got, ",") != strings.Join(tt.want, ",") {
				t.Errorf("resolveColumns(%v) = %v, want %v", tt.columns, got, tt.want)
			}
		})
	}
}

func TestFlattenKeys(t *testing.T) {
	r := gjson.Parse(`{"name":"a","metadata":{"net":"x","sig":{"rsrp":-65}},"tags":[1,2]}`)
	got := flattenKeys(&r)
	want := []string{"metadata.net", "metadata.sig.rsrp", "name", "tags"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("flattenKeys = %v, want %v", got, want)
	}
}

func TestFormatTableColumnsAndCJK(t *testing.T) {
	body := []byte(`{"result":[{"name":"默认网络","vip":"10.0.0.1","oid":"abc"},{"name":"test","vip":"10.0.0.2","oid":"def"}],"total":2}`)
	var buf bytes.Buffer
	ios := &IOStreams{Out: &buf, IsTTY: true, Columns: []string{"name", "vip"}}
	if err := formatTable(body, ios); err != nil {
		t.Fatal(err)
	}
	out := buf.String()

	if strings.Contains(out, "OID") {
		t.Errorf("excluded column OID appeared:\n%s", out)
	}
	if !strings.Contains(out, "NAME") || !strings.Contains(out, "VIP") {
		t.Errorf("expected headers NAME/VIP:\n%s", out)
	}
	if !strings.Contains(out, "默认网络") {
		t.Errorf("expected CJK value:\n%s", out)
	}

	// go-pretty must align the second column by DISPLAY width, not byte count:
	// the CJK row's prefix has more bytes but the same visual width as the ASCII
	// row, so the VIP values land at the same terminal column.
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	var vipDispCols []int
	for _, ln := range lines {
		if i := strings.Index(ln, "10.0.0."); i >= 0 {
			vipDispCols = append(vipDispCols, runewidth.StringWidth(ln[:i]))
		}
	}
	if len(vipDispCols) != 2 {
		t.Fatalf("expected 2 VIP rows, got %d:\n%s", len(vipDispCols), out)
	}
	if vipDispCols[0] != vipDispCols[1] {
		t.Errorf("VIP column misaligned by display width across CJK/ASCII rows: %d vs %d\n%s", vipDispCols[0], vipDispCols[1], out)
	}
}

func TestFormatTableEmpty(t *testing.T) {
	var buf bytes.Buffer
	ios := &IOStreams{Out: &buf, IsTTY: true}
	if err := formatTable([]byte(`{"result":[],"total":0}`), ios); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "(no results)") {
		t.Errorf("expected (no results), got: %q", buf.String())
	}
}

func TestFormatTableSingleObject(t *testing.T) {
	var buf bytes.Buffer
	ios := &IOStreams{Out: &buf, IsTTY: true, Columns: []string{"id", "name"}}
	if err := formatTable([]byte(`{"result":{"id":"1","name":"x","vip":"10.0.0.1"}}`), ios); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "id") || !strings.Contains(out, "name") {
		t.Errorf("expected KEY/VALUE rows for id/name:\n%s", out)
	}
	if strings.Contains(out, "10.0.0.1") {
		t.Errorf("vip should be filtered out:\n%s", out)
	}
}

func TestPaginationHeader(t *testing.T) {
	tests := []struct {
		name  string
		json  string
		count int
		want  string
	}{
		{"total and pages", `{"total":50,"totalPages":5,"page":1}`, 10, "Showing 10 of 50 results (Page 2 of 5)"},
		{"total only", `{"total":3}`, 3, "Showing 3 of 3 results"},
		{"no total", `{}`, 4, "Showing 4 results"},
		{"page without totalPages", `{"page":2}`, 5, "Showing 5 results (Page 3)"},
		{"zero count suppresses header", `{"total":10}`, 0, ""},
		{"non-object", `[1,2,3]`, 3, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gjson.Parse(tt.json)
			if got := paginationHeader(&r, tt.count); got != tt.want {
				t.Errorf("paginationHeader(%s, %d) = %q, want %q", tt.json, tt.count, got, tt.want)
			}
		})
	}
}

func TestFormatLocalTime(t *testing.T) {
	// Recognized timestamps: compare against an independent parse + local
	// conversion so the test is timezone-agnostic. A non-empty want also proves
	// the layout was actually matched.
	for _, in := range []string{
		"2026-06-11T12:30:45Z",
		"2026-06-11T12:30:45.123456789Z",
		"2026-06-11T12:30:45",
		"2026-06-11T12:30:45+08:00",
	} {
		var want string
		for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05"} {
			if tm, err := time.Parse(layout, in); err == nil {
				want = tm.Local().Format("2006-01-02 15:04:05")
				break
			}
		}
		if want == "" {
			t.Fatalf("test setup: %q did not match any layout", in)
		}
		if got := formatLocalTime(in); got != want {
			t.Errorf("formatLocalTime(%q) = %q, want %q", in, got, want)
		}
	}

	// Non-timestamps must return "".
	for _, in := range []string{"", "hello", "10.0.0.1", "2026-06-11", "2026/06/11T12:30:45", "abcdefghijklmnopqrs"} {
		if got := formatLocalTime(in); got != "" {
			t.Errorf("formatLocalTime(%q) = %q, want empty", in, got)
		}
	}
}

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		in   float64
		want string
	}{
		{1.5, "1.5"},
		{1.25, "1.25"},
		{3.14159, "3.142"}, // rounded to 3 decimals
		{2.0, "2"},
		{0, "0"},
		{100, "100"},
		{0.0001, "0.0001"}, // too small for 3 decimals; falls back to %g
		{0.001, "0.001"},
	}
	for _, tt := range tests {
		if got := formatFloat(tt.in); got != tt.want {
			t.Errorf("formatFloat(%v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
