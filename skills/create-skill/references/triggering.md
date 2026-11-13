# Triggering and Eval Practice

The description decides whether a skill gets invoked; this file distills the shared practices of well-known skill-creator implementations on trigger optimization and prompt evals.

## Triggering Mechanism

- The model sees only `name` + `description` and decides from those whether to load the body.
- Models tend to under-trigger: the description must state both what the skill does and when to use it, erring slightly aggressive when listing trigger words and scenarios.
- Simple single-step tasks may not trigger any skill — that is normal and not worth optimizing for.

## Writing a Good Description

- Put what it does in the first sentence, then list trigger scenarios and user phrasings ("Use when the user mentions X, Y, or wants Z").
- Cover phrasings that never name the skill: if the user says "help me tidy up this report", a report-handling skill should still trigger.
- Avoid abstract slogans; write concrete, matchable words.

## Test Prompt Design

- Start from 2-3 realistic prompts: concrete paths, field names, colloquial phrasing, even typos — not abstract instructions.
- Should-trigger cases should cover different phrasings and edge scenarios; among should-not-trigger cases, the most valuable are adjacent tasks that look similar but should not trigger.
- Overly simple cases ("read this file") test nothing.

## Eval Loop

1. Run each prompt both with the skill and without it, forming a control.
2. Judge both results and trajectory: looping around, re-reading files, and writing useless scripts all signal that the skill over-constrains.
3. Assertions must be objective, verifiable, and descriptively named; leave subjective quality to user review rather than forcing contrived assertions.
4. Generalize from individual feedback to unseen inputs; patches made for a single example are negative assets.

## Iterating to Convergence

- Fix → rerun → review the output with the user → fix again, until the user is satisfied or changes stop producing value.
- Trigger optimization of the description needs batches of query samples (roughly half should-trigger, half should-not); have the user confirm the samples before running.
- When improving an existing skill, keep the original `name` — no version-suffix renames; run trigger evals on the model and entry point the user actually uses.

## Sources

- Anthropic skill-creator's description optimization and evals design
- ZCode skill-creator's test prompts and review loop
- Anthropic's official skills writing guidance (explaining why beats MUST)
