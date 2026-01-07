# VerMAS Design - Approved Decisions

> Approved design decisions for the Python implementation

---

## Organizational Model (Approved)

### Communication Patterns

1. **Identity** - Everyone has a name/role (AGENT_ID)
2. **Department** - Groups of identities with shared expertise (engineering, marketing, QA)
3. **Assembly Line** - Where work actually happens (git worktrees)
4. **Communication** - Messages between people (mail)
5. **Work tracking** - Tasks, tickets, assignments (work orders)
6. **Processes** - Workflows, procedures (graph-based templates)
7. **Events** - "What happened" log (audit trail)
8. **Hierarchy** - CEO → Departments → Workers
9. **Artifacts** - Deliverables and documents that flow between departments

### Data Flow: What Moves Through the Organization

In a real organization, **data/artifacts** flow between departments:

| Org Concept | System Primitive | Examples |
|-------------|------------------|----------|
| Documents | Attachments, files | Specs, designs, reports |
| Deliverables | Git commits/branches | Code changes, PRs |
| Reports | Generated outputs | Test results, build logs |
| Records | Evidence files | Verification artifacts, audit trails |
| Handoff package | Work order + attachments | Everything needed for next step |

**Key insight:** Storage (filing cabinet) is where things *rest*. Artifacts are things that *move*.

```
ARTIFACT LIFECYCLE:

  Created          Attached           Handed off         Archived
     │                │                   │                  │
     ▼                ▼                   ▼                  ▼
 ┌────────┐     ┌──────────┐       ┌──────────┐       ┌──────────┐
 │ Worker │────▶│Work Order│──────▶│ Next Dept│──────▶│ Evidence │
 │produces│     │attachment│       │ receives │       │ storage  │
 └────────┘     └──────────┘       └──────────┘       └──────────┘
```

**System primitives for artifacts:**
- **Git branches/commits** - Code artifacts (the actual work product)
- **Attachments** - Files linked to work orders
- **Evidence** - Verification artifacts (.work/evidence/)
- **Logs** - Process outputs, build logs

### Business Primitives (by department)

```
FOUNDATION LAYER (infrastructure)
├── 1. Identity System - Who are you?
├── 2. Storage Layer - Where does data live?
└── 3. Event Log - What happened?

COMMUNICATION LAYER (internal comms)
├── 4. Mail System - How do agents talk?
└── 5. Assignment System - How is work dispatched?

WORK LAYER (project management)
├── 6. Work Order System - What needs to be done?
├── 7. Process System - How is work tracked? (graph-based workflows)
└── 8. Verification System - Is work correct?

HR LAYER (human resources)
├── 9. Agent Lifecycle - Spawn/monitor/kill agents
└── 10. Assembly Lines - Git worktrees for work execution

OPS LAYER (customer-facing operations)
├── 11. CLI Interface - Bring in work, check quality
└── 12. Customer Feedback - Bug reports, feature requests
```

### Org → System Mapping

| Org Concept | System Primitive | Layer |
|-------------|------------------|-------|
| Employee badge | AGENT_ID | Foundation |
| Department (Eng, QA, Marketing) | Group of identities + expertise | Foundation |
| Filing cabinet | .work/ directory | Foundation |
| Activity log | events.jsonl | Foundation |
| Email/Mailbox | messages.jsonl + mail protocol | Communication |
| Desk assignment | .assignment-{agent} | Communication |
| Ticket system | work_orders.jsonl | Work |
| Process/Workflow | Graph-based templates (LangGraph-style) | Work |
| Documents/Deliverables | Git commits, branches, attachments | Artifacts |
| Reports/Logs | Generated outputs, build logs | Artifacts |
| Records/Evidence | .work/evidence/, audit trails | Artifacts |
| HR Department | Agent lifecycle (spawn/monitor/kill) | HR |
| Assembly Line | Git worktree | HR |
| Computer terminal | Tmux session | Work Tool |
| Sales/Customer Service | CLI (bring in work) | Ops |
| Customer Success | Bug handling, feedback | Ops |
| QA department | Verification pipeline | Quality |

---

## LangGraph-Style Workflow Notation (Approved)

> Making it easy for users to define their own processes

### Design Decision

Use graph-based workflow definitions with **nodes** (steps) and **edges** (transitions).

### YAML Syntax

```yaml
# .work/templates/code-review.yaml

name: code-review
description: Standard code review workflow

nodes:
  start:
    type: entry
    next: assign_reviewer

  assign_reviewer:
    type: action
    agent: supervisor
    prompt: "Assign a reviewer for this PR"
    outputs: [reviewer_id]
    next: review_code

  review_code:
    type: action
    agent: "{reviewer_id}"  # Dynamic agent from previous step
    prompt: "Review the code changes"
    outputs: [approved, comments]
    next: check_approval

  check_approval:
    type: condition
    expression: "approved == true"
    true_next: merge
    false_next: request_changes

  request_changes:
    type: action
    agent: "{reviewer_id}"
    prompt: "Request changes based on: {comments}"
    next: wait_for_fixes

  wait_for_fixes:
    type: wait
    event: "work_order.updated"
    filter: "status == 'ready_for_review'"
    next: review_code  # Loop back

  merge:
    type: action
    agent: qa
    prompt: "Merge the approved PR"
    next: end

  end:
    type: exit
    status: completed
```

### Node Types

| Type | Purpose | Example |
|------|---------|---------|
| `entry` | Start of workflow | `start` node |
| `exit` | End of workflow | `completed`, `failed`, `cancelled` |
| `action` | Agent performs work | Review code, write tests |
| `condition` | Branch based on data | If approved → merge, else → revise |
| `wait` | Wait for external event | Wait for PR update, wait for approval |
| `parallel` | Run multiple steps concurrently | Run tests + linting simultaneously |
| `human` | Requires human input | Approval gates, escalations |

### Edge Types

```yaml
# Simple next
next: review_code

# Conditional branching
true_next: merge
false_next: request_changes

# Parallel fan-out
parallel_next: [run_tests, run_lint, run_security_scan]

# Parallel fan-in (wait for all)
join_from: [run_tests, run_lint, run_security_scan]
next: check_all_passed

# Error handling
error_next: escalate_to_human
```

### Alternative: TOML Syntax

```toml
[workflow]
name = "code-review"
description = "Standard code review workflow"

[nodes.start]
type = "entry"
next = "assign_reviewer"

[nodes.assign_reviewer]
type = "action"
agent = "supervisor"
prompt = "Assign a reviewer"
next = "review_code"

[nodes.check_approval]
type = "condition"
expression = "approved == true"
true_next = "merge"
false_next = "request_changes"
```

### Implementation Structure

```
vermas/workflows/
├── schema.py      # Pydantic models for nodes, edges, graphs
├── compiler.py    # YAML/TOML → ProcessGraph
├── executor.py    # Run compiled graphs
├── visualizer.py  # Generate ASCII/Mermaid diagrams
└── builtins/      # Built-in node type implementations
    ├── action.py
    ├── condition.py
    ├── wait.py
    └── parallel.py
```

### Benefits

1. **Declarative** - Define *what*, not *how*
2. **Visual** - Can generate diagrams automatically
3. **Composable** - Nest workflows, reuse nodes
4. **Portable** - YAML/TOML works anywhere
5. **Extensible** - Add custom node types
6. **LangGraph-compatible** - Similar mental model, easy migration

---

## Governance Model (Approved)

### How Real Organizations Govern

```
GOVERNANCE LAYERS (top to bottom)

┌─────────────────────────────────────────────────────────────────────────────┐
│ BOARD / SHAREHOLDERS                                                         │
│ Mission, Vision, Values - "Why we exist"                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│ EXECUTIVE (CEO)                                                              │
│ Strategy, Goals, Priorities - "What we're trying to achieve"               │
├─────────────────────────────────────────────────────────────────────────────┤
│ COMPLIANCE / LEGAL                                                           │
│ Rules, Constraints, Requirements - "What we MUST do / MUST NOT do"          │
├─────────────────────────────────────────────────────────────────────────────┤
│ PLANNING (Quarterly / Sprint)                                               │
│ Objectives, Milestones, Deadlines - "When things need to happen"            │
├─────────────────────────────────────────────────────────────────────────────┤
│ OPERATIONS (Managers)                                                        │
│ Routing, Assignment, Prioritization - "Who does what"                       │
├─────────────────────────────────────────────────────────────────────────────┤
│ EXECUTION (Workers)                                                          │
│ Actual Work - "Getting it done"                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Governance Concepts → System Primitives

| Org Concept | System Primitive | Enforcement |
|-------------|------------------|-------------|
| Mission/Vision | CLAUDE.md preamble | Prompt engineering |
| Strategy/Goals | Quarterly objectives file | Reference document |
| Compliance rules | Hooks + validators | Hard enforcement |
| Planning cycles | Sprint/milestone tracking | Time-boxed work |
| Work routing | Assignment algorithm | Automated dispatch |
| Prioritization | Priority field + queue | Ordering logic |
| Escalation | Mail + timeout triggers | Automated alerts |
| Audit trail | events.jsonl | Append-only log |

### File Format Convention

**Principle**: Use the format that matches how the content is consumed.

```
FORMAT DECISION TREE

Is it read by LLMs/humans as prose?
    │
    ├── YES → Use Markdown (.md)
    │         • Constitution (mission, vision, values)
    │         • Objectives and strategy
    │         • Department/team context
    │         • CLAUDE.md files can reference directly
    │
    └── NO (machine-parsed structure) → Use YAML (.yaml)
              • Workflow definitions (nodes, edges)
              • Compliance rules (patterns, actions)
              • Configuration (thresholds, settings)
```

| Content Type | Format | Why |
|-------------|--------|-----|
| Constitution (mission/vision/values) | `.md` | Referenced in CLAUDE.md, read by LLM |
| Quarterly objectives | `.md` | Strategy doc, human + LLM readable |
| Department context | `.md` | Part of agent prompt context |
| Team charter | `.md` | Referenced in factory CLAUDE.md |
| Workflow definitions | `.yaml` | Parsed by workflow engine |
| Compliance rules | `.yaml` | Parsed by validators |
| Event schemas | `.yaml` | Machine-validated structure |
| Configuration | `.yaml` | Settings with specific types |

**Example structure:**
```
.work/
├── governance/
│   ├── CONSTITUTION.md      # Mission, vision, values (prose)
│   ├── OBJECTIVES-Q1.md     # Quarterly strategy (prose)
│   └── compliance/
│       └── rules.yaml       # Machine-parsed rules
│
├── workflows/
│   └── code-review.yaml     # Parsed by executor
│
└── departments/
    └── engineering/
        └── CHARTER.md       # Department context (prose)
```

This follows the [Spec Kit](https://github.com/github/spec-kit) pattern: markdown as source of truth for AI-readable specifications.

### Cascading Mission/Vision

Mission and vision cascade through organizational levels. Each level inherits from above and adds specificity:

```
MISSION/VISION HIERARCHY

┌─────────────────────────────────────────────────────────────────────────────┐
│ COMPANY LEVEL (constitution.yaml)                                            │
│ Mission: "Why the organization exists"                                      │
│ Vision: "What success looks like"                                           │
│ Values: "How we behave"                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│ DEPARTMENT LEVEL (department CLAUDE.md)                                      │
│ Mission: How this department serves the company mission                     │
│ Vision: What excellence looks like for this department                      │
│ Focus: Department-specific values and priorities                            │
│                                                                             │
│ Example (Engineering):                                                      │
│   Mission: "Build reliable, maintainable systems that serve our users"     │
│   Vision: "Set the standard for code quality in our domain"                │
│   Focus: ["quality", "velocity", "documentation"]                          │
├─────────────────────────────────────────────────────────────────────────────┤
│ TEAM/FACTORY LEVEL (factory CLAUDE.md)                                       │
│ Mission: How this team serves the department mission                        │
│ Vision: What this team uniquely contributes                                 │
│ Focus: Team-specific expertise and responsibilities                         │
│                                                                             │
│ Example (API Team):                                                         │
│   Mission: "Provide stable, well-documented APIs for all consumers"        │
│   Vision: "Our APIs are the model for developer experience"                │
│   Focus: ["backwards compatibility", "documentation", "response times"]    │
└─────────────────────────────────────────────────────────────────────────────┘

INHERITANCE RULES:
- Lower levels MUST align with higher levels
- Lower levels add specificity, not contradiction
- Conflicts escalate upward for resolution
- Each level can only narrow scope, not expand it
```

### What MUST Be True (Invariants)

Certain things are **invariants** - they must ALWAYS be true:

```
ORGANIZATIONAL INVARIANTS

1. IDENTITY
   - Every action has a known actor
   - No anonymous work

2. ACCOUNTABILITY
   - Every work item has an owner
   - No orphaned tasks

3. TRACEABILITY
   - Every decision is logged
   - Can reconstruct "why" from history

4. COMPLIANCE
   - Rules are enforced, not suggested
   - Violations are blocked or escalated

5. CONTINUITY
   - Work survives agent failure
   - No single point of failure

6. QUALITY
   - Output meets standards before delivery
   - Verification is not optional
```

### Compliance Categories

| Category | Examples | Enforcement |
|----------|----------|-------------|
| **Security** | No secrets in code, access control | Hooks, pre-commit |
| **Quality** | Tests pass, code review done | Workflow gates |
| **Process** | PR before merge, approval required | Workflow steps |
| **Legal** | License headers, no copied code | Validators |
| **Business** | Follows spec, meets requirements | Verification |

### Work Routing: How Work Finds Workers

```
ROUTING DECISION TREE

New Work Order
     │
     ▼
┌─────────────┐
│ What type?  │
└─────┬───────┘
      │
      ├─── Bug ──────────────▶ Route to: Available worker, FIFO
      │
      ├─── Feature ──────────▶ Route to: Specialist or next available
      │
      ├─── Security ─────────▶ Route to: Security team, HIGH priority
      │
      └─── Research ─────────▶ Route to: Senior worker, TIME-BOXED
```

**Routing factors:**
1. **Type** - What kind of work?
2. **Priority** - How urgent?
3. **Expertise** - Who can do it?
4. **Availability** - Who's free?
5. **Load balancing** - Who has capacity?
6. **Affinity** - Who knows this area?

### Parallelization: Getting Speed

```
WORK SPLITTING STRATEGIES

STRATEGY 1: Task Parallelism
├── Break one large task into subtasks
├── Assign subtasks to different workers
└── Merge results when all complete

STRATEGY 2: Pipeline Parallelism
├── Work flows through stages
├── Each stage can process multiple items
└── Throughput = slowest stage

STRATEGY 3: Worker Pool
├── Pool of identical workers
├── Work queue with FIFO/priority ordering
├── Workers pull from queue
└── Scale by adding workers

STRATEGY 4: Specialization
├── Different workers for different work types
├── Route based on expertise
└── Optimize for quality over throughput
```

### Decision Making: Who Decides What

```
DECISION AUTHORITY MATRIX

Decision Type          │ Authority       │ Escalation
───────────────────────┼─────────────────┼──────────────
Technical approach     │ Worker          │ Supervisor
Priority change        │ Supervisor      │ CEO
Resource allocation    │ Supervisor      │ Operations
Cross-team dependency  │ Operations      │ CEO
Architecture change    │ CEO             │ Human
Compliance exception   │ BLOCKED         │ Human only
Mission change         │ BLOCKED         │ Human only
```

### Work Hierarchy (Complexity-Based)

Work is organized by **complexity**, not time. Larger work decomposes into smaller units:

```
WORK HIERARCHY

┌─────────────────────────────────────────────────────────────────────────────┐
│ EPIC (CEO / Human)                                                           │
│ ├── Large initiative spanning multiple sprints                              │
│ ├── Business-level goal or capability                                       │
│ ├── Decomposes into: Sprints or Stories                                     │
│ └── Example: "Implement user authentication system"                         │
├─────────────────────────────────────────────────────────────────────────────┤
│ SPRINT (Supervisor)                                                          │
│ ├── Coherent chunk of work toward an Epic                                   │
│ ├── Can be completed by a team in focused effort                            │
│ ├── Decomposes into: Stories                                                │
│ └── Example: "OAuth2 integration with Google"                               │
├─────────────────────────────────────────────────────────────────────────────┤
│ STORY (Supervisor / Worker)                                                  │
│ ├── User-visible feature or improvement                                     │
│ ├── Testable, demonstrable outcome                                          │
│ ├── Decomposes into: Tasks                                                  │
│ └── Example: "User can sign in with Google account"                         │
├─────────────────────────────────────────────────────────────────────────────┤
│ TASK (Worker)                                                                │
│ ├── Atomic unit of work                                                     │
│ ├── Single worker, single session                                           │
│ ├── Clear completion criteria                                               │
│ └── Example: "Add OAuth callback endpoint"                                  │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Decomposition rules:**
- Epic → 2-10 Sprints (or Stories if small)
- Sprint → 3-10 Stories
- Story → 1-5 Tasks
- Task → Should complete in one work session

**Why complexity over time:**
- Work takes as long as it takes
- Complexity is intrinsic; duration is emergent
- Enables parallel execution at any level
- No artificial time-boxing that doesn't fit the work

---

## Extensibility Model (Approved - Initial Scope)

Start with a constrained, safe model for plugins. Expand later based on experience.

### Core Principles

1. **Anthropic-approved only**: Plugins must come from Anthropic's curated registry
2. **Ephemeral by default**: Install for the work, remove when done
3. **Autonomous at lower scopes**: Worker/Team can self-serve
4. **Human gates at higher scopes**: Department/Company requires human approval

### Installation Scope Matrix

```
PLUGIN INSTALLATION AUTHORITY

┌──────────────┬─────────────────┬─────────────────┬─────────────────────────┐
│ Scope        │ Who Approves    │ Lifetime        │ Source Restriction      │
├──────────────┼─────────────────┼─────────────────┼─────────────────────────┤
│ Worker       │ Worker (self)   │ Until work done │ Anthropic registry only │
│ Team/Factory │ Supervisor      │ Until work done │ Anthropic registry only │
│ Department   │ Human required  │ Persistent      │ Any approved source     │
│ Company/Org  │ Human required  │ Persistent      │ Any approved source     │
└──────────────┴─────────────────┴─────────────────┴─────────────────────────┘
```

### Ephemeral Plugin Lifecycle

```
WORK-SCOPED PLUGIN LIFECYCLE

   Work Order Created
          │
          ▼
   ┌─────────────────┐
   │ Analyze work    │  ← What capabilities needed?
   │ requirements    │
   └────────┬────────┘
            │
            ▼
   ┌─────────────────┐
   │ Search Anthropic│  ← Only approved plugins
   │ registry        │
   └────────┬────────┘
            │
            ▼
   ┌─────────────────┐
   │ Install plugin  │  ← Scoped to worker/team
   │ (ephemeral)     │
   └────────┬────────┘
            │
            ▼
   ┌─────────────────┐
   │ Use plugin      │  ← During work execution
   │ tools           │
   └────────┬────────┘
            │
            ▼
   ┌─────────────────┐
   │ Work completed  │
   └────────┬────────┘
            │
            ▼
   ┌─────────────────┐
   │ Uninstall       │  ← Automatic cleanup
   │ plugin          │
   └─────────────────┘
```

### Why This Model

- **Trust**: Anthropic-approved plugins are vetted for safety and quality
- **Simplicity**: Workers don't need to evaluate arbitrary extensions
- **Clean state**: Ephemeral installation prevents plugin sprawl
- **Escalation path**: Frequently-used plugins can be promoted to persistent install with human approval

### Promotion Flow

```
EPHEMERAL → PERSISTENT PROMOTION

Worker uses plugin 10+ times
        │
        ▼
Supervisor notices pattern
        │
        ▼
Requests human approval for
Team/Department persistent install
        │
        ▼
Human approves (or rejects)
        │
        ▼
Plugin installed persistently
(available without per-work install)
```

### Future Expansion (Not Initial Scope)

These will be designed after initial plugin model is proven:
- **Skills**: Workflow automation (composite operations)
- **Experts**: Specialized agent configurations
- **Templates**: Best practice packages
- **Learning**: Knowledge promotion across scopes

See [EXTENSIBILITY.md](./EXTENSIBILITY.md) for exploration of these concepts.

---

## Approval Log

| Date | Section | Status |
|------|---------|--------|
| 2026-01-06 | Organizational Model | Approved |
| 2026-01-06 | LangGraph-Style Workflows | Approved |
| 2026-01-06 | Data Flow / Artifacts | Approved |
| 2026-01-07 | Governance Model (through Planning Cycles) | Approved |
| 2026-01-07 | Cascading Mission/Vision (Company → Dept → Team) | Approved |
| 2026-01-07 | File Format Convention (MD for prose, YAML for structures) | Approved |
| 2026-01-07 | Work Hierarchy (Epic → Sprint → Story → Task, complexity-based) | Approved |
| 2026-01-07 | Extensibility (Initial Scope - Anthropic plugins, ephemeral) | Approved |
