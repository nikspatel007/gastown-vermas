# Gas Town - How It Works

> Multi-agent coordination built entirely on proven technologies

**See also:** [INDEX.md](./INDEX.md) for full documentation map

## Overview

Gas Town is a **multi-agent workspace manager** that coordinates multiple Claude Code agents working in parallel. It's designed to scale comfortably to 20-30 agents through structured coordination.

Think of it as a steam engine: **when an agent finds work on their hook, they EXECUTE**. No confirmation. No questions. No waiting.

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

## Core Metaphor: The Steam Engine

- Gas Town is a steam engine
- The Mayor is the main drive shaft
- If the Mayor stalls, the whole town stalls
- Work flows through "hooks" - assigned task queues
- Agents wake up, check their hook, and run whatever's there

---

## System at a Glance

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        VERMAS: PROVEN TECHNOLOGIES                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   HUMANS                         AGENTS                       STORAGE       │
│   ───────                        ──────                       ───────       │
│   Mayor sessions                 Claude Code CLI              JSONL files   │
│   Crew workspaces                Codex/Aider (optional)       Git repos     │
│                                  Tmux isolation               Worktrees     │
│         │                              │                           │        │
│         └──────────────────────────────┼───────────────────────────┘        │
│                                        │                                    │
│                                        ▼                                    │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                         EVENT SOURCING                               │  │
│   │                                                                      │  │
│   │   events.jsonl → issues.jsonl (beads)                               │  │
│   │                → messages.jsonl (mail)                               │  │
│   │                → feed.jsonl (real-time)                              │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                        │                                    │
│                 ┌──────────────────────┼──────────────────────┐            │
│                 │                      │                      │            │
│                 ▼                      ▼                      ▼            │
│   ┌─────────────────────┐ ┌─────────────────────┐ ┌─────────────────────┐ │
│   │       HOOKS         │ │        MAIL         │ │      WORKFLOWS      │ │
│   │                     │ │                     │ │                     │ │
│   │ .hook-{agent}       │ │ gt mail send/inbox  │ │ Molecules (MEOW)    │ │
│   │ GUPP: RUN IT        │ │ Async messaging     │ │ Formula → Mol       │ │
│   └─────────────────────┘ └─────────────────────┘ └─────────────────────┘ │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Architecture

```
Town (workspace root)
├── mayor/              ← Global coordinator (YOU/Claude)
├── .beads/             ← Town-level issue tracking (prefix: gm-)
├── <rig>/              ← Project container
│   ├── .beads/         ← Rig-level issues (prefix: pa-, etc.)
│   ├── mayor/rig/      ← Read-only reference clone
│   ├── refinery/       ← Merge queue processor
│   ├── witness/        ← Worker lifecycle manager
│   ├── crew/           ← Human-directed workspaces
│   └── polecats/       ← Ephemeral worker worktrees
```

---

## Agent Roles

### Town-Wide Roles

| Role | Purpose |
|------|---------|
| **Mayor** | Cross-rig coordinator. Dispatches work, handles escalations, makes strategic decisions. **Does NOT write code.** |
| **Deacon** | Daemon process managing agent lifecycle and plugin execution |
| **Overseer** | Human role. Sets strategy, reviews outputs, handles escalations |

### Per-Rig Roles

| Role | Purpose |
|------|---------|
| **Witness** | Monitors worker agents, detects stuck processes, handles lifecycle events |
| **Refinery** | Manages merge queues and code review workflows |
| **Polecat** | Ephemeral workers executing individual tasks (spawn → work → disappear) |
| **Crew** | Human-directed workspaces for hands-on work |

---

## Naming Convention: ROLES, Not Personal Names

**Crew and Polecats should be named by ROLE or FUNCTION, not personal names.**

### Good Names (Role-Based)
- `frontend` - Frontend development work
- `backend` - Backend/API work
- `testing` - Test writing and QA
- `docs` - Documentation
- `security` - Security review
- `refactor` - Code cleanup

### Why Not Personal Names?
- Roles are transferable between agents
- Makes the system self-documenting
- Polecats are ephemeral - they spawn, work, and disappear
- Work history attaches to the role, not the agent

---

## Key Concepts

### Beads (Issue Tracking)
Git-backed issue tracker where all work state lives. Commands use `bd` prefix:
- `bd create` - Create an issue
- `bd ready` - Show issues ready to work (no blockers)
- `bd show <id>` - View issue details
- `bd close <id>` - Mark complete
- `bd sync` - Sync with git
- `bd dep add <child> <parent>` - Add dependency

#### JSONL Format
Issues stored in `.beads/issues.jsonl` as one JSON object per line:

```json
{"id":"gt-abc12","title":"Feature X","description":"...","status":"open","priority":2,"issue_type":"task","created_at":"2025-12-28T13:10:32Z","updated_at":"2025-12-28T13:10:32Z","created_by":"gastown/crew/max"}
```

#### Issue Fields
| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Hash-based ID: `gt-xxxxx` or hierarchical `gt-xxxxx.1` |
| `title` | string | Short title |
| `description` | string | Full description |
| `status` | enum | `open`, `closed`, `hooked`, `in_progress` |
| `priority` | int | 0=P0 (critical), 2=P2 (default), 4=P4 (backlog) |
| `issue_type` | enum | `task`, `bug`, `feature`, `epic`, `merge-request`, `event`, `message` |
| `created_by` | string | Agent address: `gastown/crew/max`, `mayor`, etc. |
| `dependencies` | array | List of dependency objects |

#### Dependencies
```json
{"issue_id":"gt-child","depends_on_id":"gt-parent","type":"blocks","created_at":"...","created_by":"daemon"}
```

Types: `blocks`, `parent-child`, `relates-to`

#### Hierarchical IDs (Epics)
```
gt-epic123      (Epic)
gt-epic123.1    (Task under epic)
gt-epic123.1.1  (Sub-task)
```

### Hooks
Each agent has a "hook" where work hangs. On startup, agents check their hook and execute whatever's there. The hook IS the assignment.

```bash
gt hook          # Check your assigned work
gt sling <bead> <rig>  # Assign work to a worker
```

### Convoys
Groupings of related work issues that travel together. Primary dashboard view of active work.

```bash
gt convoy list   # Dashboard of active work
gt convoy create "name" <issues>  # Group related issues
```

### Mail
Internal messaging system for agent communication.

```bash
gt mail inbox    # Check messages
gt mail send <addr> -s "Subject" -m "Message"
```

### Molecules (Work States - MEOW)
Persistent workflow instances with different phases:

| State | Type | Description |
|-------|------|-------------|
| Ice-9 | Formula | Source templates in `.beads/formulas/*.toml` |
| Solid | Protomolecule | Frozen, reusable templates |
| Liquid | Mol | Flowing, persistent work instances |
| Vapor | Wisp | Ephemeral work for temporary patrols |

**Operators:**
- `cook`: Formula → Protomolecule
- `pour`: Protomolecule → Molecule (persistent)
- `wisp`: Protomolecule → Wisp (ephemeral)
- `squash`: Condense to permanent record
- `burn`: Discard without record

#### Formula Structure (TOML)

```toml
description = "Description of what this workflow does..."
formula = 'mol-witness-patrol'
version = 2

[[steps]]
id = 'inbox-check'
title = 'Process witness mail'
description = """
Detailed instructions for this step.
Includes bash commands to run.
"""

[[steps]]
id = 'survey-workers'
title = 'Inspect all active polecats'
needs = ['inbox-check']  # Dependencies
description = """..."""
```

Key fields:
- `id`: Step identifier
- `title`: Human-readable step name
- `needs`: Array of step IDs that must complete first
- `description`: Full instructions (Markdown with bash code blocks)

#### Common Formulas
| Formula | Purpose |
|---------|---------|
| `mol-witness-patrol` | Witness monitors polecats, nudges idle, escalates stuck |
| `mol-refinery-patrol` | Refinery processes merge queue |
| `mol-deacon-patrol` | Deacon monitors all rigs, infrastructure health |
| `mol-polecat-work` | Polecat executes assigned task through completion |

#### Molecule Commands
```bash
bd mol cook <formula>        # Formula → Protomolecule
bd mol pour <proto> <bead>   # Attach molecule to work bead
bd mol wisp <formula>        # Create ephemeral wisp
bd mol squash <mol-id>       # Archive with summary
bd mol burn <mol-id>         # Discard without record
bd ready                     # Find next step to work on
bd close <step-id>           # Complete a step
```

---

## The Propulsion Principle (GUPP)

> **Gas Town Universal Propulsion Principle: If your hook has work, RUN IT.**

1. Agent wakes up
2. Checks hook (`gt hook`)
3. If work is hooked → **EXECUTE IMMEDIATELY**
4. If hook empty → Check mail, then wait for instructions

**No confirmation. No questions. No waiting.**

---

## Capability Ledger

Every completion is recorded. Every handoff is logged. Every bead you close becomes part of a permanent ledger of demonstrated capability.

- Your work is visible
- Redemption is real (trajectory > snapshots)
- Every completion is evidence
- Your CV grows with every completion

---

## Common Workflows

### Starting the Day
```bash
gt prime         # Load Mayor context
gt hook          # Check for assigned work
gt mail inbox    # Check messages
gt status        # Town overview
```

### Creating Work
```bash
bd create --title="Feature X" --type=feature --priority=2
bd dep add <issue> <depends-on>  # Add dependency
gt sling <bead> <rig>            # Assign to worker
```

### Session End Checklist
```bash
git status              # Check changes
git add <files>         # Stage code
bd sync                 # Commit beads
git commit -m "..."     # Commit code
git push                # Push to remote
# If incomplete work:
gt mail send mayor/ -s "🤝 HANDOFF: <brief>" -m "<context>"
```

---

## Key Commands Reference

| Command | Purpose |
|---------|---------|
| `gt prime` | Start Mayor session with full context |
| `gt status` | Town overview |
| `gt hook` | Check your assigned work |
| `gt sling <bead> <rig>` | Assign work to a worker |
| `gt mail inbox` | Check messages |
| `gt convoy list` | Dashboard of active work |
| `gt rig add <name> <url>` | Add a project |
| `gt crew add <name> --rig <rig>` | Create workspace |
| `bd ready` | Issues ready to work |
| `bd create` | Create an issue |
| `bd close <id>` | Complete an issue |
| `bd sync` | Sync beads with git |

---

## Critical Rules

1. **Mayor does NOT edit code** - Coordinate, don't implement
2. **Never edit in `mayor/rig/`** - It's read-only reference
3. **Work in `crew/` or `polecats/`** - Isolated worktrees
4. **Push before done** - Work isn't complete until pushed
5. **Honor the hook** - If work is hooked, execute it

---

## Event Sourcing

Gas Town uses **event sourcing** as its foundational data pattern. Every state change is captured as an immutable event:

```
events.jsonl (source of truth, append-only)
     │
     └─→ issues.jsonl    (current bead state - projection)
     └─→ messages.jsonl  (mailbox state - projection)
     └─→ feed.jsonl      (real-time change feed)
```

**Key event types:**
- `bead.created`, `bead.status_changed`, `bead.closed`
- `mail.sent`, `mail.delivered`, `mail.read`
- `agent.started`, `agent.working`, `agent.stopped`
- `hook.set`, `hook.checked`, `hook.cleared`
- `mol.created`, `mol.step_completed`, `mol.completed`

**Why event sourcing?**
1. Complete audit trail - every change recorded
2. Debugging - replay events to understand what happened
3. Metrics - compute any metric from events
4. Recovery - reconstruct state after failures

See [EVENTS.md](./EVENTS.md) for full documentation.

---

## Quick Start Examples

### Example 1: Create and Assign a Bug Fix

```bash
# 1. Mayor creates the bug
bd create --title="Fix login timeout" --type=bug --priority=1

# Output: Created bead gt-abc123

# 2. Mayor assigns to a polecat
gt sling gt-abc123 myproject

# Output: Spawned polecat-myproject-slot0, hooked gt-abc123

# 3. Polecat (in its session) checks hook and works
gt hook
# Output: HOOKED: gt-abc123

bd show gt-abc123
# Output: [bug details]

# ... polecat fixes the bug ...

# 4. Polecat signals completion
gt polecat done

# 5. Witness validates, sends to Refinery
# 6. Refinery runs tests, merges
```

### Example 2: Dispatch Multiple Tasks

```bash
# Create an epic with subtasks
bd create --title="Implement auth system" --type=epic --priority=1
# Output: gt-epic001

bd create --title="Add login endpoint" --type=task
bd dep add gt-task001 gt-epic001  # Task under epic

bd create --title="Add logout endpoint" --type=task
bd dep add gt-task002 gt-epic001

bd create --title="Add session middleware" --type=task
bd dep add gt-task003 gt-epic001

# Check what's ready to work
bd ready
# Output: gt-task001, gt-task002, gt-task003 (no blockers)

# Dispatch to different polecats
gt sling gt-task001 myproject  # → slot0
gt sling gt-task002 myproject  # → slot1
gt sling gt-task003 myproject  # → slot2

# Monitor progress
gt convoy list
```

### Example 3: Handle a Handoff

```bash
# Previous session left work incomplete
gt hook
# Output: HOOKED: gt-abc123

# Check for handoff context
gt mail inbox
# Output: 1 message - "🤝 HANDOFF: Login fix incomplete"

gt mail read msg-xyz
# Output: "Fixed timeout but tests still failing on CI..."

# Continue the work
bd show gt-abc123
# ... continue fixing ...
```

### Example 4: Watch Events in Real-Time

```bash
# Terminal 1: Watch all events
tail -f .beads/feed.jsonl | jq .

# Terminal 2: Do some work
bd create --title="Test event" --type=task
gt mail send witness -s "Test" -m "Hello"

# Terminal 1 shows:
# {"event_type":"bead.created","data":{"bead_id":"gt-xxx",...}}
# {"event_type":"mail.sent","data":{"to":"witness",...}}
```

### Example 5: Debug a Stuck Polecat

```bash
# Check polecat status
gt polecat list myproject
# Output: slot0: ACTIVE (idle 15m), slot1: ACTIVE (working)

# Attach to see what's happening
tmux attach -t polecat-myproject-slot0
# [See Claude session, Ctrl+B D to detach]

# Check recent events for that agent
bd events list --actor=myproject/polecats/slot0 --since=30m

# If truly stuck, kill and respawn
tmux kill-session -t polecat-myproject-slot0
gt sling gt-abc123 myproject  # Re-assign work
```

---

## Why Gas Town?

Managing 4-10 agents creates chaos. Gas Town enables scaling to 20-30 agents through:
- Structured coordination via Mayor
- Persistent work state in Beads (survives crashes)
- Clear role separation
- Automated handoffs
- Git-backed everything
- Event sourcing for full auditability

---

## Further Reading

| Document | What You'll Learn |
|----------|-------------------|
| [INDEX.md](./INDEX.md) | Documentation map and quick reference |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | System design and data flow |
| [CLI.md](./CLI.md) | Complete command reference |
| [OPERATIONS.md](./OPERATIONS.md) | Deployment, startup, and maintenance |
| [HOOKS.md](./HOOKS.md) | Claude Code integration, git worktrees |
| [EVENTS.md](./EVENTS.md) | Event sourcing patterns |
| [AGENTS.md](./AGENTS.md) | Agent roles and responsibilities |
| [MESSAGING.md](./MESSAGING.md) | Mail protocol |
| [WORKFLOWS.md](./WORKFLOWS.md) | Molecule state machine |
| [VERIFICATION.md](./VERIFICATION.md) | VerMAS Inspector pipeline |
| [EVALUATION.md](./EVALUATION.md) | Metrics and testing |
