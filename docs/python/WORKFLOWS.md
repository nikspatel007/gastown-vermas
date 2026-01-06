# VerMAS Workflows

> Process state machine and workflow patterns

## Process System

Processes are workflow instances that track multi-step work execution.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          PROCESS STATE MACHINE                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                         ┌─────────────┐                                     │
│                         │  TEMPLATE   │                                     │
│                         │             │                                     │
│                         │ .toml file  │                                     │
│                         │ Definition  │                                     │
│                         └──────┬──────┘                                     │
│                                │                                            │
│                                │ compile                                    │
│                                ▼                                            │
│                         ┌─────────────┐                                     │
│                         │   READY     │                                     │
│                         │             │                                     │
│                         │ Compiled    │                                     │
│                         │ In memory   │                                     │
│                         └──────┬──────┘                                     │
│                                │                                            │
│                                │ start                                      │
│                                ▼                                            │
│                         ┌─────────────┐                                     │
│                         │   ACTIVE    │                                     │
│                         │             │                                     │
│                         │ Persistent  │                                     │
│                         │ Tracked     │                                     │
│                         └──────┬──────┘                                     │
│                                │                                            │
│                 ┌──────────────┴──────────────┐                             │
│                 │ complete                    │ cancel                      │
│                 ▼                             ▼                              │
│          ┌─────────────┐              ┌─────────────┐                       │
│          │  ARCHIVE    │              │  DISCARDED  │                       │
│          │             │              │             │                       │
│          │ Record kept │              │ No trace   │                       │
│          └─────────────┘              └─────────────┘                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## State Definitions

| State | Name | Storage | Persistence | Use Case |
|-------|------|---------|-------------|----------|
| **Template** | Definition | `.work/templates/*.toml` | Permanent | Workflow definitions |
| **Ready** | Compiled | In memory | Session | Compiled, ready to use |
| **Active** | Process | `.work/processes/*.json` | Persistent | Active work tracking |
| **Archive** | Record | `.work/processes/*.archive.json` | Permanent | Completed workflows |

---

## Operators

| Operator | From | To | Description |
|----------|------|----|-------------|
| **compile** | Template | Ready | Compile template into executable form |
| **start** | Ready | Active | Create persistent workflow attached to work order |
| **complete** | Active | Archive | Complete with summary, create record |
| **cancel** | Active | (gone) | Discard without record |

---

## Template Structure (TOML)

Templates are workflow definitions in TOML format:

```
.work/templates/
├── supervisor-patrol.toml
├── qa-merge.toml
├── worker-execute.toml
└── verify-pipeline.toml
```

### Template Fields

| Field | Type | Description |
|-------|------|-------------|
| `template` | string | Unique identifier |
| `description` | string | What this workflow does |
| `version` | int | Schema version |
| `steps` | array | Ordered list of steps |

### Step Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique step identifier |
| `title` | string | Human-readable name |
| `description` | string | Detailed instructions (markdown) |
| `needs` | array | Step IDs that must complete first |

---

## Example Templates

### Supervisor Patrol

```
template = "supervisor-patrol"
description = "Monitor workers, nudge idle, escalate stuck"
version = 1

STEPS:
1. inbox-check
   - Check mail for WORKER_DONE messages
   - No dependencies

2. survey-workers
   - List all active workers
   - Check idle time for each
   - Needs: inbox-check

3. nudge-idle
   - Send nudge to workers idle >5min
   - Needs: survey-workers

4. escalate-stuck
   - Kill workers stuck >15min
   - Report to Operations
   - Needs: nudge-idle

FLOW:
inbox-check → survey-workers → nudge-idle → escalate-stuck
```

### Worker Execute

```
template = "worker-execute"
description = "Execute assigned work order to completion"
version = 1

STEPS:
1. understand
   - Read assignment and work order details
   - Plan approach
   - No dependencies

2. implement
   - Write code
   - Follow conventions
   - Needs: understand

3. test
   - Run tests
   - Verify requirements
   - Needs: implement

4. complete
   - Commit and push
   - Signal done
   - Needs: test

FLOW:
understand → implement → test → complete
```

### QA Verification

```
template = "verify-pipeline"
description = "VerMAS verification workflow"
version = 1

STEPS:
1. elaborate
   - Designer creates spec
   - No dependencies

2. strategize
   - Strategist creates tests
   - Needs: elaborate

3. verify
   - Verifier runs tests (no LLM)
   - Needs: strategize

4. audit
   - Auditor reviews evidence
   - Needs: verify

5. advocate
   - Argue for PASS
   - Needs: audit

6. criticize
   - Argue for FAIL
   - Needs: audit

7. judge
   - Render verdict
   - Needs: advocate, criticize

FLOW:
                    ┌─→ advocate ─┐
elaborate → strategize → verify → audit ─┤              ├─→ judge
                    └─→ criticize─┘
```

---

## Step Execution Model

### Step Status

| Status | Meaning |
|--------|---------|
| `pending` | Not started, may have unmet dependencies |
| `blocked` | Dependencies not yet complete |
| `ready` | Dependencies met, can execute |
| `in_progress` | Currently executing |
| `completed` | Finished successfully |
| `failed` | Execution failed |
| `skipped` | Explicitly skipped |

### Execution Rules

1. **Ready check:** Step is ready when all `needs` are `completed`
2. **One at a time:** Only one step `in_progress` at a time (per process)
3. **No backtrack:** Completed steps don't re-run
4. **Fail fast:** Failed step blocks downstream steps

### Finding Ready Steps

```
ready_steps = []
for each step:
    if step.status == pending:
        if all(dep.status == completed for dep in step.needs):
            ready_steps.append(step)
```

---

## Process Lifecycle

### Creating a Process

```
1. Load template from .toml file
2. Compile to ready state (compile)
3. Attach to work order (start)
4. Save to .work/processes/{id}.json
```

### Executing a Process

```
1. Load process from .json
2. Find ready steps
3. For each ready step:
   a. Mark in_progress
   b. Execute step instructions
   c. Mark completed (or failed)
4. Repeat until all done or stuck
```

### Completing a Process

```
Option A - Complete (with record):
1. Generate summary
2. Write to .archive.json
3. Delete active .json

Option B - Cancel (no record):
1. Delete active .json
2. No archive created
```

---

## Parallel Execution

Some workflows support parallel steps:

```
QA Verification:
                    ┌─→ advocate ─┐
... → audit ───────┤              ├─→ judge
                    └─→ criticize─┘

Advocate and Critic run in parallel after Audit completes.
Judge waits for both.
```

### Implementing Parallelism

1. Multiple steps can be `ready` simultaneously
2. Launch each ready step in its own context
3. Parent waits for all children
4. Next step starts when all dependencies complete

---

## Workflow Patterns

### Sequential Pipeline

```
A → B → C → D

Each step depends on the previous.
Simple, predictable, easy to debug.
```

### Fork-Join

```
    ┌─→ B ─┐
A ──┤      ├──→ D
    └─→ C ─┘

A completes, B and C run in parallel, D waits for both.
Used for: parallel testing, adversarial review.
```

### Diamond

```
    ┌─→ B ─┐
A ──┤      ├──→ D
    └─→ C ─┘
        │
        ▼
        E

D needs B and C.
E only needs C.
```

### Conditional (not yet implemented)

```
A → [condition] → B (if true)
              └─→ C (if false)

Would require step predicates.
Future enhancement.
```

---

## Integration with Agents

### Supervisor Patrol

```
Supervisor starts
    │
    ▼
start(supervisor-patrol)
    │
    ├─→ Execute steps in loop
    │
    └─→ When complete, start again (infinite patrol)
```

### Worker Execute

```
Worker spawns with work order
    │
    ▼
start(worker-execute, wo_id)
    │
    ├─→ Execute steps
    │
    ├─→ On complete: archive + signal done
    │
    └─→ On fail: leave for Supervisor
```

### QA Merge

```
QA receives READY_FOR_QA
    │
    ▼
start(verify-pipeline, wo_id)  [if VerMAS enabled]
    │
    ├─→ Execute verification
    │
    ├─→ On PASS: merge + archive
    │
    └─→ On FAIL: REWORK_REQUEST + cancel
```

---

## Debugging Workflows

### Inspecting Active Processes

```bash
wo process list              # List all active processes
wo process show {id}         # Show process details
wo process steps {id}        # Show step status
```

### Common Issues

| Symptom | Cause | Fix |
|---------|-------|-----|
| Step stuck in `pending` | Dependency not complete | Check dependency status |
| Step stuck in `in_progress` | Execution hanging | Check agent session |
| No ready steps | All blocked | Check for circular deps |
| Process not completing | Failed step | Review step output |

---

## Events Emitted

Workflow operations emit events to the event log. See [EVENTS.md](./EVENTS.md).

| Event Type | When | Data |
|------------|------|------|
| `process.created` | Process started | template, wo_id, process_id |
| `process.step_started` | Step begins | process_id, step_id |
| `process.step_completed` | Step finishes | process_id, step_id, status |
| `process.completed` | All steps done | process_id, summary |
| `process.cancelled` | Workflow discarded | process_id, reason |

### Tracking Workflow Progress

```python
# Get all events for a process
process_events = get_events(filter={"correlation_id": process_id})

# Compute step durations from events
for step_start in get_events(type="process.step_started"):
    step_end = get_event(type="process.step_completed",
                         filter={"step_id": step_start.data["step_id"]})
    duration = step_end.timestamp - step_start.timestamp
```

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [AGENTS.md](./AGENTS.md) - Agent roles
- [HOOKS.md](./HOOKS.md) - Claude Code integration and git worktrees
- [MESSAGING.md](./MESSAGING.md) - Communication patterns
- [EVENTS.md](./EVENTS.md) - Event sourcing and change feeds
- [SCHEMAS.md](./SCHEMAS.md) - Template and process data specs
- [VERIFICATION.md](./VERIFICATION.md) - VerMAS QA workflow
- [EVALUATION.md](./EVALUATION.md) - How to evaluate the system
