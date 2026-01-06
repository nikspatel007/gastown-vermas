# VerMAS Python Messaging

> Mail protocol, hooks, and agent communication

## Mail Protocol Overview

Gas Town agents communicate via an internal mail system stored in beads.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           MAIL FLOW                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Polecat ──[POLECAT_DONE]──▶ Witness ──[MERGE_READY]──▶ Refinery          │
│      │                           │                           │              │
│      │                           │                           │              │
│      ▼                           ▼                           ▼              │
│   [HELP] ───────────────▶ Escalation ◀────────────── [REWORK_REQUEST]      │
│                                │                                            │
│                                ▼                                            │
│                             Mayor                                           │
│                                │                                            │
│                                ▼                                            │
│   [HANDOFF] ◀──────────── Session End                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Message Types

| Type | From | To | Purpose |
|------|------|-----|---------|
| `POLECAT_DONE` | Polecat | Witness | Work completed, ready for review |
| `MERGE_READY` | Witness | Refinery | Forward completed work for merge |
| `MERGED` | Refinery | Author | Notify successful merge |
| `REWORK_REQUEST` | Refinery | Author | Changes needed before merge |
| `WITNESS_PING` | Deacon | Witness | Health check |
| `HELP` | Any | Witness/Mayor | Request assistance |
| `HANDOFF` | Any | Self/Next | Session continuity (🤝 prefix) |
| `NUDGE` | Witness | Polecat | Wake up idle worker |
| `ESCALATION` | Witness | Deacon/Mayor | Problem beyond local handling |

---

## Pydantic Models

```python
# vermas/models/mail.py
from enum import Enum
from typing import Optional, Dict, Any
from datetime import datetime
from pydantic import BaseModel, Field


class MessageType(str, Enum):
    """Known message types."""
    POLECAT_DONE = "polecat_done"
    MERGE_READY = "merge_ready"
    MERGED = "merged"
    REWORK_REQUEST = "rework_request"
    WITNESS_PING = "witness_ping"
    HELP = "help"
    HANDOFF = "handoff"
    NUDGE = "nudge"
    ESCALATION = "escalation"
    GENERIC = "generic"


class MessagePriority(str, Enum):
    """Message priority levels."""
    URGENT = "urgent"    # Requires immediate attention
    NORMAL = "normal"    # Process in order
    LOW = "low"          # Can wait


class Message(BaseModel):
    """
    Mail message between agents.

    Stored in .beads/messages.jsonl
    """
    id: str
    from_addr: str           # BD_ACTOR format: rig/role/name
    to_addr: str             # BD_ACTOR format
    subject: str
    body: str
    message_type: MessageType = MessageType.GENERIC
    priority: MessagePriority = MessagePriority.NORMAL
    created_at: datetime = Field(default_factory=datetime.utcnow)
    read_at: Optional[datetime] = None
    metadata: Dict[str, Any] = Field(default_factory=dict)

    @property
    def is_read(self) -> bool:
        return self.read_at is not None

    @property
    def is_handoff(self) -> bool:
        return self.subject.startswith("🤝 HANDOFF:")


class Attachment(BaseModel):
    """Message attachment (bead reference)."""
    bead_id: str
    attachment_type: str  # "hook", "reference", "context"
```

---

## Mail Store

```python
# vermas/mail/store.py
from pathlib import Path
from typing import List, Optional, Iterator
from datetime import datetime
import hashlib
import json

from vermas.models.mail import Message, MessageType, MessagePriority


class MailStore:
    """
    JSONL-based mail storage.

    Messages stored in .beads/messages.jsonl
    """

    def __init__(self, beads_dir: Path):
        self.beads_dir = beads_dir
        self.messages_file = beads_dir / "messages.jsonl"
        self._ensure_file()

    def _ensure_file(self):
        self.beads_dir.mkdir(parents=True, exist_ok=True)
        if not self.messages_file.exists():
            self.messages_file.touch()

    def _generate_id(self) -> str:
        timestamp = datetime.utcnow().isoformat()
        return f"msg-{hashlib.sha256(timestamp.encode()).hexdigest()[:8]}"

    def send(
        self,
        from_addr: str,
        to_addr: str,
        subject: str,
        body: str,
        message_type: MessageType = MessageType.GENERIC,
        priority: MessagePriority = MessagePriority.NORMAL,
        metadata: dict = None,
    ) -> Message:
        """Send a message."""
        msg = Message(
            id=self._generate_id(),
            from_addr=from_addr,
            to_addr=to_addr,
            subject=subject,
            body=body,
            message_type=message_type,
            priority=priority,
            metadata=metadata or {},
        )

        with open(self.messages_file, "a") as f:
            f.write(msg.model_dump_json() + "\n")

        return msg

    def inbox(self, addr: str, unread_only: bool = False) -> List[Message]:
        """Get messages for an address."""
        messages = []
        for msg in self._iter_messages():
            if msg.to_addr == addr or msg.to_addr.startswith(addr):
                if unread_only and msg.is_read:
                    continue
                messages.append(msg)

        # Sort by priority (urgent first), then by date
        priority_order = {"urgent": 0, "normal": 1, "low": 2}
        messages.sort(key=lambda m: (priority_order[m.priority], m.created_at))
        return messages

    def outbox(self, addr: str) -> List[Message]:
        """Get sent messages from an address."""
        return [
            msg for msg in self._iter_messages()
            if msg.from_addr == addr or msg.from_addr.startswith(addr)
        ]

    def get(self, msg_id: str) -> Optional[Message]:
        """Get message by ID."""
        for msg in self._iter_messages():
            if msg.id == msg_id:
                return msg
        return None

    def mark_read(self, msg_id: str) -> Optional[Message]:
        """Mark a message as read."""
        messages = list(self._iter_messages())
        for msg in messages:
            if msg.id == msg_id:
                msg.read_at = datetime.utcnow()
                self._write_all(messages)
                return msg
        return None

    def _iter_messages(self) -> Iterator[Message]:
        """Iterate over all messages."""
        if not self.messages_file.exists():
            return

        with open(self.messages_file) as f:
            for line in f:
                line = line.strip()
                if line:
                    yield Message.model_validate_json(line)

    def _write_all(self, messages: List[Message]):
        """Rewrite all messages."""
        with open(self.messages_file, "w") as f:
            for msg in messages:
                f.write(msg.model_dump_json() + "\n")
```

---

## Mail Protocol Handlers

```python
# vermas/mail/protocol.py
from typing import Callable, Awaitable, Dict
from vermas.models.mail import Message, MessageType


class MailProtocol:
    """
    Protocol handlers for different message types.

    Each handler processes a specific message type.
    """

    def __init__(self):
        self._handlers: Dict[MessageType, Callable] = {}

    def register(self, msg_type: MessageType):
        """Decorator to register a message handler."""
        def decorator(func: Callable[[Message], Awaitable[None]]):
            self._handlers[msg_type] = func
            return func
        return decorator

    async def handle(self, msg: Message):
        """Handle a message based on its type."""
        handler = self._handlers.get(msg.message_type)
        if handler:
            await handler(msg)


# Create protocol instance
protocol = MailProtocol()


# ─────────────────────────────────────────────────────────────
# Standard Protocol Handlers
# ─────────────────────────────────────────────────────────────

@protocol.register(MessageType.POLECAT_DONE)
async def handle_polecat_done(msg: Message):
    """
    Handle POLECAT_DONE message.

    Sent by: Polecat
    Received by: Witness
    Action: Forward to Refinery for merge processing
    """
    from vermas.mail.store import MailStore
    from pathlib import Path

    store = MailStore(Path(".beads"))

    # Extract bead ID from metadata
    bead_id = msg.metadata.get("bead_id")
    slot = msg.metadata.get("slot")

    # Forward to Refinery
    store.send(
        from_addr=msg.to_addr,  # Witness
        to_addr=msg.to_addr.replace("/witness", "/refinery"),
        subject=f"MERGE_READY: {bead_id}",
        body=f"Polecat {slot} completed work on {bead_id}. Ready for merge.\n\nOriginal message:\n{msg.body}",
        message_type=MessageType.MERGE_READY,
        metadata={"bead_id": bead_id, "slot": slot},
    )


@protocol.register(MessageType.MERGE_READY)
async def handle_merge_ready(msg: Message):
    """
    Handle MERGE_READY message.

    Sent by: Witness
    Received by: Refinery
    Action: Add to merge queue
    """
    # Refinery will pick this up in its patrol loop
    pass


@protocol.register(MessageType.REWORK_REQUEST)
async def handle_rework_request(msg: Message):
    """
    Handle REWORK_REQUEST message.

    Sent by: Refinery
    Received by: Polecat/Author
    Action: Re-open bead, notify author
    """
    from vermas.beads.store import BeadStore
    from vermas.models.bead import BeadStatus
    from pathlib import Path

    bead_id = msg.metadata.get("bead_id")
    if bead_id:
        store = BeadStore(Path(".beads"))
        bead = store.get(bead_id)
        if bead:
            bead.status = BeadStatus.OPEN
            store.update(bead)


@protocol.register(MessageType.HANDOFF)
async def handle_handoff(msg: Message):
    """
    Handle HANDOFF message.

    Sent by: Any agent ending session
    Received by: Self (next session) or designated successor
    Action: Hook the handoff context for next session
    """
    from vermas.core.hooks import Hook
    from vermas.beads.store import BeadStore
    from pathlib import Path

    store = BeadStore(Path(".beads"))
    hook = Hook(msg.to_addr, store)

    # Handoff messages can be hooked for the next session
    # The message body contains context for continuation


@protocol.register(MessageType.HELP)
async def handle_help(msg: Message):
    """
    Handle HELP message.

    Sent by: Any struggling agent
    Received by: Witness or Mayor
    Action: Assess and escalate if needed
    """
    # Log help request for monitoring
    pass


@protocol.register(MessageType.ESCALATION)
async def handle_escalation(msg: Message):
    """
    Handle ESCALATION message.

    Sent by: Witness (usually)
    Received by: Deacon or Mayor
    Action: Higher-level intervention
    """
    # Escalations require human or Mayor attention
    pass
```

---

## Mail CLI Commands

```python
# vermas/cli.py (mail commands)
import typer
from pathlib import Path
from rich.console import Console
from rich.table import Table

mail_app = typer.Typer()
console = Console()


@mail_app.command("inbox")
def inbox(unread: bool = False):
    """Check your inbox."""
    from vermas.mail.store import MailStore
    import os

    actor = os.environ.get("BD_ACTOR", "unknown")
    store = MailStore(Path(".beads"))
    messages = store.inbox(actor, unread_only=unread)

    if not messages:
        console.print("[dim]No messages[/dim]")
        return

    table = Table(title=f"Inbox for {actor}")
    table.add_column("ID", style="cyan")
    table.add_column("From")
    table.add_column("Subject")
    table.add_column("Date")
    table.add_column("Read", justify="center")

    for msg in messages:
        read_mark = "✓" if msg.is_read else ""
        table.add_row(
            msg.id,
            msg.from_addr,
            msg.subject[:40],
            msg.created_at.strftime("%Y-%m-%d %H:%M"),
            read_mark,
        )

    console.print(table)


@mail_app.command("read")
def read(msg_id: str):
    """Read a specific message."""
    from vermas.mail.store import MailStore

    store = MailStore(Path(".beads"))
    msg = store.get(msg_id)

    if not msg:
        console.print(f"[red]Message not found: {msg_id}[/red]")
        raise typer.Exit(1)

    store.mark_read(msg_id)

    console.print(f"[bold]From:[/bold] {msg.from_addr}")
    console.print(f"[bold]To:[/bold] {msg.to_addr}")
    console.print(f"[bold]Subject:[/bold] {msg.subject}")
    console.print(f"[bold]Date:[/bold] {msg.created_at}")
    console.print(f"[bold]Type:[/bold] {msg.message_type}")
    console.print()
    console.print(msg.body)


@mail_app.command("send")
def send(
    to: str,
    subject: str = typer.Option(..., "-s", "--subject"),
    message: str = typer.Option(..., "-m", "--message"),
):
    """Send a message."""
    from vermas.mail.store import MailStore
    import os

    actor = os.environ.get("BD_ACTOR", "unknown")
    store = MailStore(Path(".beads"))

    msg = store.send(
        from_addr=actor,
        to_addr=to,
        subject=subject,
        body=message,
    )

    console.print(f"[green]Sent: {msg.id}[/green]")
```

---

## Hooks System

The hook is where work hangs. GUPP: If your hook has work, RUN IT.

```python
# vermas/core/hooks.py
from pathlib import Path
from typing import Optional, Union
from enum import Enum
from datetime import datetime

from vermas.models.bead import Bead
from vermas.models.mail import Message
from vermas.beads.store import BeadStore


class HookType(str, Enum):
    """What can be hooked."""
    BEAD = "bead"       # Work bead
    MAIL = "mail"       # Mail message (for handoffs)
    MOL = "mol"         # Molecule workflow


class HookContent:
    """Content of a hook."""
    def __init__(self, hook_type: HookType, ref_id: str):
        self.hook_type = hook_type
        self.ref_id = ref_id


class Hook:
    """
    Agent hook - where work hangs.

    The hook is the assignment mechanism for Gas Town.
    Work is "slung" to a hook, and when the agent starts,
    they check the hook and EXECUTE immediately (GUPP).

    Hooks persist across sessions in .beads/.hook-{agent}
    """

    def __init__(self, agent_id: str, beads_dir: Path):
        self.agent_id = agent_id
        self.beads_dir = beads_dir
        self._hook_file = beads_dir / f".hook-{agent_id.replace('/', '-')}"

    def check(self) -> Optional[HookContent]:
        """
        Check the hook for assigned work.

        Returns HookContent if work is hooked, None otherwise.
        """
        if not self._hook_file.exists():
            return None

        content = self._hook_file.read_text().strip()
        if not content:
            return None

        # Format: type:id
        try:
            hook_type, ref_id = content.split(":", 1)
            return HookContent(HookType(hook_type), ref_id)
        except ValueError:
            return None

    def attach(self, item: Union[Bead, Message, str], hook_type: HookType = None):
        """
        Attach work to the hook.

        Args:
            item: Bead, Message, or ID string
            hook_type: Required if item is a string
        """
        if isinstance(item, Bead):
            content = f"{HookType.BEAD}:{item.id}"
        elif isinstance(item, Message):
            content = f"{HookType.MAIL}:{item.id}"
        elif isinstance(item, str):
            if not hook_type:
                raise ValueError("hook_type required when attaching by ID")
            content = f"{hook_type}:{item}"
        else:
            raise TypeError(f"Cannot hook {type(item)}")

        self._hook_file.write_text(content)

    def clear(self):
        """Clear the hook."""
        if self._hook_file.exists():
            self._hook_file.unlink()

    def get_bead(self, beads: BeadStore) -> Optional[Bead]:
        """Get the hooked bead, if any."""
        content = self.check()
        if content and content.hook_type == HookType.BEAD:
            return beads.get(content.ref_id)
        return None

    def get_message(self, mail_store) -> Optional[Message]:
        """Get the hooked message, if any."""
        content = self.check()
        if content and content.hook_type == HookType.MAIL:
            return mail_store.get(content.ref_id)
        return None


class HookManager:
    """
    Manages hooks across all agents.

    Used by Mayor to sling work to agents.
    """

    def __init__(self, beads_dir: Path):
        self.beads_dir = beads_dir

    def sling(self, bead: Bead, target: str):
        """
        Sling a bead to a target agent.

        This is the primary work dispatch mechanism.
        """
        hook = Hook(target, self.beads_dir)
        hook.attach(bead)

        # Update bead status to HOOKED
        from vermas.models.bead import BeadStatus
        bead.status = BeadStatus.HOOKED
        bead.assigned_to = target

    def check_all(self) -> dict:
        """Check hooks for all known agents."""
        results = {}
        for hook_file in self.beads_dir.glob(".hook-*"):
            agent_id = hook_file.name.replace(".hook-", "").replace("-", "/")
            hook = Hook(agent_id, self.beads_dir)
            content = hook.check()
            results[agent_id] = {
                "hooked": content is not None,
                "type": content.hook_type if content else None,
                "ref": content.ref_id if content else None,
            }
        return results
```

---

## Hook CLI Commands

```python
# vermas/cli.py (hook commands)
import typer
from pathlib import Path
from rich.console import Console

hook_app = typer.Typer()
console = Console()


@hook_app.command("check")
def check():
    """Check your hook for assigned work."""
    from vermas.core.hooks import Hook
    from vermas.beads.store import BeadStore
    import os

    actor = os.environ.get("BD_ACTOR", "unknown")
    beads = BeadStore(Path(".beads"))
    hook = Hook(actor, Path(".beads"))

    content = hook.check()
    if not content:
        console.print("[dim]Hook is empty - no assigned work[/dim]")
        return

    console.print(f"[bold green]HOOKED WORK FOUND[/bold green]")
    console.print(f"Type: {content.hook_type}")
    console.print(f"ID: {content.ref_id}")

    if content.hook_type.value == "bead":
        bead = hook.get_bead(beads)
        if bead:
            console.print(f"\n[bold]Title:[/bold] {bead.title}")
            console.print(f"[bold]Priority:[/bold] {bead.priority}")
            console.print(f"[bold]Type:[/bold] {bead.issue_type}")

    console.print("\n[yellow]GUPP: Hook has work. Execute immediately![/yellow]")


@hook_app.command("attach")
def attach(bead_id: str):
    """Attach a bead to your hook."""
    from vermas.core.hooks import Hook, HookType
    import os

    actor = os.environ.get("BD_ACTOR", "unknown")
    hook = Hook(actor, Path(".beads"))
    hook.attach(bead_id, HookType.BEAD)
    console.print(f"[green]Hooked: {bead_id}[/green]")


@hook_app.command("clear")
def clear():
    """Clear your hook."""
    from vermas.core.hooks import Hook
    import os

    actor = os.environ.get("BD_ACTOR", "unknown")
    hook = Hook(actor, Path(".beads"))
    hook.clear()
    console.print("[dim]Hook cleared[/dim]")


@hook_app.command("sling")
def sling(bead_id: str, target: str):
    """Sling a bead to another agent's hook."""
    from vermas.core.hooks import HookManager
    from vermas.beads.store import BeadStore

    beads = BeadStore(Path(".beads"))
    bead = beads.get(bead_id)

    if not bead:
        console.print(f"[red]Bead not found: {bead_id}[/red]")
        raise typer.Exit(1)

    manager = HookManager(Path(".beads"))
    manager.sling(bead, target)
    beads.update(bead)

    console.print(f"[green]Slung {bead_id} to {target}[/green]")
```

---

## Handoff Protocol

For session continuity across restarts.

```python
# vermas/mail/handoff.py
from pathlib import Path
from vermas.mail.store import MailStore
from vermas.models.mail import MessageType, MessagePriority


def create_handoff(
    from_addr: str,
    to_addr: str = None,
    context: str = "",
    hook_it: bool = True,
) -> str:
    """
    Create a handoff message for session continuity.

    Args:
        from_addr: Current agent address
        to_addr: Target (default: same agent, next session)
        context: Handoff context/notes
        hook_it: Whether to hook the message for next session

    Returns:
        Message ID
    """
    store = MailStore(Path(".beads"))

    target = to_addr or from_addr
    msg = store.send(
        from_addr=from_addr,
        to_addr=target,
        subject=f"🤝 HANDOFF: Session continuity",
        body=context,
        message_type=MessageType.HANDOFF,
        priority=MessagePriority.NORMAL,
    )

    if hook_it:
        from vermas.core.hooks import Hook, HookType
        hook = Hook(target, Path(".beads"))
        hook.attach(msg.id, HookType.MAIL)

    return msg.id


# CLI command
def handoff_command(message: str = typer.Option(..., "-m", "--message")):
    """Create a handoff for the next session."""
    import os
    actor = os.environ.get("BD_ACTOR", "unknown")

    msg_id = create_handoff(
        from_addr=actor,
        context=message,
        hook_it=True,
    )

    console.print(f"[green]Handoff created: {msg_id}[/green]")
    console.print("[dim]Message hooked for next session[/dim]")
```

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [AGENTS.md](./AGENTS.md) - Agent implementations
- [WORKFLOWS.md](./WORKFLOWS.md) - Molecule workflow system
- [LIFECYCLE.md](./LIFECYCLE.md) - Polecat lifecycle
