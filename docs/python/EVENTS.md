# VerMAS Event System

> Event sourcing, change feeds, and audit trails

## Design Philosophy

VerMAS uses **event sourcing** as its foundational data pattern. Every state change in the system is captured as an immutable event in append-only logs. The current state of any entity can be reconstructed by replaying its event history.

```
                    PROVEN TECHNOLOGY STACK
+--------------------------------------------------------------------+
|                                                                    |
|   File System          Git              CLI           Agents       |
|   +-----------+    +-----------+    +-----------+    +-----------+ |
|   | JSONL     |    | Commits   |    | gt / bd   |    | Claude    | |
|   | Append-   |    | Branches  |    | Commands  |    | Code CLI  | |
|   | Only Logs |    | Worktrees |    | Typer     |    | Profiles  | |
|   +-----------+    +-----------+    +-----------+    +-----------+ |
|        |                |                |                |        |
|        +----------------+----------------+----------------+        |
|                         |                                          |
|                         v                                          |
|                  +-------------+                                   |
|                  | Event Store |                                   |
|                  | (.events.   |                                   |
|                  |   jsonl)    |                                   |
|                  +-------------+                                   |
|                                                                    |
+--------------------------------------------------------------------+
```

---

## Event Types

### Core Events

| Event | Description | Emitter |
|-------|-------------|---------|
| `bead.created` | New bead/issue created | bd create |
| `bead.updated` | Bead field changed | bd update |
| `bead.status_changed` | Status transition | bd update/close |
| `bead.hooked` | Bead assigned to agent hook | gt sling |
| `bead.closed` | Bead marked complete | bd close |

### Message Events

| Event | Description | Emitter |
|-------|-------------|---------|
| `mail.sent` | Message dispatched | gt mail send |
| `mail.delivered` | Message arrived in inbox | mail router |
| `mail.read` | Recipient read message | gt mail read |

### Agent Lifecycle Events

| Event | Description | Emitter |
|-------|-------------|---------|
| `agent.started` | Agent session began | gt polecat spawn |
| `agent.hook_checked` | Agent checked its hook | gt hook |
| `agent.working` | Agent began work on bead | GUPP execution |
| `agent.idle` | Agent went idle | witness patrol |
| `agent.nudged` | Agent received nudge | witness |
| `agent.stopped` | Agent session ended | gt polecat done |

### Workflow Events

| Event | Description | Emitter |
|-------|-------------|---------|
| `mol.created` | Molecule instantiated | bd mol pour |
| `mol.step_started` | Workflow step began | agent |
| `mol.step_completed` | Workflow step finished | agent |
| `mol.completed` | All steps done | bd mol squash |
| `mol.abandoned` | Workflow discarded | bd mol burn |

### Verification Events (VerMAS-specific)

| Event | Description | Emitter |
|-------|-------------|---------|
| `verify.started` | Verification began | refinery |
| `verify.spec_created` | Designer produced spec | designer |
| `verify.tests_generated` | Strategist created tests | strategist |
| `verify.test_executed` | Verifier ran a test | verifier |
| `verify.audited` | Auditor reviewed evidence | auditor |
| `verify.verdict` | Final PASS/FAIL decision | judge |

---

## Event Schema

### Base Event Structure

Every event shares this schema:

```json
{
  "event_id": "evt-abc123",
  "event_type": "bead.created",
  "timestamp": "2026-01-06T12:00:00Z",
  "actor": "gastown/crew/frontend",
  "correlation_id": "conv-xyz789",
  "data": { ... }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `event_id` | string | Unique event identifier |
| `event_type` | string | Event type (namespaced) |
| `timestamp` | datetime | When event occurred (ISO 8601) |
| `actor` | string | BD_ACTOR who triggered event |
| `correlation_id` | string | Links related events (convoy, molecule) |
| `data` | object | Event-specific payload |

### Example Events

**Bead Created:**
```json
{
  "event_id": "evt-a1b2c3",
  "event_type": "bead.created",
  "timestamp": "2026-01-06T10:30:00Z",
  "actor": "mayor",
  "correlation_id": null,
  "data": {
    "bead_id": "gt-def456",
    "title": "Implement feature X",
    "issue_type": "feature",
    "priority": 2
  }
}
```

**Status Changed:**
```json
{
  "event_id": "evt-d4e5f6",
  "event_type": "bead.status_changed",
  "timestamp": "2026-01-06T11:00:00Z",
  "actor": "gastown/polecats/slot0",
  "correlation_id": "conv-123",
  "data": {
    "bead_id": "gt-def456",
    "from_status": "open",
    "to_status": "in_progress"
  }
}
```

**Agent Started:**
```json
{
  "event_id": "evt-g7h8i9",
  "event_type": "agent.started",
  "timestamp": "2026-01-06T10:45:00Z",
  "actor": "gastown/polecats/slot0",
  "correlation_id": null,
  "data": {
    "session_name": "polecat-gastown-slot0",
    "worktree": "/path/to/worktree",
    "profile": "polecat",
    "hooked_bead": "gt-def456"
  }
}
```

**Verification Verdict:**
```json
{
  "event_id": "evt-j1k2l3",
  "event_type": "verify.verdict",
  "timestamp": "2026-01-06T12:15:00Z",
  "actor": "gastown/inspector/judge",
  "correlation_id": "mol-verify-abc",
  "data": {
    "bead_id": "gt-def456",
    "verdict": "PASS",
    "criteria_passed": 5,
    "criteria_failed": 0,
    "reasoning": "All acceptance criteria met..."
  }
}
```

---

## Event Storage

### File Layout

```
.beads/
├── events.jsonl           # Main event log (append-only)
├── events/
│   ├── 2026-01-06.jsonl   # Daily partitioned logs
│   └── 2026-01-05.jsonl
└── feed.jsonl             # Real-time subscription feed
```

### JSONL Format

Events stored as newline-delimited JSON (JSONL):
- One event per line
- Append-only writes
- Git-friendly diffs
- Easy to grep/tail/process

```bash
# View recent events
tail -20 .beads/events.jsonl | jq .

# Filter by event type
grep '"event_type":"bead.created"' .beads/events.jsonl

# Count events by type
cat .beads/events.jsonl | jq -r '.event_type' | sort | uniq -c
```

---

## Event Sourcing Patterns

### State Reconstruction

Current bead state = replay all events for that bead:

```python
def get_bead_state(bead_id: str) -> Bead:
    events = load_events(filter={"data.bead_id": bead_id})
    state = {}

    for event in events:
        if event.event_type == "bead.created":
            state = event.data
        elif event.event_type == "bead.updated":
            state.update(event.data.get("changes", {}))
        elif event.event_type == "bead.status_changed":
            state["status"] = event.data["to_status"]

    return Bead(**state)
```

### Projection (Materialized Views)

For performance, maintain materialized views:

```
events.jsonl  →  project  →  issues.jsonl (current state)
                          →  messages.jsonl (mailboxes)
                          →  metrics.jsonl (aggregates)
```

**The `issues.jsonl` file is a projection of bead events**, not the source of truth. It's an optimization for fast reads.

### Temporal Queries

Event sourcing enables "time travel":

```python
def get_bead_state_at(bead_id: str, timestamp: datetime) -> Bead:
    events = load_events(
        filter={"data.bead_id": bead_id},
        before=timestamp
    )
    return replay_events(events)
```

---

## Change Feed

### Real-time Subscription

Agents can subscribe to the change feed:

```
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│   Producer  │──────▶│  feed.jsonl │──────▶│  Consumers  │
│  (bd, gt)   │       │  (append)   │       │  (agents)   │
└─────────────┘       └─────────────┘       └─────────────┘
```

**Producer** (CLI commands) appends to `feed.jsonl`
**Consumers** (Witness, Refinery) tail the feed

### Implementation

```python
# Producer: emit event
def emit_event(event: Event):
    with open(".beads/events.jsonl", "a") as f:
        f.write(event.model_dump_json() + "\n")

    # Also append to real-time feed
    with open(".beads/feed.jsonl", "a") as f:
        f.write(event.model_dump_json() + "\n")

# Consumer: watch feed
async def watch_feed():
    with open(".beads/feed.jsonl", "r") as f:
        f.seek(0, 2)  # Go to end
        while True:
            line = f.readline()
            if line:
                event = Event.model_validate_json(line)
                yield event
            else:
                await asyncio.sleep(0.1)
```

### File-based Pub/Sub

The feed file acts as a simple pub/sub mechanism:
- No external message broker needed
- Uses proven file system primitives
- Git-compatible
- Easy to debug (just `tail -f`)

---

## Correlation and Causation

### Correlation ID

Links events that belong together:
- All events in a convoy share `correlation_id`
- All events in a molecule execution share `correlation_id`
- Enables tracing work across agents

### Causation Chain

For detailed tracing, events can include `caused_by`:

```json
{
  "event_id": "evt-xyz",
  "event_type": "mail.sent",
  "caused_by": "evt-abc",
  "data": { ... }
}
```

This creates an audit trail: `bead.created` → `bead.hooked` → `agent.started` → `agent.working` → ...

---

## CLI Commands for Events

```bash
# List recent events
bd events list                    # Last 50 events
bd events list --type=bead.*      # Filter by type
bd events list --actor=mayor      # Filter by actor

# Tail events (real-time)
bd events tail                    # Watch feed.jsonl

# Replay events for a bead
bd events replay <bead-id>        # Show event history

# Event statistics
bd events stats                   # Counts by type
bd events stats --since=1d        # Last 24 hours

# Export events
bd events export --since=2026-01-01 > events-january.jsonl
```

---

## Integration with Agents

### Hooks as Events

When a hook is set, emit `bead.hooked`:

```python
def set_hook(agent: str, bead_id: str):
    # Write hook file
    hook_path = Path(f".beads/.hook-{agent}")
    hook_path.write_text(f"bead:{bead_id}")

    # Emit event
    emit_event(Event(
        event_type="bead.hooked",
        actor="system",
        data={
            "bead_id": bead_id,
            "agent": agent,
            "hook_path": str(hook_path)
        }
    ))
```

### GUPP Compliance Events

Track propulsion principle compliance:

```json
{
  "event_type": "agent.gupp_check",
  "data": {
    "agent": "gastown/polecats/slot0",
    "hook_found": true,
    "response_ms": 150,
    "action": "execute_immediately"
  }
}
```

### Witness Patrol Events

```json
{
  "event_type": "patrol.completed",
  "actor": "gastown/witness",
  "data": {
    "polecats_checked": 3,
    "idle_found": 1,
    "nudges_sent": 1,
    "kills_executed": 0
  }
}
```

---

## Evaluation via Events

Events enable precise metrics:

```python
# Completion rate
completed = count_events(type="bead.status_changed", to="closed")
total = count_events(type="bead.created")
rate = completed / total

# GUPP compliance
gupp_events = get_events(type="agent.gupp_check")
fast_responses = [e for e in gupp_events if e.data["response_ms"] < 30000]
compliance = len(fast_responses) / len(gupp_events)

# Verification accuracy
verdicts = get_events(type="verify.verdict")
passes = [v for v in verdicts if v.data["verdict"] == "PASS"]
pass_rate = len(passes) / len(verdicts)
```

See [EVALUATION.md](./EVALUATION.md) for detailed metrics.

---

## Why Event Sourcing?

### Proven Technology

Event sourcing is used by:
- Banking systems (transaction logs)
- Kafka/streaming platforms
- Git itself (commits are events)
- Redux/Flux (frontend state)

### Benefits for VerMAS

1. **Complete Audit Trail** - Every state change is recorded
2. **Debugging** - Replay events to understand what happened
3. **Metrics** - Compute any metric from events
4. **Recovery** - Reconstruct state after failures
5. **Temporal Queries** - What was the state at time T?
6. **Causation Tracking** - Why did this happen?

### Why JSONL?

1. **Append-only** - No read-modify-write races
2. **Git-friendly** - Line-based diffs, easy merges
3. **Unix-friendly** - grep, tail, head, jq
4. **Language-agnostic** - Go and Python both work
5. **Debuggable** - Human-readable text

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [OPERATIONS.md](./OPERATIONS.md) - Monitoring and maintenance
- [HOOKS.md](./HOOKS.md) - Claude Code integration and git worktrees
- [MESSAGING.md](./MESSAGING.md) - Mail protocol
- [SCHEMAS.md](./SCHEMAS.md) - Event data specifications
- [EVALUATION.md](./EVALUATION.md) - Metrics from events
- [WORKFLOWS.md](./WORKFLOWS.md) - Molecule lifecycle
