# Contributing to effective-skills

Thank you for your interest! This repository hosts agent skills — reusable instruction sets that extend AI coding agents' capabilities.

Engineering rules — how skills are structured, named, versioned, validated, and released — live in [AGENTS.md](AGENTS.md) and are enforced mechanically by the validator in `cli/`. Contributors and coding agents alike should read it before proposing changes; it lists every file a skill change must touch.

## Proposing a change

1. Fork the repository and create a branch: `git checkout -b <change>`
2. Follow [AGENTS.md](AGENTS.md)
3. Validate locally: `go -C cli run . validate` (CI runs the same check)
4. Open a pull request describing the change and its intended use

For questions, open an issue.
