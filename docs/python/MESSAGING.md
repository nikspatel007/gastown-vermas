# VerMAS Messaging

> Communication patterns, mail protocol, and hooks

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
│   │   Stored in .beads/messages.jsonl                                   │  │
│   │   Supports: notifications, requests, handoffs                       │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                           HOOKS                                     │  │
│   │                                                                     │  │
│   │   Work assignment mechanism                                         │  │
│   │   Stored in .beads/.hook-{agent}                                    │  │
│   │   Agent checks on startup, executes immediately (GUPP)              │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                          BEADS                                      │  │
│   │                                                                     │  │
│   │   Shared work state                                                 │  │
│   │   Stored in .beads/issues.jsonl                                     │  │
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
| `POLECAT_DONE` | Polecat | Witness | Work completed |
| `MERGE_READY` | Witness | Refinery | Ready for merge |
| `MERGED` | Refinery | Author | Successfully merged |
| `REWORK_REQUEST` | Refinery | Author | Changes needed |
| `NUDGE` | Witness | Polecat | Wake up idle worker |
| `WITNESS_PING` | Deacon | Witness | Health check |
| `HELP` | Any | Witness/Mayor | Request assistance |
| `HANDOFF` | Any | Self/Next | Session continuity |
| `ESCALATION` | Any | Deacon/Mayor | Problem report |

### Message Flow Diagrams

**Happy Path: Work Completion**

```
Polecat                 Witness                 Refinery
   │                       │                       │
   │ POLECAT_DONE          │                       │
   │──────────────────────▶│                       │
   │                       │                       │
   │                       │ MERGE_READY           │
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
Polecat                 Witness                 Refinery
   │                       │                       │
   │ POLECAT_DONE          │                       │
   │──────────────────────▶│                       │
   │                       │                       │
   │                       │ MERGE_READY           │
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

**Escalation Path: Stuck Polecat**

```
Polecat                 Witness                 Deacon                 Mayor
   │                       │                       │                      │
   │ (idle >5min)          │                       │                      │
   │                       │                       │                      │
   │       NUDGE           │                       │                      │
   │◀──────────────────────│                       │                      │
   │                       │                       │                      │
   │ (still idle >15min)   │                       │                      │
   │                       │                       │                      │
   │ (killed by Witness)   │                       │                      │
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
| `from` | string | Sender BD_ACTOR |
| `to` | string | Recipient BD_ACTOR |
| `subject` | string | Message subject |
| `body` | string | Message content |
| `timestamp` | datetime | When sent |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `type` | enum | Message type (POLECAT_DONE, etc.) |
| `priority` | enum | urgent/normal/low |
| `read_at` | datetime | When recipient read it |
| `metadata` | object | Additional context |

### Addressing

```
{rig}/{role}/{name}    # Full address
{rig}/{role}           # Role address (any agent in role)
{role}                 # Town-level role
```

**Examples:**
- `gastown/polecats/slot0` - Specific polecat
- `gastown/witness` - Witness for gastown rig
- `mayor` - Town-level Mayor
- `deacon` - Town-level Deacon

---

## Hook System

### What is a Hook?

A hook is where work "hangs" waiting for an agent. It's the assignment mechanism.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              HOOK MECHANISM                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Mayor/Human                                                               │
│       │                                                                     │
│       │ gt sling bead-123 gastown                                          │
│       │                                                                     │
│       ▼                                                                     │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                                                                     │  │
│   │   .beads/.hook-gastown-polecats-slot0                              │  │
│   │   ─────────────────────────────────────                            │  │
│   │   bead:bead-123                                                    │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   Polecat starts                                                            │
│       │                                                                     │
│       │ gt hook (check)                                                    │
│       │                                                                     │
│       ▼                                                                     │
│   "HOOKED: bead-123"                                                        │
│       │                                                                     │
│       │ GUPP: Execute immediately!                                         │
│       │                                                                     │
│       ▼                                                                     │
│   (start working on bead-123)                                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Hook File Format

Simple text file: `{type}:{id}`

```
bead:gt-abc12       # Bead hooked
mail:msg-xyz99      # Mail message hooked (for handoffs)
mol:mol-123abc      # Molecule hooked
```

### Hook Types

| Type | Use Case |
|------|----------|
| `bead` | Normal work assignment |
| `mail` | Handoff instructions |
| `mol` | Workflow continuation |

### GUPP (Propulsion Principle)

> **If your hook has work, RUN IT.**

```
Agent starts
    │
    ▼
Check hook
    │
    ├── Hook has work ──────▶ EXECUTE IMMEDIATELY
    │                        No confirmation
    │                        No questions
    │                        No waiting
    │
    └── Hook empty ──────────▶ Check mail
                              Then await instructions
```

---

## Handoff Protocol

For session continuity across restarts.

### Creating a Handoff

```
gt mail send mayor -s "🤝 HANDOFF: Context for next session" -m "Details..."
gt hook attach {mail-id}
```

### Receiving a Handoff

```
Agent starts
    │
    ▼
Check hook
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
3. GUPP overrides priority (hooked work runs first)

---

## Communication Patterns

### Request-Response

```
A: Request
B: Response

Example:
Witness sends NUDGE
Polecat resumes or sends HELP
```

### Notification (Fire-and-Forget)

```
A: Notification
(no response expected)

Example:
Polecat sends POLECAT_DONE
Witness processes, Polecat doesn't wait
```

### Cascade

```
A → B → C

Example:
POLECAT_DONE → MERGE_READY → MERGED
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

### Hook Events

| Event Type | When Emitted | Data |
|------------|--------------|------|
| `hook.set` | Work assigned to hook | agent, ref_type, ref_id |
| `hook.cleared` | Hook emptied | agent, previous_ref |
| `hook.checked` | Agent checked hook (GUPP) | agent, found, response_ms |

### Example Event Stream

```
mail.sent       → {from: "polecat", to: "witness", msg: "POLECAT_DONE"}
mail.delivered  → {to: "witness", msg_id: "..."}
mail.read       → {reader: "witness", msg_id: "..."}
mail.sent       → {from: "witness", to: "refinery", msg: "MERGE_READY"}
```

This enables precise timing analysis and debugging of communication flows.

---

## Logging and Audit

### Message Archive

All messages stored in `.beads/messages.jsonl`:
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
| Hook lost | Re-sling the work |
| Mail stuck | Clear and retry |

---

## Best Practices

### Message Content

1. **Be specific** - Include bead IDs, slot names
2. **Include context** - What led to this message
3. **Action oriented** - What should recipient do
4. **Structured** - Easy to parse programmatically

### Hook Usage

1. **One hook per agent** - Single assignment at a time
2. **Clear when done** - Don't leave stale hooks
3. **Check on startup** - GUPP compliance
4. **Persist through crashes** - File-based, survives restart

### Handoffs

1. **Use emoji** - 🤝 HANDOFF in subject
2. **Be thorough** - Include all context
3. **Hook it** - So next session finds it
4. **Time-bound** - Don't leave handoffs indefinitely

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [AGENTS.md](./AGENTS.md) - Agent roles
- [HOOKS.md](./HOOKS.md) - Claude Code integration and git worktrees
- [WORKFLOWS.md](./WORKFLOWS.md) - Molecule system
- [EVENTS.md](./EVENTS.md) - Event sourcing and change feeds
- [SCHEMAS.md](./SCHEMAS.md) - Message data specifications
- [CLI.md](./CLI.md) - Mail command reference
- [EVALUATION.md](./EVALUATION.md) - How to evaluate the system
