# VerMAS Design Order: The Startup Growth Model

> Build like a startup: start with 2-4 people getting work done, add complexity only when needed

---

## Philosophy: Bottom-Up Growth

**Don't design a corporation when you're a garage startup.**

The previous iterations explored top-down enterprise architecture:
- 8 layers of abstraction
- Event sourcing with projections
- Complex governance hierarchies
- 35+ files before anything works

**The problem:** This is how you build for scale you don't have yet. Real startups don't start with HR departments and compliance workflows. They start with people doing work.

**The new approach:** Start with the absolute minimum to get agents working together. Add structure only when pain emerges.

---

## Stage 1: The Garage (2-4 Agents)

**What you have:**
- 2-4 Claude instances (agents)
- A shared git repo
- The ability to talk to each other (files, mail)

**What you DON'T need yet:**
- Event sourcing
- Projections
- Workflows/templates
- Verification pipelines
- Complex governance
- Departments

**The goal:** Get work done. Period.

```
THE GARAGE SETUP

┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│   Human                                                                     │
│     │                                                                       │
│     │  "Here's what I need done"                                            │
│     │                                                                       │
│     ▼                                                                       │
│   ┌─────────────┐                                                           │
│   │   Agent A   │◄────────────────────────────────┐                         │
│   │  (Lead)     │                                 │                         │
│   └──────┬──────┘                                 │                         │
│          │                                        │                         │
│          │  "You do X, you do Y"                  │ "Done" / "Stuck"        │
│          │                                        │                         │
│          ▼                                        │                         │
│   ┌─────────────┐    ┌─────────────┐             │                         │
│   │   Agent B   │    │   Agent C   │─────────────┘                         │
│   │  (Worker)   │    │  (Worker)   │                                        │
│   └─────────────┘    └─────────────┘                                        │
│                                                                             │
│   Communication: .beads/ for tracking, files/mail for coordination          │
│   No layers. No projections. Just work.                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### What Stage 1 Actually Needs

| Need | Solution | Complexity |
|------|----------|------------|
| Track work | Beads (already have it) | Zero - exists |
| Coordinate | Mail system (gt mail) | Zero - exists |
| Know who does what | Simple CLAUDE.md roles | Minimal |
| Share code | Git worktrees | Minimal |

**Files needed:** Maybe 5-10, not 35.

### Stage 1 Anti-Patterns

DON'T do this in Stage 1:
- ❌ Build elaborate event sourcing
- ❌ Create verification pipelines
- ❌ Design governance hierarchies
- ❌ Implement complex workflows
- ❌ Add compliance rules

DO this instead:
- ✅ Get 2-4 agents actually completing tasks
- ✅ Use beads to track what needs doing
- ✅ Use mail to coordinate
- ✅ Fix problems as they emerge

---

## Stage 2: The Growing Startup (5-10 Agents)

**Trigger:** Stage 1 is working, but you're hitting problems:
- Work getting dropped
- Confusion about who's doing what
- Quality issues sneaking through
- Coordination overhead increasing

**What you ADD (only what's needed):**
- Clearer role definitions
- Basic QA step (human or simple automated)
- Maybe one "supervisor" agent

```
GROWING STARTUP

┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│   Human                                                                     │
│     │                                                                       │
│     ▼                                                                       │
│   ┌─────────────┐                                                           │
│   │   Lead      │  ← Makes decisions, breaks down work                     │
│   │   Agent     │                                                           │
│   └──────┬──────┘                                                           │
│          │                                                                  │
│          ├──────────────┬──────────────┬──────────────┐                     │
│          ▼              ▼              ▼              ▼                     │
│   ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐               │
│   │  Worker   │  │  Worker   │  │  Worker   │  │  Worker   │               │
│   │    A      │  │    B      │  │    C      │  │    D      │               │
│   └───────────┘  └───────────┘  └───────────┘  └───────────┘               │
│          │              │              │              │                     │
│          └──────────────┴──────────────┴──────────────┘                     │
│                                    │                                        │
│                                    ▼                                        │
│                            ┌───────────┐                                    │
│                            │   QA      │  ← Added when quality matters     │
│                            │  (light)  │                                    │
│                            └───────────┘                                    │
│                                                                             │
│   NEW: Basic roles, light QA                                                │
│   STILL NO: Event sourcing, complex governance                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### What Stage 2 Adds

| Problem Emerged | Solution Added |
|-----------------|----------------|
| Work getting lost | More structured beads usage |
| "Who's doing this?" | Clearer assignments |
| Quality slipping | Light QA review step |
| Lead overloaded | Maybe delegate some decisions |

---

## Stage 3: The Scaling Company (10-20+ Agents)

**Trigger:** Stage 2 is working, but:
- Coordination is breaking down
- Need parallel workstreams
- Can't rely on one Lead for everything
- Quality needs to be more systematic

**What you ADD:**
- Multiple teams/workstreams
- Supervisors per team
- More formal workflows
- Maybe verification for critical work

```
SCALING COMPANY

┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│   Human                                                                     │
│     │                                                                       │
│     ▼                                                                       │
│   ┌─────────────┐                                                           │
│   │   CEO       │  ← High-level coordination only                          │
│   │   Agent     │                                                           │
│   └──────┬──────┘                                                           │
│          │                                                                  │
│          ├───────────────────────────────────────┐                          │
│          ▼                                       ▼                          │
│   ┌─────────────────────┐               ┌─────────────────────┐             │
│   │    Team Alpha       │               │    Team Beta        │             │
│   │  ┌───────────────┐  │               │  ┌───────────────┐  │             │
│   │  │  Supervisor   │  │               │  │  Supervisor   │  │             │
│   │  └───────┬───────┘  │               │  └───────┬───────┘  │             │
│   │          │          │               │          │          │             │
│   │    ┌─────┴─────┐    │               │    ┌─────┴─────┐    │             │
│   │    ▼     ▼     ▼    │               │    ▼     ▼     ▼    │             │
│   │   W1    W2    W3    │               │   W4    W5    W6    │             │
│   │                     │               │                     │             │
│   │  └──────┬──────┘    │               │  └──────┬──────┘    │             │
│   │         ▼           │               │         ▼           │             │
│   │       QA-A          │               │       QA-B          │             │
│   └─────────────────────┘               └─────────────────────┘             │
│                                                                             │
│   NOW you might need: Workflows, verification, better governance            │
│   Because: Scale demands it, not because you designed for it               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Implementation: Start at Stage 1

### Minimum Viable Multi-Agent System

**To get Stage 1 working, you need:**

1. **Beads** - Already exists (Steve Yegge's bd)
   - Track work items
   - Dependencies handled by graph
   - Priority emerges from blocked/blocking

2. **Mail** - Already exists (gt mail)
   - Agents communicate
   - Handoffs work

3. **Roles** - Just CLAUDE.md files
   - Lead agent: "You coordinate work"
   - Worker agents: "You do the work assigned"

4. **Git worktrees** - Standard git feature
   - Each agent works in isolation
   - Merge when done

**That's it.** No event sourcing. No 8 layers. No 35 files.

### The First Working System

```
CONCRETE STAGE 1

Directory structure:
project/
├── .beads/              ← Work tracking (existing)
├── .claude/
│   └── settings.json    ← Agent configurations
├── crew/
│   ├── lead/            ← Lead agent worktree
│   │   └── CLAUDE.md    ← "You are the lead..."
│   ├── worker-1/        ← Worker worktree
│   │   └── CLAUDE.md    ← "You are a worker..."
│   └── worker-2/        ← Worker worktree
│       └── CLAUDE.md    ← "You are a worker..."
└── src/                 ← Actual code

Communication:
- bd create/show/update  ← Track work
- gt mail send/inbox     ← Coordinate
- git push/pull          ← Share code

That's the whole system.
```

### What We're NOT Building Initially

These are all deferred until pain demands them:

| Feature | Stage Added | Why Wait |
|---------|-------------|----------|
| Event sourcing | 3+ | Overkill for small teams |
| Projections | 3+ | Simple file reads work fine |
| Verification pipeline | 2-3 | Manual QA works first |
| Complex governance | 3+ | Adds overhead without value |
| Workflow templates | 2-3 | Ad-hoc coordination works |
| 8 layers of abstraction | Never? | Maybe never needed |

---

## Growth Triggers

**Move from Stage 1 to Stage 2 when:**
- [ ] Work is getting dropped/forgotten
- [ ] Quality issues are hurting output
- [ ] Lead agent is overwhelmed
- [ ] Coordination takes more time than work

**Move from Stage 2 to Stage 3 when:**
- [ ] Single lead can't coordinate everything
- [ ] Need parallel independent workstreams
- [ ] Manual QA is bottleneck
- [ ] Want to scale beyond 10 agents

**Add a feature when:**
- [ ] The pain of NOT having it exceeds the cost of building it
- [ ] NOT because "good architecture" says you need it

---

## Comparison: Top-Down vs Bottom-Up

| Aspect | Top-Down (Old) | Bottom-Up (New) |
|--------|----------------|-----------------|
| Start with | 8-layer architecture | 2-4 agents working |
| First milestone | "Event system complete" | "Work got done" |
| Files to build | 35+ | 5-10 |
| When to add complexity | "Before we scale" | "When pain emerges" |
| Governance | Design upfront | Grow organically |
| Verification | Build first | Add when quality matters |
| Risk | Over-engineering | Under-engineering |
| Recovery | Hard (sunk cost) | Easy (add as needed) |

---

## Approved Design Decisions

These principles are approved for the Python implementation:

1. **Start at Stage 1** - Get 2-4 agents completing real work
2. **Use existing tools** - Beads (bd), Mail (gt mail), Git worktrees
3. **Add complexity only when pain demands it**
4. **Priority comes from beads dependency graph** - Not a separate scoring system
5. **Grow like a startup** - Don't architect for scale you don't have

---

## Next Steps

1. **Define Stage 1 roles** - Lead + 2-3 Workers
2. **Create CLAUDE.md files** - Minimal role definitions
3. **Test with real work** - Give the system actual tasks
4. **Document pain points** - What's not working?
5. **Add solutions to pain** - Grow the system organically

**NOT next steps:**
- ~~Build 8 layers of abstraction~~
- ~~Implement event sourcing~~
- ~~Create verification pipeline~~
- ~~Design complex governance~~

---

## Archived: Top-Down Design Iterations

The previous 5 iterations of top-down design are preserved in [DESIGN_ORDER_ARCHIVE.md](./DESIGN_ORDER_ARCHIVE.md) for reference. They contain useful thinking about:
- Layer dependencies
- Event sourcing patterns
- Testing strategies
- File organization

These may be valuable when/if we grow to Stage 3 and need more sophisticated infrastructure.
