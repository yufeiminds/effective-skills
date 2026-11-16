package main

import (
	"os"

	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	opts := &validateCmdOptions{}
	cmd := &cobra.Command{
		Use:   "validate [paths...]",
		Short: "Validate skills; exits non-zero on failure",
		Long: `Validates skill directories against the agentskills.io specification
and this repository's conventions.

With no arguments, every skill under the top-level skills/ directory is
validated. Each path argument may be a repository root, a skills directory,
or a single skill directory.

Exit status: 0 when no errors are found (warnings alone pass), 1 when any
error is found. With --strict, warnings also cause exit status 1.`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dirs, err := discoverSkills(args)
			if err != nil {
				return err
			}
			results := make([]*skillResult, 0, len(dirs))
			for _, dir := range dirs {
				results = append(results, validateSkillDir(dir))
			}
			// Repository-level check: the marketplace manifest must not
			// carry a version field (drifts silently as skills change).
			// Skipped when validating a standalone tree without one.
			if mp, ok := findUpFile(".", ".claude-plugin/marketplace.json", 6); ok {
				res := &skillResult{
					Path:     displayPath(mp),
					Errors:   checkMarketplaceJSON(mp),
					Warnings: []string{},
				}
				results = append(results, res)
			}
			rep := buildReport(results)

			if opts.json {
				if err := renderJSON(os.Stdout, rep); err != nil {
					return err
				}
			} else {
				renderHuman(os.Stdout, rep)
			}

			if rep.failed(opts.strict) {
				cmd.SilenceErrors = true
				return errValidationFailed
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&opts.json, "json", false, "print a machine-readable JSON report")
	cmd.Flags().BoolVar(&opts.strict, "strict", false, "treat warnings as failures")
	return cmd
}

// validateCmdOptions holds the validate subcommand's flags.
type validateCmdOptions struct {
	json   bool
	strict bool
}
