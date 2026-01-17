# Plugin Structure

AI Assist Kit uses a canonical plugin format that can be converted to assistant-specific formats.

## Canonical Structure

```
my-plugin/
├── plugin.yaml             # Plugin metadata
├── commands/               # Slash commands
│   ├── build.md
│   └── test.md
├── skills/                 # Reusable skills
│   ├── code-review.md
│   └── documentation.md
└── agents/                 # Specialized agents
    └── security-scanner.yaml
```

## Plugin Metadata

The `plugin.yaml` file defines your plugin:

```yaml
name: my-plugin
version: 1.0.0
description: A helpful plugin for developers
author: Your Name
repository: https://github.com/yourname/my-plugin

commands:
  - build
  - test

skills:
  - code-review
  - documentation

agents:
  - security-scanner
```

## Generated Output

When you generate plugins, each assistant gets its own directory:

```
plugins/
├── claude/
│   ├── .claude-plugin/
│   │   └── plugin.json
│   ├── commands/
│   ├── skills/
│   └── agents/
├── gemini/
│   ├── plugin.yaml
│   ├── commands/
│   └── skills/
├── codex/
│   ├── agents.yaml
│   └── commands/
└── kiro/
    ├── plugin.json
    ├── commands/
    └── agents/
```

## Assistant-Specific Formats

### Claude Code

```
.claude-plugin/
├── plugin.json
├── commands/
│   └── build.md           # Markdown with YAML frontmatter
├── skills/
│   └── review.md
└── agents/
    └── scanner.json
```

### Gemini CLI

```
plugin.yaml                 # Plugin manifest
commands/
└── build.yaml              # YAML format
skills/
└── review.yaml
```

### OpenAI Codex

```
agents.yaml                 # Agent definitions
commands/
└── build.md                # Markdown format
```

### AWS Kiro

```
plugin.json                 # Plugin manifest
commands/
└── build.md
agents/
└── scanner.json            # JSON agent config
```

## File Formats

### Commands

Commands are typically Markdown files with YAML frontmatter:

```markdown
---
name: build
description: Build the project
allowed_tools:
  - Bash
  - Read
---

Build the project using the appropriate build system.
Detect whether this is a Go, Node.js, Python, or other project
and run the correct build command.
```

### Skills

Skills follow a similar format but focus on reusable capabilities:

```markdown
---
name: code-review
description: Review code for best practices
---

Review the provided code for:

- Security vulnerabilities
- Performance issues
- Code style and readability
- Best practices

Provide specific, actionable feedback.
```

### Agents

Agents are defined in YAML or JSON:

```yaml
name: security-scanner
description: Scan code for security vulnerabilities
instructions: |
  You are a security expert. Analyze code for:
  - SQL injection
  - XSS vulnerabilities
  - Authentication issues
  - Sensitive data exposure
model: sonnet
tools:
  - Read
  - Grep
  - Glob
```
