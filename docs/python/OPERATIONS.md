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

# 2. Initialize company root
mkdir ~/company && cd ~/company
co init

# 3. Add a factory (project)
co factory add myproject git@github.com:org/myproject.git

# 4. Create Claude profiles
mkdir -p ~/.claude/profiles/{ceo,supervisor,qa,worker}
# Copy CLAUDE.md files to each profile directory

# 5. Verify setup
co doctor
```

---

## Starting the System

### Boot Sequence

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           STARTUP SEQUENCE                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   1. Operations (infrastructure watchdog)                                   │
│      co operations start                                                    │
│           │                                                                 │
│           ▼                                                                 │
│   2. Per-factory Supervisors                                                │
│      (Started by Operations automatically)                                  │
│           │                                                                 │
│           ▼                                                                 │
│   3. Per-factory QA Departments                                             │
│      (Started by Operations automatically)                                  │
│           │                                                                 │
│           ▼                                                                 │
│   4. CEO (when human starts session)                                        │
│      claude --profile ceo                                                   │
│                                                                             │
│   Workers spawn on-demand via co dispatch                                   │
│                                                                             │
│   Assignment Principle: When agents find work assigned, they EXECUTE.       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Starting Operations

```bash
# Start Operations (runs in background)
co operations start

# Check Operations status
co operations status

# View Operations logs
co operations logs

# Stop Operations
co operations stop
```

### Starting CEO Session

```bash
# Navigate to company root
cd ~/company

# Start CEO with profile
claude --profile ceo

# Or use the convenience command
co ceo
```

### Manual Service Start (without Operations)

```bash
# Start Supervisor for a factory
tmux new-session -d -s supervisor-myproject \
  -e "AGENT_ID=myproject/supervisor" \
  "claude --profile supervisor"

# Start QA for a factory
tmux new-session -d -s qa-myproject \
  -e "AGENT_ID=myproject/qa" \
  "claude --profile qa"
```

---

## Stopping the System

### Graceful Shutdown

```bash
# 1. Stop Operations (stops Supervisors and QA)
co operations stop

# 2. Wait for active workers to complete
co worker list --all
# Wait for empty list

# 3. Kill any remaining sessions
tmux kill-server
```

### Emergency Stop

```bash
# Kill everything immediately
tmux kill-server

# Note: Work in progress may need recovery
# Check .work/.assignment-* files for interrupted work
```

---

## Monitoring

### Health Checks

```bash
# Overall system status
co status

# Per-factory status
co factory status myproject

# List all active sessions
tmux list-sessions

# Check event feed (real-time)
tail -f .work/feed.jsonl | jq .
```

### Monitoring Commands

| Command | What it Shows |
|---------|---------------|
| `co status` | Company overview: factories, agents, mail |
| `co worker list` | Active workers per factory |
| `co sprint list` | Work in progress |
| `wo stats` | Work order counts and health |
| `wo events stats --since=1h` | Event volume |

### Automated Monitoring

```bash
# Set up monitoring script
cat > ~/company/monitor.sh << 'EOF'
#!/bin/bash
while true; do
    clear
    echo "=== VerMAS Status $(date) ==="
    co status
    echo ""
    echo "=== Active Sessions ==="
    tmux list-sessions 2>/dev/null || echo "No sessions"
    echo ""
    echo "=== Recent Events ==="
    tail -5 .work/feed.jsonl | jq -c '{t: .event_type, a: .actor}'
    sleep 30
done
EOF
chmod +x ~/company/monitor.sh
```

### Dashboard Views

```bash
# Pane 1: Event stream
tail -f .work/feed.jsonl | jq .

# Pane 2: System status (refreshes)
watch -n 30 'co status'

# Pane 3: Active work
watch -n 60 'co sprint list'
```

---

## Common Operations

### Adding a New Project (Factory)

```bash
# Add factory with git remote
co factory add newproject git@github.com:org/newproject.git

# Verify setup
co factory status newproject

# Create team workspace if needed
co team add dev --factory newproject
```

### Spawning a Worker

```bash
# Create work
wo create --title="Fix bug X" --type=bug --priority=1
# Output: wo-abc123

# Dispatch to factory (spawns worker automatically)
co dispatch wo-abc123 myproject
```

### Checking on Stuck Work

```bash
# Find idle workers
co worker list --status=idle

# Check specific worker
co worker status myproject slot0

# Attach to see what's happening
tmux attach -t worker-myproject-slot0
# Ctrl+B D to detach

# Nudge manually if needed
co send myproject/workers/slot0 -s "NUDGE" -m "Wake up"

# Kill if truly stuck
co worker kill myproject slot0
```

### Recovering from Failure

```bash
# Check for orphaned assignments
ls .work/.assignment-*
cat .work/.assignment-myproject-workers-slot0

# Clean up stale assignments
rm .work/.assignment-myproject-workers-slot0

# Prune stale worktrees
git worktree prune

# Re-dispatch work
co dispatch wo-abc123 myproject
```

---

## Configuration

### Company Configuration

```toml
# ~/company/.vermas/config.toml

[company]
name = "my-company"

[operations]
patrol_interval_seconds = 60
restart_on_failure = true

[worker]
max_slots_per_factory = 5
idle_timeout_minutes = 5
stuck_timeout_minutes = 15

[mail]
poll_interval_seconds = 10

[events]
partition_by_day = true
retention_days = 30
```

### Factory Configuration

```toml
# myproject/.work/config.toml

[factory]
name = "myproject"
prefix = "mp"

[verification]
enabled = true
objective_only = false
require_all_pass = true

[sync]
branch = "work-sync"
auto_sync = true
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `AGENT_ID` | Agent identity | Required for agents |
| `COMPANY_ROOT` | Company root path | `~/company` |
| `LOG_LEVEL` | Logging verbosity | `INFO` |
| `DEBUG` | Enable debug mode | `false` |

---

## Maintenance

### Daily Tasks

```bash
# Check system health
co doctor

# Review any failed work
wo list --status=failed

# Check for stale workers
co worker list --status=idle

# Sync work orders
wo sync --all-factories
```

### Weekly Tasks

```bash
# Archive old events
wo events archive --older-than=7d

# Clean up worktrees
git worktree prune

# Review metrics
wo eval report --since=7d

# Check disk usage
du -sh ~/company/.work/
```

### Monthly Tasks

```bash
# Full backup
tar -czf ~/backups/company-$(date +%Y%m).tar.gz ~/company/.work/

# Review and close stale work orders
wo list --status=open --older-than=30d

# Update Claude profiles if needed
# Review CLAUDE.md files for improvements
```

---

## Backup and Recovery

### What to Back Up

| Path | Contents | Frequency |
|------|----------|-----------|
| `.work/events.jsonl` | Event log (source of truth) | Daily |
| `.work/work_orders.jsonl` | Work order state | Daily |
| `.work/messages.jsonl` | Mail archive | Daily |
| `.work/templates/*.toml` | Workflow templates | On change |
| `~/.claude/profiles/` | Agent profiles | On change |

### Backup Script

```bash
#!/bin/bash
# backup-company.sh

BACKUP_DIR=~/backups/company
DATE=$(date +%Y%m%d)

mkdir -p $BACKUP_DIR

# Backup company work data
tar -czf $BACKUP_DIR/company-$DATE.tar.gz ~/company/.work/

# Backup each factory's work data
for factory in ~/company/*/; do
    if [ -d "$factory/.work" ]; then
        name=$(basename $factory)
        tar -czf $BACKUP_DIR/factory-$name-$DATE.tar.gz $factory/.work/
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
co operations stop
tmux kill-server

# Restore work data
tar -xzf ~/backups/company/company-20260106.tar.gz -C ~/company/

# Restart
co operations start
```

**Rebuild projections from events:**
```bash
# If work_orders.jsonl is corrupted
wo rebuild-projections

# This replays events.jsonl to regenerate:
# - work_orders.jsonl
# - messages.jsonl
```

---

## Troubleshooting

### Common Issues

#### "Command not found: co"
```bash
# Ensure VerMAS is installed
pip install vermas

# Or check PATH
which co
```

#### "No factories found"
```bash
# Initialize company first
co init

# Add a factory
co factory add myproject <git-url>
```

#### "Operations won't start"
```bash
# Check for existing process
pgrep -f "co operations"

# Kill and restart
pkill -f "co operations"
co operations start
```

#### "Worker not spawning"
```bash
# Check available slots
co worker list myproject

# Check for worktree issues
git worktree list

# Try manual spawn
co worker spawn myproject --verbose
```

#### "Events not appearing"
```bash
# Check file permissions
ls -la .work/events.jsonl

# Check disk space
df -h .

# Verify file is writable
touch .work/test && rm .work/test
```

### Debug Mode

```bash
# Enable verbose logging
export LOG_LEVEL=DEBUG
export DEBUG=true

# Run command with debug
DEBUG=1 co dispatch wo-abc123 myproject

# Debug routing
DEBUG_ROUTING=1 wo show wo-abc123
```

### Log Locations

| Log | Path | Contains |
|-----|------|----------|
| Operations | `~/company/logs/operations.log` | Infrastructure events |
| Session | `~/company/logs/{session}/` | Agent output |
| Events | `.work/events.jsonl` | All state changes |
| Feed | `.work/feed.jsonl` | Real-time stream |

---

## Performance Tuning

### Optimize for Many Work Orders

```bash
# Enable daily event partitioning
# In config.toml:
[events]
partition_by_day = true

# Archive old events
wo events archive --older-than=30d
```

### Optimize for Many Workers

```bash
# Increase tmux limits
# In ~/.tmux.conf:
set -g history-limit 50000

# Monitor resource usage
htop -p $(pgrep -d, -f "claude")
```

### Reduce Disk I/O

```bash
# Use SSDs for .work directory
# Or symlink to fast storage:
mv .work /ssd/.work
ln -s /ssd/.work .work
```

---

## See Also

- [HOW_IT_WORKS.md](./HOW_IT_WORKS.md) - System overview
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Technical design
- [CLI.md](./CLI.md) - Command reference
- [HOOKS.md](./HOOKS.md) - Troubleshooting agent issues
- [EVENTS.md](./EVENTS.md) - Event log format
- [EVALUATION.md](./EVALUATION.md) - Metrics and monitoring
