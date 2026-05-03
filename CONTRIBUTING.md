# Contributing to effective-skills

Thank you for your interest in contributing! This repository hosts a collection of [agent skills](https://github.com/vercel-labs/skills) — reusable instruction sets that extend AI coding agents' capabilities.

## Adding a New Skill

### 1. Create a skill directory

Each skill lives in its own subdirectory under `skills/`:

```
skills/
└── your-skill-name/
    ├── SKILL.md
    └── assets/
        └── appendix.md
```

Use `SKILL.md` as the concise entrypoint. Put long language-, framework-, or workflow-specific details into optional companion Markdown files under `assets/` only when they materially reduce prompt bloat.

Run the following command to scaffold a new skill:

```shell
npx skills init your-skill-name
```

Or create the directory manually:

```shell
mkdir -p skills/your-skill-name
```

### 2. Write the SKILL.md file

Every skill requires a `SKILL.md` file with YAML frontmatter and Markdown content.

Keep the top-level skill concise. If detailed sub-guides are necessary, link them from `SKILL.md` and store them under `assets/` instead of expanding the entry file indefinitely.

**Frontmatter fields:**

| Field | Required | Description |
|-------|----------|-------------|
| `name` | ✅ | Unique identifier for the skill (kebab-case) |
| `description` | ✅ | What the skill does and when to use it (used by the agent to decide when to apply it) |
| `license` | ✗ | License for the skill (e.g. `MIT`) |
| `metadata.author` | ✗ | Author name or GitHub username |
| `metadata.version` | ✗ | Semantic version string (e.g. `"1.0.0"`) |

**Template:**

```markdown
---
name: your-skill-name
description: >
  A clear description of what this skill does and when it should be applied.
  The agent uses this description to decide when to activate the skill.
license: MIT
metadata:
  author: your-github-username
  version: "1.0.0"
---

# Your Skill Title

Brief overview of the skill.

## When to Apply

- Describe the scenarios or triggers for this skill
- Be specific so the agent knows exactly when to use it

## Instructions

Step-by-step instructions for the agent to follow.

### Step 1

...

### Step 2

...

## Examples

Provide examples to clarify the expected behavior.
```

### 3. Review existing skills

See [`skills/clean-code/SKILL.md`](skills/clean-code/SKILL.md) for a real skill that keeps `SKILL.md` focused and moves language-specific detail into `assets/` companion docs.

### 4. Open a Pull Request

1. Fork the repository
2. Create a branch: `git checkout -b add-your-skill-name`
3. Add your skill under `skills/your-skill-name/SKILL.md`
4. Commit and push your changes
5. Open a Pull Request describing the skill and its intended use

## Skill Naming Guidelines

- Use **kebab-case** for directory names and the `name` field (e.g., `code-review-guidelines`)
- Choose a **descriptive name** that reflects what the skill does
- Avoid overly generic names like `helper` or `tool`

## Writing Good Skill Instructions

- **Be specific**: Vague instructions lead to inconsistent agent behavior
- **Use imperative mood**: Write instructions as direct commands (e.g., "Check for X", "Always do Y")
- **Include examples**: Concrete examples help agents understand the expected output
- **Keep it focused**: Each skill should do one thing well; split complex skills into multiple focused ones
- **Test your skill**: Try installing it locally and verify the agent behaves as expected

## Installing Skills Locally

To test skills from this repository:

```shell
# Install all skills
npx skills add https://github.com/yufeiminds/effective-skills

# Install a specific skill
npx skills add https://github.com/yufeiminds/effective-skills --skill your-skill-name

# List available skills without installing
npx skills add https://github.com/yufeiminds/effective-skills --list
```

## Questions?

Open an issue if you have questions about creating or contributing skills.
