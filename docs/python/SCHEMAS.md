# VerMAS Data Schemas

> JSONL and TOML specifications for all data types

## Overview

VerMAS uses simple, proven data formats:
- **JSONL** for append-only logs and records
- **TOML** for configuration and templates
- **Plain text** for assignments and simple state

All schemas are designed to be:
- Human-readable (grep, cat, jq)
- Git-friendly (line-based diffs)
- Language-agnostic (Go and Python interoperable)

---

## File Layout

```
.work/
├── events.jsonl           # Event log (source of truth)
├── work_orders.jsonl      # Work order records (projection)
├── messages.jsonl         # Mail messages (projection)
├── routes.jsonl           # Prefix routing
├── feed.jsonl             # Real-time change feed
├── .assignment-{agent}    # Assignment files (plain text)
├── templates/
│   └── *.toml             # Workflow templates
├── processes/
│   └── *.json             # Active processes
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
  "event_type": "work_order.created",
  "timestamp": "2026-01-06T12:00:00.000Z",
  "actor": "project-a/teams/frontend",
  "correlation_id": "sprint-xyz789",
  "caused_by": "evt-previous123",
  "data": {}
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `event_id` | string | Yes | Unique ID: `evt-{random12}` |
| `event_type` | string | Yes | Namespaced type (see below) |
| `timestamp` | datetime | Yes | ISO 8601 with milliseconds |
| `actor` | string | Yes | AGENT_ID of emitter |
| `correlation_id` | string | No | Links related events |
| `caused_by` | string | No | Event that triggered this |
| `data` | object | Yes | Event-specific payload |

### Event Types

#### Work Order Events

```json
// work_order.created
{
  "event_type": "work_order.created",
  "data": {
    "work_order_id": "wo-abc123",
    "title": "Implement feature X",
    "description": "Full description...",
    "issue_type": "feature",
    "priority": 2,
    "created_by": "ceo"
  }
}

// work_order.updated
{
  "event_type": "work_order.updated",
  "data": {
    "work_order_id": "wo-abc123",
    "changes": {
      "title": "New title",
      "priority": 1
    }
  }
}

// work_order.status_changed
{
  "event_type": "work_order.status_changed",
  "data": {
    "work_order_id": "wo-abc123",
    "from_status": "open",
    "to_status": "in_progress"
  }
}

// work_order.assigned
{
  "event_type": "work_order.assigned",
  "data": {
    "work_order_id": "wo-abc123",
    "agent": "project-a/workers/slot0",
    "assignment_path": ".work/.assignment-project-a-workers-slot0"
  }
}

// work_order.closed
{
  "event_type": "work_order.closed",
  "data": {
    "work_order_id": "wo-abc123",
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
    "from": "project-a/workers/slot0",
    "to": "project-a/supervisor",
    "subject": "WORKER_DONE",
    "message_type": "WORKER_DONE"
  }
}

// mail.delivered
{
  "event_type": "mail.delivered",
  "data": {
    "message_id": "msg-xyz789",
    "recipient": "project-a/supervisor"
  }
}

// mail.read
{
  "event_type": "mail.read",
  "data": {
    "message_id": "msg-xyz789",
    "reader": "project-a/supervisor",
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
    "session_name": "worker-project-a-slot0",
    "worktree": "/path/to/worktree",
    "profile": "worker",
    "assigned_work_order": "wo-abc123"
  }
}

// agent.assignment_checked (Assignment Principle compliance)
{
  "event_type": "agent.assignment_checked",
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
    "work_order_id": "wo-abc123"
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

#### Process Events

```json
// process.created
{
  "event_type": "process.created",
  "data": {
    "process_id": "proc-abc123",
    "template": "worker-execute",
    "work_order_id": "wo-xyz789"
  }
}

// process.step_started
{
  "event_type": "process.step_started",
  "data": {
    "process_id": "proc-abc123",
    "step_id": "implement"
  }
}

// process.step_completed
{
  "event_type": "process.step_completed",
  "data": {
    "process_id": "proc-abc123",
    "step_id": "implement",
    "status": "completed"
  }
}

// process.completed
{
  "event_type": "process.completed",
  "data": {
    "process_id": "proc-abc123",
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
    "work_order_id": "wo-abc123",
    "process_id": "proc-verify-xyz"
  }
}

// verify.verdict
{
  "event_type": "verify.verdict",
  "data": {
    "work_order_id": "wo-abc123",
    "verdict": "PASS",
    "criteria_passed": 5,
    "criteria_failed": 0,
    "reasoning": "All acceptance criteria met"
  }
}
```

---

## Work Order Schema

Work orders stored in `work_orders.jsonl`:

```json
{
  "id": "wo-abc123",
  "title": "Implement feature X",
  "description": "Full description with markdown...",
  "status": "open",
  "priority": 2,
  "issue_type": "feature",
  "created_at": "2026-01-06T10:00:00.000Z",
  "updated_at": "2026-01-06T12:00:00.000Z",
  "created_by": "ceo",
  "assignee": "project-a/workers/slot0",
  "dependencies": [
    {
      "depends_on_id": "wo-dep456",
      "type": "blocks",
      "created_at": "2026-01-06T10:05:00.000Z",
      "created_by": "ceo"
    }
  ],
  "labels": ["backend", "api"],
  "metadata": {}
}
```

### Field Definitions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique ID: `wo-{hash}` or hierarchical `wo-{hash}.{n}` |
| `title` | string | Yes | Short title (max 100 chars) |
| `description` | string | No | Markdown description |
| `status` | enum | Yes | `open`, `in_progress`, `assigned`, `closed` |
| `priority` | int | Yes | 0-4 (0=P0 critical, 4=P4 backlog) |
| `issue_type` | enum | Yes | `task`, `bug`, `feature`, `epic`, `merge-request` |
| `created_at` | datetime | Yes | ISO 8601 |
| `updated_at` | datetime | Yes | ISO 8601 |
| `created_by` | string | Yes | AGENT_ID |
| `assignee` | string | No | AGENT_ID |
| `dependencies` | array | No | Dependency objects |
| `labels` | array | No | String labels |
| `metadata` | object | No | Custom key-value pairs |

### Status Values

| Status | Meaning |
|--------|---------|
| `open` | Ready to be worked |
| `in_progress` | Being worked on |
| `assigned` | Assigned to agent |
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
  "from": "project-a/workers/slot0",
  "to": "project-a/supervisor",
  "subject": "WORKER_DONE",
  "body": "Work completed for work order wo-abc123",
  "message_type": "WORKER_DONE",
  "timestamp": "2026-01-06T12:00:00.000Z",
  "read_at": null,
  "priority": "normal",
  "metadata": {
    "work_order_id": "wo-abc123"
  }
}
```

### Field Definitions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique ID: `msg-{random}` |
| `from` | string | Yes | Sender AGENT_ID |
| `to` | string | Yes | Recipient AGENT_ID |
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
| `WORKER_DONE` | Worker | Supervisor | Work completed |
| `READY_FOR_QA` | Supervisor | QA | Ready for merge |
| `MERGED` | QA | Author | Successfully merged |
| `REWORK_REQUEST` | QA | Author | Changes needed |
| `NUDGE` | Supervisor | Worker | Wake up idle |
| `HELP` | Any | Supervisor/CEO | Request assistance |
| `HANDOFF` | Any | Self | Session continuity |
| `ESCALATION` | Any | CEO | Problem report |

---

## Assignment Schema

Assignments are plain text files: `.work/.assignment-{agent-id}`

```
{type}:{reference_id}
```

### Examples

```
work_order:wo-abc123
mail:msg-xyz789
process:proc-verify-def
```

### Assignment Types

| Type | Reference | Purpose |
|------|-----------|---------|
| `work_order` | Work order ID | Work assignment |
| `mail` | Message ID | Handoff instructions |
| `process` | Process ID | Workflow continuation |

### Agent ID in Filename

Agent ID with `/` replaced by `-`:
- `project-a/workers/slot0` → `.assignment-project-a-workers-slot0`
- `ceo` → `.assignment-ceo`

---

## Template Schema (TOML)

Templates define workflow definitions in `.work/templates/*.toml`:

```toml
# worker-execute.toml

template = "worker-execute"
description = "Execute assigned work order to completion"
version = 1

[[steps]]
id = "understand"
title = "Understand the task"
description = """
Read the work order details and plan your approach.

```bash
wo show $WORK_ORDER_ID
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
co worker done
```
"""
```

### Template Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `template` | string | Yes | Unique template name |
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

## Process Schema (JSON)

Active processes stored in `.work/processes/{id}.json`:

```json
{
  "id": "proc-abc123",
  "template": "worker-execute",
  "work_order_id": "wo-xyz789",
  "state": "active",
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

### Process Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique process ID |
| `template` | string | Yes | Source template name |
| `work_order_id` | string | No | Attached work order (if any) |
| `state` | enum | Yes | `ready`, `active`, `archive` |
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

Routes map ID prefixes to work locations in `routes.jsonl`:

```json
{"prefix": "wo", "path": "/path/to/project-a/.work", "factory": "project-a"}
{"prefix": "hq", "path": "/path/to/company/.work", "factory": null}
{"prefix": "pa", "path": "/path/to/project-a/.work", "factory": "project-a"}
```

### Route Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `prefix` | string | Yes | ID prefix (e.g., `wo`) |
| `path` | string | Yes | Absolute path to .work/ |
| `factory` | string | No | Factory name (null for company-level) |

---

## Validation

### Python (Pydantic)

```python
from pydantic import BaseModel, Field
from datetime import datetime
from enum import Enum
from typing import Optional

class WorkOrderStatus(str, Enum):
    open = "open"
    in_progress = "in_progress"
    assigned = "assigned"
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

class WorkOrder(BaseModel):
    id: str
    title: str
    description: str = ""
    status: WorkOrderStatus = WorkOrderStatus.open
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
- [HOOKS.md](./HOOKS.md) - Assignment file format
- [WORKFLOWS.md](./WORKFLOWS.md) - Template and process schemas
- [GO-VS-PYTHON.md](./GO-VS-PYTHON.md) - Implementation patterns
