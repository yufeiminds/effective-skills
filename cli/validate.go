package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
	"gopkg.in/yaml.v3"
)

// Field limits from the agentskills.io specification, as enforced by the
// official skills-ref validator.
const (
	maxNameLength          = 64
	maxDescriptionLength   = 1024
	maxCompatibilityLength = 500
)

// maxBodyLines is the recommended SKILL.md body length from the progressive
// disclosure guidance (https://agentskills.io/specification). Exceeding it is
// a lint warning, not a spec error.
const maxBodyLines = 500

// allowedFields is the closed set of frontmatter fields defined by the spec.
// skills-ref rejects anything else; unknown extensions (e.g. Claude
// Code-only fields) would also break uploads to claude.ai and the Skills API.
var allowedFields = map[string]bool{
	"name":          true,
	"description":   true,
	"license":       true,
	"allowed-tools": true,
	"compatibility": true,
	"metadata":      true,
}

var sortedAllowedFields = []string{"allowed-tools", "compatibility", "description", "license", "metadata", "name"}

type skillResult struct {
	Path     string   `json:"path"`
	Name     string   `json:"name,omitempty"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

func (r *skillResult) addError(format string, args ...any) {
	r.Errors = append(r.Errors, fmt.Sprintf(format, args...))
}

func (r *skillResult) addWarning(format string, args ...any) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(format, args...))
}

func (r *skillResult) ok() bool { return len(r.Errors) == 0 }

// validateSkillDir validates one skill directory, auto-resolving the nearest
// README.md (walking up) for the repository-listing check.
func validateSkillDir(dir string) *skillResult {
	return validateSkill(dir, resolveReadme(dir))
}

// validateSkill runs all checks against a skill directory. readmeContent is
// the nearest README.md above the skill, or empty when validating a
// standalone directory outside a repository (the listing check is skipped).
func validateSkill(dir string, readmeContent string) *skillResult {
	res := &skillResult{
		Path:     displayPath(dir),
		Errors:   []string{},
		Warnings: []string{},
	}

	fields, body, err := readSkillMD(dir)
	if err != nil {
		res.addError("%s", err.Error())
		return res
	}

	checkFieldWhitelist(res, fields)
	checkName(res, fields, dir)
	checkDescription(res, fields)
	checkCompatibility(res, fields)
	checkMetadata(res, fields)
	checkReadmeListing(res, dir, readmeContent)
	lintBodyLength(res, body)
	lintTriggerGuidance(res, fields)
	lintRelativeLinks(res, dir)

	return res
}

// checkFieldWhitelist rejects frontmatter keys outside the spec's closed set.
func checkFieldWhitelist(res *skillResult, fields map[string]yaml.Node) {
	var unexpected []string
	for key := range fields {
		if !allowedFields[key] {
			unexpected = append(unexpected, key)
		}
	}
	if len(unexpected) == 0 {
		return
	}
	sort.Strings(unexpected)
	res.addError("Unexpected fields in frontmatter: [%s]. Only [%s] are allowed.",
		strings.Join(unexpected, ", "), strings.Join(sortedAllowedFields, ", "))
}

// checkName validates the required name field. Rules and message wording
// follow skills-ref's validator (validator.py):
// non-empty, NFKC length <= 64, lowercase, no leading/trailing/consecutive
// hyphens, letters/digits/hyphens only, and matching the directory name.
func checkName(res *skillResult, fields map[string]yaml.Node, dir string) {
	node, ok := fields["name"]
	if !ok {
		res.addError("Missing required field in frontmatter: name")
		return
	}
	if node.Tag == "!!null" || !isStringScalar(node) {
		res.addError("Field 'name' must be a non-empty string")
		return
	}
	name := strings.TrimSpace(node.Value)
	if name == "" {
		res.addError("Field 'name' must be a non-empty string")
		return
	}

	normalized := norm.NFKC.String(name)
	res.Name = normalized

	if r := []rune(normalized); len(r) > maxNameLength {
		res.addError("Skill name '%s' exceeds %d character limit (%d chars)", name, maxNameLength, len(r))
	}
	if !isLowercaseName(normalized) {
		res.addError("Skill name '%s' must be lowercase", name)
	}
	if strings.HasPrefix(normalized, "-") || strings.HasSuffix(normalized, "-") {
		res.addError("Skill name cannot start or end with a hyphen")
	}
	if strings.Contains(normalized, "--") {
		res.addError("Skill name cannot contain consecutive hyphens")
	}
	for _, r := range normalized {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
			res.addError("Skill name '%s' contains invalid characters. Only letters, digits, and hyphens are allowed.", name)
			break
		}
	}

	dirBase := norm.NFKC.String(filepath.Base(dir))
	if normalized != "" && dirBase != normalized {
		res.addError("Directory name '%s' must match skill name '%s'", filepath.Base(dir), normalized)
	}
}

// isLowercaseName mirrors Python's str.islower(): at least one cased
// character, and every cased character lowercase. Digits and hyphens are
// uncased; CJK ideographs are letters without case, which also mirrors the
// reference implementation's behavior (such names fail the check).
func isLowercaseName(s string) bool {
	hasCased := false
	for _, r := range s {
		if unicode.IsUpper(r) || unicode.IsTitle(r) {
			return false
		}
		if unicode.IsLower(r) {
			hasCased = true
		}
	}
	return hasCased
}

func isStringScalar(node yaml.Node) bool {
	return node.Kind == yaml.ScalarNode
}

// checkDescription validates the required description field.
func checkDescription(res *skillResult, fields map[string]yaml.Node) {
	node, ok := fields["description"]
	if !ok {
		res.addError("Missing required field in frontmatter: description")
		return
	}
	if node.Tag == "!!null" || !isStringScalar(node) {
		res.addError("Field 'description' must be a non-empty string")
		return
	}
	desc := strings.TrimSpace(node.Value)
	if desc == "" {
		res.addError("Field 'description' must be a non-empty string")
		return
	}
	if n := len([]rune(desc)); n > maxDescriptionLength {
		res.addError("Description exceeds %d character limit (%d chars)", maxDescriptionLength, n)
	}
}

// checkCompatibility validates the optional compatibility field: it must be
// a plain string (YAML null counts as missing, mirroring skills-ref where
// None is not a string) of at most 500 characters.
func checkCompatibility(res *skillResult, fields map[string]yaml.Node) {
	node, ok := fields["compatibility"]
	if !ok || node.Tag == "!!null" {
		return
	}
	if !isStringScalar(node) {
		res.addError("Field 'compatibility' must be a string")
		return
	}
	if n := len([]rune(node.Value)); n > maxCompatibilityLength {
		res.addError("Compatibility exceeds %d character limit (%d chars)", maxCompatibilityLength, n)
	}
}

// checkMetadata validates the optional metadata mapping. The spec requires a
// string-to-string map. skills-ref coerces scalars with str() (a YAML float
// 1.0 becomes "1.0"); this implementation keeps the raw scalar text, which
// preserves the same values, but rejects nested mappings/sequences instead
// of stringifying them, since a repr like "{'a': 1}" in a published skill is
// lossy and confusing.
//
// Repository convention (carried over from the previous Python validator):
// missing metadata.author / metadata.version are warnings.
func checkMetadata(res *skillResult, fields map[string]yaml.Node) {
	node, ok := fields["metadata"]
	if !ok || node.Tag == "!!null" {
		res.addWarning("no metadata.author field")
		res.addWarning("no metadata.version field")
		return
	}
	if node.Kind != yaml.MappingNode {
		res.addError("Field 'metadata' must be a mapping of string values")
		return
	}
	hasAuthor, hasVersion := false, false
	for i := 0; i+1 < len(node.Content); i += 2 {
		key, val := node.Content[i], node.Content[i+1]
		if val.Kind != yaml.ScalarNode || val.Tag == "!!null" {
			res.addError("Field 'metadata' must be a mapping of string values: value of '%s' is not a scalar", key.Value)
			return
		}
		switch key.Value {
		case "author":
			hasAuthor = true
		case "version":
			hasVersion = true
		}
	}
	if !hasAuthor {
		res.addWarning("no metadata.author field")
	}
	if !hasVersion {
		res.addWarning("no metadata.version field")
	}
}

// checkReadmeListing enforces this repository's convention that every skill
// is listed in the root README's Available Skills table (error, matching the
// previous Python validator). Skipped for standalone directories with no
// README above them.
func checkReadmeListing(res *skillResult, dir string, readmeContent string) {
	if readmeContent == "" {
		return
	}
	name := filepath.Base(dir)
	if !strings.Contains(readmeContent, fmt.Sprintf("(skills/%s/SKILL.md)", name)) {
		res.addError("not listed in README.md Available Skills table")
	}
}

// resolveReadme finds the nearest README.md walking up from dir, at most six
// levels, and returns its content; empty when none exists.
func resolveReadme(dir string) string {
	path, ok := findUpFile(dir, "README.md", 6)
	if !ok {
		return ""
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(content)
}

// findUpFile searches for name in start and up to maxLevels ancestor
// directories, returning the first match's path.
func findUpFile(start, name string, maxLevels int) (string, bool) {
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", false
	}
	for i := 0; i <= maxLevels; i++ {
		p := filepath.Join(abs, name)
		if st, err := os.Stat(p); err == nil && st.Mode().IsRegular() {
			return p, true
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return "", false
		}
		abs = parent
	}
	return "", false
}

// checkMarketplaceJSON enforces the repository's version-drift rule on the
// Claude Code marketplace manifest: it must not carry any version field. A
// marketplace-level version drifts the moment any skill changes, and nothing
// keeps it honest; per-skill versions live in each SKILL.md frontmatter
// (metadata.version), where the changelog and this validator can see them.
// (Same reasoning as mattpocock/skills, whose marketplace.json omits the
// field for exactly this reason.)
func checkMarketplaceJSON(path string) []string {
	content, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("failed to read %s: %v", displayPath(path), err)}
	}
	var doc any
	if err := json.Unmarshal(content, &doc); err != nil {
		return []string{fmt.Sprintf("invalid JSON in %s: %v", displayPath(path), err)}
	}

	var problems []string
	var walk func(prefix string, node any)
	walk = func(prefix string, node any) {
		switch v := node.(type) {
		case map[string]any:
			for k, val := range v {
				if k == "version" {
					problems = append(problems, fmt.Sprintf("%s.%s: marketplace.json must not carry a version field; skill versions live in each SKILL.md frontmatter (metadata.version)", prefix, k))
					continue
				}
				walk(prefix+"."+k, val)
			}
		case []any:
			for i, val := range v {
				walk(fmt.Sprintf("%s[%d]", prefix, i), val)
			}
		}
	}
	walk("$", doc)
	sort.Strings(problems)
	return problems
}

// lintBodyLength warns when the SKILL.md body exceeds the recommended line
// budget. Every line of the body is loaded whenever the skill triggers, so
// long bodies are a persistent context cost.
func lintBodyLength(res *skillResult, body string) {
	trimmed := strings.TrimSpace(body)
	if trimmed == "" {
		return
	}
	if lines := strings.Count(trimmed, "\n") + 1; lines > maxBodyLines {
		res.addWarning("SKILL.md body is %d lines; recommended maximum is %d (progressive disclosure)", lines, maxBodyLines)
	}
}

// triggerWordRe matches the trigger-guidance vocabulary recommended by the
// agentskills.io authoring guides ("Use this skill when...", etc.).
var triggerWordRe = regexp.MustCompile(`(?i)\b(use|used|using|when|whenever|trigger|triggers|applies|apply)\b`)

// lintTriggerGuidance warns when the description never states when the skill
// applies. The description is the agent's only trigger signal at startup.
func lintTriggerGuidance(res *skillResult, fields map[string]yaml.Node) {
	node, ok := fields["description"]
	if !ok || node.Tag == "!!null" || !isStringScalar(node) {
		return
	}
	if !triggerWordRe.MatchString(node.Value) {
		res.addWarning("description does not state when to use the skill; add trigger guidance (see https://agentskills.io/skill-creation/best-practices.md)")
	}
}

// markdownLinkRe extracts inline markdown links; targets are the second
// group. Reference-style links and images with empty targets are ignored.
var markdownLinkRe = regexp.MustCompile(`\[[^\]]*\]\(([^)\s]+)\)`)

// lintRelativeLinks warns about broken relative markdown links in SKILL.md
// and in one level of references/*.md files. Absolute URLs, anchors, and
// other schemes are skipped.
func lintRelativeLinks(res *skillResult, dir string) {
	files := []string{}
	if path, ok := findSkillMD(dir); ok {
		files = append(files, path)
	}
	refDir := filepath.Join(dir, "references")
	if entries, err := os.ReadDir(refDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(strings.ToLower(e.Name()), ".md") {
				files = append(files, filepath.Join(refDir, e.Name()))
			}
		}
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		for _, m := range markdownLinkRe.FindAllStringSubmatch(string(content), -1) {
			target := strings.Trim(m[1], "<>")
			if !isRelativeFileTarget(target) {
				continue
			}
			target = strings.Split(target, "#")[0]
			if target == "" {
				continue
			}
			resolved := filepath.Join(filepath.Dir(file), target)
			if _, err := os.Stat(resolved); err != nil {
				display := file
				if rel, err := filepath.Rel(dir, file); err == nil {
					display = rel
				}
				res.addWarning("broken relative link '%s' in %s", filepath.ToSlash(target), filepath.ToSlash(display))
			}
		}
	}
}

func isRelativeFileTarget(target string) bool {
	lower := strings.ToLower(target)
	for _, prefix := range []string{"http://", "https://", "mailto:", "ftp://", "tel:", "data:"} {
		if strings.HasPrefix(lower, prefix) {
			return false
		}
	}
	return !strings.HasPrefix(target, "/")
}

// discoverSkills resolves the paths to validate. Each argument may be a
// repository root, a skills directory, or a single skill directory; each is
// searched recursively (bounded depth) for directories containing SKILL.md.
// With no arguments, the nearest skills/ directory above the working
// directory is used, so the tool runs from anywhere inside the repository.
func discoverSkills(args []string) ([]string, error) {
	if len(args) == 0 {
		args = []string{defaultSkillsDir()}
	}

	var skills []string
	seen := map[string]bool{}
	for _, arg := range args {
		arg = resolveArgPath(arg)
		st, err := os.Stat(arg)
		if os.IsNotExist(err) {
			return nil, failf("Path does not exist: %s", arg)
		}
		if err != nil {
			return nil, failf("failed to access %s: %v", arg, err)
		}
		if !st.IsDir() {
			return nil, failf("Not a directory: %s", arg)
		}

		found := 0
		walkErr := collectSkills(arg, func(dir string) {
			if !seen[dir] {
				seen[dir] = true
				skills = append(skills, dir)
				found++
			}
		})
		if walkErr != nil {
			return nil, failf("failed to scan %s: %v", arg, walkErr)
		}
		if found == 0 {
			return nil, failf("no skills found under %s", arg)
		}
	}
	sort.Strings(skills)
	return skills, nil
}

// defaultSkillsDir finds the nearest skills/ directory at or above the
// working directory (at most six levels up), falling back to ".".
func defaultSkillsDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for i := 0; i <= 6; i++ {
		candidate := filepath.Join(dir, "skills")
		if st, err := os.Stat(candidate); err == nil && st.IsDir() {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "."
}

// resolveArgPath resolves a path argument for commands that may run with a
// working directory different from the repository root (e.g. "go -C cli
// run ."): the argument is tried against the working directory first, then
// against each ancestor directory. The original argument is returned when
// nothing matches, so error messages stay accurate.
func resolveArgPath(arg string) string {
	if filepath.IsAbs(arg) {
		return arg
	}
	if _, err := os.Stat(arg); err == nil {
		return arg
	}
	wd, err := os.Getwd()
	if err != nil {
		return arg
	}
	for i := 0; i <= 6; i++ {
		candidate := filepath.Join(wd, arg)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	return arg
}

// displayPath renders a path relative to the working directory when possible
// (including "../" forms, since the tool commonly runs from the repository's
// cli/ directory), falling back to the absolute path with forward slashes.
func displayPath(p string) string {
	if abs, err := filepath.Abs(p); err == nil {
		if cwd, err := os.Getwd(); err == nil {
			if rel, err := filepath.Rel(cwd, abs); err == nil {
				return filepath.ToSlash(rel)
			}
		}
	}
	return filepath.ToSlash(filepath.Clean(p))
}

// collectSkills walks root looking for skill directories: any directory
// containing a SKILL.md is reported to fn and not descended into further
// (its contents are that skill's resources, not more skills). The walk skips
// hidden directories and descends at most three levels below root.
func collectSkills(root string, fn func(dir string)) error {
	maxDepth := 3
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if path != root && strings.HasPrefix(d.Name(), ".") {
			return fs.SkipDir
		}
		if _, ok := findSkillMD(path); ok {
			fn(path)
			return fs.SkipDir
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}
		depth := 0
		if rel != "." {
			depth = strings.Count(rel, string(filepath.Separator)) + 1
		}
		if depth >= maxDepth {
			return fs.SkipDir
		}
		return nil
	})
}
