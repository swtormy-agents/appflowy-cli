package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/swtormy-agents/appflowy-cli/internal/appflowy"
)

const (
	exitRuntime = 1
	exitUsage   = 2
	exitAuth    = 3
	exitPolicy  = 4
)

type exitError struct {
	code int
	kind string
	err  error
	api  *appflowy.APIError
}

func (e *exitError) Error() string { return e.err.Error() }

func usageErr(format string, a ...any) error {
	return &exitError{code: exitUsage, kind: "usage", err: fmt.Errorf(format, a...)}
}

func runtimeErr(err error) error {
	return &exitError{code: exitRuntime, kind: "runtime", err: err}
}

func authErr(err error) error {
	return &exitError{code: exitAuth, kind: "auth", err: err}
}

func policyErr(msg string) error {
	return &exitError{code: exitPolicy, kind: "policy", err: errors.New(msg)}
}

var Version = "dev"
var Commit = ""

func versionString() string {
	if Version != "dev" {
		out := Version
		if Commit != "" {
			c := Commit
			if len(c) > 7 {
				c = c[:7]
			}
			out += "+" + c
		}
		return out
	}
	if bi, ok := debug.ReadBuildInfo(); ok {
		rev, dirty := "", ""
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				rev = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					dirty = "+dirty"
				}
			}
		}
		if len(rev) > 7 {
			rev = rev[:7]
		}
		if rev != "" {
			return "dev-" + rev + dirty
		}
	}
	return Version
}

type app struct {
	in     io.Reader
	out    io.Writer
	errOut io.Writer
	getenv func(string) string
	isTTY  func() bool

	cfg         appflowy.Config
	dryRun      bool
	writeFlag   bool
	writeConfirm string
	jq          string
	pretty      bool
	jsonOut     bool
	baseURL     string
	timeout     time.Duration
}

func Execute(args []string, stdin io.Reader, stdout, stderr io.Writer, getenv func(string) string) int {
	if getenv == nil {
		getenv = os.Getenv
	}
	a := &app{
		in: stdin, out: stdout, errOut: stderr, getenv: getenv,
		isTTY: defaultIsTTY, timeout: 60 * time.Second, jsonOut: true,
	}
	root := a.rootCmd()
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)
	if err := root.Execute(); err != nil {
		return a.reportError(err)
	}
	return 0
}

func (a *app) reportError(err error) int {
	ee, ok := err.(*exitError)
	if !ok {
		if u, ok := mapValidationErr(err); ok {
			ee = u
		} else {
			ee = &exitError{code: exitUsage, kind: "usage", err: err}
		}
	}
	detail := map[string]any{"type": ee.kind, "message": a.redact(ee.err.Error())}
	if ee.api != nil {
		detail["status"] = ee.api.Status
		if ee.api.Message != "" {
			detail["message"] = a.redact(ee.api.Message)
		}
	}
	b, _ := marshal(struct {
		OK    bool           `json:"ok"`
		Error map[string]any `json:"error"`
	}{false, detail})
	fmt.Fprintln(a.errOut, string(b))
	return ee.code
}

func mapValidationErr(err error) (*exitError, bool) {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "UUID"):
		return &exitError{code: exitUsage, kind: "usage", err: err}, true
	case strings.Contains(msg, "path"):
		return &exitError{code: exitUsage, kind: "usage", err: err}, true
	case strings.Contains(msg, "depth"):
		return &exitError{code: exitUsage, kind: "usage", err: err}, true
	case strings.Contains(msg, "APPFLOWY"):
		return &exitError{code: exitUsage, kind: "usage", err: err}, true
	default:
		return nil, false
	}
}

func (a *app) redact(s string) string {
	return (&appflowy.Client{Token: a.cfg.Token}).Redact(s)
}

func (a *app) loadConfig() error {
	cfg, err := appflowy.LoadConfig(a.getenv)
	if err != nil {
		return usageErr("%s", err.Error())
	}
	if a.baseURL != "" {
		n, insecure, err := appflowy.ValidateBaseURL(a.baseURL)
		if err != nil {
			return usageErr("%s", err.Error())
		}
		cfg.URL = n
		if insecure {
			fmt.Fprintln(a.errOut, "appflowy: warning: APPFLOWY_URL uses http://; prefer https://")
		}
	}
	a.cfg = cfg
	return nil
}

func (a *app) client() (*appflowy.Client, error) {
	if err := a.loadConfig(); err != nil {
		return nil, err
	}
	if a.cfg.URL == "" && !a.dryRun {
		return nil, usageErr("APPFLOWY_URL is not set")
	}
	if a.cfg.Token == "" && !a.dryRun {
		return nil, authErr(errors.New("APPFLOWY_TOKEN is not set"))
	}
	return &appflowy.Client{
		Token:   a.cfg.Token,
		BaseURL: a.cfg.URL,
		HTTP:    &http.Client{Timeout: a.timeout},
	}, nil
}

func (a *app) do(ctx context.Context, r appflowy.Request) error {
	if a.dryRun {
		if err := a.loadConfig(); err != nil {
			return err
		}
		if a.cfg.URL == "" {
			return usageErr("APPFLOWY_URL is not set")
		}
		return a.emitValue(appflowy.DryRun(a.cfg.URL, r))
	}
	c, err := a.client()
	if err != nil {
		return err
	}
	raw, err := c.Do(ctx, r)
	if err != nil {
		return a.wrapErr(err)
	}
	return a.emit(raw)
}

func (a *app) wrapErr(err error) error {
	var api *appflowy.APIError
	if errors.As(err, &api) {
		code, kind := exitRuntime, "api"
		if api.IsAuth() {
			code, kind = exitAuth, "auth"
		}
		return &exitError{code: code, kind: kind, err: api, api: api}
	}
	return runtimeErr(err)
}

func (a *app) requireWrite(action string) error {
	if !a.writeFlag {
		return policyErr(fmt.Sprintf("%s requires --write (read-only by default; no env override)", action))
	}
	if !a.isTTY() {
		return policyErr(fmt.Sprintf("%s requires an interactive terminal to confirm write", action))
	}
	if a.writeConfirm == "write" {
		return nil
	}
	fmt.Fprint(a.errOut, "Type \"write\" to confirm: ")
	var line string
	if _, err := fmt.Fscanln(a.in, &line); err != nil {
		return policyErr("write not confirmed")
	}
	if strings.TrimSpace(line) != "write" {
		return policyErr("write not confirmed")
	}
	return nil
}

func (a *app) parseUUIDArg(s string) (string, error) {
	id, err := appflowy.ParseUUID(s)
	if err != nil {
		return "", usageErr("%s", err.Error())
	}
	return id, nil
}

func (a *app) rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "appflowy",
		Short:         "AppFlowy Cloud REST API from the shell (unofficial)",
		Long:          rootLong,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	root.CompletionOptions.DisableDefaultCmd = true
	pf := root.PersistentFlags()
	pf.BoolVar(&a.dryRun, "dry-run", false, "print the HTTP request instead of sending it")
	pf.BoolVar(&a.writeFlag, "write", false, "opt in to a non-GET request on this invocation (requires TTY confirmation)")
	pf.StringVar(&a.writeConfirm, "confirm", "", `when using --write, may be the word "write" to skip the interactive prompt`)
	pf.BoolVar(&a.pretty, "pretty", false, "indent JSON output")
	pf.BoolVar(&a.jsonOut, "json", true, "emit JSON on stdout (default true)")
	pf.StringVar(&a.jq, "jq", "", "filter JSON output with a jq expression")
	pf.StringVar(&a.baseURL, "url", "", "AppFlowy Cloud base URL (default $APPFLOWY_URL)")
	pf.DurationVar(&a.timeout, "http-timeout", 60*time.Second, "HTTP timeout per request")

	root.AddCommand(
		a.versionCmd(),
		a.statusCmd(),
		a.workspacesCmd(),
		a.foldersCmd(),
		a.pageCmd(),
		a.searchCmd(),
		a.apiCmd(),
		a.skillCmd(),
	)
	return root
}

const rootLong = `appflowy calls the AppFlowy Cloud REST API. JSON on stdout; errors as JSON on stderr.
This CLI is unofficial and not affiliated with or endorsed by AppFlowy.

Auth
  export APPFLOWY_URL=https://your-appflowy-cloud.example
  export APPFLOWY_TOKEN=...              Authorization: Bearer (never printed; masked as <TOKEN>)

Safety
  Non-GET requests require --write and interactive confirmation (type "write").
  There is no environment variable or config file to enable writes.

Exit codes: 0 ok, 1 API/runtime, 2 usage, 3 auth, 4 blocked write policy`

func defaultIsTTY() bool {
	f, err := os.Stderr.Stat()
	if err != nil {
		return false
	}
	return (f.Mode() & os.ModeCharDevice) != 0
}
