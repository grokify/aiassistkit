package requirements

// Registry maps tool names to their requirement definitions.
type Registry map[string]Requirement

// DefaultRegistry contains the default set of known tools.
// Projects can extend this with project-specific requirements.
var DefaultRegistry = Registry{
	"releasekit": {
		Name:    "releasekit",
		Purpose: "language-specific validation (build, test, lint, format)",
		Check:   "releasekit --version",
		Homepage: "https://github.com/grokify/releasekit",
		InstallMethods: []InstallMethod{
			{
				Name:     "go",
				Command:  "go install github.com/grokify/releasekit/cmd/releasekit@latest",
				Requires: []string{"go"},
			},
		},
	},
	"schangelog": {
		Name:    "schangelog",
		Purpose: "changelog generation and validation",
		Check:   "schangelog --version",
		Homepage: "https://github.com/grokify/structured-changelog",
		InstallMethods: []InstallMethod{
			{
				Name:     "go",
				Command:  "go install github.com/grokify/structured-changelog/cmd/schangelog@latest",
				Requires: []string{"go"},
			},
		},
	},
	"sroadmap": {
		Name:    "sroadmap",
		Purpose: "roadmap generation from structured data",
		Check:   "sroadmap --version",
		Homepage: "https://github.com/grokify/structured-roadmap",
		InstallMethods: []InstallMethod{
			{
				Name:     "go",
				Command:  "go install github.com/grokify/structured-roadmap/cmd/sroadmap@latest",
				Requires: []string{"go"},
			},
		},
	},
	"golangci-lint": {
		Name:    "golangci-lint",
		Purpose: "Go linting and static analysis",
		Check:   "golangci-lint --version",
		Homepage: "https://golangci-lint.run",
		InstallMethods: []InstallMethod{
			{
				Name:     "go",
				Command:  "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
				Requires: []string{"go"},
			},
			{
				Name:      "brew",
				Command:   "brew install golangci-lint",
				Requires:  []string{"brew"},
				Platforms: []string{"darwin", "linux"},
			},
			{
				Name:     "curl",
				Command:  "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin",
				Requires: []string{"curl", "go"},
			},
		},
	},
	"go": {
		Name:    "go",
		Purpose: "Go programming language toolchain",
		Check:   "go version",
		Homepage: "https://go.dev",
		InstallMethods: []InstallMethod{
			{
				Name:      "brew",
				Command:   "brew install go",
				Requires:  []string{"brew"},
				Platforms: []string{"darwin"},
			},
			{
				Name:      "apt",
				Command:   "sudo apt-get install -y golang",
				Requires:  []string{"apt-get"},
				Platforms: []string{"linux"},
			},
			{
				Name:     "manual",
				Command:  "# Download from https://go.dev/dl/",
				Requires: []string{},
			},
		},
	},
	"git": {
		Name:    "git",
		Purpose: "version control system",
		Check:   "git --version",
		Homepage: "https://git-scm.com",
		InstallMethods: []InstallMethod{
			{
				Name:      "brew",
				Command:   "brew install git",
				Requires:  []string{"brew"},
				Platforms: []string{"darwin"},
			},
			{
				Name:      "apt",
				Command:   "sudo apt-get install -y git",
				Requires:  []string{"apt-get"},
				Platforms: []string{"linux"},
			},
			{
				Name:      "xcode",
				Command:   "xcode-select --install",
				Requires:  []string{},
				Platforms: []string{"darwin"},
			},
		},
	},
	"gh": {
		Name:    "gh",
		Purpose: "GitHub CLI for PR, issue, and release management",
		Check:   "gh --version",
		Homepage: "https://cli.github.com",
		InstallMethods: []InstallMethod{
			{
				Name:      "brew",
				Command:   "brew install gh",
				Requires:  []string{"brew"},
				Platforms: []string{"darwin", "linux"},
			},
			{
				Name:      "apt",
				Command:   "sudo apt-get install -y gh",
				Requires:  []string{"apt-get"},
				Platforms: []string{"linux"},
			},
			{
				Name:     "go",
				Command:  "go install github.com/cli/cli/v2/cmd/gh@latest",
				Requires: []string{"go"},
			},
		},
	},
	"helm": {
		Name:    "helm",
		Purpose: "Kubernetes package manager",
		Check:   "helm version",
		Homepage: "https://helm.sh",
		InstallMethods: []InstallMethod{
			{
				Name:      "brew",
				Command:   "brew install helm",
				Requires:  []string{"brew"},
				Platforms: []string{"darwin", "linux"},
			},
			{
				Name:     "curl",
				Command:  "curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash",
				Requires: []string{"curl"},
			},
			{
				Name:      "apt",
				Command:   "sudo apt-get install -y helm",
				Requires:  []string{"apt-get"},
				Platforms: []string{"linux"},
			},
		},
	},
	"brew": {
		Name:    "brew",
		Purpose: "Homebrew package manager",
		Check:   "brew --version",
		Homepage: "https://brew.sh",
		InstallMethods: []InstallMethod{
			{
				Name:      "curl",
				Command:   `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`,
				Requires:  []string{"curl"},
				Platforms: []string{"darwin", "linux"},
			},
		},
	},
	"kubectl": {
		Name:    "kubectl",
		Purpose: "Kubernetes command-line tool",
		Check:   "kubectl version --client",
		Homepage: "https://kubernetes.io/docs/tasks/tools/",
		InstallMethods: []InstallMethod{
			{
				Name:      "brew",
				Command:   "brew install kubectl",
				Requires:  []string{"brew"},
				Platforms: []string{"darwin", "linux"},
			},
			{
				Name:      "apt",
				Command:   "sudo apt-get install -y kubectl",
				Requires:  []string{"apt-get"},
				Platforms: []string{"linux"},
			},
			{
				Name:     "gcloud",
				Command:  "gcloud components install kubectl",
				Requires: []string{"gcloud"},
			},
		},
	},
	"docker": {
		Name:    "docker",
		Purpose: "container runtime",
		Check:   "docker --version",
		Homepage: "https://www.docker.com",
		InstallMethods: []InstallMethod{
			{
				Name:      "brew",
				Command:   "brew install --cask docker",
				Requires:  []string{"brew"},
				Platforms: []string{"darwin"},
			},
			{
				Name:      "apt",
				Command:   "sudo apt-get install -y docker.io",
				Requires:  []string{"apt-get"},
				Platforms: []string{"linux"},
			},
		},
	},
	"node": {
		Name:    "node",
		Purpose: "Node.js JavaScript runtime",
		Check:   "node --version",
		Homepage: "https://nodejs.org",
		InstallMethods: []InstallMethod{
			{
				Name:      "brew",
				Command:   "brew install node",
				Requires:  []string{"brew"},
				Platforms: []string{"darwin", "linux"},
			},
			{
				Name:      "apt",
				Command:   "sudo apt-get install -y nodejs",
				Requires:  []string{"apt-get"},
				Platforms: []string{"linux"},
			},
			{
				Name:     "nvm",
				Command:  "nvm install --lts",
				Requires: []string{"nvm"},
			},
		},
	},
	"npm": {
		Name:    "npm",
		Purpose: "Node.js package manager",
		Check:   "npm --version",
		Homepage: "https://www.npmjs.com",
		InstallMethods: []InstallMethod{
			{
				Name:     "node",
				Command:  "# npm is included with Node.js",
				Requires: []string{"node"},
			},
		},
	},
	"pnpm": {
		Name:    "pnpm",
		Purpose: "fast, disk space efficient package manager",
		Check:   "pnpm --version",
		Homepage: "https://pnpm.io",
		InstallMethods: []InstallMethod{
			{
				Name:     "npm",
				Command:  "npm install -g pnpm",
				Requires: []string{"npm"},
			},
			{
				Name:      "brew",
				Command:   "brew install pnpm",
				Requires:  []string{"brew"},
				Platforms: []string{"darwin", "linux"},
			},
			{
				Name:     "curl",
				Command:  "curl -fsSL https://get.pnpm.io/install.sh | sh -",
				Requires: []string{"curl"},
			},
		},
	},
	"govulncheck": {
		Name:    "govulncheck",
		Purpose: "Go vulnerability scanner",
		Check:   "govulncheck --version",
		Homepage: "https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck",
		InstallMethods: []InstallMethod{
			{
				Name:     "go",
				Command:  "go install golang.org/x/vuln/cmd/govulncheck@latest",
				Requires: []string{"go"},
			},
		},
	},
	"trivy": {
		Name:    "trivy",
		Purpose: "vulnerability scanner for containers and filesystems",
		Check:   "trivy --version",
		Homepage: "https://trivy.dev",
		InstallMethods: []InstallMethod{
			{
				Name:      "brew",
				Command:   "brew install trivy",
				Requires:  []string{"brew"},
				Platforms: []string{"darwin", "linux"},
			},
			{
				Name:      "apt",
				Command:   "sudo apt-get install -y trivy",
				Requires:  []string{"apt-get"},
				Platforms: []string{"linux"},
			},
			{
				Name:     "curl",
				Command:  "curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin",
				Requires: []string{"curl"},
			},
		},
	},
}

// Get returns the requirement for a tool name, or nil if not found.
func (r Registry) Get(name string) *Requirement {
	if req, ok := r[name]; ok {
		return &req
	}
	return nil
}

// Merge combines two registries, with the overlay taking precedence.
func (r Registry) Merge(overlay Registry) Registry {
	merged := make(Registry)
	for k, v := range r {
		merged[k] = v
	}
	for k, v := range overlay {
		merged[k] = v
	}
	return merged
}

// Names returns all tool names in the registry.
func (r Registry) Names() []string {
	names := make([]string, 0, len(r))
	for name := range r {
		names = append(names, name)
	}
	return names
}
