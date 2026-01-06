# VerMAS Implementation Plan

> What needs to be built to make VerMAS work

This document maps out the code, hooks, commands, mail flows, and formulas needed to implement the VerMAS (Verifiable Multi-Agent System) design.

## Executive Summary

**Existing Infrastructure We Can Use:**
- `internal/auditor/auditor.go` - Already has `Verify()`, `VerifyMR()`, verdicts (PASS/FAIL/NEEDS_HUMAN)
- `internal/cmd/verify.go` - Basic verification commands exist
- `internal/mail/` - Full mail system with routing, priorities, thread support
- `internal/cmd/gate.go` - Gate coordination with wake notifications
- `.beads/formulas/*.toml` - Workflow formula system

**What Needs to Be Built:**
1. **Inspector Coordinator** - New top-level agent (like Mayor)
2. **Test Spec System** - New bead type + commands
3. **Pre-Sling Gate** - Block work until spec approved
4. **Inspector Ecosystem Agents** - Auditor, Verifier, Advocate, Critic, Judge
5. **Two-Pane Tmux Layout** - Mayor + Inspector side by side
6. **Designer & Strategist** - AI-assisted requirements/criteria

---

## Phase 1: Test Spec System

### 1.1 New Bead Type: `test-spec`

**File:** `internal/beads/types.go`

Add new bead type constant:
```go
const (
    TypeIssue     = "issue"
    TypeMessage   = "message"
    TypeTestSpec  = "test-spec"     // NEW
    TypeVerResult = "verification-result"  // NEW
)
```

**Schema for test-spec bead:**
```yaml
id: spec-abc123
type: test-spec
parent: gt-xyz789  # Work item this spec is for
status: open | closed  # closed = approved
created_at: timestamp
acceptance_criteria:
  - id: ac-1
    description: "Human-readable criterion"
    verify_command: "bash command that exits 0 for pass"
  - id: ac-2
    # ...
```

### 1.2 New Commands: `gt inspect`

**File:** `internal/cmd/inspect.go` (NEW)

```go
// Subcommands:
//   gt inspect pending      - List specs needing criteria
//   gt inspect show <id>    - Show spec details
//   gt inspect add-criteria - Add acceptance criterion
//   gt inspect approve      - Approve a spec (closes it)
//   gt inspect run          - Run full verification
//   gt inspect verify       - Run just objective tests
//   gt inspect audit        - Run just compliance check
//   gt inspect review       - Run adversarial review
//   gt inspect feedback     - Show rejection feedback
//   gt inspect status       - Check verification status
```

**Command implementations:**

| Command | What It Does | Integration |
|---------|--------------|-------------|
| `pending` | `bd list --type=test-spec --status=open` | Beads query |
| `show` | `bd show <spec-id>` with formatted criteria | Beads query |
| `add-criteria` | Append criterion to spec's acceptance_criteria | Beads update |
| `approve` | `bd close <spec-id>` + notify Mayor | Beads + mail |
| `run` | Full Inspector workflow | Calls Verifier→Auditor→Advocate→Critic→Judge |
| `verify` | Just run test commands | Shell execution |
| `audit` | Compare against spec criteria | Auditor.Verify() |
| `review` | Adversarial evaluation | Advocate + Critic + Judge |
| `feedback` | Format rejection reasons | Beads query |
| `status` | Show verification state | Beads query |

### 1.3 Auto-Create Spec on Work Creation

**File:** `internal/cmd/create.go` (modification to `bd create` hook)

When `bd create --type=task|feature|bug` is called:
1. Create the work bead as normal
2. Auto-create a `test-spec` bead with same suffix
3. Add dependency: work bead → spec bead
4. Send mail to Inspector: "NEW SPEC NEEDED"

```go
func postCreateHook(work *beads.Issue) error {
    if work.Type != "task" && work.Type != "feature" && work.Type != "bug" {
        return nil  // Only create specs for work items
    }

    // Create test spec
    specID := "spec-" + strings.TrimPrefix(work.ID, "gt-")
    spec := &beads.Issue{
        ID:          specID,
        Type:        "test-spec",
        Title:       fmt.Sprintf("Test spec for: %s", work.Title),
        Parent:      work.ID,
        Status:      "open",
        Description: "Pending Inspector approval. Define acceptance criteria.",
    }
    if err := beads.Create(spec); err != nil {
        return err
    }

    // Add dependency
    if err := beads.AddDep(work.ID, specID); err != nil {
        return err
    }

    // Notify Inspector
    mail.Send("inspector/", &mail.Message{
        Subject: fmt.Sprintf("NEW SPEC NEEDED: %s", work.Title),
        Body:    fmt.Sprintf("Work %s needs verification criteria.\nRun: gt inspect show %s", work.ID, specID),
        Type:    mail.TypeTask,
    })

    return nil
}
```

---

## Phase 2: Pre-Sling Gate

### 2.1 Modify `gt sling` to Check Spec

**File:** `internal/cmd/sling.go` (line ~200, in `runSling()`)

Add spec check before slinging:
```go
func runSling(cmd *cobra.Command, args []string) error {
    beadID := args[0]

    // NEW: Check for approved test spec
    if !slingForce {
        if err := checkTestSpecApproved(beadID); err != nil {
            return err
        }
    }

    // ... existing sling logic
}

func checkTestSpecApproved(beadID string) error {
    // Find associated spec
    specID := "spec-" + strings.TrimPrefix(beadID, "gt-")

    spec, err := beads.Show(specID)
    if err != nil {
        if errors.Is(err, beads.ErrNotFound) {
            return fmt.Errorf("no test spec found for %s\nCreate with: gt inspect create-spec %s", beadID, beadID)
        }
        return err
    }

    if spec.Status != "closed" {
        return fmt.Errorf("test spec %s not approved (status: %s)\nAction: gt inspect approve %s", specID, spec.Status, specID)
    }

    return nil
}
```

### 2.2 Add `--force` Flag for Emergency Override

**File:** `internal/cmd/sling.go`

```go
var slingForce bool
var slingForceReason string

func init() {
    slingCmd.Flags().BoolVar(&slingForce, "force", false, "Bypass spec gate (emergency)")
    slingCmd.Flags().StringVar(&slingForceReason, "reason", "", "Reason for force bypass (required with --force)")
}
```

Log bypass to audit trail:
```go
if slingForce {
    if slingForceReason == "" {
        return fmt.Errorf("--reason required when using --force")
    }
    // Log to audit
    beads.AddNote(beadID, fmt.Sprintf("GATE_BYPASS: %s (by %s)", slingForceReason, agentID))
    fmt.Printf("⚠️  GATE BYPASS: Slinging without approved spec\n")
    fmt.Printf("   Reason: %s\n", slingForceReason)
}
```

---

## Phase 3: Inspector Ecosystem

### 3.1 Inspector Directory Structure

Create directory structure for Inspector agents:
```
<rig>/
├── inspector/
│   ├── auditor/
│   │   ├── rig/          # Worktree
│   │   ├── .beads/       # Local beads
│   │   ├── .claude/      # Settings
│   │   │   └── settings.json
│   │   └── CLAUDE.md     # Role context
│   ├── advocate/
│   │   └── ...
│   ├── critic/
│   │   └── ...
│   └── judge/
│       └── ...
```

**File:** `internal/cmd/install.go` (modification)

Add Inspector directory creation to `gt install`:
```go
func createInspectorStructure(rigPath string) error {
    inspectorRoles := []string{"auditor", "advocate", "critic", "judge"}

    for _, role := range inspectorRoles {
        rolePath := filepath.Join(rigPath, "inspector", role)
        if err := os.MkdirAll(rolePath, 0755); err != nil {
            return err
        }

        // Create CLAUDE.md with role-specific context
        if err := createRoleContext(rolePath, role); err != nil {
            return err
        }

        // Create .claude/settings.json
        if err := createClaudeSettings(rolePath, role); err != nil {
            return err
        }
    }

    return nil
}
```

### 3.2 Role Context Files (CLAUDE.md)

**File:** `internal/templates/inspector_auditor.md` (embedded)

```markdown
# Auditor Context

You are the Auditor for VerMAS. Your role is to check compliance against the test spec.

## Your Responsibilities
1. Verify each acceptance criterion is met
2. Run verification commands
3. Report compliance status
4. Flag any criteria that can't be verified

## Commands
- `gt inspect audit <bead-id>` - Run audit
- `gt inspect show <spec-id>` - View spec criteria
- `bd show <bead-id>` - View work item

## When in doubt
Mail Inspector coordinator: `gt mail send inspector/ -s "Question" -m "..."`
```

Similar files for: advocate, critic, judge

### 3.3 Inspector Coordinator

**File:** `internal/cmd/inspector.go` (NEW)

```go
var inspectorCmd = &cobra.Command{
    Use:   "inspector",
    Short: "Inspector coordinator commands",
}

var inspectorStartCmd = &cobra.Command{
    Use:   "start",
    Short: "Start Inspector session",
    RunE:  runInspectorStart,
}

var inspectorAttachCmd = &cobra.Command{
    Use:   "attach",
    Short: "Attach to Inspector session",
    RunE:  runInspectorAttach,
}

// etc.
```

### 3.4 Verifier Implementation

**File:** `internal/inspector/verifier.go` (NEW)

The Verifier runs objective tests - no LLM needed:

```go
package inspector

type Verifier struct {
    workdir string
}

type VerifyResult struct {
    CriteriaID string
    Status     string // "pass" | "fail" | "error"
    Output     string
    Duration   time.Duration
}

func (v *Verifier) RunCriteria(criteria []AcceptanceCriterion) ([]VerifyResult, error) {
    var results []VerifyResult

    for _, c := range criteria {
        result := v.runSingleCriterion(c)
        results = append(results, result)
    }

    return results, nil
}

func (v *Verifier) runSingleCriterion(c AcceptanceCriterion) VerifyResult {
    start := time.Now()

    cmd := exec.Command("bash", "-c", c.VerifyCommand)
    cmd.Dir = v.workdir
    output, err := cmd.CombinedOutput()

    status := "pass"
    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            if exitErr.ExitCode() != 0 {
                status = "fail"
            }
        } else {
            status = "error"
        }
    }

    return VerifyResult{
        CriteriaID: c.ID,
        Status:     status,
        Output:     string(output),
        Duration:   time.Since(start),
    }
}
```

### 3.5 Advocate/Critic/Judge Implementation

**File:** `internal/inspector/adversarial.go` (NEW)

These use the existing auditor runtime abstraction:

```go
package inspector

import (
    "github.com/steveyegge/gastown/internal/agent"
    "github.com/steveyegge/gastown/internal/auditor"
)

type Advocate struct {
    runtime agent.Runtime
}

func (a *Advocate) BuildDefense(evidence *Evidence) (*Brief, error) {
    prompt := buildAdvocatePrompt(evidence)
    response, err := a.runtime.Execute(context.Background(), prompt, evidence.Workdir)
    if err != nil {
        return nil, err
    }
    return parseDefenseBrief(response)
}

type Critic struct {
    runtime agent.Runtime
}

func (c *Critic) BuildProsecution(evidence *Evidence) (*Brief, error) {
    prompt := buildCriticPrompt(evidence)
    response, err := c.runtime.Execute(context.Background(), prompt, evidence.Workdir)
    if err != nil {
        return nil, err
    }
    return parseProsecutionBrief(response)
}

type Judge struct {
    runtime agent.Runtime
}

func (j *Judge) Deliberate(defenseBrief, prosecutionBrief *Brief, evidence *Evidence) (*auditor.VerificationResult, error) {
    prompt := buildJudgePrompt(defenseBrief, prosecutionBrief, evidence)
    response, err := j.runtime.Execute(context.Background(), prompt, evidence.Workdir)
    if err != nil {
        return nil, err
    }
    return parseVerdict(response)
}
```

---

## Phase 4: Refinery Integration

### 4.1 Modify Refinery to Call Inspector

**File:** `internal/refinery/process.go` (modification)

```go
func processMergeRequest(mr *mrqueue.MergeRequest) error {
    // Existing: validate MR, check branch

    // NEW: Run Inspector verification
    result, err := runInspectorVerification(mr)
    if err != nil {
        return fmt.Errorf("verification failed: %w", err)
    }

    switch result.Verdict {
    case auditor.VerdictPass:
        return merge(mr)
    case auditor.VerdictFail:
        return reject(mr, result)
    case auditor.VerdictNeedsHuman:
        return escalate(mr, result)
    }

    return nil
}

func runInspectorVerification(mr *mrqueue.MergeRequest) (*auditor.VerificationResult, error) {
    // 1. Run Verifier (objective tests)
    verifier := inspector.NewVerifier(mr.Workdir)
    spec := loadTestSpec(mr.BeadID)
    verifyResults, err := verifier.RunCriteria(spec.AcceptanceCriteria)
    if err != nil {
        return nil, err
    }

    // 2. Run Auditor (compliance check)
    aud, err := auditor.New(registry, beadsDB)
    if err != nil {
        return nil, err
    }
    auditResult, err := aud.Verify(context.Background(), mr.BeadID, mr.Workdir)
    if err != nil {
        return nil, err
    }

    // 3. Run Adversarial Review
    evidence := &inspector.Evidence{
        VerifyResults: verifyResults,
        AuditResult:   auditResult,
        Diff:          getDiff(mr.Branch),
        Workdir:       mr.Workdir,
    }

    advocate := inspector.NewAdvocate(registry)
    defenseBrief, _ := advocate.BuildDefense(evidence)

    critic := inspector.NewCritic(registry)
    prosecutionBrief, _ := critic.BuildProsecution(evidence)

    judge := inspector.NewJudge(registry)
    return judge.Deliberate(defenseBrief, prosecutionBrief, evidence)
}
```

### 4.2 Rejection Flow

**File:** `internal/refinery/reject.go` (NEW)

```go
func reject(mr *mrqueue.MergeRequest, result *auditor.VerificationResult) error {
    // Create feedback bead
    feedbackID := "feedback-" + mr.BeadID[3:]
    feedback := &beads.Issue{
        ID:          feedbackID,
        Type:        "verification-result",
        Parent:      mr.BeadID,
        Title:       fmt.Sprintf("Verification FAILED: %s", mr.Title),
        Description: formatFailureReport(result),
        Status:      "open",
    }
    if err := beads.Create(feedback); err != nil {
        return err
    }

    // Notify Mayor
    mail.Send("mayor/", &mail.Message{
        Subject:  fmt.Sprintf("❌ REJECTED: %s", mr.Title),
        Body:     formatRejectionMail(result, feedbackID),
        Type:     mail.TypeTask,
        Priority: mail.PriorityHigh,
    })

    // Update work bead status
    beads.Update(mr.BeadID, beads.UpdateOptions{Status: "rejected"})

    return nil
}
```

---

## Phase 5: Mail Flows

### 5.1 New Mail Routes

**File:** `internal/mail/router.go` (modification)

Add Inspector addresses:
```go
func (r *Router) resolveAddress(addr string) (string, error) {
    switch {
    case addr == "inspector/" || addr == "inspector":
        return r.inspectorMailbox(), nil
    case strings.HasPrefix(addr, "inspector/"):
        // inspector/auditor, inspector/judge, etc.
        role := strings.TrimPrefix(addr, "inspector/")
        return r.inspectorRoleMailbox(role), nil
    // ... existing cases
    }
}
```

### 5.2 Mail Flow Definitions

| From | To | Subject Pattern | Purpose |
|------|-----|-----------------|---------|
| `bd create` hook | `inspector/` | "NEW SPEC NEEDED: ..." | Trigger spec creation |
| Inspector | `mayor/` | "SPEC APPROVED: ..." | Unblock work |
| Inspector | `mayor/` | "NEEDS_CLARIFICATION: ..." | Question from any role |
| Mayor | `inspector/` | "RE: NEEDS_CLARIFICATION: ..." | Answer to question |
| Refinery | `mayor/` | "❌ REJECTED: ..." | Verification failure |
| Refinery | `mayor/` | "✅ MERGED: ..." | Verification success |
| Inspector | `<rig>/refinery` | "VERIFICATION COMPLETE: ..." | Signal proceed/reject |

### 5.3 Question Flow Implementation

**File:** `internal/inspector/question.go` (NEW)

```go
func askMayor(question string, beadID string, fromRole string) (string, error) {
    // Send question
    msg := mail.Send("mayor/", &mail.Message{
        From:    fmt.Sprintf("inspector/%s", fromRole),
        Subject: fmt.Sprintf("NEEDS_CLARIFICATION: %s", beadID),
        Body:    question,
        Type:    mail.TypeTask,
    })

    // Create gate to wait for response
    gateID := fmt.Sprintf("gate-question-%s", msg.ID)
    gate := gates.Create(gates.GateOptions{
        ID:   gateID,
        Type: "mail",
        Condition: fmt.Sprintf("mail reply to %s", msg.ID),
    })

    // Wait for reply (or timeout)
    reply, err := gate.Wait(5 * time.Minute)
    if err != nil {
        return "", fmt.Errorf("no response from Mayor: %w", err)
    }

    return reply.Body, nil
}
```

---

## Phase 6: Designer & Strategist Agents

### 6.1 Designer Agent

**File:** `internal/inspector/designer.go` (NEW)

```go
type Designer struct {
    runtime agent.Runtime
}

func (d *Designer) ElaborateRequirements(request string) (*Requirements, error) {
    prompt := `You are the Designer. Turn this vague request into detailed requirements.

User request: "` + request + `"

Produce requirements covering:
1. Functional - What exactly should this do?
2. Non-functional - Performance, quality, style
3. Constraints - What should it NOT do?
4. Examples - Concrete input/output examples

Be specific enough that a developer could implement without questions.
Output as JSON: {"functional": [...], "non_functional": [...], "constraints": [...], "examples": [...]}`

    response, err := d.runtime.Execute(context.Background(), prompt, "")
    if err != nil {
        return nil, err
    }

    return parseRequirements(response)
}
```

### 6.2 Strategist Agent

**File:** `internal/inspector/strategist.go` (NEW)

```go
type Strategist struct {
    runtime agent.Runtime
}

func (s *Strategist) ProposeCriteria(requirements *Requirements) ([]AcceptanceCriterion, error) {
    prompt := `You are the Strategist. Design verification criteria for these requirements.

Requirements:
` + formatRequirements(requirements) + `

For each key behavior, create a criterion:
- ID: AC-N
- Description: Human-readable check
- Command: Bash command that exits 0 for pass, non-0 for fail

Focus on correctness, edge cases, and quality. Don't over-test.
Output as JSON array: [{"id": "AC-1", "description": "...", "verify_command": "..."}, ...]`

    response, err := s.runtime.Execute(context.Background(), prompt, "")
    if err != nil {
        return nil, err
    }

    return parseCriteria(response)
}
```

### 6.3 Integration with Work Creation

**File:** `internal/cmd/create.go` (further modification)

```go
func createWorkItemWithAIAssist(title string, opts CreateOptions) (*beads.Issue, error) {
    // 1. Call Designer to elaborate requirements
    designer := inspector.NewDesigner(registry)
    requirements, err := designer.ElaborateRequirements(title)
    if err != nil {
        return nil, err
    }

    // 2. Present to user for approval
    fmt.Println("Designer suggests these requirements:")
    fmt.Println(formatRequirements(requirements))
    fmt.Print("\nApprove? [y/n/modify]: ")

    // ... handle user input

    // 3. Create work bead with full requirements
    work := &beads.Issue{
        Title:       title,
        Description: formatRequirementsForBead(requirements),
        // ...
    }

    // 4. Call Strategist to propose criteria
    strategist := inspector.NewStrategist(registry)
    criteria, err := strategist.ProposeCriteria(requirements)
    if err != nil {
        return nil, err
    }

    // 5. Present criteria to user
    fmt.Println("\nStrategist suggests these verification criteria:")
    for _, c := range criteria {
        fmt.Printf("  [%s] %s\n", c.ID, c.Description)
    }
    fmt.Print("\nApprove? [y/n/modify]: ")

    // ... handle user input

    // 6. Create spec with approved criteria
    spec := &beads.Issue{
        Type:     "test-spec",
        Parent:   work.ID,
        Metadata: map[string]interface{}{"criteria": criteria},
    }

    return work, nil
}
```

---

## Phase 7: Two-Pane Tmux Layout

### 7.1 VerMAS Start Command

**File:** `internal/cmd/vermas.go` (NEW)

```go
var vermasCmd = &cobra.Command{
    Use:   "vermas",
    Short: "VerMAS multi-pane commands",
}

var vermasStartCmd = &cobra.Command{
    Use:   "start",
    Short: "Start Mayor and Inspector in split panes",
    RunE:  runVermasStart,
}

func runVermasStart(cmd *cobra.Command, args []string) error {
    // Create tmux session with split layout
    sessionName := fmt.Sprintf("vermas-%d", os.Getpid())

    // Create session with Mayor pane
    if err := tmux.NewSession(sessionName, "mayor", townRoot); err != nil {
        return err
    }

    // Split and create Inspector pane
    if err := tmux.SplitHorizontal(sessionName, "inspector", townRoot); err != nil {
        return err
    }

    // Start Mayor agent in left pane
    if err := tmux.SendKeys(sessionName, "0", "claude --profile mayor"); err != nil {
        return err
    }

    // Start Inspector agent in right pane
    if err := tmux.SendKeys(sessionName, "1", "claude --profile inspector"); err != nil {
        return err
    }

    fmt.Printf("Started VerMAS session: %s\n", sessionName)
    fmt.Printf("Attach with: gt vermas attach\n")

    return nil
}
```

### 7.2 Layout Configuration

**File:** `config/vermas.yaml`

```yaml
layout:
  type: horizontal
  ratio: 50

panes:
  - name: mayor
    role: mayor
    position: left
    startup:
      - "gt prime"
      - "gt mail inbox"

  - name: inspector
    role: inspector
    position: right
    startup:
      - "gt prime"
      - "gt inspect pending"

status_bar:
  enabled: true
  format: "Mayor: #{mayor_status} | Inspector: #{inspector_status} | Queue: #{queue_depth}"
```

---

## Phase 8: Hooks

### 8.1 New Hook Types

**File:** `internal/hooks/types.go` (modification)

```go
const (
    // Existing hooks
    HookSessionStart     = "SessionStart"
    HookPreToolUse       = "PreToolUse"
    // ...

    // NEW: Verification hooks
    HookPreAuditPlan     = "PreAuditPlan"
    HookPostAuditPlan    = "PostAuditPlan"
    HookPreInspect       = "PreInspect"
    HookPostVerify       = "PostVerify"
    HookPostAudit        = "PostAudit"
    HookPreAdvocate      = "PreAdvocate"
    HookPostAdvocate     = "PostAdvocate"
    HookPreCritic        = "PreCritic"
    HookPostCritic       = "PostCritic"
    HookPreJudgment      = "PreJudgment"
    HookPostJudgment     = "PostJudgment"
    HookOnVerifyPass     = "OnVerifyPass"
    HookOnVerifyFail     = "OnVerifyFail"
    HookOnVerifyEscalate = "OnVerifyEscalate"
)
```

### 8.2 Hook Integration Points

| Hook | Trigger Point | Input | Purpose |
|------|---------------|-------|---------|
| `PreAuditPlan` | Before `gt sling` creates spec | bead_id | Customize spec creation |
| `PostAuditPlan` | After spec created | bead_id, spec_id | Notify, log |
| `PreInspect` | Before Inspector runs | mr_id, bead_id | Pre-verification checks |
| `PostVerify` | After Verifier | bead_id, results | Log test results |
| `PostAudit` | After Auditor | bead_id, report | Log compliance |
| `PreAdvocate` | Before Advocate | bead_id, evidence | Custom evidence |
| `PostAdvocate` | After Advocate | bead_id, brief | Review defense |
| `PreCritic` | Before Critic | bead_id, evidence | Custom evidence |
| `PostCritic` | After Critic | bead_id, brief | Review prosecution |
| `PreJudgment` | Before Judge | bead_id, briefs | Custom input |
| `PostJudgment` | After verdict | bead_id, verdict | Log verdict |
| `OnVerifyPass` | Verdict = PASS | bead_id, result | Custom merge actions |
| `OnVerifyFail` | Verdict = FAIL | bead_id, result, feedback | Custom rejection |
| `OnVerifyEscalate` | Verdict = NEEDS_HUMAN | bead_id, question | Custom escalation |

---

## Phase 9: Formulas

### 9.1 `mol-inspect.formula.toml`

**File:** `.beads/formulas/mol-inspect.formula.toml` (NEW)

```toml
description = """
Inspector verification workflow.
Runs the full adversarial verification process.
"""

formula = "mol-inspect"
version = 1

[[steps]]
id = "verify"
title = "Run objective verification"
description = """
Verifier runs test commands from test spec.
No LLM needed - just shell execution.

Commands:
  gt inspect verify {{bead_id}}

Exit: All test results captured
"""

[[steps]]
id = "audit"
title = "Check compliance"
needs = ["verify"]
description = """
Auditor checks work against acceptance criteria.

Commands:
  gt inspect audit {{bead_id}}

Exit: Compliance report generated
"""

[[steps]]
id = "advocate"
title = "Build defense"
needs = ["audit"]
description = """
Advocate reviews evidence and argues FOR the code.

Commands:
  gt advocate build {{bead_id}}

Exit: Defense brief submitted
"""

[[steps]]
id = "critic"
title = "Build prosecution"
needs = ["audit"]
description = """
Critic reviews evidence and argues AGAINST the code.
Runs in parallel with Advocate.

Commands:
  gt critic build {{bead_id}}

Exit: Prosecution brief submitted
"""

[[steps]]
id = "judge"
title = "Deliver verdict"
needs = ["advocate", "critic"]
description = """
Judge reviews all evidence and briefs.
Delivers PASS, FAIL, or NEEDS_HUMAN.

Commands:
  gt judge deliberate {{bead_id}}

If NEEDS_HUMAN:
  gt mail send mayor/ -s "NEEDS_CLARIFICATION" -m "<question>"
  gt mol await-signal mayor-response

Exit: Verdict delivered
"""

[[steps]]
id = "route"
title = "Route verdict"
needs = ["judge"]
description = """
Route based on verdict:
- PASS → signal refinery to merge
- FAIL → mail mayor with feedback
- NEEDS_HUMAN → wait for response, re-judge

Commands:
  gt inspect route {{bead_id}}
"""

[vars]
[vars.bead_id]
description = "The work item being verified"
required = true
```

### 9.2 `mol-spec-planning.formula.toml`

**File:** `.beads/formulas/mol-spec-planning.formula.toml` (NEW)

```toml
description = """
AI-assisted spec planning workflow.
Helps user define requirements and verification criteria.
"""

formula = "mol-spec-planning"
version = 1

[[steps]]
id = "design"
title = "Elaborate requirements"
description = """
Designer analyzes user request and proposes detailed requirements.

Commands:
  gt design elaborate "{{request}}"

Present to user for approval.
"""

[[steps]]
id = "approve-design"
title = "Get user approval on design"
needs = ["design"]
description = """
User reviews and approves/modifies requirements.

If rejected, loop back to design with feedback.
"""

[[steps]]
id = "strategize"
title = "Propose verification criteria"
needs = ["approve-design"]
description = """
Strategist analyzes requirements and proposes test criteria.

Commands:
  gt strategize {{bead_id}}

Present criteria to user for approval.
"""

[[steps]]
id = "approve-strategy"
title = "Get user approval on criteria"
needs = ["strategize"]
description = """
User reviews and approves/modifies criteria.

If rejected, loop back to strategize with feedback.
"""

[[steps]]
id = "create-spec"
title = "Create approved spec"
needs = ["approve-strategy"]
description = """
Create test spec bead with approved criteria.
Closes the spec, unblocking the work item.

Commands:
  gt inspect approve {{spec_id}}
"""

[vars]
[vars.request]
description = "The user's original request"
required = true
[vars.bead_id]
description = "The work item ID"
required = true
[vars.spec_id]
description = "The test spec ID"
required = true
```

---

## Implementation Order

### Phase 1 (Foundation) - 1-2 weeks
1. Test spec bead type
2. `gt inspect` command skeleton
3. Pre-sling gate check
4. Auto-create spec on work creation

### Phase 2 (Inspector Agents) - 2-3 weeks
1. Verifier (no LLM, shell execution)
2. Auditor integration (existing code)
3. Advocate implementation
4. Critic implementation
5. Judge implementation
6. Inspector coordinator

### Phase 3 (Integration) - 1-2 weeks
1. Refinery gate integration
2. Rejection flow
3. Mail flows
4. Hook integration

### Phase 4 (AI Assistance) - 1-2 weeks
1. Designer agent
2. Strategist agent
3. User approval UX

### Phase 5 (Polish) - 1 week
1. Two-pane tmux layout
2. Formulas
3. Documentation
4. Testing

---

## Files to Create

| File | Purpose |
|------|---------|
| `internal/cmd/inspect.go` | `gt inspect` commands |
| `internal/cmd/vermas.go` | `gt vermas` commands |
| `internal/inspector/verifier.go` | Objective test runner |
| `internal/inspector/adversarial.go` | Advocate/Critic/Judge |
| `internal/inspector/designer.go` | Requirements elaboration |
| `internal/inspector/strategist.go` | Criteria proposal |
| `internal/inspector/question.go` | Question flow to Mayor |
| `internal/refinery/inspect.go` | Refinery verification integration |
| `internal/refinery/reject.go` | Rejection handling |
| `.beads/formulas/mol-inspect.formula.toml` | Verification workflow |
| `.beads/formulas/mol-spec-planning.formula.toml` | Spec planning workflow |
| `roles/inspector/*.md` | Role context files |
| `config/vermas.yaml` | Layout configuration |

## Files to Modify

| File | Changes |
|------|---------|
| `internal/beads/types.go` | Add test-spec, verification-result types |
| `internal/cmd/sling.go` | Add spec gate check |
| `internal/cmd/create.go` | Auto-create spec hook |
| `internal/cmd/install.go` | Create Inspector directories |
| `internal/mail/router.go` | Add inspector/ address routing |
| `internal/hooks/types.go` | Add verification hooks |
| `internal/refinery/process.go` | Call Inspector before merge |

---

## Testing Strategy

### Unit Tests
- Verifier: Mock shell commands, verify pass/fail detection
- Adversarial: Mock runtime, verify prompt construction
- Gate: Verify spec check logic

### Integration Tests
- Full flow: Create work → spec → sling → verify → merge
- Rejection flow: Verify → fail → feedback → rework
- Question flow: Judge question → Mayor → response → verdict

### E2E Tests
- FizzBuzz dry run from design doc
- Multi-issue convoy with mixed verdicts
- Emergency bypass with audit trail
