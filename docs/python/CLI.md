# VerMAS CLI Reference

> Complete command reference for gt and bd tools

## Overview

VerMAS provides two CLI tools built on proven technologies:

| Tool | Purpose | Framework |
|------|---------|-----------|
| `gt` | Gas Town operations (agents, mail, convoys) | Typer (Python) |
| `bd` | Beads operations (issues, workflows, events) | Typer (Python) |

Both tools follow Unix philosophy:
- Small, focused commands
- Composable via pipes
- Text-based output (JSONL, plain text)
- Exit codes indicate success/failure

---

## gt - Gas Town CLI

### Status & Navigation

```bash
gt status                     # Overall town status
gt rigs                       # List all rigs
gt prime                      # Load context (run on session start)
```

### Hook Commands

```bash
gt hook                       # Check your assigned work
gt hook set <agent> <ref>     # Assign work to agent
gt hook clear <agent>         # Clear agent's hook
gt hook attach <mail-id>      # Hook mail for handoff
```

### Mail Commands

```bash
gt mail inbox                 # Check messages
gt mail read <id>             # Read specific message
gt mail send <addr> -s "Subject" -m "Message"
gt mail archive <id>          # Archive message
gt mail list --unread         # List unread only
```

### Polecat Commands

```bash
gt polecat list [rig]         # List polecats in rig
gt polecat spawn <rig>        # Spawn new polecat
gt polecat done               # Signal work complete
gt polecat status <slot>      # Check polecat status
```

### Work Dispatch

```bash
gt sling <bead> <rig>         # Assign bead to polecat
gt convoy list                # Dashboard of active work
gt convoy create "name" <ids> # Create convoy for batch
gt convoy status <id>         # Detailed convoy progress
```

### Worktree Commands

```bash
gt worktree create <rig> <slot>   # Create polecat worktree
gt worktree list <rig>            # List rig worktrees
gt worktree remove <rig> <slot>   # Clean up worktree
```

### Rig Management

```bash
gt rig add <name> <url>       # Add new rig
gt rig list                   # List all rigs
gt rig status <name>          # Rig health check
```

### Handoff

```bash
gt handoff -m "Context..."    # Create handoff for next session
```

---

## bd - Beads CLI

### Finding Work

```bash
bd ready                      # Issues ready to work (no blockers)
bd list                       # All issues
bd list --status=open         # Filter by status
bd list --type=bug            # Filter by type
bd list --priority=0          # Filter by priority (P0=critical)
bd show <id>                  # Detailed issue view
```

### Creating Issues

```bash
bd create --title="..." --type=task --priority=2
bd create --title="..." --type=bug --priority=0
bd create --title="..." --type=feature
bd create --title="..." --type=epic
```

**Issue Types:** `task`, `bug`, `feature`, `epic`, `merge-request`, `event`, `message`

**Priority:** 0-4 (0=P0 critical, 2=P2 medium, 4=P4 backlog)

### Updating Issues

```bash
bd update <id> --status=in_progress
bd update <id> --assignee=username
bd update <id> --priority=1
bd update <id> --title="New title"
```

### Closing Issues

```bash
bd close <id>                 # Mark complete
bd close <id1> <id2> ...      # Close multiple
bd close <id> --reason="..."  # Close with reason
```

### Dependencies

```bash
bd dep add <issue> <depends-on>   # Add dependency
bd dep remove <issue> <dep>       # Remove dependency
bd blocked                        # Show blocked issues
```

**Note:** Think "X needs Y" not "X before Y"
- `bd dep add task-B task-A` means "B depends on A" (A blocks B)

### Sync

```bash
bd sync                       # Sync with git remote
bd sync --status              # Check sync status
bd sync --force               # Force sync
```

### Statistics

```bash
bd stats                      # Project statistics
bd doctor                     # Check for issues
```

---

## bd mol - Molecule Commands

### Workflow Management

```bash
bd mol list                   # List active molecules
bd mol show <id>              # Show molecule details
bd mol steps <id>             # Show step status
```

### Lifecycle Operations

```bash
bd mol cook <formula>         # Formula → Protomolecule
bd mol pour <proto> <bead>    # Create persistent molecule
bd mol wisp <formula>         # Create ephemeral wisp
bd mol squash <id>            # Archive with summary
bd mol burn <id>              # Discard without record
```

---

## bd events - Event Commands

### Querying Events

```bash
bd events list                    # Last 50 events
bd events list --type=bead.*      # Filter by type
bd events list --actor=mayor      # Filter by actor
bd events list --since=1h         # Time filter
```

### Real-time

```bash
bd events tail                    # Watch feed.jsonl
```

### History

```bash
bd events replay <bead-id>        # Event history for bead
```

### Statistics

```bash
bd events stats                   # Counts by type
bd events stats --since=1d        # Last 24 hours
```

### Export

```bash
bd events export --since=2026-01-01 > events.jsonl
```

---

## bd eval - Evaluation Commands

```bash
bd eval completion --since=7d     # Completion rate
bd eval gupp --since=1d           # GUPP compliance
bd eval throughput --since=24h    # Work throughput
bd eval verify-accuracy --since=30d
bd eval report --since=7d --output=report.json
```

---

## Common Patterns

### Starting a Session (Any Agent)

```bash
gt prime                      # Load context
gt hook                       # Check for work
gt mail inbox                 # Check messages
```

### Mayor Workflow

```bash
gt status                     # Town overview
bd ready                      # Find work to dispatch
gt sling <bead> <rig>         # Assign to polecat
gt convoy list                # Monitor progress
```

### Polecat Workflow

```bash
gt hook                       # Find assigned work (GUPP)
bd show <bead-id>             # Understand the task
# ... do the work ...
git add . && git commit -m "..." && git push
gt polecat done               # Signal completion
```

### Witness Patrol

```bash
gt mail inbox                 # Check for POLECAT_DONE
gt polecat list               # Survey workers
# Nudge idle, kill stuck
```

### Session End

```bash
git status                    # Check changes
git add <files>               # Stage code
bd sync                       # Commit beads
git commit -m "..."           # Commit code
git push                      # Push to remote
gt handoff -m "..."           # If incomplete
```

---

## Environment Variables

| Variable | Description | Used By |
|----------|-------------|---------|
| `BD_ACTOR` | Agent identity | All agents |
| `BEAD_ID` | Current bead | Polecats |
| `GT_RIG` | Current rig | All agents |
| `GT_ROLE` | Agent role | All agents |
| `BD_DEBUG_ROUTING` | Debug prefix routing | Debugging |

### Setting Debug Mode

```bash
BD_DEBUG_ROUTING=1 bd show <id>   # Debug routing
GT_VERBOSE=1 gt status            # Verbose output
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
bd list
# gt-abc123  Feature X         open    P2
# gt-def456  Bug in login      open    P0
```

### JSON Output

```bash
bd list --json
# [{"id": "gt-abc123", "title": "Feature X", ...}]

bd show <id> --json
# {"id": "gt-abc123", ...}
```

### JSONL Output (for piping)

```bash
bd events list --jsonl | jq '.event_type'
```

---

## Piping and Composition

```bash
# Count open bugs
bd list --status=open --type=bug --json | jq length

# Get IDs of blocked issues
bd blocked --json | jq -r '.[].id'

# Close all completed in convoy
gt convoy status <id> --json | jq -r '.completed[].id' | xargs bd close

# Event analysis
bd events list --since=1d --jsonl | jq -r '.event_type' | sort | uniq -c
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
[beads]
default_priority = 2
sync_branch = "beads-sync"

[mail]
default_from = "mayor"

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
