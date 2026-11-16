package main

import (
	"encoding/json"
	"fmt"
	"io"

	"charm.land/lipgloss/v2"
)

// report aggregates per-skill results for rendering.
type report struct {
	Skills  []*skillResult `json:"skills"`
	Summary summary        `json:"summary"`
}

type summary struct {
	Checked  int `json:"checked"`
	Failed   int `json:"failed"`
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
}

func buildReport(results []*skillResult) report {
	rep := report{Skills: results}
	for _, r := range results {
		rep.Summary.Checked++
		rep.Summary.Errors += len(r.Errors)
		rep.Summary.Warnings += len(r.Warnings)
		if !r.ok() {
			rep.Summary.Failed++
		}
	}
	return rep
}

// failed reports whether the run should exit non-zero: any error always
// fails; warnings fail only under --strict.
func (rep report) failed(strict bool) bool {
	if rep.Summary.Failed > 0 {
		return true
	}
	return strict && rep.Summary.Warnings > 0
}

var (
	styleOK    = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	styleErr   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // red
	styleWarn  = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // yellow
	styleBold  = lipgloss.NewStyle().Bold(true)
	markerOK   = styleOK.Render("✔")
	markerWarn = styleWarn.Render("⚠")
	markerErr  = styleErr.Render("✗")
	tagError   = styleErr.Render("error:")
	tagWarning = styleWarn.Render("warning:")
)

// renderHuman prints the human-readable validation report. Warnings alone
// render as ⚠ (pass), errors as ✗ (fail).
func renderHuman(w io.Writer, rep report) {
	width := 0
	for _, r := range rep.Skills {
		if len(r.Path) > width {
			width = len(r.Path)
		}
	}
	for _, r := range rep.Skills {
		label := styleBold.Render(fmt.Sprintf("%-*s", width, r.Path))
		switch {
		case !r.ok():
			fmt.Fprintf(w, "%s %s\n", markerErr, label)
		case len(r.Warnings) > 0:
			fmt.Fprintf(w, "%s %s\n", markerWarn, label)
		default:
			fmt.Fprintf(w, "%s %s\n", markerOK, label)
		}
		for _, e := range r.Errors {
			fmt.Fprintf(w, "  %s %s\n", tagError, e)
		}
		for _, warn := range r.Warnings {
			fmt.Fprintf(w, "  %s %s\n", tagWarning, warn)
		}
	}
	s := rep.Summary
	fmt.Fprintln(w)
	fmt.Fprintf(w, "%d skill(s) checked, %d failed (%d errors, %d warnings)\n",
		s.Checked, s.Failed, s.Errors, s.Warnings)
}

// renderJSON prints the machine-readable report used by CI and automation.
func renderJSON(w io.Writer, rep report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rep)
}
