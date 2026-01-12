package core

import (
	"bytes"
	"fmt"
	"strings"
)

// OrchestrationConfig holds configuration for generating orchestration instructions.
type OrchestrationConfig struct {
	// Version is the target release version (e.g., "v1.2.0").
	Version string

	// AgentSpecsPath is the path to agent spec files (e.g., "validation/specs").
	AgentSpecsPath string

	// IncludeTasks limits generation to specific tasks (empty = all tasks).
	IncludeTasks []string
}

// GenerateOrchestrationMD generates Claude Code orchestration instructions in Markdown.
func (t *Team) GenerateOrchestrationMD(cfg OrchestrationConfig) string {
	var buf bytes.Buffer

	// Header
	buf.WriteString(fmt.Sprintf("# %s Orchestration\n\n", toTitle(t.Name)))
	if t.Description != "" {
		buf.WriteString(fmt.Sprintf("%s\n\n", t.Description))
	}

	if cfg.Version != "" {
		buf.WriteString(fmt.Sprintf("**Target Version:** `%s`\n\n", cfg.Version))
	}

	// Process type info
	buf.WriteString(fmt.Sprintf("**Process:** %s\n", t.Process))
	if t.Manager != "" {
		buf.WriteString(fmt.Sprintf("**Manager:** %s\n", t.Manager))
	}
	buf.WriteString("\n---\n\n")

	// Filter tasks if IncludeTasks is specified
	tasks := t.Tasks
	if len(cfg.IncludeTasks) > 0 {
		includeSet := make(map[string]bool)
		for _, name := range cfg.IncludeTasks {
			includeSet[name] = true
		}
		var filtered []Task
		for _, task := range tasks {
			if includeSet[task.Name] {
				filtered = append(filtered, task)
			}
		}
		tasks = filtered
	}

	// Get parallel groups
	tempTeam := &Team{Tasks: tasks}
	groups, err := tempTeam.ParallelGroups()
	if err != nil {
		// Fall back to sequential if grouping fails
		groups = [][]Task{tasks}
	}

	// Generate instructions for each group
	for i, group := range groups {
		if len(group) == 0 {
			continue
		}

		if len(group) > 1 {
			buf.WriteString(fmt.Sprintf("## Parallel Group %d\n\n", i+1))
			buf.WriteString("These tasks can run concurrently using parallel Task tool calls.\n\n")
		} else {
			buf.WriteString(fmt.Sprintf("## Group %d\n\n", i+1))
		}

		for _, task := range group {
			writeTaskInstructions(&buf, task, cfg)
		}
	}

	// Status report template
	buf.WriteString("---\n\n")
	buf.WriteString("## Expected Status Report\n\n")
	buf.WriteString("After execution, report status in this format:\n\n")
	buf.WriteString("```\n")
	writeStatusReportTemplate(&buf, tasks)
	buf.WriteString("```\n")

	return buf.String()
}

// writeTaskInstructions writes instructions for a single task.
func writeTaskInstructions(buf *bytes.Buffer, task Task, cfg OrchestrationConfig) {
	buf.WriteString(fmt.Sprintf("### Task: %s\n\n", task.Name))

	if task.Description != "" {
		buf.WriteString(fmt.Sprintf("%s\n\n", task.Description))
	}

	// Dependencies
	if len(task.DependsOn) > 0 {
		buf.WriteString(fmt.Sprintf("**Requires:** %s (must be GO)\n\n", strings.Join(task.DependsOn, ", ")))
	}

	// Agent reference
	agentPath := task.Agent + ".md"
	if cfg.AgentSpecsPath != "" {
		agentPath = cfg.AgentSpecsPath + "/" + task.Agent + ".md"
	}

	buf.WriteString("**Instructions:**\n\n")
	buf.WriteString(fmt.Sprintf("Use the Task tool to spawn subagent `%s`:\n\n", task.Agent))
	buf.WriteString("```\n")
	buf.WriteString(fmt.Sprintf("Task tool:\n"))
	buf.WriteString(fmt.Sprintf("  subagent_type: general-purpose\n"))
	buf.WriteString(fmt.Sprintf("  description: \"%s\"\n", task.Description))
	buf.WriteString(fmt.Sprintf("  prompt: |\n"))
	buf.WriteString(fmt.Sprintf("    You are the %s specialist. Read your instructions from %s\n", task.Agent, agentPath))
	buf.WriteString(fmt.Sprintf("    \n"))
	buf.WriteString(fmt.Sprintf("    Execute the following subtasks and report Go/No-Go for each:\n"))
	for _, st := range task.Subtasks {
		buf.WriteString(fmt.Sprintf("    - %s\n", st.Name))
	}
	if cfg.Version != "" {
		buf.WriteString(fmt.Sprintf("    \n"))
		buf.WriteString(fmt.Sprintf("    Target version: %s\n", cfg.Version))
	}
	buf.WriteString("```\n\n")

	// Subtasks checklist
	if len(task.Subtasks) > 0 {
		buf.WriteString("**Subtasks:**\n\n")
		buf.WriteString("| Subtask | Type | Required | Expected |\n")
		buf.WriteString("|---------|------|----------|----------|\n")
		for _, st := range task.Subtasks {
			required := "Yes"
			if !st.Required {
				required = "No"
			}
			expected := st.ExpectedOutput
			if expected == "" {
				expected = "Pass"
			}
			// Truncate long expected outputs
			if len(expected) > 30 {
				expected = expected[:27] + "..."
			}
			buf.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", st.Name, st.Type(), required, expected))
		}
		buf.WriteString("\n")
	}

	// Sign-off criteria
	requiredCount := task.RequiredSubtaskCount()
	buf.WriteString(fmt.Sprintf("**Sign-off:** GO if all %d required subtasks pass. ", requiredCount))
	buf.WriteString("Optional subtasks report WARN on failure.\n\n")
}

// writeStatusReportTemplate writes the expected status report format.
func writeStatusReportTemplate(buf *bytes.Buffer, tasks []Task) {
	// Calculate max widths
	maxTaskLen := 20
	maxSubtaskLen := 18

	for _, task := range tasks {
		if len(task.Name) > maxTaskLen {
			maxTaskLen = len(task.Name)
		}
		for _, st := range task.Subtasks {
			if len(st.Name) > maxSubtaskLen {
				maxSubtaskLen = len(st.Name)
			}
		}
	}

	width := maxTaskLen + maxSubtaskLen + 20
	if width < 60 {
		width = 60
	}

	border := strings.Repeat("═", width)
	buf.WriteString(fmt.Sprintf("╔%s╗\n", border))
	buf.WriteString(fmt.Sprintf("║%s║\n", centerText("TEAM STATUS REPORT", width)))
	buf.WriteString(fmt.Sprintf("╠%s╣\n", border))

	for _, task := range tasks {
		taskLine := fmt.Sprintf(" %s (%s)", task.Name, task.Agent)
		buf.WriteString(fmt.Sprintf("║%-*s║\n", width, taskLine))

		for _, st := range task.Subtasks {
			status := "🟢 GO"
			if !st.Required {
				status = "🟢 GO (optional)"
			}
			stLine := fmt.Sprintf("   %-*s %s", maxSubtaskLen, st.Name, status)
			buf.WriteString(fmt.Sprintf("║%-*s║\n", width, stLine))
		}
		buf.WriteString(fmt.Sprintf("╠%s╣\n", border))
	}

	buf.WriteString(fmt.Sprintf("║%s║\n", centerText("🚀 TEAM: GO 🚀", width)))
	buf.WriteString(fmt.Sprintf("╚%s╝\n", border))
}

// toTitle converts a kebab-case string to Title Case.
func toTitle(s string) string {
	words := strings.Split(s, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

// centerText centers text within a given width.
func centerText(s string, width int) string {
	if len(s) >= width {
		return s
	}
	padding := (width - len(s)) / 2
	return strings.Repeat(" ", padding) + s + strings.Repeat(" ", width-padding-len(s))
}
