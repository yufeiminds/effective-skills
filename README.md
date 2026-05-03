# effective-skills

A collection of [agent skills](https://github.com/vercel-labs/skills) by yufeiminds — reusable instruction sets that extend AI coding agents' capabilities.

## Install Skills

Use the `skills` CLI to install skills from this repository into your coding agent:

```shell
# Install all skills
npx skills add https://github.com/yufeiminds/effective-skills

# List available skills without installing
npx skills add https://github.com/yufeiminds/effective-skills --list

# Install a specific skill
npx skills add https://github.com/yufeiminds/effective-skills --skill <skill-name>
```

Supports [Claude Code](https://claude.ai/code), [Cursor](https://cursor.sh), [Codex](https://openai.com/codex), [OpenCode](https://opencode.ai), and [50+ more agents](https://github.com/vercel-labs/skills#supported-agents).

## Available Skills

| Skill | Description |
|-------|-------------|
| [clean-code](skills/clean-code/SKILL.md) | A concise refactoring skill for shrinking public APIs, building deeper modules, and cleaning structure across Rust, TypeScript, and MoonBit |

## Contributing

Want to add a new skill? See [CONTRIBUTING.md](CONTRIBUTING.md) for a step-by-step guide on how to create and submit skills to this repository.
