package core

import "fmt"

// ValidationError indicates the plugin failed validation.
type ValidationError struct {
	PluginDir string
	Message   string
	Missing   []string // Missing required files
}

func (e *ValidationError) Error() string {
	if len(e.Missing) > 0 {
		return fmt.Sprintf("validation failed for %s: missing files: %v", e.PluginDir, e.Missing)
	}
	return fmt.Sprintf("validation failed for %s: %s", e.PluginDir, e.Message)
}

// ForkError indicates a failure to fork the repository.
type ForkError struct {
	Owner string
	Repo  string
	Err   error
}

func (e *ForkError) Error() string {
	return fmt.Sprintf("failed to fork %s/%s: %v", e.Owner, e.Repo, e.Err)
}

func (e *ForkError) Unwrap() error {
	return e.Err
}

// BranchError indicates a failure to create or update a branch.
type BranchError struct {
	Branch string
	Err    error
}

func (e *BranchError) Error() string {
	return fmt.Sprintf("failed to create branch %s: %v", e.Branch, e.Err)
}

func (e *BranchError) Unwrap() error {
	return e.Err
}

// CommitError indicates a failure to create a commit.
type CommitError struct {
	Message string
	Err     error
}

func (e *CommitError) Error() string {
	return fmt.Sprintf("failed to create commit: %v", e.Err)
}

func (e *CommitError) Unwrap() error {
	return e.Err
}

// PRError indicates a failure to create a pull request.
type PRError struct {
	Title string
	Err   error
}

func (e *PRError) Error() string {
	return fmt.Sprintf("failed to create PR '%s': %v", e.Title, e.Err)
}

func (e *PRError) Unwrap() error {
	return e.Err
}

// AuthError indicates an authentication failure.
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("authentication failed: %s", e.Message)
}
