package cli

import (
	"context"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/swtormy-agents/appflowy-cli/internal/appflowy"
)

func (a *app) pageCmd() *cobra.Command {
	page := &cobra.Command{
		Use:   "page",
		Short: "Read page content",
	}
	var collabType string
	get := &cobra.Command{
		Use:   "get WORKSPACE_ID PAGE_ID",
		Short: "Get page document JSON (GET /api/workspace/v1/{ws}/collab/{id}/json)",
		Long:  `Fetches structured document JSON for a page collab object. PAGE_ID is the collab/view UUID.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			wsID, err := a.parseUUIDArg(args[0])
			if err != nil {
				return err
			}
			pageID, err := a.parseUUIDArg(args[1])
			if err != nil {
				return err
			}
			q := url.Values{}
			q.Set("collab_type", collabType)
			return a.do(context.Background(), appflowy.Request{
				Method: http.MethodGet,
				Path:   "/api/workspace/v1/" + wsID + "/collab/" + pageID + "/json",
				Query:  q,
			})
		},
	}
	get.Flags().StringVar(&collabType, "collab-type", "document", "collab_type query parameter (default document)")
	page.AddCommand(get)
	return page
}
