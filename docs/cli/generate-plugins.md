# Generate Plugins

The `generate plugins` command creates platform-specific plugins from canonical JSON specifications.

## Synopsis

```bash
assistantkit generate plugins [flags]
```

## Description

This command reads plugin definitions from a canonical spec directory and generates platform-specific plugins for Claude Code, Kiro IDE, and Gemini CLI.

The canonical spec format allows you to define your plugin once and automatically generate outputs for multiple AI coding assistants.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--spec` | `plugins/spec` | Path to canonical spec directory |
| `--output` | `plugins` | Output directory for generated plugins |
| `--platforms` | `claude,kiro` | Comma-separated list of platforms to generate |
| `--config` | *(none)* | Config file path (assistantkit.yaml if exists) |

## Supported Platforms

- **claude**: Claude Code plugins (`.claude-plugin/` directory structure)
- **kiro**: Kiro IDE Powers (POWER.md + mcp.json) or Kiro Agents (agents/*.json)
- **gemini**: Gemini CLI extensions (gemini-extension.json)

## Spec Directory Structure

The canonical spec directory should contain:

```
plugins/spec/
├── plugin.json       # Plugin metadata
├── commands/         # Command definitions (*.json)
│   └── create.json
├── skills/           # Skill definitions (*.json)
│   └── review.json
└── agents/           # Agent definitions (*.json)
    └── release.json
```

### plugin.json

The plugin metadata file defines the plugin name, version, keywords, and MCP server configurations:

```json
{
  "name": "my-plugin",
  "displayName": "My Plugin",
  "version": "1.0.0",
  "description": "A plugin for AI assistants",
  "keywords": ["keyword1", "keyword2"],
  "mcpServers": {
    "my-server": {
      "command": "my-mcp-server",
      "args": []
    }
  }
}
```

### commands/*.json

Command definitions for slash commands:

```json
{
  "name": "create",
  "description": "Create a new resource",
  "arguments": [
    {"name": "name", "description": "Resource name", "required": true}
  ],
  "instructions": "Instructions for the AI assistant...",
  "examples": ["Example usage"]
}
```

### skills/*.json

Skill definitions for reusable capabilities:

```json
{
  "name": "code-review",
  "description": "Reviews code for best practices",
  "instructions": "Instructions for performing code review...",
  "triggers": ["review code", "check code"]
}
```

### agents/*.json

Agent definitions for autonomous tasks:

```json
{
  "name": "release-agent",
  "description": "Manages release process",
  "model": "claude-sonnet-4",
  "systemPrompt": "You are a release management agent...",
  "tools": ["read_file", "write_file", "run_command"]
}
```

## Generated Output

### Claude Code

Generates a `.claude-plugin/` directory structure:

```
plugins/claude/
├── .claude-plugin/
│   └── plugin.json       # Claude plugin manifest
├── commands/
│   └── create.md         # Command instructions
└── skills/
    └── code-review/
        └── SKILL.md      # Skill instructions
```

### Kiro IDE

Generates either a Power or Agents format based on the plugin spec:

**Power format** (when keywords or mcpServers are present):

```
plugins/kiro/
├── POWER.md              # Power description
└── mcp.json              # MCP server configuration
```

**Agents format** (when no keywords/mcpServers):

```
plugins/kiro/
└── agents/
    └── release-agent.json  # Agent definition
```

### Gemini CLI

Generates a Gemini extension:

```
plugins/gemini/
├── gemini-extension.json  # Extension manifest
└── agents/
    └── release-agent.json  # Agent definition
```

## Examples

Generate plugins for all default platforms:

```bash
assistantkit generate plugins
```

Generate only for Claude:

```bash
assistantkit generate plugins --platforms=claude
```

Generate for all platforms with custom directories:

```bash
assistantkit generate plugins --spec=canonical --output=dist --platforms=claude,kiro,gemini
```

## See Also

- [Plugin Structure](../plugins/structure.md) - Learn about plugin components
- [Commands](../plugins/commands.md) - Command definition details
- [Skills](../plugins/skills.md) - Skill definition details
- [Agents](../plugins/agents.md) - Agent definition details
