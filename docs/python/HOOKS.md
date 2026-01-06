# VerMAS Hooks and Claude Code Integration

> How hooks, profiles, and Claude Code CLI work together

## Overview

VerMAS uses **Claude Code** as its agent runtime. Each agent is a Claude Code session running in tmux with:
- A **profile** defining its role (CLAUDE.md system prompt)
- A **hook** defining its assigned work
- An **identity** (BD_ACTOR environment variable)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      CLAUDE CODE INTEGRATION                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐                  │
│   │   Profile   │     │    Hook     │     │  BD_ACTOR   │                  │
│   │             │     │             │     │             │                  │
│   │  CLAUDE.md  │     │ .hook-agent │     │ Environment │                  │
│   │  (role def) │     │ (work ref)  │     │ (identity)  │                  │
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

VerMAS has **two distinct hook systems** that work together:

### 1. Gas Town Hooks (Work Assignment)

File-based hooks that assign work to agents:

```
.beads/.hook-{agent}     # Contains: bead:gt-abc123
```

- **What**: Simple text files mapping agents to work
- **Where**: `.beads/` directory
- **Format**: `{type}:{id}` (e.g., `bead:gt-abc123`, `mail:msg-xyz`)
- **Checked by**: Agent on startup via `gt hook`

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

## Gas Town Hook System

### Hook File Format

```
{type}:{reference_id}
```

Types:
- `bead` - Work issue to execute
- `mail` - Mail message (for handoffs)
- `mol` - Molecule (workflow instance)

Examples:
```
bead:gt-abc123
mail:msg-xyz789
mol:mol-verify-def456
```

### Hook Commands

```bash
# Check your hook
gt hook                    # Shows what's hooked (if anything)

# Set a hook (typically done by gt sling)
gt hook set <agent> bead:<id>

# Clear a hook
gt hook clear <agent>

# Attach mail as hook (for handoffs)
gt hook attach <mail-id>
```

### Hook Lifecycle

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           HOOK LIFECYCLE                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   1. SLING                                                                  │
│   ───────────────────────────────────────                                   │
│   gt sling bead-123 gastown                                                 │
│        │                                                                    │
│        ▼                                                                    │
│   Write .beads/.hook-gastown-polecats-slot0                                │
│   Content: "bead:bead-123"                                                  │
│        │                                                                    │
│        ▼                                                                    │
│   Emit event: bead.hooked                                                   │
│                                                                             │
│   2. AGENT STARTUP                                                          │
│   ───────────────────────────────────────                                   │
│   Claude Code session starts                                                │
│        │                                                                    │
│        ▼                                                                    │
│   Profile loads CLAUDE.md (with startup instructions)                       │
│        │                                                                    │
│        ▼                                                                    │
│   Agent runs: gt hook                                                       │
│        │                                                                    │
│        ▼                                                                    │
│   Emit event: hook.checked                                                  │
│                                                                             │
│   3. GUPP EXECUTION                                                         │
│   ───────────────────────────────────────                                   │
│   Hook found? → EXECUTE IMMEDIATELY                                         │
│        │                                                                    │
│        ▼                                                                    │
│   Emit event: agent.working                                                 │
│                                                                             │
│   4. COMPLETION                                                             │
│   ───────────────────────────────────────                                   │
│   Work done → gt polecat done                                               │
│        │                                                                    │
│        ▼                                                                    │
│   Hook cleared, slot released                                               │
│        │                                                                    │
│        ▼                                                                    │
│   Emit event: hook.cleared                                                  │
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
        "command": "gt prime"
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
| `UserPromptSubmit` | Before processing user input | Load context via `gt prime` |
| `PreToolUse` | Before tool execution | Verify intent before writes |
| `PostToolUse` | After tool execution | Emit tool usage events |
| `Stop` | Session ending | Emit session end event, sync beads |
| `Notification` | System notifications | Forward to mail system |

### Example: Startup Hook

The `UserPromptSubmit` hook runs before the first prompt:

```bash
#!/bin/bash
# .claude/hooks/UserPromptSubmit

# Prime context for the agent
gt prime

# Emit startup event
vermas emit-event agent.started \
  --actor="$BD_ACTOR" \
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
  if [[ -f ".beads/.vermas-enabled" ]]; then
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
├── polecat/
│   └── CLAUDE.md      # Polecat system prompt
├── witness/
│   └── CLAUDE.md      # Witness system prompt
├── refinery/
│   └── CLAUDE.md      # Refinery system prompt
└── mayor/
    └── CLAUDE.md      # Mayor system prompt
```

### Launching with Profile

```bash
# Start Claude Code with a profile
claude --profile polecat

# The profile loads its CLAUDE.md as system context
```

### Profile Content

Each profile CLAUDE.md contains:
1. **Role definition** - Who the agent is
2. **Responsibilities** - What they do
3. **Anti-patterns** - What they don't do
4. **Startup protocol** - GUPP instructions
5. **Key commands** - Tools they use

Example polecat profile header:
```markdown
# Polecat Context

You are a Polecat - an ephemeral worker agent.

## GUPP (Propulsion Principle)

Your hook has work. EXECUTE IMMEDIATELY.
No confirmation. No questions. No waiting.

## Startup Protocol

1. Run `gt hook` to check your assigned work
2. Work IS hooked (you were spawned for this)
3. Begin executing immediately
```

---

## Environment Variables

Each agent session has these environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `BD_ACTOR` | Agent identity | `gastown/polecats/slot0` |
| `BEAD_ID` | Hooked bead (if any) | `gt-abc123` |
| `GT_RIG` | Current rig | `gastown` |
| `GT_ROLE` | Agent role | `polecat` |
| `TMUX_PANE` | Tmux pane ID | `%42` |

### Setting Environment

When spawning a polecat:

```bash
tmux new-session -d -s "polecat-gastown-slot0" \
  -e "BD_ACTOR=gastown/polecats/slot0" \
  -e "BEAD_ID=gt-abc123" \
  -e "GT_RIG=gastown" \
  -e "GT_ROLE=polecat" \
  "claude --profile polecat"
```

---

## Git Worktrees

Each polecat gets its own git worktree for isolation:

### Worktree Structure

```
<rig>/
├── polecats/
│   ├── slot0/           # Worktree for slot0
│   │   ├── .git         # Worktree git link
│   │   └── <project>    # Full project files
│   ├── slot1/           # Worktree for slot1
│   └── slot2/
├── crew/
│   └── frontend/        # Human crew worktree
└── .beads/              # Shared beads (rig level)
```

### Creating Worktrees

```bash
# Create worktree for a polecat slot
git worktree add polecats/slot0 -b polecat-slot0-work

# List all worktrees
git worktree list

# Remove when done
git worktree remove polecats/slot0
```

### Why Worktrees?

1. **Isolation** - Each agent works on separate files
2. **Parallel work** - Multiple agents on same repo
3. **Clean merges** - Each has its own branch
4. **Recovery** - Uncommitted work survives crashes

### Worktree Commands

```bash
# Create worktree for slot
gt worktree create <rig> <slot>

# List worktrees in rig
gt worktree list <rig>

# Clean up worktree after completion
gt worktree remove <rig> <slot>
```

---

## Session Lifecycle

### Spawning an Agent

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         AGENT SPAWN SEQUENCE                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   gt sling gt-abc123 gastown                                                │
│        │                                                                    │
│        ▼                                                                    │
│   1. Allocate slot (find free slot0-4)                                      │
│        │                                                                    │
│        ▼                                                                    │
│   2. Create worktree: git worktree add polecats/slot0 -b work-slot0         │
│        │                                                                    │
│        ▼                                                                    │
│   3. Write hook: echo "bead:gt-abc123" > .beads/.hook-...-slot0            │
│        │                                                                    │
│        ▼                                                                    │
│   4. Create tmux session with environment                                   │
│      tmux new-session -d -s "polecat-gastown-slot0" \                       │
│        -e "BD_ACTOR=gastown/polecats/slot0" \                               │
│        -e "BEAD_ID=gt-abc123" \                                             │
│        "claude --profile polecat"                                           │
│        │                                                                    │
│        ▼                                                                    │
│   5. Claude Code starts, profile loads CLAUDE.md                            │
│        │                                                                    │
│        ▼                                                                    │
│   6. UserPromptSubmit hook runs: gt prime                                   │
│        │                                                                    │
│        ▼                                                                    │
│   7. Agent checks hook: gt hook → finds bead:gt-abc123                      │
│        │                                                                    │
│        ▼                                                                    │
│   8. GUPP: Execute immediately                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Session End

```bash
# Polecat signals completion
gt polecat done

# This:
# 1. Commits and pushes changes
# 2. Sends POLECAT_DONE to Witness
# 3. Clears hook
# 4. Releases slot
# 5. Session terminates
```

---

## Events Emitted

All hook operations emit events. See [EVENTS.md](./EVENTS.md).

| Event | When | Data |
|-------|------|------|
| `hook.set` | Hook written | agent, ref_type, ref_id |
| `hook.checked` | Agent checked hook | agent, found, response_ms |
| `hook.cleared` | Hook removed | agent, previous_ref |
| `agent.started` | Claude session began | actor, profile, worktree |
| `agent.stopped` | Claude session ended | actor, reason |

---

## Troubleshooting

### Common Issues

#### Agent Won't Start

**Symptom:** `gt sling` runs but no tmux session appears

**Check:**
```bash
# Is tmux running?
tmux list-sessions

# Is the worktree created?
ls -la <rig>/polecats/slot0/

# Check for errors in spawn
gt polecat spawn <rig> --verbose
```

**Common causes:**
- Worktree creation failed (check git status)
- Profile not found (check `~/.claude/profiles/<role>/`)
- Tmux server not running

---

#### Hook Not Found

**Symptom:** Agent starts but says "Hook is empty"

**Check:**
```bash
# Is hook file present?
ls -la .beads/.hook-*

# What's in the hook?
cat .beads/.hook-gastown-polecats-slot0

# Is BD_ACTOR set correctly?
echo $BD_ACTOR
```

**Common causes:**
- Hook file path mismatch (slashes vs dashes)
- BD_ACTOR environment not set
- Hook was cleared before agent started

---

#### Agent Stuck / Not Responding

**Symptom:** Agent session exists but no progress

**Check:**
```bash
# Attach and see what's happening
tmux attach -t polecat-gastown-slot0

# Check if Claude is waiting for input
# Look for [?] prompts

# Check events for last activity
bd events list --actor=gastown/polecats/slot0 --since=1h
```

**Common causes:**
- Claude waiting for confirmation (shouldn't happen with GUPP)
- Rate limiting
- Network issues
- Infinite loop in tool use

**Recovery:**
```bash
# Kill and respawn
tmux kill-session -t polecat-gastown-slot0
gt polecat spawn <rig> <slot>
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
ls -la .beads/events.jsonl

# Tail the feed
tail -f .beads/feed.jsonl

# Check for errors
cat .beads/events.jsonl | tail -5 | jq .
```

**Common causes:**
- File permissions
- Disk full
- Event emission disabled in config

---

#### Mail Not Delivered

**Symptom:** `gt mail send` succeeds but recipient doesn't see it

**Check:**
```bash
# Check sender's outbox
grep "from.*$SENDER" .beads/messages.jsonl

# Check recipient's inbox
grep "to.*$RECIPIENT" .beads/messages.jsonl

# Verify addresses
echo "Sender: $BD_ACTOR"
```

**Common causes:**
- Wrong recipient address (typo in path)
- Messages.jsonl not synced
- Recipient checking wrong beads location

---

### Debugging Commands

#### Hook Debugging

```bash
# View all hooks
ls -la .beads/.hook-*

# Check specific agent's hook
cat .beads/.hook-gastown-polecats-slot0

# View hook events
bd events list --type=hook.*

# Manually set a hook (for testing)
echo "bead:gt-test123" > .beads/.hook-test-agent
```

#### Session Debugging

```bash
# List all sessions
tmux list-sessions

# Attach to watch agent (read-only)
tmux attach -t polecat-gastown-slot0

# Detach: Ctrl+B then D

# Kill a stuck session
tmux kill-session -t polecat-gastown-slot0

# View session environment
tmux show-environment -t polecat-gastown-slot0
```

#### Worktree Debugging

```bash
# List all worktrees
git worktree list

# Check worktree status
cd <rig>/polecats/slot0 && git status

# Prune stale worktrees
git worktree prune

# Remove specific worktree
git worktree remove <path>
```

#### Event Debugging

```bash
# Tail events in real-time
tail -f .beads/feed.jsonl | jq .

# Filter by type
grep '"event_type":"bead' .beads/events.jsonl | tail -10 | jq .

# Filter by actor
grep '"actor":"mayor"' .beads/events.jsonl | jq .

# Count events by type
cat .beads/events.jsonl | jq -r '.event_type' | sort | uniq -c
```

#### Log Debugging

```bash
# Session logs
ls logs/polecat-gastown-slot0/

# Tail Claude output
tail -f logs/polecat-gastown-slot0/claude.log

# Search for errors
grep -i error logs/polecat-gastown-slot0/*.log
```

---

### Recovery Procedures

#### Full Agent Recovery

When an agent is stuck beyond repair:

```bash
# 1. Kill the session
tmux kill-session -t polecat-gastown-slot0

# 2. Clear the hook
rm .beads/.hook-gastown-polecats-slot0

# 3. Prune the worktree
git worktree remove polecats/slot0 --force 2>/dev/null
git worktree prune

# 4. Re-sling the work
gt sling <bead-id> <rig>
```

#### Beads Recovery

If `.beads/` gets corrupted:

```bash
# 1. Check git status
cd .beads && git status

# 2. Reset to last known good state
git checkout HEAD -- issues.jsonl messages.jsonl

# 3. Re-sync
bd sync --force
```

#### Event Log Recovery

If events.jsonl is corrupted:

```bash
# 1. Validate JSON lines
cat .beads/events.jsonl | while read line; do
  echo "$line" | jq . > /dev/null || echo "Bad line: $line"
done

# 2. Filter to valid lines
cat .beads/events.jsonl | jq -c . > events.jsonl.clean
mv events.jsonl.clean events.jsonl
```

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [OPERATIONS.md](./OPERATIONS.md) - Deployment and recovery procedures
- [AGENTS.md](./AGENTS.md) - Agent roles
- [EVENTS.md](./EVENTS.md) - Event sourcing
- [MESSAGING.md](./MESSAGING.md) - Mail system
- [WORKFLOWS.md](./WORKFLOWS.md) - Molecule workflows
- [SCHEMAS.md](./SCHEMAS.md) - Data specifications
- [CLI.md](./CLI.md) - Command reference
