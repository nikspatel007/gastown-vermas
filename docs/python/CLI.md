# VerMAS CLI Reference

> Complete command reference for co and wo tools

## Overview

VerMAS provides two CLI tools built on proven technologies:

| Tool | Purpose | Framework |
|------|---------|-----------|
| `co` | Company operations (agents, mail, sprints) | Typer (Python) |
| `wo` | Work Order operations (tasks, workflows, events) | Typer (Python) |

Both tools follow Unix philosophy:
- Small, focused commands
- Composable via pipes
- Text-based output (JSONL, plain text)
- Exit codes indicate success/failure

---

## co - Company CLI

### Status & Navigation

```bash
co status                     # Overall company status
co factories                  # List all factories
co prime                      # Load context (run on session start)
```

### Assignment Commands

```bash
co assignment                 # Check your assigned work
co assignment set <agent> <ref>     # Assign work to agent
co assignment clear <agent>         # Clear agent's assignment
co assignment attach <mail-id>      # Attach mail for handoff
```

### Mail Commands

```bash
co inbox                      # Check messages
co read <id>                  # Read specific message
co send <addr> -s "Subject" -m "Message"
co archive <id>               # Archive message
co mail list --unread         # List unread only
```

### Worker Commands

```bash
co workers [factory]          # List workers in factory
co worker spawn <factory>     # Spawn new worker
co worker done                # Signal work complete
co worker status <slot>       # Check worker status
```

### Work Dispatch

```bash
co dispatch <wo> <factory>    # Assign work order to worker
co sprints                    # Dashboard of active work
co sprint create "name" <ids> # Create sprint for batch
co sprint status <id>         # Detailed sprint progress
```

### Worktree Commands

```bash
co worktree create <factory> <slot>   # Create worker worktree
co worktree list <factory>            # List factory worktrees
co worktree remove <factory> <slot>   # Clean up worktree
```

### Factory Management

```bash
co factory add <name> <url>   # Add new factory
co factory list               # List all factories
co factory status <name>      # Factory health check
```

### Handoff

```bash
co handoff -m "Context..."    # Create handoff for next session
```

---

## wo - Work Order CLI

### Finding Work

```bash
wo ready                      # Work orders ready to work (no blockers)
wo list                       # All work orders
wo list --status=open         # Filter by status
wo list --type=bug            # Filter by type
wo list --priority=0          # Filter by priority (P0=critical)
wo show <id>                  # Detailed work order view
```

### Creating Work Orders

```bash
wo create --title="..." --type=task --priority=2
wo create --title="..." --type=bug --priority=0
wo create --title="..." --type=feature
wo create --title="..." --type=epic
```

**Work Order Types:** `task`, `bug`, `feature`, `epic`, `merge-request`, `event`, `message`

**Priority:** 0-4 (0=P0 critical, 2=P2 medium, 4=P4 backlog)

### Updating Work Orders

```bash
wo update <id> --status=in_progress
wo update <id> --assignee=username
wo update <id> --priority=1
wo update <id> --title="New title"
```

### Closing Work Orders

```bash
wo close <id>                 # Mark complete
wo close <id1> <id2> ...      # Close multiple
wo close <id> --reason="..."  # Close with reason
```

### Dependencies

```bash
wo dep add <wo> <depends-on>  # Add dependency
wo dep remove <wo> <dep>      # Remove dependency
wo blocked                    # Show blocked work orders
```

**Note:** Think "X needs Y" not "X before Y"
- `wo dep add task-B task-A` means "B depends on A" (A blocks B)

### Sync

```bash
wo sync                       # Sync with git remote
wo sync --status              # Check sync status
wo sync --force               # Force sync
```

### Statistics

```bash
wo stats                      # Project statistics
wo doctor                     # Check for issues
```

---

## wo process - Process Commands

### Workflow Management

```bash
wo process list               # List active processes
wo process show <id>          # Show process details
wo process steps <id>         # Show step status
```

### Lifecycle Operations

```bash
wo process compile <template> # Template → Ready
wo process start <wo>         # Ready → Active (attach to work order)
wo process complete <id>      # Active → Archive with summary
wo process cancel <id>        # Discard without record
```

---

## wo events - Event Commands

### Querying Events

```bash
wo events list                    # Last 50 events
wo events list --type=work_order.*  # Filter by type
wo events list --actor=ceo        # Filter by actor
wo events list --since=1h         # Time filter
```

### Real-time

```bash
wo events tail                    # Watch feed.jsonl
```

### History

```bash
wo events replay <wo-id>          # Event history for work order
```

### Statistics

```bash
wo events stats                   # Counts by type
wo events stats --since=1d        # Last 24 hours
```

### Export

```bash
wo events export --since=2026-01-01 > events.jsonl
```

---

## wo eval - Evaluation Commands

```bash
wo eval completion --since=7d     # Completion rate
wo eval assignment --since=1d     # Assignment principle compliance
wo eval throughput --since=24h    # Work throughput
wo eval verify-accuracy --since=30d
wo eval report --since=7d --output=report.json
```

---

## Common Patterns

### Starting a Session (Any Agent)

```bash
co prime                      # Load context
co assignment                 # Check for work
co inbox                      # Check messages
```

### CEO Workflow

```bash
co status                     # Company overview
wo ready                      # Find work to dispatch
co dispatch <wo> <factory>    # Assign to worker
co sprints                    # Monitor progress
```

### Worker Workflow

```bash
co assignment                 # Find assigned work (Assignment Principle)
wo show <wo-id>               # Understand the task
# ... do the work ...
git add . && git commit -m "..." && git push
co worker done                # Signal completion
```

### Supervisor Patrol

```bash
co inbox                      # Check for WORKER_DONE
co workers                    # Survey workers
# Nudge idle, kill stuck
```

### Session End

```bash
git status                    # Check changes
git add <files>               # Stage code
wo sync                       # Commit work orders
git commit -m "..."           # Commit code
git push                      # Push to remote
co handoff -m "..."           # If incomplete
```

---

## Environment Variables

| Variable | Description | Used By |
|----------|-------------|---------|
| `AGENT_ID` | Agent identity | All agents |
| `WORK_ORDER_ID` | Current work order | Workers |
| `CO_FACTORY` | Current factory | All agents |
| `CO_ROLE` | Agent role | All agents |
| `WO_DEBUG_ROUTING` | Debug prefix routing | Debugging |

### Setting Debug Mode

```bash
WO_DEBUG_ROUTING=1 wo show <id>   # Debug routing
CO_VERBOSE=1 co status            # Verbose output
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Not found |
| 4 | Permission denied |
| 5 | Conflict |

---

## Output Formats

### Default (Human-readable)

```bash
wo list
# wo-abc123  Feature X         open    P2
# wo-def456  Bug in login      open    P0
```

### JSON Output

```bash
wo list --json
# [{"id": "wo-abc123", "title": "Feature X", ...}]

wo show <id> --json
# {"id": "wo-abc123", ...}
```

### JSONL Output (for piping)

```bash
wo events list --jsonl | jq '.event_type'
```

---

## Piping and Composition

```bash
# Count open bugs
wo list --status=open --type=bug --json | jq length

# Get IDs of blocked work orders
wo blocked --json | jq -r '.[].id'

# Close all completed in sprint
co sprint status <id> --json | jq -r '.completed[].id' | xargs wo close

# Event analysis
wo events list --since=1d --jsonl | jq -r '.event_type' | sort | uniq -c
```

---

## Configuration

### Config File Location

```
~/.config/vermas/config.toml
.vermas/config.toml (per-project)
```

### Example Config

```toml
[work_orders]
default_priority = 2
sync_branch = "work-sync"

[mail]
default_from = "ceo"

[events]
retention_days = 30
partition_by_day = true
```

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [OPERATIONS.md](./OPERATIONS.md) - Deployment and operations
- [HOOKS.md](./HOOKS.md) - Claude Code integration
- [EVENTS.md](./EVENTS.md) - Event sourcing
- [SCHEMAS.md](./SCHEMAS.md) - Data specifications
- [HOW_IT_WORKS.md](./HOW_IT_WORKS.md) - Getting started
