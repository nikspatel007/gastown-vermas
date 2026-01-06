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
|   | JSONL     |    | Commits   |    | co / wo   |    | Claude    | |
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
| `work_order.created` | New work order created | wo create |
| `work_order.updated` | Work order field changed | wo update |
| `work_order.status_changed` | Status transition | wo update/close |
| `work_order.assigned` | Work order assigned to agent | co dispatch |
| `work_order.closed` | Work order marked complete | wo close |

### Message Events

| Event | Description | Emitter |
|-------|-------------|---------|
| `mail.sent` | Message dispatched | co send |
| `mail.delivered` | Message arrived in inbox | mail router |
| `mail.read` | Recipient read message | co read |

### Agent Lifecycle Events

| Event | Description | Emitter |
|-------|-------------|---------|
| `agent.started` | Agent session began | co worker spawn |
| `agent.assignment_checked` | Agent checked assignment | co assignment |
| `agent.working` | Agent began work on order | Assignment execution |
| `agent.idle` | Agent went idle | supervisor patrol |
| `agent.nudged` | Agent received nudge | supervisor |
| `agent.stopped` | Agent session ended | co worker done |

### Workflow Events

| Event | Description | Emitter |
|-------|-------------|---------|
| `process.created` | Process instantiated | wo process start |
| `process.step_started` | Workflow step began | agent |
| `process.step_completed` | Workflow step finished | agent |
| `process.completed` | All steps done | wo process complete |
| `process.cancelled` | Workflow discarded | wo process cancel |

### Verification Events (VerMAS-specific)

| Event | Description | Emitter |
|-------|-------------|---------|
| `verify.started` | Verification began | qa |
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
  "event_type": "work_order.created",
  "timestamp": "2026-01-06T12:00:00Z",
  "actor": "project-a/teams/frontend",
  "correlation_id": "sprint-xyz789",
  "data": { ... }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `event_id` | string | Unique event identifier |
| `event_type` | string | Event type (namespaced) |
| `timestamp` | datetime | When event occurred (ISO 8601) |
| `actor` | string | AGENT_ID who triggered event |
| `correlation_id` | string | Links related events (sprint, process) |
| `data` | object | Event-specific payload |

### Example Events

**Work Order Created:**
```json
{
  "event_id": "evt-a1b2c3",
  "event_type": "work_order.created",
  "timestamp": "2026-01-06T10:30:00Z",
  "actor": "ceo",
  "correlation_id": null,
  "data": {
    "work_order_id": "wo-def456",
    "title": "Implement feature X",
    "type": "feature",
    "priority": 2
  }
}
```

**Status Changed:**
```json
{
  "event_id": "evt-d4e5f6",
  "event_type": "work_order.status_changed",
  "timestamp": "2026-01-06T11:00:00Z",
  "actor": "project-a/workers/slot0",
  "correlation_id": "sprint-123",
  "data": {
    "work_order_id": "wo-def456",
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
  "actor": "project-a/workers/slot0",
  "correlation_id": null,
  "data": {
    "session_name": "worker-project-a-slot0",
    "worktree": "/path/to/worktree",
    "profile": "worker",
    "assigned_work_order": "wo-def456"
  }
}
```

**Verification Verdict:**
```json
{
  "event_id": "evt-j1k2l3",
  "event_type": "verify.verdict",
  "timestamp": "2026-01-06T12:15:00Z",
  "actor": "project-a/qa/judge",
  "correlation_id": "process-verify-abc",
  "data": {
    "work_order_id": "wo-def456",
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
.work/
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
tail -20 .work/events.jsonl | jq .

# Filter by event type
grep '"event_type":"work_order.created"' .work/events.jsonl

# Count events by type
cat .work/events.jsonl | jq -r '.event_type' | sort | uniq -c
```

---

## Event Sourcing Patterns

### State Reconstruction

Current work order state = replay all events for that order:

```python
def get_work_order_state(wo_id: str) -> WorkOrder:
    events = load_events(filter={"data.work_order_id": wo_id})
    state = {}

    for event in events:
        if event.event_type == "work_order.created":
            state = event.data
        elif event.event_type == "work_order.updated":
            state.update(event.data.get("changes", {}))
        elif event.event_type == "work_order.status_changed":
            state["status"] = event.data["to_status"]

    return WorkOrder(**state)
```

### Projection (Materialized Views)

For performance, maintain materialized views:

```
events.jsonl  →  project  →  work_orders.jsonl (current state)
                          →  messages.jsonl (mailboxes)
                          →  metrics.jsonl (aggregates)
```

**The `work_orders.jsonl` file is a projection of work order events**, not the source of truth. It's an optimization for fast reads.

### Temporal Queries

Event sourcing enables "time travel":

```python
def get_work_order_state_at(wo_id: str, timestamp: datetime) -> WorkOrder:
    events = load_events(
        filter={"data.work_order_id": wo_id},
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
│  (co, wo)   │       │  (append)   │       │  (agents)   │
└─────────────┘       └─────────────┘       └─────────────┘
```

**Producer** (CLI commands) appends to `feed.jsonl`
**Consumers** (Supervisor, QA) tail the feed

### Implementation

```python
# Producer: emit event
def emit_event(event: Event):
    with open(".work/events.jsonl", "a") as f:
        f.write(event.model_dump_json() + "\n")

    # Also append to real-time feed
    with open(".work/feed.jsonl", "a") as f:
        f.write(event.model_dump_json() + "\n")

# Consumer: watch feed
async def watch_feed():
    with open(".work/feed.jsonl", "r") as f:
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
- All events in a sprint share `correlation_id`
- All events in a process execution share `correlation_id`
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

This creates an audit trail: `work_order.created` → `work_order.assigned` → `agent.started` → `agent.working` → ...

---

## CLI Commands for Events

```bash
# List recent events
wo events list                    # Last 50 events
wo events list --type=work_order.*  # Filter by type
wo events list --actor=ceo        # Filter by actor

# Tail events (real-time)
wo events tail                    # Watch feed.jsonl

# Replay events for a work order
wo events replay <wo-id>          # Show event history

# Event statistics
wo events stats                   # Counts by type
wo events stats --since=1d        # Last 24 hours

# Export events
wo events export --since=2026-01-01 > events-january.jsonl
```

---

## Integration with Agents

### Assignments as Events

When an assignment is set, emit `work_order.assigned`:

```python
def set_assignment(agent: str, wo_id: str):
    # Write assignment file
    assignment_path = Path(f".work/.assignment-{agent}")
    assignment_path.write_text(f"work_order:{wo_id}")

    # Emit event
    emit_event(Event(
        event_type="work_order.assigned",
        actor="system",
        data={
            "work_order_id": wo_id,
            "agent": agent,
            "assignment_path": str(assignment_path)
        }
    ))
```

### Assignment Principle Compliance Events

Track assignment principle compliance:

```json
{
  "event_type": "agent.assignment_checked",
  "data": {
    "agent": "project-a/workers/slot0",
    "assignment_found": true,
    "response_ms": 150,
    "action": "execute_immediately"
  }
}
```

### Supervisor Patrol Events

```json
{
  "event_type": "patrol.completed",
  "actor": "project-a/supervisor",
  "data": {
    "workers_checked": 3,
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
completed = count_events(type="work_order.status_changed", to="closed")
total = count_events(type="work_order.created")
rate = completed / total

# Assignment principle compliance
assignment_events = get_events(type="agent.assignment_checked")
fast_responses = [e for e in assignment_events if e.data["response_ms"] < 30000]
compliance = len(fast_responses) / len(assignment_events)

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
- [WORKFLOWS.md](./WORKFLOWS.md) - Process lifecycle
