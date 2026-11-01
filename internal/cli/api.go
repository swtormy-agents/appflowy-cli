package cli

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/swtormy-agents/appflowy-cli/internal/appflowy"
)

func (a *app) apiCmd() *cobra.Command {
	var method string
	var query []string
	var bodyFile string
	cmd := &cobra.Command{
		Use:   "api PATH",
		Short: "Raw API request (GET by default)",
		Long:  `PATH must start with /api/ and must not contain query or fragment. Use --query key=value. Non-GET methods require --write and TTY confirmation.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := appflowy.ValidateAPIPath(args[0])
			if err != nil {
				return usageErr("%s", err.Error())
			}
			m := strings.ToUpper(strings.TrimSpace(method))
			if m == "" {
				m = http.MethodGet
			}
			if m != http.MethodGet {
				if err := a.requireWrite(m + " " + path); err != nil {
					return err
				}
			}
			q, err := parseQueryFlags(query)
			if err != nil {
				return err
			}
			var body []byte
			if bodyFile != "" {
				if m == http.MethodGet || m == http.MethodHead {
					return usageErr("request body not allowed for %s", m)
				}
				body, err = readBodyFile(bodyFile, a.in)
				if err != nil {
					return usageErr("%s", err.Error())
				}
			}
			return a.do(context.Background(), appflowy.Request{
				Method: m,
				Path:   path,
				Query:  q,
				Body:   body,
			})
		},
	}
	cmd.Flags().StringVar(&method, "method", "GET", "HTTP method")
	cmd.Flags().StringArrayVar(&query, "query", nil, "query parameter key=value (repeatable)")
	cmd.Flags().StringVar(&bodyFile, "body", "", "JSON body file (- for stdin)")
	return cmd
}

func parseQueryFlags(pairs []string) (url.Values, error) {
	if len(pairs) == 0 {
		return nil, nil
	}
	q := url.Values{}
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok || k == "" {
			return nil, usageErr("invalid --query %q (want key=value)", p)
		}
		q.Add(k, v)
	}
	return q, nil
}

func readBodyFile(path string, in io.Reader) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(in)
	}
	return os.ReadFile(path)
}
