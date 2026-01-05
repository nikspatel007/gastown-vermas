// Package refinery verification integration.
// This file adds auditor verification capabilities to the merge queue.

package refinery

import (
	"context"
	"fmt"
	"time"

	"github.com/steveyegge/gastown/internal/agent"
	"github.com/steveyegge/gastown/internal/auditor"
	"github.com/steveyegge/gastown/internal/beads"
)

// VerificationStatus represents the verification state of an MR.
type VerificationStatus string

const (
	// VerificationPending means the MR has not been verified yet.
	VerificationPending VerificationStatus = "pending"

	// VerificationVerified means the MR passed verification.
	VerificationVerified VerificationStatus = "verified"

	// VerificationRejected means the MR failed verification.
	VerificationRejected VerificationStatus = "rejected"

	// VerificationNeedsReview means the MR requires human review.
	VerificationNeedsReview VerificationStatus = "needs_review"

	// VerificationSkipped means verification was skipped (not configured).
	VerificationSkipped VerificationStatus = "skipped"
)

// VerificationInfo contains the verification details for an MR.
type VerificationInfo struct {
	Status      VerificationStatus `json:"status"`
	ReviewedBy  string             `json:"reviewed_by,omitempty"`
	Confidence  float64            `json:"confidence,omitempty"`
	Issues      []string           `json:"issues,omitempty"`
	Suggestions []string           `json:"suggestions,omitempty"`
	VerifiedAt  *time.Time         `json:"verified_at,omitempty"`
}

// VerifiableMR extends MergeRequest with verification capabilities.
type VerifiableMR struct {
	*MergeRequest
	Verification *VerificationInfo `json:"verification,omitempty"`
}

// IsVerified returns true if the MR has passed verification.
func (v *VerifiableMR) IsVerified() bool {
	if v.Verification == nil {
		return false
	}
	return v.Verification.Status == VerificationVerified
}

// NeedsVerification returns true if the MR needs to be verified.
func (v *VerifiableMR) NeedsVerification() bool {
	if v.Verification == nil {
		return true
	}
	return v.Verification.Status == VerificationPending
}

// VerificationGate handles verification of merge requests before merge.
type VerificationGate struct {
	auditor  *auditor.Auditor
	registry *agent.RuntimeRegistry
	config   auditor.VerificationConfig
}

// NewVerificationGate creates a new verification gate.
// Returns nil if no verification runtime is available (verification disabled).
func NewVerificationGate(workdir string) (*VerificationGate, error) {
	registry := agent.NewRuntimeRegistry()
	db := beads.New(workdir)

	aud, err := auditor.New(registry, db)
	if err != nil {
		// No runtime available - verification will be skipped
		return nil, err
	}

	return &VerificationGate{
		auditor:  aud,
		registry: registry,
		config:   auditor.DefaultVerificationConfig(),
	}, nil
}

// NewVerificationGateWithConfig creates a gate with custom configuration.
func NewVerificationGateWithConfig(workdir string, config auditor.VerificationConfig) (*VerificationGate, error) {
	registry := agent.NewRuntimeRegistry()
	db := beads.New(workdir)

	aud, err := auditor.New(registry, db)
	if err != nil {
		return nil, err
	}

	return &VerificationGate{
		auditor:  aud,
		registry: registry,
		config:   config,
	}, nil
}

// RuntimeName returns the name of the verification runtime being used.
func (g *VerificationGate) RuntimeName() string {
	if g.auditor == nil {
		return ""
	}
	return g.auditor.RuntimeName()
}

// VerifyMR performs verification on a merge request.
// Returns the verification result or an error.
func (g *VerificationGate) VerifyMR(ctx context.Context, mr *MergeRequest, workdir string) (*VerificationInfo, error) {
	if !g.config.Enabled {
		return &VerificationInfo{
			Status: VerificationSkipped,
		}, nil
	}

	// Create a timeout context if configured
	if g.config.TimeoutSeconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(g.config.TimeoutSeconds)*time.Second)
		defer cancel()
	}

	// Perform verification
	result, err := g.auditor.VerifyMR(ctx, mr.ID, mr.Branch, mr.TargetBranch, workdir)
	if err != nil {
		return &VerificationInfo{
			Status: VerificationNeedsReview,
			Issues: []string{fmt.Sprintf("Verification error: %v", err)},
		}, err
	}

	// Convert result to VerificationInfo
	info := &VerificationInfo{
		ReviewedBy:  result.ReviewedBy,
		Confidence:  result.Confidence,
		Issues:      result.Issues,
		Suggestions: result.Suggestions,
	}

	now := result.ReviewedAt
	info.VerifiedAt = &now

	// Determine status based on verdict and confidence
	switch result.Verdict {
	case auditor.VerdictPass:
		if result.Confidence >= g.config.RequiredConfidence {
			info.Status = VerificationVerified
		} else {
			// Pass but low confidence - needs human review
			info.Status = VerificationNeedsReview
			info.Issues = append(info.Issues,
				fmt.Sprintf("Confidence %.2f below threshold %.2f",
					result.Confidence, g.config.RequiredConfidence))
		}

	case auditor.VerdictFail:
		info.Status = VerificationRejected

	case auditor.VerdictNeedsHuman:
		info.Status = VerificationNeedsReview
	}

	return info, nil
}

// VerifyBead performs verification on a specific bead.
func (g *VerificationGate) VerifyBead(ctx context.Context, beadID string, workdir string) (*VerificationInfo, error) {
	if !g.config.Enabled {
		return &VerificationInfo{
			Status: VerificationSkipped,
		}, nil
	}

	// Create a timeout context if configured
	if g.config.TimeoutSeconds > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(g.config.TimeoutSeconds)*time.Second)
		defer cancel()
	}

	// Perform verification
	result, err := g.auditor.Verify(ctx, beadID, workdir)
	if err != nil {
		return &VerificationInfo{
			Status: VerificationNeedsReview,
			Issues: []string{fmt.Sprintf("Verification error: %v", err)},
		}, err
	}

	// Convert result to VerificationInfo
	info := &VerificationInfo{
		ReviewedBy:  result.ReviewedBy,
		Confidence:  result.Confidence,
		Issues:      result.Issues,
		Suggestions: result.Suggestions,
	}

	now := result.ReviewedAt
	info.VerifiedAt = &now

	// Determine status based on verdict and confidence
	switch result.Verdict {
	case auditor.VerdictPass:
		if result.Confidence >= g.config.RequiredConfidence {
			info.Status = VerificationVerified
		} else {
			info.Status = VerificationNeedsReview
		}
	case auditor.VerdictFail:
		info.Status = VerificationRejected
	case auditor.VerdictNeedsHuman:
		info.Status = VerificationNeedsReview
	}

	return info, nil
}

// CanProceed returns true if the MR can proceed to merge based on verification.
func (g *VerificationGate) CanProceed(info *VerificationInfo) bool {
	if info == nil {
		return false
	}
	return info.Status == VerificationVerified || info.Status == VerificationSkipped
}

// ShouldSlingBack returns true if the MR should be sent back for fixes.
func (g *VerificationGate) ShouldSlingBack(info *VerificationInfo) bool {
	if info == nil {
		return false
	}
	return info.Status == VerificationRejected
}

// ShouldEscalate returns true if the MR needs human review.
func (g *VerificationGate) ShouldEscalate(info *VerificationInfo) bool {
	if info == nil {
		return true // No info means we need to escalate
	}
	return info.Status == VerificationNeedsReview
}
