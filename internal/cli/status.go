package cli

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/swtormy-agents/appflowy-cli/internal/appflowy"
)

func (a *app) versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version and commit as JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.emitValue(map[string]string{
				"version": versionString(),
				"commit":  Commit,
			})
		},
	}
}

func (a *app) statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show configuration (never prints the token)",
		Long:  `Reports APPFLOWY_URL and whether APPFLOWY_TOKEN is set. When both are set, calls GET /api/workspace to verify auth.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.loadConfig(); err != nil {
				return err
			}
			out := map[string]any{
				"ok":        true,
				"url":       a.cfg.URL,
				"has_token": a.cfg.Token != "",
			}
			if a.cfg.Token == "" || a.cfg.URL == "" {
				return a.emitValue(out)
			}
			c, err := a.client()
			if err != nil {
				return err
			}
			raw, err := c.Do(context.Background(), appflowy.Request{
				Method: http.MethodGet,
				Path:   "/api/workspace",
			})
			if err != nil {
				return a.wrapErr(err)
			}
			out["workspaces"] = jsonRawToAny(raw)
			return a.emitValue(out)
		},
	}
}

func jsonRawToAny(raw []byte) any {
	var v any
	_ = json.Unmarshal(raw, &v)
	return v
}
