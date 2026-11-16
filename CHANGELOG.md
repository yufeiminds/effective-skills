# Changelog

All notable changes to this repository are documented here.
Skills follow semver in their `metadata.version`; format adapted from [Keep a Changelog](https://keepachangelog.com).

## 2026-11-16

### Added

- **mise-managed toolchains** — `mise.toml` pins `go` 1.27.4 (mirrors `cli/go.mod`), `node` 24, and `python` 3.13; CI installs them via `jdx/mise-action` so local and CI versions stay identical.
- **`cli/` skill validator (Go, Charm Fang)** — replaces the Python script. Spec rules aligned with the official `skills-ref` reference implementation (field whitelist, name ≤64 chars/lowercase/hyphen rules, description ≤1024, compatibility ≤500) plus repository conventions (README listing) and lint warnings (SKILL.md body over 500 lines, broken relative links, descriptions without trigger guidance). Ships `validate` (human/JSON output, `--strict`) and `read-properties` subcommands.
- **`.claude-plugin/marketplace.json`** — the skill collection can now be installed as a Claude Code plugin in addition to the `skills` CLI.
- **evals scaffolds** — `skills/clean-code/evals/evals.json` and `skills/diataxis/evals/evals.json` with qualitative test prompts (both skills have subjective outputs; no assertions).

### Changed

- **Single source of truth for engineering rules** — new `AGENTS.md` carries the invariants agents must follow; README rewritten with per-agent collapsible install blocks, verbatim-copied install commands, and a mattpocock-style Reference section; CONTRIBUTING slimmed to the contribution workflow.
- **Version-drift handling in the validator** — `.claude-plugin/marketplace.json` no longer carries a version field, and the validator rejects any version key in it; per-skill versions live in each SKILL.md frontmatter (metadata.version).
- **diataxis 2.0.0** (renamed from `docs-engineering`) — the skill is named after the framework it teaches. Description unchanged, so task-keyword triggering is unaffected. **Breaking** for consumers installing by the old name (`--skill docs-engineering`).
- Dependencies rolled to latest stable: Go toolchain 1.27.4, module deps (`fang` v1.0.0, `charm.land/lipgloss/v2` v2.1.1, `x/sys` v0.50.0, …); `govulncheck` reports zero known vulnerabilities. Upstream `skill-creator` verified current (1780c5e).
- CI runs `go test ./cli/...` and the Go validator instead of the Python script; toolchain setup (setup-go/setup-node) replaced by `jdx/mise-action`.
- CONTRIBUTING: frontmatter table now covers `compatibility` and `allowed-tools`; new Official Resources section (agentskills.io spec and authoring guides); evaluation guidance split by output type (objective/subjective/triggering); validation commands updated to the Go CLI; documented that `skills-lock.json` pointing at the local mirror is expected.
- **clean-code 1.3.2** — References section now states when to read each companion doc; the stray "Deep Module" note moved out of the link list.

### Removed

- `scripts/validate_skills.py` (superseded by `cli/`).

## 2026-11-14

### Removed

- **create-skill** — retired in favor of Anthropic's upstream `skill-creator`, which now covers skill authoring and evaluation as a development dependency.
- Vendored `tooling/skill-creator/` copy (superseded by the restored dependency).

### Changed

- Development tooling: Anthropic's `skill-creator` (Apache-2.0) is installed again with the `skills` CLI at `.agents/skills/skill-creator/` and tracked in `skills-lock.json`, which is no longer gitignored. CONTRIBUTING documents bootstrapping the upstream mirror, syncing it, and running `npx skills update skill-creator`.
- CONTRIBUTING: companion docs now use `references/` (plus `scripts/` and `assets/`), and `npx skills init` is run from `skills/`.

## 2026-11-13

### Changed

- **clean-code 1.3.1** — Rewritten to be language-agnostic (removed the Rust, TypeScript, and MoonBit appendices plus the language-bound eval cases) and translated to English. Supersedes 1.2.x.

### Added

- **create-skill 1.0.0** — Skill-authoring workflow: intent capture, SKILL.md drafting with progressive disclosure, realistic test prompts, and feedback-driven iteration. Distilled from well-known skill-creator implementations.
- **docs-engineering 1.0.0** — Diátaxis-based documentation methodology for deciding where content belongs, per-quadrant writing rules, and docs-site information architecture. Extracted and generalized from an internal project.
- Repository engineering: skill structure validator (`scripts/validate_skills.py`), GitHub Actions workflow, release checklist plus evals and versioning conventions in CONTRIBUTING, and this changelog.

### Removed

- Stale machine-local `skills-lock.json` (now gitignored — the lock belongs to consumers, not to this distribution repo).
