# VerMAS Python Workflows

> Molecule system and LangGraph workflows for Gas Town + VerMAS

## Molecule State Machine (MEOW)

Gas Town uses the MEOW (Molecule states) system for workflow management:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        MOLECULE STATE MACHINE                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────┐                                                           │
│   │   Ice-9     │  Formula - Source templates (.toml files)                 │
│   │  (Formula)  │  Location: .beads/formulas/*.formula.toml                 │
│   └──────┬──────┘                                                           │
│          │ cook                                                              │
│          ▼                                                                   │
│   ┌─────────────┐                                                           │
│   │   Solid     │  Protomolecule - Frozen, reusable templates              │
│   │  (Proto)    │  Created from formulas, ready to instantiate              │
│   └──────┬──────┘                                                           │
│          │                                                                   │
│          ├────────────────────────────────┐                                 │
│          │ pour                           │ wisp                            │
│          ▼                                ▼                                  │
│   ┌─────────────┐                  ┌─────────────┐                          │
│   │   Liquid    │                  │   Vapor     │                          │
│   │   (Mol)     │                  │   (Wisp)    │                          │
│   │             │                  │             │                          │
│   │ Persistent  │                  │ Ephemeral   │                          │
│   │ workflow    │                  │ patrol      │                          │
│   └──────┬──────┘                  └──────┬──────┘                          │
│          │                                │                                  │
│          ├────────────────────────────────┤                                 │
│          │ squash              │ burn     │ (auto-evaporate)                │
│          ▼                     ▼          ▼                                  │
│   ┌─────────────┐       ┌─────────────┐                                     │
│   │  Archive    │       │  Discarded  │                                     │
│   │  (Record)   │       │  (No trace) │                                     │
│   └─────────────┘       └─────────────┘                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

Operators:
- cook:   Formula → Protomolecule (compile template)
- pour:   Protomolecule → Molecule (instantiate persistent workflow)
- wisp:   Protomolecule → Wisp (instantiate ephemeral workflow)
- squash: Mol/Wisp → Archive (condense to permanent record)
- burn:   Mol/Wisp → Nothing (discard without record)
```

---

## Python Implementation

### Pydantic Models

```python
# vermas/models/molecule.py
from enum import Enum
from typing import List, Optional, Dict, Any
from datetime import datetime
from pydantic import BaseModel, Field


class MoleculeState(str, Enum):
    """Molecule lifecycle states."""
    ICE9 = "ice9"       # Formula (source template)
    SOLID = "solid"     # Protomolecule (compiled)
    LIQUID = "liquid"   # Molecule (active workflow)
    VAPOR = "vapor"     # Wisp (ephemeral)
    ARCHIVE = "archive" # Squashed (complete)


class StepStatus(str, Enum):
    """Workflow step status."""
    PENDING = "pending"
    IN_PROGRESS = "in_progress"
    COMPLETED = "completed"
    BLOCKED = "blocked"
    SKIPPED = "skipped"


class Step(BaseModel):
    """Single workflow step."""
    id: str
    title: str
    description: str
    needs: List[str] = Field(default_factory=list)  # Dependencies
    status: StepStatus = StepStatus.PENDING
    started_at: Optional[datetime] = None
    completed_at: Optional[datetime] = None
    output: Optional[str] = None


class Formula(BaseModel):
    """
    Ice-9 Formula - Source workflow template.

    Stored in .beads/formulas/*.formula.toml
    """
    name: str
    description: str
    version: int = 1
    steps: List[Step]

    @classmethod
    def from_toml(cls, path: str) -> "Formula":
        """Load formula from TOML file."""
        import tomllib
        with open(path, "rb") as f:
            data = tomllib.load(f)

        steps = [
            Step(
                id=s["id"],
                title=s["title"],
                description=s.get("description", ""),
                needs=s.get("needs", []),
            )
            for s in data.get("steps", [])
        ]

        return cls(
            name=data.get("formula", path),
            description=data.get("description", ""),
            version=data.get("version", 1),
            steps=steps,
        )


class Protomolecule(BaseModel):
    """
    Solid Protomolecule - Compiled template ready for instantiation.
    """
    formula_name: str
    steps: List[Step]
    created_at: datetime = Field(default_factory=datetime.utcnow)


class Molecule(BaseModel):
    """
    Liquid Molecule - Active, persistent workflow instance.

    Attached to a bead for work tracking.
    """
    id: str
    proto_name: str
    bead_id: str
    state: MoleculeState = MoleculeState.LIQUID
    steps: List[Step]
    current_step: Optional[str] = None
    created_at: datetime = Field(default_factory=datetime.utcnow)
    updated_at: datetime = Field(default_factory=datetime.utcnow)
    metadata: Dict[str, Any] = Field(default_factory=dict)

    def get_ready_steps(self) -> List[Step]:
        """Get steps that are ready to work (dependencies met)."""
        completed_ids = {
            s.id for s in self.steps
            if s.status == StepStatus.COMPLETED
        }

        ready = []
        for step in self.steps:
            if step.status != StepStatus.PENDING:
                continue
            if all(dep in completed_ids for dep in step.needs):
                ready.append(step)

        return ready

    def complete_step(self, step_id: str, output: str = None):
        """Mark a step as completed."""
        for step in self.steps:
            if step.id == step_id:
                step.status = StepStatus.COMPLETED
                step.completed_at = datetime.utcnow()
                step.output = output
                break
        self.updated_at = datetime.utcnow()


class Wisp(Molecule):
    """
    Vapor Wisp - Ephemeral workflow for temporary patrols.

    Auto-evaporates after completion or timeout.
    """
    state: MoleculeState = MoleculeState.VAPOR
    ttl_seconds: int = 3600  # 1 hour default
    expires_at: Optional[datetime] = None

    def __init__(self, **data):
        super().__init__(**data)
        if not self.expires_at:
            self.expires_at = datetime.utcnow() + timedelta(seconds=self.ttl_seconds)

    @property
    def is_expired(self) -> bool:
        return datetime.utcnow() > self.expires_at if self.expires_at else False
```

### Molecule Lifecycle Operations

```python
# vermas/molecules/lifecycle.py
from pathlib import Path
from typing import Optional
import json
import hashlib
from datetime import datetime

from vermas.models.molecule import (
    Formula, Protomolecule, Molecule, Wisp,
    MoleculeState, StepStatus
)


class MoleculeManager:
    """
    Manages molecule lifecycle operations.

    Operations:
    - cook: Formula → Protomolecule
    - pour: Protomolecule → Molecule
    - wisp: Protomolecule → Wisp
    - squash: Mol → Archive
    - burn: Mol → Discard
    """

    def __init__(self, beads_dir: Path):
        self.beads_dir = beads_dir
        self.formulas_dir = beads_dir / "formulas"
        self.mols_dir = beads_dir / "mols"
        self._ensure_dirs()

    def _ensure_dirs(self):
        self.formulas_dir.mkdir(parents=True, exist_ok=True)
        self.mols_dir.mkdir(parents=True, exist_ok=True)

    def _generate_id(self, prefix: str = "mol") -> str:
        timestamp = datetime.utcnow().isoformat()
        short_hash = hashlib.sha256(timestamp.encode()).hexdigest()[:8]
        return f"{prefix}-{short_hash}"

    # ─────────────────────────────────────────────────────────────
    # COOK: Formula → Protomolecule
    # ─────────────────────────────────────────────────────────────

    def cook(self, formula_name: str) -> Protomolecule:
        """
        Cook a formula into a protomolecule.

        Loads the TOML formula and creates a compiled template.
        """
        formula_path = self.formulas_dir / f"{formula_name}.formula.toml"
        if not formula_path.exists():
            raise FileNotFoundError(f"Formula not found: {formula_name}")

        formula = Formula.from_toml(str(formula_path))

        return Protomolecule(
            formula_name=formula.name,
            steps=formula.steps,
        )

    # ─────────────────────────────────────────────────────────────
    # POUR: Protomolecule → Molecule (persistent)
    # ─────────────────────────────────────────────────────────────

    def pour(self, proto: Protomolecule, bead_id: str) -> Molecule:
        """
        Pour a protomolecule into a molecule attached to a bead.

        Creates a persistent workflow instance.
        """
        mol_id = self._generate_id("mol")

        mol = Molecule(
            id=mol_id,
            proto_name=proto.formula_name,
            bead_id=bead_id,
            steps=proto.steps.copy(),
        )

        # Persist to disk
        self._save_mol(mol)
        return mol

    # ─────────────────────────────────────────────────────────────
    # WISP: Protomolecule → Wisp (ephemeral)
    # ─────────────────────────────────────────────────────────────

    def wisp(
        self,
        proto: Protomolecule,
        bead_id: str = "ephemeral",
        ttl_seconds: int = 3600
    ) -> Wisp:
        """
        Create an ephemeral wisp from a protomolecule.

        Wisps auto-evaporate after TTL expires.
        """
        wisp_id = self._generate_id("wisp")

        wisp = Wisp(
            id=wisp_id,
            proto_name=proto.formula_name,
            bead_id=bead_id,
            steps=proto.steps.copy(),
            ttl_seconds=ttl_seconds,
        )

        # Don't persist wisps to disk by default
        return wisp

    # ─────────────────────────────────────────────────────────────
    # SQUASH: Mol → Archive (permanent record)
    # ─────────────────────────────────────────────────────────────

    def squash(self, mol: Molecule, summary: str = None) -> dict:
        """
        Squash a molecule into a permanent archive record.

        Returns the archive record.
        """
        archive = {
            "id": mol.id,
            "proto_name": mol.proto_name,
            "bead_id": mol.bead_id,
            "steps": [
                {
                    "id": s.id,
                    "title": s.title,
                    "status": s.status.value,
                    "output": s.output,
                }
                for s in mol.steps
            ],
            "summary": summary,
            "squashed_at": datetime.utcnow().isoformat(),
        }

        # Write to archive
        archive_file = self.mols_dir / f"{mol.id}.archive.json"
        archive_file.write_text(json.dumps(archive, indent=2))

        # Remove active mol file
        mol_file = self.mols_dir / f"{mol.id}.json"
        if mol_file.exists():
            mol_file.unlink()

        return archive

    # ─────────────────────────────────────────────────────────────
    # BURN: Mol → Discard (no record)
    # ─────────────────────────────────────────────────────────────

    def burn(self, mol: Molecule):
        """
        Burn a molecule without creating a record.

        Use for failed/abandoned workflows.
        """
        mol_file = self.mols_dir / f"{mol.id}.json"
        if mol_file.exists():
            mol_file.unlink()

    # ─────────────────────────────────────────────────────────────
    # Persistence
    # ─────────────────────────────────────────────────────────────

    def _save_mol(self, mol: Molecule):
        """Save molecule to disk."""
        mol_file = self.mols_dir / f"{mol.id}.json"
        mol_file.write_text(mol.model_dump_json(indent=2))

    def load_mol(self, mol_id: str) -> Optional[Molecule]:
        """Load molecule from disk."""
        mol_file = self.mols_dir / f"{mol_id}.json"
        if not mol_file.exists():
            return None
        return Molecule.model_validate_json(mol_file.read_text())

    def list_active_mols(self) -> list[Molecule]:
        """List all active molecules."""
        mols = []
        for f in self.mols_dir.glob("*.json"):
            if not f.name.endswith(".archive.json"):
                mol = Molecule.model_validate_json(f.read_text())
                mols.append(mol)
        return mols
```

### Workflow Executor

```python
# vermas/molecules/executor.py
import asyncio
from typing import Optional, Callable, Awaitable
from vermas.models.molecule import Molecule, Step, StepStatus


class WorkflowExecutor:
    """
    Executes molecule workflows step by step.

    The executor:
    1. Finds ready steps (dependencies met)
    2. Executes step via callback
    3. Updates step status
    4. Repeats until complete
    """

    def __init__(self, mol: Molecule):
        self.mol = mol
        self._running = False

    async def execute(
        self,
        step_handler: Callable[[Step], Awaitable[str]],
        on_complete: Callable[[Molecule], Awaitable[None]] = None,
    ):
        """
        Execute the workflow.

        Args:
            step_handler: Async function to execute each step
            on_complete: Called when workflow finishes
        """
        self._running = True

        while self._running:
            ready_steps = self.mol.get_ready_steps()

            if not ready_steps:
                # Check if all done
                if all(s.status == StepStatus.COMPLETED for s in self.mol.steps):
                    break
                # Some steps blocked - wait
                await asyncio.sleep(1)
                continue

            # Execute first ready step
            step = ready_steps[0]
            self.mol.current_step = step.id
            step.status = StepStatus.IN_PROGRESS
            step.started_at = datetime.utcnow()

            try:
                output = await step_handler(step)
                self.mol.complete_step(step.id, output)
            except Exception as e:
                step.status = StepStatus.BLOCKED
                step.output = str(e)

        self._running = False
        if on_complete:
            await on_complete(self.mol)

    def stop(self):
        """Stop execution."""
        self._running = False
```

---

## Common Formulas

### Witness Patrol

```toml
# .beads/formulas/mol-witness-patrol.formula.toml
description = "Witness monitors polecats, nudges idle, escalates stuck"
formula = 'mol-witness-patrol'
version = 2

[[steps]]
id = 'inbox-check'
title = 'Process witness mail'
description = """
Check for incoming mail:
- POLECAT_DONE messages
- Escalations from Deacon
- Help requests

Run: `gt mail inbox`
"""

[[steps]]
id = 'survey-workers'
title = 'Inspect all active polecats'
needs = ['inbox-check']
description = """
Survey all polecats in this rig:
- Check session status
- Measure idle time
- Identify stuck workers

Run: `gt polecat list`
"""

[[steps]]
id = 'nudge-idle'
title = 'Nudge idle polecats'
needs = ['survey-workers']
description = """
Send nudge to polecats idle > 5 min.

For each idle polecat:
`gt mail send <polecat> -s "NUDGE" -m "Please continue work or report status"`
"""

[[steps]]
id = 'escalate-stuck'
title = 'Escalate stuck polecats'
needs = ['nudge-idle']
description = """
For polecats stuck > 15 min:
1. Kill the session
2. Report to Deacon
3. Return slot to pool

`gt mail send deacon -s "STUCK: <polecat>" -m "Details..."`
"""
```

### Polecat Work

```toml
# .beads/formulas/mol-polecat-work.formula.toml
description = "Polecat executes assigned bead through completion"
formula = 'mol-polecat-work'
version = 1

[[steps]]
id = 'understand'
title = 'Understand the task'
description = """
Read and understand the assigned bead:

1. Run: `gt hook` to see assigned work
2. Run: `bd show <bead-id>` for details
3. Identify requirements
4. Plan implementation approach
"""

[[steps]]
id = 'implement'
title = 'Implement the solution'
needs = ['understand']
description = """
Write the code:

1. Create/modify necessary files
2. Follow project conventions
3. Add tests if required
4. Commit changes locally
"""

[[steps]]
id = 'test'
title = 'Run tests'
needs = ['implement']
description = """
Verify the implementation:

1. Run project tests
2. Check for regressions
3. Verify acceptance criteria
"""

[[steps]]
id = 'complete'
title = 'Signal completion'
needs = ['test']
description = """
Wrap up and notify:

1. Push changes to branch
2. Run: `gt polecat done`
3. This notifies Witness for merge processing
"""
```

### VerMAS Inspector Workflow

```toml
# .beads/formulas/mol-inspector-verify.formula.toml
description = "VerMAS verification workflow for merge requests"
formula = 'mol-inspector-verify'
version = 1

[[steps]]
id = 'elaborate'
title = 'Designer elaborates requirements'
description = """
Designer agent:
1. Parse the original requirements
2. Identify edge cases
3. Create detailed specification
4. Output acceptance criteria
"""

[[steps]]
id = 'strategize'
title = 'Strategist creates test plan'
needs = ['elaborate']
description = """
Strategist agent:
1. Review specification
2. Design objective tests
3. Create shell commands for each criterion
4. Output test specifications
"""

[[steps]]
id = 'verify'
title = 'Verifier runs tests'
needs = ['strategize']
description = """
Verifier (NO LLM):
1. Execute each test command
2. Capture output as evidence
3. Record pass/fail for each criterion
4. Pass evidence to Auditor
"""

[[steps]]
id = 'audit'
title = 'Auditor reviews evidence'
needs = ['verify']
description = """
Auditor agent (if shell tests inconclusive):
1. Review evidence from Verifier
2. Apply LLM judgment where needed
3. Document reasoning
4. Forward to adversarial review
"""

[[steps]]
id = 'advocate'
title = 'Advocate argues for PASS'
needs = ['audit']
description = """
Advocate agent:
1. Review evidence
2. Build case for PASS
3. Address potential objections
4. Submit argument to Judge
"""

[[steps]]
id = 'criticize'
title = 'Critic argues for FAIL'
needs = ['audit']
description = """
Critic agent:
1. Review evidence
2. Build case for FAIL
3. Identify gaps and weaknesses
4. Submit argument to Judge
"""

[[steps]]
id = 'judge'
title = 'Judge renders verdict'
needs = ['advocate', 'criticize']
description = """
Judge agent:
1. Review evidence
2. Consider both arguments
3. Render PASS or FAIL verdict
4. Document reasoning
"""
```

---

## LangGraph Integration (Optional)

For complex workflows, LangGraph can manage state transitions:

```python
# vermas/molecules/langgraph_executor.py
"""
Optional LangGraph integration for complex workflows.

Note: This uses LangGraph for state management only.
LLM calls still go through Claude Code CLI (no API costs).
"""

from typing import TypedDict, Annotated
from langgraph.graph import StateGraph, END


class WorkflowState(TypedDict):
    """State for LangGraph workflow."""
    mol_id: str
    current_step: str
    step_outputs: dict
    completed_steps: list
    error: str | None


def create_workflow_graph(mol: Molecule) -> StateGraph:
    """
    Create a LangGraph StateGraph from a molecule.

    Each step becomes a node in the graph.
    Dependencies become edges.
    """
    graph = StateGraph(WorkflowState)

    # Add nodes for each step
    for step in mol.steps:
        async def execute_step(state: WorkflowState, step=step) -> WorkflowState:
            # Execute step via Claude Code CLI
            from vermas.claude.cli import run_prompt
            output = await run_prompt(step.description)

            state["step_outputs"][step.id] = output
            state["completed_steps"].append(step.id)
            return state

        graph.add_node(step.id, execute_step)

    # Add edges based on dependencies
    for step in mol.steps:
        if not step.needs:
            # No dependencies - connect from START
            graph.set_entry_point(step.id)
        else:
            # Connect from each dependency
            for dep in step.needs:
                graph.add_edge(dep, step.id)

    # Find terminal nodes and connect to END
    terminal_steps = []
    all_deps = set()
    for step in mol.steps:
        all_deps.update(step.needs)

    for step in mol.steps:
        if step.id not in all_deps:
            terminal_steps.append(step.id)

    for term in terminal_steps:
        graph.add_edge(term, END)

    return graph.compile()


async def run_langgraph_workflow(mol: Molecule):
    """Run a molecule workflow using LangGraph."""
    graph = create_workflow_graph(mol)

    initial_state: WorkflowState = {
        "mol_id": mol.id,
        "current_step": "",
        "step_outputs": {},
        "completed_steps": [],
        "error": None,
    }

    final_state = await graph.ainvoke(initial_state)
    return final_state
```

---

## Workflow Commands

CLI commands for molecule management:

```python
# vermas/cli.py (molecule commands)
import typer
from pathlib import Path

mol_app = typer.Typer()


@mol_app.command("cook")
def cook(formula: str):
    """Cook a formula into a protomolecule."""
    from vermas.molecules.lifecycle import MoleculeManager

    mgr = MoleculeManager(Path(".beads"))
    proto = mgr.cook(formula)
    typer.echo(f"Cooked: {proto.formula_name} ({len(proto.steps)} steps)")


@mol_app.command("pour")
def pour(formula: str, bead_id: str):
    """Pour a protomolecule into a molecule attached to a bead."""
    from vermas.molecules.lifecycle import MoleculeManager

    mgr = MoleculeManager(Path(".beads"))
    proto = mgr.cook(formula)
    mol = mgr.pour(proto, bead_id)
    typer.echo(f"Poured: {mol.id} attached to {bead_id}")


@mol_app.command("wisp")
def wisp(formula: str, ttl: int = 3600):
    """Create an ephemeral wisp from a formula."""
    from vermas.molecules.lifecycle import MoleculeManager

    mgr = MoleculeManager(Path(".beads"))
    proto = mgr.cook(formula)
    w = mgr.wisp(proto, ttl_seconds=ttl)
    typer.echo(f"Wisped: {w.id} (expires in {ttl}s)")


@mol_app.command("squash")
def squash(mol_id: str, summary: str = None):
    """Archive a molecule with optional summary."""
    from vermas.molecules.lifecycle import MoleculeManager

    mgr = MoleculeManager(Path(".beads"))
    mol = mgr.load_mol(mol_id)
    if not mol:
        typer.echo(f"Molecule not found: {mol_id}", err=True)
        raise typer.Exit(1)

    archive = mgr.squash(mol, summary)
    typer.echo(f"Squashed: {mol_id}")


@mol_app.command("burn")
def burn(mol_id: str, force: bool = False):
    """Discard a molecule without record."""
    from vermas.molecules.lifecycle import MoleculeManager

    if not force:
        typer.confirm(f"Burn {mol_id}? This cannot be undone.", abort=True)

    mgr = MoleculeManager(Path(".beads"))
    mol = mgr.load_mol(mol_id)
    if mol:
        mgr.burn(mol)
        typer.echo(f"Burned: {mol_id}")


@mol_app.command("list")
def list_mols():
    """List active molecules."""
    from vermas.molecules.lifecycle import MoleculeManager

    mgr = MoleculeManager(Path(".beads"))
    mols = mgr.list_active_mols()

    for mol in mols:
        completed = sum(1 for s in mol.steps if s.status.value == "completed")
        total = len(mol.steps)
        typer.echo(f"{mol.id}: {mol.proto_name} ({completed}/{total} steps) - {mol.bead_id}")
```

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [AGENTS.md](./AGENTS.md) - Agent implementations
- [MESSAGING.md](./MESSAGING.md) - Mail protocol
- [LIFECYCLE.md](./LIFECYCLE.md) - Polecat lifecycle
