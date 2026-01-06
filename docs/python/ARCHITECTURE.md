# VerMAS Python Architecture

> Python implementation of Gas Town + VerMAS using CLI orchestration

## Overview

This document maps Gas Town concepts to Python implementations using:
- **libtmux**: Session and pane management
- **subprocess**: Claude Code CLI invocation
- **Pydantic**: Data models and validation
- **asyncio**: Concurrent agent execution
- **JSONL files**: Beads-compatible persistence

**Key principle**: No API costs. All LLM interactions go through Claude Code CLI.

---

## Two-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              TOWN LEVEL                                      │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                        Python Orchestrator                          │  │
│   │                                                                     │  │
│   │   TownManager ──────────────────────────────────────────────────   │  │
│   │        │                                                            │  │
│   │        ├── RigManager(rig1) ─┬── WitnessSession                    │  │
│   │        │                     ├── RefinerySession                   │  │
│   │        │                     └── PolecatPool                       │  │
│   │        │                                                            │  │
│   │        ├── RigManager(rig2) ─┬── WitnessSession                    │  │
│   │        │                     ├── RefinerySession                   │  │
│   │        │                     └── PolecatPool                       │  │
│   │        │                                                            │  │
│   │        └── DeaconService (monitors all rigs)                       │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                      Town Beads (.beads/)                           │  │
│   │    issues.jsonl  |  messages.jsonl  |  routes.jsonl                 │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                              RIG LEVEL                                       │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                        Tmux Sessions                                │  │
│   │                                                                     │  │
│   │   witness-{rig}    refinery-{rig}    polecat-{rig}-{slot}          │  │
│   │        │                │                    │                      │  │
│   │        │                │                    │                      │  │
│   │        ▼                ▼                    ▼                      │  │
│   │   ┌─────────┐     ┌─────────┐          ┌─────────┐                 │  │
│   │   │ Claude  │     │ Claude  │          │ Claude  │                 │  │
│   │   │  Code   │     │  Code   │          │  Code   │                 │  │
│   │   │ --role  │     │ --role  │          │ --role  │                 │  │
│   │   │ witness │     │ refinery│          │ polecat │                 │  │
│   │   └─────────┘     └─────────┘          └─────────┘                 │  │
│   │                                                                     │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                      Rig Beads (.beads/)                            │  │
│   │    issues.jsonl  |  formulas/*.toml  |  mols/*.json                 │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Directory Structure

```
vermas-py/
├── vermas/
│   ├── __init__.py
│   ├── cli.py                    # Typer CLI entry point
│   │
│   ├── models/                   # Pydantic data models
│   │   ├── __init__.py
│   │   ├── bead.py               # Issue, Dependency, Message
│   │   ├── molecule.py           # Formula, Protomol, Mol, Wisp
│   │   ├── agent.py              # AgentConfig, AgentState
│   │   ├── convoy.py             # Convoy tracking
│   │   └── verification.py       # VerMAS Inspector models
│   │
│   ├── core/                     # Core orchestration
│   │   ├── __init__.py
│   │   ├── town.py               # TownManager
│   │   ├── rig.py                # RigManager
│   │   ├── session.py            # TmuxSessionManager
│   │   └── hooks.py              # Hook system
│   │
│   ├── agents/                   # Agent implementations
│   │   ├── __init__.py
│   │   ├── base.py               # BaseAgent ABC
│   │   ├── mayor.py              # Mayor coordinator
│   │   ├── deacon.py             # Deacon daemon
│   │   ├── witness.py            # Witness monitor
│   │   ├── refinery.py           # Refinery merge processor
│   │   ├── polecat.py            # Polecat worker
│   │   └── inspector/            # VerMAS Inspector ecosystem
│   │       ├── __init__.py
│   │       ├── designer.py
│   │       ├── strategist.py
│   │       ├── verifier.py
│   │       ├── auditor.py
│   │       ├── advocate.py
│   │       ├── critic.py
│   │       └── judge.py
│   │
│   ├── beads/                    # Beads persistence layer
│   │   ├── __init__.py
│   │   ├── store.py              # JSONL read/write
│   │   ├── router.py             # Prefix-based routing
│   │   └── sync.py               # Git sync operations
│   │
│   ├── molecules/                # Molecule workflow system
│   │   ├── __init__.py
│   │   ├── formula.py            # TOML formula parser
│   │   ├── lifecycle.py          # cook/pour/wisp/squash/burn
│   │   └── executor.py           # Step execution
│   │
│   ├── mail/                     # Mail protocol
│   │   ├── __init__.py
│   │   ├── inbox.py              # Message retrieval
│   │   ├── send.py               # Message sending
│   │   └── protocol.py           # Message type handlers
│   │
│   └── tmux/                     # Tmux management
│       ├── __init__.py
│       ├── session.py            # Session CRUD
│       ├── pane.py               # Pane management
│       └── layout.py             # Window layouts
│
├── templates/                    # CLAUDE.md templates per role
│   ├── mayor.md
│   ├── witness.md
│   ├── refinery.md
│   ├── polecat.md
│   └── inspector/
│       ├── designer.md
│       ├── strategist.md
│       ├── verifier.md
│       └── ...
│
├── tests/
│   ├── test_beads.py
│   ├── test_molecules.py
│   ├── test_agents.py
│   └── test_verification.py
│
└── pyproject.toml
```

---

## Core Components

### 1. TownManager

Central orchestrator that manages all rigs and town-level operations.

```python
# vermas/core/town.py
from pathlib import Path
from typing import Dict, Optional
import asyncio

from pydantic import BaseModel
from vermas.core.rig import RigManager
from vermas.agents.deacon import DeaconService
from vermas.beads.store import BeadStore


class TownConfig(BaseModel):
    """Town configuration."""
    root: Path
    beads_prefix: str = "hq"


class TownManager:
    """
    Central coordinator for Gas Town.

    Responsibilities:
    - Manage rig lifecycle
    - Route beads by prefix
    - Coordinate cross-rig operations
    - Run Deacon service
    """

    def __init__(self, config: TownConfig):
        self.config = config
        self.rigs: Dict[str, RigManager] = {}
        self.beads = BeadStore(config.root / ".beads")
        self.deacon: Optional[DeaconService] = None

    async def add_rig(self, name: str, repo_url: str) -> RigManager:
        """Add a new rig to the town."""
        rig_path = self.config.root / name
        rig = RigManager(name=name, path=rig_path, repo_url=repo_url)
        await rig.initialize()
        self.rigs[name] = rig

        # Register prefix routing
        self.beads.register_route(rig.prefix, rig.beads_path)
        return rig

    async def start_deacon(self):
        """Start the Deacon daemon service."""
        self.deacon = DeaconService(self)
        await self.deacon.start()

    async def status(self) -> dict:
        """Get town-wide status."""
        return {
            "rigs": {
                name: await rig.status()
                for name, rig in self.rigs.items()
            },
            "deacon_running": self.deacon is not None and self.deacon.running,
            "total_polecats": sum(
                len(rig.polecats) for rig in self.rigs.values()
            ),
        }
```

### 2. RigManager

Per-project container managing witness, refinery, and polecats.

```python
# vermas/core/rig.py
from pathlib import Path
from typing import Dict, List, Optional
import asyncio

from pydantic import BaseModel
from vermas.agents.witness import WitnessAgent
from vermas.agents.refinery import RefineryAgent
from vermas.agents.polecat import PolecatAgent
from vermas.beads.store import BeadStore
from vermas.tmux.session import TmuxSession


class RigConfig(BaseModel):
    """Rig configuration."""
    name: str
    path: Path
    repo_url: str
    prefix: str  # Beads ID prefix (e.g., "gt")
    max_polecats: int = 5


class RigManager:
    """
    Per-rig manager.

    Manages:
    - Witness agent (monitors polecats)
    - Refinery agent (processes merges)
    - Polecat pool (ephemeral workers)
    - Rig-level beads
    """

    def __init__(self, config: RigConfig):
        self.config = config
        self.beads = BeadStore(config.path / ".beads")
        self.witness: Optional[WitnessAgent] = None
        self.refinery: Optional[RefineryAgent] = None
        self.polecats: Dict[str, PolecatAgent] = {}
        self._available_slots: List[str] = []

    async def initialize(self):
        """Initialize rig directory structure."""
        dirs = ["polecats", "crew", "refinery", "witness", "mayor/rig"]
        for d in dirs:
            (self.config.path / d).mkdir(parents=True, exist_ok=True)

        # Initialize slot pool
        self._available_slots = [f"slot{i}" for i in range(self.config.max_polecats)]

    async def start_witness(self):
        """Start the Witness agent for this rig."""
        self.witness = WitnessAgent(rig=self)
        await self.witness.start()

    async def start_refinery(self):
        """Start the Refinery agent for this rig."""
        self.refinery = RefineryAgent(rig=self)
        await self.refinery.start()

    async def spawn_polecat(self, bead_id: str) -> PolecatAgent:
        """
        Spawn a new polecat worker.

        Three-layer lifecycle:
        1. Slot: Allocate from pool
        2. Sandbox: Create git worktree
        3. Session: Launch tmux with Claude Code
        """
        if not self._available_slots:
            raise RuntimeError("No available polecat slots")

        slot = self._available_slots.pop(0)
        polecat = PolecatAgent(
            rig=self,
            slot=slot,
            bead_id=bead_id,
        )
        await polecat.spawn()
        self.polecats[slot] = polecat
        return polecat

    async def release_polecat(self, slot: str):
        """Release a polecat slot back to the pool."""
        if slot in self.polecats:
            await self.polecats[slot].cleanup()
            del self.polecats[slot]
            self._available_slots.append(slot)

    async def status(self) -> dict:
        """Get rig status."""
        return {
            "name": self.config.name,
            "witness_running": self.witness is not None and self.witness.running,
            "refinery_running": self.refinery is not None and self.refinery.running,
            "polecats": {
                slot: pc.status for slot, pc in self.polecats.items()
            },
            "available_slots": len(self._available_slots),
        }
```

### 3. TmuxSessionManager

Manages tmux sessions for Claude Code agents.

```python
# vermas/tmux/session.py
import libtmux
from pathlib import Path
from typing import Optional
from pydantic import BaseModel


class SessionConfig(BaseModel):
    """Tmux session configuration."""
    name: str
    working_dir: Path
    claude_profile: str
    env: dict = {}


class TmuxSessionManager:
    """
    Manages tmux sessions for Claude Code agents.

    Each agent runs in its own tmux session with:
    - Dedicated working directory
    - Claude Code profile (--profile flag)
    - Environment variables for identity
    """

    def __init__(self):
        self.server = libtmux.Server()

    def create_session(self, config: SessionConfig) -> libtmux.Session:
        """
        Create a new tmux session with Claude Code.

        The session runs:
        claude --profile {profile} --dangerously-skip-permissions
        """
        # Build environment
        env_exports = " ".join(
            f'{k}="{v}"' for k, v in config.env.items()
        )

        # Create session
        session = self.server.new_session(
            session_name=config.name,
            start_directory=str(config.working_dir),
            attach=False,
        )

        # Launch Claude Code with profile
        cmd = f"{env_exports} claude --profile {config.claude_profile}"
        if config.env.get("DANGEROUSLY_SKIP_PERMISSIONS"):
            cmd += " --dangerously-skip-permissions"

        session.active_window.active_pane.send_keys(cmd)
        return session

    def get_session(self, name: str) -> Optional[libtmux.Session]:
        """Get existing session by name."""
        try:
            return self.server.sessions.get(session_name=name)
        except Exception:
            return None

    def kill_session(self, name: str) -> bool:
        """Kill a session by name."""
        session = self.get_session(name)
        if session:
            session.kill()
            return True
        return False

    def send_prompt(self, session_name: str, prompt: str):
        """
        Send a prompt to a running Claude Code session.

        Uses tmux send-keys to inject the prompt.
        """
        session = self.get_session(session_name)
        if session:
            pane = session.active_window.active_pane
            # Escape special characters
            escaped = prompt.replace("'", "'\\''")
            pane.send_keys(escaped)

    def list_sessions(self, prefix: str = "") -> list:
        """List all sessions, optionally filtered by prefix."""
        sessions = [s.name for s in self.server.sessions]
        if prefix:
            sessions = [s for s in sessions if s.startswith(prefix)]
        return sessions
```

---

## Beads Integration

### BeadStore

JSONL-based persistence compatible with Gas Town beads.

```python
# vermas/beads/store.py
import json
from pathlib import Path
from typing import List, Optional, Iterator
from datetime import datetime
import hashlib

from vermas.models.bead import Bead, BeadStatus, BeadType


class BeadStore:
    """
    JSONL-based bead storage.

    Compatible with Gas Town's .beads/issues.jsonl format.
    """

    def __init__(self, beads_dir: Path):
        self.beads_dir = beads_dir
        self.issues_file = beads_dir / "issues.jsonl"
        self.messages_file = beads_dir / "messages.jsonl"
        self._ensure_dirs()

    def _ensure_dirs(self):
        """Create beads directory structure."""
        self.beads_dir.mkdir(parents=True, exist_ok=True)
        (self.beads_dir / "formulas").mkdir(exist_ok=True)
        (self.beads_dir / "mols").mkdir(exist_ok=True)

    def _generate_id(self, prefix: str) -> str:
        """Generate unique bead ID."""
        timestamp = datetime.utcnow().isoformat()
        hash_input = f"{timestamp}-{id(self)}"
        short_hash = hashlib.sha256(hash_input.encode()).hexdigest()[:5]
        return f"{prefix}-{short_hash}"

    def create(self, bead: Bead) -> Bead:
        """Create a new bead."""
        if not bead.id:
            bead.id = self._generate_id(bead.prefix or "gt")
        bead.created_at = datetime.utcnow()
        bead.updated_at = bead.created_at

        with open(self.issues_file, "a") as f:
            f.write(bead.model_dump_json() + "\n")

        return bead

    def get(self, bead_id: str) -> Optional[Bead]:
        """Get bead by ID."""
        for bead in self._iter_beads():
            if bead.id == bead_id:
                return bead
        return None

    def update(self, bead: Bead) -> Bead:
        """Update existing bead."""
        bead.updated_at = datetime.utcnow()
        beads = list(self._iter_beads())

        # Replace the bead
        for i, b in enumerate(beads):
            if b.id == bead.id:
                beads[i] = bead
                break

        # Rewrite file
        self._write_all(beads)
        return bead

    def list(
        self,
        status: Optional[BeadStatus] = None,
        bead_type: Optional[BeadType] = None,
    ) -> List[Bead]:
        """List beads with optional filters."""
        beads = list(self._iter_beads())

        if status:
            beads = [b for b in beads if b.status == status]
        if bead_type:
            beads = [b for b in beads if b.issue_type == bead_type]

        return beads

    def ready(self) -> List[Bead]:
        """Get beads ready to work (open, no blockers)."""
        all_beads = list(self._iter_beads())
        blocked_ids = set()

        # Find all blocked beads
        for bead in all_beads:
            for dep in bead.dependencies:
                if dep.type == "blocks":
                    blocked_ids.add(bead.id)

        return [
            b for b in all_beads
            if b.status == BeadStatus.OPEN and b.id not in blocked_ids
        ]

    def _iter_beads(self) -> Iterator[Bead]:
        """Iterate over all beads."""
        if not self.issues_file.exists():
            return

        with open(self.issues_file) as f:
            for line in f:
                line = line.strip()
                if line:
                    yield Bead.model_validate_json(line)

    def _write_all(self, beads: List[Bead]):
        """Write all beads to file."""
        with open(self.issues_file, "w") as f:
            for bead in beads:
                f.write(bead.model_dump_json() + "\n")
```

---

## Hook System

Implements GUPP (Gas Town Universal Propulsion Principle).

```python
# vermas/core/hooks.py
from pathlib import Path
from typing import Optional, Callable, Awaitable
from enum import Enum

from vermas.models.bead import Bead
from vermas.beads.store import BeadStore


class HookType(str, Enum):
    """Types of hooks."""
    BEAD = "bead"      # Work bead hooked
    MAIL = "mail"      # Mail message hooked
    MOL = "mol"        # Molecule workflow hooked


class Hook:
    """
    Agent hook - where work hangs.

    GUPP: If your hook has work, RUN IT.
    """

    def __init__(self, agent_id: str, beads: BeadStore):
        self.agent_id = agent_id
        self.beads = beads
        self._hook_file = beads.beads_dir / f".hook-{agent_id}"

    def check(self) -> Optional[Bead]:
        """
        Check for hooked work.

        Returns the hooked bead if any, None otherwise.
        """
        if not self._hook_file.exists():
            return None

        content = self._hook_file.read_text().strip()
        if not content:
            return None

        # Parse hook file: type:id
        hook_type, hook_id = content.split(":", 1)

        if hook_type == HookType.BEAD:
            return self.beads.get(hook_id)

        return None

    def attach(self, bead: Bead):
        """Attach a bead to the hook."""
        self._hook_file.write_text(f"{HookType.BEAD}:{bead.id}")

    def clear(self):
        """Clear the hook."""
        if self._hook_file.exists():
            self._hook_file.unlink()

    async def run_if_hooked(
        self,
        executor: Callable[[Bead], Awaitable[None]]
    ) -> bool:
        """
        GUPP implementation: Run executor if work is hooked.

        Returns True if work was found and executed.
        """
        bead = self.check()
        if bead:
            await executor(bead)
            return True
        return False
```

---

## Agent Identity

BD_ACTOR format for attribution.

```python
# vermas/models/agent.py
from pydantic import BaseModel, Field
from typing import Optional
from enum import Enum


class AgentRole(str, Enum):
    """Agent roles in Gas Town."""
    MAYOR = "mayor"
    DEACON = "deacon"
    WITNESS = "witness"
    REFINERY = "refinery"
    POLECAT = "polecat"
    CREW = "crew"
    # VerMAS Inspector roles
    DESIGNER = "designer"
    STRATEGIST = "strategist"
    VERIFIER = "verifier"
    AUDITOR = "auditor"
    ADVOCATE = "advocate"
    CRITIC = "critic"
    JUDGE = "judge"


class AgentIdentity(BaseModel):
    """
    Agent identity in BD_ACTOR format.

    Format: {rig}/{role}/{name} or {role} for town-level

    Examples:
    - mayor (town-level)
    - gastown/witness (per-rig)
    - gastown/polecats/slot1 (polecat)
    - gastown/crew/frontend (crew member)
    """

    rig: Optional[str] = None
    role: AgentRole
    name: Optional[str] = None

    @property
    def actor_string(self) -> str:
        """Get BD_ACTOR format string."""
        parts = []
        if self.rig:
            parts.append(self.rig)
        parts.append(self.role.value)
        if self.name:
            parts.append(self.name)
        return "/".join(parts)

    @classmethod
    def from_string(cls, s: str) -> "AgentIdentity":
        """Parse BD_ACTOR string."""
        parts = s.split("/")

        if len(parts) == 1:
            return cls(role=AgentRole(parts[0]))
        elif len(parts) == 2:
            return cls(rig=parts[0], role=AgentRole(parts[1]))
        else:
            return cls(rig=parts[0], role=AgentRole(parts[1]), name=parts[2])


class AgentConfig(BaseModel):
    """Agent configuration."""
    identity: AgentIdentity
    working_dir: Path
    claude_profile: str
    env: dict = Field(default_factory=dict)

    @property
    def session_name(self) -> str:
        """Get tmux session name."""
        return self.identity.actor_string.replace("/", "-")
```

---

## See Also

- [AGENTS.md](./AGENTS.md) - Detailed agent implementations
- [WORKFLOWS.md](./WORKFLOWS.md) - Molecule workflow system
- [MESSAGING.md](./MESSAGING.md) - Mail protocol
- [LIFECYCLE.md](./LIFECYCLE.md) - Polecat lifecycle management
