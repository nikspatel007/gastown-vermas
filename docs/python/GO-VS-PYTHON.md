# Go vs Python Implementation Comparison

> Same architecture, different languages

## Same Architecture, Different Language

Both Go and Python implementations use **identical architecture**:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SHARED ARCHITECTURE (Both Languages)                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌────────────┐     ┌────────────┐     ┌────────────┐                     │
│   │   Claude   │     │   CLI      │     │   Tmux     │                     │
│   │   Code     │────▶│  Commands  │────▶│  Sessions  │                     │
│   │  (Agent)   │     │  (co/wo)   │     │  (Workers) │                     │
│   └────────────┘     └────────────┘     └────────────┘                     │
│         │                  │                  │                             │
│         └──────────────────┼──────────────────┘                             │
│                            ▼                                                │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                     WORK ORDERS (JSONL Files)                       │  │
│   │    work_orders.jsonl  |  messages.jsonl  |  templates/*.toml        │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   Orchestration: tmux sessions + Claude Code CLI + assignments + CLAUDE.md │
│   Storage: JSONL work orders (git-backed)                                  │
│   LLM: Claude Code CLI (no API costs)                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Both implementations share:**
- Tmux-based agent sessions
- Claude Code CLI for LLM (no API costs)
- JSONL work orders for persistence
- TOML templates for workflows
- Assignment-based work dispatch (Assignment Principle)
- Mail protocol for agent communication
- Git worktrees for isolation

---

## Library Comparison

| Component | Go | Python |
|-----------|-----|--------|
| **CLI Framework** | Cobra | Typer |
| **Tmux Bindings** | go-tmux | libtmux |
| **Data Models** | Go structs | Pydantic |
| **JSON Handling** | encoding/json | pydantic + json |
| **TOML Parsing** | BurntSushi/toml | tomllib (stdlib) |
| **Async/Concurrency** | goroutines | asyncio |
| **Process Control** | os/exec | subprocess/asyncio |
| **File Watching** | fsnotify | watchdog |
| **Testing** | go test | pytest |
| **Type Safety** | Compile-time | Runtime (Pydantic) |
| **Distribution** | Single binary | pip package |

---

## Code Comparison

### CLI Command Definition

**Go (Cobra):**
```go
var assignmentCmd = &cobra.Command{
    Use:   "assignment",
    Short: "Check your assigned work",
    Run: func(cmd *cobra.Command, args []string) {
        actor := os.Getenv("AGENT_ID")
        assignment := LoadAssignment(actor)
        if assignment != nil {
            fmt.Printf("ASSIGNED: %s\n", assignment.WorkOrderID)
        } else {
            fmt.Println("Assignment is empty")
        }
    },
}
```

**Python (Typer):**
```python
@app.command()
def assignment():
    """Check your assigned work."""
    actor = os.environ.get("AGENT_ID")
    assign = Assignment(actor, Path(".work"))
    content = assign.check()
    if content:
        console.print(f"ASSIGNED: {content.ref_id}")
    else:
        console.print("Assignment is empty")
```

### Data Models

**Go (structs):**
```go
type WorkOrder struct {
    ID          string       `json:"id"`
    Title       string       `json:"title"`
    Description string       `json:"description"`
    Status      WorkOrderStatus   `json:"status"`
    Priority    int          `json:"priority"`
    IssueType   IssueType    `json:"issue_type"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
    CreatedBy   string       `json:"created_by"`
    Dependencies []Dependency `json:"dependencies,omitempty"`
}
```

**Python (Pydantic):**
```python
class WorkOrder(BaseModel):
    id: str
    title: str
    description: str
    status: WorkOrderStatus
    priority: int
    issue_type: IssueType
    created_at: datetime
    updated_at: datetime
    created_by: str
    dependencies: List[Dependency] = Field(default_factory=list)
```

### Tmux Session Creation

**Go (go-tmux):**
```go
func CreateSession(name, workDir, profile string) error {
    server := gotmux.NewServer()
    session, err := server.NewSession(
        gotmux.SessionName(name),
        gotmux.StartDirectory(workDir),
    )
    if err != nil {
        return err
    }

    cmd := fmt.Sprintf("claude --profile %s", profile)
    return session.Windows[0].Panes[0].SendKeys(cmd)
}
```

**Python (libtmux):**
```python
def create_session(name: str, work_dir: str, profile: str):
    server = libtmux.Server()
    session = server.new_session(
        session_name=name,
        start_directory=work_dir,
        attach=False,
    )

    cmd = f"claude --profile {profile}"
    session.active_window.active_pane.send_keys(cmd)
    return session
```

### Async Subprocess

**Go (goroutines):**
```go
func RunClaudePrompt(prompt string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    cmd := exec.CommandContext(ctx, "claude", "-p", prompt)
    output, err := cmd.Output()
    return string(output), err
}
```

**Python (asyncio):**
```python
async def run_claude_prompt(prompt: str) -> str:
    process = await asyncio.create_subprocess_exec(
        "claude", "-p", prompt,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
    )
    stdout, _ = await asyncio.wait_for(
        process.communicate(),
        timeout=120,
    )
    return stdout.decode()
```

---

## When to Use Each

### Choose Go When:

1. **Production deployment** - Single binary, zero dependencies
2. **Performance critical** - Lower memory, faster startup
3. **Existing Go codebase** - Already using Go implementation
4. **Static typing** - Catch errors at compile time
5. **Cross-compilation** - Easy builds for any OS/arch

### Choose Python When:

1. **Rapid prototyping** - Faster iteration cycles
2. **Team familiarity** - More developers know Python
3. **Rich ecosystem** - Access to ML/data libraries
4. **Runtime flexibility** - Duck typing, dynamic modification
5. **Testing ease** - pytest fixtures, mocking

---

## Interoperability

Both implementations are **fully interoperable** because they share:

1. **JSONL format** - Same `.work/work_orders.jsonl` schema
2. **TOML templates** - Same `.work/templates/*.toml` structure
3. **Mail protocol** - Same message types and routing
4. **Assignment files** - Same `.work/.assignment-{agent}` format
5. **Git integration** - Same worktree and branch conventions

**You can mix and match:**
```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        MIXED DEPLOYMENT                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Go CLI (co, wo)                Python CLI (vermas)                        │
│        │                              │                                     │
│        └──────────────┬───────────────┘                                     │
│                       │                                                     │
│                       ▼                                                     │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                   SHARED WORK ORDERS (JSONL)                        │  │
│   │                                                                     │  │
│   │   - Go creates work order → Python reads it                         │  │
│   │   - Python sends mail → Go receives it                              │  │
│   │   - Go spawns worker → Python Supervisor monitors it                │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Migration Path

### Go → Python

1. Install Python CLI alongside Go CLI
2. Both read/write same work orders
3. Gradually replace `co` commands with `vermas` commands
4. Keep using same tmux sessions, assignments, templates

### Python → Go

1. Compile Go binary
2. Both read/write same work orders
3. Gradually replace `vermas` commands with `co` commands
4. No data migration needed - same JSONL format

---

## Recommendation

**For VerMAS development:**

Start with **Python** because:
1. Faster iteration on verification logic
2. Pydantic validation catches schema issues early
3. pytest makes testing verification agents easier
4. libtmux has cleaner API than go-tmux
5. asyncio handles concurrent agents naturally

**Later, if needed:**
- Port performance-critical parts to Go
- Create single-binary distribution
- Both can coexist via shared work orders

The choice is about **developer ergonomics**, not architecture. Both implementations do the same thing the same way.

---

## Additional Python Implementation Patterns

### Event Emission

```python
from pydantic import BaseModel
from datetime import datetime
from pathlib import Path
import json

class Event(BaseModel):
    event_id: str
    event_type: str
    timestamp: datetime
    actor: str
    correlation_id: str | None = None
    data: dict

def emit_event(event: Event, work_path: Path = Path(".work")):
    """Append event to both event log and feed."""
    event_line = event.model_dump_json() + "\n"

    # Append to main event log (source of truth)
    with open(work_path / "events.jsonl", "a") as f:
        f.write(event_line)

    # Append to real-time feed (for watchers)
    with open(work_path / "feed.jsonl", "a") as f:
        f.write(event_line)
```

### Change Feed Watcher

```python
import asyncio
from pathlib import Path

async def watch_feed(work_path: Path = Path(".work")):
    """Async generator that yields events as they arrive."""
    feed_path = work_path / "feed.jsonl"

    with open(feed_path, "r") as f:
        # Start at end of file
        f.seek(0, 2)

        while True:
            line = f.readline()
            if line:
                event = Event.model_validate_json(line.strip())
                yield event
            else:
                await asyncio.sleep(0.1)

# Usage
async def main():
    async for event in watch_feed():
        if event.event_type == "work_order.status_changed":
            print(f"Work order {event.data['work_order_id']} → {event.data['to_status']}")
```

### Assignment Management

```python
from pathlib import Path
from dataclasses import dataclass

@dataclass
class AssignmentContent:
    ref_type: str  # work_order, mail, process
    ref_id: str

class Assignment:
    def __init__(self, agent: str, work_path: Path):
        self.agent = agent
        self.path = work_path / f".assignment-{agent.replace('/', '-')}"

    def check(self) -> AssignmentContent | None:
        """Check if assignment has content."""
        if not self.path.exists():
            return None
        content = self.path.read_text().strip()
        if not content:
            return None
        ref_type, ref_id = content.split(":", 1)
        return AssignmentContent(ref_type, ref_id)

    def set(self, ref_type: str, ref_id: str):
        """Set assignment content."""
        self.path.write_text(f"{ref_type}:{ref_id}")

    def clear(self):
        """Clear assignment."""
        if self.path.exists():
            self.path.unlink()
```

### JSONL Store

```python
from pathlib import Path
from typing import TypeVar, Type, Iterator
from pydantic import BaseModel
import json

T = TypeVar('T', bound=BaseModel)

class JSONLStore:
    """Generic JSONL file store."""

    def __init__(self, path: Path, model: Type[T]):
        self.path = path
        self.model = model

    def append(self, item: T):
        """Append item to store."""
        with open(self.path, "a") as f:
            f.write(item.model_dump_json() + "\n")

    def iter_all(self) -> Iterator[T]:
        """Iterate all items."""
        if not self.path.exists():
            return
        with open(self.path, "r") as f:
            for line in f:
                if line.strip():
                    yield self.model.model_validate_json(line)

    def find(self, **filters) -> T | None:
        """Find first item matching filters."""
        for item in self.iter_all():
            if all(getattr(item, k) == v for k, v in filters.items()):
                return item
        return None

    def filter(self, **filters) -> list[T]:
        """Get all items matching filters."""
        return [
            item for item in self.iter_all()
            if all(getattr(item, k) == v for k, v in filters.items())
        ]
```

### Tmux Session Manager

```python
import libtmux
from dataclasses import dataclass

@dataclass
class AgentSession:
    name: str
    role: str
    factory: str
    worktree: str
    profile: str

class TmuxManager:
    def __init__(self):
        self.server = libtmux.Server()

    def spawn_worker(self, factory: str, slot: int, wo_id: str) -> AgentSession:
        """Spawn a new worker session."""
        name = f"worker-{factory}-slot{slot}"
        worktree = f"{factory}/workers/slot{slot}"

        session = self.server.new_session(
            session_name=name,
            start_directory=worktree,
            attach=False,
            environment={
                "AGENT_ID": f"{factory}/workers/slot{slot}",
                "WORK_ORDER_ID": wo_id,
                "FACTORY": factory,
                "ROLE": "worker",
            }
        )

        # Start Claude Code with worker profile
        session.active_window.active_pane.send_keys("claude --profile worker")

        return AgentSession(
            name=name,
            role="worker",
            factory=factory,
            worktree=worktree,
            profile="worker"
        )

    def list_sessions(self, prefix: str = None) -> list[str]:
        """List tmux sessions, optionally filtered by prefix."""
        sessions = self.server.sessions
        names = [s.name for s in sessions]
        if prefix:
            names = [n for n in names if n.startswith(prefix)]
        return names

    def kill_session(self, name: str):
        """Kill a session by name."""
        session = self.server.find_where({"session_name": name})
        if session:
            session.kill_session()
```

---

## See Also

- [INDEX.md](./INDEX.md) - Documentation map
- [ARCHITECTURE.md](./ARCHITECTURE.md) - System design
- [CLI.md](./CLI.md) - Command reference
- [HOOKS.md](./HOOKS.md) - Claude Code integration
- [EVENTS.md](./EVENTS.md) - Event sourcing
- [SCHEMAS.md](./SCHEMAS.md) - Data specifications
