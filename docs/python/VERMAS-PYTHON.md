# VerMAS Python Implementation

> Building the Verifiable Multi-Agent System in Python with LangGraph & Pydantic

## Overview

This document describes how to implement VerMAS in Python using:
- **LangGraph** - Stateful multi-agent orchestration
- **Pydantic** - Data models and validation
- **LangChain** - LLM abstractions
- **SQLite/SQLModel** - Persistent state (beads)
- **FastAPI** - Optional REST API
- **Rich** - CLI output formatting

## Why Python?

| Aspect | Go (Current) | Python (Proposed) |
|--------|--------------|-------------------|
| LLM ecosystem | Custom runtime abstraction | Native LangChain/LangGraph |
| Agent orchestration | Shell + tmux + hooks | LangGraph graphs |
| State machines | Custom molecule system | LangGraph StateGraph |
| Typing | Compile-time | Pydantic runtime validation |
| Iteration speed | Compile cycle | Instant reload |
| AI tooling | Limited | Extensive (LangChain, etc.) |

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         VERMAS PYTHON ARCHITECTURE                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         LANGGRAPH LAYER                              │   │
│  │                                                                      │   │
│  │   ┌──────────────┐   ┌──────────────┐   ┌──────────────┐           │   │
│  │   │ MayorGraph   │   │InspectorGraph│   │PolecatGraph  │           │   │
│  │   │              │   │              │   │              │           │   │
│  │   │  Designer    │   │  Strategist  │   │  Implement   │           │   │
│  │   │  Coordinate  │   │  Verify      │   │  Test        │           │   │
│  │   │  Dispatch    │   │  Judge       │   │  Submit      │           │   │
│  │   └──────────────┘   └──────────────┘   └──────────────┘           │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                      │                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                         PYDANTIC LAYER                               │   │
│  │                                                                      │   │
│  │   WorkItem    TestSpec    Verdict    Message    AgentState          │   │
│  │   Criterion   Brief       Evidence   Config     ...                  │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                      │                                      │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                       PERSISTENCE LAYER                              │   │
│  │                                                                      │   │
│  │   SQLite + SQLModel (beads)    Redis (optional, for pub/sub)        │   │
│  │                                                                      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Project Structure

```
vermas-py/
├── pyproject.toml              # Project config (poetry/uv)
├── vermas/
│   ├── __init__.py
│   ├── cli.py                  # Typer CLI (gt equivalent)
│   │
│   ├── models/                 # Pydantic models
│   │   ├── __init__.py
│   │   ├── bead.py            # WorkItem, TestSpec, Message
│   │   ├── agent.py           # AgentState, AgentConfig
│   │   ├── verification.py    # Verdict, Evidence, Brief
│   │   └── config.py          # Settings, RuntimeConfig
│   │
│   ├── graphs/                 # LangGraph definitions
│   │   ├── __init__.py
│   │   ├── mayor.py           # Mayor workflow graph
│   │   ├── inspector.py       # Inspector workflow graph
│   │   ├── polecat.py         # Polecat workflow graph
│   │   ├── verification.py    # Verification subgraph
│   │   └── spec_planning.py   # AI-assisted spec planning
│   │
│   ├── agents/                 # Agent implementations
│   │   ├── __init__.py
│   │   ├── base.py            # BaseAgent class
│   │   ├── mayor.py           # Mayor agent
│   │   ├── inspector.py       # Inspector coordinator
│   │   ├── designer.py        # Requirements elaboration
│   │   ├── strategist.py      # Criteria proposal
│   │   ├── verifier.py        # Objective test runner
│   │   ├── auditor.py         # Compliance checker
│   │   ├── advocate.py        # Defense builder
│   │   ├── critic.py          # Prosecution builder
│   │   └── judge.py           # Verdict deliverer
│   │
│   ├── db/                     # Persistence (beads)
│   │   ├── __init__.py
│   │   ├── models.py          # SQLModel definitions
│   │   ├── repository.py      # CRUD operations
│   │   └── migrations/        # Alembic migrations
│   │
│   ├── messaging/              # Agent communication
│   │   ├── __init__.py
│   │   ├── mailbox.py         # Message queue
│   │   └── router.py          # Address routing
│   │
│   ├── tools/                  # LangChain tools
│   │   ├── __init__.py
│   │   ├── git.py             # Git operations
│   │   ├── shell.py           # Shell execution
│   │   ├── beads.py           # Beads operations
│   │   └── mail.py            # Messaging tools
│   │
│   └── utils/
│       ├── __init__.py
│       ├── llm.py             # LLM factory
│       └── prompts.py         # Prompt templates
│
├── tests/
│   ├── test_models.py
│   ├── test_graphs.py
│   └── test_agents.py
│
└── examples/
    ├── fizzbuzz.py            # FizzBuzz dry run
    └── simple_task.py         # Simple task example
```

---

## Pydantic Models

### Core Data Models

```python
# vermas/models/bead.py
from datetime import datetime
from enum import Enum
from typing import Optional
from pydantic import BaseModel, Field


class BeadType(str, Enum):
    TASK = "task"
    FEATURE = "feature"
    BUG = "bug"
    TEST_SPEC = "test-spec"
    VERIFICATION_RESULT = "verification-result"
    MESSAGE = "message"


class BeadStatus(str, Enum):
    OPEN = "open"
    IN_PROGRESS = "in_progress"
    BLOCKED = "blocked"
    CLOSED = "closed"
    REJECTED = "rejected"


class Priority(int, Enum):
    CRITICAL = 0
    HIGH = 1
    MEDIUM = 2
    LOW = 3
    BACKLOG = 4


class WorkItem(BaseModel):
    """A work item (task, feature, bug)."""

    id: str = Field(..., pattern=r"^[a-z]+-[a-z0-9]+$")
    type: BeadType
    title: str = Field(..., min_length=1, max_length=200)
    description: str = ""
    status: BeadStatus = BeadStatus.OPEN
    priority: Priority = Priority.MEDIUM
    assignee: Optional[str] = None
    parent: Optional[str] = None
    depends_on: list[str] = Field(default_factory=list)
    labels: list[str] = Field(default_factory=list)
    created_at: datetime = Field(default_factory=datetime.utcnow)
    updated_at: datetime = Field(default_factory=datetime.utcnow)

    class Config:
        use_enum_values = True


class AcceptanceCriterion(BaseModel):
    """A single testable acceptance criterion."""

    id: str = Field(..., pattern=r"^AC-\d+$")
    description: str = Field(..., min_length=1)
    verify_command: str = Field(..., min_length=1)
    expected_exit_code: int = 0
    timeout_seconds: int = 60


class TestSpec(BaseModel):
    """Test specification for a work item."""

    id: str = Field(..., pattern=r"^spec-[a-z0-9]+$")
    parent_work_item: str
    status: BeadStatus = BeadStatus.OPEN
    criteria: list[AcceptanceCriterion] = Field(default_factory=list)
    success_threshold: float = Field(default=1.0, ge=0.0, le=1.0)  # 1.0 = all must pass
    created_at: datetime = Field(default_factory=datetime.utcnow)
    approved_at: Optional[datetime] = None
    approved_by: Optional[str] = None


class Requirements(BaseModel):
    """Structured requirements from Designer."""

    functional: list[str] = Field(default_factory=list)
    non_functional: list[str] = Field(default_factory=list)
    constraints: list[str] = Field(default_factory=list)
    examples: list[dict] = Field(default_factory=list)
```

### Verification Models

```python
# vermas/models/verification.py
from datetime import datetime
from enum import Enum
from typing import Optional
from pydantic import BaseModel, Field


class Verdict(str, Enum):
    PASS = "PASS"
    FAIL = "FAIL"
    NEEDS_HUMAN = "NEEDS_HUMAN"


class CriterionResult(BaseModel):
    """Result of running a single criterion."""

    criterion_id: str
    status: str  # "pass", "fail", "error", "timeout"
    output: str
    exit_code: int
    duration_ms: int


class Evidence(BaseModel):
    """All evidence for adversarial review."""

    work_item: "WorkItem"
    test_spec: "TestSpec"
    criterion_results: list[CriterionResult]
    diff: str  # Git diff
    commit_messages: list[str]
    files_changed: list[str]


class Brief(BaseModel):
    """Argument brief from Advocate or Critic."""

    role: str  # "advocate" or "critic"
    position: str  # "FOR" or "AGAINST"
    arguments: list[str]
    evidence_cited: list[str]
    confidence: float = Field(..., ge=0.0, le=1.0)
    summary: str


class VerificationResult(BaseModel):
    """Complete verification result."""

    work_item_id: str
    spec_id: str
    verdict: Verdict
    confidence: float = Field(..., ge=0.0, le=1.0)

    criterion_results: list[CriterionResult]
    advocate_brief: Optional[Brief] = None
    critic_brief: Optional[Brief] = None
    judge_reasoning: str = ""

    issues: list[str] = Field(default_factory=list)
    suggestions: list[str] = Field(default_factory=list)
    required_fixes: list[str] = Field(default_factory=list)

    reviewed_at: datetime = Field(default_factory=datetime.utcnow)
    duration_ms: int = 0

    def is_pass(self) -> bool:
        return self.verdict == Verdict.PASS
```

### Agent State Models

```python
# vermas/models/agent.py
from datetime import datetime
from enum import Enum
from typing import Any, Optional
from pydantic import BaseModel, Field


class AgentRole(str, Enum):
    MAYOR = "mayor"
    INSPECTOR = "inspector"
    POLECAT = "polecat"
    DESIGNER = "designer"
    STRATEGIST = "strategist"
    VERIFIER = "verifier"
    AUDITOR = "auditor"
    ADVOCATE = "advocate"
    CRITIC = "critic"
    JUDGE = "judge"


class AgentStatus(str, Enum):
    IDLE = "idle"
    WORKING = "working"
    WAITING = "waiting"
    BLOCKED = "blocked"
    DONE = "done"


class AgentState(BaseModel):
    """Shared state for LangGraph agents."""

    # Identity
    agent_id: str
    role: AgentRole

    # Current work
    current_work_item: Optional[str] = None
    current_spec: Optional[str] = None
    status: AgentStatus = AgentStatus.IDLE

    # Communication
    messages: list["Message"] = Field(default_factory=list)
    pending_questions: list[str] = Field(default_factory=list)

    # Context
    working_directory: str = "."
    conversation_history: list[dict] = Field(default_factory=list)

    # Metadata
    started_at: datetime = Field(default_factory=datetime.utcnow)
    last_activity: datetime = Field(default_factory=datetime.utcnow)

    class Config:
        arbitrary_types_allowed = True


class Message(BaseModel):
    """Inter-agent message."""

    id: str
    from_agent: str
    to_agent: str
    subject: str
    body: str
    message_type: str = "notification"  # notification, task, reply
    priority: str = "normal"
    thread_id: Optional[str] = None
    reply_to: Optional[str] = None
    created_at: datetime = Field(default_factory=datetime.utcnow)
    read: bool = False
```

### Configuration Models

```python
# vermas/models/config.py
from typing import Optional
from pydantic import BaseModel, Field
from pydantic_settings import BaseSettings


class LLMConfig(BaseModel):
    """Configuration for an LLM."""

    provider: str = "anthropic"  # anthropic, openai, ollama
    model: str = "claude-sonnet-4-20250514"
    temperature: float = 0.0
    max_tokens: int = 4096
    api_key: Optional[str] = None


class RoleConfig(BaseModel):
    """Configuration for an agent role."""

    llm: LLMConfig
    system_prompt: str = ""
    tools: list[str] = Field(default_factory=list)


class VermasSettings(BaseSettings):
    """Global VerMAS settings."""

    # Database
    database_url: str = "sqlite:///vermas.db"

    # LLM defaults
    default_provider: str = "anthropic"
    anthropic_api_key: Optional[str] = None
    openai_api_key: Optional[str] = None

    # Role-specific LLMs (for independence)
    roles: dict[str, RoleConfig] = Field(default_factory=lambda: {
        "mayor": RoleConfig(llm=LLMConfig(model="claude-sonnet-4-20250514")),
        "designer": RoleConfig(llm=LLMConfig(model="claude-sonnet-4-20250514")),
        "inspector": RoleConfig(llm=LLMConfig(model="claude-sonnet-4-20250514")),
        "strategist": RoleConfig(llm=LLMConfig(provider="openai", model="gpt-4o")),
        "advocate": RoleConfig(llm=LLMConfig(model="claude-sonnet-4-20250514")),
        "critic": RoleConfig(llm=LLMConfig(provider="openai", model="gpt-4o")),
        "judge": RoleConfig(llm=LLMConfig(model="claude-sonnet-4-20250514")),
    })

    # Verification
    require_spec_approval: bool = True
    default_confidence_threshold: float = 0.7
    verification_timeout_seconds: int = 300

    class Config:
        env_prefix = "VERMAS_"
        env_file = ".env"
```

---

## LangGraph Workflows

### Mayor Graph

```python
# vermas/graphs/mayor.py
from typing import TypedDict, Annotated
from langgraph.graph import StateGraph, END
from langgraph.prebuilt import ToolNode
from langchain_core.messages import HumanMessage, AIMessage

from vermas.models.agent import AgentState, AgentRole
from vermas.models.bead import WorkItem, Requirements
from vermas.agents.designer import DesignerAgent
from vermas.tools import beads_tools, mail_tools, git_tools


class MayorState(TypedDict):
    """State for Mayor graph."""
    messages: list
    current_request: str
    requirements: Requirements | None
    work_item: WorkItem | None
    spec_approved: bool
    next_action: str


def create_mayor_graph() -> StateGraph:
    """Create the Mayor workflow graph."""

    # Define the graph
    graph = StateGraph(MayorState)

    # Add nodes
    graph.add_node("receive_request", receive_request)
    graph.add_node("design_requirements", design_requirements)
    graph.add_node("approve_requirements", approve_requirements)
    graph.add_node("create_work_item", create_work_item)
    graph.add_node("wait_for_spec", wait_for_spec)
    graph.add_node("sling_work", sling_work)
    graph.add_node("handle_rejection", handle_rejection)
    graph.add_node("tools", ToolNode(tools=beads_tools + mail_tools))

    # Add edges
    graph.set_entry_point("receive_request")

    graph.add_edge("receive_request", "design_requirements")
    graph.add_edge("design_requirements", "approve_requirements")

    graph.add_conditional_edges(
        "approve_requirements",
        lambda s: "create" if s["requirements"] else "design",
        {
            "create": "create_work_item",
            "design": "design_requirements",  # Loop back if rejected
        }
    )

    graph.add_edge("create_work_item", "wait_for_spec")

    graph.add_conditional_edges(
        "wait_for_spec",
        lambda s: "sling" if s["spec_approved"] else "wait",
        {
            "sling": "sling_work",
            "wait": "wait_for_spec",
        }
    )

    graph.add_edge("sling_work", END)
    graph.add_edge("handle_rejection", "sling_work")  # Re-sling after fix

    return graph.compile()


async def receive_request(state: MayorState) -> MayorState:
    """Receive and parse user request."""
    # Extract the request from messages
    last_message = state["messages"][-1]
    state["current_request"] = last_message.content
    state["next_action"] = "design"
    return state


async def design_requirements(state: MayorState) -> MayorState:
    """Use Designer to elaborate requirements."""
    designer = DesignerAgent()
    requirements = await designer.elaborate(state["current_request"])
    state["requirements"] = requirements

    # Add message showing requirements to user
    state["messages"].append(AIMessage(content=f"""
I've designed these requirements for your request:

**Functional:**
{chr(10).join(f'- {r}' for r in requirements.functional)}

**Non-functional:**
{chr(10).join(f'- {r}' for r in requirements.non_functional)}

**Constraints:**
{chr(10).join(f'- {r}' for r in requirements.constraints)}

Do these look correct? Reply 'yes' to approve or provide modifications.
"""))

    return state


async def approve_requirements(state: MayorState) -> MayorState:
    """Get user approval on requirements."""
    last_message = state["messages"][-1]

    if isinstance(last_message, HumanMessage):
        content = last_message.content.lower()
        if content in ("yes", "y", "approve", "looks good", "lgtm"):
            # Approved, keep requirements
            pass
        elif content in ("no", "n", "reject"):
            # Rejected, clear requirements to loop back
            state["requirements"] = None
        else:
            # Modification request - update requirements
            # TODO: Parse modifications and update
            pass

    return state


async def create_work_item(state: MayorState) -> MayorState:
    """Create work item and trigger spec creation."""
    from vermas.db.repository import BeadsRepository

    repo = BeadsRepository()

    # Create work item
    work_item = WorkItem(
        id=repo.generate_id("gt"),
        type="task",
        title=state["current_request"][:100],
        description=format_requirements(state["requirements"]),
    )
    work_item = await repo.create_work_item(work_item)
    state["work_item"] = work_item

    # Spec is auto-created by repository hook
    # Wait for Inspector to approve it
    state["spec_approved"] = False

    state["messages"].append(AIMessage(content=f"""
Created work item: {work_item.id}
Title: {work_item.title}

Waiting for Inspector to define verification criteria...
"""))

    return state


async def wait_for_spec(state: MayorState) -> MayorState:
    """Wait for spec approval from Inspector."""
    from vermas.db.repository import BeadsRepository

    repo = BeadsRepository()
    spec_id = f"spec-{state['work_item'].id[3:]}"

    spec = await repo.get_test_spec(spec_id)
    if spec and spec.status == "closed":
        state["spec_approved"] = True
        state["messages"].append(AIMessage(content=f"""
✅ Test spec {spec_id} approved by Inspector!
Criteria defined: {len(spec.criteria)}

Ready to dispatch work.
"""))

    return state


async def sling_work(state: MayorState) -> MayorState:
    """Dispatch work to a polecat."""
    from vermas.agents.polecat import dispatch_to_polecat

    await dispatch_to_polecat(state["work_item"])

    state["messages"].append(AIMessage(content=f"""
✅ Work {state['work_item'].id} dispatched to polecat.
Monitoring for completion...
"""))

    return state
```

### Inspector Graph

```python
# vermas/graphs/inspector.py
from typing import TypedDict
from langgraph.graph import StateGraph, END

from vermas.models.bead import TestSpec, AcceptanceCriterion
from vermas.models.verification import Evidence, Brief, VerificationResult, Verdict
from vermas.agents.strategist import StrategistAgent
from vermas.agents.verifier import VerifierAgent
from vermas.agents.auditor import AuditorAgent
from vermas.agents.advocate import AdvocateAgent
from vermas.agents.critic import CriticAgent
from vermas.agents.judge import JudgeAgent


class InspectorState(TypedDict):
    """State for Inspector graph."""
    messages: list
    work_item_id: str
    test_spec: TestSpec | None
    proposed_criteria: list[AcceptanceCriterion]
    evidence: Evidence | None
    advocate_brief: Brief | None
    critic_brief: Brief | None
    verdict: VerificationResult | None
    phase: str  # "planning" or "verification"


def create_inspector_graph() -> StateGraph:
    """Create the Inspector workflow graph."""

    graph = StateGraph(InspectorState)

    # Planning phase nodes
    graph.add_node("receive_spec_request", receive_spec_request)
    graph.add_node("propose_criteria", propose_criteria)
    graph.add_node("approve_criteria", approve_criteria)
    graph.add_node("finalize_spec", finalize_spec)

    # Verification phase nodes
    graph.add_node("receive_verification_request", receive_verification_request)
    graph.add_node("run_verifier", run_verifier)
    graph.add_node("run_auditor", run_auditor)
    graph.add_node("run_advocate", run_advocate)
    graph.add_node("run_critic", run_critic)
    graph.add_node("run_judge", run_judge)
    graph.add_node("route_verdict", route_verdict)

    # Entry point depends on phase
    graph.set_conditional_entry_point(
        lambda s: "receive_spec_request" if s["phase"] == "planning" else "receive_verification_request"
    )

    # Planning phase edges
    graph.add_edge("receive_spec_request", "propose_criteria")
    graph.add_edge("propose_criteria", "approve_criteria")

    graph.add_conditional_edges(
        "approve_criteria",
        lambda s: "finalize" if s["proposed_criteria"] else "propose",
        {
            "finalize": "finalize_spec",
            "propose": "propose_criteria",
        }
    )

    graph.add_edge("finalize_spec", END)

    # Verification phase edges
    graph.add_edge("receive_verification_request", "run_verifier")
    graph.add_edge("run_verifier", "run_auditor")
    graph.add_edge("run_auditor", "run_advocate")
    graph.add_edge("run_auditor", "run_critic")  # Parallel with advocate
    graph.add_edge("run_advocate", "run_judge")
    graph.add_edge("run_critic", "run_judge")
    graph.add_edge("run_judge", "route_verdict")
    graph.add_edge("route_verdict", END)

    return graph.compile()


async def propose_criteria(state: InspectorState) -> InspectorState:
    """Use Strategist to propose verification criteria."""
    from vermas.db.repository import BeadsRepository

    repo = BeadsRepository()
    work_item = await repo.get_work_item(state["work_item_id"])

    strategist = StrategistAgent()
    criteria = await strategist.propose_criteria(work_item)

    state["proposed_criteria"] = criteria
    state["messages"].append({
        "role": "assistant",
        "content": format_criteria_proposal(criteria)
    })

    return state


async def run_verifier(state: InspectorState) -> InspectorState:
    """Run objective tests."""
    verifier = VerifierAgent(working_directory=state.get("working_directory", "."))
    results = await verifier.run_criteria(state["test_spec"].criteria)

    state["evidence"] = Evidence(
        work_item=state["work_item"],
        test_spec=state["test_spec"],
        criterion_results=results,
        diff=await get_git_diff(),
        commit_messages=await get_commit_messages(),
        files_changed=await get_files_changed(),
    )

    return state


async def run_advocate(state: InspectorState) -> InspectorState:
    """Build defense brief."""
    advocate = AdvocateAgent()
    brief = await advocate.build_defense(state["evidence"])
    state["advocate_brief"] = brief
    return state


async def run_critic(state: InspectorState) -> InspectorState:
    """Build prosecution brief."""
    critic = CriticAgent()
    brief = await critic.build_prosecution(state["evidence"])
    state["critic_brief"] = brief
    return state


async def run_judge(state: InspectorState) -> InspectorState:
    """Deliver verdict."""
    judge = JudgeAgent()
    verdict = await judge.deliberate(
        evidence=state["evidence"],
        advocate_brief=state["advocate_brief"],
        critic_brief=state["critic_brief"],
    )
    state["verdict"] = verdict
    return state


async def route_verdict(state: InspectorState) -> InspectorState:
    """Route based on verdict."""
    from vermas.messaging.mailbox import send_message

    verdict = state["verdict"]

    if verdict.verdict == Verdict.PASS:
        await send_message(
            to="refinery",
            subject=f"✅ VERIFICATION PASSED: {state['work_item_id']}",
            body=f"Proceed with merge.\nConfidence: {verdict.confidence}",
        )
    elif verdict.verdict == Verdict.FAIL:
        await send_message(
            to="mayor",
            subject=f"❌ VERIFICATION FAILED: {state['work_item_id']}",
            body=format_failure_report(verdict),
        )
    else:  # NEEDS_HUMAN
        await send_message(
            to="mayor",
            subject=f"❓ NEEDS_CLARIFICATION: {state['work_item_id']}",
            body=verdict.judge_reasoning,
        )

    return state
```

### Verification Subgraph

```python
# vermas/graphs/verification.py
from typing import TypedDict, Annotated
from langgraph.graph import StateGraph, END
import operator

from vermas.models.verification import Evidence, Brief, VerificationResult, Verdict


class VerificationState(TypedDict):
    """State for verification subgraph."""
    evidence: Evidence
    advocate_brief: Annotated[Brief | None, operator.or_]
    critic_brief: Annotated[Brief | None, operator.or_]
    verdict: VerificationResult | None


def create_verification_subgraph() -> StateGraph:
    """
    Create verification subgraph with parallel Advocate/Critic.

    This demonstrates LangGraph's ability to run nodes in parallel
    and then converge to a single node (Judge).
    """

    graph = StateGraph(VerificationState)

    graph.add_node("advocate", run_advocate_node)
    graph.add_node("critic", run_critic_node)
    graph.add_node("judge", run_judge_node)

    # Both advocate and critic run from entry
    graph.set_entry_point("advocate")
    graph.set_entry_point("critic")  # Parallel entry

    # Both converge to judge
    graph.add_edge("advocate", "judge")
    graph.add_edge("critic", "judge")

    graph.add_edge("judge", END)

    return graph.compile()


async def run_advocate_node(state: VerificationState) -> dict:
    """Run Advocate agent."""
    from vermas.agents.advocate import AdvocateAgent

    advocate = AdvocateAgent()
    brief = await advocate.build_defense(state["evidence"])
    return {"advocate_brief": brief}


async def run_critic_node(state: VerificationState) -> dict:
    """Run Critic agent."""
    from vermas.agents.critic import CriticAgent

    critic = CriticAgent()
    brief = await critic.build_prosecution(state["evidence"])
    return {"critic_brief": brief}


async def run_judge_node(state: VerificationState) -> dict:
    """Run Judge agent - waits for both briefs."""
    from vermas.agents.judge import JudgeAgent

    judge = JudgeAgent()
    verdict = await judge.deliberate(
        evidence=state["evidence"],
        advocate_brief=state["advocate_brief"],
        critic_brief=state["critic_brief"],
    )
    return {"verdict": verdict}
```

---

## Agent Implementations

### Base Agent

```python
# vermas/agents/base.py
from abc import ABC, abstractmethod
from typing import Any

from langchain_core.language_models import BaseChatModel
from langchain_anthropic import ChatAnthropic
from langchain_openai import ChatOpenAI

from vermas.models.config import VermasSettings, LLMConfig


class BaseAgent(ABC):
    """Base class for all VerMAS agents."""

    def __init__(self, role: str):
        self.role = role
        self.settings = VermasSettings()
        self.llm = self._create_llm()

    def _create_llm(self) -> BaseChatModel:
        """Create LLM based on role configuration."""
        role_config = self.settings.roles.get(self.role)
        if not role_config:
            # Default to Claude
            return ChatAnthropic(
                model="claude-sonnet-4-20250514",
                api_key=self.settings.anthropic_api_key,
            )

        llm_config = role_config.llm

        if llm_config.provider == "anthropic":
            return ChatAnthropic(
                model=llm_config.model,
                temperature=llm_config.temperature,
                max_tokens=llm_config.max_tokens,
                api_key=llm_config.api_key or self.settings.anthropic_api_key,
            )
        elif llm_config.provider == "openai":
            return ChatOpenAI(
                model=llm_config.model,
                temperature=llm_config.temperature,
                max_tokens=llm_config.max_tokens,
                api_key=llm_config.api_key or self.settings.openai_api_key,
            )
        else:
            raise ValueError(f"Unknown provider: {llm_config.provider}")

    @abstractmethod
    async def run(self, *args, **kwargs) -> Any:
        """Run the agent's main task."""
        pass
```

### Designer Agent

```python
# vermas/agents/designer.py
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import PydanticOutputParser

from vermas.agents.base import BaseAgent
from vermas.models.bead import Requirements


DESIGNER_PROMPT = """You are the Designer. Your job is to turn vague requests into detailed requirements.

User request: {request}

Produce requirements covering:
1. Functional - What exactly should this do? Be specific.
2. Non-functional - Performance, quality, style requirements
3. Constraints - What should it NOT do?
4. Examples - Concrete input/output examples

Be specific enough that a developer could implement without asking questions.
Be concise - don't over-engineer simple requests.

{format_instructions}
"""


class DesignerAgent(BaseAgent):
    """Agent that elaborates vague requests into detailed requirements."""

    def __init__(self):
        super().__init__(role="designer")
        self.parser = PydanticOutputParser(pydantic_object=Requirements)
        self.prompt = ChatPromptTemplate.from_template(DESIGNER_PROMPT)

    async def elaborate(self, request: str) -> Requirements:
        """Elaborate a vague request into detailed requirements."""
        chain = self.prompt | self.llm | self.parser

        result = await chain.ainvoke({
            "request": request,
            "format_instructions": self.parser.get_format_instructions(),
        })

        return result

    async def run(self, request: str) -> Requirements:
        return await self.elaborate(request)
```

### Strategist Agent

```python
# vermas/agents/strategist.py
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import PydanticOutputParser

from vermas.agents.base import BaseAgent
from vermas.models.bead import WorkItem, AcceptanceCriterion


STRATEGIST_PROMPT = """You are the Strategist. Your job is to design verification criteria.

Work Item: {title}
Description:
{description}

For each key behavior, create a verification criterion:
- ID: AC-N (e.g., AC-1, AC-2)
- Description: What we're checking (human readable)
- Verify Command: Exact bash command that exits 0 for pass, non-0 for fail

Focus on:
- Correctness (does it work?)
- Edge cases (does it handle boundaries?)
- Quality (is the code clean?)

Don't over-test. Cover the critical paths. 5-10 criteria is usually enough.

{format_instructions}
"""


class StrategistAgent(BaseAgent):
    """Agent that proposes verification criteria from requirements."""

    def __init__(self):
        super().__init__(role="strategist")
        self.parser = PydanticOutputParser(pydantic_object=list[AcceptanceCriterion])
        self.prompt = ChatPromptTemplate.from_template(STRATEGIST_PROMPT)

    async def propose_criteria(self, work_item: WorkItem) -> list[AcceptanceCriterion]:
        """Propose verification criteria for a work item."""
        chain = self.prompt | self.llm | self.parser

        result = await chain.ainvoke({
            "title": work_item.title,
            "description": work_item.description,
            "format_instructions": self.parser.get_format_instructions(),
        })

        return result

    async def run(self, work_item: WorkItem) -> list[AcceptanceCriterion]:
        return await self.propose_criteria(work_item)
```

### Verifier Agent (No LLM)

```python
# vermas/agents/verifier.py
import asyncio
import subprocess
from dataclasses import dataclass
from typing import Optional

from vermas.models.bead import AcceptanceCriterion
from vermas.models.verification import CriterionResult


class VerifierAgent:
    """
    Agent that runs objective verification tests.

    This agent does NOT use an LLM - it simply executes shell commands
    and reports pass/fail based on exit codes.
    """

    def __init__(self, working_directory: str = "."):
        self.working_directory = working_directory

    async def run_criteria(
        self,
        criteria: list[AcceptanceCriterion]
    ) -> list[CriterionResult]:
        """Run all acceptance criteria and collect results."""
        results = []

        for criterion in criteria:
            result = await self.run_single_criterion(criterion)
            results.append(result)

        return results

    async def run_single_criterion(
        self,
        criterion: AcceptanceCriterion
    ) -> CriterionResult:
        """Run a single criterion and return the result."""
        import time

        start_time = time.time()

        try:
            process = await asyncio.create_subprocess_shell(
                criterion.verify_command,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.STDOUT,
                cwd=self.working_directory,
            )

            try:
                stdout, _ = await asyncio.wait_for(
                    process.communicate(),
                    timeout=criterion.timeout_seconds,
                )
                output = stdout.decode("utf-8", errors="replace")
                exit_code = process.returncode
                status = "pass" if exit_code == criterion.expected_exit_code else "fail"

            except asyncio.TimeoutError:
                process.kill()
                output = f"TIMEOUT after {criterion.timeout_seconds}s"
                exit_code = -1
                status = "timeout"

        except Exception as e:
            output = f"ERROR: {str(e)}"
            exit_code = -1
            status = "error"

        duration_ms = int((time.time() - start_time) * 1000)

        return CriterionResult(
            criterion_id=criterion.id,
            status=status,
            output=output,
            exit_code=exit_code,
            duration_ms=duration_ms,
        )

    async def run(self, criteria: list[AcceptanceCriterion]) -> list[CriterionResult]:
        return await self.run_criteria(criteria)
```

### Advocate Agent

```python
# vermas/agents/advocate.py
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import PydanticOutputParser

from vermas.agents.base import BaseAgent
from vermas.models.verification import Evidence, Brief


ADVOCATE_PROMPT = """You are the Advocate. Your job is to DEFEND this code change.

## Evidence

### Work Item
{work_item_title}
{work_item_description}

### Test Results
{test_results}

### Code Changes
{diff}

## Your Task

Argue WHY this code should be merged:
- Highlight strengths
- Explain design decisions
- Mitigate any failing tests or concerns
- Show requirement compliance

Be persuasive but honest. Don't defend genuinely bad code.

{format_instructions}
"""


class AdvocateAgent(BaseAgent):
    """Agent that builds a defense for the code."""

    def __init__(self):
        super().__init__(role="advocate")
        self.parser = PydanticOutputParser(pydantic_object=Brief)
        self.prompt = ChatPromptTemplate.from_template(ADVOCATE_PROMPT)

    async def build_defense(self, evidence: Evidence) -> Brief:
        """Build a defense brief for the code."""
        chain = self.prompt | self.llm | self.parser

        result = await chain.ainvoke({
            "work_item_title": evidence.work_item.title,
            "work_item_description": evidence.work_item.description,
            "test_results": format_test_results(evidence.criterion_results),
            "diff": evidence.diff[:10000],  # Truncate large diffs
            "format_instructions": self.parser.get_format_instructions(),
        })

        result.role = "advocate"
        result.position = "FOR"
        return result

    async def run(self, evidence: Evidence) -> Brief:
        return await self.build_defense(evidence)


def format_test_results(results: list) -> str:
    """Format test results for prompt."""
    lines = []
    for r in results:
        status_emoji = "✅" if r.status == "pass" else "❌"
        lines.append(f"{status_emoji} {r.criterion_id}: {r.status}")
        if r.status != "pass":
            lines.append(f"   Output: {r.output[:200]}")
    return "\n".join(lines)
```

### Critic Agent

```python
# vermas/agents/critic.py
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import PydanticOutputParser

from vermas.agents.base import BaseAgent
from vermas.models.verification import Evidence, Brief


CRITIC_PROMPT = """You are the Critic. Your job is to ATTACK this code change.

## Evidence

### Work Item
{work_item_title}
{work_item_description}

### Test Results
{test_results}

### Code Changes
{diff}

## Your Task

Argue WHY this code should NOT be merged:
- Find bugs
- Identify security issues
- Question design decisions
- Note missing tests
- Flag performance concerns
- Point out any failing criteria

Be thorough but fair. Don't manufacture false concerns.

{format_instructions}
"""


class CriticAgent(BaseAgent):
    """Agent that builds a prosecution against the code."""

    def __init__(self):
        super().__init__(role="critic")
        self.parser = PydanticOutputParser(pydantic_object=Brief)
        self.prompt = ChatPromptTemplate.from_template(CRITIC_PROMPT)

    async def build_prosecution(self, evidence: Evidence) -> Brief:
        """Build a prosecution brief against the code."""
        chain = self.prompt | self.llm | self.parser

        result = await chain.ainvoke({
            "work_item_title": evidence.work_item.title,
            "work_item_description": evidence.work_item.description,
            "test_results": format_test_results(evidence.criterion_results),
            "diff": evidence.diff[:10000],
            "format_instructions": self.parser.get_format_instructions(),
        })

        result.role = "critic"
        result.position = "AGAINST"
        return result

    async def run(self, evidence: Evidence) -> Brief:
        return await self.build_prosecution(evidence)
```

### Judge Agent

```python
# vermas/agents/judge.py
from langchain_core.prompts import ChatPromptTemplate
from langchain_core.output_parsers import PydanticOutputParser

from vermas.agents.base import BaseAgent
from vermas.models.verification import Evidence, Brief, VerificationResult, Verdict


JUDGE_PROMPT = """You are the Judge. Your job is to deliver a VERDICT.

## Evidence

### Test Results
{test_results}

### Advocate's Defense
{advocate_brief}

### Critic's Concerns
{critic_brief}

## Your Task

Weigh the evidence and decide:
- **PASS**: Code meets requirements, concerns adequately addressed
- **FAIL**: Critical issues must be fixed before merge
- **NEEDS_HUMAN**: Cannot decide autonomously, escalate to human

Provide:
1. Your verdict (PASS, FAIL, or NEEDS_HUMAN)
2. Confidence score (0.0-1.0)
3. Clear reasoning for your decision
4. If FAIL: List of required fixes
5. If NEEDS_HUMAN: The specific question you need answered

Be strict but fair. Only PASS if genuinely production-ready.

{format_instructions}
"""


class JudgeAgent(BaseAgent):
    """Agent that delivers the final verdict."""

    def __init__(self):
        super().__init__(role="judge")
        self.parser = PydanticOutputParser(pydantic_object=VerificationResult)
        self.prompt = ChatPromptTemplate.from_template(JUDGE_PROMPT)

    async def deliberate(
        self,
        evidence: Evidence,
        advocate_brief: Brief,
        critic_brief: Brief,
    ) -> VerificationResult:
        """Deliberate and deliver a verdict."""
        chain = self.prompt | self.llm | self.parser

        result = await chain.ainvoke({
            "test_results": format_test_results(evidence.criterion_results),
            "advocate_brief": format_brief(advocate_brief),
            "critic_brief": format_brief(critic_brief),
            "format_instructions": self.parser.get_format_instructions(),
        })

        # Attach the briefs for audit trail
        result.work_item_id = evidence.work_item.id
        result.spec_id = evidence.test_spec.id
        result.advocate_brief = advocate_brief
        result.critic_brief = critic_brief
        result.criterion_results = evidence.criterion_results

        return result

    async def run(
        self,
        evidence: Evidence,
        advocate_brief: Brief,
        critic_brief: Brief
    ) -> VerificationResult:
        return await self.deliberate(evidence, advocate_brief, critic_brief)


def format_brief(brief: Brief) -> str:
    """Format a brief for the judge prompt."""
    return f"""
Position: {brief.position}
Confidence: {brief.confidence}

Arguments:
{chr(10).join(f'- {a}' for a in brief.arguments)}

Summary: {brief.summary}
"""
```

---

## CLI Interface

```python
# vermas/cli.py
import typer
from rich.console import Console
from rich.table import Table

app = typer.Typer(name="vermas", help="Verifiable Multi-Agent System")
console = Console()

# Sub-apps
inspect_app = typer.Typer(help="Inspector commands")
app.add_typer(inspect_app, name="inspect")


@app.command()
def create(
    title: str = typer.Argument(..., help="Work item title"),
    type: str = typer.Option("task", help="Type: task, feature, bug"),
):
    """Create a new work item with AI-assisted requirements."""
    import asyncio
    from vermas.graphs.mayor import create_work_item_flow

    asyncio.run(create_work_item_flow(title, type))


@app.command()
def sling(
    bead_id: str = typer.Argument(..., help="Work item ID"),
    force: bool = typer.Option(False, "--force", help="Bypass spec gate"),
    reason: str = typer.Option("", "--reason", help="Reason for force bypass"),
):
    """Dispatch work to a polecat."""
    import asyncio
    from vermas.db.repository import BeadsRepository

    async def do_sling():
        repo = BeadsRepository()

        # Check spec gate
        if not force:
            spec_id = f"spec-{bead_id[3:]}"
            spec = await repo.get_test_spec(spec_id)

            if not spec:
                console.print(f"[red]ERROR: No test spec found for {bead_id}[/red]")
                console.print(f"Create with: vermas inspect create-spec {bead_id}")
                raise typer.Exit(1)

            if spec.status != "closed":
                console.print(f"[red]ERROR: Test spec {spec_id} not approved[/red]")
                console.print(f"Status: {spec.status}")
                console.print(f"Action: vermas inspect approve {spec_id}")
                raise typer.Exit(1)

        elif not reason:
            console.print("[red]ERROR: --reason required when using --force[/red]")
            raise typer.Exit(1)

        else:
            console.print(f"[yellow]⚠️  GATE BYPASS: {reason}[/yellow]")

        # Dispatch
        from vermas.agents.polecat import dispatch_to_polecat
        await dispatch_to_polecat(bead_id)
        console.print(f"[green]✅ Work {bead_id} dispatched[/green]")

    asyncio.run(do_sling())


@inspect_app.command("pending")
def inspect_pending():
    """List specs needing criteria."""
    import asyncio
    from vermas.db.repository import BeadsRepository

    async def show_pending():
        repo = BeadsRepository()
        specs = await repo.get_pending_specs()

        if not specs:
            console.print("[dim]No pending specs[/dim]")
            return

        table = Table(title="Pending Test Specs")
        table.add_column("Spec ID")
        table.add_column("Work Item")
        table.add_column("Criteria")
        table.add_column("Status")

        for spec in specs:
            table.add_row(
                spec.id,
                spec.parent_work_item,
                str(len(spec.criteria)),
                spec.status,
            )

        console.print(table)

    asyncio.run(show_pending())


@inspect_app.command("approve")
def inspect_approve(spec_id: str = typer.Argument(..., help="Spec ID")):
    """Approve a test spec."""
    import asyncio
    from vermas.db.repository import BeadsRepository
    from vermas.messaging.mailbox import send_message

    async def do_approve():
        repo = BeadsRepository()
        spec = await repo.get_test_spec(spec_id)

        if not spec:
            console.print(f"[red]Spec {spec_id} not found[/red]")
            raise typer.Exit(1)

        if not spec.criteria:
            console.print(f"[red]Cannot approve spec with no criteria[/red]")
            raise typer.Exit(1)

        await repo.approve_spec(spec_id)

        # Notify Mayor
        await send_message(
            to="mayor",
            subject=f"✅ SPEC APPROVED: {spec_id}",
            body=f"Test spec approved with {len(spec.criteria)} criteria.\n"
                 f"Work item {spec.parent_work_item} is now READY.",
        )

        console.print(f"[green]✅ Spec {spec_id} approved[/green]")
        console.print(f"Work {spec.parent_work_item} is now ready for dispatch")

    asyncio.run(do_approve())


@inspect_app.command("add-criteria")
def inspect_add_criteria(
    spec_id: str = typer.Argument(..., help="Spec ID"),
    description: str = typer.Option(..., "-d", "--description", help="Criterion description"),
    verify: str = typer.Option(..., "-v", "--verify", help="Verification command"),
):
    """Add an acceptance criterion to a spec."""
    import asyncio
    from vermas.db.repository import BeadsRepository
    from vermas.models.bead import AcceptanceCriterion

    async def do_add():
        repo = BeadsRepository()
        spec = await repo.get_test_spec(spec_id)

        if not spec:
            console.print(f"[red]Spec {spec_id} not found[/red]")
            raise typer.Exit(1)

        # Generate criterion ID
        criterion_id = f"AC-{len(spec.criteria) + 1}"

        criterion = AcceptanceCriterion(
            id=criterion_id,
            description=description,
            verify_command=verify,
        )

        await repo.add_criterion(spec_id, criterion)
        console.print(f"[green]✅ Added criterion {criterion_id}[/green]")

    asyncio.run(do_add())


@inspect_app.command("run")
def inspect_run(bead_id: str = typer.Argument(..., help="Work item ID")):
    """Run full verification workflow."""
    import asyncio
    from vermas.graphs.inspector import create_inspector_graph

    async def do_run():
        graph = create_inspector_graph()

        result = await graph.ainvoke({
            "work_item_id": bead_id,
            "phase": "verification",
            "messages": [],
        })

        verdict = result["verdict"]

        if verdict.verdict == "PASS":
            console.print(f"[green]✅ PASS (confidence: {verdict.confidence})[/green]")
        elif verdict.verdict == "FAIL":
            console.print(f"[red]❌ FAIL (confidence: {verdict.confidence})[/red]")
            console.print("\nRequired fixes:")
            for fix in verdict.required_fixes:
                console.print(f"  - {fix}")
        else:
            console.print(f"[yellow]❓ NEEDS_HUMAN[/yellow]")
            console.print(f"\n{verdict.judge_reasoning}")

    asyncio.run(do_run())


@app.command()
def start():
    """Start VerMAS with Mayor and Inspector."""
    import asyncio
    from vermas.graphs.mayor import create_mayor_graph
    from vermas.graphs.inspector import create_inspector_graph

    console.print("[bold]Starting VerMAS...[/bold]")
    console.print("  Mayor: [green]ready[/green]")
    console.print("  Inspector: [green]ready[/green]")
    console.print("\nType 'help' for commands, or enter a task to begin.")

    # Interactive loop
    mayor_graph = create_mayor_graph()

    async def interactive():
        state = {"messages": [], "phase": "idle"}

        while True:
            try:
                user_input = console.input("[bold cyan]You:[/bold cyan] ")

                if user_input.lower() in ("exit", "quit", "q"):
                    break

                if user_input.lower() == "help":
                    console.print("""
Commands:
  create <task>  - Create a new work item
  status         - Show system status
  inbox          - Check mail
  exit           - Exit VerMAS
""")
                    continue

                # Run through Mayor graph
                state["messages"].append({"role": "user", "content": user_input})
                state = await mayor_graph.ainvoke(state)

                # Print last AI message
                for msg in reversed(state["messages"]):
                    if msg.get("role") == "assistant":
                        console.print(f"\n[bold green]Mayor:[/bold green] {msg['content']}\n")
                        break

            except KeyboardInterrupt:
                break

    asyncio.run(interactive())


if __name__ == "__main__":
    app()
```

---

## Database Layer

```python
# vermas/db/models.py
from datetime import datetime
from typing import Optional
from sqlmodel import SQLModel, Field, Relationship


class WorkItemDB(SQLModel, table=True):
    """Database model for work items."""

    __tablename__ = "work_items"

    id: str = Field(primary_key=True)
    type: str
    title: str
    description: str = ""
    status: str = "open"
    priority: int = 2
    assignee: Optional[str] = None
    parent: Optional[str] = None
    created_at: datetime = Field(default_factory=datetime.utcnow)
    updated_at: datetime = Field(default_factory=datetime.utcnow)

    # Relationships
    test_spec: Optional["TestSpecDB"] = Relationship(back_populates="work_item")


class TestSpecDB(SQLModel, table=True):
    """Database model for test specs."""

    __tablename__ = "test_specs"

    id: str = Field(primary_key=True)
    work_item_id: str = Field(foreign_key="work_items.id")
    status: str = "open"
    success_threshold: float = 1.0
    created_at: datetime = Field(default_factory=datetime.utcnow)
    approved_at: Optional[datetime] = None
    approved_by: Optional[str] = None

    # Relationships
    work_item: Optional[WorkItemDB] = Relationship(back_populates="test_spec")
    criteria: list["CriterionDB"] = Relationship(back_populates="test_spec")


class CriterionDB(SQLModel, table=True):
    """Database model for acceptance criteria."""

    __tablename__ = "criteria"

    id: str = Field(primary_key=True)
    spec_id: str = Field(foreign_key="test_specs.id")
    description: str
    verify_command: str
    expected_exit_code: int = 0
    timeout_seconds: int = 60

    # Relationships
    test_spec: Optional[TestSpecDB] = Relationship(back_populates="criteria")


class MessageDB(SQLModel, table=True):
    """Database model for messages."""

    __tablename__ = "messages"

    id: str = Field(primary_key=True)
    from_agent: str
    to_agent: str
    subject: str
    body: str
    message_type: str = "notification"
    priority: str = "normal"
    thread_id: Optional[str] = None
    reply_to: Optional[str] = None
    created_at: datetime = Field(default_factory=datetime.utcnow)
    read: bool = False
```

```python
# vermas/db/repository.py
from datetime import datetime
from typing import Optional
import uuid

from sqlmodel import Session, select
from sqlalchemy import create_engine

from vermas.db.models import WorkItemDB, TestSpecDB, CriterionDB, MessageDB
from vermas.models.bead import WorkItem, TestSpec, AcceptanceCriterion
from vermas.models.config import VermasSettings


class BeadsRepository:
    """Repository for beads (work items, specs, etc.)."""

    def __init__(self):
        settings = VermasSettings()
        self.engine = create_engine(settings.database_url)

    def generate_id(self, prefix: str) -> str:
        """Generate a unique ID with prefix."""
        suffix = uuid.uuid4().hex[:8]
        return f"{prefix}-{suffix}"

    async def create_work_item(self, work_item: WorkItem) -> WorkItem:
        """Create a work item and auto-create its test spec."""
        with Session(self.engine) as session:
            # Create work item
            db_item = WorkItemDB(**work_item.model_dump())
            session.add(db_item)

            # Auto-create test spec
            spec_id = f"spec-{work_item.id[3:]}"
            db_spec = TestSpecDB(
                id=spec_id,
                work_item_id=work_item.id,
                status="open",
            )
            session.add(db_spec)

            session.commit()
            session.refresh(db_item)

            return WorkItem.model_validate(db_item)

    async def get_work_item(self, id: str) -> Optional[WorkItem]:
        """Get a work item by ID."""
        with Session(self.engine) as session:
            item = session.get(WorkItemDB, id)
            if item:
                return WorkItem.model_validate(item)
            return None

    async def get_test_spec(self, id: str) -> Optional[TestSpec]:
        """Get a test spec by ID."""
        with Session(self.engine) as session:
            spec = session.get(TestSpecDB, id)
            if spec:
                criteria = [
                    AcceptanceCriterion.model_validate(c)
                    for c in spec.criteria
                ]
                return TestSpec(
                    id=spec.id,
                    parent_work_item=spec.work_item_id,
                    status=spec.status,
                    criteria=criteria,
                    success_threshold=spec.success_threshold,
                    created_at=spec.created_at,
                    approved_at=spec.approved_at,
                    approved_by=spec.approved_by,
                )
            return None

    async def get_pending_specs(self) -> list[TestSpec]:
        """Get all pending (unapproved) test specs."""
        with Session(self.engine) as session:
            statement = select(TestSpecDB).where(TestSpecDB.status == "open")
            specs = session.exec(statement).all()
            return [await self.get_test_spec(s.id) for s in specs]

    async def add_criterion(self, spec_id: str, criterion: AcceptanceCriterion):
        """Add a criterion to a test spec."""
        with Session(self.engine) as session:
            db_criterion = CriterionDB(
                id=criterion.id,
                spec_id=spec_id,
                description=criterion.description,
                verify_command=criterion.verify_command,
                expected_exit_code=criterion.expected_exit_code,
                timeout_seconds=criterion.timeout_seconds,
            )
            session.add(db_criterion)
            session.commit()

    async def approve_spec(self, spec_id: str, approved_by: str = "inspector"):
        """Approve a test spec."""
        with Session(self.engine) as session:
            spec = session.get(TestSpecDB, spec_id)
            if spec:
                spec.status = "closed"
                spec.approved_at = datetime.utcnow()
                spec.approved_by = approved_by
                session.add(spec)
                session.commit()
```

---

## Example: FizzBuzz

```python
# examples/fizzbuzz.py
"""
FizzBuzz example demonstrating the full VerMAS flow.
"""
import asyncio
from rich.console import Console

from vermas.graphs.mayor import create_mayor_graph
from vermas.graphs.inspector import create_inspector_graph


console = Console()


async def main():
    """Run the FizzBuzz example."""

    console.print("[bold]VerMAS FizzBuzz Example[/bold]\n")

    # Step 1: User request to Mayor
    console.print("[cyan]User:[/cyan] Create a FizzBuzz program in Python")

    mayor_graph = create_mayor_graph()
    mayor_state = {
        "messages": [{"role": "user", "content": "Create a FizzBuzz program in Python"}],
        "requirements": None,
        "work_item": None,
        "spec_approved": False,
    }

    # Step 2: Designer elaborates requirements
    console.print("\n[green]Mayor:[/green] Let me get Designer to elaborate on that...\n")
    mayor_state = await mayor_graph.ainvoke(mayor_state)

    # Show requirements
    console.print("[bold]Designer's Requirements:[/bold]")
    console.print(mayor_state["requirements"])

    # Step 3: User approves
    console.print("\n[cyan]User:[/cyan] Looks good!")
    mayor_state["messages"].append({"role": "user", "content": "yes"})
    mayor_state = await mayor_graph.ainvoke(mayor_state)

    # Step 4: Inspector defines criteria
    console.print("\n[yellow]Inspector:[/yellow] Strategist is analyzing requirements...\n")

    inspector_graph = create_inspector_graph()
    inspector_state = {
        "messages": [],
        "work_item_id": mayor_state["work_item"].id,
        "test_spec": None,
        "proposed_criteria": [],
        "phase": "planning",
    }
    inspector_state = await inspector_graph.ainvoke(inspector_state)

    # Show criteria
    console.print("[bold]Strategist's Criteria:[/bold]")
    for c in inspector_state["proposed_criteria"]:
        console.print(f"  [{c.id}] {c.description}")

    # Step 5: User approves criteria
    console.print("\n[cyan]User:[/cyan] Approve all")

    # Finalize spec
    inspector_state["messages"].append({"role": "user", "content": "approve"})
    inspector_state = await inspector_graph.ainvoke(inspector_state)

    # Step 6: Work can now be dispatched
    console.print(f"\n[green]✅ Work {mayor_state['work_item'].id} ready for dispatch[/green]")

    # ... Polecat implements, Refinery triggers verification, etc.


if __name__ == "__main__":
    asyncio.run(main())
```

---

## Dependencies

```toml
# pyproject.toml
[project]
name = "vermas"
version = "0.1.0"
description = "Verifiable Multi-Agent System"
requires-python = ">=3.11"
dependencies = [
    "langgraph>=0.2.0",
    "langchain>=0.3.0",
    "langchain-anthropic>=0.2.0",
    "langchain-openai>=0.2.0",
    "pydantic>=2.0",
    "pydantic-settings>=2.0",
    "sqlmodel>=0.0.18",
    "typer>=0.12.0",
    "rich>=13.0",
    "aiosqlite>=0.20.0",
]

[project.optional-dependencies]
dev = [
    "pytest>=8.0",
    "pytest-asyncio>=0.23",
    "ruff>=0.4.0",
    "mypy>=1.10",
]

[project.scripts]
vermas = "vermas.cli:app"

[tool.ruff]
target-version = "py311"
line-length = 100

[tool.mypy]
python_version = "3.11"
strict = true
```

---

## Advantages of Python Implementation

1. **Native LangGraph**: State machines with parallel execution, human-in-the-loop
2. **Pydantic Validation**: Runtime type checking, clear schemas
3. **LangChain Ecosystem**: Ready-made tools, prompt templates, output parsers
4. **Async Native**: Full async/await support throughout
5. **Rich CLI**: Beautiful terminal output with Rich
6. **Fast Iteration**: No compile step, instant testing
7. **AI Tooling**: Direct access to all Python ML/AI libraries

## Trade-offs vs Go Implementation

| Aspect | Go | Python |
|--------|-----|--------|
| Performance | Faster runtime | Slower, but LLM calls dominate |
| Type Safety | Compile-time | Runtime (Pydantic) |
| Concurrency | goroutines | asyncio |
| Distribution | Single binary | Requires Python environment |
| Shell integration | Native | subprocess |
| LLM ecosystem | Limited | Extensive |

## Recommended Approach

Build the Python version as a **parallel implementation**, not a replacement:

1. **Phase 1**: Build core models and graphs in Python
2. **Phase 2**: Test with FizzBuzz and simple examples
3. **Phase 3**: Compare results with Go implementation
4. **Phase 4**: Decide on primary implementation based on results

The Python version can serve as:
- Rapid prototyping environment
- Reference implementation for verification logic
- Alternative runtime for specific use cases
