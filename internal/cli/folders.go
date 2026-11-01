package cli

import (
	"context"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/swtormy-agents/appflowy-cli/internal/appflowy"
)

func (a *app) foldersCmd() *cobra.Command {
	folders := &cobra.Command{
		Use:     "folders",
		Short:   "Folder tree for a workspace",
		Aliases: []string{"views"},
	}
	folders.AddCommand(a.folderListCmd())
	return folders
}

func (a *app) folderListCmd() *cobra.Command {
	var depth string
	var rootViewID string
	cmd := &cobra.Command{
		Use:   "list WORKSPACE_ID",
		Short: "List folder/view tree (GET /api/workspace/{id}/folder)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			wsID, err := a.parseUUIDArg(args[0])
			if err != nil {
				return err
			}
			q := url.Values{}
			if depth != "" {
				d, err := appflowy.ParseDepth(depth)
				if err != nil {
					return usageErr("%s", err.Error())
				}
				q.Set("depth", d)
			}
			if rootViewID != "" {
				rid, err := a.parseUUIDArg(rootViewID)
				if err != nil {
					return err
				}
				q.Set("root_view_id", rid)
			}
			return a.do(context.Background(), appflowy.Request{
				Method: http.MethodGet,
				Path:   "/api/workspace/" + wsID + "/folder",
				Query:  q,
			})
		},
	}
	cmd.Flags().StringVar(&depth, "depth", "", "tree depth (1-10)")
	cmd.Flags().StringVar(&rootViewID, "root-view-id", "", "root view UUID")
	return cmd
}
