# Changelog

All notable changes to this repository are documented here.
Skills follow semver in their `metadata.version`; format adapted from [Keep a Changelog](https://keepachangelog.com).

## 2026-11-13

### Changed

- **clean-code 1.3.1** — Rewritten to be language-agnostic (removed the Rust, TypeScript, and MoonBit appendices plus the language-bound eval cases) and translated to English. Supersedes 1.2.x.

### Added

- **create-skill 1.0.0** — Skill-authoring workflow: intent capture, SKILL.md drafting with progressive disclosure, realistic test prompts, and feedback-driven iteration. Distilled from well-known skill-creator implementations.
- **docs-engineering 1.0.0** — Diátaxis-based documentation methodology for deciding where content belongs, per-quadrant writing rules, and docs-site information architecture. Extracted and generalized from an internal project.
- Repository engineering: skill structure validator (`scripts/validate_skills.py`), GitHub Actions workflow, release checklist plus evals and versioning conventions in CONTRIBUTING, and this changelog.

### Removed

- Stale machine-local `skills-lock.json` (now gitignored — the lock belongs to consumers, not to this distribution repo).
