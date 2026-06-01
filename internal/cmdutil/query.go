package cmdutil

import (
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
)

type ListFlags struct {
	Cursor int
	Limit  int
	Sort   string
	Name   string
}

func (lf *ListFlags) RegisterPagination(cmd *cobra.Command) {
	cmd.Flags().IntVar(&lf.Cursor, "cursor", 0, "Pagination cursor / offset")
	cmd.Flags().IntVar(&lf.Limit, "limit", 20, "Maximum number of records to return")
	cmd.Flags().StringVar(&lf.Sort, "sort", "", "Sort order (e.g. createdAt,desc)")

	cmd.Flags().IntVar(&lf.Limit, "page-size", 20, "Alias for --limit")
	cmd.Flags().IntVar(&lf.Limit, "per-page", 20, "Alias for --limit")
	_ = cmd.Flags().MarkHidden("page-size")
	_ = cmd.Flags().MarkHidden("per-page")
}

func (lf *ListFlags) Register(cmd *cobra.Command) {
	lf.RegisterPagination(cmd)
	cmd.Flags().StringVar(&lf.Name, "name", "", "Filter by name")
}

func NewQuery(cmd *cobra.Command, lf *ListFlags) url.Values {
	return buildQuery(cmd, lf, "skip")
}

func NewCursorQuery(cmd *cobra.Command, lf *ListFlags) url.Values {
	return buildQuery(cmd, lf, "cursor")
}

func buildQuery(cmd *cobra.Command, lf *ListFlags, cursorParam string) url.Values {
	q := url.Values{}

	if lf.Cursor > 0 {
		q.Set(cursorParam, strconv.Itoa(lf.Cursor))
	}
	if lf.Limit > 0 {
		q.Set("limit", strconv.Itoa(lf.Limit))
	}
	if lf.Sort != "" {
		q.Set("sort", lf.Sort)
	}
	if lf.Name != "" {
		q.Set("name", lf.Name)
	}

	if oid, _ := cmd.Flags().GetString("oid"); oid != "" {
		q.Set("oid", oid)
	}

	return q
}
