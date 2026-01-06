# VerMAS Python Quick Start

> Get the Python implementation running in 5 minutes

## Prerequisites

- Python 3.11+
- uv (recommended) or pip
- **Claude Code CLI** installed (`claude` command)
- **tmux** installed
- Claude Pro/Max subscription (no API costs!)

## Setup

```bash
# Create project
mkdir vermas-py && cd vermas-py

# Initialize with uv (recommended)
uv init
uv add pydantic pydantic-settings typer rich libtmux

# Or with pip
python -m venv .venv
source .venv/bin/activate
pip install pydantic pydantic-settings typer rich libtmux
```

## Verify Prerequisites

```bash
# Check Claude Code is installed
claude --version

# Check tmux is installed
tmux -V

# Should see something like:
# claude version 1.x.x
# tmux 3.x
```

## Minimal Working Example

```python
# main.py
"""
VerMAS CLI-based demo.

Uses Claude Code CLI and tmux - no API costs!
"""
import asyncio
import subprocess
import os
from pathlib import Path

import libtmux
from pydantic import BaseModel
from rich.console import Console

console = Console()


class TmuxManager:
    """Simple tmux manager."""

    def __init__(self):
        self.server = libtmux.Server()

    def create_session(self, name: str, cwd: str, command: str):
        """Create a tmux session running a command."""
        session = self.server.new_session(
            session_name=name,
            start_directory=cwd,
            attach=False,
        )
        session.active_window.active_pane.send_keys(command)
        return session

    def list_sessions(self):
        """List all sessions."""
        return [s.name for s in self.server.sessions]


async def run_claude_prompt(prompt: str, profile: str = None) -> str:
    """
    Run a prompt through Claude Code CLI.

    Uses: claude -p "prompt" --output-format json
    No API costs - uses your Claude subscription!
    """
    cmd = ["claude", "-p", prompt]
    if profile:
        cmd.extend(["--profile", profile])

    process = await asyncio.create_subprocess_exec(
        *cmd,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
    )

    stdout, stderr = await process.communicate()
    return stdout.decode()


async def run_shell_test(command: str) -> tuple[bool, str]:
    """
    Run a shell command and return pass/fail.

    This is the Verifier - no LLM, just shell.
    """
    process = await asyncio.create_subprocess_shell(
        command,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.STDOUT,
    )

    stdout, _ = await process.communicate()
    passed = process.returncode == 0
    return passed, stdout.decode()


async def main():
    """Run minimal VerMAS demo."""

    console.print("[bold]VerMAS CLI Demo[/bold]")
    console.print("Using Claude Code CLI - no API costs!\n")

    # Step 1: Use Claude CLI to elaborate requirements
    console.print("[cyan]Step 1: Designer elaborates requirements...[/cyan]")

    requirements = await run_claude_prompt(
        "List 5 specific requirements for a FizzBuzz program in Python. "
        "Output as a numbered list, nothing else.",
        profile="designer"  # Optional: use a custom profile
    )

    console.print(requirements)

    # Step 2: Use Claude CLI to propose criteria
    console.print("\n[cyan]Step 2: Strategist proposes verification criteria...[/cyan]")

    criteria = await run_claude_prompt(
        f"Given these requirements:\n{requirements}\n\n"
        "Propose 5 bash commands that verify them. "
        "Each command should exit 0 for pass. Format: AC-N: description\\nCommand: ...",
    )

    console.print(criteria)

    # Step 3: Run objective tests (no LLM)
    console.print("\n[cyan]Step 3: Verifier runs objective tests...[/cyan]")

    # Example test - would normally come from criteria
    test_cmd = "python -c \"print('FizzBuzz')\" | grep -q FizzBuzz"
    passed, output = await run_shell_test(test_cmd)

    status = "[green]PASS[/green]" if passed else "[red]FAIL[/red]"
    console.print(f"  Test: FizzBuzz output check - {status}")

    # Step 4: Demonstrate tmux session creation
    console.print("\n[cyan]Step 4: Create Mayor + Inspector tmux layout...[/cyan]")

    tmux = TmuxManager()

    # Check if we should actually create sessions
    if "--create-sessions" in os.sys.argv:
        session_name = f"vermas-demo-{os.getpid()}"
        tmux.create_session(
            session_name,
            str(Path.cwd()),
            "claude --profile mayor"
        )
        console.print(f"  Created session: {session_name}")
        console.print(f"  Attach with: tmux attach -t {session_name}")
    else:
        console.print("  (Skipped - run with --create-sessions to create tmux session)")

    console.print("\n[green]✅ Demo complete![/green]")


if __name__ == "__main__":
    asyncio.run(main())
```

## Run It

```bash
# No API key needed! Uses Claude Code CLI

# Run the demo
python main.py

# Or with tmux session creation
python main.py --create-sessions
```

## Expected Output

```
VerMAS CLI Demo
Using Claude Code CLI - no API costs!

Step 1: Designer elaborates requirements...
1. Print numbers 1 through 100
2. For multiples of 3, print "Fizz"
3. For multiples of 5, print "Buzz"
4. For multiples of both, print "FizzBuzz"
5. Output one item per line

Step 2: Strategist proposes verification criteria...
AC-1: Output has 100 lines
Command: python fizzbuzz.py | wc -l | grep -q 100

AC-2: Line 15 is FizzBuzz
Command: python fizzbuzz.py | sed -n '15p' | grep -q FizzBuzz
...

Step 3: Verifier runs objective tests...
  Test: FizzBuzz output check - PASS

Step 4: Create Mayor + Inspector tmux layout...
  (Skipped - run with --create-sessions to create tmux session)

✅ Demo complete!
```

## Full Two-Pane Startup

```bash
# Start VerMAS with Mayor and Inspector
python -c "
import os
import libtmux

server = libtmux.Server()
session = server.new_session('vermas', attach=False)

# Left pane: Mayor
session.active_window.active_pane.send_keys('claude --profile mayor')

# Right pane: Inspector
session.active_window.split(vertical=True)
session.active_window.panes[1].send_keys('claude --profile inspector')

print(f'Started! Attach with: tmux attach -t vermas')
"

# Attach to it
tmux attach -t vermas
```

## Next Steps

1. Add Pydantic models for WorkItem, TestSpec, etc.
2. Implement full Inspector workflow (Advocate/Critic/Judge)
3. Add beads integration (JSONL files)
4. Build the CLI with Typer

## Project Structure

```
vermas-py/
├── vermas/
│   ├── models/       # Pydantic models
│   ├── tmux/         # Tmux session management
│   ├── claude/       # Claude CLI wrapper
│   ├── agents/       # Agent implementations
│   ├── db/           # Beads (JSONL) persistence
│   └── cli.py        # Typer CLI
├── tests/
├── main.py
└── pyproject.toml
```

See `VERMAS-PYTHON-CLI.md` for the complete CLI-based architecture.
