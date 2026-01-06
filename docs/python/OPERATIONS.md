# VerMAS Operations Guide

> Deployment, startup, monitoring, and maintenance

## Quick Start

### Prerequisites

| Component | Version | Check Command |
|-----------|---------|---------------|
| Python | 3.10+ | `python --version` |
| Git | 2.20+ | `git --version` |
| Tmux | 3.0+ | `tmux -V` |
| Claude Code | Latest | `claude --version` |

### Initial Setup

```bash
# 1. Install VerMAS
pip install vermas

# 2. Initialize town root
mkdir ~/gt && cd ~/gt
gt init

# 3. Add a rig (project)
gt rig add myproject git@github.com:org/myproject.git

# 4. Create Claude profiles
mkdir -p ~/.claude/profiles/{mayor,witness,refinery,polecat}
# Copy CLAUDE.md files to each profile directory

# 5. Verify setup
gt doctor
```

---

## Starting the System

### Boot Sequence

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           STARTUP SEQUENCE                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   1. Deacon (infrastructure watchdog)                                       │
│      gt deacon start                                                        │
│           │                                                                 │
│           ▼                                                                 │
│   2. Per-rig Witnesses                                                      │
│      (Started by Deacon automatically)                                      │
│           │                                                                 │
│           ▼                                                                 │
│   3. Per-rig Refineries                                                     │
│      (Started by Deacon automatically)                                      │
│           │                                                                 │
│           ▼                                                                 │
│   4. Mayor (when human starts session)                                      │
│      claude --profile mayor                                                 │
│                                                                             │
│   Polecats spawn on-demand via gt sling                                     │
│                                                                             │
│   GUPP: When agents find work on their hook, they EXECUTE immediately.      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Starting Deacon

```bash
# Start Deacon (runs in background)
gt deacon start

# Check Deacon status
gt deacon status

# View Deacon logs
gt deacon logs

# Stop Deacon
gt deacon stop
```

### Starting Mayor Session

```bash
# Navigate to town root
cd ~/gt

# Start Mayor with profile
claude --profile mayor

# Or use the convenience command
gt mayor
```

### Manual Service Start (without Deacon)

```bash
# Start Witness for a rig
tmux new-session -d -s witness-myproject \
  -e "BD_ACTOR=myproject/witness" \
  "claude --profile witness"

# Start Refinery for a rig
tmux new-session -d -s refinery-myproject \
  -e "BD_ACTOR=myproject/refinery" \
  "claude --profile refinery"
```

---

## Stopping the System

### Graceful Shutdown

```bash
# 1. Stop Deacon (stops Witnesses and Refineries)
gt deacon stop

# 2. Wait for active polecats to complete
gt polecat list --all
# Wait for empty list

# 3. Kill any remaining sessions
tmux kill-server
```

### Emergency Stop

```bash
# Kill everything immediately
tmux kill-server

# Note: Work in progress may need recovery
# Check .beads/.hook-* files for interrupted work
```

---

## Monitoring

### Health Checks

```bash
# Overall system status
gt status

# Per-rig status
gt rig status myproject

# List all active sessions
tmux list-sessions

# Check event feed (real-time)
tail -f .beads/feed.jsonl | jq .
```

### Monitoring Commands

| Command | What it Shows |
|---------|---------------|
| `gt status` | Town overview: rigs, agents, mail |
| `gt polecat list` | Active polecats per rig |
| `gt convoy list` | Work in progress |
| `bd stats` | Bead counts and health |
| `bd events stats --since=1h` | Event volume |

### Automated Monitoring

```bash
# Set up monitoring script
cat > ~/gt/monitor.sh << 'EOF'
#!/bin/bash
while true; do
    clear
    echo "=== Gas Town Status $(date) ==="
    gt status
    echo ""
    echo "=== Active Sessions ==="
    tmux list-sessions 2>/dev/null || echo "No sessions"
    echo ""
    echo "=== Recent Events ==="
    tail -5 .beads/feed.jsonl | jq -c '{t: .event_type, a: .actor}'
    sleep 30
done
EOF
chmod +x ~/gt/monitor.sh
```

### Dashboard Views

```bash
# Pane 1: Event stream
tail -f .beads/feed.jsonl | jq .

# Pane 2: System status (refreshes)
watch -n 30 'gt status'

# Pane 3: Active work
watch -n 60 'gt convoy list'
```

---

## Common Operations

### Adding a New Project (Rig)

```bash
# Add rig with git remote
gt rig add newproject git@github.com:org/newproject.git

# Verify setup
gt rig status newproject

# Create crew workspace if needed
gt crew add dev --rig newproject
```

### Spawning a Polecat

```bash
# Create work
bd create --title="Fix bug X" --type=bug --priority=1
# Output: gt-abc123

# Sling to rig (spawns polecat automatically)
gt sling gt-abc123 myproject
```

### Checking on Stuck Work

```bash
# Find idle polecats
gt polecat list --status=idle

# Check specific polecat
gt polecat status myproject slot0

# Attach to see what's happening
tmux attach -t polecat-myproject-slot0
# Ctrl+B D to detach

# Nudge manually if needed
gt mail send myproject/polecats/slot0 -s "NUDGE" -m "Wake up"

# Kill if truly stuck
gt polecat kill myproject slot0
```

### Recovering from Failure

```bash
# Check for orphaned hooks
ls .beads/.hook-*
cat .beads/.hook-myproject-polecats-slot0

# Clean up stale hooks
rm .beads/.hook-myproject-polecats-slot0

# Prune stale worktrees
git worktree prune

# Re-sling work
gt sling gt-abc123 myproject
```

---

## Configuration

### Town Configuration

```toml
# ~/gt/.vermas/config.toml

[town]
name = "my-gas-town"

[deacon]
patrol_interval_seconds = 60
restart_on_failure = true

[polecat]
max_slots_per_rig = 5
idle_timeout_minutes = 5
stuck_timeout_minutes = 15

[mail]
poll_interval_seconds = 10

[events]
partition_by_day = true
retention_days = 30
```

### Rig Configuration

```toml
# myproject/.beads/config.toml

[rig]
name = "myproject"
prefix = "mp"

[verification]
enabled = true
objective_only = false
require_all_pass = true

[sync]
branch = "beads-sync"
auto_sync = true
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `BD_ACTOR` | Agent identity | Required for agents |
| `GT_TOWN` | Town root path | `~/gt` |
| `GT_LOG_LEVEL` | Logging verbosity | `INFO` |
| `GT_DEBUG` | Enable debug mode | `false` |

---

## Maintenance

### Daily Tasks

```bash
# Check system health
gt doctor

# Review any failed work
bd list --status=failed

# Check for stale polecats
gt polecat list --status=idle

# Sync beads
bd sync --all-rigs
```

### Weekly Tasks

```bash
# Archive old events
bd events archive --older-than=7d

# Clean up worktrees
git worktree prune

# Review metrics
bd eval report --since=7d

# Check disk usage
du -sh ~/gt/.beads/
```

### Monthly Tasks

```bash
# Full backup
tar -czf ~/backups/gt-$(date +%Y%m).tar.gz ~/gt/.beads/

# Review and close stale beads
bd list --status=open --older-than=30d

# Update Claude profiles if needed
# Review CLAUDE.md files for improvements
```

---

## Backup and Recovery

### What to Back Up

| Path | Contents | Frequency |
|------|----------|-----------|
| `.beads/events.jsonl` | Event log (source of truth) | Daily |
| `.beads/issues.jsonl` | Bead state | Daily |
| `.beads/messages.jsonl` | Mail archive | Daily |
| `.beads/formulas/*.toml` | Workflow templates | On change |
| `~/.claude/profiles/` | Agent profiles | On change |

### Backup Script

```bash
#!/bin/bash
# backup-gt.sh

BACKUP_DIR=~/backups/gt
DATE=$(date +%Y%m%d)

mkdir -p $BACKUP_DIR

# Backup town beads
tar -czf $BACKUP_DIR/town-$DATE.tar.gz ~/gt/.beads/

# Backup each rig's beads
for rig in ~/gt/*/; do
    if [ -d "$rig/.beads" ]; then
        name=$(basename $rig)
        tar -czf $BACKUP_DIR/rig-$name-$DATE.tar.gz $rig/.beads/
    fi
done

# Cleanup old backups (keep 30 days)
find $BACKUP_DIR -mtime +30 -delete

echo "Backup completed: $BACKUP_DIR"
```

### Recovery Procedures

**Restore from backup:**
```bash
# Stop all agents
gt deacon stop
tmux kill-server

# Restore beads
tar -xzf ~/backups/gt/town-20260106.tar.gz -C ~/gt/

# Restart
gt deacon start
```

**Rebuild projections from events:**
```bash
# If issues.jsonl is corrupted
bd rebuild-projections

# This replays events.jsonl to regenerate:
# - issues.jsonl
# - messages.jsonl
```

---

## Troubleshooting

### Common Issues

#### "Command not found: gt"
```bash
# Ensure VerMAS is installed
pip install vermas

# Or check PATH
which gt
```

#### "No rigs found"
```bash
# Initialize town first
gt init

# Add a rig
gt rig add myproject <git-url>
```

#### "Deacon won't start"
```bash
# Check for existing process
pgrep -f "gt deacon"

# Kill and restart
pkill -f "gt deacon"
gt deacon start
```

#### "Polecat not spawning"
```bash
# Check available slots
gt polecat list myproject

# Check for worktree issues
git worktree list

# Try manual spawn
gt polecat spawn myproject --verbose
```

#### "Events not appearing"
```bash
# Check file permissions
ls -la .beads/events.jsonl

# Check disk space
df -h .

# Verify file is writable
touch .beads/test && rm .beads/test
```

### Debug Mode

```bash
# Enable verbose logging
export GT_LOG_LEVEL=DEBUG
export GT_DEBUG=true

# Run command with debug
GT_DEBUG=1 gt sling gt-abc123 myproject

# Debug routing
BD_DEBUG_ROUTING=1 bd show gt-abc123
```

### Log Locations

| Log | Path | Contains |
|-----|------|----------|
| Deacon | `~/gt/logs/deacon.log` | Infrastructure events |
| Session | `~/gt/logs/{session}/` | Agent output |
| Events | `.beads/events.jsonl` | All state changes |
| Feed | `.beads/feed.jsonl` | Real-time stream |

---

## Performance Tuning

### Optimize for Many Beads

```bash
# Enable daily event partitioning
# In config.toml:
[events]
partition_by_day = true

# Archive old events
bd events archive --older-than=30d
```

### Optimize for Many Polecats

```bash
# Increase tmux limits
# In ~/.tmux.conf:
set -g history-limit 50000

# Monitor resource usage
htop -p $(pgrep -d, -f "claude")
```

### Reduce Disk I/O

```bash
# Use SSDs for .beads directory
# Or symlink to fast storage:
mv .beads /ssd/.beads
ln -s /ssd/.beads .beads
```

---

## See Also

- [HOW_IT_WORKS.md](./HOW_IT_WORKS.md) - System overview
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Technical design
- [CLI.md](./CLI.md) - Command reference
- [HOOKS.md](./HOOKS.md) - Troubleshooting agent issues
- [EVENTS.md](./EVENTS.md) - Event log format
- [EVALUATION.md](./EVALUATION.md) - Metrics and monitoring
