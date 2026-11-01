package cli

import (
	"context"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/swtormy-agents/appflowy-cli/internal/appflowy"
)

func (a *app) workspacesCmd() *cobra.Command {
	ws := &cobra.Command{
		Use:   "workspaces",
		Short: "Workspace operations",
	}
	ws.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List workspaces (GET /api/workspace)",
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if includeMembers, _ := cmd.Flags().GetBool("include-member-count"); includeMembers {
				q.Set("include_member_count", "true")
			}
			return a.do(context.Background(), appflowy.Request{
				Method: http.MethodGet,
				Path:   "/api/workspace",
				Query:  q,
			})
		},
	})
	list := ws.Commands()[0]
	list.Flags().Bool("include-member-count", false, "include member counts in the response")
	return ws
}
