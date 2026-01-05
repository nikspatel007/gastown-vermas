// Package agent provides runtime abstraction for executing agent prompts
// across different CLI tools (Claude, Codex, OpenCode).
package agent

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"
)

// Runtime represents an agent execution environment.
// Different runtimes (Claude, Codex, OpenCode) implement this interface
// to provide a consistent way to execute prompts.
type Runtime interface {
	// Name returns the runtime identifier (e.g., "claude", "codex", "opencode")
	Name() string

	// Execute runs a prompt in the given working directory and returns the output.
	Execute(ctx context.Context, prompt string, workdir string) (string, error)

	// Available returns true if this runtime is installed and usable.
	Available() bool
}

// ClaudeRuntime executes prompts using the Claude CLI.
// This is the default runtime for all Gas Town agents.
type ClaudeRuntime struct{}

// Name returns "claude".
func (r *ClaudeRuntime) Name() string { return "claude" }

// Execute runs a prompt using the claude CLI with the -p flag.
func (r *ClaudeRuntime) Execute(ctx context.Context, prompt, workdir string) (string, error) {
	cmd := exec.CommandContext(ctx, "claude", "-p", prompt) //nolint:gosec // G204: claude is a trusted CLI tool
	cmd.Dir = workdir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("claude execution failed: %s", stderr.String())
		}
		return "", fmt.Errorf("claude execution failed: %w", err)
	}

	return stdout.String(), nil
}

// Available returns true if the claude CLI is installed.
func (r *ClaudeRuntime) Available() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}

// CodexRuntime executes prompts using the OpenAI Codex CLI.
// Codex is used for independent verification by an alternate model.
type CodexRuntime struct{}

// Name returns "codex".
func (r *CodexRuntime) Name() string { return "codex" }

// Execute runs a prompt using the codex CLI with the -q (quiet) flag.
func (r *CodexRuntime) Execute(ctx context.Context, prompt, workdir string) (string, error) {
	cmd := exec.CommandContext(ctx, "codex", "-q", prompt) //nolint:gosec // G204: codex is a trusted CLI tool
	cmd.Dir = workdir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("codex execution failed: %s", stderr.String())
		}
		return "", fmt.Errorf("codex execution failed: %w", err)
	}

	return stdout.String(), nil
}

// Available returns true if the codex CLI is installed.
func (r *CodexRuntime) Available() bool {
	_, err := exec.LookPath("codex")
	return err == nil
}

// OpenCodeRuntime executes prompts using the OpenCode CLI.
// OpenCode is an open-source alternative for local/self-hosted verification.
type OpenCodeRuntime struct{}

// Name returns "opencode".
func (r *OpenCodeRuntime) Name() string { return "opencode" }

// Execute runs a prompt using the opencode CLI with the -p flag.
func (r *OpenCodeRuntime) Execute(ctx context.Context, prompt, workdir string) (string, error) {
	cmd := exec.CommandContext(ctx, "opencode", "-p", prompt) //nolint:gosec // G204: opencode is a trusted CLI tool
	cmd.Dir = workdir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return "", fmt.Errorf("opencode execution failed: %s", stderr.String())
		}
		return "", fmt.Errorf("opencode execution failed: %w", err)
	}

	return stdout.String(), nil
}

// Available returns true if the opencode CLI is installed.
func (r *OpenCodeRuntime) Available() bool {
	_, err := exec.LookPath("opencode")
	return err == nil
}

// RuntimeRegistry manages available runtimes and provides role-based selection.
type RuntimeRegistry struct {
	mu       sync.RWMutex
	runtimes map[string]Runtime
}

// NewRuntimeRegistry creates a registry and discovers available runtimes.
func NewRuntimeRegistry() *RuntimeRegistry {
	r := &RuntimeRegistry{runtimes: make(map[string]Runtime)}

	// Register all runtimes
	allRuntimes := []Runtime{
		&ClaudeRuntime{},
		&CodexRuntime{},
		&OpenCodeRuntime{},
	}

	for _, rt := range allRuntimes {
		if rt.Available() {
			r.runtimes[rt.Name()] = rt
		}
	}

	return r
}

// Get returns a runtime by name.
// Returns the runtime and true if found, nil and false otherwise.
func (r *RuntimeRegistry) Get(name string) (Runtime, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rt, ok := r.runtimes[name]
	return rt, ok
}

// List returns all available runtime names.
func (r *RuntimeRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.runtimes))
	for name := range r.runtimes {
		names = append(names, name)
	}
	return names
}

// GetForRole returns the appropriate runtime for a given role.
// For the "auditor" role, it prefers Codex > OpenCode > Claude.
// For all other roles, it returns Claude.
func (r *RuntimeRegistry) GetForRole(role string) Runtime {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Role-based runtime selection
	switch role {
	case "auditor":
		// Prefer Codex for independent verification, fallback chain
		if rt, ok := r.runtimes["codex"]; ok {
			return rt
		}
		if rt, ok := r.runtimes["opencode"]; ok {
			return rt
		}
		// Last resort: same model (Claude) - still useful for syntax checks
		if rt, ok := r.runtimes["claude"]; ok {
			return rt
		}
	}

	// Default to Claude for everything else
	if rt, ok := r.runtimes["claude"]; ok {
		return rt
	}

	// No runtimes available
	return nil
}

// HasRuntime returns true if a runtime with the given name is available.
func (r *RuntimeRegistry) HasRuntime(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.runtimes[name]
	return ok
}

// Register adds a custom runtime to the registry.
// This can be used to add new runtimes dynamically.
func (r *RuntimeRegistry) Register(rt Runtime) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.runtimes[rt.Name()] = rt
}
