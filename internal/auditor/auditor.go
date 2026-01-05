// Package auditor provides independent code verification using alternate AI models.
// The Auditor agent reviews work before merge to catch issues that the original
// model might have missed.
package auditor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/steveyegge/gastown/internal/agent"
	"github.com/steveyegge/gastown/internal/beads"
)

// Common errors
var (
	ErrNoRuntime      = errors.New("no verification runtime available")
	ErrBeadNotFound   = errors.New("bead not found")
	ErrInvalidBead    = errors.New("invalid bead for verification")
	ErrParseResponse  = errors.New("failed to parse verification response")
)

// Verdict represents the outcome of a verification.
type Verdict string

const (
	// VerdictPass indicates the work passed verification.
	VerdictPass Verdict = "PASS"

	// VerdictFail indicates the work failed verification.
	VerdictFail Verdict = "FAIL"

	// VerdictNeedsHuman indicates human review is required.
	VerdictNeedsHuman Verdict = "NEEDS_HUMAN"
)

// VerificationResult contains the outcome of a verification check.
type VerificationResult struct {
	// BeadID is the ID of the bead that was verified.
	BeadID string `json:"bead_id"`

	// Verdict is the verification outcome (PASS, FAIL, NEEDS_HUMAN).
	Verdict Verdict `json:"verdict"`

	// Confidence is how confident the auditor is in the verdict (0.0-1.0).
	Confidence float64 `json:"confidence"`

	// Issues contains any problems found during verification.
	Issues []string `json:"issues,omitempty"`

	// Suggestions contains improvement suggestions (not blocking).
	Suggestions []string `json:"suggestions,omitempty"`

	// ReviewedBy is the name of the runtime that performed the review.
	ReviewedBy string `json:"reviewed_by"`

	// ReviewedAt is when the verification was performed.
	ReviewedAt time.Time `json:"reviewed_at"`

	// Duration is how long the verification took.
	Duration time.Duration `json:"duration"`
}

// IsPass returns true if the verification passed.
func (r *VerificationResult) IsPass() bool {
	return r.Verdict == VerdictPass
}

// IsFail returns true if the verification failed.
func (r *VerificationResult) IsFail() bool {
	return r.Verdict == VerdictFail
}

// NeedsHuman returns true if human review is required.
func (r *VerificationResult) NeedsHuman() bool {
	return r.Verdict == VerdictNeedsHuman
}

// Auditor performs independent verification of work using an alternate AI model.
type Auditor struct {
	runtime  agent.Runtime
	beadsDB  *beads.Beads
	registry *agent.RuntimeRegistry
}

// New creates a new Auditor with the appropriate runtime for verification.
// Uses the registry's GetForRole("auditor") to select the best available runtime.
func New(registry *agent.RuntimeRegistry, db *beads.Beads) (*Auditor, error) {
	runtime := registry.GetForRole("auditor")
	if runtime == nil {
		return nil, ErrNoRuntime
	}

	return &Auditor{
		runtime:  runtime,
		beadsDB:  db,
		registry: registry,
	}, nil
}

// NewWithRuntime creates an Auditor with a specific runtime.
// This is useful for testing or when a specific runtime is required.
func NewWithRuntime(runtime agent.Runtime, db *beads.Beads) *Auditor {
	return &Auditor{
		runtime: runtime,
		beadsDB: db,
	}
}

// RuntimeName returns the name of the runtime being used for verification.
func (a *Auditor) RuntimeName() string {
	if a.runtime == nil {
		return ""
	}
	return a.runtime.Name()
}

// Verify performs verification on a bead's associated work.
// It fetches the bead details, constructs a verification prompt,
// and executes it using the configured runtime.
func (a *Auditor) Verify(ctx context.Context, beadID string, workdir string) (*VerificationResult, error) {
	if a.runtime == nil {
		return nil, ErrNoRuntime
	}

	startTime := time.Now()

	// Get bead details
	bead, err := a.beadsDB.Show(beadID)
	if err != nil {
		if errors.Is(err, beads.ErrNotFound) {
			return nil, fmt.Errorf("%w: %s", ErrBeadNotFound, beadID)
		}
		return nil, fmt.Errorf("fetching bead: %w", err)
	}

	// Build verification prompt
	prompt := a.buildPrompt(bead, workdir)

	// Execute with configured runtime
	response, err := a.runtime.Execute(ctx, prompt, workdir)
	if err != nil {
		return &VerificationResult{
			BeadID:     beadID,
			Verdict:    VerdictNeedsHuman,
			Confidence: 0,
			Issues:     []string{fmt.Sprintf("Verification execution failed: %v", err)},
			ReviewedBy: a.runtime.Name(),
			ReviewedAt: time.Now(),
			Duration:   time.Since(startTime),
		}, nil
	}

	// Parse response
	result, err := a.parseResponse(response, beadID)
	if err != nil {
		// If parsing fails, treat as needs human review
		result = &VerificationResult{
			BeadID:     beadID,
			Verdict:    VerdictNeedsHuman,
			Confidence: 0,
			Issues:     []string{"Failed to parse verification response", err.Error()},
		}
	}

	result.ReviewedBy = a.runtime.Name()
	result.ReviewedAt = time.Now()
	result.Duration = time.Since(startTime)

	return result, nil
}

// VerifyMR performs verification specifically for a merge request.
// It uses additional context about the MR such as the branch and target.
func (a *Auditor) VerifyMR(ctx context.Context, mrID string, branch string, targetBranch string, workdir string) (*VerificationResult, error) {
	if a.runtime == nil {
		return nil, ErrNoRuntime
	}

	startTime := time.Now()

	// Build MR-specific verification prompt
	prompt := a.buildMRPrompt(mrID, branch, targetBranch, workdir)

	// Execute with configured runtime
	response, err := a.runtime.Execute(ctx, prompt, workdir)
	if err != nil {
		return &VerificationResult{
			BeadID:     mrID,
			Verdict:    VerdictNeedsHuman,
			Confidence: 0,
			Issues:     []string{fmt.Sprintf("Verification execution failed: %v", err)},
			ReviewedBy: a.runtime.Name(),
			ReviewedAt: time.Now(),
			Duration:   time.Since(startTime),
		}, nil
	}

	// Parse response
	result, err := a.parseResponse(response, mrID)
	if err != nil {
		result = &VerificationResult{
			BeadID:     mrID,
			Verdict:    VerdictNeedsHuman,
			Confidence: 0,
			Issues:     []string{"Failed to parse verification response"},
		}
	}

	result.ReviewedBy = a.runtime.Name()
	result.ReviewedAt = time.Now()
	result.Duration = time.Since(startTime)

	return result, nil
}

// buildPrompt constructs the verification prompt for a bead.
func (a *Auditor) buildPrompt(bead *beads.Issue, workdir string) string {
	return fmt.Sprintf(`You are a code reviewer performing independent verification.
Review the work in this directory for the following task.

Task ID: %s
Title: %s
Description: %s

Review the changes for:
1. Does the code meet the requirements described in the task?
2. Are there any bugs, edge cases, or error handling issues?
3. Is the code well-structured and follows best practices?
4. Are tests adequate and do they cover the main scenarios?
5. Are there any security vulnerabilities (injection, auth bypass, etc.)?

Be thorough but practical. Focus on issues that would affect correctness,
security, or maintainability.

Respond with ONLY valid JSON in this exact format:
{
  "verdict": "PASS" | "FAIL" | "NEEDS_HUMAN",
  "confidence": 0.0-1.0,
  "issues": ["issue1", "issue2"],
  "suggestions": ["suggestion1", "suggestion2"]
}

Where:
- PASS: Work meets requirements and has no blocking issues
- FAIL: Work has bugs, security issues, or doesn't meet requirements
- NEEDS_HUMAN: Unable to determine, requires human review

Only include issues that are actual problems. Suggestions are for improvements
that don't block approval.`, bead.ID, bead.Title, bead.Description)
}

// buildMRPrompt constructs a verification prompt specifically for MRs.
func (a *Auditor) buildMRPrompt(mrID, branch, targetBranch, workdir string) string {
	return fmt.Sprintf(`You are a code reviewer performing independent verification of a merge request.
Review the changes between the source and target branches.

MR ID: %s
Source Branch: %s
Target Branch: %s
Working Directory: %s

To see the changes, run:
  git diff %s...%s

Review the changes for:
1. Does the code meet the requirements?
2. Are there any bugs, edge cases, or error handling issues?
3. Is the code well-structured and follows best practices?
4. Are tests adequate and do they cover the main scenarios?
5. Are there any security vulnerabilities?
6. Will this merge cleanly with no conflicts?

Be thorough but practical. Focus on issues that would affect correctness,
security, or maintainability.

Respond with ONLY valid JSON in this exact format:
{
  "verdict": "PASS" | "FAIL" | "NEEDS_HUMAN",
  "confidence": 0.0-1.0,
  "issues": ["issue1", "issue2"],
  "suggestions": ["suggestion1", "suggestion2"]
}

Where:
- PASS: Changes are ready to merge
- FAIL: Changes have blocking issues
- NEEDS_HUMAN: Unable to determine, requires human review`, mrID, branch, targetBranch, workdir, targetBranch, branch)
}

// parseResponse extracts a VerificationResult from the runtime's response.
func (a *Auditor) parseResponse(response, beadID string) (*VerificationResult, error) {
	result := &VerificationResult{
		BeadID: beadID,
	}

	// Try to extract JSON from the response
	// The response might have text before/after the JSON
	jsonStr := extractJSON(response)
	if jsonStr == "" {
		return nil, fmt.Errorf("%w: no JSON found in response", ErrParseResponse)
	}

	// Parse the JSON
	var parsed struct {
		Verdict     string   `json:"verdict"`
		Confidence  float64  `json:"confidence"`
		Issues      []string `json:"issues"`
		Suggestions []string `json:"suggestions"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseResponse, err)
	}

	// Validate and convert verdict
	switch strings.ToUpper(parsed.Verdict) {
	case "PASS":
		result.Verdict = VerdictPass
	case "FAIL":
		result.Verdict = VerdictFail
	case "NEEDS_HUMAN":
		result.Verdict = VerdictNeedsHuman
	default:
		// Unknown verdict, treat as needs human
		result.Verdict = VerdictNeedsHuman
		result.Issues = append(result.Issues, fmt.Sprintf("Unknown verdict: %s", parsed.Verdict))
	}

	// Clamp confidence to 0-1 range
	result.Confidence = parsed.Confidence
	if result.Confidence < 0 {
		result.Confidence = 0
	}
	if result.Confidence > 1 {
		result.Confidence = 1
	}

	result.Issues = parsed.Issues
	result.Suggestions = parsed.Suggestions

	return result, nil
}

// extractJSON attempts to find and extract a JSON object from text.
// It looks for the first { and last } to extract the JSON.
func extractJSON(text string) string {
	start := strings.Index(text, "{")
	if start == -1 {
		return ""
	}

	// Find matching closing brace
	depth := 0
	for i := start; i < len(text); i++ {
		switch text[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return text[start : i+1]
			}
		}
	}

	return ""
}

// VerificationConfig holds configuration for the verification process.
type VerificationConfig struct {
	// Enabled controls whether verification is active.
	Enabled bool `json:"enabled" yaml:"enabled"`

	// RequiredConfidence is the minimum confidence for auto-approval.
	// Results below this threshold will be escalated to NEEDS_HUMAN.
	RequiredConfidence float64 `json:"required_confidence" yaml:"required_confidence"`

	// PreferredRuntime is the name of the preferred runtime for verification.
	// If not available, fallback order is: codex > opencode > claude.
	PreferredRuntime string `json:"preferred_runtime" yaml:"preferred_runtime"`

	// TimeoutSeconds is the maximum time for a verification operation.
	TimeoutSeconds int `json:"timeout_seconds" yaml:"timeout_seconds"`
}

// DefaultVerificationConfig returns the default verification configuration.
func DefaultVerificationConfig() VerificationConfig {
	return VerificationConfig{
		Enabled:            true,
		RequiredConfidence: 0.7,
		PreferredRuntime:   "codex",
		TimeoutSeconds:     300, // 5 minutes
	}
}
