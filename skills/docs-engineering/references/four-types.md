# The Four Types and Their Writing Rules

Source: https://diataxis.fr/ (Daniele Procida). This file is an operational distillation of the Diátaxis methodology. Before writing any document, find here the rules for the quadrant it belongs to.

## The Core Mental Model: Four Docs, Four Needs

Documentation serves exactly four user needs, corresponding to four types. **Judge the need first, then choose the type**:

|                | Serves study (acquiring skill) | Serves work (applying skill) |
| -------------- | ------------------------------ | ---------------------------- |
| **Action** (doing) | **Tutorial**               | **How-to guide**             |
| **Cognition** (knowing) | **Explanation**       | **Reference**                |

- **Tutorial**: a lesson. Guides a beginner through a meaningful learning experience; the goal is to **learn**, not to get the task done.
- **How-to guide**: a recipe. Helps a user who **already has the basics** accomplish a specific real-world goal.
- **Reference**: a map. A neutral, accurate, authoritative description of the machine, for consultation rather than reading.
- **Explanation**: a discussion. Provides background, reasons, and perspective to support understanding of "why".

## The Compass: Judging Before You Write

For any piece of content, ask two questions first:

```
Is it guiding action, or conveying knowledge?
Does it serve the user's study, or the user's work?
```

| If the content... | And it serves...    | Then it must be |
| ----------------- | ------------------- | --------------- |
| Guides action     | Acquiring skill     | Tutorial        |
| Guides action     | Applying skill      | How-to guide    |
| Conveys knowledge | Applying skill      | Reference       |
| Conveys knowledge | Acquiring skill     | Explanation     |

**Mixing types is the root cause of the vast majority of documentation problems.** For content you cannot clearly classify: split it apart, or delete it.
Each type does only its own job; route the rest away with links.

## Tutorial (Study-Oriented)

- **Don't "teach"** — design an experience that lets learners learn on their own; learners learn by doing.
- Show the destination at the start ("In this tutorial we will..."), but never say "You will learn...".
- **Deliver visible results early and often**; every step has a perceptible outcome.
- **Maintain the "expected narrative"**: tell learners what they should see at each step ("the output should look like...") and what a failure signal looks like.
- Point out what learners should notice ("Did you notice...?").
- **Keep explanation to a minimum** ("we use HTTPS here because it's more secure; see the security model"), and link the rest to explanation.
- Follow a single path: ignore all options, branches, and alternatives.
- Allow and encourage repetition; every step is reversible and repeatable.
- Aim for **perfect reliability**: following it must succeed. Use "we" ("let's do...") and short imperative sentences.
- Tutorials are the part of the site where changes cascade most easily — the end-to-end story must always hold.

## How-to Guide (Goal-Oriented)

- **Start from the user's real goal**, not from what the machine can do. "How to configure a custom domain for your service" ✓; "explanation of the domains API endpoint" ✗.
- The title states precisely how to do something: **How to do X**. Avoid noun-phrase titles.
- Stay focused on the goal: no teaching, no explaining, no sprawling into reference — link out if needed.
- Non-linearity is allowed: steps can branch, with multiple entries and exits; acknowledge judgment ("if you are..., then...").
- Adapt to the complexity of real situations; **usability beats completeness** — reasonable starting and ending points suffice.
- Describe a logically ordered sequence of steps; mind the user's flow — the rhythm of thought and action must not be interrupted.

## Reference (Information-Oriented)

- **Describe only, and do nothing else**: no instruction, no explanation, no evaluation. Users come here to look up facts.
- Be concise, ordered, unambiguous, and authoritative; use **consistent, standard formatting patterns** so users find the expected information in the expected place.
- **Structure mirrors the product's structure** (APIs by endpoint, CLIs by command tree, listings by field).
- A moderate number of examples is fine, but they must serve the description, not turn into a tutorial.
- Factual statements: "X is defined in Y", "the subcommands are a, b, c", "must...; must not... unless...".
- Whatever can be auto-generated from code (type docs, API references) should be, to guarantee consistency with the code.
- Writing reference content from memory is taboo; before writing, verify against source code, type definitions, or actual runtime behavior.

## Explanation (Understanding-Oriented)

- Answer "Can you tell me about X?"; it is the only doc type suited to **reading away from the product**.
- Make connections: to related topics, history, design trade-offs, cross-domain analogies.
- Provide background and context: why it was designed this way, what the constraints are, what the alternatives are.
- Titles may take an "About" prefix: "About configuration sync", "About session lifecycle".
- **Opinion, judgment, and trade-offs are allowed** — this is the only type where weighing alternatives is on-topic.
- Hold the boundary: don't absorb tutorial steps or reference detail; draw a reasonable line around scope and keep it.
- Anchor the writing on a real "why" question.

## One-Line Self-Check for Each Type

- Tutorial: can a newcomer follow it and succeed on the first try?
- How-to guide: can an experienced user with a clear goal go straight to the answer?
- Reference: does the fact-looker find the expected information in the expected place?
- Explanation: after closing the doc, can the reader explain the "why"?
