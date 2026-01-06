# VerMAS Workflows

> Molecule state machine and workflow patterns

## Molecule System (MEOW)

MEOW = Molecule states for workflow management

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          MOLECULE STATE MACHINE                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│                         ┌─────────────┐                                     │
│                         │   ICE-9     │                                     │
│                         │  (Formula)  │                                     │
│                         │             │                                     │
│                         │ .toml file  │                                     │
│                         │ Template    │                                     │
│                         └──────┬──────┘                                     │
│                                │                                            │
│                                │ cook                                       │
│                                ▼                                            │
│                         ┌─────────────┐                                     │
│                         │   SOLID     │                                     │
│                         │  (Proto)    │                                     │
│                         │             │                                     │
│                         │ Compiled    │                                     │
│                         │ Ready       │                                     │
│                         └──────┬──────┘                                     │
│                                │                                            │
│                 ┌──────────────┼──────────────┐                             │
│                 │ pour                        │ wisp                        │
│                 ▼                             ▼                              │
│          ┌─────────────┐              ┌─────────────┐                       │
│          │   LIQUID    │              │   VAPOR     │                       │
│          │   (Mol)     │              │   (Wisp)    │                       │
│          │             │              │             │                       │
│          │ Persistent  │              │ Ephemeral   │                       │
│          │ Tracked     │              │ Auto-expire │                       │
│          └──────┬──────┘              └──────┬──────┘                       │
│                 │                            │                              │
│      ┌──────────┴──────────┐                 │                              │
│      │ squash              │ burn            │ (evaporate)                  │
│      ▼                     ▼                 ▼                              │
│ ┌─────────┐          ┌─────────┐       ┌─────────┐                         │
│ │ ARCHIVE │          │ DISCARD │       │  GONE   │                         │
│ │ Record  │          │ No trace│       │         │                         │
│ └─────────┘          └─────────┘       └─────────┘                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## State Definitions

| State | Name | Storage | Persistence | Use Case |
|-------|------|---------|-------------|----------|
| **Ice-9** | Formula | `.beads/formulas/*.toml` | Permanent | Template definitions |
| **Solid** | Protomolecule | In memory | Session | Compiled, ready to use |
| **Liquid** | Molecule | `.beads/mols/*.json` | Persistent | Active work tracking |
| **Vapor** | Wisp | In memory | Ephemeral | Patrol loops, temp tasks |
| **Archive** | Record | `.beads/mols/*.archive.json` | Permanent | Completed workflows |

---

## Operators

| Operator | From | To | Description |
|----------|------|----|-------------|
| **cook** | Ice-9 | Solid | Compile formula into protomolecule |
| **pour** | Solid | Liquid | Create persistent workflow attached to bead |
| **wisp** | Solid | Vapor | Create ephemeral workflow (auto-expires) |
| **squash** | Liquid/Vapor | Archive | Complete with summary, create record |
| **burn** | Liquid/Vapor | (gone) | Discard without record |

---

## Formula Structure (TOML)

Formulas are workflow templates in TOML format:

```
.beads/formulas/
├── mol-witness-patrol.formula.toml
├── mol-refinery-merge.formula.toml
├── mol-polecat-work.formula.toml
└── mol-inspector-verify.formula.toml
```

### Formula Fields

| Field | Type | Description |
|-------|------|-------------|
| `formula` | string | Unique identifier |
| `description` | string | What this workflow does |
| `version` | int | Schema version |
| `steps` | array | Ordered list of steps |

### Step Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique step identifier |
| `title` | string | Human-readable name |
| `description` | string | Detailed instructions (markdown) |
| `needs` | array | Step IDs that must complete first |

---

## Example Formulas

### Witness Patrol

```
formula = "mol-witness-patrol"
description = "Monitor polecats, nudge idle, escalate stuck"
version = 1

STEPS:
1. inbox-check
   - Check mail for POLECAT_DONE messages
   - No dependencies

2. survey-workers
   - List all active polecats
   - Check idle time for each
   - Needs: inbox-check

3. nudge-idle
   - Send nudge to polecats idle >5min
   - Needs: survey-workers

4. escalate-stuck
   - Kill polecats stuck >15min
   - Report to Deacon
   - Needs: nudge-idle

FLOW:
inbox-check → survey-workers → nudge-idle → escalate-stuck
```

### Polecat Work

```
formula = "mol-polecat-work"
description = "Execute assigned bead to completion"
version = 1

STEPS:
1. understand
   - Read hook and bead details
   - Plan approach
   - No dependencies

2. implement
   - Write code
   - Follow conventions
   - Needs: understand

3. test
   - Run tests
   - Verify requirements
   - Needs: implement

4. complete
   - Commit and push
   - Signal done
   - Needs: test

FLOW:
understand → implement → test → complete
```

### Inspector Verification

```
formula = "mol-inspector-verify"
description = "VerMAS verification workflow"
version = 1

STEPS:
1. elaborate
   - Designer creates spec
   - No dependencies

2. strategize
   - Strategist creates tests
   - Needs: elaborate

3. verify
   - Verifier runs tests (no LLM)
   - Needs: strategize

4. audit
   - Auditor reviews evidence
   - Needs: verify

5. advocate
   - Argue for PASS
   - Needs: audit

6. criticize
   - Argue for FAIL
   - Needs: audit

7. judge
   - Render verdict
   - Needs: advocate, criticize

FLOW:
                    ┌─→ advocate ─┐
elaborate → strategize → verify → audit ─┤              ├─→ judge
                    └─→ criticize─┘
```

---

## Step Execution Model

### Step Status

| Status | Meaning |
|--------|---------|
| `pending` | Not started, may have unmet dependencies |
| `blocked` | Dependencies not yet complete |
| `ready` | Dependencies met, can execute |
| `in_progress` | Currently executing |
| `completed` | Finished successfully |
| `failed` | Execution failed |
| `skipped` | Explicitly skipped |

### Execution Rules

1. **Ready check:** Step is ready when all `needs` are `completed`
2. **One at a time:** Only one step `in_progress` at a time (per molecule)
3. **No backtrack:** Completed steps don't re-run
4. **Fail fast:** Failed step blocks downstream steps

### Finding Ready Steps

```
ready_steps = []
for each step:
    if step.status == pending:
        if all(dep.status == completed for dep in step.needs):
            ready_steps.append(step)
```

---

## Molecule Lifecycle

### Creating a Molecule

```
1. Load formula from .toml file
2. Compile to protomolecule (cook)
3. Attach to bead (pour)
4. Save to .beads/mols/{id}.json
```

### Executing a Molecule

```
1. Load molecule from .json
2. Find ready steps
3. For each ready step:
   a. Mark in_progress
   b. Execute step instructions
   c. Mark completed (or failed)
4. Repeat until all done or stuck
```

### Completing a Molecule

```
Option A - Squash (with record):
1. Generate summary
2. Write to .archive.json
3. Delete active .json

Option B - Burn (no record):
1. Delete active .json
2. No archive created
```

---

## Wisp vs Molecule

| Aspect | Molecule (Liquid) | Wisp (Vapor) |
|--------|-------------------|--------------|
| Persistence | Saved to disk | Memory only |
| Lifetime | Until squash/burn | Until TTL expires |
| Recovery | Survives crashes | Lost on crash |
| Use case | Important work | Patrol loops |
| Tracking | Full audit trail | No record |

**When to use Molecule:**
- Bead execution
- Verification workflows
- Anything that needs audit trail

**When to use Wisp:**
- Witness patrol loops
- Deacon health checks
- Temporary coordination

---

## Parallel Execution

Some workflows support parallel steps:

```
Inspector Verification:
                    ┌─→ advocate ─┐
... → audit ───────┤              ├─→ judge
                    └─→ criticize─┘

Advocate and Critic run in parallel after Audit completes.
Judge waits for both.
```

### Implementing Parallelism

1. Multiple steps can be `ready` simultaneously
2. Launch each ready step in its own context
3. Parent waits for all children
4. Next step starts when all dependencies complete

---

## Workflow Patterns

### Sequential Pipeline

```
A → B → C → D

Each step depends on the previous.
Simple, predictable, easy to debug.
```

### Fork-Join

```
    ┌─→ B ─┐
A ──┤      ├──→ D
    └─→ C ─┘

A completes, B and C run in parallel, D waits for both.
Used for: parallel testing, adversarial review.
```

### Diamond

```
    ┌─→ B ─┐
A ──┤      ├──→ D
    └─→ C ─┘
        │
        ▼
        E

D needs B and C.
E only needs C.
```

### Conditional (not yet implemented)

```
A → [condition] → B (if true)
              └─→ C (if false)

Would require step predicates.
Future enhancement.
```

---

## Integration with Agents

### Witness Patrol

```
Witness starts
    │
    ▼
wisp(mol-witness-patrol)
    │
    ├─→ Execute steps in loop
    │
    └─→ When complete, wisp again (infinite patrol)
```

### Polecat Work

```
Polecat spawns with bead
    │
    ▼
pour(mol-polecat-work, bead_id)
    │
    ├─→ Execute steps
    │
    ├─→ On complete: squash + signal done
    │
    └─→ On fail: leave for Witness
```

### Refinery Merge

```
Refinery receives MERGE_READY
    │
    ▼
pour(mol-inspector-verify, bead_id)  [if VerMAS enabled]
    │
    ├─→ Execute verification
    │
    ├─→ On PASS: merge + squash
    │
    └─→ On FAIL: REWORK_REQUEST + burn
```

---

## Debugging Workflows

### Inspecting Active Molecules

```
bd mol list              # List all active molecules
bd mol show {id}         # Show molecule details
bd mol steps {id}        # Show step status
```

### Common Issues

| Symptom | Cause | Fix |
|---------|-------|-----|
| Step stuck in `pending` | Dependency not complete | Check dependency status |
| Step stuck in `in_progress` | Execution hanging | Check agent session |
| No ready steps | All blocked | Check for circular deps |
| Molecule not completing | Failed step | Review step output |

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [AGENTS.md](./AGENTS.md) - Agent roles
- [MESSAGING.md](./MESSAGING.md) - Communication patterns
- [EVALUATION.md](./EVALUATION.md) - How to evaluate the system
