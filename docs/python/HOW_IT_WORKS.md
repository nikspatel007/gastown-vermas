# VerMAS - How It Works

> Multi-agent coordination built entirely on proven technologies

**See also:** [INDEX.md](./INDEX.md) for full documentation map

## Overview

VerMAS is a **multi-agent workspace manager** that coordinates multiple Claude Code agents working in parallel. It's designed to scale comfortably to 20-30 agents through structured coordination.

Think of it as a **company with factories**: the CEO coordinates, supervisors monitor, workers execute, and QA ensures quality.

## Proven Technology Foundation

Every component uses battle-tested technology:

| Need | Technology | Why Proven |
|------|------------|------------|
| State storage | **File system** (JSONL) | 50+ years, universal, debuggable |
| Version control | **Git** (worktrees) | Industry standard, distributed |
| Agent runtime | **Claude Code CLI** | No API costs, profiles, hooks |
| Process isolation | **Tmux** sessions | Decades old, rock solid |
| CLI framework | **Typer** (Python) | Modern, type-safe |
| Data flow | **Event sourcing** | Banking, Kafka, Redux |

**No databases. No message brokers. No containers. Just files, git, and processes.**

## The Company Metaphor

- **Company** = Your workspace, the organization
- **CEO** = Strategic coordinator (you/Claude), doesn't write code
- **Factory** = A project/repository
- **Supervisor** = Monitors workers in a factory
- **QA** = Quality control, verification, merge queue
- **Worker** = Ephemeral task executor (spawn → work → done)
- **Team** = Human-directed workspace

When workers have assignments, they **execute immediately**. No confirmation. No questions.

---

## System at a Glance

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        VERMAS: COMPANY STRUCTURE                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   HEADQUARTERS                   FACTORIES                   STORAGE        │
│   ────────────                   ─────────                   ───────        │
│   CEO (coordinator)              Workers (Claude CLI)        JSONL files    │
│   Operations (daemon)            Codex/Aider (optional)      Git repos      │
│   Board (human)                  Tmux isolation              Worktrees      │
│         │                              │                           │        │
│         └──────────────────────────────┼───────────────────────────┘        │
│                                        │                                    │
│                                        ▼                                    │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                         EVENT SOURCING                               │  │
│   │                                                                      │  │
│   │   events.jsonl → work_orders.jsonl (tasks)                          │  │
│   │                → messages.jsonl (comms)                              │  │
│   │                → feed.jsonl (real-time)                              │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                        │                                    │
│                 ┌──────────────────────┼──────────────────────┐            │
│                 │                      │                      │            │
│                 ▼                      ▼                      ▼            │
│   ┌─────────────────────┐ ┌─────────────────────┐ ┌─────────────────────┐ │
│   │    ASSIGNMENTS      │ │   COMMUNICATIONS    │ │     PROCESSES       │ │
│   │                     │ │                     │ │                     │ │
│   │ .assignment-{agent} │ │ co inbox/send       │ │ Templates → Active  │ │
│   │ Execute immediately │ │ Async messaging     │ │ Workflow tracking   │ │
│   └─────────────────────┘ └─────────────────────┘ └─────────────────────┘ │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Architecture

```
Company (workspace root)
├── ceo/                ← Strategic coordinator (YOU/Claude)
├── .work/              ← Company-level work orders
├── <factory>/          ← Project container
│   ├── .work/          ← Factory-level work orders
│   ├── ceo/factory/    ← Read-only reference clone
│   ├── qa/             ← Quality assurance, merge queue
│   ├── supervisor/     ← Worker lifecycle manager
│   ├── teams/          ← Human-directed workspaces
│   └── workers/        ← Ephemeral worker worktrees
```

---

## Roles

### Headquarters

| Role | Purpose |
|------|---------|
| **CEO** | Cross-factory coordinator. Dispatches work, handles escalations, makes strategic decisions. **Does NOT write code.** |
| **Operations** | Daemon process managing agent lifecycle and health |
| **Board** | Human role. Sets strategy, reviews outputs, handles escalations |

### Per-Factory

| Role | Purpose |
|------|---------|
| **Supervisor** | Monitors workers, detects stuck processes, handles lifecycle events |
| **QA Department** | Manages merge queues, verification, code review |
| **Worker** | Ephemeral executor (spawn → work → disappear) |
| **Team** | Human-directed workspaces for hands-on work |

---

## Work Orders

Git-backed work tracking where all state lives. Commands use `wo` prefix:
- `wo create` - Create a work order
- `wo ready` - Show work orders ready to execute (no blockers)
- `wo show <id>` - View work order details
- `wo close <id>` - Mark complete
- `wo sync` - Sync with git
- `wo dep add <child> <parent>` - Add dependency

### JSONL Format

Work orders stored in `.work/work_orders.jsonl` as one JSON object per line:

```json
{"id":"wo-abc12","title":"Feature X","description":"...","status":"open","priority":2,"type":"task","created_at":"2025-12-28T13:10:32Z","created_by":"factory-a/teams/frontend"}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Hash-based ID: `wo-xxxxx` or hierarchical `wo-xxxxx.1` |
| `title` | string | Short title |
| `description` | string | Full description |
| `status` | enum | `open`, `closed`, `assigned`, `in_progress` |
| `priority` | int | 0=P0 (critical), 2=P2 (default), 4=P4 (backlog) |
| `type` | enum | `task`, `bug`, `feature`, `epic` |
| `created_by` | string | Agent address |

---

## The Assignment Principle

> **If you have an assignment, EXECUTE IT.**

1. Agent starts up
2. Checks assignment (`co assignment`)
3. If work is assigned → **EXECUTE IMMEDIATELY**
4. If empty → Check messages, then wait for instructions

**No confirmation. No questions. No waiting.**

---

## Processes (Workflows)

Workflow state machine:

| State | Name | Storage | Use Case |
|-------|------|---------|----------|
| **Template** | TOML file | `.work/templates/*.toml` | Workflow definitions |
| **Ready** | Compiled | In memory | Ready to execute |
| **Active** | Running | `.work/processes/*.json` | Work in progress |
| **Archive** | Completed | `.work/processes/*.archive.json` | Historical record |

### Operators

| Operator | Description |
|----------|-------------|
| **compile** | Template → Ready |
| **start** | Ready → Active (attached to work order) |
| **complete** | Active → Archive (with summary) |
| **cancel** | Active → (gone) |

---

## Event Sourcing

All state changes are immutable events:

```
events.jsonl (source of truth, append-only)
     │
     └─→ work_orders.jsonl  (current state - projection)
     └─→ messages.jsonl     (communications - projection)
     └─→ feed.jsonl         (real-time change feed)
```

**Key event types:**
- `work_order.created`, `work_order.status_changed`, `work_order.closed`
- `message.sent`, `message.delivered`, `message.read`
- `agent.started`, `agent.working`, `agent.stopped`
- `assignment.set`, `assignment.checked`, `assignment.cleared`
- `process.created`, `process.step_completed`, `process.completed`

See [EVENTS.md](./EVENTS.md) for full documentation.

---

## Quick Start Examples

### Example 1: Create and Assign a Bug Fix

```bash
# 1. CEO creates the bug
wo create --title="Fix login timeout" --type=bug --priority=1
# Output: Created work order wo-abc123

# 2. CEO assigns to a factory
co dispatch wo-abc123 myproject
# Output: Spawned worker-myproject-slot0, assigned wo-abc123

# 3. Worker (in its session) checks assignment and works
co assignment
# Output: ASSIGNED: wo-abc123

wo show wo-abc123
# Output: [bug details]

# ... worker fixes the bug ...

# 4. Worker signals completion
co worker done

# 5. Supervisor validates, sends to QA
# 6. QA runs tests, verifies, merges
```

### Example 2: Monitor Progress

```bash
# Company-wide status
co status

# Active workers
co workers

# Work in progress
co sprints

# Event stream
tail -f .work/feed.jsonl | jq .
```

### Example 3: Handle a Handoff

```bash
# Previous session left work incomplete
co assignment
# Output: ASSIGNED: wo-abc123

# Check for handoff context
co inbox
# Output: 1 message - "HANDOFF: Login fix incomplete"

co read msg-xyz
# Output: "Fixed timeout but tests still failing on CI..."

# Continue the work
wo show wo-abc123
# ... continue fixing ...
```

---

## Common Workflows

### Starting the Day
```bash
co ops start         # Start infrastructure
co assignment        # Check for assigned work
co inbox             # Check messages
co status            # Company overview
```

### Creating Work
```bash
wo create --title="Feature X" --type=feature --priority=2
wo dep add <wo> <depends-on>     # Add dependency
co dispatch <wo> <factory>        # Assign to worker
```

### Session End Checklist
```bash
git status              # Check changes
git add <files>         # Stage code
wo sync                 # Commit work order changes
git commit -m "..."     # Commit code
git push                # Push to remote
# If incomplete work:
co send ceo -s "HANDOFF: <brief>" -m "<context>"
```

---

## Why This Architecture?

Managing 4-10 agents creates chaos. VerMAS enables scaling to 20-30 agents through:
- Structured coordination via CEO
- Persistent work state (survives crashes)
- Clear role separation
- Automated handoffs
- Git-backed everything
- Event sourcing for full auditability

---

## Further Reading

| Document | What You'll Learn |
|----------|-------------------|
| [INDEX.md](./INDEX.md) | Documentation map and glossary |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | System design and data flow |
| [CLI.md](./CLI.md) | Complete command reference |
| [OPERATIONS.md](./OPERATIONS.md) | Deployment, startup, and maintenance |
| [HOOKS.md](./HOOKS.md) | Claude Code integration, git worktrees |
| [EVENTS.md](./EVENTS.md) | Event sourcing patterns |
| [AGENTS.md](./AGENTS.md) | Role responsibilities |
| [MESSAGING.md](./MESSAGING.md) | Internal communications |
| [WORKFLOWS.md](./WORKFLOWS.md) | Process system |
| [VERIFICATION.md](./VERIFICATION.md) | QA pipeline |
| [EVALUATION.md](./EVALUATION.md) | Metrics and testing |
