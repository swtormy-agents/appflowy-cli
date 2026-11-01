package cli

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed SKILL.md
var skillMD string

var skills = map[string]*string{"appflowy": &skillMD}

func skillNames() []string { return []string{"appflowy"} }

var skillDirs = map[string]string{
	"claude": ".claude/skills",
	"cursor": ".cursor/skills",
	"grok":   ".grok/skills",
	"agents": ".agents/skills",
}

func (a *app) skillCmd() *cobra.Command {
	skill := &cobra.Command{
		Use:   "skill",
		Short: "Print or install the agent skill bundled in this binary",
	}
	skill.AddCommand(&cobra.Command{
		Use:   "print [appflowy]",
		Short: "Print SKILL.md to stdout (default: appflowy)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, _, err := pickSkill(args)
			if err != nil {
				return err
			}
			_, err = fmt.Fprint(a.out, body)
			return err
		},
	})

	var agent, dir string
	var global bool
	install := &cobra.Command{
		Use:   "install [appflowy]",
		Short: "Write SKILL.md for an agent (project-local by default)",
		Example: `  appflowy skill install
  appflowy skill install --global
  appflowy skill install --agent cursor --global`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, name, err := pickSkill(args)
			if err != nil {
				return err
			}
			target := dir
			if target == "" {
				sub, ok := skillDirs[agent]
				if !ok {
					return usageErr("--agent must be one of claude, cursor, grok, agents")
				}
				base, err := os.Getwd()
				if global {
					base, err = os.UserHomeDir()
				}
				if err != nil {
					return runtimeErr(err)
				}
				target = filepath.Join(base, sub, name)
			}
			if err := os.MkdirAll(target, 0o755); err != nil {
				return runtimeErr(err)
			}
			file := filepath.Join(target, "SKILL.md")
			if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
				return runtimeErr(err)
			}
			return a.emitValue(map[string]any{"ok": true, "path": file})
		},
	}
	install.Flags().StringVar(&agent, "agent", "cursor", "agent layout: claude, cursor, grok, agents")
	install.Flags().StringVar(&dir, "dir", "", "override install directory")
	install.Flags().BoolVar(&global, "global", false, "install under user home instead of cwd")
	skill.AddCommand(install)
	return skill
}

func pickSkill(args []string) (body string, name string, err error) {
	name = "appflowy"
	if len(args) > 0 {
		name = strings.ToLower(args[0])
	}
	ptr, ok := skills[name]
	if !ok {
		return "", "", usageErr("unknown skill %q (available: %s)", name, strings.Join(skillNames(), ", "))
	}
	return *ptr, name, nil
}
