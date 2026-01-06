# VerMAS Hooks and Claude Code Integration

> How assignments, profiles, and Claude Code CLI work together

## Overview

VerMAS uses **Claude Code** as its agent runtime. Each agent is a Claude Code session running in tmux with:
- A **profile** defining its role (CLAUDE.md system prompt)
- An **assignment** defining its assigned work
- An **identity** (AGENT_ID environment variable)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      CLAUDE CODE INTEGRATION                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐                  │
│   │   Profile   │     │ Assignment  │     │  AGENT_ID   │                  │
│   │             │     │             │     │             │                  │
│   │  CLAUDE.md  │     │.assignment- │     │ Environment │                  │
│   │  (role def) │     │   agent     │     │ (identity)  │                  │
│   └──────┬──────┘     └──────┬──────┘     └──────┬──────┘                  │
│          │                   │                   │                          │
│          └───────────────────┼───────────────────┘                          │
│                              │                                              │
│                              ▼                                              │
│                    ┌─────────────────┐                                      │
│                    │  Claude Code    │                                      │
│                    │    Session      │                                      │
│                    │  (tmux pane)    │                                      │
│                    └─────────────────┘                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Two Types of Hooks

VerMAS has **two distinct systems** that work together:

### 1. Assignment Files (Work Assignment)

File-based assignments that give work to agents:

```
.work/.assignment-{agent}     # Contains: work_order:wo-abc123
```

- **What**: Simple text files mapping agents to work
- **Where**: `.work/` directory
- **Format**: `{type}:{id}` (e.g., `work_order:wo-abc123`, `mail:msg-xyz`)
- **Checked by**: Agent on startup via `co assignment`

### 2. Claude Code Hooks (Lifecycle Events)

Claude Code's built-in hook system for executing code at lifecycle points:

```
.claude/hooks/
├── PreToolUse           # Before tool execution
├── PostToolUse          # After tool execution
├── Notification         # On notifications
├── Stop                 # On session stop
└── UserPromptSubmit     # On user input
```

- **What**: Shell scripts/commands executed by Claude Code
- **Where**: `.claude/hooks/` or `~/.claude/hooks/`
- **Format**: Executable scripts or JSONL config
- **Triggered by**: Claude Code runtime events

---

## Assignment System

### Assignment File Format

```
{type}:{reference_id}
```

Types:
- `work_order` - Work order to execute
- `mail` - Mail message (for handoffs)
- `process` - Process (workflow instance)

Examples:
```
work_order:wo-abc123
mail:msg-xyz789
process:proc-verify-def456
```

### Assignment Commands

```bash
# Check your assignment
co assignment                    # Shows what's assigned (if anything)

# Set an assignment (typically done by co dispatch)
co assignment set <agent> work_order:<id>

# Clear an assignment
co assignment clear <agent>

# Attach mail as assignment (for handoffs)
co assignment attach <mail-id>
```

### Assignment Lifecycle

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       ASSIGNMENT LIFECYCLE                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   1. DISPATCH                                                               │
│   ───────────────────────────────────────                                   │
│   co dispatch wo-123 project-a                                              │
│        │                                                                    │
│        ▼                                                                    │
│   Write .work/.assignment-project-a-workers-slot0                          │
│   Content: "work_order:wo-123"                                              │
│        │                                                                    │
│        ▼                                                                    │
│   Emit event: work_order.assigned                                           │
│                                                                             │
│   2. AGENT STARTUP                                                          │
│   ───────────────────────────────────────                                   │
│   Claude Code session starts                                                │
│        │                                                                    │
│        ▼                                                                    │
│   Profile loads CLAUDE.md (with startup instructions)                       │
│        │                                                                    │
│        ▼                                                                    │
│   Agent runs: co assignment                                                 │
│        │                                                                    │
│        ▼                                                                    │
│   Emit event: agent.assignment_checked                                      │
│                                                                             │
│   3. ASSIGNMENT PRINCIPLE                                                   │
│   ───────────────────────────────────────                                   │
│   Assignment found? → EXECUTE IMMEDIATELY                                   │
│        │                                                                    │
│        ▼                                                                    │
│   Emit event: agent.working                                                 │
│                                                                             │
│   4. COMPLETION                                                             │
│   ───────────────────────────────────────                                   │
│   Work done → co worker done                                                │
│        │                                                                    │
│        ▼                                                                    │
│   Assignment cleared, slot released                                         │
│        │                                                                    │
│        ▼                                                                    │
│   Emit event: assignment.cleared                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Claude Code Hooks

Claude Code provides lifecycle hooks that VerMAS uses for:
- Injecting context on startup
- Emitting events on tool use
- Enforcing verification

### Hook Configuration

Located in `.claude/settings.json` or `.claude/hooks/`:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Write|Edit",
        "command": "vermas verify-intent"
      }
    ],
    "PostToolUse": [
      {
        "matcher": "*",
        "command": "vermas emit-event tool.used"
      }
    ],
    "UserPromptSubmit": [
      {
        "command": "co prime"
      }
    ],
    "Stop": [
      {
        "command": "vermas session-end"
      }
    ]
  }
}
```

### Common Claude Code Hooks

| Hook | Trigger | VerMAS Usage |
|------|---------|--------------|
| `UserPromptSubmit` | Before processing user input | Load context via `co prime` |
| `PreToolUse` | Before tool execution | Verify intent before writes |
| `PostToolUse` | After tool execution | Emit tool usage events |
| `Stop` | Session ending | Emit session end event, sync work orders |
| `Notification` | System notifications | Forward to mail system |

### Example: Startup Hook

The `UserPromptSubmit` hook runs before the first prompt:

```bash
#!/bin/bash
# .claude/hooks/UserPromptSubmit

# Prime context for the agent
co prime

# Emit startup event
vermas emit-event agent.started \
  --actor="$AGENT_ID" \
  --data="{\"session\": \"$TMUX_PANE\"}"
```

### Example: Verification Hook

Enforce verification before code changes:

```bash
#!/bin/bash
# .claude/hooks/PreToolUse

TOOL="$1"

if [[ "$TOOL" == "Write" || "$TOOL" == "Edit" ]]; then
  # Check if verification is enabled
  if [[ -f ".work/.vermas-enabled" ]]; then
    vermas check-intent "$@"
    exit $?
  fi
fi

exit 0
```

---

## Claude Code Profiles

Profiles define agent roles via CLAUDE.md files:

### Profile Structure

```
~/.claude/profiles/
├── worker/
│   └── CLAUDE.md      # Worker system prompt
├── supervisor/
│   └── CLAUDE.md      # Supervisor system prompt
├── qa/
│   └── CLAUDE.md      # QA system prompt
└── ceo/
    └── CLAUDE.md      # CEO system prompt
```

### Launching with Profile

```bash
# Start Claude Code with a profile
claude --profile worker

# The profile loads its CLAUDE.md as system context
```

### Profile Content

Each profile CLAUDE.md contains:
1. **Role definition** - Who the agent is
2. **Responsibilities** - What they do
3. **Anti-patterns** - What they don't do
4. **Startup protocol** - Assignment Principle instructions
5. **Key commands** - Tools they use

Example worker profile header:
```markdown
# Worker Context

You are a Worker - an ephemeral task executor.

## Assignment Principle

Your assignment has work. EXECUTE IMMEDIATELY.
No confirmation. No questions. No waiting.

## Startup Protocol

1. Run `co assignment` to check your assigned work
2. Work IS assigned (you were spawned for this)
3. Begin executing immediately
```

---

## Environment Variables

Each agent session has these environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `AGENT_ID` | Agent identity | `project-a/workers/slot0` |
| `WORK_ORDER_ID` | Assigned work order (if any) | `wo-abc123` |
| `CO_FACTORY` | Current factory | `project-a` |
| `CO_ROLE` | Agent role | `worker` |
| `TMUX_PANE` | Tmux pane ID | `%42` |

### Setting Environment

When spawning a worker:

```bash
tmux new-session -d -s "worker-project-a-slot0" \
  -e "AGENT_ID=project-a/workers/slot0" \
  -e "WORK_ORDER_ID=wo-abc123" \
  -e "CO_FACTORY=project-a" \
  -e "CO_ROLE=worker" \
  "claude --profile worker"
```

---

## Git Worktrees

Each worker gets its own git worktree for isolation:

### Worktree Structure

```
<factory>/
├── workers/
│   ├── slot0/           # Worktree for slot0
│   │   ├── .git         # Worktree git link
│   │   └── <project>    # Full project files
│   ├── slot1/           # Worktree for slot1
│   └── slot2/
├── teams/
│   └── frontend/        # Human team worktree
└── .work/               # Shared work orders (factory level)
```

### Creating Worktrees

```bash
# Create worktree for a worker slot
git worktree add workers/slot0 -b worker-slot0-work

# List all worktrees
git worktree list

# Remove when done
git worktree remove workers/slot0
```

### Why Worktrees?

1. **Isolation** - Each agent works on separate files
2. **Parallel work** - Multiple agents on same repo
3. **Clean merges** - Each has its own branch
4. **Recovery** - Uncommitted work survives crashes

### Worktree Commands

```bash
# Create worktree for slot
co worktree create <factory> <slot>

# List worktrees in factory
co worktree list <factory>

# Clean up worktree after completion
co worktree remove <factory> <slot>
```

---

## Session Lifecycle

### Spawning an Agent

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         AGENT SPAWN SEQUENCE                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   co dispatch wo-abc123 project-a                                           │
│        │                                                                    │
│        ▼                                                                    │
│   1. Allocate slot (find free slot0-4)                                      │
│        │                                                                    │
│        ▼                                                                    │
│   2. Create worktree: git worktree add workers/slot0 -b work-slot0          │
│        │                                                                    │
│        ▼                                                                    │
│   3. Write assignment: echo "work_order:wo-abc123" > .work/.assignment-...  │
│        │                                                                    │
│        ▼                                                                    │
│   4. Create tmux session with environment                                   │
│      tmux new-session -d -s "worker-project-a-slot0" \                      │
│        -e "AGENT_ID=project-a/workers/slot0" \                              │
│        -e "WORK_ORDER_ID=wo-abc123" \                                       │
│        "claude --profile worker"                                            │
│        │                                                                    │
│        ▼                                                                    │
│   5. Claude Code starts, profile loads CLAUDE.md                            │
│        │                                                                    │
│        ▼                                                                    │
│   6. UserPromptSubmit hook runs: co prime                                   │
│        │                                                                    │
│        ▼                                                                    │
│   7. Agent checks assignment: co assignment → finds work_order:wo-abc123    │
│        │                                                                    │
│        ▼                                                                    │
│   8. Assignment Principle: Execute immediately                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Session End

```bash
# Worker signals completion
co worker done

# This:
# 1. Commits and pushes changes
# 2. Sends WORKER_DONE to Supervisor
# 3. Clears assignment
# 4. Releases slot
# 5. Session terminates
```

---

## Events Emitted

All assignment operations emit events. See [EVENTS.md](./EVENTS.md).

| Event | When | Data |
|-------|------|------|
| `assignment.set` | Assignment written | agent, ref_type, ref_id |
| `assignment.checked` | Agent checked assignment | agent, found, response_ms |
| `assignment.cleared` | Assignment removed | agent, previous_ref |
| `agent.started` | Claude session began | actor, profile, worktree |
| `agent.stopped` | Claude session ended | actor, reason |

---

## Troubleshooting

### Common Issues

#### Agent Won't Start

**Symptom:** `co dispatch` runs but no tmux session appears

**Check:**
```bash
# Is tmux running?
tmux list-sessions

# Is the worktree created?
ls -la <factory>/workers/slot0/

# Check for errors in spawn
co worker spawn <factory> --verbose
```

**Common causes:**
- Worktree creation failed (check git status)
- Profile not found (check `~/.claude/profiles/<role>/`)
- Tmux server not running

---

#### Assignment Not Found

**Symptom:** Agent starts but says "Assignment is empty"

**Check:**
```bash
# Is assignment file present?
ls -la .work/.assignment-*

# What's in the assignment?
cat .work/.assignment-project-a-workers-slot0

# Is AGENT_ID set correctly?
echo $AGENT_ID
```

**Common causes:**
- Assignment file path mismatch (slashes vs dashes)
- AGENT_ID environment not set
- Assignment was cleared before agent started

---

#### Agent Stuck / Not Responding

**Symptom:** Agent session exists but no progress

**Check:**
```bash
# Attach and see what's happening
tmux attach -t worker-project-a-slot0

# Check if Claude is waiting for input
# Look for [?] prompts

# Check events for last activity
wo events list --actor=project-a/workers/slot0 --since=1h
```

**Common causes:**
- Claude waiting for confirmation (shouldn't happen with Assignment Principle)
- Rate limiting
- Network issues
- Infinite loop in tool use

**Recovery:**
```bash
# Kill and respawn
tmux kill-session -t worker-project-a-slot0
co worker spawn <factory> <slot>
```

---

#### Worktree Conflicts

**Symptom:** "fatal: '<path>' is already checked out"

**Check:**
```bash
# List all worktrees
git worktree list

# Check for stale worktrees
git worktree prune --dry-run
```

**Fix:**
```bash
# Remove stale worktree reference
git worktree prune

# Or force remove
git worktree remove <path> --force
```

---

#### Events Not Appearing

**Symptom:** Actions happen but no events in `events.jsonl`

**Check:**
```bash
# Is events.jsonl writable?
ls -la .work/events.jsonl

# Tail the feed
tail -f .work/feed.jsonl

# Check for errors
cat .work/events.jsonl | tail -5 | jq .
```

**Common causes:**
- File permissions
- Disk full
- Event emission disabled in config

---

#### Mail Not Delivered

**Symptom:** `co send` succeeds but recipient doesn't see it

**Check:**
```bash
# Check sender's outbox
grep "from.*$SENDER" .work/messages.jsonl

# Check recipient's inbox
grep "to.*$RECIPIENT" .work/messages.jsonl

# Verify addresses
echo "Sender: $AGENT_ID"
```

**Common causes:**
- Wrong recipient address (typo in path)
- Messages.jsonl not synced
- Recipient checking wrong work location

---

### Debugging Commands

#### Assignment Debugging

```bash
# View all assignments
ls -la .work/.assignment-*

# Check specific agent's assignment
cat .work/.assignment-project-a-workers-slot0

# View assignment events
wo events list --type=assignment.*

# Manually set an assignment (for testing)
echo "work_order:wo-test123" > .work/.assignment-test-agent
```

#### Session Debugging

```bash
# List all sessions
tmux list-sessions

# Attach to watch agent (read-only)
tmux attach -t worker-project-a-slot0

# Detach: Ctrl+B then D

# Kill a stuck session
tmux kill-session -t worker-project-a-slot0

# View session environment
tmux show-environment -t worker-project-a-slot0
```

#### Worktree Debugging

```bash
# List all worktrees
git worktree list

# Check worktree status
cd <factory>/workers/slot0 && git status

# Prune stale worktrees
git worktree prune

# Remove specific worktree
git worktree remove <path>
```

#### Event Debugging

```bash
# Tail events in real-time
tail -f .work/feed.jsonl | jq .

# Filter by type
grep '"event_type":"work_order' .work/events.jsonl | tail -10 | jq .

# Filter by actor
grep '"actor":"ceo"' .work/events.jsonl | jq .

# Count events by type
cat .work/events.jsonl | jq -r '.event_type' | sort | uniq -c
```

#### Log Debugging

```bash
# Session logs
ls logs/worker-project-a-slot0/

# Tail Claude output
tail -f logs/worker-project-a-slot0/claude.log

# Search for errors
grep -i error logs/worker-project-a-slot0/*.log
```

---

### Recovery Procedures

#### Full Agent Recovery

When an agent is stuck beyond repair:

```bash
# 1. Kill the session
tmux kill-session -t worker-project-a-slot0

# 2. Clear the assignment
rm .work/.assignment-project-a-workers-slot0

# 3. Prune the worktree
git worktree remove workers/slot0 --force 2>/dev/null
git worktree prune

# 4. Re-dispatch the work
co dispatch <wo-id> <factory>
```

#### Work Order Recovery

If `.work/` gets corrupted:

```bash
# 1. Check git status
cd .work && git status

# 2. Reset to last known good state
git checkout HEAD -- work_orders.jsonl messages.jsonl

# 3. Re-sync
wo sync --force
```

#### Event Log Recovery

If events.jsonl is corrupted:

```bash
# 1. Validate JSON lines
cat .work/events.jsonl | while read line; do
  echo "$line" | jq . > /dev/null || echo "Bad line: $line"
done

# 2. Filter to valid lines
cat .work/events.jsonl | jq -c . > events.jsonl.clean
mv events.jsonl.clean events.jsonl
```

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [OPERATIONS.md](./OPERATIONS.md) - Deployment and recovery procedures
- [AGENTS.md](./AGENTS.md) - Agent roles
- [EVENTS.md](./EVENTS.md) - Event sourcing
- [MESSAGING.md](./MESSAGING.md) - Mail system
- [WORKFLOWS.md](./WORKFLOWS.md) - Process workflows
- [SCHEMAS.md](./SCHEMAS.md) - Data specifications
- [CLI.md](./CLI.md) - Command reference
