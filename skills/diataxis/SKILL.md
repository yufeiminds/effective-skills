---
name: diataxis
description: >
  Documentation engineering methodology based on the Diátaxis framework. Use it when planning,
  writing, reviewing, or restructuring documentation for a technical product, to answer where
  a piece of content belongs, what form it should take, and who it is for. Trigger it when a
  document mixes responsibilities and needs splitting, when a docs site's information
  architecture needs organizing, or when judging which quadrant a document belongs to — even
  if the user doesn't name Diátaxis.
license: MIT
metadata:
  author: yufeiminds
  version: "2.0.0"
---

# Diátaxis

A documentation methodology with the Diátaxis four quadrants as its core mental model: judge the need first, then pick the type. Each type does only its own job; mixed types get split apart.

## When to Apply

- Planning a docs site structure for a technical product, or creating any new document
- A document mixes responsibilities (theory inside a tutorial, steps inside a reference) and needs to be split and restructured
- Reviewing documentation quality, deciding which quadrant content belongs in
- Organizing a docs site's directories, landing pages, or navigation

## Do Not Apply

- One-off communication text (weekly reports, announcements, emails) that is not technical documentation
- Purely mechanical edits — fixing typos, adjusting formatting — that don't touch type placement
- The documentation system already has explicit conventions and the task is just filling in content per them, with no type judgment needed

## Instructions

1. Before writing, run the compass: is this content guiding action or conveying knowledge? Does it serve study or work? Place it in one of the four quadrants accordingly — Tutorial (learning), How-to guide (getting it done), Reference (looking up), Explanation (understanding).
2. Split apart or delete content you cannot clearly classify; one document serves one need, and everything else gets routed away with links.
3. Write according to the rules of the assigned type: read the relevant quadrant's section in [The Four Types and Their Writing Rules](./references/four-types.md) before writing — don't go on instinct.
4. For reference content (CLIs, APIs, fields, ports, config options), verify against the source of truth (source code, type definitions, actual runtime behavior) before writing; never write from memory. A reference's structure must mirror the product's structure.
5. After writing, self-verify against the type: a tutorial succeeds end to end when followed exactly, with a perceptible outcome at every step; a how-to guide fits real goals and its steps are never interrupted; a reference has consistent formatting and predictable locations; an explanation holds its boundaries and does not absorb steps or details.
6. Organize top-level directories by the four quadrants; write landing pages as overviews, not lists; keep every list to at most 7 items. Read [Information Architecture and Anti-patterns](./references/information-architecture.md) before planning a directory structure.
7. Improve existing documentation incrementally: pick one small improvement you can make right now, make it, repeat; don't try to finish in one pass.
8. Before delivering, scan the anti-pattern list; wherever one appears, split it apart or rewrite.

## Writing Heuristics

- Mixing types is the root cause of the vast majority of documentation problems; one page serving two needs (a "quick reference + getting started" or "tutorial + reference" hybrid) is the most common pathology.
- Tutorials don't "teach": design an experience that lets learners learn on their own; show the destination at the start, and never say "You will learn".
- Title how-to guides "How to do X"; start from the user's real goal, teach nothing, explain nothing; usability beats completeness.
- References only describe: no instruction, no explanation, no evaluation; factual statements in consistent, standard format patterns.
- Explanation is the only place where opinion, trade-offs, and analogy are allowed; anchor the writing on a real "why" question.
- Fact drift begins the moment the code changes; after a product change, check reference documentation item by item against the source.

## References

- [The Four Types and Their Writing Rules](./references/four-types.md)
- [Information Architecture and Anti-patterns](./references/information-architecture.md)

Before writing body text, read the rules for the type's quadrant; read the information architecture doc when planning or tidying directories. Load only the one you need — don't load everything at once.

## Output

The final response should include:

1. Type judgment and its basis (position within the four quadrants)
2. Summary of content or structure changes
3. Consistency check against the source of truth (when reference content is involved)
4. Anti-pattern scan results
5. Residual risks and suggested follow-up improvements

## Guardrails

- Never mix two types in one document; splitting beats accommodating
- Never write reference content from memory; always verify against the source of truth
- Never delete the user's existing valid content; when restructuring, judge first, then move, then clean up
- Never pass off fabricated example output, screenshots, or run results as real content
