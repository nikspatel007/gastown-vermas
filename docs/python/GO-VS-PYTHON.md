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
│   │  (Agent)   │     │  (gt/bd)   │     │  (Workers) │                     │
│   └────────────┘     └────────────┘     └────────────┘                     │
│         │                  │                  │                             │
│         └──────────────────┼──────────────────┘                             │
│                            ▼                                                │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                      BEADS (JSONL Files)                             │  │
│   │    issues.jsonl  |  messages.jsonl  |  formulas/*.toml               │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   Orchestration: tmux sessions + Claude Code CLI + hooks + CLAUDE.md       │
│   Storage: JSONL beads (git-backed)                                        │
│   LLM: Claude Code CLI (no API costs)                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Both implementations share:**
- Tmux-based agent sessions
- Claude Code CLI for LLM (no API costs)
- JSONL beads for persistence
- TOML formulas for workflows
- Hook-based work dispatch (GUPP)
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
var hookCmd = &cobra.Command{
    Use:   "hook",
    Short: "Check your assigned work",
    Run: func(cmd *cobra.Command, args []string) {
        actor := os.Getenv("BD_ACTOR")
        hook := LoadHook(actor)
        if hook != nil {
            fmt.Printf("HOOKED: %s\n", hook.BeadID)
        } else {
            fmt.Println("Hook is empty")
        }
    },
}
```

**Python (Typer):**
```python
@app.command()
def hook():
    """Check your assigned work."""
    actor = os.environ.get("BD_ACTOR")
    hook = Hook(actor, Path(".beads"))
    content = hook.check()
    if content:
        console.print(f"HOOKED: {content.ref_id}")
    else:
        console.print("Hook is empty")
```

### Data Models

**Go (structs):**
```go
type Bead struct {
    ID          string       `json:"id"`
    Title       string       `json:"title"`
    Description string       `json:"description"`
    Status      BeadStatus   `json:"status"`
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
class Bead(BaseModel):
    id: str
    title: str
    description: str
    status: BeadStatus
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
3. **Existing Gas Town** - Already using Go implementation
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

1. **JSONL format** - Same `.beads/issues.jsonl` schema
2. **TOML formulas** - Same `.beads/formulas/*.toml` structure
3. **Mail protocol** - Same message types and routing
4. **Hook files** - Same `.beads/.hook-{agent}` format
5. **Git integration** - Same worktree and branch conventions

**You can mix and match:**
```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        MIXED DEPLOYMENT                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Go CLI (gt, bd)                Python CLI (vermas)                        │
│        │                              │                                     │
│        └──────────────┬───────────────┘                                     │
│                       │                                                     │
│                       ▼                                                     │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                    SHARED BEADS (JSONL)                              │  │
│   │                                                                     │  │
│   │   - Go creates bead → Python reads it                               │  │
│   │   - Python sends mail → Go receives it                              │  │
│   │   - Go spawns polecat → Python Witness monitors it                  │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Migration Path

### Go → Python

1. Install Python CLI alongside Go CLI
2. Both read/write same beads
3. Gradually replace `gt` commands with `vermas` commands
4. Keep using same tmux sessions, hooks, formulas

### Python → Go

1. Compile Go binary
2. Both read/write same beads
3. Gradually replace `vermas` commands with `gt` commands
4. No data migration needed - same JSONL format

---

## Recommendation

**For VerMAS development:**

Start with **Python** because:
1. Faster iteration on verification logic
2. Pydantic validation catches schema issues early
3. pytest makes testing Inspector agents easier
4. libtmux has cleaner API than go-tmux
5. asyncio handles concurrent agents naturally

**Later, if needed:**
- Port performance-critical parts to Go
- Create single-binary distribution
- Both can coexist via shared beads

The choice is about **developer ergonomics**, not architecture. Both implementations do the same thing the same way.
