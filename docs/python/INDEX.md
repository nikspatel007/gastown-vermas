# VerMAS Python Documentation

> Multi-agent verification system built on proven technologies

## Quick Start

1. Read [HOW_IT_WORKS.md](./HOW_IT_WORKS.md) - Core concepts and quick start
2. Review [ARCHITECTURE.md](./ARCHITECTURE.md) - System design
3. Reference [CLI.md](./CLI.md) - Command reference

---

## The Company Metaphor

VerMAS uses a **Company with Factories** model:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              THE COMPANY                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   HEADQUARTERS                                                              │
│   ─────────────                                                             │
│   CEO             Strategic leader, coordinates across factories            │
│   Operations      Infrastructure, keeps systems running                     │
│   Board           Human oversight, escalations                              │
│                                                                             │
│   FACTORY (per repository/project)                                          │
│   ─────────────────────────────────                                         │
│   Supervisor      Monitors workers, handles issues                          │
│   QA Department   Quality control, merge queue, verification                │
│   Workers         Ephemeral task executors (spawn → work → done)            │
│   Teams           Human-directed workspaces                                 │
│                                                                             │
│   WORK MANAGEMENT                                                           │
│   ───────────────                                                           │
│   Work Orders     Tasks, bugs, features (stored in JSONL)                   │
│   Assignments     Work assigned to a worker                                 │
│   Processes       Workflow templates and instances                          │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Documentation Map

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         DOCUMENTATION STRUCTURE                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   FOUNDATIONS                                                               │
│   ────────────────────────────────────────────────────────                  │
│   HOW_IT_WORKS.md     Quick start, core concepts                            │
│   ARCHITECTURE.md     System design, data flow                              │
│   GO-VS-PYTHON.md     Language comparison, code examples                    │
│                                                                             │
│   PROVEN TECHNOLOGIES                                                       │
│   ────────────────────────────────────────────────────────                  │
│   EVENTS.md           Event sourcing (JSONL append-only logs)               │
│   HOOKS.md            Claude Code integration, git worktrees                │
│   CLI.md              Unix-style CLI tools                                  │
│   SCHEMAS.md          JSONL/TOML data specifications                        │
│                                                                             │
│   ORGANIZATION                                                              │
│   ────────────────────────────────────────────────────────                  │
│   AGENTS.md           Roles (CEO, Supervisor, Worker, etc.)                 │
│   MESSAGING.md        Internal communications                               │
│   WORKFLOWS.md        Process templates and execution                       │
│                                                                             │
│   VERIFICATION                                                              │
│   ────────────────────────────────────────────────────────                  │
│   VERIFICATION.md     Quality Assurance pipeline                            │
│   EVALUATION.md       Metrics, testing, benchmarks                          │
│                                                                             │
│   OPERATIONS                                                                │
│   ────────────────────────────────────────────────────────                  │
│   OPERATIONS.md       Deployment, startup, monitoring                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Proven Technology Stack

VerMAS is built entirely on battle-tested technologies:

| Layer | Technology | Why |
|-------|------------|-----|
| **Storage** | File system (JSONL) | Append-only, git-friendly, debuggable |
| **Version Control** | Git (worktrees, branches) | Isolation, history, collaboration |
| **Agent Runtime** | Claude Code CLI | No API costs, profiles, hooks |
| **Process Isolation** | Tmux sessions | Persistent, observable, simple |
| **CLI Framework** | Typer (Python) | Type-safe, auto-docs, composable |
| **Data Validation** | Pydantic | Runtime type checking, serialization |
| **Async** | asyncio | Concurrent agent coordination |

### Why These Choices?

**File System over Database**
- No external dependencies
- Human-readable (grep, cat, jq)
- Git-native (commits, diffs, history)
- Crash recovery via append-only logs

**Git over Custom VCS**
- Worktrees for agent isolation
- Branches for parallel work
- Built-in conflict resolution
- Universal tooling

**CLI over API**
- Zero runtime costs
- Unix composability (pipes, scripts)
- Debuggable (just run the command)
- Existing muscle memory

**Tmux over Containers**
- Instant startup
- Easy observation (attach/detach)
- Session persistence
- Minimal overhead

---

## Core Concepts

### Event Sourcing

All state changes are immutable events. See [EVENTS.md](./EVENTS.md).

```
events.jsonl (source of truth)
     │
     ├─→ work_orders.jsonl (projection)
     ├─→ messages.jsonl (projection)
     └─→ feed.jsonl (change feed)
```

### Assignment Principle

> If you have an assignment, EXECUTE IT.

Agents check their assignment on startup and execute immediately. No confirmation, no questions. See [HOOKS.md](./HOOKS.md).

### Processes (Workflows)

Workflow state machine with phases:
- **Template** (TOML) → **Ready** (compiled) → **Active** (running) → **Archive**

See [WORKFLOWS.md](./WORKFLOWS.md).

### Internal Communications

Async messaging between agents via JSONL mailboxes. See [MESSAGING.md](./MESSAGING.md).

---

## By Role

### For CEOs
- [HOW_IT_WORKS.md](./HOW_IT_WORKS.md) - Overview
- [CLI.md](./CLI.md) - Command reference
- [AGENTS.md](./AGENTS.md) - Role responsibilities

### For Developers
- [ARCHITECTURE.md](./ARCHITECTURE.md) - System design
- [EVENTS.md](./EVENTS.md) - Event sourcing
- [HOOKS.md](./HOOKS.md) - Claude Code integration

### For Evaluators
- [EVALUATION.md](./EVALUATION.md) - Metrics and testing
- [EVENTS.md](./EVENTS.md) - Computing metrics from events

---

## Document Inventory

| Document | Description |
|----------|-------------|
| [AGENTS.md](./AGENTS.md) | Roles and responsibilities |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | System architecture |
| [CLI.md](./CLI.md) | Command reference |
| [EVALUATION.md](./EVALUATION.md) | Metrics and evaluation |
| [EVENTS.md](./EVENTS.md) | Event sourcing |
| [GO-VS-PYTHON.md](./GO-VS-PYTHON.md) | Language comparison |
| [HOOKS.md](./HOOKS.md) | Claude Code integration |
| [HOW_IT_WORKS.md](./HOW_IT_WORKS.md) | Quick start guide |
| [MESSAGING.md](./MESSAGING.md) | Internal communications |
| [OPERATIONS.md](./OPERATIONS.md) | Deployment and operations |
| [SCHEMAS.md](./SCHEMAS.md) | JSONL/TOML data specs |
| [VERIFICATION.md](./VERIFICATION.md) | QA pipeline |
| [WORKFLOWS.md](./WORKFLOWS.md) | Process system |

**Total: 14 documents**

---

## Design Principles

1. **No API Costs** - All LLM via Claude Code CLI
2. **Git-Backed Everything** - State lives in version control
3. **Event Sourced** - Append-only logs, derived state
4. **Unix Philosophy** - Small tools, text streams, composable
5. **Observable** - Attach to any agent, grep any log
6. **Recoverable** - Assignments persist, worktrees survive crashes

---

## Glossary

| Term | Definition |
|------|------------|
| **Assignment** | Work assigned to an agent (stored in `.assignment-{agent}` file). |
| **CEO** | Company coordinator. Dispatches work, handles escalations. Does NOT write code. |
| **Company** | The workspace root containing all factories and the CEO. |
| **Event** | An immutable record of a state change. Stored in `events.jsonl`. |
| **Factory** | A project/repository container. Has its own workers, supervisor, QA. |
| **Feed** | Real-time stream of events (`feed.jsonl`). Agents can tail this. |
| **JSONL** | JSON Lines format. One JSON object per line, append-only. |
| **Operations** | Daemon process managing agent lifecycle and system health. |
| **Process** | A running workflow instance. Tracks step progress. |
| **Projection** | Derived state from events (e.g., `work_orders.jsonl` is a projection). |
| **QA Department** | Per-factory agent that manages merge queues and verification. |
| **Slot** | A worker position (slot0-slot4). Limited per factory. |
| **Sprint** | A group of related work orders traveling together. |
| **Supervisor** | Per-factory agent that monitors workers. Nudges idle, escalates stuck. |
| **Team** | Human-directed workspace within a factory. Persistent, not ephemeral. |
| **Template** | A TOML file defining a workflow. |
| **VerMAS** | Verification Multi-Agent System. The QA pipeline for code verification. |
| **Work Order** | A work item (task, bug, feature). Stored in `work_orders.jsonl`. |
| **Worker** | Ephemeral agent. Spawns, executes one work order, disappears. |
| **Worktree** | A git worktree. Each worker gets its own for isolation. |
| **Adversarial Review** | QA pattern: Advocate argues PASS, Critic argues FAIL, Judge decides. |
| **Change Feed** | Real-time event stream for agent coordination. |
| **Correlation ID** | Links related events across agents and workflows. |
| **LLM Backend** | Abstraction allowing Claude, Codex, Aider, or custom CLI tools. |

---

## Quick Reference Card

### Startup Sequence
```bash
co ops start             # Start infrastructure
claude --profile ceo     # Start CEO session
```

### Daily Workflow
```bash
co assignment            # Check assigned work
co inbox                 # Check messages
wo ready                 # Find available work
co dispatch <wo> <factory>  # Dispatch work
co sprints               # Monitor progress
```

### Session End
```bash
git add . && git commit -m "..."
wo sync && git push
co handoff -m "..."      # If incomplete
```

### Debugging
```bash
co status                # Company overview
co workers               # Active workers
tail -f .work/feed.jsonl | jq .  # Event stream
tmux attach -t <session> # Watch agent
```

---

## Getting Help

```bash
co --help              # Company commands
wo --help              # Work Order commands
co <command> --help    # Specific command help
wo <command> --help
```
