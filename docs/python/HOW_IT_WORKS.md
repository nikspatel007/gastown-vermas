# Gas Town - How It Works

> Source: [Steve Yegge's Blog](https://steve-yegge.medium.com/welcome-to-gas-town-4f25ee16dd04) | [GitHub Repo](https://github.com/steveyegge/gastown)

## Overview

Gas Town is a **multi-agent workspace manager** that coordinates multiple Claude Code agents working in parallel. It's designed to scale comfortably to 20-30 agents through structured coordination.

Think of it as a steam engine: **when an agent finds work on their hook, they EXECUTE**. No confirmation. No questions. No waiting.

## Core Metaphor: The Steam Engine

- Gas Town is a steam engine
- The Mayor is the main drive shaft
- If the Mayor stalls, the whole town stalls
- Work flows through "hooks" - assigned task queues
- Agents wake up, check their hook, and run whatever's there

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

## Why Gas Town?

Managing 4-10 agents creates chaos. Gas Town enables scaling to 20-30 agents through:
- Structured coordination via Mayor
- Persistent work state in Beads (survives crashes)
- Clear role separation
- Automated handoffs
- Git-backed everything
