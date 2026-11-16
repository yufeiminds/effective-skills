// Command skill-validator validates Agent Skills (SKILL.md packages) against
// the agentskills.io specification and this repository's conventions.
//
// It re-implements the validation rules of the official skills-ref reference
// implementation (github.com/agentskills/agentskills, Apache-2.0) in Go, and
// adds repository-specific checks (README listing, metadata conventions) plus
// lint warnings (body length, broken links, trigger guidance).
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

const version = "0.1.0"

// errValidationFailed signals that a validation run reported at least one
// error (or a warning under --strict). The report has already been printed by
// the validate command, so the fang error handler silences it; the error
// exists only to turn into a non-zero exit code.
var errValidationFailed = errors.New("validation failed")

func main() {
	root := &cobra.Command{
		Use:   "skill-validator",
		Short: "Validate Agent Skills against the agentskills.io specification",
		Long: `skill-validator checks skills for spec compliance and repository conventions.

Rules are aligned with the official skills-ref reference implementation
(github.com/agentskills/agentskills). Repository-specific checks and lint
warnings are layered on top:

  spec errors       SKILL.md presence, frontmatter format, field whitelist,
                    name/description/compatibility constraints
  repo conventions  README Available Skills listing; marketplace.json must
                    not carry a version field (skill versions live in each
                    SKILL.md frontmatter)
  lint warnings     SKILL.md body length, broken relative links,
                    missing trigger guidance in the description

Run "skill-validator validate" from the repository root to check every skill
under skills/. Engineering rules live in AGENTS.md.`,
		SilenceUsage: true,
	}

	root.AddCommand(newValidateCmd(), newReadPropertiesCmd())

	err := fang.Execute(context.Background(), root,
		fang.WithVersion(version),
		fang.WithErrorHandler(silenceValidationErrors))
	if err != nil {
		os.Exit(1)
	}
}

// silenceValidationErrors hides errValidationFailed from fang's styled error
// output because the validate command has already printed a full report.
func silenceValidationErrors(w io.Writer, styles fang.Styles, err error) {
	if errors.Is(err, errValidationFailed) {
		return
	}
	fang.DefaultErrorHandler(w, styles, err)
}

// failf returns a user-facing error for fang to render.
func failf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
