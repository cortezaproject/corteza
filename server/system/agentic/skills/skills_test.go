package skills

import (
	"testing"
)

func TestLoadLibrary(t *testing.T) {
	reg, err := LoadLibrary()
	if err != nil {
		t.Fatalf("load library: %v", err)
	}

	all := reg.All()
	if len(all) == 0 {
		t.Fatal("expected at least one skill")
	}

	found := false
	for _, s := range all {
		if s.Name == "page_building" {
			found = true
		}
	}
	if !found {
		t.Fatal("page_building skill missing")
	}

	hits := reg.ForTool("compose_page_create")
	if len(hits) == 0 {
		t.Fatal("compose_page_create should activate at least one skill")
	}
}

func TestParseSkill(t *testing.T) {
	src := []byte(`---
name: demo
description: A test
triggers:
  - foo_tool
---

Body text here.`)

	s, err := parseSkill(src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if s.Name != "demo" {
		t.Errorf("name = %q, want demo", s.Name)
	}
	if len(s.Triggers) != 1 || s.Triggers[0] != "foo_tool" {
		t.Errorf("triggers = %v", s.Triggers)
	}
	if s.Body != "Body text here." {
		t.Errorf("body = %q", s.Body)
	}
}

func TestParseSkillRejectsMissingName(t *testing.T) {
	src := []byte(`---
triggers:
  - foo
---

body`)
	if _, err := parseSkill(src); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestParseSkillRejectsNoTriggers(t *testing.T) {
	src := []byte(`---
name: x
---

body`)
	if _, err := parseSkill(src); err == nil {
		t.Fatal("expected error for missing triggers")
	}
}
