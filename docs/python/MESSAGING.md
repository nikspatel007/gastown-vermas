# VerMAS Messaging

> Communication patterns, mail protocol, and assignments

## Communication Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        COMMUNICATION CHANNELS                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                            MAIL                                     │  │
│   │                                                                     │  │
│   │   Async messages between agents                                     │  │
│   │   Stored in .work/messages.jsonl                                    │  │
│   │   Supports: notifications, requests, handoffs                       │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                       ASSIGNMENTS                                   │  │
│   │                                                                     │  │
│   │   Work assignment mechanism                                         │  │
│   │   Stored in .work/.assignment-{agent}                               │  │
│   │   Agent checks on startup, executes immediately                     │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                      WORK ORDERS                                    │  │
│   │                                                                     │  │
│   │   Shared work state                                                 │  │
│   │   Stored in .work/work_orders.jsonl                                 │  │
│   │   All agents can read; owners can update                           │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Mail Protocol

### Message Types

| Type | From | To | Purpose |
|------|------|----|---------|
| `WORKER_DONE` | Worker | Supervisor | Work completed |
| `READY_FOR_QA` | Supervisor | QA | Ready for merge |
| `MERGED` | QA | Author | Successfully merged |
| `REWORK_REQUEST` | QA | Author | Changes needed |
| `NUDGE` | Supervisor | Worker | Wake up idle worker |
| `SUPERVISOR_PING` | Operations | Supervisor | Health check |
| `HELP` | Any | Supervisor/CEO | Request assistance |
| `HANDOFF` | Any | Self/Next | Session continuity |
| `ESCALATION` | Any | Operations/CEO | Problem report |

### Message Flow Diagrams

**Happy Path: Work Completion**

```
Worker                 Supervisor               QA
   │                       │                       │
   │ WORKER_DONE           │                       │
   │──────────────────────▶│                       │
   │                       │                       │
   │                       │ READY_FOR_QA          │
   │                       │──────────────────────▶│
   │                       │                       │
   │                       │                       │ (run tests, verify)
   │                       │                       │
   │                       │       MERGED          │
   │◀──────────────────────│◀──────────────────────│
   │                       │                       │
```

**Failure Path: Rework Required**

```
Worker                 Supervisor               QA
   │                       │                       │
   │ WORKER_DONE           │                       │
   │──────────────────────▶│                       │
   │                       │                       │
   │                       │ READY_FOR_QA          │
   │                       │──────────────────────▶│
   │                       │                       │
   │                       │                       │ (tests fail)
   │                       │                       │
   │   REWORK_REQUEST      │                       │
   │◀──────────────────────│◀──────────────────────│
   │                       │                       │
   │ (fix issues, retry)   │                       │
   │                       │                       │
```

**Escalation Path: Stuck Worker**

```
Worker                 Supervisor             Operations               CEO
   │                       │                       │                      │
   │ (idle >5min)          │                       │                      │
   │                       │                       │                      │
   │       NUDGE           │                       │                      │
   │◀──────────────────────│                       │                      │
   │                       │                       │                      │
   │ (still idle >15min)   │                       │                      │
   │                       │                       │                      │
   │ (killed by Supervisor)│                       │                      │
   │                       │                       │                      │
   │                       │ ESCALATION (if >30min)│                      │
   │                       │──────────────────────▶│                      │
   │                       │                       │                      │
   │                       │                       │ ESCALATION           │
   │                       │                       │─────────────────────▶│
   │                       │                       │                      │
```

---

## Message Format

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique message ID |
| `from` | string | Sender AGENT_ID |
| `to` | string | Recipient AGENT_ID |
| `subject` | string | Message subject |
| `body` | string | Message content |
| `timestamp` | datetime | When sent |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `type` | enum | Message type (WORKER_DONE, etc.) |
| `priority` | enum | urgent/normal/low |
| `read_at` | datetime | When recipient read it |
| `metadata` | object | Additional context |

### Addressing

```
{factory}/{role}/{name}    # Full address
{factory}/{role}           # Role address (any agent in role)
{role}                     # Company-level role
```

**Examples:**
- `project-a/workers/slot0` - Specific worker
- `project-a/supervisor` - Supervisor for project-a factory
- `ceo` - Company-level CEO
- `operations` - Company-level Operations

---

## Assignment System

### What is an Assignment?

An assignment is where work "waits" for an agent. It's the work assignment mechanism.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          ASSIGNMENT MECHANISM                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   CEO/Human                                                                 │
│       │                                                                     │
│       │ co dispatch wo-123 project-a                                       │
│       │                                                                     │
│       ▼                                                                     │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                                                                     │  │
│   │   .work/.assignment-project-a-workers-slot0                        │  │
│   │   ─────────────────────────────────────────                        │  │
│   │   work_order:wo-123                                                │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   Worker starts                                                             │
│       │                                                                     │
│       │ co assignment (check)                                              │
│       │                                                                     │
│       ▼                                                                     │
│   "ASSIGNED: wo-123"                                                        │
│       │                                                                     │
│       │ Assignment Principle: Execute immediately!                         │
│       │                                                                     │
│       ▼                                                                     │
│   (start working on wo-123)                                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Assignment File Format

Simple text file: `{type}:{id}`

```
work_order:wo-abc12     # Work order assigned
mail:msg-xyz99          # Mail message assigned (for handoffs)
process:proc-123abc     # Process assigned
```

### Assignment Types

| Type | Use Case |
|------|----------|
| `work_order` | Normal work assignment |
| `mail` | Handoff instructions |
| `process` | Workflow continuation |

### Assignment Principle

> **If you have an assignment, EXECUTE IT.**

```
Agent starts
    │
    ▼
Check assignment
    │
    ├── Assignment has work ──────▶ EXECUTE IMMEDIATELY
    │                               No confirmation
    │                               No questions
    │                               No waiting
    │
    └── Assignment empty ──────────▶ Check mail
                                     Then await instructions
```

---

## Handoff Protocol

For session continuity across restarts.

### Creating a Handoff

```
co send ceo -s "🤝 HANDOFF: Context for next session" -m "Details..."
co assignment attach {mail-id}
```

### Receiving a Handoff

```
Agent starts
    │
    ▼
Check assignment
    │
    ▼
Find mail:msg-123
    │
    ▼
Read mail content
    │
    ▼
Execute instructions from mail body
```

### Handoff Subject Convention

```
🤝 HANDOFF: Brief description
```

The 🤝 emoji indicates this is a handoff message for session continuity.

---

## Priority System

### Priority Levels

| Priority | Meaning | Processing |
|----------|---------|------------|
| `urgent` | Immediate attention | Process first |
| `normal` | Standard work | Process in order |
| `low` | Can wait | Process when idle |

### Priority Rules

1. Urgent messages processed before checking regular queue
2. Within same priority, process oldest first
3. Assignment Principle overrides priority (assigned work runs first)

---

## Communication Patterns

### Request-Response

```
A: Request
B: Response

Example:
Supervisor sends NUDGE
Worker resumes or sends HELP
```

### Notification (Fire-and-Forget)

```
A: Notification
(no response expected)

Example:
Worker sends WORKER_DONE
Supervisor processes, Worker doesn't wait
```

### Cascade

```
A → B → C

Example:
WORKER_DONE → READY_FOR_QA → MERGED
```

### Broadcast (not implemented)

```
A → B, C, D

Would require distribution lists.
Future enhancement.
```

---

## Events Emitted

All mail operations emit events to the event log. See [EVENTS.md](./EVENTS.md) for full event sourcing documentation.

### Mail Events

| Event Type | When Emitted | Data |
|------------|--------------|------|
| `mail.sent` | Message dispatched | from, to, subject, message_id |
| `mail.delivered` | Message written to inbox | message_id, recipient |
| `mail.read` | Recipient opened message | message_id, reader, read_at |
| `mail.archived` | Message moved to archive | message_id |

### Assignment Events

| Event Type | When Emitted | Data |
|------------|--------------|------|
| `assignment.set` | Work assigned | agent, ref_type, ref_id |
| `assignment.cleared` | Assignment emptied | agent, previous_ref |
| `assignment.checked` | Agent checked assignment | agent, found, response_ms |

### Example Event Stream

```
mail.sent       → {from: "worker", to: "supervisor", msg: "WORKER_DONE"}
mail.delivered  → {to: "supervisor", msg_id: "..."}
mail.read       → {reader: "supervisor", msg_id: "..."}
mail.sent       → {from: "supervisor", to: "qa", msg: "READY_FOR_QA"}
```

This enables precise timing analysis and debugging of communication flows.

---

## Logging and Audit

### Message Archive

All messages stored in `.work/messages.jsonl`:
- Full message content
- Sender and recipient
- Timestamps
- Read status

**Note:** `messages.jsonl` is a projection of `mail.*` events. The event log is the source of truth.

### Audit Trail

Can reconstruct:
- What messages were sent
- Who sent them
- When they were read
- Full communication history

Event-based audit provides:
- Precise timing (millisecond accuracy)
- Causation chains (which event caused which)
- Correlation IDs (link messages to workflows)

### Privacy Considerations

All messages are stored in plain text in git.
Assume all communication is visible to:
- All agents in the system
- Anyone with repo access

---

## Error Handling

### Undeliverable Messages

If recipient doesn't exist:
1. Message stored with `undeliverable` flag
2. No automatic retry
3. Sender not notified (async)

### Lost Messages

If message file corrupted:
1. Agent continues without it
2. Work may need manual re-dispatch
3. Check logs for what was lost

### Recovery Strategies

| Scenario | Recovery |
|----------|----------|
| Message lost | Re-send or check logs |
| Assignment lost | Re-dispatch the work |
| Mail stuck | Clear and retry |

---

## Best Practices

### Message Content

1. **Be specific** - Include work order IDs, slot names
2. **Include context** - What led to this message
3. **Action oriented** - What should recipient do
4. **Structured** - Easy to parse programmatically

### Assignment Usage

1. **One assignment per agent** - Single work item at a time
2. **Clear when done** - Don't leave stale assignments
3. **Check on startup** - Assignment Principle compliance
4. **Persist through crashes** - File-based, survives restart

### Handoffs

1. **Use emoji** - 🤝 HANDOFF in subject
2. **Be thorough** - Include all context
3. **Assign it** - So next session finds it
4. **Time-bound** - Don't leave handoffs indefinitely

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [AGENTS.md](./AGENTS.md) - Agent roles
- [HOOKS.md](./HOOKS.md) - Claude Code integration and git worktrees
- [WORKFLOWS.md](./WORKFLOWS.md) - Process system
- [EVENTS.md](./EVENTS.md) - Event sourcing and change feeds
- [SCHEMAS.md](./SCHEMAS.md) - Message data specifications
- [CLI.md](./CLI.md) - Mail command reference
- [EVALUATION.md](./EVALUATION.md) - How to evaluate the system
