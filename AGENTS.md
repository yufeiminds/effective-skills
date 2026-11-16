# AGENTS.md

Engineering rules for this repository. Agents maintaining or extending the
skill collection must follow them. Humans contributing see
[CONTRIBUTING.md](CONTRIBUTING.md).

## Layout

- `skills/<name>/` — published skills: `SKILL.md`, optional `references/`,
  `scripts/`, `assets/`, and `evals/evals.json`
- `cli/` — Go structure validator (its own module); CI runs its tests and
  `go -C cli run . validate` on every push
- `.claude-plugin/marketplace.json` — the repo is its own Claude Code plugin
  marketplace
- `mise.toml` — pinned toolchains; `cli/go.mod` mirrors the `go` pin, bump
  both together
- `.tmp/anthropics-skills` — local shallow mirror of upstream skills
  (gitignored)

## Skill invariants

Mechanically enforced by the validator:

1. Frontmatter follows the [agentskills.io spec](https://agentskills.io/specification):
   closed field set (`name`, `description`, `license`, `compatibility`,
   `allowed-tools`, `metadata`); `name` kebab-case matching the directory;
   `description` ≤1024 chars stating what the skill does AND when to use it,
   written for the model with trigger phrasing ("Use when…, mentions…,
   asks for…").
2. Listed in README's Reference section, skill name linked to its `SKILL.md`.
3. `metadata.author` and `metadata.version` present; version is semver —
   patch for wording, minor for new guidance, major for behavior or identity
   changes (renames).
4. `license: MIT`.
5. `SKILL.md` body ≤500 lines, relative links resolve, description carries
   trigger guidance.
6. `.claude-plugin/marketplace.json` carries **no** version field. A
   marketplace-level version drifts the moment any skill changes; per-skill
   versions live in each `SKILL.md` frontmatter, where the changelog and
   validator can see them.

Checked by hand when touching files (not yet validator-enforced):

- `.claude-plugin/marketplace.json`'s `plugins[].skills` array lists every
  published skill, and only published skills.
- Every skill has `evals/evals.json` with 2+ realistic prompts. Subjective
  outputs (methodology, writing) get qualitative `expected_output` and no
  assertions; objective outputs may add verifiable `expectations`.

## Versioning and CHANGELOG

Bump the skill's `metadata.version` on every change to it, and record the
change under a dated section in [CHANGELOG.md](CHANGELOG.md).

## Install commands (single source)

The commands below are canonical. README copies them verbatim; update both
together.

```text
# Claude Code plugin (updates itself when the marketplace source moves)
/plugin marketplace add yufeiminds/effective-skills
/plugin install effective-skills@yufeiminds-skills

# skills CLI (editable copies; update by hand)
npx skills add https://github.com/yufeiminds/effective-skills
npx skills add https://github.com/yufeiminds/effective-skills --skill <skill-name>
```

Pick one route per agent; installing both loads every skill twice.

## Validation

```shell
go -C cli run . validate            # structure + conventions
go -C cli test ./...                # validator unit tests
mise install                        # pinned toolchains (go mirrors cli/go.mod)
```

After changing a skill, a manifest, or the validator: run validate, run the
tests, and record user-visible changes in CHANGELOG.md.

## Skill authoring

- Keep `SKILL.md` under 500 lines; push detail into `references/` and say
  when to read each file. Progressive disclosure: the description is always
  in context, the body loads on trigger, references load on demand.
- The description is the only trigger mechanism. Write it for the model:
  what it does, when to use it, and the words users actually say.
- Explain why over heavy-handed MUSTs; bundle deterministic or repeated work
  as `scripts/` with useful error messages.
- Authoring and eval loop: use the vendored
  [`skill-creator`](.agents/skills/skill-creator/SKILL.md) (Apache-2.0
  development dependency; not part of the published skill set and not
  discovered by `npx skills add . --list`).

## Official resources

- Spec: <https://agentskills.io/specification>
- Authoring guides: [best practices](https://agentskills.io/skill-creation/best-practices.md),
  [optimizing descriptions](https://agentskills.io/skill-creation/optimizing-descriptions.md),
  [evaluating skills](https://agentskills.io/skill-creation/evaluating-skills.md);
  index at <https://agentskills.io/llms.txt>
- Reference validator: [skills-ref](https://github.com/agentskills/agentskills/tree/main/skills-ref)
  (Python, demo-grade) — `cli/` re-implements its rules in Go

## Development tooling: skill-creator

`skills-lock.json` points at the local mirror under `.tmp/` by design; the
upstream repository is <https://github.com/anthropics/skills>. First time, or
after deleting `.tmp/`:

```shell
git clone --depth 1 https://github.com/anthropics/skills.git .tmp/anthropics-skills
```

To sync to upstream's latest commit and reinstall the tracked copy:

```shell
git -C .tmp/anthropics-skills fetch --depth 1 origin main
git -C .tmp/anthropics-skills reset --hard FETCH_HEAD
npx skills update skill-creator
```

Review the diff under `.agents/skills/skill-creator/` and commit it together
with any change to `skills-lock.json`.
