# Skill Anatomy and Progressive Disclosure

The subject here is not any single platform's spec but the structural conventions shared by mainstream agent skill implementations: how to organize files so that resident context is minimal and on-demand loading is maximal.

## Key Tenets

- A skill's main cost is context: the description is resident, the body loads on trigger, and supporting files are read on demand.
- Structure serves load order: the more often content is needed, the closer it sits to the entry point.
- One directory, one skill; keep `name` consistent with the directory name.

## Three-Layer Loading Model

1. **Metadata** (`name` + `description`): always in context, roughly 100 words or fewer.
2. **SKILL.md body**: loaded only on trigger, target under 500 lines.
3. **Supporting resources**: read on demand, with no size limit in principle.

## Directory Anatomy

```
skill-name/
├── SKILL.md          # Required; frontmatter + body
├── references/       # Documents loaded on demand (language specifics, domain background, long examples)
├── scripts/          # Executable scripts the model can call directly
└── assets/           # Templates and material used for output
```

## Frontmatter Fields

| Field | Required | Notes |
|------|------|------|
| `name` | ✅ | kebab-case, 1-64 characters, matches the directory name |
| `description` | ✅ | What it does + when to trigger; the primary signal deciding whether the body loads |
| `license` | ✗ | e.g. `MIT` |
| `metadata.author` | ✗ | Author name or GitHub username |
| `metadata.version` | ✗ | Semantic version string |

## When to Split

- Body approaching 500 lines: split by domain or scenario into `references/<domain>.md`, keeping only selection logic and "when to read which file" in the body.
- Multi-target or multi-framework support: split by variant (e.g. `aws.md` / `gcp.md` / `azure.md`) and read only the relevant variant after triggering.
- Reference files over 300 lines get a table of contents.
- Multi-step operations that recur on every run: solidify them as scripts under `scripts/` — writing once beats reinventing every time.
- Output requires fixed templates, icons, or fonts: put them in `assets/`.

## Common Misreadings

- "The finer the split, the better" does not hold: if a split-out file has to be read on every trigger, it lengthens the body under another name.
- "The fuller the supporting directories, the better" does not hold: files that add no on-demand value are pure liabilities.
- Support for fields beyond `name` and `description` varies by platform; do not depend on non-required fields.

## Sources

- Anthropic skill-creator public implementation (the skills repo and the agentskills.io spec)
- ZCode skill-creator plugin (bundled-skills)
- SKILL.md conventions of vercel-labs/skills
