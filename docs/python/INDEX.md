# VerMAS Python Documentation

> Multi-agent verification system built on proven technologies

## Quick Start

1. Read [HOW_IT_WORKS.md](./HOW_IT_WORKS.md) - Core concepts and quick start
2. Review [ARCHITECTURE.md](./ARCHITECTURE.md) - System design
3. Reference [CLI.md](./CLI.md) - Command reference

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
│   CLI.md              Unix-style CLI tools (gt, bd)                         │
│   SCHEMAS.md          JSONL/TOML data specifications                        │
│                                                                             │
│   AGENT SYSTEM                                                              │
│   ────────────────────────────────────────────────────────                  │
│   AGENTS.md           Agent roles (Mayor, Witness, Polecat, etc.)           │
│   MESSAGING.md        Mail protocol between agents                          │
│   WORKFLOWS.md        Molecule state machine (MEOW)                         │
│                                                                             │
│   VERIFICATION                                                              │
│   ────────────────────────────────────────────────────────                  │
│   VERIFICATION.md     VerMAS Inspector pipeline                             │
│   EVALUATION.md       Metrics, testing, benchmarks                          │
│                                                                             │
│   OPERATIONS                                                                │
│   ────────────────────────────────────────────────────────                  │
│   OPERATIONS.md       Deployment, startup, monitoring, maintenance          │
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
     ├─→ issues.jsonl (projection)
     ├─→ messages.jsonl (projection)
     └─→ feed.jsonl (change feed)
```

### GUPP (Propulsion Principle)

> If your hook has work, RUN IT.

Agents check their hook on startup and execute immediately. No confirmation, no questions. See [HOOKS.md](./HOOKS.md).

### Molecules (MEOW)

Workflow state machine with phases:
- **Ice-9** (Formula) → **Solid** (Proto) → **Liquid** (Mol) → **Archive**

See [WORKFLOWS.md](./WORKFLOWS.md).

### Mail Protocol

Async messaging between agents via JSONL mailboxes. See [MESSAGING.md](./MESSAGING.md).

---

## By Role

### For Mayors
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
| [AGENTS.md](./AGENTS.md) | Agent roles and responsibilities |
| [ARCHITECTURE.md](./ARCHITECTURE.md) | System architecture |
| [CLI.md](./CLI.md) | Command reference |
| [EVALUATION.md](./EVALUATION.md) | Metrics and evaluation |
| [EVENTS.md](./EVENTS.md) | Event sourcing |
| [GO-VS-PYTHON.md](./GO-VS-PYTHON.md) | Language comparison |
| [HOOKS.md](./HOOKS.md) | Claude Code integration |
| [HOW_IT_WORKS.md](./HOW_IT_WORKS.md) | Quick start guide |
| [MESSAGING.md](./MESSAGING.md) | Mail protocol |
| [OPERATIONS.md](./OPERATIONS.md) | Deployment and operations |
| [SCHEMAS.md](./SCHEMAS.md) | JSONL/TOML data specs |
| [VERIFICATION.md](./VERIFICATION.md) | VerMAS Inspector pipeline |
| [WORKFLOWS.md](./WORKFLOWS.md) | Molecule system |

**Total: 14 documents**

---

## Design Principles

1. **No API Costs** - All LLM via Claude Code CLI
2. **Git-Backed Everything** - State lives in version control
3. **Event Sourced** - Append-only logs, derived state
4. **Unix Philosophy** - Small tools, text streams, composable
5. **Observable** - Attach to any agent, grep any log
6. **Recoverable** - Hooks persist, worktrees survive crashes

---

## Glossary

| Term | Definition |
|------|------------|
| **Bead** | A work item (issue, task, bug, feature). Stored in `issues.jsonl`. |
| **BD_ACTOR** | Environment variable identifying an agent (e.g., `gastown/polecats/slot0`). |
| **Convoy** | A group of related beads traveling together through the system. |
| **Crew** | Human-directed workspace within a rig. Persistent, not ephemeral. |
| **Deacon** | Daemon process managing agent lifecycle and system health. |
| **Event** | An immutable record of a state change. Stored in `events.jsonl`. |
| **Feed** | Real-time stream of events (`feed.jsonl`). Agents can tail this. |
| **Formula** | A TOML template defining a workflow (Ice-9 state). |
| **GUPP** | Gas Town Universal Propulsion Principle: "If your hook has work, RUN IT." |
| **Hook** | A file (`.hook-{agent}`) containing an agent's assigned work reference. |
| **Inspector** | Verification agent roles: Designer, Strategist, Verifier, Auditor, Advocate, Critic, Judge. |
| **JSONL** | JSON Lines format. One JSON object per line, append-only. |
| **Mayor** | Global coordinator. Dispatches work, handles escalations. Does NOT write code. |
| **MEOW** | Molecule state machine: Ice-9 → Solid → Liquid → Archive. |
| **Molecule (Mol)** | A running workflow instance (Liquid state). Tracks step progress. |
| **Polecat** | Ephemeral worker agent. Spawns, executes one bead, disappears. |
| **Projection** | Derived state from events (e.g., `issues.jsonl` is a projection of bead events). |
| **Protomolecule** | A frozen, reusable workflow template (Solid state). |
| **Refinery** | Per-rig agent that manages merge queues and code review. |
| **Rig** | A project container. Has its own beads, polecats, witness, refinery. |
| **Sling** | To assign work to an agent (`gt sling <bead> <rig>`). |
| **Town** | The workspace root containing all rigs and the Mayor. |
| **VerMAS** | Verification Multi-Agent System. The Inspector pipeline for code verification. |
| **Wisp** | An ephemeral workflow instance (Vapor state). No permanent record. |
| **Witness** | Per-rig agent that monitors polecats. Nudges idle, kills stuck, escalates. |
| **Worktree** | A git worktree. Each polecat gets its own for isolation. |
| **Adversarial Review** | Verification pattern: Advocate argues PASS, Critic argues FAIL, Judge decides. |
| **Change Feed** | Real-time event stream (`feed.jsonl`) for agent coordination. |
| **Correlation ID** | Links related events across agents and workflows. |
| **LLM Backend** | Abstraction allowing Claude, Codex, Aider, or custom CLI tools. |
| **Slot** | A polecat work position (slot0-slot4). Limited per rig. |
| **Squash** | Archive a molecule with summary. Creates permanent record. |
| **Burn** | Discard a molecule without record. |

---

## Quick Reference Card

### Startup Sequence
```bash
gt deacon start          # Start infrastructure
claude --profile mayor   # Start Mayor session
```

### Daily Workflow
```bash
gt hook                  # Check assigned work
gt mail inbox            # Check messages
bd ready                 # Find available work
gt sling <bead> <rig>    # Dispatch work
gt convoy list           # Monitor progress
```

### Session End
```bash
git add . && git commit -m "..."
bd sync && git push
gt handoff -m "..."      # If incomplete
```

### Debugging
```bash
gt status                # Town overview
gt polecat list          # Active workers
tail -f .beads/feed.jsonl | jq .  # Event stream
tmux attach -t <session> # Watch agent
```

---

## Getting Help

```bash
gt --help              # Gas Town commands
bd --help              # Beads commands
gt <command> --help    # Specific command help
bd <command> --help
```
