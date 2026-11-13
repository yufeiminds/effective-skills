# Key Takeaways from "A Philosophy of Software Design"

The book's center of gravity is not code style but a more upstream question: how to make a system easier to understand and change. Ousterhout collapses nearly every judgment of design quality into one question: did complexity go down?

## Key Claims

- The first goal of good design is reducing complexity, not adding structural nouns.
- Complexity is not "a lot of code"; it is too much that those who come later must know at once.
- Design should prioritize reducing the reader's cognitive load, not merely making local code look tidy.

## The Three Telltale Symptoms of Complexity

- Change amplification: one requirement change forces edits in many places.
- Cognitive load: changing one piece of code requires holding too much prerequisite knowledge in your head at once.
- Unknown unknowns: you do not even know where you still need to look for the key facts.

## The Main Sources of Complexity

- Too many dependencies: one change drags many modules along.
- Information leakage: knowledge that should be encapsulated inside a module is forced onto its callers to understand.
- Inconsistency: similar concepts expressed with different rules, different names, and different error semantics in different places.

## Key Content

### 1. Deep Module is the core structural criterion

- A good module should hide substantial implementation complexity behind a small interface.
- If extracting a module merely adds a layer of names and forwarding, it is a shallow module.
- Judge a split's success not by file count but by whether callers genuinely need to know less.

### 2. Push complexity downward, not upward

- Callers should depend on stable interfaces, not have to understand internal flows, timing, or directory structure.
- Pass-through methods, mirror APIs, and layers that only forward without abstracting usually keep spreading complexity outward.
- A truly good abstraction lets users remember only the interface, without having to remember internal details.

### 3. General-purpose mechanisms are usually deeper than specialized ones

- Modules over-tuned to the current scenario tend to have narrow interfaces and short lifespans.
- If a mechanism covers a broader class of problems while keeping a simple interface, it is usually deeper than an ad-hoc specialized solution.
- The second edition strengthens this point, expanding "General-Purpose Modules are Deeper" as a highlighted topic.

### 4. Do not over-split related logic

- A split is only justified when the interface is simpler than the original code.
- If two functions must be read by jumping back and forth, memorized together, and reasoned about together, they are entangled — splitting them apart makes them harder to understand.
- So "more functions" does not automatically mean "better design"; often keeping related logic together is clearer.

### 5. Comments are a design tool, not just a patch

- Interface comments define the contract, telling callers what they may depend on.
- Implementation comments capture key facts not obvious from the code itself, design intent, and counterintuitive constraints.
- Data structure comments explain field meanings, lifetimes, and invariants.
- The book's stance is not "the more comments the better" but "missing key comments directly increase complexity."

### 6. Reduce complexity by defining special cases out of existence

- Special branches, exception protocols, and piecemeal patches rapidly amplify cognitive cost.
- One important strategy is "define errors out of existence": redesign the interface, data model, or flow so that situations that once needed special handling become ordinary paths.
- Good design often shows up as: callers never need to know a class of exceptional detail ever existed.

### 7. Design requires strategic thinking, not purely tactical progress

- Tactical programming focuses on "get the immediate change done first"; Strategic programming invests more time thinking about long-term complexity.
- Design is not cleanup after the code is written; it means considering boundaries, interfaces, dependencies, and information hiding before writing.
- The second edition adds "Decide What Matters," stressing that you must first identify the factors that truly drive complexity, then focus your firepower on them.

### 8. Consistency is itself a complexity-reduction tool

- Consistent naming, error semantics, interface style, and structural conventions lower the reader's context-switching cost.
- Inconsistency creates extra judgment: is this the same thing, but with yet another set of rules here?

## The Most Valuable Questions to Take Away

- After this split, do callers genuinely need to know less, or have they merely been introduced to a few more files?
- Does this interface hide complexity, or expose it in a new disguise?
- When one requirement changes, why do so many places need edits? Has change amplification already set in?
- Is this code hard to read because the logic is inherently complex, or because key information is scattered across too many places?
- Does this need one more layer of splitting, or should the tightly coupled logic be pulled back together?

## The Public Points of Disagreement with "Clean Code"

- It is warier of "small functions" being mechanically executed into over-splitting.
- It places more emphasis on the role of interface comments and implementation comments in reducing complexity.
- It prefers making design judgments first and then implementing, rather than over-tacticalizing the development process.

## Public Sources

- The book page and second-edition notes on John Ousterhout's author homepage
- Publications and course/design-related pages on the author's homepage
- The `johnousterhout/aposd-vs-clean-code` public discussion repository
