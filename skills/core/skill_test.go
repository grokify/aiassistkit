package core

import "testing"

func TestNewSkill(t *testing.T) {
	skill := NewSkill("version-analysis", "Analyze git history")

	if skill.Name != "version-analysis" {
		t.Errorf("expected Name 'version-analysis', got '%s'", skill.Name)
	}
	if skill.Description != "Analyze git history" {
		t.Errorf("expected Description 'Analyze git history', got '%s'", skill.Description)
	}
}

func TestSkillAddScript(t *testing.T) {
	skill := NewSkill("test", "test")

	skill.AddScript("scripts/analyze.sh")
	skill.AddScript("scripts/validate.sh")

	if len(skill.Scripts) != 2 {
		t.Errorf("expected 2 scripts, got %d", len(skill.Scripts))
	}
}

func TestSkillAddReference(t *testing.T) {
	skill := NewSkill("test", "test")

	skill.AddReference("docs/semver.md")

	if len(skill.References) != 1 {
		t.Errorf("expected 1 reference, got %d", len(skill.References))
	}
}

func TestSkillAddAsset(t *testing.T) {
	skill := NewSkill("test", "test")

	skill.AddAsset("templates/changelog.tmpl")

	if len(skill.Assets) != 1 {
		t.Errorf("expected 1 asset, got %d", len(skill.Assets))
	}
}

func TestSkillAddTrigger(t *testing.T) {
	skill := NewSkill("test", "test")

	skill.AddTrigger("version")
	skill.AddTrigger("semver")

	if len(skill.Triggers) != 2 {
		t.Errorf("expected 2 triggers, got %d", len(skill.Triggers))
	}
}

func TestSkillAddDependency(t *testing.T) {
	skill := NewSkill("test", "test")

	skill.AddDependency("git")
	skill.AddDependency("schangelog")

	if len(skill.Dependencies) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(skill.Dependencies))
	}
}
