# VerMAS Python Architecture

> System design for CLI-based multi-agent verification

**See also:** [INDEX.md](./INDEX.md) for documentation map

## Design Principles

1. **No API costs** - All LLM interactions through Claude Code CLI
2. **Shared data formats** - Interoperable with Go implementation via JSONL/TOML
3. **Tmux isolation** - Each agent runs in its own tmux session
4. **Git-backed state** - All persistent state lives in git
5. **Hook-driven execution** - GUPP: "If your hook has work, RUN IT"
6. **Event sourced** - All state derived from append-only event logs

## Proven Technology Stack

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      PROVEN TECHNOLOGY STACK                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │  STORAGE: File System                                                 │ │
│  │  ─────────────────────                                                │ │
│  │  • JSONL files (append-only, one record per line)                    │ │
│  │  • Git for version control and sync                                  │ │
│  │  • No database required                                              │ │
│  │  • Human-readable (grep, cat, jq)                                    │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │  ISOLATION: Git Worktrees + Tmux                                      │ │
│  │  ───────────────────────────────                                      │ │
│  │  • Each agent gets own worktree (parallel development)               │ │
│  │  • Each agent runs in tmux session (process isolation)               │ │
│  │  • Sessions survive disconnects                                       │ │
│  │  • Easy observation (tmux attach)                                    │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │  AGENTS: Claude Code CLI                                              │ │
│  │  ───────────────────────                                              │ │
│  │  • Profiles define agent roles (CLAUDE.md)                           │ │
│  │  • Hooks for lifecycle events                                        │ │
│  │  • No API costs (uses subscription)                                  │ │
│  │  • Full Claude capabilities                                          │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │  CLI: Typer (Python)                                                  │ │
│  │  ───────────────────                                                  │ │
│  │  • Type-safe commands                                                │ │
│  │  • Auto-generated help                                               │ │
│  │  • Unix composable (pipes, scripts)                                  │ │
│  │  • JSON/JSONL output modes                                           │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐ │
│  │  DATA FLOW: Event Sourcing                                            │ │
│  │  ─────────────────────────                                            │ │
│  │  • All changes as immutable events                                   │ │
│  │  • Projections for current state                                     │ │
│  │  • Change feed for real-time updates                                 │ │
│  │  • Temporal queries ("what was state at T?")                         │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Why These Technologies?

| Technology | Alternative | Why We Chose This |
|------------|-------------|-------------------|
| JSONL files | PostgreSQL, SQLite | No setup, git-native, grep-able |
| Git worktrees | Docker containers | Instant, no daemon, git-integrated |
| Tmux | Kubernetes, systemd | Simple, observable, universal |
| Claude Code CLI | API calls | No per-request costs, profiles |
| Event sourcing | CRUD | Audit trail, debugging, recovery |

---

## LLM Backend Abstraction

VerMAS supports multiple LLM backends through a unified interface. The architecture is **LLM-agnostic** - agents interact through CLI tools, not direct API calls.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         LLM BACKEND ABSTRACTION                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Agent (Tmux Session)                                                      │
│        │                                                                    │
│        │  Runs CLI tool (e.g., "claude", "codex", "aider")                 │
│        ▼                                                                    │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                      LLM CLI Abstraction                             │  │
│   │                                                                     │  │
│   │   Interface: stdin/stdout, exit codes, file I/O                    │  │
│   │   Contract: Read context, produce output, use tools                │  │
│   │                                                                     │  │
│   └───────────────────────────┬─────────────────────────────────────────┘  │
│                               │                                            │
│         ┌─────────────────────┼─────────────────────┐                      │
│         │                     │                     │                      │
│         ▼                     ▼                     ▼                      │
│   ┌───────────┐        ┌───────────┐        ┌───────────┐                  │
│   │  Claude   │        │  Codex    │        │   Aider   │                  │
│   │  Code     │        │  CLI      │        │   CLI     │                  │
│   │           │        │           │        │           │                  │
│   │ Anthropic │        │  OpenAI   │        │  Any LLM  │                  │
│   └───────────┘        └───────────┘        └───────────┘                  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Supported Backends

| Backend | CLI Tool | Configuration |
|---------|----------|---------------|
| **Claude Code** | `claude` | `--profile <role>` |
| **OpenAI Codex** | `codex` | Via environment |
| **Aider** | `aider` | `--model <model>` |
| **Custom** | Any CLI | Implements contract |

### Backend Contract

Any LLM CLI tool can be used if it:
1. Accepts work context via stdin, files, or arguments
2. Produces output to stdout/files
3. Returns exit code 0 on success
4. Can read/write to the filesystem (for tools)

### Configuring Backend

```toml
# .beads/config.toml

[llm]
# Default backend for all agents
backend = "claude"
command = "claude --profile {role}"

# Per-role overrides
[llm.roles.polecat]
backend = "claude"
command = "claude --profile polecat"

[llm.roles.verifier]
# Verifier uses no LLM - just shell execution
backend = "none"

[llm.roles.inspector]
# Use different model for verification
backend = "claude"
command = "claude --profile inspector --model opus"
```

### Why CLI-Based?

1. **No API key management** - CLIs handle auth
2. **Cost control** - Subscription vs per-token
3. **Switchable** - Change backend without code changes
4. **Observable** - See exactly what runs in tmux
5. **Debuggable** - Replay commands manually

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
| `events.jsonl` | JSONL | **Event log (source of truth)** |
| `issues.jsonl` | JSONL | Beads - projection of events |
| `messages.jsonl` | JSONL | Mail - projection of events |
| `feed.jsonl` | JSONL | Real-time change feed |
| `routes.jsonl` | JSONL | Prefix → rig routing |
| `formulas/*.toml` | TOML | Workflow templates |
| `mols/*.json` | JSON | Active workflow instances |
| `.hook-{agent}` | Plain text | Current hook assignment |

### Event Sourcing Model

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          EVENT SOURCING MODEL                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Commands (gt, bd)                                                         │
│        │                                                                    │
│        ▼                                                                    │
│   ┌─────────────┐                                                           │
│   │   events.   │  ← Source of truth (append-only)                         │
│   │   jsonl     │                                                           │
│   └──────┬──────┘                                                           │
│          │                                                                  │
│          │ project                                                          │
│          ▼                                                                  │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                    │
│   │  issues.    │    │ messages.   │    │  feed.      │                    │
│   │  jsonl      │    │ jsonl       │    │  jsonl      │                    │
│   │             │    │             │    │             │                    │
│   │ (current    │    │ (mailbox    │    │ (real-time  │                    │
│   │  state)     │    │  state)     │    │  stream)    │                    │
│   └─────────────┘    └─────────────┘    └─────────────┘                    │
│                                                                             │
│   All state is derived from events. Events are immutable.                  │
│   See EVENTS.md for full event sourcing documentation.                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

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

## Security Considerations

### Threat Model

VerMAS runs on a single machine or trusted cluster. The threat model assumes:
- **Trusted agents** - All LLM agents are under your control
- **Trusted filesystem** - No malicious writes to `.beads/`
- **Network isolation** - Agents don't expose network services

### Security Boundaries

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         SECURITY ARCHITECTURE                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   TRUSTED ZONE (your machine)                                               │
│   ┌───────────────────────────────────────────────────────────────────────┐ │
│   │                                                                       │ │
│   │   Agents (tmux)     Beads (files)      Git (worktrees)               │ │
│   │        │                 │                    │                       │ │
│   │        └─────────────────┼────────────────────┘                       │ │
│   │                          │                                            │ │
│   │   All run as YOUR user with YOUR permissions                         │ │
│   │                                                                       │ │
│   └───────────────────────────────────────────────────────────────────────┘ │
│                                    │                                        │
│                                    │ git push/pull                          │
│                                    ▼                                        │
│   EXTERNAL (network)                                                        │
│   ┌───────────────────────────────────────────────────────────────────────┐ │
│   │   GitHub/GitLab         LLM APIs (Claude, OpenAI)                     │ │
│   │   - Code sync           - Model inference                             │ │
│   │   - No secrets in repo  - API keys in env                             │ │
│   └───────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Best Practices

**File Permissions:**
```bash
# Beads should be readable/writable by your user only
chmod 700 .beads
chmod 600 .beads/*.jsonl
```

**Secrets Management:**
- Never store API keys in beads or git
- Use environment variables or secure vaults
- LLM CLIs handle their own auth

**Git Security:**
```bash
# Add to .gitignore
.beads/.hook-*      # Hook files (local state)
.beads/feed.jsonl   # Real-time feed (local)
logs/               # Session logs (may contain sensitive output)
```

**Agent Isolation:**
- Each polecat runs in isolated worktree
- Agents can't access each other's worktrees directly
- All communication through mail (auditable)

**Audit Trail:**
- All state changes are events
- Events include actor (who did it)
- Git history provides additional audit

### What Agents Can Do

Agents run with your user permissions. They CAN:
- Read/write files in their worktree
- Execute shell commands
- Make network requests (LLM API)
- Read environment variables

Agents are LIMITED by:
- Worktree isolation (can't see other worktrees)
- CLAUDE.md instructions (behavioral constraints)
- Hook-based work assignment (explicit scope)

### Verification as Security

The VerMAS verification pipeline provides security benefits:
- Code review before merge (Inspector roles)
- Automated testing (Verifier)
- Adversarial review (Advocate/Critic)
- Audit trail of all decisions

---

## Scaling and Performance

### Scaling Dimensions

| Dimension | Typical Range | Bottleneck |
|-----------|---------------|------------|
| **Agents per rig** | 5-10 polecats | Tmux sessions, disk I/O |
| **Rigs per town** | 3-10 | Memory, coordination overhead |
| **Total agents** | 20-50 | LLM rate limits, human oversight |
| **Beads per rig** | 1000s | JSONL scan time |

### Performance Characteristics

**Fast Operations (< 100ms):**
- Hook check (`gt hook`)
- Mail send (`gt mail send`)
- Event emit (append to JSONL)

**Medium Operations (100ms - 1s):**
- Bead lookup by ID (`bd show`)
- List beads with filter (`bd list`)
- Sync status check (`bd sync --status`)

**Slow Operations (> 1s):**
- Full beads sync (`bd sync`)
- Worktree creation (`git worktree add`)
- LLM agent spawn (Claude startup)

### Optimizations

**Event Log Partitioning:**
```
.beads/
├── events.jsonl           # Current day
└── events/
    ├── 2026-01-06.jsonl   # Archived by day
    └── 2026-01-05.jsonl
```

**Index Files (optional):**
```python
# For large bead counts, maintain index
# .beads/index/by-status.json
{
  "open": ["gt-abc", "gt-def", ...],
  "closed": ["gt-xyz", ...]
}
```

**Projection Caching:**
- `issues.jsonl` is a projection, not source of truth
- Can be regenerated from events
- Cache for fast reads, events for writes

### Scaling Recommendations

**< 10 agents:** Default configuration works fine

**10-30 agents:**
- Partition events by day
- Use separate rigs for independent projects
- Monitor disk I/O

**30+ agents:**
- Consider multiple machines
- Implement event archival
- Add monitoring (Deacon metrics)

### Resource Usage

| Component | Memory | Disk | CPU |
|-----------|--------|------|-----|
| Tmux session | ~5MB | - | Idle |
| Claude Code | ~200MB | - | Varies |
| Beads sync | ~50MB | ~1MB/1000 beads | Low |
| Event tail | ~10MB | Append only | Low |

---

## See Also

- [INDEX.md](./INDEX.md) - Documentation map and glossary
- [HOW_IT_WORKS.md](./HOW_IT_WORKS.md) - Quick start guide
- [CLI.md](./CLI.md) - Command reference
- [OPERATIONS.md](./OPERATIONS.md) - Deployment and maintenance
- [AGENTS.md](./AGENTS.md) - Agent roles and responsibilities
- [HOOKS.md](./HOOKS.md) - Claude Code integration and git worktrees
- [WORKFLOWS.md](./WORKFLOWS.md) - Molecule state machine
- [MESSAGING.md](./MESSAGING.md) - Communication patterns
- [EVENTS.md](./EVENTS.md) - Event sourcing and change feeds
- [SCHEMAS.md](./SCHEMAS.md) - Data specifications
- [VERIFICATION.md](./VERIFICATION.md) - VerMAS Inspector pipeline
- [EVALUATION.md](./EVALUATION.md) - How to evaluate the system
