# VerMAS Python Architecture

> System design for CLI-based multi-agent verification

## Design Principles

1. **No API costs** - All LLM interactions through Claude Code CLI
2. **Shared data formats** - Interoperable with Go implementation via JSONL/TOML
3. **Tmux isolation** - Each agent runs in its own tmux session
4. **Git-backed state** - All persistent state lives in git
5. **Hook-driven execution** - GUPP: "If your hook has work, RUN IT"

---

## System Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              TOWN (Workspace Root)                           │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                         COORDINATION LAYER                          │  │
│   │                                                                     │  │
│   │   Mayor ◄────────────────────────────────────────────► Deacon      │  │
│   │   (Human-directed)                                    (Daemon)     │  │
│   │                                                                     │  │
│   │   - Cross-rig decisions                              - Health mon  │  │
│   │   - Strategic planning                               - Restarts    │  │
│   │   - Escalation handling                              - Watchdog    │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                    │                                        │
│                    ┌───────────────┴───────────────┐                       │
│                    ▼                               ▼                        │
│   ┌────────────────────────────┐   ┌────────────────────────────┐         │
│   │         RIG A              │   │         RIG B              │         │
│   │                            │   │                            │         │
│   │   ┌────────┐ ┌─────────┐  │   │   ┌────────┐ ┌─────────┐  │         │
│   │   │Witness │ │Refinery │  │   │   │Witness │ │Refinery │  │         │
│   │   └───┬────┘ └────┬────┘  │   │   └───┬────┘ └────┬────┘  │         │
│   │       │           │       │   │       │           │       │         │
│   │       ▼           ▼       │   │       ▼           ▼       │         │
│   │   ┌─────────────────┐     │   │   ┌─────────────────┐     │         │
│   │   │    Polecats     │     │   │   │    Polecats     │     │         │
│   │   │  slot0..slot4   │     │   │   │  slot0..slot4   │     │         │
│   │   └─────────────────┘     │   │   └─────────────────┘     │         │
│   │                            │   │                            │         │
│   │   .beads/ (rig-level)     │   │   .beads/ (rig-level)     │         │
│   └────────────────────────────┘   └────────────────────────────┘         │
│                                                                             │
│   .beads/ (town-level: mayor mail, HQ coordination)                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Data Flow

### Work Assignment (Sling)

```
Mayor                    Rig                     Polecat
  │                       │                        │
  │  1. gt sling bead rig │                        │
  │──────────────────────▶│                        │
  │                       │                        │
  │                       │  2. Allocate slot      │
  │                       │  3. Create worktree    │
  │                       │  4. Write hook file    │
  │                       │  5. Start tmux session │
  │                       │───────────────────────▶│
  │                       │                        │
  │                       │                        │  6. GUPP: Check hook
  │                       │                        │  7. Find work
  │                       │                        │  8. EXECUTE
  │                       │                        │
```

### Work Completion

```
Polecat                 Witness                Refinery
  │                        │                       │
  │  1. gt polecat done    │                       │
  │  (sends POLECAT_DONE)  │                       │
  │───────────────────────▶│                       │
  │                        │                       │
  │                        │  2. Validate work     │
  │                        │  3. Send MERGE_READY  │
  │                        │──────────────────────▶│
  │                        │                       │
  │                        │                       │  4. Run tests
  │                        │                       │  5. Run verification
  │                        │                       │  6. Merge or REWORK
  │                        │                       │
  │  7. Release slot       │                       │
  │◀───────────────────────│                       │
  │                        │                       │
```

---

## Storage Architecture

### Two-Level Beads

| Level | Location | Purpose | Git Behavior |
|-------|----------|---------|--------------|
| **Town** | `~/.beads/` | Mayor mail, HQ coordination | Commits to main |
| **Rig** | `<rig>/.beads/` | Project issues, workflows | Uses beads-sync branch |

### File Types

| File | Format | Contents |
|------|--------|----------|
| `issues.jsonl` | JSONL | Beads (issues, tasks, bugs) |
| `messages.jsonl` | JSONL | Mail between agents |
| `routes.jsonl` | JSONL | Prefix → rig routing |
| `formulas/*.toml` | TOML | Workflow templates |
| `mols/*.json` | JSON | Active workflow instances |
| `.hook-{agent}` | Plain text | Current hook assignment |

### Prefix-Based Routing

Every bead ID has a prefix (e.g., `gt-abc12`). The router maps prefixes to rigs:

```
gt-* → gastown/.beads/
hq-* → town/.beads/ (mayor level)
pa-* → project-a/.beads/
```

---

## Session Architecture

### Tmux Session Naming

```
{role}-{rig}           # Witness, Refinery
{role}-{rig}-{slot}    # Polecats
mayor                  # Town-level Mayor
deacon                 # Town-level Deacon
```

Examples:
- `witness-gastown`
- `refinery-gastown`
- `polecat-gastown-slot0`
- `mayor`

### Session Environment

Each session sets:
- `BD_ACTOR` - Agent identity (e.g., `gastown/polecats/slot0`)
- `BEAD_ID` - Assigned work (for polecats)
- Working directory - Appropriate worktree or rig path

### Claude Code Profile

Each role has a profile that loads its CLAUDE.md:
- `claude --profile witness`
- `claude --profile refinery`
- `claude --profile polecat`
- `claude --profile mayor`

---

## Verification Integration (VerMAS)

### Where Verification Happens

```
Polecat completes work
        │
        ▼
    Witness
        │
        ▼
    Refinery ────────────────────────────────────────┐
        │                                            │
        ▼                                            ▼
    Run Tests                               Run VerMAS Inspector
        │                                            │
        │                                            ▼
        │                                   ┌─────────────────┐
        │                                   │    Designer     │
        │                                   │  (elaborate)    │
        │                                   └────────┬────────┘
        │                                            │
        │                                   ┌────────▼────────┐
        │                                   │   Strategist    │
        │                                   │  (plan tests)   │
        │                                   └────────┬────────┘
        │                                            │
        │                                   ┌────────▼────────┐
        │                                   │    Verifier     │
        │                                   │ (run shell/no LLM)│
        │                                   └────────┬────────┘
        │                                            │
        │                                   ┌────────▼────────┐
        │                                   │    Auditor      │
        │                                   │ (LLM if needed) │
        │                                   └────────┬────────┘
        │                                            │
        │                                   ┌────────▼────────┐
        │                                   │   Adversarial   │
        │                                   │ Advocate/Critic │
        │                                   │     Judge       │
        │                                   └────────┬────────┘
        │                                            │
        ▼                                            ▼
    All pass? ◄──────────────────────────── Verdict: PASS/FAIL
        │
        ├── Yes → Merge
        └── No  → REWORK_REQUEST
```

### Inspector as Separate Sessions

Each Inspector role can run as its own tmux session:
- `inspector-designer-{rig}`
- `inspector-strategist-{rig}`
- `inspector-verifier-{rig}` (no LLM - just runs shell)
- `inspector-advocate-{rig}`
- `inspector-critic-{rig}`
- `inspector-judge-{rig}`

Or as a single orchestrated workflow within Refinery.

---

## Logging Strategy

### What Gets Logged

| Log Type | Location | Contents |
|----------|----------|----------|
| **Session logs** | `logs/{session}/` | Full Claude Code output |
| **Mail archive** | `.beads/messages.jsonl` | All agent communication |
| **Bead history** | `.beads/issues.jsonl` | Work state changes |
| **Molecule traces** | `.beads/mols/*.json` | Workflow step execution |
| **Verification evidence** | `.beads/evidence/` | Test outputs, verdicts |

### Log Levels

1. **Trace** - Every Claude Code interaction (large, debugging only)
2. **Debug** - Step-by-step workflow execution
3. **Info** - Major state changes (bead status, merges)
4. **Warn** - Nudges, retries, recoverable issues
5. **Error** - Failures, escalations, stuck agents

### Structured Logging Fields

- `timestamp` - When it happened
- `actor` - BD_ACTOR of the agent
- `event` - What happened (bead_created, mail_sent, step_completed)
- `bead_id` - Related bead if any
- `mol_id` - Related molecule if any
- `details` - Event-specific data

---

## Failure Modes and Recovery

### Agent Failures

| Failure | Detection | Recovery |
|---------|-----------|----------|
| Polecat stuck | Witness patrol (idle >15min) | Kill session, release slot |
| Witness down | Deacon patrol | Restart Witness |
| Refinery down | Deacon patrol | Restart Refinery |
| Deacon down | Boot process | Restart Deacon |

### Work Recovery

| Scenario | State | Recovery |
|----------|-------|----------|
| Polecat killed mid-work | Sandbox has uncommitted changes | New polecat can resume from worktree |
| Session crash | Hook file persists | New session reads hook, continues |
| Merge failed | Bead still open | REWORK_REQUEST sent, work continues |

### Watchdog Chain

```
OS (systemd/launchd)
        │
        ▼
     Boot
        │
        ▼
    Deacon ──────────────────┬──────────────────┐
        │                    │                  │
        ▼                    ▼                  ▼
   Witness(A)           Witness(B)         Witness(C)
        │                    │                  │
        ▼                    ▼                  ▼
   Polecats              Polecats           Polecats
```

---

## Key Design Decisions

### Why Tmux?

1. **Isolation** - Each agent has clean environment
2. **Persistence** - Sessions survive disconnects
3. **Observability** - Can attach and watch any agent
4. **Simplicity** - No container orchestration needed

### Why JSONL?

1. **Append-only** - Easy concurrent writes
2. **Git-friendly** - Line-based diffs
3. **Interoperable** - Go and Python read same files
4. **Debuggable** - Human-readable, greppable

### Why Hooks?

1. **Stateless startup** - Agent reads hook, knows what to do
2. **Crash recovery** - Hook persists, work doesn't get lost
3. **Explicit assignment** - No ambiguity about who does what
4. **GUPP enforcement** - Work on hook = immediate execution

### Why Molecules?

1. **Reusable workflows** - Define once, instantiate many
2. **Step tracking** - Know exactly where work stopped
3. **Dependencies** - Steps can depend on other steps
4. **Audit trail** - Archive completed workflows

---

## See Also

- [AGENTS.md](./AGENTS.md) - Agent roles and responsibilities
- [WORKFLOWS.md](./WORKFLOWS.md) - Molecule state machine
- [MESSAGING.md](./MESSAGING.md) - Communication patterns
- [EVALUATION.md](./EVALUATION.md) - How to evaluate the system
