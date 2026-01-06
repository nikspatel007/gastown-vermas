# VerMAS Python Lifecycle

> Polecat lifecycle, session management, and watchdog chain

## Polecat Three-Layer Architecture

Polecats have three distinct lifecycle layers:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      POLECAT THREE-LAYER LIFECYCLE                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │  LAYER 1: SESSION (Ephemeral)                                       │  │
│   │                                                                     │  │
│   │  - Tmux session running Claude Code                                 │  │
│   │  - Lives for duration of single task                                │  │
│   │  - Dies when work completes or is killed                           │  │
│   │  - Session name: polecat-{rig}-{slot}                              │  │
│   │                                                                     │  │
│   │  State: Running / Idle / Stuck / Dead                              │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                              │                                              │
│                              ▼                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │  LAYER 2: SANDBOX (Persistent)                                      │  │
│   │                                                                     │  │
│   │  - Git worktree for isolated work                                  │  │
│   │  - Survives session death                                          │  │
│   │  - Contains uncommitted changes                                    │  │
│   │  - Path: {rig}/polecats/{slot}/                                    │  │
│   │                                                                     │  │
│   │  State: Clean / Dirty / Conflicted                                 │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                              │                                              │
│                              ▼                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │  LAYER 3: SLOT (Allocation)                                         │  │
│   │                                                                     │  │
│   │  - Name from pool: slot0, slot1, slot2, etc.                       │  │
│   │  - Limited pool per rig (default: 5)                               │  │
│   │  - Allocated on spawn, released on cleanup                         │  │
│   │  - Prevents resource exhaustion                                    │  │
│   │                                                                     │  │
│   │  State: Available / Allocated                                      │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

Lifecycle: spawn(slot) → create_worktree(sandbox) → start_session(session)
                                                            ↓
           release(slot) ← remove_worktree(sandbox) ← kill_session(session)
```

---

## Python Implementation

### Slot Pool Manager

```python
# vermas/lifecycle/slots.py
from typing import List, Optional, Set
from threading import Lock


class SlotPool:
    """
    Manages polecat slot allocation.

    Prevents resource exhaustion by limiting concurrent polecats.
    Thread-safe for concurrent allocation/release.
    """

    def __init__(self, max_slots: int = 5, prefix: str = "slot"):
        self.max_slots = max_slots
        self.prefix = prefix
        self._available: Set[str] = {
            f"{prefix}{i}" for i in range(max_slots)
        }
        self._allocated: Set[str] = set()
        self._lock = Lock()

    def allocate(self) -> Optional[str]:
        """
        Allocate a slot from the pool.

        Returns slot name or None if pool exhausted.
        """
        with self._lock:
            if not self._available:
                return None
            slot = self._available.pop()
            self._allocated.add(slot)
            return slot

    def release(self, slot: str):
        """Release a slot back to the pool."""
        with self._lock:
            if slot in self._allocated:
                self._allocated.remove(slot)
                self._available.add(slot)

    @property
    def available_count(self) -> int:
        return len(self._available)

    @property
    def allocated_count(self) -> int:
        return len(self._allocated)

    def list_available(self) -> List[str]:
        return sorted(self._available)

    def list_allocated(self) -> List[str]:
        return sorted(self._allocated)
```

### Sandbox Manager

```python
# vermas/lifecycle/sandbox.py
import asyncio
import subprocess
from pathlib import Path
from typing import Optional
from enum import Enum


class SandboxState(str, Enum):
    """Sandbox (worktree) state."""
    CLEAN = "clean"
    DIRTY = "dirty"
    CONFLICTED = "conflicted"
    MISSING = "missing"


class Sandbox:
    """
    Git worktree sandbox for isolated polecat work.

    Each polecat gets its own worktree to prevent conflicts.
    """

    def __init__(self, rig_path: Path, slot: str, branch_prefix: str = "polecat"):
        self.rig_path = rig_path
        self.slot = slot
        self.branch_name = f"{branch_prefix}-{slot}"
        self.worktree_path = rig_path / "polecats" / slot

    async def create(self, base_branch: str = "main") -> bool:
        """
        Create the worktree sandbox.

        Creates a new branch and worktree for this polecat.
        """
        # Ensure polecats directory exists
        self.worktree_path.parent.mkdir(parents=True, exist_ok=True)

        # Create worktree with new branch
        cmd = [
            "git", "worktree", "add",
            str(self.worktree_path),
            "-b", self.branch_name,
            base_branch,
        ]

        try:
            process = await asyncio.create_subprocess_exec(
                *cmd,
                cwd=str(self.rig_path / "mayor" / "rig"),  # Source repo
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE,
            )
            await process.communicate()
            return process.returncode == 0
        except Exception:
            return False

    async def remove(self, force: bool = False) -> bool:
        """
        Remove the worktree sandbox.

        Warning: Uncommitted changes will be lost if force=True.
        """
        cmd = ["git", "worktree", "remove", str(self.worktree_path)]
        if force:
            cmd.append("--force")

        try:
            process = await asyncio.create_subprocess_exec(
                *cmd,
                cwd=str(self.rig_path / "mayor" / "rig"),
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE,
            )
            await process.communicate()

            # Also delete the branch
            if process.returncode == 0:
                await self._delete_branch()

            return process.returncode == 0
        except Exception:
            return False

    async def _delete_branch(self):
        """Delete the polecat branch."""
        cmd = ["git", "branch", "-D", self.branch_name]
        process = await asyncio.create_subprocess_exec(
            *cmd,
            cwd=str(self.rig_path / "mayor" / "rig"),
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
        )
        await process.communicate()

    async def get_state(self) -> SandboxState:
        """Check sandbox state."""
        if not self.worktree_path.exists():
            return SandboxState.MISSING

        # Check for uncommitted changes
        cmd = ["git", "status", "--porcelain"]
        process = await asyncio.create_subprocess_exec(
            *cmd,
            cwd=str(self.worktree_path),
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
        )
        stdout, _ = await process.communicate()

        if stdout:
            output = stdout.decode()
            if any(line.startswith("UU") for line in output.splitlines()):
                return SandboxState.CONFLICTED
            return SandboxState.DIRTY
        return SandboxState.CLEAN

    async def has_uncommitted_changes(self) -> bool:
        """Check if sandbox has uncommitted changes."""
        state = await self.get_state()
        return state in (SandboxState.DIRTY, SandboxState.CONFLICTED)
```

### Session Manager

```python
# vermas/lifecycle/session.py
import asyncio
from datetime import datetime
from typing import Optional
from enum import Enum
import libtmux

from vermas.models.agent import AgentIdentity, AgentRole


class SessionState(str, Enum):
    """Tmux session state."""
    RUNNING = "running"
    IDLE = "idle"
    STUCK = "stuck"
    DEAD = "dead"


class PolecatSession:
    """
    Tmux session for a polecat worker.

    Manages the ephemeral Claude Code session.
    """

    def __init__(
        self,
        rig_name: str,
        slot: str,
        working_dir: Path,
        bead_id: str,
    ):
        self.rig_name = rig_name
        self.slot = slot
        self.working_dir = working_dir
        self.bead_id = bead_id
        self.session_name = f"polecat-{rig_name}-{slot}"
        self._server = libtmux.Server()
        self._session: Optional[libtmux.Session] = None
        self._started_at: Optional[datetime] = None
        self._last_activity: Optional[datetime] = None

    @property
    def identity(self) -> AgentIdentity:
        return AgentIdentity(
            rig=self.rig_name,
            role=AgentRole.POLECAT,
            name=self.slot,
        )

    async def start(self) -> bool:
        """
        Start the tmux session with Claude Code.

        Sets up environment and launches Claude with polecat profile.
        """
        try:
            # Create session
            self._session = self._server.new_session(
                session_name=self.session_name,
                start_directory=str(self.working_dir),
                attach=False,
            )

            # Set environment
            pane = self._session.active_window.active_pane
            pane.send_keys(f"export BD_ACTOR='{self.identity.actor_string}'")
            pane.send_keys(f"export BEAD_ID='{self.bead_id}'")

            # Launch Claude Code
            pane.send_keys("claude --profile polecat")

            self._started_at = datetime.utcnow()
            self._last_activity = self._started_at
            return True

        except Exception as e:
            print(f"Failed to start session: {e}")
            return False

    async def stop(self) -> bool:
        """Kill the tmux session."""
        try:
            if self._session:
                self._session.kill()
            return True
        except Exception:
            return False

    def send_prompt(self, prompt: str):
        """Send a prompt to the Claude Code session."""
        if self._session:
            pane = self._session.active_window.active_pane
            # Escape for tmux
            escaped = prompt.replace("'", "'\\''")
            pane.send_keys(escaped)
            self._last_activity = datetime.utcnow()

    async def get_state(self) -> SessionState:
        """Determine session state."""
        # Check if session exists
        try:
            session = self._server.sessions.get(session_name=self.session_name)
            if not session:
                return SessionState.DEAD
        except Exception:
            return SessionState.DEAD

        # Check activity
        if not self._last_activity:
            return SessionState.RUNNING

        idle_seconds = (datetime.utcnow() - self._last_activity).total_seconds()

        if idle_seconds > 900:  # 15 min
            return SessionState.STUCK
        elif idle_seconds > 300:  # 5 min
            return SessionState.IDLE
        else:
            return SessionState.RUNNING

    @property
    def idle_seconds(self) -> float:
        if not self._last_activity:
            return 0
        return (datetime.utcnow() - self._last_activity).total_seconds()
```

### Polecat Lifecycle Manager

```python
# vermas/lifecycle/polecat.py
import asyncio
from pathlib import Path
from typing import Dict, Optional
from datetime import datetime

from vermas.lifecycle.slots import SlotPool
from vermas.lifecycle.sandbox import Sandbox, SandboxState
from vermas.lifecycle.session import PolecatSession, SessionState
from vermas.models.bead import Bead, BeadStatus
from vermas.beads.store import BeadStore
from vermas.core.hooks import Hook


class PolecatLifecycle:
    """
    Complete lifecycle manager for polecats.

    Coordinates: Slot → Sandbox → Session
    """

    def __init__(self, rig_name: str, rig_path: Path, max_polecats: int = 5):
        self.rig_name = rig_name
        self.rig_path = rig_path
        self.slots = SlotPool(max_polecats)
        self.beads = BeadStore(rig_path / ".beads")
        self._polecats: Dict[str, PolecatInstance] = {}

    async def spawn(self, bead: Bead) -> Optional["PolecatInstance"]:
        """
        Spawn a new polecat for a bead.

        1. Allocate slot
        2. Create sandbox (worktree)
        3. Start session
        4. Hook the bead
        """
        # 1. Allocate slot
        slot = self.slots.allocate()
        if not slot:
            return None  # Pool exhausted

        try:
            # 2. Create sandbox
            sandbox = Sandbox(self.rig_path, slot)
            if not await sandbox.create():
                self.slots.release(slot)
                return None

            # 3. Start session
            session = PolecatSession(
                rig_name=self.rig_name,
                slot=slot,
                working_dir=sandbox.worktree_path,
                bead_id=bead.id,
            )
            if not await session.start():
                await sandbox.remove(force=True)
                self.slots.release(slot)
                return None

            # 4. Hook the bead
            hook = Hook(session.identity.actor_string, self.rig_path / ".beads")
            hook.attach(bead)
            bead.status = BeadStatus.HOOKED
            bead.assigned_to = session.identity.actor_string
            self.beads.update(bead)

            # Create instance
            instance = PolecatInstance(
                slot=slot,
                sandbox=sandbox,
                session=session,
                bead_id=bead.id,
            )
            self._polecats[slot] = instance

            # Send startup prompt (GUPP)
            session.send_prompt(self._build_startup_prompt(bead))

            return instance

        except Exception as e:
            # Cleanup on failure
            self.slots.release(slot)
            raise

    async def cleanup(self, slot: str, force: bool = False):
        """
        Clean up a polecat.

        1. Kill session
        2. Remove sandbox (unless dirty and not force)
        3. Release slot
        """
        instance = self._polecats.get(slot)
        if not instance:
            return

        # 1. Kill session
        await instance.session.stop()

        # 2. Check sandbox state
        state = await instance.sandbox.get_state()
        if state == SandboxState.DIRTY and not force:
            # Don't remove dirty sandbox - might have work
            print(f"Warning: Sandbox {slot} has uncommitted changes")
        else:
            await instance.sandbox.remove(force=force)

        # 3. Release slot
        self.slots.release(slot)
        del self._polecats[slot]

    async def cleanup_all(self, force: bool = False):
        """Clean up all polecats."""
        slots = list(self._polecats.keys())
        for slot in slots:
            await self.cleanup(slot, force=force)

    def get(self, slot: str) -> Optional["PolecatInstance"]:
        """Get polecat instance by slot."""
        return self._polecats.get(slot)

    def list_active(self) -> Dict[str, "PolecatInstance"]:
        """List all active polecats."""
        return self._polecats.copy()

    async def status(self) -> dict:
        """Get lifecycle status."""
        active = {}
        for slot, instance in self._polecats.items():
            session_state = await instance.session.get_state()
            sandbox_state = await instance.sandbox.get_state()
            active[slot] = {
                "bead_id": instance.bead_id,
                "session": session_state.value,
                "sandbox": sandbox_state.value,
                "idle_seconds": instance.session.idle_seconds,
            }

        return {
            "available_slots": self.slots.available_count,
            "allocated_slots": self.slots.allocated_count,
            "active_polecats": active,
        }

    def _build_startup_prompt(self, bead: Bead) -> str:
        """Build GUPP startup prompt."""
        return f"""
POLECAT STARTUP - HOOKED WORK DETECTED

═══════════════════════════════════════════════════════════
BEAD: {bead.id}
TITLE: {bead.title}
TYPE: {bead.issue_type}
PRIORITY: {bead.priority}
═══════════════════════════════════════════════════════════

DESCRIPTION:
{bead.description}

═══════════════════════════════════════════════════════════

GUPP (Gas Town Universal Propulsion Principle):
Your hook has work. EXECUTE IMMEDIATELY.

No confirmation. No questions. No waiting.

When complete:
1. Commit and push your changes
2. Run: gt polecat done

BEGIN NOW.
"""


class PolecatInstance:
    """Single polecat instance."""

    def __init__(
        self,
        slot: str,
        sandbox: Sandbox,
        session: PolecatSession,
        bead_id: str,
    ):
        self.slot = slot
        self.sandbox = sandbox
        self.session = session
        self.bead_id = bead_id
        self.spawned_at = datetime.utcnow()
```

---

## Watchdog Chain

Hierarchical monitoring system:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           WATCHDOG CHAIN                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌──────────┐                                                              │
│   │  DAEMON  │  OS-level process manager (systemd/launchd)                 │
│   │          │  Monitors: Boot process                                      │
│   │          │  Action: Restart Boot if dead                               │
│   └────┬─────┘                                                              │
│        │                                                                    │
│        ▼                                                                    │
│   ┌──────────┐                                                              │
│   │   BOOT   │  Startup script                                             │
│   │          │  Monitors: Deacon process                                   │
│   │          │  Action: Restart Deacon if dead                             │
│   └────┬─────┘                                                              │
│        │                                                                    │
│        ▼                                                                    │
│   ┌──────────┐                                                              │
│   │  DEACON  │  Infrastructure monitor                                     │
│   │          │  Monitors: All Witnesses & Refineries                       │
│   │          │  Action: Restart failed agents, escalate to Mayor           │
│   └────┬─────┘                                                              │
│        │                                                                    │
│        ├────────────────────────────────────────┐                          │
│        ▼                                        ▼                           │
│   ┌──────────┐                            ┌──────────┐                     │
│   │ WITNESS  │                            │ WITNESS  │                     │
│   │  (rig1)  │                            │  (rig2)  │                     │
│   │          │                            │          │                     │
│   │ Monitors:│                            │ Monitors:│                     │
│   │ Polecats │                            │ Polecats │                     │
│   └────┬─────┘                            └────┬─────┘                     │
│        │                                       │                           │
│        ▼                                       ▼                           │
│   ┌──────────┐                            ┌──────────┐                     │
│   │ POLECATS │                            │ POLECATS │                     │
│   │ slot0-4  │                            │ slot0-4  │                     │
│   └──────────┘                            └──────────┘                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Watchdog Implementation

```python
# vermas/lifecycle/watchdog.py
import asyncio
from typing import Dict, List, Callable, Awaitable
from datetime import datetime


class WatchdogLevel:
    """Single level in the watchdog chain."""

    def __init__(
        self,
        name: str,
        check_interval: int = 60,
        restart_action: Callable[[], Awaitable[bool]] = None,
    ):
        self.name = name
        self.check_interval = check_interval
        self.restart_action = restart_action
        self._last_check = None
        self._targets: Dict[str, WatchedTarget] = {}

    def add_target(self, name: str, health_check: Callable[[], Awaitable[bool]]):
        """Add a target to watch."""
        self._targets[name] = WatchedTarget(name, health_check)

    async def patrol(self) -> List[str]:
        """
        Check all targets.

        Returns list of failed targets that were restarted.
        """
        self._last_check = datetime.utcnow()
        restarted = []

        for name, target in self._targets.items():
            healthy = await target.check()
            if not healthy:
                if self.restart_action:
                    success = await self.restart_action()
                    if success:
                        restarted.append(name)

        return restarted


class WatchedTarget:
    """Target being watched."""

    def __init__(self, name: str, health_check: Callable[[], Awaitable[bool]]):
        self.name = name
        self.health_check = health_check
        self.last_healthy = None
        self.consecutive_failures = 0

    async def check(self) -> bool:
        """Check target health."""
        try:
            healthy = await self.health_check()
            if healthy:
                self.last_healthy = datetime.utcnow()
                self.consecutive_failures = 0
            else:
                self.consecutive_failures += 1
            return healthy
        except Exception:
            self.consecutive_failures += 1
            return False


class WatchdogChain:
    """
    Complete watchdog chain implementation.

    Implements: Daemon → Boot → Deacon → Witnesses → Polecats
    """

    def __init__(self):
        self.levels: Dict[str, WatchdogLevel] = {}
        self._running = False

    def add_level(self, level: WatchdogLevel):
        """Add a watchdog level."""
        self.levels[level.name] = level

    async def start(self):
        """Start the watchdog chain."""
        self._running = True
        asyncio.create_task(self._patrol_loop())

    async def stop(self):
        """Stop the watchdog chain."""
        self._running = False

    async def _patrol_loop(self):
        """Main patrol loop."""
        while self._running:
            for level_name, level in self.levels.items():
                try:
                    restarted = await level.patrol()
                    if restarted:
                        print(f"Watchdog [{level_name}]: Restarted {restarted}")
                except Exception as e:
                    print(f"Watchdog [{level_name}] error: {e}")

            await asyncio.sleep(30)  # Base interval
```

---

## CLI Commands

```python
# vermas/cli.py (lifecycle commands)
import typer
from pathlib import Path
from rich.console import Console
from rich.table import Table

lifecycle_app = typer.Typer()
console = Console()


@lifecycle_app.command("spawn")
def spawn_polecat(bead_id: str, rig: str = None):
    """Spawn a new polecat for a bead."""
    import asyncio
    from vermas.lifecycle.polecat import PolecatLifecycle
    from vermas.beads.store import BeadStore

    rig_path = Path(rig) if rig else Path.cwd()
    beads = BeadStore(rig_path / ".beads")
    bead = beads.get(bead_id)

    if not bead:
        console.print(f"[red]Bead not found: {bead_id}[/red]")
        raise typer.Exit(1)

    lifecycle = PolecatLifecycle(rig_path.name, rig_path)
    instance = asyncio.run(lifecycle.spawn(bead))

    if instance:
        console.print(f"[green]Spawned polecat: {instance.slot}[/green]")
        console.print(f"  Sandbox: {instance.sandbox.worktree_path}")
        console.print(f"  Session: {instance.session.session_name}")
    else:
        console.print("[red]Failed to spawn polecat (pool may be exhausted)[/red]")


@lifecycle_app.command("list")
def list_polecats(rig: str = None):
    """List active polecats."""
    import asyncio
    from vermas.lifecycle.polecat import PolecatLifecycle

    rig_path = Path(rig) if rig else Path.cwd()
    lifecycle = PolecatLifecycle(rig_path.name, rig_path)
    status = asyncio.run(lifecycle.status())

    table = Table(title=f"Polecats in {rig_path.name}")
    table.add_column("Slot")
    table.add_column("Bead")
    table.add_column("Session")
    table.add_column("Sandbox")
    table.add_column("Idle")

    for slot, info in status["active_polecats"].items():
        idle_str = f"{int(info['idle_seconds'])}s"
        table.add_row(
            slot,
            info["bead_id"],
            info["session"],
            info["sandbox"],
            idle_str,
        )

    console.print(table)
    console.print(f"\nAvailable slots: {status['available_slots']}")


@lifecycle_app.command("kill")
def kill_polecat(slot: str, force: bool = False, rig: str = None):
    """Kill a polecat."""
    import asyncio
    from vermas.lifecycle.polecat import PolecatLifecycle

    rig_path = Path(rig) if rig else Path.cwd()
    lifecycle = PolecatLifecycle(rig_path.name, rig_path)

    if not force:
        typer.confirm(f"Kill polecat {slot}?", abort=True)

    asyncio.run(lifecycle.cleanup(slot, force=force))
    console.print(f"[yellow]Killed polecat: {slot}[/yellow]")


@lifecycle_app.command("done")
def polecat_done():
    """Signal polecat completion (run from within polecat session)."""
    import os
    from vermas.mail.store import MailStore
    from vermas.models.mail import MessageType

    actor = os.environ.get("BD_ACTOR")
    bead_id = os.environ.get("BEAD_ID")

    if not actor or "polecat" not in actor:
        console.print("[red]This command must be run from a polecat session[/red]")
        raise typer.Exit(1)

    # Extract rig and slot from actor
    parts = actor.split("/")
    rig = parts[0]
    slot = parts[2] if len(parts) > 2 else "unknown"

    # Send POLECAT_DONE to Witness
    store = MailStore(Path(".beads"))
    store.send(
        from_addr=actor,
        to_addr=f"{rig}/witness",
        subject=f"POLECAT_DONE: {bead_id}",
        body="Work completed. Ready for merge.",
        message_type=MessageType.POLECAT_DONE,
        metadata={"bead_id": bead_id, "slot": slot},
    )

    console.print("[green]Completion signaled to Witness[/green]")
```

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [AGENTS.md](./AGENTS.md) - Agent implementations
- [WORKFLOWS.md](./WORKFLOWS.md) - Molecule workflow system
- [MESSAGING.md](./MESSAGING.md) - Mail protocol
