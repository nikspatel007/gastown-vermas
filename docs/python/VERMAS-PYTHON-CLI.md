# VerMAS Python Implementation (CLI-Based)

> Using Claude Code, tmux, and subprocess - No API costs

## Philosophy

**Don't call APIs. Orchestrate CLIs.**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         WRONG: API-based ($$$)                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Python ──► Anthropic API ──► $0.003/1K tokens ──► 💸💸💸                  │
│   Python ──► OpenAI API ──► $0.01/1K tokens ──► 💸💸💸                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                         RIGHT: CLI-based (Free*)                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Python ──► subprocess ──► claude (CLI) ──► Uses your subscription        │
│   Python ──► subprocess ──► codex (CLI) ──► Uses your subscription         │
│   Python ──► libtmux ──► tmux sessions ──► Claude Code in each pane        │
│                                                                             │
│   * Requires Claude Pro/Max subscription                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    VERMAS PYTHON (CLI ORCHESTRATION)                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      ORCHESTRATION LAYER                             │   │
│  │                                                                      │   │
│  │   ┌──────────────┐   ┌──────────────┐   ┌──────────────┐           │   │
│  │   │   Typer CLI  │   │   libtmux    │   │  subprocess  │           │   │
│  │   │  (commands)  │   │  (sessions)  │   │  (claude)    │           │   │
│  │   └──────────────┘   └──────────────┘   └──────────────┘           │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                      │                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        TMUX SESSIONS                                 │   │
│  │                                                                      │   │
│  │   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐           │   │
│  │   │  Mayor   │  │Inspector │  │ Polecat  │  │ Polecat  │           │   │
│  │   │ (claude) │  │ (claude) │  │ (claude) │  │ (claude) │           │   │
│  │   └──────────┘  └──────────┘  └──────────┘  └──────────┘           │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                      │                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      STATE MANAGEMENT                                │   │
│  │                                                                      │   │
│  │   Pydantic Models ──► SQLite/JSONL ──► Shared across sessions       │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Gas Town Architecture

Before diving into the Python structure, understand the Gas Town model:

```
Town (workspace root)
├── mayor/              ← Global coordinator (does NOT write code)
├── .beads/             ← Town-level issue tracking (prefix: hq-)
├── <rig>/              ← Project container (e.g., "myproject")
│   ├── .beads/         ← Rig-level issues (prefix: mp-, etc.)
│   ├── mayor/rig/      ← Read-only reference clone
│   ├── refinery/       ← Merge queue processor
│   ├── witness/        ← Worker lifecycle manager
│   ├── crew/           ← Human-directed workspaces
│   │   ├── frontend/   ← Role-based naming
│   │   └── backend/
│   └── polecats/       ← Ephemeral worker worktrees
│       ├── polecat-1/
│       └── polecat-2/
```

## Agent Roles

### Town-Wide Agents

| Role | Purpose | Writes Code? |
|------|---------|--------------|
| **Mayor** | Cross-rig coordinator. Dispatches work, handles escalations | **NO** |
| **Deacon** | Daemon managing agent lifecycle and plugins | No |
| **Overseer** | Human role. Strategy, reviews, escalations | Manual |

### Per-Rig Agents

| Role | Purpose | Writes Code? |
|------|---------|--------------|
| **Witness** | Monitors workers, detects stuck processes, nudges | No |
| **Refinery** | Manages merge queues, code review workflows | No |
| **Polecat** | Ephemeral workers (spawn → work → disappear) | **YES** |
| **Crew** | Human-directed workspaces for hands-on work | **YES** |

### VerMAS Inspector Ecosystem (NEW)

| Role | Purpose | Under |
|------|---------|-------|
| **Inspector** | Quality gate coordinator | Town |
| **Designer** | Elaborates vague requests into requirements | Mayor |
| **Strategist** | Proposes verification criteria from requirements | Inspector |
| **Verifier** | Runs objective tests (NO LLM - just shell) | Inspector |
| **Auditor** | Checks compliance against spec | Inspector |
| **Advocate** | Argues FOR the code (defense attorney) | Inspector |
| **Critic** | Argues AGAINST the code (prosecutor) | Inspector |
| **Judge** | Delivers verdict (PASS/FAIL/NEEDS_HUMAN) | Inspector |

## Key Concepts

### Hooks
Each agent has a "hook" where work hangs. **GUPP: If your hook has work, RUN IT.**

### Beads
Git-backed issue tracking. Commands: `bd create`, `bd ready`, `bd close`, `bd sync`

### Molecules (MEOW)
Workflow states: Ice-9 (formula) → Solid (proto) → Liquid (mol) → Vapor (wisp)

### Convoys
Groupings of related work issues that travel together.

### Mail
Internal messaging: `gt mail inbox`, `gt mail send <addr> -s "Subject" -m "Body"`

---

## Project Structure

```
vermas-py/
├── pyproject.toml
├── vermas/
│   ├── __init__.py
│   ├── cli.py                  # Main CLI (vermas command)
│   │
│   ├── models/                 # Pydantic models
│   │   ├── __init__.py
│   │   ├── bead.py            # WorkItem, TestSpec, Message
│   │   ├── agent.py           # AgentState, AgentRole
│   │   ├── molecule.py        # Molecule states (MEOW)
│   │   ├── verification.py    # Verdict, Evidence, Brief
│   │   └── config.py          # Settings
│   │
│   ├── tmux/                   # Tmux session management
│   │   ├── __init__.py
│   │   ├── server.py          # libtmux server wrapper
│   │   ├── session.py         # Create/manage sessions
│   │   ├── pane.py            # Pane operations
│   │   └── layout.py          # Multi-pane layouts
│   │
│   ├── claude/                 # Claude Code CLI wrapper
│   │   ├── __init__.py
│   │   ├── runner.py          # Run claude -p "prompt"
│   │   ├── session.py         # Interactive sessions in tmux
│   │   └── profiles.py        # Profile management
│   │
│   ├── agents/                 # Agent orchestration
│   │   ├── __init__.py
│   │   ├── base.py            # Base agent (tmux + claude)
│   │   │
│   │   ├── town/              # Town-wide agents
│   │   │   ├── __init__.py
│   │   │   ├── mayor.py       # Global coordinator
│   │   │   ├── deacon.py      # Daemon/lifecycle manager
│   │   │   └── overseer.py    # Human interface
│   │   │
│   │   ├── rig/               # Per-rig agents
│   │   │   ├── __init__.py
│   │   │   ├── witness.py     # Worker monitor, stuck detection
│   │   │   ├── refinery.py    # Merge queue processor
│   │   │   ├── polecat.py     # Ephemeral worker
│   │   │   └── crew.py        # Human-directed workspace
│   │   │
│   │   ├── inspector/         # VerMAS Inspector ecosystem
│   │   │   ├── __init__.py
│   │   │   ├── coordinator.py # Inspector main
│   │   │   ├── designer.py    # Requirements elaboration
│   │   │   ├── strategist.py  # Criteria proposal
│   │   │   ├── verifier.py    # Objective tests (NO LLM)
│   │   │   ├── auditor.py     # Compliance check
│   │   │   ├── advocate.py    # Defense builder
│   │   │   ├── critic.py      # Prosecution builder
│   │   │   └── judge.py       # Verdict deliverer
│   │   │
│   │   └── roles/             # CLAUDE.md templates
│   │       ├── mayor.md
│   │       ├── deacon.md
│   │       ├── witness.md
│   │       ├── refinery.md
│   │       ├── polecat.md
│   │       ├── crew.md
│   │       ├── inspector.md
│   │       ├── designer.md
│   │       ├── strategist.md
│   │       ├── advocate.md
│   │       ├── critic.md
│   │       └── judge.md
│   │
│   ├── beads/                  # Git-backed issue tracking
│   │   ├── __init__.py
│   │   ├── db.py              # JSONL read/write (Go-compatible)
│   │   ├── issue.py           # Issue CRUD
│   │   ├── molecule.py        # Molecule operations (MEOW)
│   │   ├── formula.py         # Formula parsing (TOML)
│   │   └── sync.py            # Git sync
│   │
│   ├── mail/                   # Inter-agent messaging
│   │   ├── __init__.py
│   │   ├── mailbox.py         # Send/receive messages
│   │   └── router.py          # Address routing
│   │
│   ├── convoy/                 # Work groupings
│   │   ├── __init__.py
│   │   └── manager.py         # Convoy operations
│   │
│   └── hooks/                  # Hook system
│       ├── __init__.py
│       └── manager.py         # Hook checking (GUPP)
│
└── tests/
```

## Core Components

### 1. Tmux Session Manager

```python
# vermas/tmux/session.py
import libtmux
from pathlib import Path
from typing import Optional
from pydantic import BaseModel


class TmuxSession(BaseModel):
    """Represents a tmux session for an agent."""
    name: str
    working_dir: Path
    role: str
    pid: Optional[int] = None

    class Config:
        arbitrary_types_allowed = True


class TmuxManager:
    """Manage tmux sessions for VerMAS agents."""

    def __init__(self):
        self.server = libtmux.Server()

    def create_session(
        self,
        name: str,
        working_dir: Path,
        start_command: Optional[str] = None,
    ) -> libtmux.Session:
        """Create a new tmux session."""
        session = self.server.new_session(
            session_name=name,
            start_directory=str(working_dir),
            attach=False,
        )

        if start_command:
            window = session.active_window
            pane = window.active_pane
            pane.send_keys(start_command)

        return session

    def get_session(self, name: str) -> Optional[libtmux.Session]:
        """Get an existing session by name."""
        try:
            return self.server.sessions.get(session_name=name)
        except Exception:
            return None

    def kill_session(self, name: str) -> bool:
        """Kill a session."""
        session = self.get_session(name)
        if session:
            session.kill()
            return True
        return False

    def list_sessions(self) -> list[str]:
        """List all VerMAS sessions."""
        return [
            s.name for s in self.server.sessions
            if s.name.startswith(("mayor-", "inspector-", "polecat-", "vermas-"))
        ]

    def send_keys(self, session_name: str, keys: str, pane_index: int = 0):
        """Send keys to a session's pane."""
        session = self.get_session(session_name)
        if session:
            pane = session.active_window.panes[pane_index]
            pane.send_keys(keys)

    def capture_pane(self, session_name: str, pane_index: int = 0) -> str:
        """Capture output from a pane."""
        session = self.get_session(session_name)
        if session:
            pane = session.active_window.panes[pane_index]
            return "\n".join(pane.capture_pane())
        return ""

    def create_split_layout(
        self,
        session_name: str,
        panes: list[dict],
        vertical: bool = True,
    ):
        """Create a multi-pane layout."""
        session = self.get_session(session_name)
        if not session:
            return

        window = session.active_window

        for i, pane_config in enumerate(panes[1:], 1):
            window.split(
                vertical=vertical,
                start_directory=pane_config.get("cwd"),
            )

        # Send start commands
        for i, pane_config in enumerate(panes):
            if "command" in pane_config:
                window.panes[i].send_keys(pane_config["command"])
```

### 2. Claude Code Runner

```python
# vermas/claude/runner.py
import subprocess
import asyncio
import json
from pathlib import Path
from typing import Optional, AsyncIterator
from pydantic import BaseModel


class ClaudeResult(BaseModel):
    """Result from a Claude Code execution."""
    success: bool
    output: str
    exit_code: int
    session_id: Optional[str] = None


class ClaudeRunner:
    """Run Claude Code CLI commands."""

    def __init__(self, working_dir: Path = Path(".")):
        self.working_dir = working_dir

    async def run_prompt(
        self,
        prompt: str,
        profile: Optional[str] = None,
        timeout: int = 300,
    ) -> ClaudeResult:
        """
        Run a single prompt through Claude Code.

        Uses: claude -p "prompt" --output-format json
        """
        cmd = ["claude"]

        if profile:
            cmd.extend(["--profile", profile])

        cmd.extend(["-p", prompt, "--output-format", "json"])

        try:
            process = await asyncio.create_subprocess_exec(
                *cmd,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE,
                cwd=self.working_dir,
            )

            stdout, stderr = await asyncio.wait_for(
                process.communicate(),
                timeout=timeout,
            )

            return ClaudeResult(
                success=process.returncode == 0,
                output=stdout.decode(),
                exit_code=process.returncode,
            )

        except asyncio.TimeoutError:
            process.kill()
            return ClaudeResult(
                success=False,
                output=f"Timeout after {timeout}s",
                exit_code=-1,
            )
        except Exception as e:
            return ClaudeResult(
                success=False,
                output=str(e),
                exit_code=-1,
            )

    async def run_with_context(
        self,
        prompt: str,
        context_files: list[Path],
        profile: Optional[str] = None,
    ) -> ClaudeResult:
        """
        Run prompt with file context.

        Uses: claude -p "prompt" file1.py file2.py
        """
        cmd = ["claude"]

        if profile:
            cmd.extend(["--profile", profile])

        cmd.extend(["-p", prompt])
        cmd.extend([str(f) for f in context_files])
        cmd.append("--output-format")
        cmd.append("json")

        process = await asyncio.create_subprocess_exec(
            *cmd,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
            cwd=self.working_dir,
        )

        stdout, stderr = await process.communicate()

        return ClaudeResult(
            success=process.returncode == 0,
            output=stdout.decode(),
            exit_code=process.returncode,
        )

    def start_interactive(
        self,
        profile: Optional[str] = None,
        system_prompt: Optional[str] = None,
    ) -> subprocess.Popen:
        """
        Start an interactive Claude Code session.

        Returns the process handle for the tmux pane to manage.
        """
        cmd = ["claude"]

        if profile:
            cmd.extend(["--profile", profile])

        if system_prompt:
            cmd.extend(["--system-prompt", system_prompt])

        return subprocess.Popen(
            cmd,
            cwd=self.working_dir,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )


class ClaudeSession:
    """
    Manage an interactive Claude Code session in tmux.

    This wraps a tmux pane running `claude` interactively.
    """

    def __init__(
        self,
        session_name: str,
        working_dir: Path,
        role: str,
        tmux_manager: "TmuxManager",
    ):
        self.session_name = session_name
        self.working_dir = working_dir
        self.role = role
        self.tmux = tmux_manager
        self._started = False

    def start(self, profile: Optional[str] = None):
        """Start the Claude session in tmux."""
        if self._started:
            return

        # Build the claude command
        cmd = "claude"
        if profile:
            cmd = f"claude --profile {profile}"

        # Create or get the tmux session
        session = self.tmux.get_session(self.session_name)
        if not session:
            session = self.tmux.create_session(
                self.session_name,
                self.working_dir,
                start_command=cmd,
            )

        self._started = True

    def send_message(self, message: str):
        """Send a message to the Claude session."""
        if not self._started:
            raise RuntimeError("Session not started")

        # Send to tmux pane
        self.tmux.send_keys(self.session_name, message)
        self.tmux.send_keys(self.session_name, "Enter")

    def get_output(self) -> str:
        """Get recent output from the session."""
        return self.tmux.capture_pane(self.session_name)

    def stop(self):
        """Stop the Claude session."""
        self.tmux.kill_session(self.session_name)
        self._started = False
```

### 3. Agent Base Class

```python
# vermas/agents/base.py
import os
from pathlib import Path
from typing import Optional
from abc import ABC, abstractmethod

from vermas.tmux.session import TmuxManager
from vermas.claude.runner import ClaudeRunner, ClaudeSession
from vermas.models.agent import AgentRole, AgentState
from vermas.db.beads import BeadsDB
from vermas.mail.mailbox import Mailbox


class BaseAgent(ABC):
    """
    Base class for VerMAS agents.

    Each agent runs as a Claude Code session in a tmux pane.
    Communication happens via the beads mail system.
    """

    def __init__(
        self,
        role: AgentRole,
        working_dir: Path,
        agent_id: Optional[str] = None,
    ):
        self.role = role
        self.working_dir = working_dir
        self.agent_id = agent_id or f"{role.value}-{os.getpid()}"

        # Components
        self.tmux = TmuxManager()
        self.claude = ClaudeRunner(working_dir)
        self.beads = BeadsDB(working_dir / ".beads")
        self.mailbox = Mailbox(self.agent_id, self.beads)

        # Session
        self.session: Optional[ClaudeSession] = None

    @property
    def session_name(self) -> str:
        """Tmux session name for this agent."""
        return self.agent_id

    @property
    def role_context_path(self) -> Path:
        """Path to the CLAUDE.md for this role."""
        return self.working_dir / "CLAUDE.md"

    def setup_role_context(self):
        """Write the CLAUDE.md file for this agent's role."""
        from vermas.agents.roles import get_role_context

        context = get_role_context(self.role)
        self.role_context_path.write_text(context)

    def start(self):
        """Start the agent's Claude session in tmux."""
        self.setup_role_context()

        self.session = ClaudeSession(
            session_name=self.session_name,
            working_dir=self.working_dir,
            role=self.role.value,
            tmux_manager=self.tmux,
        )

        self.session.start(profile=self.role.value)

    def stop(self):
        """Stop the agent's session."""
        if self.session:
            self.session.stop()

    def send_to_claude(self, message: str):
        """Send a message to this agent's Claude session."""
        if self.session:
            self.session.send_message(message)

    def check_mail(self) -> list:
        """Check for new mail."""
        return self.mailbox.get_unread()

    def send_mail(self, to: str, subject: str, body: str):
        """Send mail to another agent."""
        self.mailbox.send(to, subject, body)

    @abstractmethod
    async def run_task(self, task: str) -> str:
        """
        Run a task through Claude Code.

        This is the main entry point for agent work.
        Subclasses implement role-specific behavior.
        """
        pass
```

### 4. Mayor Agent

```python
# vermas/agents/mayor.py
from pathlib import Path
from typing import Optional

from vermas.agents.base import BaseAgent
from vermas.models.agent import AgentRole
from vermas.models.bead import WorkItem


class MayorAgent(BaseAgent):
    """
    Mayor agent - coordinates work and talks to user.

    Runs in an interactive tmux session with Claude Code.
    """

    def __init__(self, working_dir: Path):
        super().__init__(
            role=AgentRole.MAYOR,
            working_dir=working_dir,
            agent_id=f"mayor-{os.getpid()}",
        )

    async def run_task(self, task: str) -> str:
        """
        Process a user request.

        Instead of calling APIs, we send the task to the Claude session
        and let it handle the interaction.
        """
        # Send task to Claude session
        prompt = f"""
User request: {task}

Please:
1. Analyze this request
2. Call Designer to elaborate requirements (use /design skill)
3. Create a work item with bd create
4. Wait for Inspector to approve the spec
5. Sling to a polecat when ready

Start by understanding what the user wants.
"""
        self.send_to_claude(prompt)

        # The actual work happens in the Claude session
        # We just orchestrate and monitor
        return "Task sent to Mayor session"

    async def create_work_item(self, title: str, description: str) -> WorkItem:
        """Create a work item via Claude session."""
        prompt = f"""
Create a work item:
Title: {title}
Description: {description}

Use: bd create --title="{title}" --type=task

Then report the created bead ID.
"""
        result = await self.claude.run_prompt(prompt, profile="mayor")
        # Parse result to get bead ID
        return self._parse_work_item(result.output)

    def _parse_work_item(self, output: str) -> WorkItem:
        """Parse bd create output to get WorkItem."""
        # Extract bead ID from output like "Created: gt-abc123"
        import re
        match = re.search(r"Created: (gt-[a-z0-9]+)", output)
        if match:
            return WorkItem(
                id=match.group(1),
                type="task",
                title="",  # Would need to parse more
            )
        raise ValueError(f"Could not parse work item from: {output}")
```

### 5. Inspector Agent

```python
# vermas/agents/inspector.py
import asyncio
from pathlib import Path
from typing import Optional

from vermas.agents.base import BaseAgent
from vermas.models.agent import AgentRole
from vermas.models.verification import VerificationResult, Verdict


class InspectorAgent(BaseAgent):
    """
    Inspector agent - coordinates verification.

    Manages sub-agents: Strategist, Verifier, Advocate, Critic, Judge.
    Each sub-agent runs in its own Claude session.
    """

    def __init__(self, working_dir: Path):
        super().__init__(
            role=AgentRole.INSPECTOR,
            working_dir=working_dir,
            agent_id=f"inspector-{os.getpid()}",
        )

        # Sub-agent sessions (created on demand)
        self.strategist_session: Optional[ClaudeSession] = None
        self.advocate_session: Optional[ClaudeSession] = None
        self.critic_session: Optional[ClaudeSession] = None
        self.judge_session: Optional[ClaudeSession] = None

    async def run_task(self, task: str) -> str:
        """Process an inspection request."""
        self.send_to_claude(task)
        return "Task sent to Inspector session"

    async def propose_criteria(self, work_item_id: str) -> list:
        """
        Use Strategist to propose verification criteria.

        Runs in a separate Claude session with strategist profile.
        """
        prompt = f"""
You are the Strategist. Analyze work item {work_item_id} and propose verification criteria.

Steps:
1. Read the work item: bd show {work_item_id}
2. Understand the requirements
3. Propose 5-10 testable criteria
4. For each criterion, provide:
   - ID (AC-1, AC-2, etc.)
   - Description
   - Bash command that exits 0 for pass

Output the criteria in a clear format.
"""
        result = await self.claude.run_prompt(prompt, profile="strategist")
        return self._parse_criteria(result.output)

    async def run_verification(self, work_item_id: str) -> VerificationResult:
        """
        Run full adversarial verification.

        1. Verifier runs objective tests (no LLM, just shell)
        2. Advocate argues FOR (Claude session)
        3. Critic argues AGAINST (Claude session, in parallel)
        4. Judge delivers verdict (Claude session)
        """
        # Step 1: Run verifier (shell commands, no LLM)
        test_results = await self._run_verifier(work_item_id)

        # Step 2 & 3: Run advocate and critic in parallel
        advocate_task = self._run_advocate(work_item_id, test_results)
        critic_task = self._run_critic(work_item_id, test_results)

        advocate_brief, critic_brief = await asyncio.gather(
            advocate_task,
            critic_task,
        )

        # Step 4: Run judge
        verdict = await self._run_judge(
            work_item_id,
            test_results,
            advocate_brief,
            critic_brief,
        )

        return verdict

    async def _run_verifier(self, work_item_id: str) -> list:
        """
        Run objective tests.

        This does NOT use an LLM - it just runs shell commands.
        """
        from vermas.agents.verifier import Verifier

        verifier = Verifier(self.working_dir)
        spec = self.beads.get_spec_for_work_item(work_item_id)
        return await verifier.run_criteria(spec.criteria)

    async def _run_advocate(self, work_item_id: str, test_results: list) -> str:
        """Run Advocate to build defense."""
        prompt = f"""
You are the Advocate. Build a defense for work item {work_item_id}.

Test Results:
{self._format_test_results(test_results)}

Your job:
1. Review the test results
2. Review the code changes: git diff main...HEAD
3. Argue WHY this code should be merged
4. Address any failing tests
5. Highlight strengths

Be persuasive but honest.
"""
        result = await self.claude.run_prompt(prompt, profile="advocate")
        return result.output

    async def _run_critic(self, work_item_id: str, test_results: list) -> str:
        """Run Critic to build prosecution."""
        prompt = f"""
You are the Critic. Attack work item {work_item_id}.

Test Results:
{self._format_test_results(test_results)}

Your job:
1. Review the test results
2. Review the code changes: git diff main...HEAD
3. Argue WHY this code should NOT be merged
4. Find bugs, security issues, missing tests
5. Be thorough but fair

Don't manufacture false concerns.
"""
        result = await self.claude.run_prompt(prompt, profile="critic")
        return result.output

    async def _run_judge(
        self,
        work_item_id: str,
        test_results: list,
        advocate_brief: str,
        critic_brief: str,
    ) -> VerificationResult:
        """Run Judge to deliver verdict."""
        prompt = f"""
You are the Judge. Deliver a verdict for work item {work_item_id}.

Test Results:
{self._format_test_results(test_results)}

Advocate's Defense:
{advocate_brief}

Critic's Concerns:
{critic_brief}

Deliver your verdict:
- PASS: Code meets requirements, concerns addressed
- FAIL: Critical issues must be fixed
- NEEDS_HUMAN: Cannot decide, need human input

Output JSON:
{{"verdict": "PASS|FAIL|NEEDS_HUMAN", "confidence": 0.0-1.0, "reasoning": "..."}}
"""
        result = await self.claude.run_prompt(prompt, profile="judge")
        return self._parse_verdict(result.output)

    def _format_test_results(self, results: list) -> str:
        """Format test results for prompts."""
        lines = []
        for r in results:
            status = "✅" if r.status == "pass" else "❌"
            lines.append(f"{status} {r.criterion_id}: {r.status}")
        return "\n".join(lines)

    def _parse_verdict(self, output: str) -> VerificationResult:
        """Parse judge output to VerificationResult."""
        import json
        import re

        # Find JSON in output
        match = re.search(r'\{.*\}', output, re.DOTALL)
        if match:
            data = json.loads(match.group())
            return VerificationResult(
                verdict=Verdict(data["verdict"]),
                confidence=data.get("confidence", 0.5),
                judge_reasoning=data.get("reasoning", ""),
            )
        raise ValueError(f"Could not parse verdict from: {output}")
```

### 6. Verifier (No LLM)

```python
# vermas/agents/verifier.py
import asyncio
import subprocess
from pathlib import Path
from dataclasses import dataclass

from vermas.models.bead import AcceptanceCriterion


@dataclass
class TestResult:
    """Result of running a single test."""
    criterion_id: str
    status: str  # pass, fail, error, timeout
    output: str
    exit_code: int
    duration_ms: int


class Verifier:
    """
    Objective test runner.

    This does NOT use an LLM - it simply runs shell commands
    and reports pass/fail based on exit codes.
    """

    def __init__(self, working_dir: Path):
        self.working_dir = working_dir

    async def run_criteria(
        self,
        criteria: list[AcceptanceCriterion],
    ) -> list[TestResult]:
        """Run all acceptance criteria."""
        results = []
        for criterion in criteria:
            result = await self.run_single(criterion)
            results.append(result)
        return results

    async def run_single(self, criterion: AcceptanceCriterion) -> TestResult:
        """Run a single criterion."""
        import time

        start = time.time()

        try:
            process = await asyncio.create_subprocess_shell(
                criterion.verify_command,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.STDOUT,
                cwd=self.working_dir,
            )

            try:
                stdout, _ = await asyncio.wait_for(
                    process.communicate(),
                    timeout=criterion.timeout_seconds,
                )
                output = stdout.decode()
                exit_code = process.returncode
                status = "pass" if exit_code == 0 else "fail"

            except asyncio.TimeoutError:
                process.kill()
                output = f"TIMEOUT after {criterion.timeout_seconds}s"
                exit_code = -1
                status = "timeout"

        except Exception as e:
            output = str(e)
            exit_code = -1
            status = "error"

        duration_ms = int((time.time() - start) * 1000)

        return TestResult(
            criterion_id=criterion.id,
            status=status,
            output=output,
            exit_code=exit_code,
            duration_ms=duration_ms,
        )
```

### 7. CLI Interface

```python
# vermas/cli.py
import os
from pathlib import Path
import typer
from rich.console import Console
from rich.table import Table

app = typer.Typer(name="vermas", help="Verifiable Multi-Agent System (CLI-based)")
console = Console()


@app.command()
def start(
    working_dir: Path = typer.Option(Path("."), "--dir", "-d", help="Working directory"),
):
    """Start VerMAS with Mayor and Inspector in tmux split."""
    from vermas.tmux.session import TmuxManager

    tmux = TmuxManager()

    session_name = f"vermas-{os.getpid()}"

    # Create session with Mayor
    session = tmux.create_session(
        session_name,
        working_dir,
        start_command="claude --profile mayor",
    )

    # Split and add Inspector
    window = session.active_window
    window.split(vertical=True, start_directory=str(working_dir))
    window.panes[1].send_keys("claude --profile inspector")

    console.print(f"[green]✅ VerMAS started: {session_name}[/green]")
    console.print(f"\nAttach with: tmux attach -t {session_name}")
    console.print("Left pane: Mayor | Right pane: Inspector")


@app.command()
def attach(
    session_name: str = typer.Argument(None, help="Session name (default: most recent)"),
):
    """Attach to a VerMAS session."""
    from vermas.tmux.session import TmuxManager

    tmux = TmuxManager()

    if not session_name:
        # Find most recent vermas session
        sessions = tmux.list_sessions()
        vermas_sessions = [s for s in sessions if s.startswith("vermas-")]
        if not vermas_sessions:
            console.print("[red]No VerMAS sessions found[/red]")
            raise typer.Exit(1)
        session_name = vermas_sessions[-1]

    os.execvp("tmux", ["tmux", "attach", "-t", session_name])


@app.command()
def status():
    """Show status of all VerMAS sessions."""
    from vermas.tmux.session import TmuxManager

    tmux = TmuxManager()
    sessions = tmux.list_sessions()

    if not sessions:
        console.print("[dim]No active sessions[/dim]")
        return

    table = Table(title="VerMAS Sessions")
    table.add_column("Session")
    table.add_column("Role")
    table.add_column("Status")

    for name in sessions:
        role = name.split("-")[0]
        table.add_row(name, role, "running")

    console.print(table)


@app.command()
def stop(
    session_name: str = typer.Argument(..., help="Session to stop"),
    all_sessions: bool = typer.Option(False, "--all", "-a", help="Stop all sessions"),
):
    """Stop VerMAS sessions."""
    from vermas.tmux.session import TmuxManager

    tmux = TmuxManager()

    if all_sessions:
        for name in tmux.list_sessions():
            tmux.kill_session(name)
            console.print(f"Stopped: {name}")
    else:
        if tmux.kill_session(session_name):
            console.print(f"[green]Stopped: {session_name}[/green]")
        else:
            console.print(f"[red]Session not found: {session_name}[/red]")


# Inspector subcommands
inspect_app = typer.Typer(help="Inspector commands")
app.add_typer(inspect_app, name="inspect")


@inspect_app.command("run")
def inspect_run(
    bead_id: str = typer.Argument(..., help="Work item ID"),
):
    """Run verification on a work item."""
    import asyncio
    from vermas.agents.inspector import InspectorAgent

    async def do_inspect():
        inspector = InspectorAgent(Path("."))
        result = await inspector.run_verification(bead_id)

        if result.verdict.value == "PASS":
            console.print(f"[green]✅ PASS (confidence: {result.confidence})[/green]")
        elif result.verdict.value == "FAIL":
            console.print(f"[red]❌ FAIL[/red]")
            console.print(f"\nReasoning: {result.judge_reasoning}")
        else:
            console.print(f"[yellow]❓ NEEDS_HUMAN[/yellow]")
            console.print(f"\n{result.judge_reasoning}")

    asyncio.run(do_inspect())


@inspect_app.command("verify")
def inspect_verify(
    bead_id: str = typer.Argument(..., help="Work item ID"),
):
    """Run just the objective tests (no LLM)."""
    import asyncio
    from vermas.agents.verifier import Verifier
    from vermas.db.beads import BeadsDB

    async def do_verify():
        beads = BeadsDB(Path(".beads"))
        spec = beads.get_spec_for_work_item(bead_id)

        if not spec:
            console.print(f"[red]No spec found for {bead_id}[/red]")
            raise typer.Exit(1)

        verifier = Verifier(Path("."))
        results = await verifier.run_criteria(spec.criteria)

        table = Table(title=f"Test Results for {bead_id}")
        table.add_column("Criterion")
        table.add_column("Status")
        table.add_column("Duration")

        passed = 0
        for r in results:
            status_str = "[green]PASS[/green]" if r.status == "pass" else "[red]FAIL[/red]"
            table.add_row(r.criterion_id, status_str, f"{r.duration_ms}ms")
            if r.status == "pass":
                passed += 1

        console.print(table)
        console.print(f"\n{passed}/{len(results)} criteria passed")

    asyncio.run(do_verify())


if __name__ == "__main__":
    app()
```

---

## Dependencies

```toml
# pyproject.toml
[project]
name = "vermas"
version = "0.1.0"
description = "Verifiable Multi-Agent System (CLI-based)"
requires-python = ">=3.11"
dependencies = [
    "pydantic>=2.0",
    "pydantic-settings>=2.0",
    "typer>=0.12.0",
    "rich>=13.0",
    "libtmux>=0.37.0",
]

[project.optional-dependencies]
dev = [
    "pytest>=8.0",
    "pytest-asyncio>=0.23",
    "ruff>=0.4.0",
]

[project.scripts]
vermas = "vermas.cli:app"
```

**Note**: No LangChain, no LangGraph, no API clients. Just:
- **libtmux** - Tmux session management
- **pydantic** - Data models
- **typer/rich** - CLI
- **subprocess/asyncio** - Run `claude` CLI

---

## How It Works

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              EXECUTION FLOW                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   User: vermas start                                                        │
│      │                                                                      │
│      ▼                                                                      │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                        TMUX SESSION                                  │  │
│   │                                                                      │  │
│   │   ┌─────────────────────┐   ┌─────────────────────┐                 │  │
│   │   │       MAYOR         │   │      INSPECTOR      │                 │  │
│   │   │                     │   │                     │                 │  │
│   │   │  $ claude           │   │  $ claude           │                 │  │
│   │   │    --profile mayor  │   │    --profile        │                 │  │
│   │   │                     │   │    inspector        │                 │  │
│   │   │  (interactive)      │   │  (interactive)      │                 │  │
│   │   │                     │   │                     │                 │  │
│   │   └─────────────────────┘   └─────────────────────┘                 │  │
│   │                                                                      │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   User types in Mayor pane:                                                 │
│   > Create a FizzBuzz program                                              │
│                                                                             │
│   Claude Code (Mayor profile) responds:                                     │
│   - Calls Designer (via /design skill or subprocess)                       │
│   - Creates work item (bd create)                                          │
│   - Sends mail to Inspector (gt mail send)                                 │
│                                                                             │
│   Inspector pane receives mail, proposes criteria...                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Cost: $0 (beyond subscription)

- All LLM work goes through `claude` CLI
- Uses your Claude Pro/Max subscription
- No per-token API charges
- Unlimited usage within subscription limits
