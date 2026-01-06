# VerMAS Data Schemas

> JSONL and TOML specifications for all data types

## Overview

VerMAS uses simple, proven data formats:
- **JSONL** for append-only logs and records
- **TOML** for configuration and templates
- **Plain text** for hooks and simple state

All schemas are designed to be:
- Human-readable (grep, cat, jq)
- Git-friendly (line-based diffs)
- Language-agnostic (Go and Python interoperable)

---

## File Layout

```
.beads/
├── events.jsonl           # Event log (source of truth)
├── issues.jsonl           # Bead records (projection)
├── messages.jsonl         # Mail messages (projection)
├── routes.jsonl           # Prefix routing
├── feed.jsonl             # Real-time change feed
├── .hook-{agent}          # Hook files (plain text)
├── formulas/
│   └── *.toml             # Workflow templates
├── mols/
│   └── *.json             # Active molecules
└── evidence/
    └── *.json             # Verification evidence
```

---

## Event Schema

Events are the source of truth. All other data is derived from events.

### Base Event

```json
{
  "event_id": "evt-abc123def456",
  "event_type": "bead.created",
  "timestamp": "2026-01-06T12:00:00.000Z",
  "actor": "gastown/crew/frontend",
  "correlation_id": "conv-xyz789",
  "caused_by": "evt-previous123",
  "data": {}
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `event_id` | string | Yes | Unique ID: `evt-{random12}` |
| `event_type` | string | Yes | Namespaced type (see below) |
| `timestamp` | datetime | Yes | ISO 8601 with milliseconds |
| `actor` | string | Yes | BD_ACTOR of emitter |
| `correlation_id` | string | No | Links related events |
| `caused_by` | string | No | Event that triggered this |
| `data` | object | Yes | Event-specific payload |

### Event Types

#### Bead Events

```json
// bead.created
{
  "event_type": "bead.created",
  "data": {
    "bead_id": "gt-abc123",
    "title": "Implement feature X",
    "description": "Full description...",
    "issue_type": "feature",
    "priority": 2,
    "created_by": "mayor"
  }
}

// bead.updated
{
  "event_type": "bead.updated",
  "data": {
    "bead_id": "gt-abc123",
    "changes": {
      "title": "New title",
      "priority": 1
    }
  }
}

// bead.status_changed
{
  "event_type": "bead.status_changed",
  "data": {
    "bead_id": "gt-abc123",
    "from_status": "open",
    "to_status": "in_progress"
  }
}

// bead.hooked
{
  "event_type": "bead.hooked",
  "data": {
    "bead_id": "gt-abc123",
    "agent": "gastown/polecats/slot0",
    "hook_path": ".beads/.hook-gastown-polecats-slot0"
  }
}

// bead.closed
{
  "event_type": "bead.closed",
  "data": {
    "bead_id": "gt-abc123",
    "reason": "Completed successfully"
  }
}
```

#### Mail Events

```json
// mail.sent
{
  "event_type": "mail.sent",
  "data": {
    "message_id": "msg-xyz789",
    "from": "gastown/polecats/slot0",
    "to": "gastown/witness",
    "subject": "POLECAT_DONE",
    "message_type": "POLECAT_DONE"
  }
}

// mail.delivered
{
  "event_type": "mail.delivered",
  "data": {
    "message_id": "msg-xyz789",
    "recipient": "gastown/witness"
  }
}

// mail.read
{
  "event_type": "mail.read",
  "data": {
    "message_id": "msg-xyz789",
    "reader": "gastown/witness",
    "read_at": "2026-01-06T12:05:00.000Z"
  }
}
```

#### Agent Events

```json
// agent.started
{
  "event_type": "agent.started",
  "data": {
    "session_name": "polecat-gastown-slot0",
    "worktree": "/path/to/worktree",
    "profile": "polecat",
    "hooked_bead": "gt-abc123"
  }
}

// agent.hook_checked (GUPP compliance)
{
  "event_type": "agent.hook_checked",
  "data": {
    "found": true,
    "response_ms": 150,
    "action": "execute_immediately"
  }
}

// agent.working
{
  "event_type": "agent.working",
  "data": {
    "bead_id": "gt-abc123"
  }
}

// agent.idle
{
  "event_type": "agent.idle",
  "data": {
    "idle_since": "2026-01-06T12:10:00.000Z",
    "last_activity": "tool_use"
  }
}

// agent.stopped
{
  "event_type": "agent.stopped",
  "data": {
    "reason": "completed",
    "exit_code": 0
  }
}
```

#### Workflow Events

```json
// mol.created
{
  "event_type": "mol.created",
  "data": {
    "mol_id": "mol-abc123",
    "formula": "mol-polecat-work",
    "bead_id": "gt-xyz789"
  }
}

// mol.step_started
{
  "event_type": "mol.step_started",
  "data": {
    "mol_id": "mol-abc123",
    "step_id": "implement"
  }
}

// mol.step_completed
{
  "event_type": "mol.step_completed",
  "data": {
    "mol_id": "mol-abc123",
    "step_id": "implement",
    "status": "completed"
  }
}

// mol.completed
{
  "event_type": "mol.completed",
  "data": {
    "mol_id": "mol-abc123",
    "summary": "All steps completed successfully"
  }
}
```

#### Verification Events

```json
// verify.started
{
  "event_type": "verify.started",
  "data": {
    "bead_id": "gt-abc123",
    "mol_id": "mol-verify-xyz"
  }
}

// verify.verdict
{
  "event_type": "verify.verdict",
  "data": {
    "bead_id": "gt-abc123",
    "verdict": "PASS",
    "criteria_passed": 5,
    "criteria_failed": 0,
    "reasoning": "All acceptance criteria met"
  }
}
```

---

## Bead Schema

Beads (issues) stored in `issues.jsonl`:

```json
{
  "id": "gt-abc123",
  "title": "Implement feature X",
  "description": "Full description with markdown...",
  "status": "open",
  "priority": 2,
  "issue_type": "feature",
  "created_at": "2026-01-06T10:00:00.000Z",
  "updated_at": "2026-01-06T12:00:00.000Z",
  "created_by": "mayor",
  "assignee": "gastown/polecats/slot0",
  "dependencies": [
    {
      "depends_on_id": "gt-dep456",
      "type": "blocks",
      "created_at": "2026-01-06T10:05:00.000Z",
      "created_by": "mayor"
    }
  ],
  "labels": ["backend", "api"],
  "metadata": {}
}
```

### Field Definitions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique ID: `{prefix}-{hash}` or hierarchical `{prefix}-{hash}.{n}` |
| `title` | string | Yes | Short title (max 100 chars) |
| `description` | string | No | Markdown description |
| `status` | enum | Yes | `open`, `in_progress`, `hooked`, `closed` |
| `priority` | int | Yes | 0-4 (0=P0 critical, 4=P4 backlog) |
| `issue_type` | enum | Yes | `task`, `bug`, `feature`, `epic`, `merge-request` |
| `created_at` | datetime | Yes | ISO 8601 |
| `updated_at` | datetime | Yes | ISO 8601 |
| `created_by` | string | Yes | BD_ACTOR |
| `assignee` | string | No | BD_ACTOR |
| `dependencies` | array | No | Dependency objects |
| `labels` | array | No | String labels |
| `metadata` | object | No | Custom key-value pairs |

### Status Values

| Status | Meaning |
|--------|---------|
| `open` | Ready to be worked |
| `in_progress` | Being worked on |
| `hooked` | Assigned to agent hook |
| `closed` | Completed |

### Priority Values

| Value | Label | Meaning |
|-------|-------|---------|
| 0 | P0 | Critical, drop everything |
| 1 | P1 | High priority |
| 2 | P2 | Medium (default) |
| 3 | P3 | Low priority |
| 4 | P4 | Backlog |

### Issue Types

| Type | Purpose |
|------|---------|
| `task` | General work item |
| `bug` | Defect to fix |
| `feature` | New functionality |
| `epic` | Large body of work (has children) |
| `merge-request` | Code review request |

---

## Message Schema

Messages stored in `messages.jsonl`:

```json
{
  "id": "msg-xyz789",
  "from": "gastown/polecats/slot0",
  "to": "gastown/witness",
  "subject": "POLECAT_DONE",
  "body": "Work completed for bead gt-abc123",
  "message_type": "POLECAT_DONE",
  "timestamp": "2026-01-06T12:00:00.000Z",
  "read_at": null,
  "priority": "normal",
  "metadata": {
    "bead_id": "gt-abc123"
  }
}
```

### Field Definitions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique ID: `msg-{random}` |
| `from` | string | Yes | Sender BD_ACTOR |
| `to` | string | Yes | Recipient BD_ACTOR |
| `subject` | string | Yes | Message subject |
| `body` | string | Yes | Message content |
| `message_type` | enum | No | Structured type |
| `timestamp` | datetime | Yes | When sent |
| `read_at` | datetime | No | When read (null if unread) |
| `priority` | enum | No | `urgent`, `normal`, `low` |
| `metadata` | object | No | Additional context |

### Message Types

| Type | From | To | Purpose |
|------|------|----|---------|
| `POLECAT_DONE` | Polecat | Witness | Work completed |
| `MERGE_READY` | Witness | Refinery | Ready for merge |
| `MERGED` | Refinery | Author | Successfully merged |
| `REWORK_REQUEST` | Refinery | Author | Changes needed |
| `NUDGE` | Witness | Polecat | Wake up idle |
| `HELP` | Any | Witness/Mayor | Request assistance |
| `HANDOFF` | Any | Self | Session continuity |
| `ESCALATION` | Any | Mayor | Problem report |

---

## Hook Schema

Hooks are plain text files: `.beads/.hook-{agent-id}`

```
{type}:{reference_id}
```

### Examples

```
bead:gt-abc123
mail:msg-xyz789
mol:mol-verify-def
```

### Hook Types

| Type | Reference | Purpose |
|------|-----------|---------|
| `bead` | Bead ID | Work assignment |
| `mail` | Message ID | Handoff instructions |
| `mol` | Molecule ID | Workflow continuation |

### Agent ID in Filename

Agent ID with `/` replaced by `-`:
- `gastown/polecats/slot0` → `.hook-gastown-polecats-slot0`
- `mayor` → `.hook-mayor`

---

## Formula Schema (TOML)

Formulas define workflow templates in `.beads/formulas/*.toml`:

```toml
# mol-polecat-work.formula.toml

formula = "mol-polecat-work"
description = "Execute assigned bead to completion"
version = 1

[[steps]]
id = "understand"
title = "Understand the task"
description = """
Read the bead details and plan your approach.

```bash
bd show $BEAD_ID
```

Understand the requirements before starting.
"""

[[steps]]
id = "implement"
title = "Implement the solution"
needs = ["understand"]
description = """
Write the code to solve the problem.

Follow project conventions and best practices.
"""

[[steps]]
id = "test"
title = "Run tests"
needs = ["implement"]
description = """
Run the test suite to verify your changes.

```bash
pytest
```

Fix any failing tests before proceeding.
"""

[[steps]]
id = "complete"
title = "Complete and signal done"
needs = ["test"]
description = """
Commit, push, and signal completion.

```bash
git add .
git commit -m "Implement feature"
git push
gt polecat done
```
"""
```

### Formula Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `formula` | string | Yes | Unique formula name |
| `description` | string | Yes | What this workflow does |
| `version` | int | Yes | Schema version |
| `steps` | array | Yes | Workflow steps |

### Step Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique step ID |
| `title` | string | Yes | Human-readable name |
| `description` | string | Yes | Markdown instructions |
| `needs` | array | No | Step IDs that must complete first |

---

## Molecule Schema (JSON)

Active molecules stored in `.beads/mols/{id}.json`:

```json
{
  "id": "mol-abc123",
  "formula": "mol-polecat-work",
  "bead_id": "gt-xyz789",
  "state": "liquid",
  "created_at": "2026-01-06T10:00:00.000Z",
  "steps": [
    {
      "id": "understand",
      "status": "completed",
      "started_at": "2026-01-06T10:00:00.000Z",
      "completed_at": "2026-01-06T10:05:00.000Z"
    },
    {
      "id": "implement",
      "status": "in_progress",
      "started_at": "2026-01-06T10:05:00.000Z",
      "completed_at": null
    },
    {
      "id": "test",
      "status": "pending",
      "started_at": null,
      "completed_at": null
    },
    {
      "id": "complete",
      "status": "pending",
      "started_at": null,
      "completed_at": null
    }
  ]
}
```

### Molecule Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique molecule ID |
| `formula` | string | Yes | Source formula name |
| `bead_id` | string | No | Attached bead (if any) |
| `state` | enum | Yes | `solid`, `liquid`, `vapor` |
| `created_at` | datetime | Yes | When instantiated |
| `steps` | array | Yes | Step instances |

### Step Status Values

| Status | Meaning |
|--------|---------|
| `pending` | Not started |
| `blocked` | Waiting on dependencies |
| `ready` | Dependencies met, can start |
| `in_progress` | Currently executing |
| `completed` | Finished successfully |
| `failed` | Execution failed |
| `skipped` | Explicitly skipped |

---

## Route Schema

Routes map ID prefixes to beads locations in `routes.jsonl`:

```json
{"prefix": "gt", "path": "/path/to/gastown/.beads", "rig": "gastown"}
{"prefix": "hq", "path": "/path/to/town/.beads", "rig": null}
{"prefix": "pa", "path": "/path/to/project-a/.beads", "rig": "project-a"}
```

### Route Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `prefix` | string | Yes | ID prefix (e.g., `gt`) |
| `path` | string | Yes | Absolute path to .beads/ |
| `rig` | string | No | Rig name (null for town-level) |

---

## Validation

### Python (Pydantic)

```python
from pydantic import BaseModel, Field
from datetime import datetime
from enum import Enum
from typing import Optional

class BeadStatus(str, Enum):
    open = "open"
    in_progress = "in_progress"
    hooked = "hooked"
    closed = "closed"

class IssueType(str, Enum):
    task = "task"
    bug = "bug"
    feature = "feature"
    epic = "epic"
    merge_request = "merge-request"

class Dependency(BaseModel):
    depends_on_id: str
    type: str = "blocks"
    created_at: datetime
    created_by: str

class Bead(BaseModel):
    id: str
    title: str
    description: str = ""
    status: BeadStatus = BeadStatus.open
    priority: int = Field(ge=0, le=4, default=2)
    issue_type: IssueType
    created_at: datetime
    updated_at: datetime
    created_by: str
    assignee: Optional[str] = None
    dependencies: list[Dependency] = Field(default_factory=list)
    labels: list[str] = Field(default_factory=list)
    metadata: dict = Field(default_factory=dict)
```

---

## See Also

- [EVENTS.md](./EVENTS.md) - Event sourcing patterns
- [CLI.md](./CLI.md) - Commands that use these schemas
- [ARCHITECTURE.md](./ARCHITECTURE.md) - System design
- [HOOKS.md](./HOOKS.md) - Hook file format
- [WORKFLOWS.md](./WORKFLOWS.md) - Formula and molecule schemas
- [GO-VS-PYTHON.md](./GO-VS-PYTHON.md) - Implementation patterns
