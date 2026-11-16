# effective-skills

Curated agent skills by [yufeiminds](https://github.com/yufeiminds) — small, language-agnostic, and composable. Take them, hack on them, make them yours.

## Installation

A plugin updates itself. The `skills` CLI copies editable files into your project, and you update them by hand. **Pick one route per agent** — installing both gives you every skill twice.

<details>
<summary><strong>Claude Code</strong></summary>

```text
/plugin marketplace add yufeiminds/effective-skills
/plugin install effective-skills@yufeiminds-skills
```

Updates itself when the marketplace source moves.

</details>

<details>
<summary><strong>Any other agent</strong> (editable files)</summary>

```shell
npx skills add https://github.com/yufeiminds/effective-skills
# or a single skill
npx skills add https://github.com/yufeiminds/effective-skills --skill <skill-name>
```

Update by hand with `npx skills update`. Works across Claude Code, Codex, Cursor, OpenCode, and [50+ more agents](https://github.com/vercel-labs/skills#supported-agents).

</details>

## Reference

These split on one axis: who can invoke them. All skills in this repo are **model-invoked** — reachable by you or reached for automatically by the agent when the task fits; they hold the reusable discipline.

**Model-invoked**

- **[clean-code](skills/clean-code/SKILL.md)**: Language-agnostic refactoring: deep modules, smaller public interfaces, converged structure.
- **[diataxis](skills/diataxis/SKILL.md)**: Diátaxis-based documentation methodology: where content belongs, per-quadrant writing rules, docs-site information architecture.

## Development

Engineering rules — structure, versioning, validation, release flow — live in [AGENTS.md](AGENTS.md). The validator in `cli/` enforces them, and CI runs it on every push. Toolchains are pinned in `mise.toml`; run `mise install` once.

## License

[MIT](LICENSE) © [yufeiminds](https://github.com/yufeiminds). The development dependency under `.agents/` keeps its own license (Apache-2.0).
