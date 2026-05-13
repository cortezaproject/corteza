package skills

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed library/*.md
var libraryFS embed.FS

type (
	Skill struct {
		Name        string
		Description string
		Triggers    []string
		Body        string
	}

	Registry interface {
		ForTool(toolName string) []*Skill
		All() []*Skill
	}

	registry struct {
		skills    []*Skill
		byTrigger map[string][]*Skill
	}

	frontmatter struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Triggers    []string `yaml:"triggers"`
	}
)

func LoadLibrary() (Registry, error) {
	return load(libraryFS, "library")
}

func (r *registry) ForTool(toolName string) []*Skill {
	return r.byTrigger[toolName]
}

func (r *registry) All() []*Skill {
	return r.skills
}

func load(fsys fs.FS, dir string) (*registry, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("read skill dir: %w", err)
	}

	reg := &registry{byTrigger: make(map[string][]*Skill)}
	seen := make(map[string]bool)

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := dir + "/" + e.Name()
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
		s, err := parseSkill(data)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		if seen[s.Name] {
			return nil, fmt.Errorf("duplicate skill name %q", s.Name)
		}
		seen[s.Name] = true
		reg.skills = append(reg.skills, s)
		for _, t := range s.Triggers {
			reg.byTrigger[t] = append(reg.byTrigger[t], s)
		}
	}

	sort.Slice(reg.skills, func(i, j int) bool { return reg.skills[i].Name < reg.skills[j].Name })
	return reg, nil
}

func parseSkill(data []byte) (*Skill, error) {
	text := string(data)
	if !strings.HasPrefix(text, "---") {
		return nil, fmt.Errorf("missing frontmatter delimiter")
	}
	rest := strings.TrimPrefix(text[3:], "\n")
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return nil, fmt.Errorf("frontmatter not terminated")
	}

	var fm frontmatter
	if err := yaml.Unmarshal([]byte(rest[:end]), &fm); err != nil {
		return nil, fmt.Errorf("frontmatter yaml: %w", err)
	}
	if fm.Name == "" {
		return nil, fmt.Errorf("skill missing name")
	}
	if len(fm.Triggers) == 0 {
		return nil, fmt.Errorf("skill %q has no triggers", fm.Name)
	}

	body := strings.TrimSpace(strings.TrimPrefix(rest[end+4:], "\n"))

	return &Skill{
		Name:        fm.Name,
		Description: fm.Description,
		Triggers:    fm.Triggers,
		Body:        body,
	}, nil
}
