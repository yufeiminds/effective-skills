package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// properties mirrors the JSON output of the official skills-ref
// "read-properties" command: present fields only, metadata values coerced
// to strings (raw scalar text, so YAML floats keep their written form).
type properties struct {
	Name          string            `json:"name,omitempty"`
	Description   string            `json:"description,omitempty"`
	License       string            `json:"license,omitempty"`
	Compatibility string            `json:"compatibility,omitempty"`
	AllowedTools  any               `json:"allowed-tools,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

func newReadPropertiesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "read-properties <skill-dir>",
		Short: "Print a skill's frontmatter properties as JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := resolveArgPath(args[0])
			st, err := os.Stat(dir)
			if os.IsNotExist(err) {
				return failf("Path does not exist: %s", dir)
			}
			if err != nil {
				return failf("failed to access %s: %v", dir, err)
			}
			if !st.IsDir() {
				return failf("Not a directory: %s", dir)
			}

			fields, _, err := readSkillMD(dir)
			if err != nil {
				return failf("%s", err.Error())
			}

			props := properties{}
			if n, ok := scalarField(fields, "name"); ok {
				props.Name = strings.TrimSpace(n.Value)
			}
			if n, ok := scalarField(fields, "description"); ok {
				props.Description = strings.TrimSpace(n.Value)
			}
			if n, ok := scalarField(fields, "license"); ok {
				props.License = strings.TrimSpace(n.Value)
			}
			if n, ok := scalarField(fields, "compatibility"); ok {
				props.Compatibility = strings.TrimSpace(n.Value)
			}
			if n, ok := fields["allowed-tools"]; ok && n.Tag != "!!null" {
				props.AllowedTools = nodeToJSON(n)
			}
			if n, ok := fields["metadata"]; ok && n.Kind == yaml.MappingNode {
				meta := map[string]string{}
				for i := 0; i+1 < len(n.Content); i += 2 {
					if v := n.Content[i+1]; v.Kind == yaml.ScalarNode {
						meta[n.Content[i].Value] = v.Value
					}
				}
				props.Metadata = meta
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(props)
		},
	}
}

// scalarField returns the field's node when it is a non-null scalar.
func scalarField(fields map[string]yaml.Node, key string) (yaml.Node, bool) {
	n, ok := fields[key]
	if !ok || n.Kind != yaml.ScalarNode || n.Tag == "!!null" {
		return yaml.Node{}, false
	}
	return n, true
}

// nodeToJSON converts a YAML node to a JSON-compatible Go value so fields
// like allowed-tools pass through verbatim regardless of YAML type.
func nodeToJSON(node yaml.Node) any {
	switch node.Kind {
	case yaml.ScalarNode:
		return node.Value
	case yaml.SequenceNode:
		items := make([]any, 0, len(node.Content))
		for _, c := range node.Content {
			items = append(items, nodeToJSON(*c))
		}
		return items
	case yaml.MappingNode:
		obj := map[string]any{}
		for i := 0; i+1 < len(node.Content); i += 2 {
			obj[node.Content[i].Value] = nodeToJSON(*node.Content[i+1])
		}
		return obj
	default:
		return nil
	}
}
