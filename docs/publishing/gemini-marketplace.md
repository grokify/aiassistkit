# Gemini CLI Marketplace

Gemini CLI has an [Extensions Gallery](https://geminicli.com/extensions/browse/) where users can discover and install extensions.

## Extension Structure

Gemini CLI extensions use a specific structure with `gemini-extension.json` as the manifest:

```
my-extension/
├── gemini-extension.json    # Required: Extension manifest
├── GEMINI.md                # Optional: Context for the model
├── commands/                # Optional: Custom commands (TOML files)
│   └── build.toml
└── hooks/                   # Optional: Lifecycle hooks
    └── hooks.json
```

## gemini-extension.json

The manifest file defines your extension:

```json
{
  "name": "my-extension",
  "version": "1.0.0",
  "mcpServers": {
    "my-server": {
      "command": "node",
      "args": ["${extensionPath}${/}dist${/}server.js"],
      "cwd": "${extensionPath}"
    }
  },
  "contextFileName": "GEMINI.md",
  "excludeTools": []
}
```

| Field | Description |
|-------|-------------|
| `name` | Unique extension name (lowercase, dashes) |
| `version` | Semantic version |
| `mcpServers` | MCP server definitions |
| `contextFileName` | Context file (defaults to GEMINI.md) |
| `excludeTools` | Tools to disable |
| `settings` | User-configurable settings |

## Custom Commands

Commands are TOML files in the `commands/` directory:

```toml
# commands/build.toml
prompt = """
Build the project using the appropriate build system.
Detect the project type and run the correct command.
"""
```

Commands can include arguments and shell interpolation:

```toml
# commands/fs/grep-code.toml
prompt = """
Please summarize the findings for the pattern `{{args}}`.

Search Results:
!{grep -r {{args}} .}
"""
```

## Releasing Extensions

There are two primary methods for releasing extensions:

### Method 1: Git Repository (Recommended)

The simplest approach - just create a public GitHub repository:

```bash
# Users install directly from GitHub
gemini extensions install https://github.com/yourname/my-extension
```

Users can specify a branch, tag, or commit:

```bash
gemini extensions install https://github.com/yourname/my-extension --ref=v1.0.0
gemini extensions install https://github.com/yourname/my-extension --ref=stable
```

#### Release Channels

Manage multiple release channels using branches:

| Branch | Purpose | Install Command |
|--------|---------|-----------------|
| `main` | Stable releases | `gemini extensions install <url>` |
| `preview` | Preview/beta | `gemini extensions install <url> --ref=preview` |
| `dev` | Development | `gemini extensions install <url> --ref=dev` |

### Method 2: GitHub Releases

For faster initial installs, use [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases) with archive files.

#### Archive Naming Convention

Platform-specific archives must follow this naming:

| Pattern | Example |
|---------|---------|
| `{platform}.{arch}.{name}.{ext}` | `darwin.arm64.my-tool.tar.gz` |
| `{platform}.{name}.{ext}` | `linux.my-tool.tar.gz` |
| Generic (single asset) | `my-tool.tar.gz` |

Supported platforms: `darwin`, `linux`, `win32`
Supported architectures: `x64`, `arm64`

#### GitHub Actions Workflow

```yaml
name: Release Extension

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '20'

      - name: Install dependencies
        run: npm ci

      - name: Build extension
        run: npm run build

      - name: Create release assets
        run: |
          npm run package -- --platform=darwin --arch=arm64
          npm run package -- --platform=linux --arch=x64
          npm run package -- --platform=win32 --arch=x64

      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            release/darwin.arm64.my-tool.tar.gz
            release/linux.x64.my-tool.tar.gz
            release/win32.x64.my-tool.zip
```

## Extensions Gallery

The [Extensions Gallery](https://geminicli.com/extensions/browse/) lists available extensions sourced from public GitHub repositories.

### How Google Discovers Extensions

Google scans GitHub for repositories containing a **`gemini-extension.json`** file. This manifest file is the primary identifier that marks a repository as a valid Gemini CLI extension.

```
GitHub Repositories
        ↓
Scan for gemini-extension.json files
        ↓
Validate JSON structure (name, version)
        ↓
Index to geminicli.com/extensions
        ↓
Rank by GitHub stars
```

!!! note "Topics Are Not Required"
    Only ~8 repositories use the `gemini-cli-extension` GitHub topic, but the gallery has 286+ extensions. The `gemini-extension.json` file is what triggers discovery, not GitHub topics.

### How to Get Listed

The gallery **automatically indexes extensions from GitHub** - there's no manual submission process:

1. **Create a public GitHub repository**
2. **Add a valid `gemini-extension.json`** at the repository root
3. **Wait for automatic indexing** - takes a few days
4. **If not discovered**, open a [GitHub issue](https://github.com/google-gemini/gemini-cli/issues) and assign it to the team

A Google collaborator confirmed: *"We index based on what's on GitHub, you should get picked up in the near future"* ([source](https://github.com/google-gemini/gemini-cli/discussions/10718))

### Minimum Required Fields

Your `gemini-extension.json` must have at least:

```json
{
  "name": "my-extension",
  "version": "1.0.0"
}
```

| Field | Requirements |
|-------|--------------|
| `name` | Lowercase letters, numbers, and dashes only |
| `version` | Semver format recommended (e.g., `1.0.0`) |

### Gallery Features

- Extensions are **ranked by GitHub stars**
- Includes community, partner, and Google-built extensions
- Google does not vet or endorse third-party extensions

!!! tip "Improve Discoverability"
    - Use a clear, descriptive repository name
    - Add the `gemini-cli-extension` GitHub topic (optional but helpful)
    - Write a comprehensive README
    - Gain GitHub stars for higher ranking

## Local Development

### Create from Template

```bash
# Available templates: context, custom-commands, exclude-tools, mcp-server
gemini extensions new my-extension mcp-server
```

### Link for Development

```bash
cd my-extension
npm install
npm run build
gemini extensions link .
```

This creates a symlink so changes are reflected immediately without reinstalling.

### Test Your Extension

```bash
# Restart Gemini CLI to load changes
gemini

# List installed extensions
/extensions list
```

## Extension Management Commands

| Command | Description |
|---------|-------------|
| `gemini extensions install <url>` | Install from GitHub |
| `gemini extensions uninstall <name>` | Remove extension |
| `gemini extensions update <name>` | Update to latest |
| `gemini extensions update --all` | Update all extensions |
| `gemini extensions enable <name>` | Enable extension |
| `gemini extensions disable <name>` | Disable extension |
| `gemini extensions link <path>` | Link local extension |
| `gemini extensions list` | List installed |

## User Settings

Extensions can define user-configurable settings:

```json
{
  "name": "my-api-extension",
  "version": "1.0.0",
  "settings": [
    {
      "name": "API Key",
      "description": "Your API key for the service.",
      "envVar": "MY_API_KEY",
      "sensitive": true
    }
  ]
}
```

Users are prompted for settings on install. Sensitive values are stored in the system keychain.

## Hooks

Extensions can intercept CLI lifecycle events:

```
my-extension/
└── hooks/
    └── hooks.json
```

```json
{
  "hooks": {
    "before_agent": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "node ${extensionPath}/scripts/setup.js",
            "name": "Extension Setup"
          }
        ]
      }
    ]
  }
}
```

## Variables

Use these variables in `gemini-extension.json` and `hooks.json`:

| Variable | Description |
|----------|-------------|
| `${extensionPath}` | Full path to extension directory |
| `${workspacePath}` | Current workspace path |
| `${/}` or `${pathSeparator}` | OS-specific path separator |
| `${process.execPath}` | Path to Node.js binary |

## Best Practices

1. **Use semantic versioning** - Follow semver for version numbers
2. **Document your extension** - Include a comprehensive README
3. **Test on multiple platforms** - Verify cross-platform compatibility
4. **Handle errors gracefully** - Provide helpful error messages
5. **Keep dependencies minimal** - Reduce installation size and conflicts

## Sources

- [Gemini CLI Extensions Documentation](https://github.com/google-gemini/gemini-cli/blob/main/docs/extensions/index.md)
- [Extension Releasing Guide](https://github.com/google-gemini/gemini-cli/blob/main/docs/extensions/extension-releasing.md)
- [Getting Started with Extensions](https://github.com/google-gemini/gemini-cli/blob/main/docs/extensions/getting-started-extensions.md)
- [Extensions Gallery](https://geminicli.com/extensions/browse/)
