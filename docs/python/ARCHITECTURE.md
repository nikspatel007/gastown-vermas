# VerMAS Python Architecture

> System design for CLI-based multi-agent verification

**See also:** [INDEX.md](./INDEX.md) for documentation map

## Design Principles

1. **No API costs** - All LLM interactions through Claude Code CLI
2. **Shared data formats** - Interoperable with Go implementation via JSONL/TOML
3. **Tmux isolation** - Each agent runs in its own tmux session
4. **Git-backed state** - All persistent state lives in git
5. **Assignment-driven execution** - "If you have an assignment, EXECUTE IT"
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
# .work/config.toml

[llm]
# Default backend for all agents
backend = "claude"
command = "claude --profile {role}"

# Per-role overrides
[llm.roles.worker]
backend = "claude"
command = "claude --profile worker"

[llm.roles.verifier]
# Verifier uses no LLM - just shell execution
backend = "none"

[llm.roles.auditor]
# Use different model for verification
backend = "claude"
command = "claude --profile auditor --model opus"
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
│                              COMPANY (Workspace Root)                        │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                         COORDINATION LAYER                          │  │
│   │                                                                     │  │
│   │   CEO ◄────────────────────────────────────────────► Operations    │  │
│   │   (Human-directed)                                   (Daemon)      │  │
│   │                                                                     │  │
│   │   - Cross-factory decisions                         - Health mon   │  │
│   │   - Strategic planning                              - Restarts     │  │
│   │   - Escalation handling                             - Watchdog     │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                    │                                        │
│                    ┌───────────────┴───────────────┐                       │
│                    ▼                               ▼                        │
│   ┌────────────────────────────┐   ┌────────────────────────────┐         │
│   │       FACTORY A            │   │       FACTORY B            │         │
│   │                            │   │                            │         │
│   │   ┌──────────┐ ┌───────┐  │   │   ┌──────────┐ ┌───────┐  │         │
│   │   │Supervisor│ │  QA   │  │   │   │Supervisor│ │  QA   │  │         │
│   │   └────┬─────┘ └───┬───┘  │   │   └────┬─────┘ └───┬───┘  │         │
│   │        │           │      │   │        │           │      │         │
│   │        ▼           ▼      │   │        ▼           ▼      │         │
│   │   ┌─────────────────┐     │   │   ┌─────────────────┐     │         │
│   │   │     Workers     │     │   │   │     Workers     │     │         │
│   │   │  slot0..slot4   │     │   │   │  slot0..slot4   │     │         │
│   │   └─────────────────┘     │   │   └─────────────────┘     │         │
│   │                            │   │                            │         │
│   │   .work/ (factory-level)  │   │   .work/ (factory-level)  │         │
│   └────────────────────────────┘   └────────────────────────────┘         │
│                                                                             │
│   .work/ (company-level: CEO mail, HQ coordination)                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Data Flow

### Work Assignment (Dispatch)

```
CEO                     Factory                  Worker
  │                       │                        │
  │  1. co dispatch wo factory                     │
  │──────────────────────▶│                        │
  │                       │                        │
  │                       │  2. Allocate slot      │
  │                       │  3. Create worktree    │
  │                       │  4. Write assignment   │
  │                       │  5. Start tmux session │
  │                       │───────────────────────▶│
  │                       │                        │
  │                       │                        │  6. Check assignment
  │                       │                        │  7. Find work
  │                       │                        │  8. EXECUTE
  │                       │                        │
```

### Work Completion

```
Worker                 Supervisor               QA
  │                        │                       │
  │  1. co worker done     │                       │
  │  (sends WORKER_DONE)   │                       │
  │───────────────────────▶│                       │
  │                        │                       │
  │                        │  2. Validate work     │
  │                        │  3. Send READY_FOR_QA │
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

### Two-Level Work Orders

| Level | Location | Purpose | Git Behavior |
|-------|----------|---------|--------------|
| **Company** | `~/.work/` | CEO mail, HQ coordination | Commits to main |
| **Factory** | `<factory>/.work/` | Project work orders, workflows | Uses work-sync branch |

### File Types

| File | Format | Contents |
|------|--------|----------|
| `events.jsonl` | JSONL | **Event log (source of truth)** |
| `work_orders.jsonl` | JSONL | Work orders - projection of events |
| `messages.jsonl` | JSONL | Mail - projection of events |
| `feed.jsonl` | JSONL | Real-time change feed |
| `routes.jsonl` | JSONL | Prefix → factory routing |
| `templates/*.toml` | TOML | Workflow templates |
| `processes/*.json` | JSON | Active workflow instances |
| `.assignment-{agent}` | Plain text | Current assignment |

### Event Sourcing Model

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          EVENT SOURCING MODEL                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Commands (co, wo)                                                         │
│        │                                                                    │
│        ▼                                                                    │
│   ┌─────────────┐                                                           │
│   │   events.   │  ← Source of truth (append-only)                         │
│   │   jsonl     │                                                           │
│   └──────┬──────┘                                                           │
│          │                                                                  │
│          │ project                                                          │
│          ▼                                                                  │
│   ┌──────────────┐   ┌─────────────┐    ┌─────────────┐                    │
│   │ work_orders. │   │ messages.   │    │  feed.      │                    │
│   │ jsonl        │   │ jsonl       │    │  jsonl      │                    │
│   │              │   │             │    │             │                    │
│   │ (current     │   │ (mailbox    │    │ (real-time  │                    │
│   │  state)      │   │  state)     │    │  stream)    │                    │
│   └──────────────┘   └─────────────┘    └─────────────┘                    │
│                                                                             │
│   All state is derived from events. Events are immutable.                  │
│   See EVENTS.md for full event sourcing documentation.                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Prefix-Based Routing

Every work order ID has a prefix (e.g., `wo-abc12`). The router maps prefixes to factories:

```
wo-*  → default-factory/.work/
hq-*  → company/.work/ (CEO level)
pa-*  → project-a/.work/
```

---

## Session Architecture

### Tmux Session Naming

```
{role}-{factory}           # Supervisor, QA
{role}-{factory}-{slot}    # Workers
ceo                        # Company-level CEO
operations                 # Company-level Operations
```

Examples:
- `supervisor-project-a`
- `qa-project-a`
- `worker-project-a-slot0`
- `ceo`

### Session Environment

Each session sets:
- `AGENT_ID` - Agent identity (e.g., `project-a/workers/slot0`)
- `WORK_ORDER_ID` - Assigned work (for workers)
- Working directory - Appropriate worktree or factory path

### Claude Code Profile

Each role has a profile that loads its CLAUDE.md:
- `claude --profile supervisor`
- `claude --profile qa`
- `claude --profile worker`
- `claude --profile ceo`

---

## Verification Integration (VerMAS)

### Where Verification Happens

```
Worker completes work
        │
        ▼
    Supervisor
        │
        ▼
    QA Department ──────────────────────────────────────┐
        │                                               │
        ▼                                               ▼
    Run Tests                              Run VerMAS Verification
        │                                               │
        │                                               ▼
        │                                      ┌─────────────────┐
        │                                      │    Designer     │
        │                                      │  (elaborate)    │
        │                                      └────────┬────────┘
        │                                               │
        │                                      ┌────────▼────────┐
        │                                      │   Strategist    │
        │                                      │  (plan tests)   │
        │                                      └────────┬────────┘
        │                                               │
        │                                      ┌────────▼────────┐
        │                                      │    Verifier     │
        │                                      │ (run shell/no LLM)│
        │                                      └────────┬────────┘
        │                                               │
        │                                      ┌────────▼────────┐
        │                                      │    Auditor      │
        │                                      │ (LLM if needed) │
        │                                      └────────┬────────┘
        │                                               │
        │                                      ┌────────▼────────┐
        │                                      │   Adversarial   │
        │                                      │ Advocate/Critic │
        │                                      │     Judge       │
        │                                      └────────┬────────┘
        │                                               │
        ▼                                               ▼
    All pass? ◄─────────────────────────────── Verdict: PASS/FAIL
        │
        ├── Yes → Merge
        └── No  → REWORK_REQUEST
```

### QA Pipeline Sessions

Each QA role can run as its own tmux session:
- `qa-designer-{factory}`
- `qa-strategist-{factory}`
- `qa-verifier-{factory}` (no LLM - just runs shell)
- `qa-advocate-{factory}`
- `qa-critic-{factory}`
- `qa-judge-{factory}`

Or as a single orchestrated workflow within QA Department.

---

## Logging Strategy

### What Gets Logged

| Log Type | Location | Contents |
|----------|----------|----------|
| **Session logs** | `logs/{session}/` | Full Claude Code output |
| **Mail archive** | `.work/messages.jsonl` | All agent communication |
| **Work order history** | `.work/work_orders.jsonl` | Work state changes |
| **Process traces** | `.work/processes/*.json` | Workflow step execution |
| **Verification evidence** | `.work/evidence/` | Test outputs, verdicts |

### Log Levels

1. **Trace** - Every Claude Code interaction (large, debugging only)
2. **Debug** - Step-by-step workflow execution
3. **Info** - Major state changes (work order status, merges)
4. **Warn** - Nudges, retries, recoverable issues
5. **Error** - Failures, escalations, stuck agents

### Structured Logging Fields

- `timestamp` - When it happened
- `actor` - AGENT_ID of the agent
- `event` - What happened (work_order_created, mail_sent, step_completed)
- `work_order_id` - Related work order if any
- `process_id` - Related process if any
- `details` - Event-specific data

---

## Failure Modes and Recovery

### Agent Failures

| Failure | Detection | Recovery |
|---------|-----------|----------|
| Worker stuck | Supervisor patrol (idle >15min) | Kill session, release slot |
| Supervisor down | Operations patrol | Restart Supervisor |
| QA down | Operations patrol | Restart QA |
| Operations down | Boot process | Restart Operations |

### Work Recovery

| Scenario | State | Recovery |
|----------|-------|----------|
| Worker killed mid-work | Worktree has uncommitted changes | New worker can resume from worktree |
| Session crash | Assignment file persists | New session reads assignment, continues |
| Merge failed | Work order still open | REWORK_REQUEST sent, work continues |

### Watchdog Chain

```
OS (systemd/launchd)
        │
        ▼
     Boot
        │
        ▼
  Operations ──────────────────┬──────────────────┐
        │                      │                  │
        ▼                      ▼                  ▼
 Supervisor(A)          Supervisor(B)      Supervisor(C)
        │                      │                  │
        ▼                      ▼                  ▼
    Workers                Workers             Workers
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

### Why Assignments?

1. **Stateless startup** - Agent reads assignment, knows what to do
2. **Crash recovery** - Assignment persists, work doesn't get lost
3. **Explicit assignment** - No ambiguity about who does what
4. **Principle enforcement** - Assignment = immediate execution

### Why Processes?

1. **Reusable workflows** - Define once, instantiate many
2. **Step tracking** - Know exactly where work stopped
3. **Dependencies** - Steps can depend on other steps
4. **Audit trail** - Archive completed workflows

---

## Security Considerations

### Threat Model

VerMAS runs on a single machine or trusted cluster. The threat model assumes:
- **Trusted agents** - All LLM agents are under your control
- **Trusted filesystem** - No malicious writes to `.work/`
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
│   │   Agents (tmux)     Work Orders (files)   Git (worktrees)            │ │
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
# Work directory should be readable/writable by your user only
chmod 700 .work
chmod 600 .work/*.jsonl
```

**Secrets Management:**
- Never store API keys in work orders or git
- Use environment variables or secure vaults
- LLM CLIs handle their own auth

**Git Security:**
```bash
# Add to .gitignore
.work/.assignment-*     # Assignment files (local state)
.work/feed.jsonl        # Real-time feed (local)
logs/                   # Session logs (may contain sensitive output)
```

**Agent Isolation:**
- Each worker runs in isolated worktree
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
- Assignment-based work (explicit scope)

### Verification as Security

The VerMAS verification pipeline provides security benefits:
- Code review before merge (QA roles)
- Automated testing (Verifier)
- Adversarial review (Advocate/Critic)
- Audit trail of all decisions

---

## Scaling and Performance

### Scaling Dimensions

| Dimension | Typical Range | Bottleneck |
|-----------|---------------|------------|
| **Agents per factory** | 5-10 workers | Tmux sessions, disk I/O |
| **Factories per company** | 3-10 | Memory, coordination overhead |
| **Total agents** | 20-50 | LLM rate limits, human oversight |
| **Work orders per factory** | 1000s | JSONL scan time |

### Performance Characteristics

**Fast Operations (< 100ms):**
- Assignment check (`co assignment`)
- Mail send (`co send`)
- Event emit (append to JSONL)

**Medium Operations (100ms - 1s):**
- Work order lookup by ID (`wo show`)
- List work orders with filter (`wo list`)
- Sync status check (`wo sync --status`)

**Slow Operations (> 1s):**
- Full work orders sync (`wo sync`)
- Worktree creation (`git worktree add`)
- LLM agent spawn (Claude startup)

### Optimizations

**Event Log Partitioning:**
```
.work/
├── events.jsonl           # Current day
└── events/
    ├── 2026-01-06.jsonl   # Archived by day
    └── 2026-01-05.jsonl
```

**Index Files (optional):**
```python
# For large work order counts, maintain index
# .work/index/by-status.json
{
  "open": ["wo-abc", "wo-def", ...],
  "closed": ["wo-xyz", ...]
}
```

**Projection Caching:**
- `work_orders.jsonl` is a projection, not source of truth
- Can be regenerated from events
- Cache for fast reads, events for writes

### Scaling Recommendations

**< 10 agents:** Default configuration works fine

**10-30 agents:**
- Partition events by day
- Use separate factories for independent projects
- Monitor disk I/O

**30+ agents:**
- Consider multiple machines
- Implement event archival
- Add monitoring (Operations metrics)

### Resource Usage

| Component | Memory | Disk | CPU |
|-----------|--------|------|-----|
| Tmux session | ~5MB | - | Idle |
| Claude Code | ~200MB | - | Varies |
| Work order sync | ~50MB | ~1MB/1000 WOs | Low |
| Event tail | ~10MB | Append only | Low |

---

## See Also

- [INDEX.md](./INDEX.md) - Documentation map and glossary
- [HOW_IT_WORKS.md](./HOW_IT_WORKS.md) - Quick start guide
- [CLI.md](./CLI.md) - Command reference
- [OPERATIONS.md](./OPERATIONS.md) - Deployment and maintenance
- [AGENTS.md](./AGENTS.md) - Agent roles and responsibilities
- [HOOKS.md](./HOOKS.md) - Claude Code integration and git worktrees
- [WORKFLOWS.md](./WORKFLOWS.md) - Process state machine
- [MESSAGING.md](./MESSAGING.md) - Communication patterns
- [EVENTS.md](./EVENTS.md) - Event sourcing and change feeds
- [SCHEMAS.md](./SCHEMAS.md) - Data specifications
- [VERIFICATION.md](./VERIFICATION.md) - VerMAS QA pipeline
- [EVALUATION.md](./EVALUATION.md) - How to evaluate the system
