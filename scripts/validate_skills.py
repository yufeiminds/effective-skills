#!/usr/bin/env python3
"""Structural validator for skills in this repository.

Checks every directory under skills/ for:
- SKILL.md presence and a well-formed frontmatter block
- required frontmatter fields: name, description
- name is kebab-case and matches the directory name
- the skill is listed in README.md's Available Skills table

Warnings (non-fatal): missing license / metadata.version, description over 1024 chars.

Usage: python3 scripts/validate_skills.py
Exit code 0 if all skills pass, 1 otherwise.
"""

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
SKILLS_DIR = ROOT / "skills"
README = ROOT / "README.md"

KEBAB = re.compile(r"^[a-z0-9]+(-[a-z0-9]+)*$")


def parse_frontmatter(text):
    """Return (fields, raw_frontmatter, error). fields is a flat dict; folded and
    nested scalars are joined into a single string per top-level key."""
    if not text.startswith("---"):
        return None, None, "missing frontmatter block"
    match = re.match(r"^---\n(.*?)\n---\s*\n", text, re.DOTALL)
    if not match:
        return None, None, "unterminated frontmatter block"
    raw = match.group(1)

    fields = {}
    key = None
    value = []
    for line in raw.splitlines():
        if line[:1] in (" ", "\t") and key is not None:
            value.append(line.strip())
            continue
        if key is not None:
            fields[key] = " ".join(value)
        m = re.match(r"^([A-Za-z0-9_-]+):\s*(.*)$", line)
        if not m:
            return None, raw, f"malformed frontmatter line: {line!r}"
        key, rest = m.group(1), m.group(2)
        value = [rest] if rest else []
    if key is not None:
        fields[key] = " ".join(value)
    return fields, raw, None


def validate_skill(skill_dir):
    errors, warnings = [], []
    skill_md = skill_dir / "SKILL.md"
    if not skill_md.is_file():
        return ["SKILL.md not found"], []

    fields, raw, err = parse_frontmatter(skill_md.read_text(encoding="utf-8"))
    if err:
        return [err], []

    name = fields.get("name", "")
    if not name:
        errors.append("frontmatter is missing 'name'")
    elif not KEBAB.match(name):
        errors.append(f"'name' is not kebab-case: {name!r}")
    elif name != skill_dir.name:
        errors.append(f"'name' ({name}) does not match directory ({skill_dir.name})")

    desc = fields.get("description", "").lstrip("> ").strip()
    if not desc:
        errors.append("frontmatter is missing 'description'")
    elif len(desc) > 1024:
        warnings.append(f"description is {len(desc)} chars; some platforms truncate at 1024")

    if not re.search(r"^\s+author:", raw, re.MULTILINE):
        warnings.append("no metadata.author field")
    if not re.search(r"^\s+version:", raw, re.MULTILINE):
        warnings.append("no metadata.version field")

    readme = README.read_text(encoding="utf-8") if README.is_file() else ""
    if f"(skills/{skill_dir.name}/SKILL.md)" not in readme:
        errors.append("not listed in README.md Available Skills table")

    return errors, warnings


def main():
    if not SKILLS_DIR.is_dir():
        print("error: skills/ directory not found")
        return 1

    failed = False
    count = 0
    for skill_dir in sorted(p for p in SKILLS_DIR.iterdir() if p.is_dir()):
        count += 1
        errors, warnings = validate_skill(skill_dir)
        print(f"{skill_dir.name:24s} {'FAIL' if errors else 'ok'}")
        for e in errors:
            print(f"  error:   {e}")
            failed = True
        for w in warnings:
            print(f"  warning: {w}")

    if count == 0:
        print("error: no skills found under skills/")
        return 1

    print(f"\n{count} skill(s) checked")
    return 1 if failed else 0


if __name__ == "__main__":
    sys.exit(main())
