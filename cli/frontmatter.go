package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// parseError reports a malformed SKILL.md. Message wording follows the
// official skills-ref parser (skills-ref/src/skills_ref/parser.py).
type parseError struct{ msg string }

func (e *parseError) Error() string { return e.msg }

// findSkillMD locates the skill's entry file. SKILL.md takes precedence; the
// lowercase form is accepted as a fallback, mirroring skills-ref.
func findSkillMD(dir string) (string, bool) {
	for _, name := range []string{"SKILL.md", "skill.md"} {
		p := filepath.Join(dir, name)
		if st, err := os.Stat(p); err == nil && st.Mode().IsRegular() {
			return p, true
		}
	}
	return "", false
}

// splitFrontmatter extracts the frontmatter YAML text and the body. It
// mirrors skills-ref: the file must start with "---" and the frontmatter is
// everything up to the next "---"; anything after that is body and may
// itself contain more "---" markers.
func splitFrontmatter(content string) (frontmatterText, body string, err error) {
	if !strings.HasPrefix(content, "---") {
		return "", "", &parseError{"SKILL.md must start with YAML frontmatter (---)"}
	}
	rest := content[len("---"):]
	i := strings.Index(rest, "---")
	if i < 0 {
		return "", "", &parseError{"SKILL.md frontmatter not properly closed with ---"}
	}
	return rest[:i], rest[i+len("---"):], nil
}

// parseFrontmatter decodes the frontmatter into a top-level YAML mapping of
// raw nodes. Keys and values are kept as nodes so the validator can check
// the field whitelist and per-field types.
func parseFrontmatter(text string) (map[string]yaml.Node, error) {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(text), &root); err != nil {
		return nil, &parseError{fmt.Sprintf("invalid YAML frontmatter: %v", err)}
	}
	// Unmarshaling into a Node yields a DocumentNode wrapping the content;
	// an empty frontmatter leaves the node untouched (zero Kind).
	if len(root.Content) == 0 {
		return nil, &parseError{"SKILL.md frontmatter must be a YAML mapping"}
	}
	content := root.Content[0]
	if content.Kind != yaml.MappingNode {
		return nil, &parseError{"SKILL.md frontmatter must be a YAML mapping"}
	}

	fields := make(map[string]yaml.Node, len(content.Content)/2)
	for i := 0; i+1 < len(content.Content); i += 2 {
		key := content.Content[i]
		fields[key.Value] = *content.Content[i+1]
	}
	return fields, nil
}

// readSkillMD reads and parses a skill's SKILL.md into raw frontmatter fields
// and the body. It returns nil fields when the file or frontmatter is broken;
// the errors are appended to the result by the caller.
func readSkillMD(dir string) (fields map[string]yaml.Node, body string, err error) {
	path, ok := findSkillMD(dir)
	if !ok {
		return nil, "", &parseError{"Missing required file: SKILL.md"}
	}
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		return nil, "", &parseError{fmt.Sprintf("failed to read SKILL.md: %v", readErr)}
	}
	fmText, body, splitErr := splitFrontmatter(string(content))
	if splitErr != nil {
		return nil, "", splitErr
	}
	fields, parseErr := parseFrontmatter(fmText)
	if parseErr != nil {
		return nil, "", parseErr
	}
	return fields, body, nil
}
