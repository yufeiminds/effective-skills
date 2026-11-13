---
name: create-skill
description: >
  Turn recurring workflows or domain knowledge into a reusable agent skill: intent clarification, minimal SKILL.md drafting, progressive disclosure splitting, testing with real prompts, and iterative refinement. Use this skill whenever the user wants to create a new skill, improve an existing skill's triggering or wording, or decide whether content belongs in the body or in reference files.
license: MIT
metadata:
  author: yufeiminds
  version: "1.0.0"
---

# Create Skill

A workflow for creating and iteratively improving agent skills: clarify intent first, write a minimal viable draft, validate it with real prompts, and converge based on feedback.

## When to Apply

- The user wants to create a new skill from scratch, or wants to capture a manual process they have just repeated several times as a skill
- An existing skill triggers unreliably, produces unexpected output, or over-constrains the model into wasted work
- The user asks to optimize a skill's description, wording, or structure
- A decision is needed on whether a piece of detail belongs in the SKILL.md body or should be split into references, scripts, or assets

## Do Not Apply

- The task is one-off, with no expectation of repeated reuse
- The user just wants an answer or a single script, not a discoverable, reusable skill
- The user explicitly says no formal process is needed, just a quick draft

## Instructions

1. Clarify intent first: what should this skill enable the model to do, which user phrasings or scenarios should trigger it, what output shape is expected, and are there input/output examples that can lock behavior down. Prefer extracting workflows already seen in the current conversation (tools used, step order, corrections the user made) from history; ask the user only about the gaps, and get their confirmation before drafting.
2. Check first for an existing skill with the same name or the same domain (this repo's `skills/` directory plus user- and project-level skill directories). When improving an existing skill, keep its original `name` and directory name; do not add suffixes like `-v2`.
3. Draft a minimal SKILL.md: keep only `name` and `description` in the frontmatter; the description is the primary trigger signal, so state both what the skill does and when to trigger it — err toward listing trigger scenarios slightly aggressively rather than writing it vaguely.
4. Organize content around progressive disclosure: keep only the core workflow needed across scenarios in the body (under 500 lines); split language specifics, domain background, and long examples into `references/`, put executable scripts in `scripts/`, put template material in `assets/`, and state in the body when to read which file.
5. Write instructions as imperatives, explaining the "why" instead of piling up MUSTs and NEVERs; give a reason whenever a rule is not obvious — a model that understands the motivation is steadier than one that obeys mechanically.
6. Lock behavior with examples: if the output has a fixed structure, provide a literal template; if a specific call is expected, write out the call itself.
7. After drafting, construct 2-3 realistic test prompts (with concrete file paths, colloquial phrasing, even typos), confirm them with the user, then run them one by one. Judge both results and trajectory: if the model repeatedly reads files or loops around while following the skill's guidance, that signals over-constraining — subtract rather than add rules.
8. Iterate on user feedback: generalize from individual examples to unseen inputs, and delete instructions that produce no value; solidify recurring multi-step operations as scripts.
9. Once converged, confirm with the user whether the description's trigger hit rate still needs optimizing, then deliver.

## Authoring Heuristics

- The description is the only metadata resident in context (roughly 100 words or fewer); put all trigger information here, not hidden in the body.
- When the body approaches 500 lines, add another layer of splitting instead of piling on more.
- Reference files over 300 lines get a table of contents; when the body references a reference file, state clearly under what circumstances to read it.
- One skill does one thing; if it does several, split it into multiple focused skills.
- Split multi-target or multi-framework support into per-variant reference files (e.g. `aws.md`, `gcp.md`), keeping only selection logic in the body.
- Patches made for individual test examples are negative assets; when repeated explanations don't land, change the wording or structure instead of adding more constraints.
- Solidify multi-step processes that recur on every run as scripts — writing once beats reinventing every time.

## References

- [Skill anatomy and progressive disclosure](./references/anatomy.md)
- [Triggering and eval practice](./references/triggering.md)

Expand the corresponding document only when you need directory-structure details, context budget, or trigger eval methodology; do not load everything at once.

## Output

The final reply should include:

1. Intent summary (what it does, when it triggers)
2. Skill file structure and key design decisions
3. Test prompts and run results
4. Iteration history and how feedback was absorbed
5. Remaining risks and follow-up suggestions

## Guardrails

- Never overwrite or delete the user's existing skills; when improving, keep the original `name`, directory name, and established conventions
- Never create skills that contradict their stated purpose or that are misleading or malicious
- Never batch-run test prompts without user confirmation
- Never cram one-off details into the body; splitting beats bloat
