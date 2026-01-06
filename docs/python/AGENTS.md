# VerMAS Python Agents

> Detailed agent implementations for Gas Town + VerMAS

## Agent Taxonomy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            TOWN-WIDE AGENTS                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐                    │
│   │   MAYOR     │    │   DEACON    │    │  OVERSEER   │                    │
│   │             │    │             │    │   (Human)   │                    │
│   │ Cross-rig   │    │ Daemon      │    │             │                    │
│   │ coordinator │    │ lifecycle   │    │ Strategy    │                    │
│   │ NO code!    │    │ monitor     │    │ review      │                    │
│   └─────────────┘    └─────────────┘    └─────────────┘                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                            PER-RIG AGENTS                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌──────────┐   │
│   │  WITNESS    │    │  REFINERY   │    │  POLECAT    │    │   CREW   │   │
│   │             │    │             │    │             │    │          │   │
│   │ Worker      │    │ Merge queue │    │ Ephemeral   │    │ Human-   │   │
│   │ lifecycle   │    │ processor   │    │ worker      │    │ directed │   │
│   │ monitor     │    │             │    │ spawn→work  │    │ workspace│   │
│   └─────────────┘    └─────────────┘    │ →disappear  │    └──────────┘   │
│                                         └─────────────┘                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                       VERMAS INSPECTOR ECOSYSTEM                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Requirements         Verification         Adversarial Review              │
│   ────────────         ────────────         ─────────────────              │
│   ┌──────────┐        ┌───────────┐        ┌──────────┐                    │
│   │ DESIGNER │───────▶│ STRATEGIST│───────▶│ VERIFIER │                    │
│   │ Elaborate│        │ Plan tests│        │ Run tests│                    │
│   │ specs    │        │           │        │ (shell)  │                    │
│   └──────────┘        └───────────┘        └────┬─────┘                    │
│                                                  │                          │
│                       ┌──────────┐               ▼                          │
│                       │ AUDITOR  │◀──────[Evidence]                        │
│                       │ LLM check│                                          │
│                       │ (if shell│                                          │
│                       │  fails)  │                                          │
│                       └────┬─────┘                                          │
│                            │                                                │
│         ┌──────────────────┼──────────────────┐                            │
│         ▼                  ▼                  ▼                             │
│   ┌──────────┐       ┌──────────┐       ┌──────────┐                       │
│   │ ADVOCATE │       │  CRITIC  │       │  JUDGE   │                       │
│   │ Argue    │       │ Argue    │       │ Decide   │                       │
│   │ PASS     │       │ FAIL     │       │ verdict  │                       │
│   └──────────┘       └──────────┘       └──────────┘                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Base Agent

All agents inherit from `BaseAgent`.

```python
# vermas/agents/base.py
from abc import ABC, abstractmethod
from pathlib import Path
from typing import Optional
import asyncio

from vermas.models.agent import AgentIdentity, AgentConfig
from vermas.tmux.session import TmuxSessionManager, SessionConfig
from vermas.core.hooks import Hook
from vermas.beads.store import BeadStore


class BaseAgent(ABC):
    """
    Base class for all Gas Town agents.

    Implements:
    - Tmux session lifecycle
    - Hook checking (GUPP)
    - Beads access
    - Mail integration
    """

    def __init__(self, config: AgentConfig, beads: BeadStore):
        self.config = config
        self.beads = beads
        self.hook = Hook(config.identity.actor_string, beads)
        self.tmux = TmuxSessionManager()
        self.session = None
        self.running = False

    @property
    def identity(self) -> AgentIdentity:
        return self.config.identity

    @property
    def session_name(self) -> str:
        return self.config.session_name

    async def start(self):
        """
        Start the agent.

        1. Create tmux session
        2. Check hook (GUPP)
        3. If hooked, execute immediately
        4. Otherwise, enter patrol/wait mode
        """
        # Create session
        session_config = SessionConfig(
            name=self.session_name,
            working_dir=self.config.working_dir,
            claude_profile=self.config.claude_profile,
            env={
                "BD_ACTOR": self.identity.actor_string,
                **self.config.env,
            },
        )
        self.session = self.tmux.create_session(session_config)
        self.running = True

        # GUPP: Check hook and run if work exists
        await self.hook.run_if_hooked(self.execute_hooked_work)

        # If no hooked work, start default behavior
        if not self.hook.check():
            await self.on_idle()

    async def stop(self):
        """Stop the agent."""
        if self.session:
            self.tmux.kill_session(self.session_name)
        self.running = False

    @abstractmethod
    async def execute_hooked_work(self, bead):
        """Execute hooked work. Subclasses implement specific behavior."""
        pass

    @abstractmethod
    async def on_idle(self):
        """Called when no work is hooked. Subclasses implement patrol/wait."""
        pass

    def send_prompt(self, prompt: str):
        """Send a prompt to the Claude Code session."""
        if self.session:
            self.tmux.send_prompt(self.session_name, prompt)

    @property
    def status(self) -> dict:
        """Get agent status."""
        hooked = self.hook.check()
        return {
            "identity": self.identity.actor_string,
            "running": self.running,
            "hooked_work": hooked.id if hooked else None,
        }
```

---

## Town-Wide Agents

### Mayor

Cross-rig coordinator. Does NOT write code.

```python
# vermas/agents/mayor.py
from vermas.agents.base import BaseAgent
from vermas.models.bead import Bead


class MayorAgent(BaseAgent):
    """
    Mayor - Global coordinator.

    Responsibilities:
    - Cross-rig coordination
    - Work dispatch (sling beads to rigs)
    - Escalation handling
    - Strategic decisions

    DOES NOT:
    - Write code
    - Edit files in mayor/rig/
    - Directly manage polecats (Witness does that)
    """

    async def execute_hooked_work(self, bead: Bead):
        """
        Execute hooked work.

        Mayor work is typically coordination tasks:
        - Epic breakdown and dispatch
        - Cross-rig dependency resolution
        - Escalation handling
        """
        prompt = f"""
You have hooked work: {bead.id}
Title: {bead.title}
Description: {bead.description}

Execute this coordination task. Remember:
- You are the Mayor - coordinate, don't implement
- Dispatch work to appropriate rigs using `gt sling`
- Check convoy status with `gt convoy list`
- Do NOT edit code directly
"""
        self.send_prompt(prompt)

    async def on_idle(self):
        """
        Idle behavior: Check mail, await user instructions.

        Unlike polecats, Mayor doesn't patrol. It waits for:
        - User instructions
        - Escalations from Witnesses
        - Mail from other agents
        """
        prompt = """
No hooked work found. Checking mail...

Run: gt mail inbox

If you have mail, process it. Otherwise, await user instructions.
"""
        self.send_prompt(prompt)


# CLAUDE.md template for Mayor
MAYOR_TEMPLATE = """
# Mayor Context

## Your Role: MAYOR (Global Coordinator)

You are the **Mayor** - the global coordinator of Gas Town.

**Responsibilities:**
- Cross-rig coordination
- Work dispatch (gt sling)
- Escalation handling
- Strategic decisions

**NOT your job:**
- Writing code
- Editing files in mayor/rig/
- Direct polecat management (Witness handles that)

## Key Commands

### Coordination
- `gt sling <bead> <rig>` - Assign work to a rig
- `gt convoy list` - Dashboard of active work
- `gt status` - Town overview

### Communication
- `gt mail inbox` - Check messages
- `gt mail send <addr> -s "Subject" -m "Message"`

## GUPP (Propulsion Principle)

If your hook has work, RUN IT. No confirmation. No waiting.

Check: `gt hook`
"""
```

### Deacon

Daemon process managing agent lifecycle.

```python
# vermas/agents/deacon.py
import asyncio
from typing import Dict, List
from vermas.core.town import TownManager


class DeaconService:
    """
    Deacon - Daemon process for agent lifecycle.

    Monitors:
    - Witness health across all rigs
    - Refinery status
    - Infrastructure health

    Part of the watchdog chain:
    Daemon → Boot → Deacon → Witnesses → Polecats
    """

    def __init__(self, town: TownManager):
        self.town = town
        self.running = False
        self._patrol_interval = 60  # seconds

    async def start(self):
        """Start the Deacon patrol loop."""
        self.running = True
        asyncio.create_task(self._patrol_loop())

    async def stop(self):
        """Stop the Deacon."""
        self.running = False

    async def _patrol_loop(self):
        """
        Main patrol loop.

        Checks:
        1. Witness health in each rig
        2. Refinery queue status
        3. Stuck polecats (escalate if Witness can't handle)
        """
        while self.running:
            try:
                await self._patrol()
            except Exception as e:
                # Log error but keep running
                print(f"Deacon patrol error: {e}")

            await asyncio.sleep(self._patrol_interval)

    async def _patrol(self):
        """Single patrol iteration."""
        for rig_name, rig in self.town.rigs.items():
            # Check Witness
            if not rig.witness or not rig.witness.running:
                await self._restart_witness(rig)

            # Check Refinery
            if not rig.refinery or not rig.refinery.running:
                await self._restart_refinery(rig)

            # Check for stuck polecats that Witness couldn't handle
            stuck = await self._find_stuck_polecats(rig)
            if stuck:
                await self._escalate_stuck(rig, stuck)

    async def _restart_witness(self, rig):
        """Restart a failed Witness."""
        print(f"Deacon: Restarting Witness for {rig.config.name}")
        await rig.start_witness()

    async def _restart_refinery(self, rig):
        """Restart a failed Refinery."""
        print(f"Deacon: Restarting Refinery for {rig.config.name}")
        await rig.start_refinery()

    async def _find_stuck_polecats(self, rig) -> List[str]:
        """Find polecats stuck beyond Witness tolerance."""
        stuck = []
        for slot, polecat in rig.polecats.items():
            if polecat.is_stuck and polecat.stuck_duration > 1800:  # 30 min
                stuck.append(slot)
        return stuck

    async def _escalate_stuck(self, rig, stuck_slots: List[str]):
        """Escalate stuck polecats to Mayor."""
        for slot in stuck_slots:
            # Send mail to Mayor
            from vermas.mail.send import send_mail
            await send_mail(
                to="mayor",
                subject=f"ESCALATION: Stuck polecat {rig.config.name}/{slot}",
                body=f"Polecat {slot} in {rig.config.name} has been stuck for >30 min.",
            )
```

---

## Per-Rig Agents

### Witness

Monitors worker agents within a rig.

```python
# vermas/agents/witness.py
import asyncio
from typing import Dict, List
from vermas.agents.base import BaseAgent
from vermas.models.bead import Bead


class WitnessAgent(BaseAgent):
    """
    Witness - Per-rig worker monitor.

    Responsibilities:
    - Monitor polecat health
    - Detect stuck/idle workers
    - Nudge or kill unresponsive polecats
    - Escalate to Deacon if needed

    Patrol loop runs continuously when idle.
    """

    def __init__(self, rig, **kwargs):
        self.rig = rig
        config = self._build_config()
        super().__init__(config, rig.beads)
        self._patrol_interval = 30  # seconds
        self._nudge_threshold = 300  # 5 min idle
        self._kill_threshold = 900  # 15 min stuck

    def _build_config(self):
        from vermas.models.agent import AgentConfig, AgentIdentity, AgentRole
        return AgentConfig(
            identity=AgentIdentity(
                rig=self.rig.config.name,
                role=AgentRole.WITNESS,
            ),
            working_dir=self.rig.config.path / "witness",
            claude_profile="witness",
        )

    async def execute_hooked_work(self, bead: Bead):
        """
        Witnesses don't typically have hooked work.
        If they do, it's a patrol/check request.
        """
        await self._patrol()

    async def on_idle(self):
        """Start patrol loop when idle."""
        asyncio.create_task(self._patrol_loop())

    async def _patrol_loop(self):
        """Continuous patrol loop."""
        while self.running:
            await self._patrol()
            await asyncio.sleep(self._patrol_interval)

    async def _patrol(self):
        """
        Single patrol iteration.

        1. Check mail for POLECAT_DONE messages
        2. Survey all active polecats
        3. Nudge idle workers
        4. Kill stuck workers
        5. Report to Deacon if needed
        """
        # Process mail
        await self._process_mail()

        # Survey workers
        for slot, polecat in self.rig.polecats.items():
            await self._check_polecat(slot, polecat)

    async def _process_mail(self):
        """Process incoming mail."""
        # Check for POLECAT_DONE messages
        from vermas.mail.inbox import get_inbox
        messages = await get_inbox(self.identity.actor_string)

        for msg in messages:
            if msg.subject.startswith("POLECAT_DONE:"):
                await self._handle_polecat_done(msg)

    async def _handle_polecat_done(self, msg):
        """Handle polecat completion."""
        # Extract slot from message
        slot = msg.metadata.get("slot")
        if slot and slot in self.rig.polecats:
            polecat = self.rig.polecats[slot]
            # Forward to Refinery for merge
            await self._notify_refinery(polecat)
            # Clean up polecat
            await self.rig.release_polecat(slot)

    async def _check_polecat(self, slot: str, polecat):
        """Check individual polecat health."""
        if polecat.is_idle:
            idle_time = polecat.idle_duration
            if idle_time > self._kill_threshold:
                await self._kill_polecat(slot, "stuck")
            elif idle_time > self._nudge_threshold:
                await self._nudge_polecat(slot)

    async def _nudge_polecat(self, slot: str):
        """Nudge an idle polecat."""
        polecat = self.rig.polecats.get(slot)
        if polecat:
            polecat.send_prompt("WITNESS: You appear idle. Please continue work or report status.")

    async def _kill_polecat(self, slot: str, reason: str):
        """Kill a stuck polecat."""
        await self.rig.release_polecat(slot)
        print(f"Witness: Killed polecat {slot} ({reason})")

    async def _notify_refinery(self, polecat):
        """Notify Refinery of completed work."""
        from vermas.mail.send import send_mail
        await send_mail(
            to=f"{self.rig.config.name}/refinery",
            subject=f"MERGE_READY: {polecat.bead_id}",
            body=f"Polecat {polecat.slot} completed work on {polecat.bead_id}. Ready for merge.",
        )


# CLAUDE.md template for Witness
WITNESS_TEMPLATE = """
# Witness Context

## Your Role: WITNESS (Worker Monitor)

You monitor polecats in this rig.

**Responsibilities:**
- Track polecat health
- Nudge idle workers
- Kill stuck processes
- Forward completions to Refinery

## Patrol Loop

1. Check mail for POLECAT_DONE
2. Survey all active polecats
3. Nudge idle (>5 min)
4. Kill stuck (>15 min)
5. Report to Deacon if needed

## Key Commands

- `gt polecat list` - List active polecats
- `gt polecat status <slot>` - Check specific polecat
- `gt mail inbox` - Check messages
"""
```

### Refinery

Processes merge queue and code reviews.

```python
# vermas/agents/refinery.py
import asyncio
from typing import List
from vermas.agents.base import BaseAgent
from vermas.models.bead import Bead, BeadStatus


class RefineryAgent(BaseAgent):
    """
    Refinery - Merge queue processor.

    Responsibilities:
    - Process merge requests
    - Run CI/tests
    - Coordinate code review
    - Merge approved changes
    - Handle conflicts

    With VerMAS: Triggers Inspector verification before merge.
    """

    def __init__(self, rig, **kwargs):
        self.rig = rig
        config = self._build_config()
        super().__init__(config, rig.beads)
        self._queue: List[str] = []  # Bead IDs awaiting merge

    def _build_config(self):
        from vermas.models.agent import AgentConfig, AgentIdentity, AgentRole
        return AgentConfig(
            identity=AgentIdentity(
                rig=self.rig.config.name,
                role=AgentRole.REFINERY,
            ),
            working_dir=self.rig.config.path / "refinery",
            claude_profile="refinery",
        )

    async def execute_hooked_work(self, bead: Bead):
        """Process hooked merge request."""
        await self._process_merge(bead)

    async def on_idle(self):
        """Start queue processing loop."""
        asyncio.create_task(self._process_loop())

    async def _process_loop(self):
        """Continuous queue processing."""
        while self.running:
            await self._check_mail()
            await self._process_queue()
            await asyncio.sleep(10)

    async def _check_mail(self):
        """Check for MERGE_READY messages."""
        from vermas.mail.inbox import get_inbox
        messages = await get_inbox(self.identity.actor_string)

        for msg in messages:
            if msg.subject.startswith("MERGE_READY:"):
                bead_id = msg.subject.split(":")[1].strip()
                if bead_id not in self._queue:
                    self._queue.append(bead_id)

    async def _process_queue(self):
        """Process merge queue."""
        if not self._queue:
            return

        bead_id = self._queue[0]
        bead = self.beads.get(bead_id)

        if bead:
            success = await self._process_merge(bead)
            if success:
                self._queue.pop(0)

    async def _process_merge(self, bead: Bead) -> bool:
        """
        Process a single merge request.

        Steps:
        1. Run tests (CI)
        2. Run VerMAS verification (if enabled)
        3. Check for conflicts
        4. Merge or request rework
        """
        # Step 1: Run tests
        tests_pass = await self._run_tests(bead)
        if not tests_pass:
            await self._request_rework(bead, "Tests failed")
            return False

        # Step 2: VerMAS verification
        verification_pass = await self._run_verification(bead)
        if not verification_pass:
            await self._request_rework(bead, "Verification failed")
            return False

        # Step 3: Check conflicts
        has_conflicts = await self._check_conflicts(bead)
        if has_conflicts:
            await self._request_rework(bead, "Merge conflicts")
            return False

        # Step 4: Merge
        await self._merge(bead)
        return True

    async def _run_tests(self, bead: Bead) -> bool:
        """Run tests for the bead's changes."""
        # Implementation would run actual tests
        return True

    async def _run_verification(self, bead: Bead) -> bool:
        """
        Run VerMAS Inspector verification.

        This is the key integration point with VerMAS.
        """
        from vermas.agents.inspector import run_verification
        result = await run_verification(bead)
        return result.verdict == "PASS"

    async def _check_conflicts(self, bead: Bead) -> bool:
        """Check for merge conflicts."""
        # Implementation would check git
        return False

    async def _merge(self, bead: Bead):
        """Perform the merge."""
        # Update bead status
        bead.status = BeadStatus.CLOSED
        self.beads.update(bead)

        # Send MERGED notification
        from vermas.mail.send import send_mail
        await send_mail(
            to=bead.created_by,
            subject=f"MERGED: {bead.id}",
            body=f"Your work on {bead.title} has been merged.",
        )

    async def _request_rework(self, bead: Bead, reason: str):
        """Request rework from the original author."""
        from vermas.mail.send import send_mail
        await send_mail(
            to=bead.created_by,
            subject=f"REWORK_REQUEST: {bead.id}",
            body=f"Merge blocked: {reason}. Please address and resubmit.",
        )


# CLAUDE.md template for Refinery
REFINERY_TEMPLATE = """
# Refinery Context

## Your Role: REFINERY (Merge Queue Processor)

You process merge requests for this rig.

**Responsibilities:**
- Run tests/CI
- Coordinate verification (VerMAS)
- Handle conflicts
- Merge approved changes

## Merge Flow

1. Receive MERGE_READY from Witness
2. Run tests
3. Run VerMAS verification
4. Check conflicts
5. Merge or REWORK_REQUEST

## Key Commands

- `gt merge queue` - View merge queue
- `gt merge process <bead>` - Process specific merge
- `gt mail inbox` - Check messages
"""
```

### Polecat

Ephemeral worker agent.

```python
# vermas/agents/polecat.py
import asyncio
from pathlib import Path
from datetime import datetime
from vermas.agents.base import BaseAgent
from vermas.models.bead import Bead, BeadStatus


class PolecatAgent(BaseAgent):
    """
    Polecat - Ephemeral worker.

    Lifecycle:
    1. SPAWN: Allocate slot, create worktree, launch session
    2. WORK: Execute assigned bead
    3. DISAPPEAR: Clean up, release slot

    Three-layer architecture:
    - Session: Ephemeral tmux session
    - Sandbox: Persistent git worktree
    - Slot: Name allocation from pool
    """

    def __init__(self, rig, slot: str, bead_id: str):
        self.rig = rig
        self.slot = slot
        self.bead_id = bead_id
        self._spawn_time = None
        self._last_activity = None
        config = self._build_config()
        super().__init__(config, rig.beads)

    def _build_config(self):
        from vermas.models.agent import AgentConfig, AgentIdentity, AgentRole
        return AgentConfig(
            identity=AgentIdentity(
                rig=self.rig.config.name,
                role=AgentRole.POLECAT,
                name=self.slot,
            ),
            working_dir=self.rig.config.path / "polecats" / self.slot,
            claude_profile="polecat",
        )

    async def spawn(self):
        """
        Spawn the polecat.

        1. Create git worktree (sandbox)
        2. Hook the assigned bead
        3. Start tmux session (launch)
        """
        # Create worktree
        await self._create_worktree()

        # Hook the bead
        bead = self.beads.get(self.bead_id)
        if bead:
            bead.status = BeadStatus.HOOKED
            self.beads.update(bead)
            self.hook.attach(bead)

        # Start session
        self._spawn_time = datetime.utcnow()
        await self.start()

    async def _create_worktree(self):
        """Create git worktree for this polecat."""
        worktree_path = self.config.working_dir
        worktree_path.mkdir(parents=True, exist_ok=True)

        # Git worktree add command would go here
        # subprocess: git worktree add {path} -b polecat-{slot}

    async def execute_hooked_work(self, bead: Bead):
        """
        Execute the hooked bead.

        GUPP: Run immediately, no confirmation.
        """
        self._last_activity = datetime.utcnow()

        prompt = f"""
POLECAT STARTUP - HOOKED WORK DETECTED

Bead: {bead.id}
Title: {bead.title}
Type: {bead.issue_type}
Priority: {bead.priority}

Description:
{bead.description}

---

GUPP (Propulsion Principle): Work is hooked. Execute immediately.

1. Understand the task
2. Implement the solution
3. Run tests
4. When complete, run: gt polecat done

DO NOT wait for confirmation. BEGIN NOW.
"""
        self.send_prompt(prompt)

    async def on_idle(self):
        """
        Polecats should not be idle - they have hooked work.
        If idle, something went wrong.
        """
        self._last_activity = datetime.utcnow()
        self.send_prompt("ERROR: Polecat started without hooked work. Check hook status.")

    async def cleanup(self):
        """Clean up polecat resources."""
        # Stop session
        await self.stop()

        # Remove worktree
        # subprocess: git worktree remove {path}

    @property
    def is_idle(self) -> bool:
        """Check if polecat appears idle."""
        if not self._last_activity:
            return False
        idle_seconds = (datetime.utcnow() - self._last_activity).total_seconds()
        return idle_seconds > 300  # 5 min

    @property
    def idle_duration(self) -> float:
        """Get idle duration in seconds."""
        if not self._last_activity:
            return 0
        return (datetime.utcnow() - self._last_activity).total_seconds()

    @property
    def is_stuck(self) -> bool:
        """Check if polecat appears stuck."""
        return self.idle_duration > 900  # 15 min

    @property
    def stuck_duration(self) -> float:
        """Get stuck duration in seconds."""
        return self.idle_duration


# CLAUDE.md template for Polecat
POLECAT_TEMPLATE = """
# Polecat Context

## Your Role: POLECAT (Ephemeral Worker)

You are an ephemeral worker. Spawn → Work → Disappear.

**Your only job:** Complete the hooked bead.

## GUPP (Propulsion Principle)

Your hook has work. RUN IT. No confirmation. No waiting.

Check hook: `gt hook`

## Completion

When done:
1. Ensure tests pass
2. Commit and push changes
3. Run: `gt polecat done`

This notifies Witness, who forwards to Refinery for merge.

## Key Commands

- `gt hook` - View your assigned work
- `bd show <id>` - Bead details
- `gt polecat done` - Signal completion
- `gt mail send witness -s "HELP" -m "..."` - Ask for help
"""
```

---

## VerMAS Inspector Agents

### Designer

Elaborates requirements into specifications.

```python
# vermas/agents/inspector/designer.py
from vermas.agents.base import BaseAgent
from vermas.models.verification import Specification


class DesignerAgent(BaseAgent):
    """
    Designer - Elaborates requirements into specifications.

    Input: Raw requirements/user story
    Output: Detailed specification with acceptance criteria

    Uses Claude Code CLI for LLM reasoning.
    """

    async def elaborate(self, requirements: str) -> Specification:
        """
        Elaborate requirements into a specification.

        Prompts Claude to:
        1. Parse the requirements
        2. Identify edge cases
        3. Define acceptance criteria
        4. Output structured specification
        """
        prompt = f"""
You are the DESIGNER in a verification system.

INPUT REQUIREMENTS:
{requirements}

OUTPUT: A detailed specification with:
1. Clear functional requirements
2. Edge cases to handle
3. Acceptance criteria (AC-1, AC-2, etc.)
4. Expected behaviors

Format as JSON:
{{
  "title": "...",
  "requirements": ["R1: ...", "R2: ..."],
  "acceptance_criteria": [
    {{"id": "AC-1", "description": "...", "test_hint": "..."}},
    ...
  ]
}}
"""
        # Run through Claude CLI
        from vermas.claude.cli import run_prompt
        response = await run_prompt(prompt)

        # Parse response into Specification
        return Specification.model_validate_json(response)
```

### Strategist

Creates verification strategy from specifications.

```python
# vermas/agents/inspector/strategist.py
from typing import List
from vermas.agents.base import BaseAgent
from vermas.models.verification import Specification, TestSpec


class StrategistAgent(BaseAgent):
    """
    Strategist - Creates verification strategy.

    Input: Specification from Designer
    Output: Test specifications (shell commands)

    Key principle: Tests must be OBJECTIVE.
    Shell commands that exit 0 for pass, non-zero for fail.
    """

    async def create_tests(self, spec: Specification) -> List[TestSpec]:
        """
        Create test specifications for each acceptance criterion.

        Each test is a shell command that:
        - Exits 0 if criterion is met
        - Exits non-zero if criterion fails
        - Produces clear output for evidence
        """
        prompt = f"""
You are the STRATEGIST in a verification system.

SPECIFICATION:
{spec.model_dump_json(indent=2)}

Create OBJECTIVE tests for each acceptance criterion.

Rules:
1. Each test MUST be a shell command
2. Exit 0 = PASS, non-zero = FAIL
3. Tests must be deterministic
4. Capture evidence in output

Format as JSON array:
[
  {{
    "criterion_id": "AC-1",
    "command": "...",
    "description": "What this tests",
    "timeout": 30
  }},
  ...
]
"""
        from vermas.claude.cli import run_prompt
        response = await run_prompt(prompt)

        import json
        tests_data = json.loads(response)
        return [TestSpec.model_validate(t) for t in tests_data]
```

### Verifier

Runs objective tests (NO LLM - pure shell).

```python
# vermas/agents/inspector/verifier.py
import asyncio
import subprocess
from typing import List
from vermas.models.verification import TestSpec, TestResult, Evidence


class Verifier:
    """
    Verifier - Runs objective tests.

    IMPORTANT: This is NOT an LLM agent.
    It runs shell commands and captures output.

    This ensures objective, reproducible verification.
    """

    async def run_tests(self, tests: List[TestSpec]) -> List[TestResult]:
        """Run all tests and collect results."""
        results = []
        for test in tests:
            result = await self._run_single_test(test)
            results.append(result)
        return results

    async def _run_single_test(self, test: TestSpec) -> TestResult:
        """Run a single test."""
        try:
            process = await asyncio.create_subprocess_shell(
                test.command,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.STDOUT,
            )

            try:
                stdout, _ = await asyncio.wait_for(
                    process.communicate(),
                    timeout=test.timeout,
                )
                output = stdout.decode()
                passed = process.returncode == 0
            except asyncio.TimeoutError:
                process.kill()
                output = f"TIMEOUT after {test.timeout}s"
                passed = False

            return TestResult(
                criterion_id=test.criterion_id,
                passed=passed,
                evidence=Evidence(
                    command=test.command,
                    output=output,
                    exit_code=process.returncode or -1,
                ),
            )

        except Exception as e:
            return TestResult(
                criterion_id=test.criterion_id,
                passed=False,
                evidence=Evidence(
                    command=test.command,
                    output=str(e),
                    exit_code=-1,
                ),
            )
```

### Auditor

LLM verification for cases where shell tests are insufficient.

```python
# vermas/agents/inspector/auditor.py
from vermas.agents.base import BaseAgent
from vermas.models.verification import Evidence, AuditResult


class AuditorAgent(BaseAgent):
    """
    Auditor - LLM verification for subjective criteria.

    Used when:
    1. Shell tests cannot verify (e.g., "code is readable")
    2. Shell test failed but evidence is ambiguous
    3. Human-like judgment is required

    Still must provide EVIDENCE and REASONING.
    """

    async def audit(self, evidence: Evidence, criterion: str) -> AuditResult:
        """
        Audit evidence against a criterion.

        The Auditor must:
        1. Examine the evidence
        2. Apply the criterion
        3. Provide clear reasoning
        4. Return PASS/FAIL with justification
        """
        prompt = f"""
You are the AUDITOR in a verification system.

CRITERION:
{criterion}

EVIDENCE:
Command: {evidence.command}
Output:
{evidence.output}
Exit code: {evidence.exit_code}

Evaluate whether this evidence demonstrates the criterion is met.

You must provide:
1. Your assessment (PASS or FAIL)
2. Clear reasoning based on the evidence
3. Specific quotes from evidence supporting your assessment

Format:
ASSESSMENT: [PASS/FAIL]
REASONING: [Your detailed reasoning]
EVIDENCE_QUOTES: [Specific quotes from the output]
"""
        from vermas.claude.cli import run_prompt
        response = await run_prompt(prompt)

        # Parse response
        passed = "ASSESSMENT: PASS" in response
        return AuditResult(
            criterion=criterion,
            passed=passed,
            reasoning=response,
        )
```

### Advocate, Critic, Judge

Adversarial review panel.

```python
# vermas/agents/inspector/adversarial.py
from vermas.agents.base import BaseAgent
from vermas.models.verification import Evidence, Verdict


class AdvocateAgent(BaseAgent):
    """Advocate - Argues for PASS."""

    async def argue(self, evidence: Evidence, criterion: str) -> str:
        prompt = f"""
You are the ADVOCATE. Your job is to argue that the criterion IS met.

CRITERION: {criterion}
EVIDENCE: {evidence.output}

Make the strongest possible case for PASS.
Be specific. Quote evidence. Address potential objections.
"""
        from vermas.claude.cli import run_prompt
        return await run_prompt(prompt)


class CriticAgent(BaseAgent):
    """Critic - Argues for FAIL."""

    async def argue(self, evidence: Evidence, criterion: str) -> str:
        prompt = f"""
You are the CRITIC. Your job is to argue that the criterion is NOT met.

CRITERION: {criterion}
EVIDENCE: {evidence.output}

Make the strongest possible case for FAIL.
Be specific. Identify gaps. Challenge assumptions.
"""
        from vermas.claude.cli import run_prompt
        return await run_prompt(prompt)


class JudgeAgent(BaseAgent):
    """Judge - Decides verdict after hearing both sides."""

    async def decide(
        self,
        criterion: str,
        evidence: Evidence,
        advocate_arg: str,
        critic_arg: str,
    ) -> Verdict:
        prompt = f"""
You are the JUDGE. You must decide: PASS or FAIL.

CRITERION: {criterion}

EVIDENCE:
{evidence.output}

ADVOCATE'S ARGUMENT:
{advocate_arg}

CRITIC'S ARGUMENT:
{critic_arg}

Consider both arguments carefully. Your verdict must be based on:
1. The actual evidence
2. The specific criterion
3. The strength of each argument

VERDICT: [PASS/FAIL]
REASONING: [Your reasoning]
"""
        from vermas.claude.cli import run_prompt
        response = await run_prompt(prompt)

        passed = "VERDICT: PASS" in response
        return Verdict(
            passed=passed,
            reasoning=response,
        )
```

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - Overall system architecture
- [WORKFLOWS.md](./WORKFLOWS.md) - Molecule workflow system
- [MESSAGING.md](./MESSAGING.md) - Mail protocol
- [LIFECYCLE.md](./LIFECYCLE.md) - Polecat lifecycle
