---
name: clean-code
description: >
  Refactor working code into deeper modules with smaller public interfaces by inventorying, converging, and verifying
  all relevant code in the user-defined scope. Use this skill whenever files grow overly long, responsibilities are
  tangled, import paths are fragile, naming is noisy, or error semantics leak.
license: MIT
metadata:
  author: yufeiminds
  version: "1.3.1"
---

# Clean Code

Converge complexity through structural refactoring — a full inventory first, then item-by-item execution — without breaking behavior.

## When to Apply

- Hand-written source files have grown too long, or a single module serves multiple reasons to change
- A single function mixes business orchestration with low-level details
- Intermediate types, helper functions, or pure pass-through layers expose implementation details to the outside
- Names repeat their parent's semantics, or the same concept exists as both a flat file and a directory module
- Import paths are overly sensitive to directory depth, requiring chains of relative fallbacks or climbs across layers to locate a module
- Errors directly leak IO, network, parsing, or framework details
- Test files are bloated and already hard to navigate by behavior

## Do Not Apply

- The issue is just a one-off local bug with no structural debt
- Structural changes are forbidden in the current iteration
- There is no minimal means of verification to confirm behavior stays unchanged

## Instructions

1. Inventory all relevant code in the user-defined scope first, rather than fixating on the current file or a minimal slice: implementations, entry points, export points, import paths, callers, tests, directory hierarchy, dead code, and empty directories are all in scope.
2. After completing the full inventory, build a todo list covering the entry points, implementations, callers, exports, tests, and cleanup items for this round; proceed in dependency and risk order.
3. Target deep modules: hide more implementation complexity behind fewer, more stable interfaces, instead of splitting for splitting's sake.
4. Converge visibility, export surface, and import paths by default; callers should depend on the root package's stable public paths, not relative directory climbing or layer-by-layer parent fallbacks.
5. Delete shallow wrappers, mirror modules, duplicated logic, empty directories, and empty compatibility stubs; keep exactly one consistent hierarchical expression per concept. If an old path must be kept temporarily, make it a facade with a clear responsibility and migration notes, not a hollow shell file.
6. Make names express responsibility directly; if a file or module name spans multiple words, first evaluate whether renaming, sinking into a parent module, or further splitting can converge it to a shorter, more local name — keep two-word or multi-word naming only when it significantly reduces ambiguity.
7. Map low-level errors to domain errors, replacing bare primitive types in core paths with semantically meaningful types.
8. Treat large test files together with their implementations: keep white-box tests close to the implementation, split combinatorial tests along a single uniform dimension, and do not keep flat piles that are hard to navigate.
9. After each completed item, run the verification closest to that slice; when you cannot run it, state the reason, the affected scope, and the residual risk.

## Restructuring Heuristics

- Split files along responsibility boundaries, not by averaging line counts; entry files only declare, re-export, and act as facades.
- Keep a parent module only when it aggregates multiple submodules, hides policy, or unifies a contract; otherwise delete the forwarding layer.
- Do not hollow out old entry points into responsibility-free compatibility stubs; either delete the dead path or keep a stable facade carrying migration responsibility.
- Prefer the root package's stable paths or stable aliases for imports over directory-relative climbing; refactored paths should be more stable, not more dependent on current directory depth.
- Two-word or multi-word file names are not the default option; first evaluate whether renaming, sinking, or further splitting can converge them into shorter, local semantics.
- Split large test files along a single dimension such as workflow, error cases, or input family; keep white-box tests close to the implementation.
- Add doc comments to stable exports, spelling out purpose, examples, and caveats.

## References

- [Clean Code summary](./references/clean-code.md)
- [A Philosophy of Software Design summary](./references/a-philosophy-of-software-design.md)
- Deep Module: when judging whether a split is justified, first check whether it genuinely reduces the number of concepts callers must understand.

These documents are supporting references only; the current repository's structural convergence and verification requirements take priority.


## Output

The final reply must include:

1. Change summary
2. Key structural decisions
3. Compatibility notes
4. Verification results
5. Residual risks

## Guardrails

- Prefer minimal necessary changes, but the analysis must cover all relevant code in the user-defined scope
- Do not stuff business implementation into entry files
- Do not make breaking API changes without authorization
- Do not mistake an increase in file count for design improvement
