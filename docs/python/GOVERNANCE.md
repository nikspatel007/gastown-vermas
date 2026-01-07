# VerMAS Governance & Operations Design

> Ralph Wiggum iteration 1/7: How organizations govern work, ensure compliance, and achieve velocity

## The Questions

1. **How does work get done?** - From request to delivery
2. **What are compliance requirements?** - Rules that MUST be followed
3. **How do we route and split work?** - Parallelization for speed
4. **How are decisions made?** - Authority, escalation, consensus
5. **How does planning work?** - Quarterly, sprint, daily cycles
6. **What MUST be true?** - Mission, vision, invariants

---

## Iteration 1: Mapping Governance Concepts

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

### What MUST Be True (Invariants)

In a well-run organization, certain things are **invariants** - they must ALWAYS be true:

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

### Planning Cycles

```
PLANNING HIERARCHY

┌─────────────────────────────────────────────────────────────────────────────┐
│ QUARTERLY PLANNING (CEO)                                                     │
│ ├── Set objectives for quarter                                              │
│ ├── Allocate resources to objectives                                        │
│ └── Output: Quarterly plan document                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│ SPRINT PLANNING (Supervisor)                                                 │
│ ├── Break objectives into work orders                                       │
│ ├── Prioritize and sequence                                                 │
│ └── Output: Sprint backlog                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│ DAILY EXECUTION (Workers)                                                    │
│ ├── Pull from assignment                                                    │
│ ├── Execute and report                                                      │
│ └── Output: Completed work                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│ RETROSPECTIVE (All)                                                          │
│ ├── What worked? What didn't?                                               │
│ ├── Process improvements                                                    │
│ └── Output: Updated procedures                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Questions for Iteration 2

1. **How do we encode compliance rules?**
   - TOML/YAML config?
   - Code validators?
   - LLM-evaluated rules?

2. **What's the work routing algorithm?**
   - Simple round-robin?
   - Priority queue?
   - Skill-based routing?

3. **How do we handle blocked work?**
   - Dependencies
   - Waiting for external input
   - Stuck workers

4. **Where does planning state live?**
   - Quarterly objectives
   - Sprint backlogs
   - Resource allocation

---

## Iteration 1 Key Insights

1. **Governance is layered**: Board → Executive → Compliance → Planning → Operations → Execution

2. **Invariants are enforced, not hoped for**: Use hooks, validators, workflow gates

3. **Work routing needs multiple factors**: Type, priority, expertise, availability, affinity

4. **Planning is hierarchical**: Quarterly → Sprint → Daily

5. **Decision authority is explicit**: Clear matrix of who decides what

---

## Iteration 2: Compliance & Constraints

### Answering: How do we encode compliance rules?

**Answer: Three-tier enforcement system**

```
COMPLIANCE ENFORCEMENT TIERS

TIER 1: HARD BLOCKS (Cannot proceed)
├── Mechanism: Hooks, pre-commit, validators
├── Examples:
│   ├── No secrets in code (regex scan)
│   ├── Tests must pass before merge
│   ├── Required approvals not obtained
│   └── License violations detected
└── Response: Block action, return error

TIER 2: SOFT BLOCKS (Can override with justification)
├── Mechanism: Workflow gates with escape hatch
├── Examples:
│   ├── Code coverage below threshold
│   ├── Documentation missing
│   ├── Style violations
│   └── Deprecated API usage
└── Response: Warn, require acknowledgment to proceed

TIER 3: ADVISORY (Log and continue)
├── Mechanism: Post-action validators, auditors
├── Examples:
│   ├── Performance regression detected
│   ├── Unusual access pattern
│   ├── Large file committed
│   └── Off-hours activity
└── Response: Log event, may trigger review
```

### Compliance Rule Definition

```yaml
# .work/compliance/rules.yaml

rules:
  - id: SEC-001
    name: no-secrets-in-code
    tier: hard
    description: "Prevent secrets from being committed"
    check:
      type: regex
      patterns:
        - "(?i)(api[_-]?key|secret|password)\\s*=\\s*['\"][^'\"]+['\"]"
        - "(?i)Bearer\\s+[A-Za-z0-9-_=]+"
    applies_to:
      - "*.py"
      - "*.js"
      - "*.yaml"
    exclude:
      - "*_test.py"
      - "*.example.*"

  - id: QA-001
    name: tests-must-pass
    tier: hard
    description: "All tests must pass before merge"
    check:
      type: command
      command: "pytest --tb=short"
      success_exit_code: 0

  - id: DOC-001
    name: readme-required
    tier: soft
    description: "README must exist for new modules"
    check:
      type: file_exists
      pattern: "**/README.md"
    escape:
      requires: supervisor_approval
      justification_required: true

  - id: PERF-001
    name: no-performance-regression
    tier: advisory
    description: "Alert on performance regression"
    check:
      type: benchmark
      threshold: 10  # percent slower
    action:
      type: log_event
      event_type: compliance.advisory
```

### Enforcement Points

```
WHERE COMPLIANCE IS CHECKED

┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│   Work Created          Work In Progress          Work Complete            │
│        │                      │                        │                   │
│        ▼                      ▼                        ▼                   │
│   ┌─────────┐           ┌─────────┐              ┌─────────┐              │
│   │ CREATE  │           │ MODIFY  │              │ SUBMIT  │              │
│   │ HOOKS   │           │ HOOKS   │              │ HOOKS   │              │
│   └────┬────┘           └────┬────┘              └────┬────┘              │
│        │                     │                        │                   │
│   - Valid work order    - File changes valid     - All tests pass        │
│   - Required fields     - No secrets             - Required approvals    │
│   - Assignee exists     - Style OK               - Verification done     │
│                                                                             │
│        │                      │                        │                   │
│        ▼                      ▼                        ▼                   │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                     CONTINUOUS MONITORING                           │  │
│   │   - Audit log analysis                                              │  │
│   │   - Anomaly detection                                               │  │
│   │   - Compliance dashboard                                            │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Compliance State Machine

```
WORK ORDER COMPLIANCE STATES

                    ┌───────────────┐
                    │   PENDING     │
                    │  COMPLIANCE   │
                    └───────┬───────┘
                            │
              ┌─────────────┼─────────────┐
              │             │             │
              ▼             ▼             ▼
        ┌──────────┐  ┌──────────┐  ┌──────────┐
        │  PASSED  │  │  FAILED  │  │  WAIVED  │
        └────┬─────┘  └────┬─────┘  └────┬─────┘
             │             │             │
             │             │             │
             ▼             ▼             ▼
        Proceed       Block +        Proceed +
        normally      require        log waiver
                      remediation    + auditor
```

### Constraint Types

| Type | Description | Example | Enforcement |
|------|-------------|---------|-------------|
| **Temporal** | Time-based rules | No deploys on Friday | Schedule check |
| **Capacity** | Resource limits | Max 5 concurrent workers | Counter check |
| **Sequential** | Order requirements | Review before merge | Workflow gate |
| **Approval** | Human sign-off needed | Security review | Human node |
| **Quality** | Metric thresholds | 80% test coverage | Metric check |
| **Structural** | Format requirements | PR description template | Schema validation |

### Escape Hatches

Not all rules are absolute. Some can be overridden:

```yaml
# Escape hatch definition
escape:
  - rule_id: DOC-001
    conditions:
      - type: approval
        approver: supervisor
      - type: justification
        min_length: 50
      - type: time_limit
        expires_in: 24h  # Override expires

    audit:
      - log_event: compliance.waived
      - notify: compliance_team
      - review_required: true  # Must be reviewed within 7 days
```

### Compliance Events

```python
# Events emitted for compliance

class ComplianceEvent(BaseModel):
    event_type: Literal[
        "compliance.check_passed",
        "compliance.check_failed",
        "compliance.waived",
        "compliance.violation_detected",
        "compliance.remediation_started",
        "compliance.remediation_completed",
    ]
    rule_id: str
    work_order_id: str
    actor: str
    details: dict
    waiver: Optional[WaiverInfo]
```

---

## Questions for Iteration 3

1. **What's the concrete work routing algorithm?**
   - Priority queue implementation
   - Load balancing strategy
   - Affinity/expertise matching

2. **How do we split large work items?**
   - Automatic decomposition?
   - Manual breakdown required?
   - Dependency management

3. **How do workers signal capacity?**
   - Pull-based (worker requests work)
   - Push-based (supervisor assigns)
   - Hybrid approach

---

## Iteration 2 Key Insights

1. **Three-tier enforcement**: Hard blocks, soft blocks (with escape), advisory

2. **Rules are declarative**: YAML/TOML config, not hardcoded

3. **Multiple enforcement points**: Create, modify, submit, continuous

4. **Escape hatches exist**: But require approval, justification, audit

5. **Everything is logged**: Compliance events are first-class citizens

---

## Iteration 3: Work Routing & Parallelization

### The Routing Problem

How do we get work to the right worker at the right time for maximum throughput?

```
ROUTING DIMENSIONS

                    ┌─────────────────────────────────────────┐
                    │           INCOMING WORK                  │
                    └────────────────┬────────────────────────┘
                                     │
        ┌────────────────────────────┼────────────────────────────┐
        │                            │                            │
        ▼                            ▼                            ▼
   ┌─────────┐                 ┌─────────┐                 ┌─────────┐
   │ URGENCY │                 │  TYPE   │                 │  SIZE   │
   │         │                 │         │                 │         │
   │ P0-P4   │                 │ bug     │                 │ small   │
   │         │                 │ feature │                 │ medium  │
   │         │                 │ research│                 │ large   │
   └────┬────┘                 └────┬────┘                 └────┬────┘
        │                           │                           │
        └───────────────────────────┼───────────────────────────┘
                                    │
                                    ▼
                           ┌───────────────┐
                           │ ROUTING       │
                           │ ALGORITHM     │
                           └───────┬───────┘
                                   │
        ┌──────────────────────────┼──────────────────────────┐
        │                          │                          │
        ▼                          ▼                          ▼
   ┌─────────┐                ┌─────────┐                ┌─────────┐
   │ Worker  │                │ Worker  │                │ Worker  │
   │  Slot 0 │                │  Slot 1 │                │  Slot 2 │
   └─────────┘                └─────────┘                └─────────┘
```

### Routing Algorithm: Weighted Priority Queue

```python
# vermas/routing/algorithm.py

from dataclasses import dataclass
from typing import List, Optional
import heapq

@dataclass
class WorkItem:
    id: str
    priority: int          # 0 (highest) to 4 (lowest)
    work_type: str         # bug, feature, research, etc.
    estimated_size: str    # small, medium, large
    required_skills: List[str]
    affinity_hints: List[str]  # Previous workers, related code areas
    created_at: datetime
    deadline: Optional[datetime]

@dataclass
class Worker:
    id: str
    skills: List[str]
    current_load: int      # 0-100
    affinity_areas: List[str]
    available: bool

def calculate_routing_score(work: WorkItem, worker: Worker) -> float:
    """
    Higher score = better match

    Factors:
    - Priority weight (urgent work gets priority)
    - Skill match (does worker have required skills?)
    - Load balance (prefer less loaded workers)
    - Affinity bonus (has worker touched this area before?)
    - Age penalty (older work gets priority boost)
    """
    score = 0.0

    # Priority: P0 = 100, P1 = 80, P2 = 60, P3 = 40, P4 = 20
    priority_scores = {0: 100, 1: 80, 2: 60, 3: 40, 4: 20}
    score += priority_scores.get(work.priority, 20)

    # Skill match: +50 for perfect match, +25 for partial
    skill_match = len(set(work.required_skills) & set(worker.skills))
    if skill_match == len(work.required_skills):
        score += 50
    elif skill_match > 0:
        score += 25

    # Load balance: Prefer workers with lower load
    score += (100 - worker.current_load) * 0.3

    # Affinity: +20 if worker has touched related code
    if set(work.affinity_hints) & set(worker.affinity_areas):
        score += 20

    # Age: +1 per hour waiting (max +24)
    age_hours = (datetime.now() - work.created_at).total_seconds() / 3600
    score += min(age_hours, 24)

    # Deadline urgency: +50 if deadline within 24h
    if work.deadline:
        hours_to_deadline = (work.deadline - datetime.now()).total_seconds() / 3600
        if hours_to_deadline < 24:
            score += 50
        elif hours_to_deadline < 72:
            score += 25

    return score

def route_work(work_queue: List[WorkItem], workers: List[Worker]) -> List[tuple]:
    """
    Match work items to workers using weighted scoring.
    Returns list of (work_item, worker) assignments.
    """
    available_workers = [w for w in workers if w.available]
    assignments = []

    for work in sorted(work_queue, key=lambda w: w.priority):
        if not available_workers:
            break

        # Score each worker for this work item
        scores = [(calculate_routing_score(work, w), w) for w in available_workers]
        scores.sort(reverse=True)

        best_worker = scores[0][1]
        assignments.append((work, best_worker))
        available_workers.remove(best_worker)

    return assignments
```

### Work Splitting: Breaking Large Items

```
WORK DECOMPOSITION STRATEGIES

STRATEGY 1: MANUAL BREAKDOWN (default)
┌─────────────────────────────────────────────────────────────────────────────┐
│ Large Feature Request                                                        │
│ "Implement user authentication system"                                       │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼ (Supervisor breaks down)
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ Design auth     │  │ Implement       │  │ Write tests     │
│ schema          │  │ login/logout    │  │                 │
└────────┬────────┘  └────────┬────────┘  └────────┬────────┘
         │                    │                    │
         │        depends on ─┘                    │
         └──────────────────────────── depends on ─┘

STRATEGY 2: TEMPLATE-BASED DECOMPOSITION
┌─────────────────────────────────────────────────────────────────────────────┐
│ Work order type: "feature"                                                   │
│ Template: feature-development                                                │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼ (Auto-generate from template)
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ Step 1: Design  │  │ Step 2: Impl    │  │ Step 3: Test    │
│ (auto-created)  │  │ (auto-created)  │  │ (auto-created)  │
└─────────────────┘  └─────────────────┘  └─────────────────┘

STRATEGY 3: LLM-ASSISTED DECOMPOSITION
┌─────────────────────────────────────────────────────────────────────────────┐
│ Large work order detected (>4 hours estimated)                               │
│ Trigger: LLM decomposition agent                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼ (LLM suggests breakdown)
┌─────────────────────────────────────────────────────────────────────────────┐
│ Decomposition proposal:                                                      │
│ 1. Database schema changes (2h)                                              │
│ 2. API endpoint implementation (3h)                                          │
│ 3. Frontend integration (2h)                                                 │
│ 4. Testing and documentation (1h)                                            │
│                                                                              │
│ Dependencies: 2→1, 3→2, 4→3                                                  │
│                                                                              │
│ [Accept] [Modify] [Reject - do manually]                                     │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Worker Capacity Signaling

**Answer: Hybrid pull/push with capacity limits**

```
CAPACITY MODEL

┌─────────────────────────────────────────────────────────────────────────────┐
│                           WORKER STATE                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   State          │ Can Accept Work? │ Behavior                              │
│   ───────────────┼──────────────────┼─────────────────────────────────────  │
│   IDLE           │ Yes              │ Pull from queue or receive push       │
│   WORKING        │ No               │ Executing assigned work               │
│   BLOCKED        │ No               │ Waiting for dependency/input          │
│   OVERLOADED     │ No               │ Above capacity threshold              │
│   OFFLINE        │ No               │ Session ended or crashed              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

CAPACITY CALCULATION:

Worker Capacity = MAX_CONCURRENT_ITEMS - current_in_progress_count

Example:
- MAX_CONCURRENT_ITEMS = 1 (for focused work)
- current_in_progress_count = 0
- Capacity = 1 (can accept one work item)
```

### Parallel Execution Patterns

```yaml
# .work/templates/parallel-ci.yaml

name: parallel-ci-pipeline
description: Run CI checks in parallel

nodes:
  start:
    type: entry
    next: fan_out

  fan_out:
    type: parallel
    parallel_next:
      - run_tests
      - run_lint
      - run_security_scan
      - run_type_check

  run_tests:
    type: action
    agent: ci_worker
    command: "pytest"
    next: join

  run_lint:
    type: action
    agent: ci_worker
    command: "ruff check ."
    next: join

  run_security_scan:
    type: action
    agent: ci_worker
    command: "bandit -r src/"
    next: join

  run_type_check:
    type: action
    agent: ci_worker
    command: "mypy src/"
    next: join

  join:
    type: join
    join_from: [run_tests, run_lint, run_security_scan, run_type_check]
    strategy: all_must_pass  # or any_pass, majority_pass
    next: check_results

  check_results:
    type: condition
    expression: "all_passed == true"
    true_next: success
    false_next: failure

  success:
    type: exit
    status: completed

  failure:
    type: exit
    status: failed
```

### Throughput Optimization

```
BOTTLENECK ANALYSIS

                    Input Rate: 10 work orders / hour
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         PIPELINE STAGES                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Stage           │ Workers │ Throughput  │ Bottleneck?                     │
│   ────────────────┼─────────┼─────────────┼───────────────────────────────  │
│   Triage          │ 1       │ 20/hr       │ No                              │
│   Implementation  │ 3       │ 6/hr        │ YES (3 × 2/hr each)             │
│   Code Review     │ 1       │ 12/hr       │ No                              │
│   QA              │ 1       │ 8/hr        │ Maybe (close to limit)          │
│                                                                             │
│   RECOMMENDATION: Add more implementation workers                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

SCALING RULES:

1. Monitor queue depth per stage
2. If queue depth > threshold for > 15 min:
   - Alert supervisor
   - Suggest scaling up workers for that stage
3. If workers idle > threshold:
   - Suggest scaling down
   - Reassign to bottleneck stage if skills match
```

---

## Questions for Iteration 4

1. **Who has authority to make what decisions?**
   - Priority changes
   - Resource reallocation
   - Scope changes

2. **How do escalations work?**
   - When to escalate
   - Escalation paths
   - Timeout triggers

3. **How do we handle conflicts?**
   - Resource contention
   - Priority disputes
   - Cross-team dependencies

---

## Iteration 3 Key Insights

1. **Weighted scoring for routing**: Priority + Skills + Load + Affinity + Age

2. **Three decomposition strategies**: Manual, template-based, LLM-assisted

3. **Hybrid capacity model**: Workers signal state, supervisor respects limits

4. **Parallel patterns in workflows**: Fan-out/join with configurable strategies

5. **Monitor for bottlenecks**: Queue depth → scaling recommendations

---

## Iteration 4: Decision Making & Escalation

### Decision Authority Matrix

Who can decide what? This must be explicit.

```
DECISION AUTHORITY LEVELS

┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│   Level    │ Role           │ Can Decide                                    │
│   ─────────┼────────────────┼─────────────────────────────────────────────  │
│   L0       │ Human          │ Everything (override any agent decision)      │
│   L1       │ CEO            │ Strategy, cross-factory, architecture         │
│   L2       │ Operations     │ Resource allocation, scaling, emergencies     │
│   L3       │ Supervisor     │ Priority, assignment, within-factory routing  │
│   L4       │ QA             │ Quality gates, verification thresholds        │
│   L5       │ Worker         │ Technical approach within assigned scope      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Decision Types and Authority

```yaml
# .work/governance/decisions.yaml

decisions:
  # Technical decisions
  - type: technical_approach
    description: "How to implement a feature"
    authority: worker
    escalate_if:
      - condition: affects_architecture
        escalate_to: ceo
      - condition: cross_team_dependency
        escalate_to: supervisor

  - type: code_merge
    description: "Merge code to main branch"
    authority: qa
    requires:
      - tests_pass: true
      - review_approved: true
      - compliance_passed: true

  # Operational decisions
  - type: priority_change
    description: "Change work order priority"
    authority: supervisor
    escalate_if:
      - condition: priority_to_p0
        escalate_to: operations
    audit: true

  - type: resource_allocation
    description: "Add/remove workers"
    authority: operations
    escalate_if:
      - condition: exceeds_budget
        escalate_to: ceo
    requires:
      - justification: true

  # Strategic decisions
  - type: scope_change
    description: "Change project scope"
    authority: ceo
    requires:
      - impact_analysis: true
      - stakeholder_notification: true

  # Blocked decisions (human only)
  - type: compliance_exception
    description: "Override compliance rule"
    authority: human_only
    audit: mandatory

  - type: mission_change
    description: "Change company mission/vision"
    authority: human_only
    audit: mandatory
```

### Escalation Paths

```
ESCALATION FLOWCHART

Worker stuck on technical issue
         │
         ▼
    Can resolve in 30 min?
         │
    ┌────┴────┐
    │ Yes     │ No
    ▼         ▼
 Continue   Escalate to Supervisor
              │
              ▼
         Supervisor can resolve?
              │
         ┌────┴────┐
         │ Yes     │ No
         ▼         ▼
      Resolve   Escalate to Operations/CEO
                  │
                  ▼
             Operations can resolve?
                  │
             ┌────┴────┐
             │ Yes     │ No
             ▼         ▼
          Resolve   Escalate to Human
                      │
                      ▼
                 Human decides
```

### Escalation Triggers

```python
# vermas/governance/escalation.py

from enum import Enum
from datetime import timedelta

class EscalationTrigger(Enum):
    TIMEOUT = "timeout"           # Work stuck too long
    BLOCKED = "blocked"           # Dependency not resolved
    CONFLICT = "conflict"         # Resource contention
    FAILURE = "failure"           # Repeated failures
    THRESHOLD = "threshold"       # Metric exceeded threshold
    EXPLICIT = "explicit"         # Agent requested escalation

class EscalationRule:
    trigger: EscalationTrigger
    condition: str
    timeout: timedelta
    escalate_to: str
    auto_escalate: bool

ESCALATION_RULES = [
    EscalationRule(
        trigger=EscalationTrigger.TIMEOUT,
        condition="work_order.status == 'in_progress'",
        timeout=timedelta(hours=4),
        escalate_to="supervisor",
        auto_escalate=True,
    ),
    EscalationRule(
        trigger=EscalationTrigger.TIMEOUT,
        condition="work_order.priority == 0",  # P0
        timeout=timedelta(minutes=30),
        escalate_to="operations",
        auto_escalate=True,
    ),
    EscalationRule(
        trigger=EscalationTrigger.BLOCKED,
        condition="blocked_by.status == 'open'",
        timeout=timedelta(hours=2),
        escalate_to="supervisor",
        auto_escalate=True,
    ),
    EscalationRule(
        trigger=EscalationTrigger.FAILURE,
        condition="failure_count >= 3",
        timeout=timedelta(minutes=0),  # Immediate
        escalate_to="supervisor",
        auto_escalate=True,
    ),
]
```

### Conflict Resolution

```
CONFLICT TYPES AND RESOLUTION

┌─────────────────────────────────────────────────────────────────────────────┐
│ CONFLICT TYPE        │ RESOLUTION STRATEGY                                  │
├──────────────────────┼──────────────────────────────────────────────────────┤
│                      │                                                      │
│ Resource Contention  │ 1. Higher priority wins                              │
│ (same resource,      │ 2. Earlier request wins (if same priority)           │
│  multiple requests)  │ 3. Escalate to supervisor if still tied              │
│                      │                                                      │
├──────────────────────┼──────────────────────────────────────────────────────┤
│                      │                                                      │
│ Priority Dispute     │ 1. CEO sets overall priority                         │
│ (disagreement on     │ 2. Supervisor resolves within factory                │
│  importance)         │ 3. Operations resolves cross-factory                 │
│                      │                                                      │
├──────────────────────┼──────────────────────────────────────────────────────┤
│                      │                                                      │
│ Technical Approach   │ 1. Worker proposes, supervisor approves              │
│ (disagreement on     │ 2. If impasse, CEO decides                           │
│  how to implement)   │ 3. Document decision for future reference            │
│                      │                                                      │
├──────────────────────┼──────────────────────────────────────────────────────┤
│                      │                                                      │
│ Cross-Team Deps      │ 1. Supervisors negotiate directly                    │
│ (team A needs        │ 2. Operations arbitrates if no agreement             │
│  something from B)   │ 3. CEO escalation for strategic conflicts            │
│                      │                                                      │
├──────────────────────┼──────────────────────────────────────────────────────┤
│                      │                                                      │
│ Scope Creep          │ 1. Supervisor gates scope changes                    │
│ (work growing        │ 2. Requires explicit approval for expansion          │
│  beyond original)    │ 3. May trigger new work order instead                │
│                      │                                                      │
└──────────────────────┴──────────────────────────────────────────────────────┘
```

### Decision Events

```python
# All decisions are logged as events

class DecisionEvent(BaseModel):
    event_type: Literal[
        "decision.made",
        "decision.escalated",
        "decision.overridden",
        "decision.delegated",
    ]
    decision_type: str          # From decisions.yaml
    actor: str                  # Who made the decision
    authority_level: int        # L0-L5
    context: dict               # What was decided
    justification: Optional[str]
    escalated_from: Optional[str]
    overrides: Optional[str]    # Previous decision if override
```

### Delegation and Proxy

```
DELEGATION MODEL

CEO can delegate to:
├── Operations (resource decisions)
├── Supervisor (priority decisions)
└── QA (quality decisions)

Delegation rules:
1. Delegation must be explicit (logged)
2. Delegatee cannot further delegate
3. Delegator retains override authority
4. Delegation can be time-limited
5. Delegation can be scope-limited

Example:
┌─────────────────────────────────────────────────────────────────────────────┐
│ DELEGATION RECORD                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│ From: ceo                                                                    │
│ To: alpha/supervisor                                                         │
│ Decision types: [priority_change, resource_allocation]                       │
│ Scope: factory=alpha                                                         │
│ Expires: 2026-01-31                                                          │
│ Justification: "Supervisor has context for Q1 sprint"                        │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Questions for Iteration 5

1. **How does quarterly planning work?**
   - OKR-style objectives
   - Sprint planning
   - Capacity allocation

2. **How do we track progress against plans?**
   - Burndown/burnup
   - Velocity metrics
   - Forecast adjustments

3. **What happens when plans change?**
   - Re-planning triggers
   - Impact analysis
   - Communication

---

## Iteration 4 Key Insights

1. **Authority is explicit**: Decision matrix defines who decides what

2. **Escalation is automatic**: Timeouts and thresholds trigger escalation

3. **Conflicts have resolution strategies**: Each conflict type has a playbook

4. **Everything is logged**: Decision events create audit trail

5. **Delegation is formal**: Explicit, scoped, time-limited

---

## Iteration 5: Work Hierarchy (Complexity-Based)

> **Design Decision**: Work is organized by **complexity**, not time. We use Epic → Sprint → Story → Task hierarchy instead of annual/quarterly/daily time-boxes.

### Why Complexity Over Time

```
TIME-BASED (rejected)              COMPLEXITY-BASED (adopted)
─────────────────────              ─────────────────────────
Annual plan → Quarterly → Sprint   Epic → Sprint → Story → Task
├── Artificial deadlines          ├── Natural decomposition
├── Work crammed to fit boxes     ├── Work sized by complexity
├── "Deadline-driven" culture     ├── "Completion-driven" culture
└── Often mismatches reality      └── Duration is emergent
```

**Key insight**: Complexity is intrinsic; duration is emergent. You can estimate complexity upfront, but duration depends on unknowns discovered during execution.

### Work Hierarchy

```
WORK DECOMPOSITION HIERARCHY

┌─────────────────────────────────────────────────────────────────────────────┐
│ EPIC (CEO / Human)                                                           │
│                                                                             │
│ Definition: Large initiative, business-level goal                           │
│ Size: Decomposes into 2-10 Sprints (or Stories if small)                   │
│ Owner: CEO or Human creates and owns                                        │
│                                                                             │
│ Example: "Implement user authentication system"                             │
│          "Build Python VerMAS implementation"                               │
├─────────────────────────────────────────────────────────────────────────────┤
│ SPRINT (Supervisor)                                                          │
│                                                                             │
│ Definition: Coherent chunk of work toward an Epic                           │
│ Size: Decomposes into 3-10 Stories                                          │
│ Owner: Supervisor decomposes and tracks                                     │
│                                                                             │
│ Example: "OAuth2 integration with Google"                                   │
│          "Implement Layer 0: Identity System"                               │
├─────────────────────────────────────────────────────────────────────────────┤
│ STORY (Supervisor / Worker)                                                  │
│                                                                             │
│ Definition: User-visible feature or improvement                             │
│ Size: Decomposes into 1-5 Tasks                                             │
│ Owner: Worker owns execution                                                │
│ Criteria: Testable, demonstrable outcome                                    │
│                                                                             │
│ Example: "User can sign in with Google account"                             │
│          "AGENT_ID generation and validation works"                         │
├─────────────────────────────────────────────────────────────────────────────┤
│ TASK (Worker)                                                                │
│                                                                             │
│ Definition: Atomic unit of work                                             │
│ Size: Single worker, single session                                         │
│ Owner: Worker executes                                                      │
│ Criteria: Clear completion criteria, no decomposition needed                │
│                                                                             │
│ Example: "Add OAuth callback endpoint"                                      │
│          "Write test_agent_id.py"                                           │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Decomposition Rules

| From | To | Count | Guidance |
|------|----|-------|----------|
| Epic | Sprints | 2-10 | Each Sprint is a coherent deliverable |
| Sprint | Stories | 3-10 | Each Story is user-visible |
| Story | Tasks | 1-5 | Each Task is atomic |

**Size heuristics:**
- Task: Should complete in one work session (single context)
- Story: Should be demonstrable to a user
- Sprint: Should deliver measurable progress toward Epic
- Epic: Should be a complete capability when done

### Work Hierarchy in YAML

```yaml
# .work/epics/auth-system.yaml

epic:
  id: epic-auth-001
  title: "Implement user authentication system"
  owner: ceo
  status: in_progress

  # Success criteria (how we know Epic is done)
  done_when:
    - "Users can register and log in"
    - "OAuth2 with Google and GitHub works"
    - "Password reset flow implemented"
    - "All auth endpoints have >90% test coverage"

sprints:
  - id: sprint-auth-001
    title: "OAuth2 integration with Google"
    status: in_progress
    stories:
      - id: story-001
        title: "User can sign in with Google"
        status: in_progress
        tasks:
          - id: task-001
            title: "Add Google OAuth client config"
            status: completed
          - id: task-002
            title: "Implement OAuth callback endpoint"
            status: in_progress
          - id: task-003
            title: "Create user from OAuth profile"
            status: pending
      - id: story-002
        title: "User can link Google to existing account"
        status: pending

  - id: sprint-auth-002
    title: "Password-based authentication"
    status: pending
```

### Progress Tracking

```
PROGRESS ROLLUP

Task completion → Story progress → Sprint progress → Epic progress

Example:
┌─────────────────────────────────────────────────────────────────────────────┐
│ Epic: Implement authentication (40% complete)                                │
│                                                                             │
│ ├── Sprint 1: OAuth2 with Google (75% complete)                            │
│ │   ├── Story: Sign in with Google (66% - 2/3 tasks done)                  │
│ │   └── Story: Link Google account (0% - not started)                      │
│ │                                                                           │
│ ├── Sprint 2: Password auth (0% - not started)                             │
│ │   ├── Story: User registration                                          │
│ │   ├── Story: User login                                                  │
│ │   └── Story: Password reset                                              │
│ │                                                                           │
│ └── Sprint 3: GitHub OAuth (0% - not started)                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Progress Tracking

```python
# vermas/planning/metrics.py

@dataclass
class SprintMetrics:
    sprint_id: str
    committed_points: int
    completed_points: int
    velocity: float  # completed / committed

    # Burndown
    ideal_burndown: List[int]  # Daily ideal remaining
    actual_burndown: List[int]  # Daily actual remaining

    # Health indicators
    scope_changes: int  # Work added/removed mid-sprint
    blockers_encountered: int
    blockers_resolved: int

def calculate_velocity(sprints: List[SprintMetrics], window: int = 3) -> float:
    """Calculate rolling average velocity over last N sprints."""
    recent = sprints[-window:]
    if not recent:
        return 1.0  # Default assumption
    return sum(s.velocity for s in recent) / len(recent)

def forecast_completion(
    remaining_points: int,
    velocity: float,
    sprint_capacity: int
) -> int:
    """Forecast number of sprints to complete remaining work."""
    points_per_sprint = sprint_capacity * velocity
    if points_per_sprint <= 0:
        return float('inf')
    return math.ceil(remaining_points / points_per_sprint)
```

### Progress Dashboard

```
QUARTERLY PROGRESS DASHBOARD

Quarter: Q1-2026                                     Week 2 of 12

┌─────────────────────────────────────────────────────────────────────────────┐
│ OBJECTIVE: Launch Python implementation of VerMAS                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ KR-001: MVP complete with Layers 0-4                                        │
│ ████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  25%                  │
│ Target: 100% | Current: 25% | On Track: ✓                                   │
│                                                                             │
│ KR-002: 3 factories running in production                                   │
│ ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  0%                   │
│ Target: 3 | Current: 0 | On Track: ⚠️ (behind schedule)                     │
│                                                                             │
│ KR-003: Documentation coverage > 80%                                        │
│ ████████████████████████░░░░░░░░░░░░░░░░░░░░░░░░░░░░  50%                  │
│ Target: 80% | Current: 50% | On Track: ✓                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

CURRENT SPRINT: sprint-2026-w01

┌─────────────────────────────────────────────────────────────────────────────┐
│ Burndown                                                                     │
│                                                                             │
│ Points │                                                                    │
│   60   │ ●                                                                  │
│   50   │   ╲  ●                                                             │
│   40   │     ╲   ●                                                          │
│   30   │       ╲    ●                                                       │
│   20   │         ╲     ●                                                    │
│   10   │           ╲      ●                                                 │
│    0   │─────────────╲───────────────────────────                           │
│        │  D1  D2  D3  D4  D5  D6  D7  D8  D9  D10                          │
│                                                                             │
│ ── Ideal    ● Actual                                                        │
│                                                                             │
│ Velocity: 0.85 (3-sprint avg)                                               │
│ Forecast: On track to complete sprint commitment                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Re-Planning Triggers

```yaml
# .work/governance/replan-triggers.yaml

triggers:
  - id: scope_creep
    condition: "sprint.added_points > sprint.committed_points * 0.2"
    action: supervisor_review
    message: "Sprint scope increased by >20%"

  - id: velocity_drop
    condition: "velocity < historical_velocity * 0.7"
    action: retrospective
    message: "Velocity dropped significantly"

  - id: blocker_pileup
    condition: "active_blockers > 3"
    action: escalate_to_operations
    message: "Multiple blockers need resolution"

  - id: objective_at_risk
    condition: "forecast_completion > quarter_remaining_weeks"
    action: ceo_review
    message: "Objective may not be achievable this quarter"

  - id: resource_change
    condition: "available_workers != planned_workers"
    action: capacity_replan
    message: "Team capacity changed"
```

### Plan Change Protocol

```
PLAN CHANGE WORKFLOW

┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. DETECT CHANGE NEED                                                        │
│    - Trigger fires (scope creep, velocity drop, etc.)                        │
│    - Explicit request from stakeholder                                       │
│    - New information (dependency discovered)                                 │
└────────────────────────────────────────────────────────────────────────────┬┘
                                                                              │
┌─────────────────────────────────────────────────────────────────────────────▼┐
│ 2. IMPACT ANALYSIS                                                           │
│    - What objectives are affected?                                           │
│    - What's the capacity impact?                                             │
│    - What dependencies change?                                               │
│    - What's the timeline impact?                                             │
└────────────────────────────────────────────────────────────────────────────┬┘
                                                                              │
┌─────────────────────────────────────────────────────────────────────────────▼┐
│ 3. PROPOSE NEW PLAN                                                          │
│    - Revised commitments                                                     │
│    - Trade-off options (cut scope, extend timeline, add resources)           │
│    - Recommendation                                                          │
└────────────────────────────────────────────────────────────────────────────┬┘
                                                                              │
┌─────────────────────────────────────────────────────────────────────────────▼┐
│ 4. APPROVAL                                                                  │
│    - Sprint change: Supervisor approves                                      │
│    - Quarterly objective change: CEO approves                                │
│    - Annual plan change: Human approves                                      │
└────────────────────────────────────────────────────────────────────────────┬┘
                                                                              │
┌─────────────────────────────────────────────────────────────────────────────▼┐
│ 5. COMMUNICATE                                                               │
│    - Notify affected parties                                                 │
│    - Update planning documents                                               │
│    - Log decision event                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Questions for Iteration 6

1. **How do we encode mission/vision?**
   - Where does it live?
   - How is it referenced?
   - How is it enforced?

2. **What are the invariants agents must preserve?**
   - Non-negotiable rules
   - Self-correction mechanisms

3. **How do we prevent drift from objectives?**
   - Alignment checks
   - Objective-aware routing

---

## Iteration 5 Key Insights

1. **Work hierarchy is complexity-based**: Epic → Sprint → Story → Task (not time-based)

2. **OKR-style objectives**: Clear targets with measurable key results

3. **Decomposition has bounds**: Epic (2-10 Sprints), Sprint (3-10 Stories), Story (1-5 Tasks)

4. **Task is atomic**: Single worker, single session, clear completion criteria

5. **Duration is emergent**: Complexity determines size, time follows naturally

---

## Iteration 6: Mission, Vision & Invariants

### Where Mission/Vision Lives

```
MISSION/VISION HIERARCHY

┌─────────────────────────────────────────────────────────────────────────────┐
│                       COMPANY CONSTITUTION                                   │
│                   .work/governance/constitution.yaml                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   mission: "Build reliable multi-agent systems that augment human work"     │
│                                                                             │
│   vision: "A world where AI agents and humans collaborate seamlessly"       │
│                                                                             │
│   values:                                                                   │
│     - name: quality                                                         │
│       description: "Ship verified, tested work"                             │
│     - name: transparency                                                    │
│       description: "All decisions logged, all work visible"                 │
│     - name: autonomy                                                        │
│       description: "Agents act independently within bounds"                 │
│     - name: human_oversight                                                 │
│       description: "Humans can always intervene"                            │
│                                                                             │
│   invariants:                                                               │
│     - "Every action has a known actor"                                      │
│     - "Every work item has an owner"                                        │
│     - "Compliance rules cannot be disabled by agents"                       │
│     - "Human can override any agent decision"                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                    References and enforces
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         AGENT PROFILES                                       │
│                     ~/.claude/profiles/*/CLAUDE.md                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Each agent's CLAUDE.md includes:                                          │
│   - Reference to constitution                                               │
│   - Role-specific interpretation                                            │
│   - Behavioral constraints                                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Constitution Schema

```yaml
# .work/governance/constitution.yaml

version: 1
last_updated: 2026-01-01
updated_by: human  # Only human can update

mission:
  statement: "Build reliable multi-agent systems that augment human work"
  rationale: |
    We exist to prove that AI agents can work alongside humans productively,
    handling routine work while humans focus on creative and strategic tasks.

vision:
  statement: "A world where AI agents and humans collaborate seamlessly"
  horizon: 5_years
  milestones:
    - year: 2026
      goal: "Production-ready VerMAS with verification"
    - year: 2027
      goal: "Self-improving agent teams"
    - year: 2028
      goal: "Cross-organization agent collaboration"

values:
  - name: quality
    description: "Ship verified, tested work"
    enforcement:
      - "All merges require verification"
      - "Tests must pass before PR"
    metrics:
      - name: defect_rate
        threshold: "< 5%"

  - name: transparency
    description: "All decisions logged, all work visible"
    enforcement:
      - "All events logged to events.jsonl"
      - "No secret decision making"
    metrics:
      - name: audit_coverage
        threshold: "100%"

  - name: autonomy
    description: "Agents act independently within bounds"
    enforcement:
      - "Assignment Principle: execute without asking"
      - "Clear authority matrix"
    metrics:
      - name: escalation_rate
        threshold: "< 20%"  # Most work handled autonomously

  - name: human_oversight
    description: "Humans can always intervene"
    enforcement:
      - "Human can override any decision"
      - "Escalation paths always lead to human"
    metrics:
      - name: human_response_time
        threshold: "< 4h for escalations"

invariants:
  - id: INV-001
    name: "Known Actor"
    statement: "Every action has a known actor"
    enforcement: "AGENT_ID required on all events"
    violation_action: reject

  - id: INV-002
    name: "Owned Work"
    statement: "Every work item has an owner"
    enforcement: "assignee field required on work orders"
    violation_action: reject

  - id: INV-003
    name: "Immutable Compliance"
    statement: "Compliance rules cannot be disabled by agents"
    enforcement: "Compliance config is human-only editable"
    violation_action: escalate_to_human

  - id: INV-004
    name: "Human Override"
    statement: "Human can override any agent decision"
    enforcement: "L0 authority always available"
    violation_action: n/a  # This is always true

  - id: INV-005
    name: "Append-Only Audit"
    statement: "Audit log is append-only"
    enforcement: "events.jsonl never truncated"
    violation_action: alert_and_restore

amendments:
  - id: AMD-001
    date: 2026-01-01
    description: "Initial constitution"
    approved_by: human
```

### Enforcing Invariants

```python
# vermas/governance/invariants.py

from enum import Enum
from typing import Callable

class ViolationAction(Enum):
    REJECT = "reject"           # Block the action
    ESCALATE = "escalate"       # Allow but escalate to human
    ALERT = "alert"             # Allow but log alert
    RESTORE = "restore"         # Attempt to fix automatically

@dataclass
class Invariant:
    id: str
    name: str
    check: Callable[..., bool]  # Returns True if invariant holds
    violation_action: ViolationAction

INVARIANTS = [
    Invariant(
        id="INV-001",
        name="Known Actor",
        check=lambda event: event.actor is not None and event.actor != "",
        violation_action=ViolationAction.REJECT,
    ),
    Invariant(
        id="INV-002",
        name="Owned Work",
        check=lambda wo: wo.assignee is not None or wo.status == "open",
        violation_action=ViolationAction.REJECT,
    ),
    Invariant(
        id="INV-003",
        name="Immutable Compliance",
        check=lambda event: not (
            event.event_type == "config.changed"
            and "compliance" in event.data.get("path", "")
            and event.actor != "human"
        ),
        violation_action=ViolationAction.ESCALATE,
    ),
    Invariant(
        id="INV-005",
        name="Append-Only Audit",
        check=lambda: audit_log_line_count() >= last_known_count(),
        violation_action=ViolationAction.RESTORE,
    ),
]

def check_invariants(context: dict) -> List[InvariantViolation]:
    """Check all invariants and return violations."""
    violations = []
    for inv in INVARIANTS:
        try:
            if not inv.check(context):
                violations.append(InvariantViolation(
                    invariant_id=inv.id,
                    action=inv.violation_action,
                    context=context,
                ))
        except Exception as e:
            # Invariant check failed - treat as violation
            violations.append(InvariantViolation(
                invariant_id=inv.id,
                action=ViolationAction.ALERT,
                error=str(e),
            ))
    return violations
```

### Alignment Checks

```yaml
# .work/governance/alignment-checks.yaml

# These checks run periodically to ensure system stays aligned

checks:
  - id: ALN-001
    name: "Work aligns to objectives"
    schedule: daily
    check: |
      # All in-progress work should link to an objective
      for wo in work_orders.in_progress():
        if not wo.contributes_to:
          yield AlignmentWarning(
            work_order=wo.id,
            message="Work not linked to any objective"
          )

  - id: ALN-002
    name: "Resource allocation matches plan"
    schedule: weekly
    check: |
      # Actual time spent should match planned allocation
      for obj in objectives.active():
        actual = time_spent_on(obj)
        planned = capacity_allocated_to(obj)
        if abs(actual - planned) > planned * 0.2:
          yield AlignmentWarning(
            objective=obj.id,
            message=f"Allocation drift: {actual}h vs {planned}h planned"
          )

  - id: ALN-003
    name: "Values reflected in decisions"
    schedule: weekly
    check: |
      # Sample recent decisions and check for value alignment
      decisions = recent_decisions(days=7)
      for decision in sample(decisions, 10):
        if not references_value(decision):
          yield AlignmentWarning(
            decision=decision.id,
            message="Decision doesn't reference company values"
          )

  - id: ALN-004
    name: "Escalation paths functioning"
    schedule: daily
    check: |
      # All escalations should be resolved within SLA
      for escalation in escalations.unresolved():
        if escalation.age > escalation.sla:
          yield AlignmentCritical(
            escalation=escalation.id,
            message=f"Escalation overdue: {escalation.age}"
          )
```

### Self-Correction Mechanisms

```
SELF-CORRECTION HIERARCHY

┌─────────────────────────────────────────────────────────────────────────────┐
│ LEVEL 1: AUTOMATIC CORRECTION                                                │
│                                                                             │
│ When: Invariant violation detected with known fix                            │
│ Action: Apply fix automatically, log event                                   │
│                                                                             │
│ Examples:                                                                    │
│ - Missing actor → Reject action                                              │
│ - Orphaned work order → Assign to supervisor                                 │
│ - Truncated audit log → Restore from backup                                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ LEVEL 2: SUPERVISED CORRECTION                                               │
│                                                                             │
│ When: Violation requires judgment                                            │
│ Action: Supervisor reviews and approves fix                                  │
│                                                                             │
│ Examples:                                                                    │
│ - Alignment drift detected → Supervisor adjusts priorities                   │
│ - Value conflict in decision → Supervisor arbitrates                         │
│ - Resource contention → Supervisor reallocates                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ LEVEL 3: HUMAN CORRECTION                                                    │
│                                                                             │
│ When: Violation affects invariants or mission                                │
│ Action: Human reviews and decides                                            │
│                                                                             │
│ Examples:                                                                    │
│ - Compliance rule needs exception                                            │
│ - Mission interpretation unclear                                             │
│ - System-wide failure                                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Preventing Drift

```python
# vermas/governance/drift_detection.py

class DriftDetector:
    """Detect when system drifts from intended state."""

    def check_objective_drift(self) -> List[DriftWarning]:
        """Check if work is drifting from objectives."""
        warnings = []

        for objective in self.active_objectives():
            # Calculate alignment score
            work_items = self.work_contributing_to(objective)
            if not work_items:
                warnings.append(DriftWarning(
                    type="no_progress",
                    objective=objective.id,
                    message="No work items contributing to objective"
                ))
                continue

            # Check velocity toward objective
            progress_rate = self.progress_rate(objective)
            required_rate = self.required_rate_to_complete(objective)

            if progress_rate < required_rate * 0.8:
                warnings.append(DriftWarning(
                    type="velocity_drift",
                    objective=objective.id,
                    message=f"Progress rate {progress_rate} < required {required_rate}"
                ))

        return warnings

    def check_value_drift(self) -> List[DriftWarning]:
        """Check if decisions align with values."""
        warnings = []

        # Sample recent decisions
        decisions = self.recent_decisions(days=7)

        # Check each value has been referenced
        for value in self.company_values():
            references = [d for d in decisions if value.name in d.justification]
            if len(references) < len(decisions) * 0.1:
                warnings.append(DriftWarning(
                    type="value_underrepresented",
                    value=value.name,
                    message=f"Value '{value.name}' rarely referenced in decisions"
                ))

        return warnings
```

---

## Questions for Iteration 7

1. **How do we synthesize all governance into a coherent model?**
   - Single source of truth
   - Clear hierarchy
   - Enforcement mechanisms

2. **What's the governance file structure?**
   - Where does each config live?
   - What references what?

3. **What's the runtime governance?**
   - Hooks and validators
   - Event handlers
   - Periodic checks

---

## Iteration 6 Key Insights

1. **Constitution is the root**: Mission, vision, values, invariants in one file

2. **Invariants are non-negotiable**: Checked on every action, violations handled

3. **Alignment is monitored**: Periodic checks detect drift from objectives/values

4. **Self-correction is hierarchical**: Automatic → Supervised → Human

5. **Only humans can change the constitution**: Ultimate authority preserved

---

## Iteration 7: Governance Synthesis

### The Complete Governance Model

```
GOVERNANCE ARCHITECTURE (Unified View)

┌─────────────────────────────────────────────────────────────────────────────┐
│                           CONSTITUTION LAYER                                 │
│                      (Human-owned, rarely changes)                          │
├─────────────────────────────────────────────────────────────────────────────┤
│  Mission │ Vision │ Values │ Invariants                                    │
│                                                                             │
│  "Why we exist and what must always be true"                                │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                    ┌───────────────┼───────────────┐
                    ▼               ▼               ▼
┌───────────────────────┐ ┌───────────────────────┐ ┌───────────────────────┐
│    COMPLIANCE LAYER   │ │    PLANNING LAYER     │ │    AUTHORITY LAYER    │
│   (Rules & Constraints)│ │  (Goals & Timelines)  │ │ (Decisions & Escalation)│
├───────────────────────┤ ├───────────────────────┤ ├───────────────────────┤
│ • Compliance rules    │ │ • Annual plan         │ │ • Decision matrix     │
│ • Enforcement tiers   │ │ • Quarterly OKRs      │ │ • Escalation rules    │
│ • Escape hatches      │ │ • Sprint backlogs     │ │ • Delegation records  │
│ • Audit requirements  │ │ • Capacity allocation │ │ • Conflict resolution │
└───────────────────────┘ └───────────────────────┘ └───────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          OPERATIONS LAYER                                    │
│                    (Routing, Assignment, Execution)                          │
├─────────────────────────────────────────────────────────────────────────────┤
│  Routing Algorithm │ Work Decomposition │ Capacity Management              │
│                                                                             │
│  "How work gets to workers and gets done"                                   │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          MONITORING LAYER                                    │
│                     (Alignment, Drift, Metrics)                              │
├─────────────────────────────────────────────────────────────────────────────┤
│  Invariant Checks │ Alignment Checks │ Drift Detection │ Progress Tracking │
│                                                                             │
│  "Is the system behaving as intended?"                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Governance File Structure

```
.work/
├── governance/
│   ├── constitution.yaml         # Mission, vision, values, invariants
│   ├── compliance/
│   │   ├── rules.yaml            # Compliance rule definitions
│   │   ├── escape-hatches.yaml   # Override conditions
│   │   └── waivers/              # Logged waivers
│   │       └── waiver-001.yaml
│   ├── authority/
│   │   ├── decisions.yaml        # Decision authority matrix
│   │   ├── escalation.yaml       # Escalation rules and paths
│   │   └── delegations/          # Active delegations
│   │       └── del-001.yaml
│   └── alignment/
│       ├── checks.yaml           # Alignment check definitions
│       └── drift-thresholds.yaml # When to alert on drift
│
├── planning/
│   ├── annual-2026.yaml          # Annual strategic plan
│   ├── q1-2026.yaml              # Quarterly objectives
│   ├── sprints/
│   │   ├── sprint-2026-w01.yaml
│   │   └── sprint-2026-w02.yaml
│   └── retrospectives/
│       └── retro-2026-w01.yaml
│
└── events.jsonl                   # Audit log (append-only)
```

### Configuration Hierarchy

```
CONFIGURATION LOADING ORDER

1. constitution.yaml (root - human only)
       │
       ├── Sets: mission, vision, values, invariants
       └── Referenced by: everything else

2. compliance/rules.yaml
       │
       ├── Imports: invariants from constitution
       └── Enforced by: hooks, validators

3. authority/decisions.yaml
       │
       ├── Imports: values from constitution
       └── References: compliance rules for blocked decisions

4. planning/*.yaml
       │
       ├── Aligned to: mission, vision
       └── Constrained by: compliance rules

5. Agent CLAUDE.md files
       │
       ├── Reference: constitution
       └── Role-specific interpretation of values
```

### Runtime Governance Components

```python
# vermas/governance/runtime.py

from typing import List, Callable
from dataclasses import dataclass

@dataclass
class GovernanceRuntime:
    """Runtime governance enforcement."""

    # Loaded from config
    constitution: Constitution
    compliance_rules: List[ComplianceRule]
    decision_authority: DecisionMatrix
    escalation_rules: List[EscalationRule]
    alignment_checks: List[AlignmentCheck]

    # Runtime hooks
    pre_action_hooks: List[Callable]  # Before any action
    post_action_hooks: List[Callable]  # After any action
    periodic_checks: List[ScheduledCheck]  # Cron-style checks
    event_handlers: Dict[str, Callable]  # React to specific events

    def validate_action(self, action: Action) -> ValidationResult:
        """Check if action is allowed before execution."""
        # 1. Check invariants
        for inv in self.constitution.invariants:
            if not inv.check(action):
                return ValidationResult(
                    allowed=False,
                    reason=f"Invariant violation: {inv.name}",
                    action=inv.violation_action
                )

        # 2. Check compliance rules
        for rule in self.applicable_rules(action):
            result = rule.check(action)
            if not result.passed:
                if rule.tier == "hard":
                    return ValidationResult(
                        allowed=False,
                        reason=f"Compliance violation: {rule.id}",
                        action="reject"
                    )
                elif rule.tier == "soft":
                    if not action.has_waiver(rule.id):
                        return ValidationResult(
                            allowed=False,
                            reason=f"Requires waiver: {rule.id}",
                            action="request_waiver"
                        )

        # 3. Check decision authority
        if action.is_decision:
            authority = self.decision_authority.get(action.decision_type)
            if action.actor.level < authority.required_level:
                return ValidationResult(
                    allowed=False,
                    reason=f"Insufficient authority for {action.decision_type}",
                    action="escalate"
                )

        return ValidationResult(allowed=True)

    def handle_event(self, event: Event):
        """Process event through governance handlers."""
        # Log event (append-only audit)
        self.audit_log.append(event)

        # Run event handlers
        for event_type, handler in self.event_handlers.items():
            if event.matches(event_type):
                handler(event)

        # Check for escalation triggers
        for rule in self.escalation_rules:
            if rule.matches(event):
                self.trigger_escalation(rule, event)

    def run_periodic_checks(self):
        """Run scheduled governance checks."""
        for check in self.periodic_checks:
            if check.is_due():
                results = check.run()
                for warning in results.warnings:
                    self.handle_alignment_warning(warning)
                for critical in results.critical:
                    self.handle_alignment_critical(critical)
```

### Governance Hooks

```python
# Hooks that enforce governance at runtime

# Hook: Pre-action invariant check
def pre_action_hook(action: Action) -> HookResult:
    """Run before any action."""
    runtime = get_governance_runtime()
    result = runtime.validate_action(action)

    if not result.allowed:
        return HookResult(
            block=True,
            message=result.reason,
            suggested_action=result.action
        )
    return HookResult(block=False)

# Hook: Post-action audit logging
def post_action_hook(action: Action, result: ActionResult):
    """Run after any action."""
    event = Event(
        event_type=f"action.{action.type}",
        actor=action.actor,
        data=action.to_dict(),
        result=result.to_dict(),
        timestamp=now()
    )
    get_governance_runtime().handle_event(event)

# Hook: Work order compliance check
def work_order_compliance_hook(wo: WorkOrder) -> HookResult:
    """Check work order against compliance rules."""
    runtime = get_governance_runtime()

    for rule in runtime.compliance_rules:
        if rule.applies_to_work_orders:
            result = rule.check_work_order(wo)
            if not result.passed:
                return HookResult(
                    block=(rule.tier == "hard"),
                    message=f"Compliance: {rule.name} - {result.message}"
                )

    return HookResult(block=False)
```

### Periodic Governance Checks

```yaml
# .work/governance/periodic-checks.yaml

checks:
  # Invariant monitoring
  - id: invariant-sweep
    schedule: "*/5 * * * *"  # Every 5 minutes
    check: check_all_invariants
    on_violation: escalate_immediately

  # Alignment checks
  - id: objective-alignment
    schedule: "0 9 * * *"  # Daily at 9 AM
    check: check_objective_drift
    on_warning: notify_supervisor
    on_critical: notify_ceo

  - id: value-alignment
    schedule: "0 10 * * 1"  # Weekly on Monday
    check: check_value_drift
    on_warning: log_and_report
    on_critical: notify_human

  # Health checks
  - id: escalation-health
    schedule: "0 * * * *"  # Hourly
    check: check_escalation_sla
    on_violation: auto_escalate

  - id: audit-integrity
    schedule: "*/15 * * * *"  # Every 15 minutes
    check: verify_audit_log_integrity
    on_violation: alert_and_restore

  # Capacity monitoring
  - id: capacity-check
    schedule: "0 */4 * * *"  # Every 4 hours
    check: check_worker_capacity
    on_imbalance: suggest_rebalance
```

### Event-Driven Governance

```python
# Event handlers for governance

EVENT_HANDLERS = {
    # Compliance events
    "compliance.violation": lambda e: log_violation_and_maybe_block(e),
    "compliance.waived": lambda e: schedule_waiver_review(e),

    # Decision events
    "decision.escalated": lambda e: notify_escalation_target(e),
    "decision.overridden": lambda e: log_override_audit(e),

    # Planning events
    "sprint.scope_changed": lambda e: check_scope_creep_threshold(e),
    "objective.at_risk": lambda e: trigger_ceo_review(e),

    # Work events
    "work_order.stuck": lambda e: trigger_escalation_if_overdue(e),
    "work_order.blocked": lambda e: notify_blocker_owner(e),

    # Agent events
    "agent.offline": lambda e: reassign_work(e),
    "agent.failure": lambda e: increment_failure_count(e),

    # Invariant events
    "invariant.violated": lambda e: execute_violation_action(e),
    "invariant.restored": lambda e: log_restoration(e),
}
```

### Governance Dashboard

```
GOVERNANCE HEALTH DASHBOARD

Last updated: 2026-01-07 00:30:00 UTC

┌─────────────────────────────────────────────────────────────────────────────┐
│ INVARIANTS                                                          ✓ ALL OK │
├─────────────────────────────────────────────────────────────────────────────┤
│ INV-001 Known Actor        ✓  Last checked: 2 min ago                       │
│ INV-002 Owned Work         ✓  Last checked: 2 min ago                       │
│ INV-003 Immutable Compliance ✓  Last checked: 2 min ago                     │
│ INV-004 Human Override     ✓  (always true)                                 │
│ INV-005 Append-Only Audit  ✓  Last checked: 2 min ago                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ COMPLIANCE                                                      ⚠️ 1 WAIVER  │
├─────────────────────────────────────────────────────────────────────────────┤
│ Hard blocks:     0 active                                                   │
│ Soft blocks:     0 active                                                   │
│ Active waivers:  1 (DOC-001, expires 2026-01-08)                           │
│ Advisory alerts: 3 (see details)                                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ ESCALATIONS                                                     ✓ ALL CLEAR │
├─────────────────────────────────────────────────────────────────────────────┤
│ Active escalations:  0                                                      │
│ Overdue:             0                                                      │
│ Resolved (24h):      2                                                      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ ALIGNMENT                                                       ⚠️ 1 WARNING │
├─────────────────────────────────────────────────────────────────────────────┤
│ Objective drift:   1 warning (OBJ-002 behind schedule)                      │
│ Value drift:       None detected                                            │
│ Resource drift:    Within 10% of plan                                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ AUDIT LOG                                                           ✓ HEALTHY│
├─────────────────────────────────────────────────────────────────────────────┤
│ Events (24h):    1,247                                                      │
│ Size:            4.2 MB                                                     │
│ Integrity:       Verified                                                   │
│ Last backup:     2026-01-06 23:00 UTC                                       │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Governance Bootstrap

```python
# vermas/governance/bootstrap.py

def bootstrap_governance(work_dir: Path) -> GovernanceRuntime:
    """Initialize governance system from config files."""

    # 1. Load constitution (required)
    constitution_path = work_dir / "governance" / "constitution.yaml"
    if not constitution_path.exists():
        raise GovernanceError("Constitution not found - cannot start")
    constitution = Constitution.from_yaml(constitution_path)

    # 2. Load compliance rules
    compliance_rules = load_compliance_rules(
        work_dir / "governance" / "compliance"
    )

    # 3. Load decision authority
    decision_authority = DecisionMatrix.from_yaml(
        work_dir / "governance" / "authority" / "decisions.yaml"
    )

    # 4. Load escalation rules
    escalation_rules = load_escalation_rules(
        work_dir / "governance" / "authority" / "escalation.yaml"
    )

    # 5. Load alignment checks
    alignment_checks = load_alignment_checks(
        work_dir / "governance" / "alignment"
    )

    # 6. Load periodic checks schedule
    periodic_checks = load_periodic_checks(
        work_dir / "governance" / "periodic-checks.yaml"
    )

    # 7. Initialize audit log
    audit_log = AuditLog(work_dir / "events.jsonl")

    # 8. Wire up runtime
    runtime = GovernanceRuntime(
        constitution=constitution,
        compliance_rules=compliance_rules,
        decision_authority=decision_authority,
        escalation_rules=escalation_rules,
        alignment_checks=alignment_checks,
        periodic_checks=periodic_checks,
        audit_log=audit_log,
        pre_action_hooks=[pre_action_hook],
        post_action_hooks=[post_action_hook],
        event_handlers=EVENT_HANDLERS,
    )

    # 9. Run initial invariant check
    violations = runtime.check_all_invariants()
    if violations:
        for v in violations:
            logger.critical(f"Invariant violation at startup: {v}")
        raise GovernanceError("Cannot start with invariant violations")

    # 10. Start periodic check scheduler
    runtime.start_scheduler()

    logger.info("Governance runtime initialized successfully")
    return runtime
```

---

## Governance Summary

### The Complete Model

| Layer | Purpose | Owner | Changes |
|-------|---------|-------|---------|
| **Constitution** | Mission, vision, values, invariants | Human | Rarely (amendments) |
| **Compliance** | Rules and constraints | Human + CEO | Occasionally |
| **Authority** | Decision matrix, escalation | CEO | As needed |
| **Planning** | Objectives, sprints, capacity | CEO + Supervisor | Regular cycles |
| **Operations** | Routing, assignment, execution | Supervisor + Workers | Continuous |
| **Monitoring** | Alignment, drift, health | Automated | Always running |

### Key Enforcement Points

1. **Pre-action**: Invariant checks, compliance validation, authority verification
2. **Post-action**: Audit logging, event handling, escalation triggering
3. **Periodic**: Alignment checks, drift detection, health monitoring
4. **On-demand**: Waiver requests, delegation grants, plan changes

### Governance Guarantees

1. **Invariants are always checked** - No action proceeds without validation
2. **Everything is logged** - Append-only audit trail for all events
3. **Escalation paths exist** - Every situation has a resolution path
4. **Human authority preserved** - L0 can always override
5. **Self-correction operates** - Automatic → Supervised → Human hierarchy
6. **Alignment is monitored** - Drift from objectives/values is detected

---

## Iteration 7 Key Insights

1. **Governance is layered**: Constitution → Compliance → Authority → Planning → Operations → Monitoring

2. **File structure mirrors layers**: Clear separation of concerns in `.work/governance/`

3. **Runtime is event-driven**: Hooks, handlers, and periodic checks enforce governance

4. **Bootstrap validates invariants**: System refuses to start with violations

5. **Dashboard provides visibility**: Single view of governance health

---

## Approval Status

| Section | Status |
|---------|--------|
| Iteration 1: Governance Mapping | Pending Review |
| Iteration 2: Compliance & Constraints | Pending Review |
| Iteration 3: Work Routing | Pending Review |
| Iteration 4: Decision Making | Pending Review |
| Iteration 5: Planning Cycles | Pending Review |
| Iteration 6: Mission/Vision | Pending Review |
| Iteration 7: Governance Synthesis | Pending Review |
