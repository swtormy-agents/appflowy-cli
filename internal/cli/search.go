package cli

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/swtormy-agents/appflowy-cli/internal/appflowy"
)

func (a *app) searchCmd() *cobra.Command {
	search := &cobra.Command{
		Use:   "search",
		Short: "Search within a workspace (quick notes API)",
		Long:  `Uses GET /api/workspace/{id}/quick-note?search_term=... when the server exposes quick-note search.`,
	}
	var limit, offset int
	search.AddCommand(&cobra.Command{
		Use:   "quick-notes WORKSPACE_ID QUERY",
		Short: "Search quick notes by term",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			wsID, err := a.parseUUIDArg(args[0])
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("search_term", args[1])
			if limit > 0 {
				q.Set("limit", strconv.Itoa(limit))
			}
			if offset > 0 {
				q.Set("offset", strconv.Itoa(offset))
			}
			return a.do(context.Background(), appflowy.Request{
				Method: http.MethodGet,
				Path:   "/api/workspace/" + wsID + "/quick-note",
				Query:  q,
			})
		},
	})
	c := search.Commands()[0]
	c.Flags().IntVar(&limit, "limit", 0, "max results (server default if 0)")
	c.Flags().IntVar(&offset, "offset", 0, "pagination offset")
	return search
}
