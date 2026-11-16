# Information Architecture and Anti-patterns

Source: https://diataxis.fr/ (Daniele Procida). This file covers directory organization for a docs site, how to write landing pages, and the anti-patterns that call for immediate splitting or rewriting.

## Directory Organization Rules

- One top-level directory per section: `tutorials/` `how-to-guides/` `reference/` `explanation/`, each with its own landing page.
- Hierarchy template (at every level, the landing page is likewise an overview):

  ```
  Home (landing)
  ├── Tutorials        ← landing + each lesson
  ├── How-to guides    ← landing + each guide (may be grouped further)
  ├── Reference        ← landing + each reference (mirrors the product structure)
  └── Explanation      ← landing + each topic
  ```

- **A landing page is an overview, not a list**: set the reader's expectations with a short introductory paragraph, then list the content.
- Keep any list to **7 items or fewer**; beyond that, group it or add another directory level.
- The docs home page states what the product is and routes readers through the four quadrants (one entry point each for study / work / look up / understand).

## Apply Flexibly, Don't Force the Directory Shape

- For two-dimensional problems (multiple user types x the four quadrants), split by "the product as the user sees it" first, rather than forcing the four boxes.
- Diátaxis is a **way of thinking**, not a directory shape that must be followed exactly — but never mix different types in one place.
- A concrete product may freely add dimensions below the four quadrants — product domains, languages, versions, and so on.

## Keeping Facts in Sync

- A reference's structure must mirror the product's real structure (types, interfaces, command trees); verify before writing, never from memory.
- After a code or product change, run reference docs through the same item-by-item check against the source to prevent drift.
- Example output, screenshots, and demo material must come from real runs; if a run isn't possible, say so honestly — never fabricate.

## Anti-patterns (Split Apart or Rewrite on Sight)

- Tutorials crammed with theory; how-to guides smuggling in beginner lessons; references containing procedural steps; explanations containing command lines.
- "You will learn..." tutorial openings; noun-phrase how-to titles; marketing copy mixed into references.
- Overlong lists; landing pages that are nothing but link dumps; tutorial steps with no expected output.
- One page serving two needs (a "quick reference + getting started" or "tutorial + reference" hybrid).
- Writing reference content from memory instead of checking source code or type definitions.

## Incremental Improvement Workflow

1. Look at the document in front of you (even if it's empty).
2. Ask: what could be improved?
3. Pick **one** small improvement you can make right now.
4. Do it. Repeat.

Don't try to finish in one pass; each improvement points to the next.
