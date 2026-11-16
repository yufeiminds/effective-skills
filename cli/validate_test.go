package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const (
	validDescription = "Use when validating test fixtures against the spec."
	validFrontmatter = `name: sample-skill
description: Use when validating test fixtures against the spec.
license: MIT
metadata:
  author: tester
  version: "1.0.0"
`
)

// writeSkill creates dir/SKILL.md with the given frontmatter and body.
func writeSkill(t *testing.T, dir, frontmatter, body string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\n"+frontmatter+"---\n"+body), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// writeRaw creates dir/SKILL.md with arbitrary content.
func writeRaw(t *testing.T, dir, content string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func hasMessage(res *skillResult, substr string) bool {
	for _, e := range res.Errors {
		if strings.Contains(e, substr) {
			return true
		}
	}
	for _, w := range res.Warnings {
		if strings.Contains(w, substr) {
			return true
		}
	}
	return false
}

func TestValidSkill(t *testing.T) {
	dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), validFrontmatter, "# Sample\n\nBody.\n")
	res := validateSkill(dir, "")
	if len(res.Errors) != 0 {
		t.Fatalf("expected no errors, got %v", res.Errors)
	}
	if len(res.Warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", res.Warnings)
	}
	if res.Name != "sample-skill" {
		t.Fatalf("expected name sample-skill, got %q", res.Name)
	}
}

func TestValidSkillMinimal(t *testing.T) {
	fm := "name: sample-skill\ndescription: Use when validating test fixtures.\n"
	dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "Body.\n")
	res := validateSkill(dir, "")
	if len(res.Errors) != 0 {
		t.Fatalf("expected no errors, got %v", res.Errors)
	}
	if !hasMessage(res, "no metadata.author field") || !hasMessage(res, "no metadata.version field") {
		t.Fatalf("expected metadata warnings, got %v", res.Warnings)
	}
}

func TestMissingName(t *testing.T) {
	fm := "description: Use when validating.\n"
	dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
	res := validateSkill(dir, "")
	if !hasMessage(res, "Missing required field in frontmatter: name") {
		t.Fatalf("expected missing name error, got %v", res.Errors)
	}
}

func TestEmptyName(t *testing.T) {
	fm := "name: \"\"\ndescription: Use when validating.\n"
	dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
	res := validateSkill(dir, "")
	if !hasMessage(res, "Field 'name' must be a non-empty string") {
		t.Fatalf("expected empty name error, got %v", res.Errors)
	}
}

func TestNameTooLong(t *testing.T) {
	long := strings.Repeat("a", 65)
	fm := "name: " + long + "\ndescription: Use when validating.\n"
	dir := writeSkill(t, filepath.Join(t.TempDir(), long), fm, "")
	res := validateSkill(dir, "")
	if !hasMessage(res, "exceeds 64 character limit (65 chars)") {
		t.Fatalf("expected length error, got %v", res.Errors)
	}
}

func TestNameUppercase(t *testing.T) {
	fm := "name: Sample-Skill\ndescription: Use when validating.\n"
	dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
	res := validateSkill(dir, "")
	if !hasMessage(res, "Skill name 'Sample-Skill' must be lowercase") {
		t.Fatalf("expected lowercase error, got %v", res.Errors)
	}
}

func TestNameLeadingHyphen(t *testing.T) {
	fm := "name: -sample\ndescription: Use when validating.\n"
	dir := writeSkill(t, filepath.Join(t.TempDir(), "sample"), fm, "")
	res := validateSkill(dir, "")
	if !hasMessage(res, "Skill name cannot start or end with a hyphen") {
		t.Fatalf("expected hyphen error, got %v", res.Errors)
	}
}

func TestNameConsecutiveHyphens(t *testing.T) {
	fm := "name: sample--skill\ndescription: Use when validating.\n"
	dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
	res := validateSkill(dir, "")
	if !hasMessage(res, "Skill name cannot contain consecutive hyphens") {
		t.Fatalf("expected consecutive hyphens error, got %v", res.Errors)
	}
}

func TestNameInvalidCharacters(t *testing.T) {
	fm := "name: sample_skill\ndescription: Use when validating.\n"
	dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
	res := validateSkill(dir, "")
	if !hasMessage(res, "contains invalid characters. Only letters, digits, and hyphens are allowed.") {
		t.Fatalf("expected invalid characters error, got %v", res.Errors)
	}
}

func TestNameDirectoryMismatch(t *testing.T) {
	fm := "name: sample-skill\ndescription: Use when validating.\n"
	dir := writeSkill(t, filepath.Join(t.TempDir(), "other-name"), fm, "")
	res := validateSkill(dir, "")
	if !hasMessage(res, "Directory name 'other-name' must match skill name 'sample-skill'") {
		t.Fatalf("expected directory mismatch error, got %v", res.Errors)
	}
}

func TestDescriptionChecks(t *testing.T) {
	t.Run("missing", func(t *testing.T) {
		fm := "name: sample-skill\n"
		dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
		res := validateSkill(dir, "")
		if !hasMessage(res, "Missing required field in frontmatter: description") {
			t.Fatalf("expected missing description error, got %v", res.Errors)
		}
	})
	t.Run("empty", func(t *testing.T) {
		fm := "name: sample-skill\ndescription: \"\"\n"
		dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
		res := validateSkill(dir, "")
		if !hasMessage(res, "Field 'description' must be a non-empty string") {
			t.Fatalf("expected empty description error, got %v", res.Errors)
		}
	})
	t.Run("too long", func(t *testing.T) {
		fm := "name: sample-skill\ndescription: " + strings.Repeat("a", 1025) + "\n"
		dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
		res := validateSkill(dir, "")
		if !hasMessage(res, "Description exceeds 1024 character limit (1025 chars)") {
			t.Fatalf("expected length error, got %v", res.Errors)
		}
	})
}

func TestCompatibilityChecks(t *testing.T) {
	t.Run("not a string", func(t *testing.T) {
		fm := "name: sample-skill\ndescription: Use when validating.\ncompatibility:\n  - requires python\n"
		dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
		res := validateSkill(dir, "")
		if !hasMessage(res, "Field 'compatibility' must be a string") {
			t.Fatalf("expected type error, got %v", res.Errors)
		}
	})
	t.Run("null is allowed", func(t *testing.T) {
		fm := "name: sample-skill\ndescription: Use when validating.\ncompatibility:\n"
		dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
		res := validateSkill(dir, "")
		if hasMessage(res, "compatibility") {
			t.Fatalf("null compatibility should pass, got %v", res.Errors)
		}
	})
	t.Run("too long", func(t *testing.T) {
		fm := "name: sample-skill\ndescription: Use when validating.\ncompatibility: " + strings.Repeat("a", 501) + "\n"
		dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
		res := validateSkill(dir, "")
		if !hasMessage(res, "Compatibility exceeds 500 character limit (501 chars)") {
			t.Fatalf("expected length error, got %v", res.Errors)
		}
	})
}

func TestUnexpectedField(t *testing.T) {
	fm := "name: sample-skill\ndescription: Use when validating.\nfoo: bar\nmodel: haiku\n"
	dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
	res := validateSkill(dir, "")
	if !hasMessage(res, "Unexpected fields in frontmatter: [foo, model]") {
		t.Fatalf("expected unexpected fields error, got %v", res.Errors)
	}
}

func TestMetadataChecks(t *testing.T) {
	t.Run("nested mapping rejected", func(t *testing.T) {
		fm := "name: sample-skill\ndescription: Use when validating.\nmetadata:\n  author: tester\n  extra:\n    nested: true\n"
		dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
		res := validateSkill(dir, "")
		if !hasMessage(res, "Field 'metadata' must be a mapping of string values") {
			t.Fatalf("expected metadata error, got %v", res.Errors)
		}
	})
	t.Run("scalar coercion keeps written form", func(t *testing.T) {
		fm := "name: sample-skill\ndescription: Use when validating.\nmetadata:\n  version: 1.0\n  author: tester\n"
		dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
		res := validateSkill(dir, "")
		if len(res.Errors) != 0 {
			t.Fatalf("expected no errors, got %v", res.Errors)
		}
	})
}

func TestFrontmatterParsing(t *testing.T) {
	t.Run("no frontmatter", func(t *testing.T) {
		dir := writeRaw(t, filepath.Join(t.TempDir(), "sample-skill"), "Just text.\n")
		res := validateSkill(dir, "")
		if !hasMessage(res, "SKILL.md must start with YAML frontmatter (---)") {
			t.Fatalf("expected frontmatter error, got %v", res.Errors)
		}
	})
	t.Run("unclosed", func(t *testing.T) {
		dir := writeRaw(t, filepath.Join(t.TempDir(), "sample-skill"), "---\nname: sample-skill\n")
		res := validateSkill(dir, "")
		if !hasMessage(res, "SKILL.md frontmatter not properly closed with ---") {
			t.Fatalf("expected unclosed error, got %v", res.Errors)
		}
	})
	t.Run("not a mapping", func(t *testing.T) {
		dir := writeRaw(t, filepath.Join(t.TempDir(), "sample-skill"), "---\n- a\n- b\n---\n")
		res := validateSkill(dir, "")
		if !hasMessage(res, "SKILL.md frontmatter must be a YAML mapping") {
			t.Fatalf("expected mapping error, got %v", res.Errors)
		}
	})
	t.Run("invalid yaml", func(t *testing.T) {
		dir := writeRaw(t, filepath.Join(t.TempDir(), "sample-skill"), "---\n: : :\n---\n")
		res := validateSkill(dir, "")
		if !hasMessage(res, "invalid YAML frontmatter") {
			t.Fatalf("expected yaml error, got %v", res.Errors)
		}
	})
	t.Run("body containing --- is body", func(t *testing.T) {
		dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), "name: sample-skill\ndescription: Use when validating.\n", "## Rules\n\n---\n\nNot frontmatter.\n")
		res := validateSkill(dir, "")
		if len(res.Errors) != 0 {
			t.Fatalf("body with --- should parse, got %v", res.Errors)
		}
	})
}

func TestMissingSkillMD(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "empty-skill")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	res := validateSkill(dir, "")
	if !hasMessage(res, "Missing required file: SKILL.md") {
		t.Fatalf("expected missing SKILL.md error, got %v", res.Errors)
	}
}

func TestReadmeListing(t *testing.T) {
	dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), validFrontmatter, "")
	readme := "| [sample-skill](skills/sample-skill/SKILL.md) | desc |\n"

	t.Run("listed", func(t *testing.T) {
		res := validateSkill(dir, readme)
		if hasMessage(res, "README") {
			t.Fatalf("listing should pass, got %v / %v", res.Errors, res.Warnings)
		}
	})
	t.Run("not listed", func(t *testing.T) {
		res := validateSkill(dir, "| [other](skills/other/SKILL.md) | desc |\n")
		if !hasMessage(res, "not listed in README.md Available Skills table") {
			t.Fatalf("expected listing error, got %v", res.Errors)
		}
	})
	t.Run("no readme skips check", func(t *testing.T) {
		res := validateSkill(dir, "")
		if hasMessage(res, "README") {
			t.Fatalf("standalone skill should skip the check, got %v", res.Errors)
		}
	})
}

func TestLintBodyLength(t *testing.T) {
	long := strings.Repeat("line\n", 501)
	dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), validFrontmatter, long)
	res := validateSkill(dir, "")
	if !hasMessage(res, "SKILL.md body is 501 lines; recommended maximum is 500") {
		t.Fatalf("expected body length warning, got %v", res.Warnings)
	}
}

func TestLintTriggerGuidance(t *testing.T) {
	fm := "name: sample-skill\ndescription: Processes files somehow.\n"
	dir := writeSkill(t, filepath.Join(t.TempDir(), "sample-skill"), fm, "")
	res := validateSkill(dir, "")
	if !hasMessage(res, "description does not state when to use the skill") {
		t.Fatalf("expected trigger guidance warning, got %v", res.Warnings)
	}
}

func TestLintRelativeLinks(t *testing.T) {
	dir := t.TempDir()
	skillDir := writeSkill(t, filepath.Join(dir, "sample-skill"),
		validFrontmatter,
		"[good](references/exists.md) [bad](references/missing.md)\n")
	if err := os.MkdirAll(filepath.Join(skillDir, "references"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "references", "exists.md"), []byte("[back](../SKILL.md)\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	res := validateSkill(skillDir, "")
	if !hasMessage(res, "broken relative link 'references/missing.md' in SKILL.md") {
		t.Fatalf("expected broken link warning, got %v", res.Warnings)
	}
	if hasMessage(res, "references/exists.md") || hasMessage(res, "../SKILL.md") {
		t.Fatalf("existing links should not warn, got %v", res.Warnings)
	}
}

func TestDiscoverSkills(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, filepath.Join(root, "skills", "beta"), validFrontmatter, "")
	writeSkill(t, filepath.Join(root, "skills", "alpha"), validFrontmatter, "")
	// Hidden directories and nested SKILL.md inside a skill are ignored.
	writeSkill(t, filepath.Join(root, "skills", ".hidden", "x"), validFrontmatter, "")
	writeSkill(t, filepath.Join(root, "skills", "alpha", "references", "nested"), validFrontmatter, "")

	dirs, err := discoverSkills([]string{root})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		filepath.Join(root, "skills", "alpha"),
		filepath.Join(root, "skills", "beta"),
	}
	if len(dirs) != 2 || dirs[0] != want[0] || dirs[1] != want[1] {
		t.Fatalf("expected %v, got %v", want, dirs)
	}
}

func TestDiscoverErrors(t *testing.T) {
	t.Run("missing path", func(t *testing.T) {
		_, err := discoverSkills([]string{filepath.Join(t.TempDir(), "nope")})
		if err == nil || !strings.Contains(err.Error(), "Path does not exist") {
			t.Fatalf("expected path error, got %v", err)
		}
	})
	t.Run("file path", func(t *testing.T) {
		f := filepath.Join(t.TempDir(), "file.txt")
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := discoverSkills([]string{f})
		if err == nil || !strings.Contains(err.Error(), "Not a directory") {
			t.Fatalf("expected directory error, got %v", err)
		}
	})
	t.Run("no skills", func(t *testing.T) {
		_, err := discoverSkills([]string{t.TempDir()})
		if err == nil || !strings.Contains(err.Error(), "no skills found under") {
			t.Fatalf("expected no skills error, got %v", err)
		}
	})
}

func TestCheckMarketplaceJSON(t *testing.T) {
	dir := t.TempDir()

	t.Run("version field is rejected", func(t *testing.T) {
		p := filepath.Join(dir, "with-version.json")
		content := `{"name":"x","metadata":{"description":"d","version":"1.0.0"},"plugins":[{"name":"p","version":"2.0.0"}]}`
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		problems := checkMarketplaceJSON(p)
		if len(problems) != 2 {
			t.Fatalf("expected 2 version problems, got %v", problems)
		}
		for _, want := range []string{"$.metadata.version", "$.plugins[0].version"} {
			if !hasMessage(&skillResult{Errors: problems}, want) {
				t.Fatalf("expected a problem under %s, got %v", want, problems)
			}
		}
	})

	t.Run("no version field passes", func(t *testing.T) {
		p := filepath.Join(dir, "clean.json")
		content := `{"name":"x","metadata":{"description":"d"},"plugins":[{"name":"p","skills":["./skills/a"]}]}`
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		if problems := checkMarketplaceJSON(p); len(problems) != 0 {
			t.Fatalf("expected no problems, got %v", problems)
		}
	})

	t.Run("invalid json is reported", func(t *testing.T) {
		p := filepath.Join(dir, "broken.json")
		if err := os.WriteFile(p, []byte("{not json"), 0o644); err != nil {
			t.Fatal(err)
		}
		if problems := checkMarketplaceJSON(p); len(problems) != 1 {
			t.Fatalf("expected 1 problem, got %v", problems)
		}
	})
}

func TestReportFailed(t *testing.T) {
	allOK := []*skillResult{{Path: "a", Errors: []string{}, Warnings: []string{}}}
	onlyWarnings := []*skillResult{{Path: "a", Errors: []string{}, Warnings: []string{"w"}}}
	withErrors := []*skillResult{{Path: "a", Errors: []string{"e"}, Warnings: []string{}}}

	if buildReport(allOK).failed(false) || buildReport(allOK).failed(true) {
		t.Fatal("clean report must pass")
	}
	if buildReport(onlyWarnings).failed(false) {
		t.Fatal("warnings must pass without strict")
	}
	if !buildReport(onlyWarnings).failed(true) {
		t.Fatal("warnings must fail with strict")
	}
	if !buildReport(withErrors).failed(false) {
		t.Fatal("errors must fail")
	}
}

func TestNodeToJSON(t *testing.T) {
	var seq yaml.Node
	if err := yaml.Unmarshal([]byte("[a, b]"), &seq); err != nil {
		t.Fatal(err)
	}
	got := nodeToJSON(*seq.Content[0])
	items, ok := got.([]any)
	if !ok || len(items) != 2 {
		t.Fatalf("expected sequence, got %#v", got)
	}
}
