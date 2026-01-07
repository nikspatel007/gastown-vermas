# VerMAS Priority & Extensibility Design

> Ralph Wiggum iteration 1/10: Understanding what matters NOW and bringing in external expertise

## The Questions

### Priority
1. **How do we know what's urgent?** - Signals that indicate "do this NOW"
2. **How do we distinguish urgent from important?** - Not everything urgent is important
3. **How does priority change over time?** - Decay, escalation, context shifts
4. **How do we prevent priority blindness?** - When everything is P0, nothing is

### Extensibility
5. **How do we bring in external expertise?** - Consultants, specialists, plugins
6. **What's the extension model?** - Tools, commands, sub-agents, skills
7. **How do extensions integrate?** - Discovery, installation, invocation
8. **Who learns from extensions?** - Individual, team, department, ecosystem
9. **How do we trust extensions?** - Verification, sandboxing, audit
10. **How do extensions evolve?** - Lifecycle, deprecation, replacement

---

## Iteration 1: Priority Fundamentals

### The Eisenhower Matrix Applied

```
THE PRIORITY QUADRANT

                    URGENT
                      │
         ┌────────────┼────────────┐
         │            │            │
         │   DO NOW   │   SCHEDULE │
         │            │            │
         │  Urgent +  │  Important │
         │  Important │  Not Urgent│
IMPORTANT├────────────┼────────────┤NOT IMPORTANT
         │            │            │
         │  DELEGATE  │   DROP     │
         │            │            │
         │  Urgent    │  Neither   │
         │  Not Import│            │
         │            │            │
         └────────────┼────────────┘
                      │
                 NOT URGENT
```

### What Makes Something Urgent?

Urgency comes from **time pressure**:

```
URGENCY SIGNALS

1. DEADLINE PROXIMITY
   - Deadline in < 24h → HIGH urgency
   - Deadline in < 72h → MEDIUM urgency
   - Deadline in < 1 week → LOW urgency
   - No deadline → BASE priority only

2. BLOCKING OTHERS
   - Work items blocked by this → URGENT
   - More blockers = more urgent
   - Critical path items → HIGHEST urgency

3. EXTERNAL PRESSURE
   - Customer waiting → URGENT
   - Stakeholder escalation → URGENT
   - Compliance deadline → URGENT
   - SLA breach imminent → CRITICAL

4. DECAY / STALENESS
   - Work aging in queue → urgency increases
   - Prevents indefinite deferral
   - "Oldest unaddressed" gets attention

5. EXPLICIT ESCALATION
   - Human marked as urgent
   - Supervisor escalated
   - Auto-escalation triggered
```

### What Makes Something Important?

Importance comes from **impact**:

```
IMPORTANCE SIGNALS

1. ALIGNMENT TO OBJECTIVES
   - Contributes to quarterly OKR → HIGH importance
   - On critical path to milestone → HIGH importance
   - Nice-to-have / not in plan → LOW importance

2. BUSINESS VALUE
   - Revenue impact → Quantifiable importance
   - User impact (# affected) → Scale importance
   - Strategic value → Long-term importance

3. RISK MITIGATION
   - Security vulnerability → HIGH importance
   - Data loss risk → HIGH importance
   - Compliance requirement → HIGH importance

4. TECHNICAL DEBT
   - Blocking future work → Important
   - Degrading velocity → Important
   - Pure cleanup → Low importance

5. ORGANIZATIONAL PRIORITY
   - CEO directive → HIGH importance
   - Company-wide initiative → HIGH importance
   - Department initiative → MEDIUM importance
```

### Priority Score Calculation

```python
# vermas/priority/calculator.py

from dataclasses import dataclass
from datetime import datetime, timedelta
from typing import Optional, List

@dataclass
class PriorityFactors:
    # Base priority (P0-P4)
    base_priority: int  # 0=critical, 4=backlog

    # Urgency factors
    deadline: Optional[datetime]
    blocked_items: List[str]  # IDs of items blocked by this
    external_pressure: bool
    age_hours: float
    escalated: bool

    # Importance factors
    objective_alignment: float  # 0.0 - 1.0
    business_value: float  # 0.0 - 1.0
    risk_level: float  # 0.0 - 1.0
    ceo_directive: bool

def calculate_priority_score(factors: PriorityFactors) -> float:
    """
    Calculate dynamic priority score.
    Higher score = higher priority (should be done first).
    """
    score = 0.0

    # Base priority contribution (P0=100, P1=80, P2=60, P3=40, P4=20)
    base_scores = {0: 100, 1: 80, 2: 60, 3: 40, 4: 20}
    score += base_scores.get(factors.base_priority, 20)

    # === URGENCY FACTORS ===

    # Deadline urgency (up to +50)
    if factors.deadline:
        hours_remaining = (factors.deadline - datetime.now()).total_seconds() / 3600
        if hours_remaining < 0:
            score += 50  # Overdue!
        elif hours_remaining < 24:
            score += 40
        elif hours_remaining < 72:
            score += 25
        elif hours_remaining < 168:  # 1 week
            score += 10

    # Blocking others (up to +30)
    blocker_count = len(factors.blocked_items)
    score += min(blocker_count * 10, 30)

    # External pressure (+20)
    if factors.external_pressure:
        score += 20

    # Age decay (+1 per day, max +14)
    age_days = factors.age_hours / 24
    score += min(age_days, 14)

    # Escalated (+25)
    if factors.escalated:
        score += 25

    # === IMPORTANCE FACTORS ===

    # Objective alignment (up to +30)
    score += factors.objective_alignment * 30

    # Business value (up to +20)
    score += factors.business_value * 20

    # Risk level (up to +25)
    score += factors.risk_level * 25

    # CEO directive (+40)
    if factors.ceo_directive:
        score += 40

    return score
```

### Priority Classes

Rather than just P0-P4, we have **dynamic priority classes**:

```
PRIORITY CLASSES

┌─────────────────────────────────────────────────────────────────────────────┐
│ CLASS: CRITICAL (Score > 180)                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│ Characteristics:                                                            │
│ - Drop everything, do this NOW                                              │
│ - May interrupt in-progress work                                            │
│ - Supervisor notified immediately                                           │
│ - SLA: Response within 15 minutes                                           │
│                                                                             │
│ Examples:                                                                   │
│ - Production outage                                                         │
│ - Security breach                                                           │
│ - CEO-escalated blocker                                                     │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ CLASS: HIGH (Score 140-180)                                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│ Characteristics:                                                            │
│ - Next thing to work on                                                     │
│ - Complete current task, then switch                                        │
│ - Supervisor aware                                                          │
│ - SLA: Start within 2 hours                                                 │
│                                                                             │
│ Examples:                                                                   │
│ - Deadline tomorrow                                                         │
│ - Blocking 3+ other items                                                   │
│ - Customer escalation                                                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ CLASS: MEDIUM (Score 100-140)                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│ Characteristics:                                                            │
│ - Standard work queue                                                       │
│ - FIFO within class                                                         │
│ - Normal processing                                                         │
│ - SLA: Start within 1 day                                                   │
│                                                                             │
│ Examples:                                                                   │
│ - Regular feature work                                                      │
│ - Bug fixes (non-critical)                                                  │
│ - Planned improvements                                                      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ CLASS: LOW (Score 60-100)                                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│ Characteristics:                                                            │
│ - Fill-in work                                                              │
│ - When nothing higher priority available                                    │
│ - May be deferred                                                           │
│ - SLA: Start within 1 week                                                  │
│                                                                             │
│ Examples:                                                                   │
│ - Tech debt cleanup                                                         │
│ - Documentation                                                             │
│ - Nice-to-have features                                                     │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ CLASS: BACKLOG (Score < 60)                                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│ Characteristics:                                                            │
│ - Not actively scheduled                                                    │
│ - Reviewed periodically                                                     │
│ - May be closed as stale                                                    │
│ - SLA: None (explicitly)                                                    │
│                                                                             │
│ Examples:                                                                   │
│ - Future ideas                                                              │
│ - Someday/maybe                                                             │
│ - Requires more research                                                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Priority Recalculation

Priority is **dynamic**, not static:

```
WHEN PRIORITY IS RECALCULATED

1. TIME-BASED (automatic)
   - Every hour for all open items
   - Deadline proximity changes
   - Age increases

2. EVENT-TRIGGERED
   - New blocker relationship added
   - Escalation received
   - Objective priority changed
   - CEO directive issued

3. CONTEXT CHANGE
   - Sprint planning (batch recalc)
   - Quarterly planning (reset baselines)
   - Team capacity change

RECALCULATION EVENTS:
┌─────────────────────────────────────────────────────────────────────────────┐
│ Event                    │ Recalc Scope      │ Trigger                      │
├──────────────────────────┼───────────────────┼──────────────────────────────┤
│ Hourly tick              │ All open items    │ Cron: 0 * * * *              │
│ Blocker added            │ Single item       │ Event: dependency.created    │
│ Escalation               │ Single item       │ Event: work_order.escalated  │
│ Objective reprioritized  │ Linked items      │ Event: objective.updated     │
│ Sprint start             │ Sprint items      │ Event: sprint.started        │
│ Capacity change          │ Assigned items    │ Event: worker.capacity_changed│
└──────────────────────────┴───────────────────┴──────────────────────────────┘
```

### Preventing Priority Inflation

The problem: Over time, everything becomes "urgent" and "critical".

```
PRIORITY INFLATION CONTROLS

1. PRIORITY BUDGET
   - Max 5% of items can be P0
   - Max 15% can be P0 + P1
   - Enforced at creation time
   - To add P0, must demote something

2. PRIORITY DECAY FOR OVER-USE
   - If team has >10% CRITICAL, scores dampened by 20%
   - Forces prioritization decisions
   - "If everything is urgent, nothing is"

3. REQUIRED JUSTIFICATION
   - P0/P1 requires justification text
   - Justification is logged
   - Can be audited

4. EXPIRING URGENCY
   - External pressure flag expires after 48h
   - Must be re-confirmed to maintain
   - Prevents "forever urgent"

5. PRIORITY REVIEW
   - Weekly review of P0/P1 items
   - Supervisor must confirm or demote
   - Stale high-priority items auto-demote
```

---

## Questions for Iteration 2

1. **What signals indicate priority changes?**
   - How do we detect "this just became urgent"?
   - What events should trigger re-prioritization?

2. **How do we surface priority to workers?**
   - Priority queue visualization
   - "What should I work on next?"

3. **How do priorities interact across teams?**
   - Cross-team dependencies
   - Competing priorities

---

## Iteration 1 Key Insights

1. **Urgency ≠ Importance**: Time pressure vs impact are orthogonal

2. **Priority is dynamic**: Recalculated based on time, events, context

3. **Score-based classification**: Continuous score maps to priority classes

4. **Inflation must be controlled**: Budget, decay, expiration mechanisms

5. **Multiple signals combine**: Base priority + urgency factors + importance factors

---

## Iteration 2: Priority Signals & Detection

### Signal Sources

Where do priority signals come from?

```
PRIORITY SIGNAL SOURCES

┌─────────────────────────────────────────────────────────────────────────────┐
│                           EXTERNAL SIGNALS                                   │
│                     (Outside the system boundary)                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Source              │ Signal Type         │ Detection Method              │
│   ────────────────────┼─────────────────────┼─────────────────────────────  │
│   Customer            │ Escalation          │ Email, ticket, support        │
│   Stakeholder         │ Deadline            │ Calendar, meeting notes       │
│   Market              │ Competitive pressure│ Human input, news             │
│   Regulatory          │ Compliance deadline │ Calendar, legal input         │
│   Production          │ Incident            │ Monitoring, alerts            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           INTERNAL SIGNALS                                   │
│                      (Within the system boundary)                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Source              │ Signal Type         │ Detection Method              │
│   ────────────────────┼─────────────────────┼─────────────────────────────  │
│   Dependency graph    │ Blocker count       │ Automatic (graph analysis)    │
│   Time                │ Age, deadline prox  │ Automatic (clock)             │
│   Workflow            │ Stage timeout       │ Automatic (timer)             │
│   Verification        │ Failure count       │ Automatic (test results)      │
│   Worker              │ Explicit escalation │ Agent request                 │
│   Supervisor          │ Priority override   │ Manual decision               │
│   CEO                 │ Directive           │ Mail, explicit command        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Signal Detection Pipeline

```
SIGNAL DETECTION FLOW

External World                    VerMAS Boundary
     │                                  │
     │  ┌────────────────────────────────────────────────────────────────┐
     │  │                                                                │
     ▼  │  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐     │
  ───────┤  │   INGEST    │───▶│   CLASSIFY   │───▶│   ATTACH     │     │
 Signals │  │             │    │              │    │              │     │
         │  │ - Webhooks  │    │ - Type       │    │ - Find work  │     │
         │  │ - Email     │    │ - Severity   │    │   order      │     │
         │  │ - API       │    │ - Source     │    │ - Update     │     │
         │  │ - CLI       │    │              │    │   priority   │     │
         │  └──────────────┘    └──────────────┘    └──────────────┘     │
         │                                                │              │
         │                                                ▼              │
         │                                    ┌──────────────┐           │
         │                                    │   TRIGGER    │           │
         │                                    │              │           │
         │                                    │ - Recalc     │           │
         │                                    │ - Notify     │           │
         │                                    │ - Escalate   │           │
         │                                    └──────────────┘           │
         │                                                               │
         └───────────────────────────────────────────────────────────────┘
```

### Event-Driven Priority Updates

```python
# vermas/priority/signals.py

from enum import Enum
from dataclasses import dataclass

class SignalType(Enum):
    DEADLINE_APPROACHING = "deadline_approaching"
    BLOCKER_ADDED = "blocker_added"
    BLOCKER_RESOLVED = "blocker_resolved"
    ESCALATION_RECEIVED = "escalation_received"
    INCIDENT_REPORTED = "incident_reported"
    CUSTOMER_WAITING = "customer_waiting"
    WORK_STALE = "work_stale"
    VERIFICATION_FAILED = "verification_failed"
    OBJECTIVE_REPRIORITIZED = "objective_reprioritized"
    CEO_DIRECTIVE = "ceo_directive"

@dataclass
class PrioritySignal:
    signal_type: SignalType
    source: str  # Where it came from
    target_work_order: str  # Which work order affected
    magnitude: float  # How much to adjust (1.0 = normal)
    expires_at: Optional[datetime]  # When signal decays
    justification: str

# Signal handlers
SIGNAL_HANDLERS = {
    SignalType.DEADLINE_APPROACHING: lambda s: adjust_deadline_urgency(s),
    SignalType.BLOCKER_ADDED: lambda s: recalc_blocked_items(s),
    SignalType.ESCALATION_RECEIVED: lambda s: apply_escalation_boost(s),
    SignalType.INCIDENT_REPORTED: lambda s: create_critical_work_order(s),
    SignalType.CUSTOMER_WAITING: lambda s: apply_external_pressure(s),
    SignalType.CEO_DIRECTIVE: lambda s: apply_ceo_boost(s),
}

def process_signal(signal: PrioritySignal):
    """Process incoming priority signal."""
    handler = SIGNAL_HANDLERS.get(signal.signal_type)
    if handler:
        handler(signal)

    # Always log the signal
    log_priority_event(signal)

    # Trigger recalculation for affected work order
    recalculate_priority(signal.target_work_order)

    # Check if priority class changed
    check_priority_class_change(signal.target_work_order)
```

### Surfacing Priority to Workers

How do workers know what to work on?

```
WORKER PRIORITY VIEW

┌─────────────────────────────────────────────────────────────────────────────┐
│ YOUR WORK QUEUE                                        Updated: 2 min ago   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ ⚡ CRITICAL (do now)                                                        │
│ ────────────────────────────────────────────────────────────────────────── │
│ [wo-abc123] Fix authentication bypass (Score: 195)                          │
│             Deadline: 2h │ Blocks: 3 items │ CEO directive                  │
│                                                                             │
│ 🔴 HIGH (next up)                                                           │
│ ────────────────────────────────────────────────────────────────────────── │
│ [wo-def456] Customer data export feature (Score: 155)                       │
│             Deadline: tomorrow │ Customer waiting                           │
│                                                                             │
│ [wo-ghi789] API rate limiting (Score: 142)                                  │
│             Blocks: 2 items │ Security                                      │
│                                                                             │
│ 🟡 MEDIUM (standard queue)                                                  │
│ ────────────────────────────────────────────────────────────────────────── │
│ [wo-jkl012] Refactor user service (Score: 118)                              │
│             Contributes to: OBJ-001                                         │
│                                                                             │
│ [wo-mno345] Add logging to payment flow (Score: 105)                        │
│             Age: 5 days                                                     │
│                                                                             │
│ 🟢 LOW (when available)                                                     │
│ ────────────────────────────────────────────────────────────────────────── │
│ [wo-pqr678] Update README (Score: 72)                                       │
│             No blockers, no deadline                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

COMMANDS:
  next        - Start working on highest priority item
  why <id>    - Explain why item has this priority
  bump <id>   - Request priority increase (needs justification)
```

### "Why This Priority?" Explainability

```python
# vermas/priority/explainer.py

def explain_priority(work_order_id: str) -> PriorityExplanation:
    """Explain why a work order has its current priority."""
    wo = get_work_order(work_order_id)
    factors = calculate_priority_factors(wo)
    score = calculate_priority_score(factors)

    explanation = PriorityExplanation(
        work_order_id=work_order_id,
        total_score=score,
        priority_class=score_to_class(score),
        breakdown=[],
    )

    # Base priority contribution
    explanation.breakdown.append(
        FactorContribution(
            factor="Base Priority",
            value=f"P{factors.base_priority}",
            points=base_scores[factors.base_priority],
            reason="Set at creation"
        )
    )

    # Deadline contribution
    if factors.deadline:
        hours = hours_until(factors.deadline)
        deadline_points = deadline_to_points(hours)
        explanation.breakdown.append(
            FactorContribution(
                factor="Deadline",
                value=f"{hours:.0f}h remaining",
                points=deadline_points,
                reason=f"Deadline: {factors.deadline}"
            )
        )

    # Blockers contribution
    if factors.blocked_items:
        blocker_points = min(len(factors.blocked_items) * 10, 30)
        explanation.breakdown.append(
            FactorContribution(
                factor="Blocking Others",
                value=f"{len(factors.blocked_items)} items",
                points=blocker_points,
                reason=f"Blocks: {', '.join(factors.blocked_items)}"
            )
        )

    # ... more factors

    return explanation
```

### Cross-Team Priority Conflicts

When teams have competing priorities:

```
CROSS-TEAM PRIORITY RESOLUTION

Scenario: Team A needs work from Team B, but Team B has different priorities

┌─────────────────────────────────────────────────────────────────────────────┐
│ TEAM A                                    TEAM B                             │
│                                                                             │
│ [wo-a1] Feature X (P1)                   [wo-b1] Bug fix (P0)               │
│    └── Depends on ───────────────────────▶ [wo-b2] API change (P3)          │
│                                                                             │
│ Team A sees wo-b2 as blocking            Team B sees wo-b2 as low priority  │
│ their P1 work                                                               │
└─────────────────────────────────────────────────────────────────────────────┘

RESOLUTION STRATEGIES:

1. BLOCKER ESCALATION
   - wo-b2's priority inherits from blocked item
   - If wo-a1 is P1 and blocked by wo-b2, wo-b2 gets +30 points
   - Automatic, no negotiation needed

2. CROSS-TEAM VISIBILITY
   - Team B sees "Blocking external: Team A / wo-a1 (P1)"
   - Creates social pressure to address

3. SUPERVISOR NEGOTIATION
   - If conflict persists, supervisors negotiate
   - Can agree on priority or timeline

4. OPERATIONS ARBITRATION
   - If supervisors can't agree, escalate to operations
   - Operations makes binding decision

5. CEO OVERRIDE
   - For strategic conflicts
   - CEO directive supersedes all
```

### Priority Signal Events

```yaml
# Events emitted for priority changes

priority.signal_received:
  signal_type: string
  source: string
  target: string  # work order ID
  magnitude: float

priority.recalculated:
  work_order_id: string
  old_score: float
  new_score: float
  old_class: string
  new_class: string
  changed_factors: list

priority.class_changed:
  work_order_id: string
  old_class: string
  new_class: string
  reason: string
  notifications_sent: list

priority.conflict_detected:
  teams: list
  work_orders: list
  resolution_path: string
```

---

## Questions for Iteration 3

1. **How does priority decay over time?**
   - What happens to work that sits too long?
   - How do we prevent eternal deferral?

2. **What about priority "freshness"?**
   - New work vs old work
   - Preventing stale priorities

3. **How do we handle priority exhaustion?**
   - When workers are burned out on CRITICAL items
   - Capacity for sustained urgency

---

## Iteration 2 Key Insights

1. **Signals come from many sources**: External (customers, incidents) and internal (dependencies, time)

2. **Detection is pipelined**: Ingest → Classify → Attach → Trigger

3. **Workers need clear views**: Priority queue with explanations

4. **Priority is explainable**: "Why this priority?" with factor breakdown

5. **Cross-team conflicts are resolved**: Escalation ladder from automatic to CEO

---

## Iteration 3: Priority Decay & Staleness

### The Staleness Problem

Work that sits in the queue too long creates multiple problems:

```
PROBLEMS WITH STALE WORK

1. CONTEXT LOSS
   - Original requirements may be outdated
   - People who understood it may have moved on
   - Codebase has changed underneath

2. RELEVANCE DECAY
   - The problem may have been solved another way
   - The feature may no longer be needed
   - Business context has shifted

3. HIDDEN COST
   - Queue management overhead
   - Mental load of "open items"
   - False sense of progress

4. PRIORITY BLINDNESS
   - Old P3 items never get done
   - Creates culture of ignoring low priority
   - Backlog becomes graveyard
```

### Two Types of Decay

```
PRIORITY DECAY MECHANISMS

┌─────────────────────────────────────────────────────────────────────────────┐
│ TYPE 1: URGENCY BOOST (Positive Decay)                                       │
│                                                                             │
│ Old work becomes MORE urgent over time                                      │
│                                                                             │
│ Score │                                    ┌───── Urgency boost kicks in    │
│       │                                ┌───┘                                │
│       │                            ┌───┘                                    │
│       │                        ┌───┘                                        │
│       │────────────────────────┘ Base priority                              │
│       │                                                                     │
│       └──────────────────────────────────────────────────▶ Age              │
│           0d        7d        14d       21d       28d                       │
│                                                                             │
│ Purpose: Prevent indefinite deferral of low-priority work                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ TYPE 2: RELEVANCE DECAY (Negative Decay)                                     │
│                                                                             │
│ Very old work becomes LESS relevant / potentially stale                     │
│                                                                             │
│ Relevance │                                                                 │
│    100%   │────────────┐                                                    │
│           │            └────────┐                                           │
│           │                     └────────┐                                  │
│           │                              └────────┐                         │
│           │                                       └──── → Review required   │
│           └──────────────────────────────────────────────▶ Age              │
│               0d       30d      60d      90d     120d                       │
│                                                                             │
│ Purpose: Force review of ancient items - close or refresh                   │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Age-Based Priority Boost

```python
# vermas/priority/decay.py

def calculate_age_boost(created_at: datetime, base_priority: int) -> float:
    """
    Calculate priority boost based on age.
    Low-priority items get more boost (to prevent eternal deferral).
    """
    age_days = (datetime.now() - created_at).days

    # Base boost: 0.5 points per day, capped at 14 points
    base_boost = min(age_days * 0.5, 14)

    # Priority multiplier: Lower priority gets more boost
    # P0: 0.5x, P1: 0.75x, P2: 1x, P3: 1.5x, P4: 2x
    priority_multipliers = {0: 0.5, 1: 0.75, 2: 1.0, 3: 1.5, 4: 2.0}
    multiplier = priority_multipliers.get(base_priority, 1.0)

    return base_boost * multiplier

# Example:
# P4 item aged 30 days: min(30 * 0.5, 14) * 2.0 = 14 * 2.0 = 28 point boost
# P1 item aged 30 days: min(30 * 0.5, 14) * 0.75 = 14 * 0.75 = 10.5 point boost
```

### Staleness Detection & Alerts

```yaml
# .work/governance/staleness-rules.yaml

staleness:
  thresholds:
    warning: 30d    # 30 days without activity
    critical: 60d   # 60 days without activity
    stale: 90d      # 90 days - requires action

  actions:
    warning:
      - notify: assignee
      - add_label: "needs-attention"

    critical:
      - notify: [assignee, supervisor]
      - add_label: "stale-risk"
      - create_review_task: true

    stale:
      - notify: [assignee, supervisor, operations]
      - options:
          - close_as_stale
          - reassign
          - refresh_requirements
          - escalate_to_human

  exceptions:
    # Some work types don't go stale
    - type: "documentation"
      threshold_multiplier: 2.0  # 180 days before stale

    - type: "research"
      threshold_multiplier: 1.5

    - type: "blocked_by_external"
      exempt: true  # Don't mark as stale while blocked
```

### Staleness Review Workflow

```
STALE ITEM REVIEW PROCESS

Item reaches 90-day threshold
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ STALENESS REVIEW NOTIFICATION                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ [wo-abc123] Implement caching layer                                         │
│                                                                             │
│ Status: STALE (93 days without activity)                                    │
│ Last activity: 2024-10-05 (code review comment)                             │
│ Assigned to: worker-1                                                        │
│                                                                             │
│ ⚠️ This item requires a decision:                                           │
│                                                                             │
│ Options:                                                                    │
│   [1] CLOSE - No longer needed                                              │
│   [2] REFRESH - Update requirements and restart                             │
│   [3] REASSIGN - Give to someone else                                       │
│   [4] DEFER - Move to backlog with new target date                          │
│   [5] ESCALATE - Need human decision                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ Supervisor decides within 7 days, or auto-escalates to Operations           │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Preventing Priority Exhaustion

When everything is urgent for too long:

```
PRIORITY EXHAUSTION DETECTION

Symptoms:
- High percentage of CRITICAL items (>10%)
- Workers constantly interrupted
- Sustained urgency >2 weeks
- Completion rate dropping
- Quality metrics declining

┌─────────────────────────────────────────────────────────────────────────────┐
│ EXHAUSTION DASHBOARD                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Current State:                                                              │
│   CRITICAL items:     8 (12%)  ⚠️ Above threshold                           │
│   HIGH items:        15 (23%)  ⚠️ Above threshold                           │
│   Avg time in CRITICAL: 3.2 days                                            │
│   Interruptions/day:  4.5                                                   │
│                                                                             │
│ Trend (last 14 days):                                                       │
│   CRITICAL │ ▄▄▄▅▅▆▆▆▇▇▇███  ← Increasing (bad)                            │
│   Velocity │ ████▇▇▆▆▅▅▄▄▃▃  ← Decreasing (bad)                            │
│   Quality  │ ███▇▇▆▆▅▅▄▄▃▃▃  ← Decreasing (bad)                            │
│                                                                             │
│ RECOMMENDATION: Priority reset needed                                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Priority Reset Protocol

When priority inflation gets out of control:

```
PRIORITY RESET PROTOCOL

Triggered when:
- CRITICAL > 10% for > 7 days
- Velocity dropped > 30%
- Supervisor or Operations requests

Steps:

1. FREEZE NEW WORK
   - No new items can be marked CRITICAL/HIGH
   - Exception: CEO directive or production incident

2. TRIAGE SESSION
   - Supervisor reviews all CRITICAL/HIGH items
   - Each item: Confirm priority or demote
   - Must provide justification for keeping high

3. BATCH DEMOTION
   - Items not confirmed are demoted one level
   - CRITICAL → HIGH
   - HIGH → MEDIUM

4. ROOT CAUSE ANALYSIS
   - Why did inflation happen?
   - Process change needed?
   - Staffing issue?

5. RESUME NORMAL OPERATIONS
   - Lift freeze
   - Monitor for recurrence

LOG ENTRY:
┌─────────────────────────────────────────────────────────────────────────────┐
│ PRIORITY RESET EVENT                                                         │
│                                                                             │
│ Date: 2026-01-15                                                            │
│ Trigger: CRITICAL > 10% for 9 days                                          │
│                                                                             │
│ Before:                                                                     │
│   CRITICAL: 12 items (15%)                                                  │
│   HIGH: 18 items (22%)                                                      │
│                                                                             │
│ After:                                                                      │
│   CRITICAL: 3 items (4%)                                                    │
│   HIGH: 9 items (11%)                                                       │
│                                                                             │
│ Demoted items: 18                                                           │
│ Reviewed by: supervisor-alpha                                               │
│ Root cause: External deadline pressure from 3 customers simultaneously      │
│ Action: Implement customer deadline coordination                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Backlog Grooming Automation

```python
# vermas/priority/grooming.py

class BacklogGroomer:
    """Automated backlog maintenance."""

    def run_grooming(self):
        """Weekly backlog grooming routine."""
        results = GroomingResults()

        # 1. Identify stale items
        stale_items = self.find_stale_items(threshold_days=90)
        for item in stale_items:
            results.stale.append(self.create_review_task(item))

        # 2. Find duplicates
        duplicates = self.detect_duplicates()
        for dup_group in duplicates:
            results.duplicates.append(self.suggest_merge(dup_group))

        # 3. Check for orphaned items (no objective alignment)
        orphans = self.find_orphaned_items()
        for orphan in orphans:
            results.orphans.append(self.request_alignment(orphan))

        # 4. Verify blocked items (are blockers still valid?)
        blocked = self.find_blocked_items()
        for item in blocked:
            if not self.verify_blocker_exists(item.blocked_by):
                results.fixed.append(self.remove_invalid_blocker(item))

        # 5. Close items that missed their deadline by > 30 days
        missed_deadlines = self.find_missed_deadlines(grace_days=30)
        for item in missed_deadlines:
            results.auto_closed.append(self.close_as_obsolete(item))

        return results

    def generate_report(self, results: GroomingResults) -> str:
        """Generate weekly grooming report."""
        return f"""
WEEKLY BACKLOG GROOMING REPORT
==============================

Stale items requiring review: {len(results.stale)}
Duplicate groups detected: {len(results.duplicates)}
Orphaned items (no objective): {len(results.orphans)}
Invalid blockers fixed: {len(results.fixed)}
Auto-closed (obsolete): {len(results.auto_closed)}

Action items created: {results.total_action_items}
        """
```

---

## Questions for Iteration 4

Now transitioning to **Extensibility**:

1. **What is the extension model?**
   - Plugins, tools, skills, sub-agents
   - How does Claude Code do it?

2. **How do extensions integrate?**
   - Discovery, installation, invocation
   - Permission model

3. **Who owns extensions?**
   - Individual, team, organization, ecosystem

---

## Iteration 3 Key Insights

1. **Two types of decay**: Urgency boost (old work rises) and relevance decay (ancient work gets reviewed)

2. **Staleness has thresholds**: 30d warning, 60d critical, 90d requires action

3. **Priority exhaustion is real**: Detect and reset when everything is CRITICAL

4. **Backlog grooming can be automated**: Stale items, duplicates, orphans, invalid blockers

5. **Decay prevents eternal deferral**: Low-priority work eventually rises

---

## Iteration 4: Extensibility Model

### The Need for External Expertise

Sometimes your team doesn't have the skills. You need to bring in:

```
TYPES OF EXTERNAL EXPERTISE

1. SPECIALIST CONSULTANTS
   - Security auditor
   - Performance expert
   - Accessibility specialist
   - Legal/compliance reviewer

2. DOMAIN EXPERTS
   - Machine learning for a specific task
   - Payment processing integration
   - Regulatory compliance knowledge
   - Industry-specific requirements

3. TOOL INTEGRATIONS
   - Code analysis tools (SonarQube, etc.)
   - Testing frameworks
   - Deployment pipelines
   - Monitoring systems

4. REUSABLE PATTERNS
   - Common workflows
   - Best practices
   - Organizational templates
   - Industry standards
```

### Learning from Claude Code's Model

Claude Code has an extensibility model we can learn from:

```
CLAUDE CODE EXTENSION TYPES

┌─────────────────────────────────────────────────────────────────────────────┐
│ MCP SERVERS (Model Context Protocol)                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│ • External services providing tools to the agent                            │
│ • Examples: filesystem, database, API integrations                          │
│ • Installed via configuration                                               │
│ • Provides: tools, resources, prompts                                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ SLASH COMMANDS / SKILLS                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│ • User-invokable actions (/commit, /review-pr)                              │
│ • Defined in .claude/commands/ or skills files                              │
│ • Can be custom prompts or complex workflows                                │
│ • Scoped: user, project, organization                                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ HOOKS                                                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│ • Event-triggered actions                                                   │
│ • Run before/after tool calls                                               │
│ • Can modify, block, or augment behavior                                    │
│ • Examples: pre-commit checks, post-edit formatting                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ SUB-AGENTS                                                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│ • Specialized agents for specific tasks                                     │
│ • Invoked via Task tool                                                     │
│ • Have their own tool access and context                                    │
│ • Examples: Explore agent, Plan agent                                       │
└─────────────────────────────────────────────────────────────────────────────┘
```

### VerMAS Extension Types

Mapping to VerMAS organizational model:

```
VERMAS EXTENSION MODEL

┌─────────────────────────────────────────────────────────────────────────────┐
│ PLUGINS (like MCP Servers)                                                   │
│ "Hiring external consultants with specific tools"                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ What: External services providing specialized capabilities                  │
│                                                                             │
│ Examples:                                                                   │
│ • Security scanner plugin (provides: scan_code, check_vulns tools)          │
│ • Translation service (provides: translate, detect_language tools)          │
│ • Code quality analyzer (provides: lint, complexity_analysis tools)         │
│                                                                             │
│ Installation:                                                               │
│ • Organization-level: Available to all factories                            │
│ • Factory-level: Available to one factory                                   │
│                                                                             │
│ Trust model:                                                                │
│ • Sandboxed execution                                                       │
│ • Declared permissions                                                      │
│ • Audit logging of all calls                                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ SKILLS (like Slash Commands)                                                 │
│ "Standard operating procedures anyone can invoke"                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ What: Named workflows or procedures that workers can invoke                 │
│                                                                             │
│ Examples:                                                                   │
│ • /security-review - Run security checklist                                 │
│ • /deploy-staging - Deploy to staging environment                           │
│ • /create-migration - Create database migration                             │
│ • /onboard-service - Set up new microservice                                │
│                                                                             │
│ Definition:                                                                 │
│ • YAML workflow files in .work/skills/                                      │
│ • Or prompts that guide agent behavior                                      │
│                                                                             │
│ Scopes:                                                                     │
│ • Worker-level: Personal shortcuts                                          │
│ • Factory-level: Team procedures                                            │
│ • Organization-level: Company standards                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ EXPERTS (like Sub-Agents)                                                    │
│ "Specialist consultants you can call in for specific problems"              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ What: Specialized agent configurations for domain expertise                 │
│                                                                             │
│ Examples:                                                                   │
│ • security-expert: Trained on security best practices                       │
│ • api-designer: Specialized in REST/GraphQL design                          │
│ • database-optimizer: Query performance specialist                          │
│ • accessibility-auditor: WCAG compliance expert                             │
│                                                                             │
│ Invocation:                                                                 │
│ • Worker requests expert consultation                                       │
│ • Expert reviews work and provides feedback                                 │
│ • Feedback attached to work order                                           │
│                                                                             │
│ Trust model:                                                                │
│ • Experts don't commit code directly                                        │
│ • Provide recommendations that worker implements                            │
│ • Or provide approval gates                                                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ TEMPLATES (like Starter Kits)                                                │
│ "Best practices packages you can adopt"                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ What: Pre-packaged configurations, workflows, and standards                 │
│                                                                             │
│ Examples:                                                                   │
│ • hipaa-compliance-template: Healthcare compliance setup                    │
│ • startup-velocity-template: Fast-moving startup configuration              │
│ • enterprise-audit-template: Large enterprise audit requirements            │
│                                                                             │
│ Contents:                                                                   │
│ • Compliance rules                                                          │
│ • Workflow definitions                                                      │
│ • Skill definitions                                                         │
│ • Expert configurations                                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Extension Discovery & Installation

```
EXTENSION LIFECYCLE

┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. DISCOVERY                                                                 │
│                                                                             │
│ Sources:                                                                    │
│ • Official VerMAS extension registry                                        │
│ • Organization's private registry                                           │
│ • Git repositories                                                          │
│ • Local files                                                               │
│                                                                             │
│ CLI:                                                                        │
│   vermas extension search "security"                                        │
│   vermas extension list --source=registry                                   │
│   vermas extension info security-scanner@1.2.3                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. INSTALLATION                                                              │
│                                                                             │
│ Scoped installation:                                                        │
│   vermas extension install security-scanner --scope=org                     │
│   vermas extension install code-quality --scope=factory                     │
│   vermas extension install my-shortcuts --scope=worker                      │
│                                                                             │
│ From source:                                                                │
│   vermas extension install git@github.com:org/extension.git                 │
│   vermas extension install ./local-extension/                               │
│                                                                             │
│ Approval flow:                                                              │
│   - Worker install: Immediate (personal scope only)                         │
│   - Factory install: Supervisor approval                                    │
│   - Org install: CEO approval                                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. CONFIGURATION                                                             │
│                                                                             │
│ Extension manifest:                                                         │
│   # .work/extensions/security-scanner/manifest.yaml                         │
│   name: security-scanner                                                    │
│   version: 1.2.3                                                            │
│   permissions:                                                              │
│     - read:code                                                             │
│     - write:reports                                                         │
│   config:                                                                   │
│     severity_threshold: medium                                              │
│     ignore_paths: [test/*, vendor/*]                                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. INVOCATION                                                                │
│                                                                             │
│ Plugin tools:                                                               │
│   Worker uses `security_scan` tool in normal workflow                       │
│                                                                             │
│ Skills:                                                                     │
│   Worker invokes `/security-review` skill                                   │
│                                                                             │
│ Experts:                                                                    │
│   Worker requests `security-expert` review                                  │
│                                                                             │
│ Automatic (hooks):                                                          │
│   Pre-merge hook invokes security scanner                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Extension Manifest Schema

```yaml
# .work/extensions/security-scanner/manifest.yaml

# Identity
name: security-scanner
version: 1.2.3
description: "Static analysis for security vulnerabilities"
author: "security-tools-org"
license: MIT

# Type
type: plugin  # plugin | skill | expert | template

# Compatibility
requires:
  vermas: ">=1.0.0"
  python: ">=3.10"

# Permissions (what the extension can do)
permissions:
  - read:code          # Can read source files
  - read:config        # Can read configuration
  - write:reports      # Can write to reports directory
  - invoke:git         # Can invoke git commands
  # Cannot: write:code, invoke:network, etc.

# For plugins: what tools are provided
provides:
  tools:
    - name: scan_code
      description: "Scan code for security vulnerabilities"
      parameters:
        - name: path
          type: string
          required: true
        - name: severity
          type: string
          enum: [low, medium, high, critical]
          default: medium
      returns: ScanResult

    - name: check_dependencies
      description: "Check dependencies for known vulnerabilities"
      parameters:
        - name: manifest_path
          type: string
      returns: DependencyReport

# For skills: workflow definition
skill:
  invocation: "/security-review"
  workflow: |
    1. Run scan_code on changed files
    2. Run check_dependencies
    3. Generate report
    4. If critical findings, block merge
    5. Else, attach report to work order

# For experts: agent configuration
expert:
  name: security-expert
  system_prompt: |
    You are a security expert specializing in...
  tools: [scan_code, check_dependencies]
  review_mode: true  # Can only advise, not commit

# Configuration schema
config_schema:
  severity_threshold:
    type: string
    enum: [low, medium, high, critical]
    default: medium
  ignore_paths:
    type: array
    items: { type: string }
    default: []
  report_format:
    type: string
    enum: [json, markdown, sarif]
    default: markdown
```

### Permission Model

```
EXTENSION PERMISSION SYSTEM

┌─────────────────────────────────────────────────────────────────────────────┐
│ PERMISSION CATEGORIES                                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ READ PERMISSIONS:                                                           │
│   read:code         - Read source code files                                │
│   read:config       - Read configuration files                              │
│   read:work_orders  - Read work order data                                  │
│   read:events       - Read event log                                        │
│   read:secrets      - Read secrets (dangerous!)                             │
│                                                                             │
│ WRITE PERMISSIONS:                                                          │
│   write:code        - Modify source files (rare for extensions)             │
│   write:reports     - Write to reports directory                            │
│   write:config      - Modify configuration                                  │
│   write:work_orders - Create/update work orders                             │
│                                                                             │
│ INVOKE PERMISSIONS:                                                         │
│   invoke:git        - Run git commands                                      │
│   invoke:shell      - Run arbitrary shell commands (dangerous!)             │
│   invoke:network    - Make network requests                                 │
│   invoke:llm        - Call LLM APIs                                         │
│                                                                             │
│ SPECIAL PERMISSIONS:                                                        │
│   block:merge       - Can block merge operations                            │
│   approve:work      - Can approve work orders                               │
│   create:workers    - Can spawn new workers                                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

PERMISSION LEVELS:

┌──────────────┬─────────────────────────────────────────────────────────────┐
│ Level        │ Allowed Permissions                                         │
├──────────────┼─────────────────────────────────────────────────────────────┤
│ SAFE         │ read:code, read:config, write:reports                       │
│              │ (Default for new extensions)                                │
├──────────────┼─────────────────────────────────────────────────────────────┤
│ STANDARD     │ SAFE + invoke:git, invoke:network, read:work_orders         │
│              │ (Requires supervisor approval)                              │
├──────────────┼─────────────────────────────────────────────────────────────┤
│ ELEVATED     │ STANDARD + write:code, block:merge, invoke:llm              │
│              │ (Requires CEO approval)                                     │
├──────────────┼─────────────────────────────────────────────────────────────┤
│ PRIVILEGED   │ ELEVATED + invoke:shell, read:secrets, create:workers       │
│              │ (Requires human approval + security review)                 │
└──────────────┴─────────────────────────────────────────────────────────────┘
```

---

## Questions for Iteration 5

1. **How do experts integrate into workflows?**
   - Consultation requests
   - Review gates
   - Feedback loops

2. **How do we pay for/allocate expert time?**
   - Rate limiting
   - Budget constraints
   - Priority access

3. **How do experts build reputation?**
   - Quality of reviews
   - Accuracy of recommendations

---

## Iteration 4 Key Insights

1. **Four extension types**: Plugins (tools), Skills (workflows), Experts (agents), Templates (packages)

2. **Scoped installation**: Worker, Factory, Organization levels

3. **Permission model**: Read/Write/Invoke categories with approval levels

4. **Discovery via registry**: Official, private, git, local sources

5. **Claude Code parallel**: MCP Servers → Plugins, Skills → Skills, Sub-agents → Experts

---

## Iteration 5: Expert/Consultant Integration

### Expert Workflow Integration

How do experts fit into the work lifecycle?

```
EXPERT INTEGRATION PATTERNS

┌─────────────────────────────────────────────────────────────────────────────┐
│ PATTERN 1: CONSULTATION (On-Demand)                                          │
│ "I need advice on how to approach this"                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Worker                                 Expert                              │
│      │                                     │                                │
│      │  Request consultation               │                                │
│      ├────────────────────────────────────▶│                                │
│      │                                     │ Review context                 │
│      │                                     ├───┐                            │
│      │                                     │◀──┘                            │
│      │◀────────────────────────────────────┤ Provide recommendations        │
│      │                                     │                                │
│      │  Implement (or not)                 │                                │
│      ├───┐                                 │                                │
│      │◀──┘                                 │                                │
│                                                                             │
│ Trigger: Worker requests via /consult security-expert                       │
│ Output: Recommendations attached to work order                              │
│ Authority: Advisory only                                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ PATTERN 2: REVIEW GATE (Mandatory)                                           │
│ "This must be reviewed before proceeding"                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Worker                    Gate                    Expert                  │
│      │                        │                        │                    │
│      │  Submit for review     │                        │                    │
│      ├───────────────────────▶│                        │                    │
│      │                        │  Trigger review        │                    │
│      │                        ├───────────────────────▶│                    │
│      │                        │                        │ Review             │
│      │                        │                        ├───┐                │
│      │                        │                        │◀──┘                │
│      │                        │◀───────────────────────┤ PASS/FAIL         │
│      │◀───────────────────────┤                        │                    │
│      │  Gate opens/blocks     │                        │                    │
│                                                                             │
│ Trigger: Work order reaches specific state (e.g., ready_for_security)       │
│ Output: PASS/FAIL decision, findings attached                               │
│ Authority: Can block progress                                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ PATTERN 3: PAIR WORK (Collaborative)                                         │
│ "Work alongside an expert"                                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Worker                                 Expert                              │
│      │                                     │                                │
│      ├────────────────────────────────────▶│ Join session                  │
│      │                                     │                                │
│      │◀───────────────────────────────────▶│ Collaborative work            │
│      │       Real-time feedback            │                                │
│      │       Guidance                      │                                │
│      │       Education                     │                                │
│      │                                     │                                │
│      ├────────────────────────────────────▶│ Session ends                  │
│      │                                     │                                │
│                                                                             │
│ Trigger: Work order tagged with #needs-expert-support                       │
│ Output: Work completed with expert guidance, learning recorded              │
│ Authority: Expert advises, worker implements                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ PATTERN 4: DELEGATION (Handoff)                                              │
│ "This requires specialist skills I don't have"                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Worker                 Supervisor                   Expert                │
│      │                        │                        │                    │
│      │  Request escalation    │                        │                    │
│      ├───────────────────────▶│                        │                    │
│      │                        │  Approve delegation    │                    │
│      │                        ├───────────────────────▶│                    │
│      │                        │                        │ Take over work     │
│      │                        │                        ├───┐                │
│      │                        │                        │◀──┘                │
│      │                        │◀───────────────────────┤ Complete           │
│      │◀───────────────────────┤                        │                    │
│      │  Resume downstream     │                        │                    │
│                                                                             │
│ Trigger: Worker requests, supervisor approves                               │
│ Output: Expert completes the work order or subtask                          │
│ Authority: Expert becomes temporary owner                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Expert Budget & Rate Limiting

Experts are a limited (and potentially expensive) resource:

```yaml
# .work/governance/expert-budget.yaml

budget:
  # Organization-wide limits
  org:
    monthly_expert_hours: 100
    max_concurrent_consultations: 5
    priority_access:
      - role: ceo
      - role: operations
      - priority: P0

  # Per-factory limits
  factory:
    monthly_expert_hours: 20
    max_queue_depth: 10

  # Per-expert limits
  experts:
    security-expert:
      hourly_rate: 2.0  # 2x normal worker cost
      max_concurrent: 3
      specialties: [security, compliance, auth]
      response_sla:
        P0: 15m
        P1: 1h
        P2: 4h
        P3: 24h

    performance-expert:
      hourly_rate: 1.5
      max_concurrent: 2
      specialties: [performance, database, caching]

# Rate limiting
rate_limits:
  consultation_per_worker_per_day: 5
  review_gate_timeout: 4h  # Auto-escalate if no response
  pair_session_max_duration: 2h
```

### Expert Request Flow

```python
# vermas/experts/request.py

@dataclass
class ExpertRequest:
    request_id: str
    requester: str  # Worker or supervisor
    expert_type: str  # security-expert, etc.
    pattern: ExpertPattern  # CONSULTATION, REVIEW_GATE, etc.
    work_order_id: str
    priority: int
    context: str  # What they need help with
    urgency_justification: Optional[str]

class ExpertRequestHandler:
    def request_expert(self, request: ExpertRequest) -> ExpertRequestResult:
        # 1. Check budget
        if not self.budget_available(request):
            return ExpertRequestResult(
                status="rejected",
                reason="Budget exhausted for this period",
                alternative="Try again next month or request budget increase"
            )

        # 2. Check rate limits
        if self.rate_limited(request):
            return ExpertRequestResult(
                status="rejected",
                reason="Rate limit exceeded",
                retry_after=self.next_available_slot(request)
            )

        # 3. Check expert availability
        expert = self.find_available_expert(request.expert_type)
        if not expert:
            return ExpertRequestResult(
                status="queued",
                position=self.queue_position(request),
                estimated_wait=self.estimate_wait_time(request)
            )

        # 4. Create expert session
        session = self.create_session(request, expert)

        # 5. Notify expert
        self.notify_expert(expert, session)

        return ExpertRequestResult(
            status="accepted",
            session_id=session.id,
            expert=expert.id,
            expected_start=session.scheduled_start
        )
```

### Expert Reputation System

How do we know which experts are good?

```
EXPERT REPUTATION MODEL

┌─────────────────────────────────────────────────────────────────────────────┐
│ REPUTATION FACTORS                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Factor                    │ Weight │ Measurement                            │
│ ──────────────────────────┼────────┼─────────────────────────────────────── │
│ Review accuracy           │ 30%    │ % of issues found vs missed            │
│ Response time             │ 20%    │ Time to first response                 │
│ Recommendation quality    │ 25%    │ Were recommendations followed?         │
│ Worker satisfaction       │ 15%    │ Post-review feedback rating            │
│ False positive rate       │ 10%    │ Flags that weren't real issues         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

REPUTATION SCORE CALCULATION:

```python
def calculate_reputation(expert: Expert, period_days: int = 90) -> float:
    sessions = get_sessions(expert, days=period_days)

    # Review accuracy: Issues found that were valid
    valid_issues = sum(s.valid_issues_found for s in sessions)
    total_issues = sum(s.total_issues_reported for s in sessions)
    accuracy = valid_issues / total_issues if total_issues else 1.0

    # Response time: Average within SLA
    response_times = [s.first_response_time for s in sessions]
    avg_response = mean(response_times)
    response_score = 1.0 - (avg_response / SLA_TARGET)

    # Recommendation quality: % followed
    followed = sum(1 for s in sessions if s.recommendation_followed)
    rec_quality = followed / len(sessions) if sessions else 0.5

    # Worker satisfaction: Average rating
    ratings = [s.worker_rating for s in sessions if s.worker_rating]
    satisfaction = mean(ratings) / 5.0 if ratings else 0.5

    # False positive rate: Lower is better
    false_positives = sum(s.false_positive_count for s in sessions)
    fp_rate = 1.0 - (false_positives / total_issues if total_issues else 0)

    # Weighted combination
    score = (
        0.30 * accuracy +
        0.20 * response_score +
        0.25 * rec_quality +
        0.15 * satisfaction +
        0.10 * fp_rate
    )

    return min(max(score, 0.0), 1.0)  # Clamp 0-1
```

### Expert Profile

```yaml
# .work/experts/security-expert/profile.yaml

expert:
  id: security-expert
  name: "Security Expert"
  description: "Specialized in application security, OWASP, auth"

  # Specialties (for matching to requests)
  specialties:
    - security
    - authentication
    - authorization
    - owasp
    - compliance/soc2
    - compliance/hipaa

  # System prompt for the expert agent
  system_prompt: |
    You are a security expert reviewing code and designs.
    Focus on:
    - OWASP Top 10 vulnerabilities
    - Authentication and authorization flaws
    - Data exposure risks
    - Injection attacks
    - Security misconfigurations

    When reviewing:
    1. Identify specific issues with file:line references
    2. Classify severity (critical/high/medium/low)
    3. Provide concrete remediation steps
    4. Note security best practices not followed

  # Tools available to this expert
  tools:
    - read_code
    - search_code
    - security_scan  # From security-scanner plugin
    - check_dependencies

  # Metrics
  metrics:
    reviews_completed: 142
    avg_issues_per_review: 3.2
    accuracy_rate: 0.94
    avg_response_time_hours: 1.8
    reputation_score: 0.91

  # Availability
  availability:
    max_concurrent: 3
    queue_limit: 10
    response_sla:
      P0: 15m
      P1: 1h
      P2: 4h
```

---

## Questions for Iteration 6

1. **How do we manage a skills/capabilities registry?**
   - What skills exist in the organization?
   - Skill matching to work orders

2. **How do skills evolve and improve?**
   - Version control for skills
   - A/B testing skills

3. **How do we know what capabilities we need vs have?**
   - Gap analysis
   - Skill investment decisions

---

## Iteration 5 Key Insights

1. **Four integration patterns**: Consultation, Review Gate, Pair Work, Delegation

2. **Budget controls**: Monthly limits, rate limiting, priority access

3. **Reputation is multi-dimensional**: Accuracy, speed, quality, satisfaction, false positives

4. **Experts are agents with specialized prompts**: Configuration defines expertise

5. **Request flow includes availability check**: Queue when busy, reject when over budget

---

## Iteration 6: Skills & Capabilities Registry

### What is a Skill?

A skill is a **named, reusable capability** that can be invoked by workers.

```
SKILL TAXONOMY

┌─────────────────────────────────────────────────────────────────────────────┐
│ ATOMIC SKILLS                                                                │
│ Single, focused operations                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│ • /lint - Run linter on code                                                │
│ • /test - Run test suite                                                    │
│ • /format - Format code                                                     │
│ • /build - Compile/build project                                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ COMPOSITE SKILLS                                                             │
│ Workflows combining multiple steps                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│ • /deploy-staging - Build → Test → Deploy to staging → Smoke test           │
│ • /release - Version bump → Changelog → Tag → Build → Deploy                │
│ • /security-review - Scan → Audit → Report → Gate                           │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ KNOWLEDGE SKILLS                                                             │
│ Domain expertise encoded as prompts                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│ • /explain-auth - Explain our authentication architecture                   │
│ • /api-standards - Our API design conventions                               │
│ • /onboarding-checklist - New service setup requirements                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Skill Registry

```yaml
# .work/registry/skills.yaml

registry:
  version: 1
  last_updated: 2026-01-07

skills:
  # Atomic skill
  - id: lint
    name: "Code Linting"
    type: atomic
    scope: org
    invocation: "/lint"
    description: "Run linter on changed files"
    implementation:
      type: command
      command: "ruff check {files}"
    tags: [code-quality, automated]
    metrics:
      invocations_30d: 342
      avg_duration_sec: 12
      success_rate: 0.98

  # Composite skill
  - id: deploy-staging
    name: "Deploy to Staging"
    type: composite
    scope: factory
    invocation: "/deploy-staging"
    description: "Full deployment pipeline to staging environment"
    implementation:
      type: workflow
      steps:
        - skill: build
        - skill: test
        - command: "kubectl apply -f k8s/staging/"
        - skill: smoke-test
    requires:
      - permission: invoke:kubernetes
      - approval: supervisor  # First deployment of day
    tags: [deployment, staging]

  # Knowledge skill
  - id: api-standards
    name: "API Design Standards"
    type: knowledge
    scope: org
    invocation: "/api-standards"
    description: "Our REST API design conventions"
    implementation:
      type: prompt
      content: |
        Our API standards:
        1. Use plural nouns for resources (/users, not /user)
        2. Version in URL (/v1/users)
        3. Use HTTP verbs correctly (GET=read, POST=create, etc.)
        4. Return 201 for creates, 204 for deletes
        5. Use consistent error format...
    tags: [knowledge, api, standards]
```

### Skill Matching

How do we match work orders to required skills?

```python
# vermas/skills/matching.py

class SkillMatcher:
    """Match work orders to required skills."""

    def analyze_work_order(self, wo: WorkOrder) -> SkillRequirements:
        """Determine what skills a work order needs."""
        requirements = SkillRequirements()

        # 1. Explicit skill tags
        for tag in wo.tags:
            if tag.startswith("needs:"):
                skill_id = tag.split(":")[1]
                requirements.add(skill_id, source="explicit_tag")

        # 2. Work type inference
        type_skill_map = {
            "security": ["security-review", "dependency-check"],
            "api": ["api-standards", "openapi-validation"],
            "database": ["migration-check", "query-optimization"],
            "deployment": ["deploy-staging", "rollback-plan"],
        }
        for wo_type in wo.types:
            if wo_type in type_skill_map:
                for skill in type_skill_map[wo_type]:
                    requirements.add(skill, source="type_inference")

        # 3. Content analysis (LLM-based)
        if wo.description:
            inferred = self.llm_analyze_skills(wo.description)
            for skill in inferred:
                requirements.add(skill, source="content_analysis")

        return requirements

    def find_capable_workers(self, requirements: SkillRequirements) -> List[Worker]:
        """Find workers who have the required skills."""
        candidates = []
        for worker in self.all_workers():
            coverage = self.skill_coverage(worker, requirements)
            if coverage >= 0.8:  # 80% skill match
                candidates.append((worker, coverage))
        return sorted(candidates, key=lambda x: x[1], reverse=True)
```

### Skill Gap Analysis

```
SKILL GAP DASHBOARD

┌─────────────────────────────────────────────────────────────────────────────┐
│ ORGANIZATION SKILL COVERAGE                               Week of Jan 6     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Skill Category          │ Coverage │ Gap    │ Recommendation               │
│ ────────────────────────┼──────────┼────────┼───────────────────────────── │
│ Core Development        │ 95%      │ 5%     │ ✓ Adequate                   │
│ Security                │ 60%      │ 40%    │ ⚠️ Hire/train or add expert  │
│ Performance             │ 45%      │ 55%    │ 🔴 Critical gap - add expert │
│ DevOps                  │ 80%      │ 20%    │ ⚠️ Consider training         │
│ Documentation           │ 70%      │ 30%    │ ⚠️ Consider skill building   │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│ RECENT SKILL REQUESTS (Not Met)                                              │
│                                                                             │
│ • kubernetes-expert: 5 requests, 0 available → Add k8s expert              │
│ • graphql-design: 3 requests, 0 available → Train or hire                  │
│ • ml-review: 2 requests, 0 available → Partner with ML team                │
│                                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│ SUGGESTED ACTIONS                                                            │
│                                                                             │
│ 1. Install 'performance-expert' extension (addresses 55% gap)               │
│ 2. Create '/security-review' skill from existing tools                      │
│ 3. Train 2 workers on kubernetes (reduces devops gap to 5%)                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Iteration 6 Key Insights

1. **Three skill types**: Atomic (single ops), Composite (workflows), Knowledge (prompts)

2. **Registry tracks all skills**: With metrics, scope, and implementation

3. **Skill matching is multi-source**: Explicit tags, type inference, content analysis

4. **Gap analysis guides investment**: Shows where to add experts or training

5. **Skills have scopes**: Worker, Factory, Organization

---

## Iteration 7: Learning from Extensions

### What Can Be Learned?

```
LEARNING OPPORTUNITIES

┌─────────────────────────────────────────────────────────────────────────────┐
│ FROM SKILLS                                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│ • Which skills are most used?                                               │
│ • Which skills have highest success rate?                                   │
│ • Which skill sequences work well together?                                 │
│ • What new skills are being requested but don't exist?                      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ FROM EXPERTS                                                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│ • What issues do experts commonly find?                                     │
│ • Which recommendations are most followed?                                  │
│ • What patterns could become automated checks?                              │
│ • Which expert advice could become skills?                                  │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ FROM PLUGINS                                                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│ • Which tools are most valuable?                                            │
│ • What tool combinations are common?                                        │
│ • Which tool outputs need post-processing?                                  │
│ • What new tools are needed?                                                │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Learning Pipeline

```
EXTENSION LEARNING FLOW

┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. COLLECT                                                                   │
│    Gather usage data, outcomes, feedback                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Usage Events:                                                             │
│   • skill.invoked                                                           │
│   • expert.consulted                                                        │
│   • plugin.tool_called                                                      │
│                                                                             │
│   Outcome Events:                                                           │
│   • skill.completed / skill.failed                                          │
│   • expert.recommendation_followed / expert.recommendation_ignored          │
│   • verification.passed / verification.failed                               │
│                                                                             │
│   Feedback Events:                                                          │
│   • worker.rated_skill                                                      │
│   • worker.suggested_improvement                                            │
│   • supervisor.endorsed_skill                                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. ANALYZE                                                                   │
│    Identify patterns, anomalies, opportunities                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Pattern Detection:                                                        │
│   • "Workers always run /lint before /test" → Create composite skill       │
│   • "Security expert flags X pattern 80% of the time" → Automate check     │
│   • "Skill Y fails 40% of the time after skill X" → Add dependency         │
│                                                                             │
│   Anomaly Detection:                                                        │
│   • "Skill success rate dropped 20% this week" → Investigate               │
│   • "Expert response time increased" → Check capacity                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. PROPOSE                                                                   │
│    Suggest improvements                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Improvement Proposals:                                                    │
│   • "Create composite skill '/lint-test' (saves 2min per invocation)"      │
│   • "Add automated check for SQL injection (expert finds in 60% reviews)" │
│   • "Retire skill '/old-deploy' (0 uses in 30 days)"                       │
│                                                                             │
│   Approval Required:                                                        │
│   • New skill creation: Supervisor                                          │
│   • Skill modification: Skill owner                                         │
│   • Skill retirement: CEO                                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. APPLY                                                                     │
│    Implement approved improvements                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Application Methods:                                                      │
│   • Auto-generate skill definition                                          │
│   • Add check to verification pipeline                                      │
│   • Update skill parameters                                                 │
│   • Retire/archive unused skills                                            │
│                                                                             │
│   Rollout:                                                                  │
│   • A/B test new skills                                                     │
│   • Gradual rollout (10% → 50% → 100%)                                     │
│   • Monitor success rate                                                    │
│   • Rollback if metrics degrade                                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Skill Evolution

```yaml
# Example: Skill evolved from expert feedback

# Before: Expert manually checks for this
expert_finding:
  pattern: "SQL queries built with string concatenation"
  frequency: "Found in 60% of database-related reviews"
  severity: high
  recommendation: "Use parameterized queries"

# After: Automated skill created
skill:
  id: sql-injection-check
  type: atomic
  source: learned_from_expert
  learned_from:
    expert: security-expert
    finding_count: 47
    accuracy: 0.95
  implementation:
    type: command
    command: "semgrep --config=p/sql-injection {files}"
  evolution:
    created: 2026-01-07
    version: 1.0
    next_review: 2026-04-07
```

---

## Iteration 7 Key Insights

1. **Learn from usage, outcomes, and feedback**: Three data sources

2. **Pattern detection drives improvement**: Identify common sequences, failures, gaps

3. **Proposals require approval**: Different levels for create/modify/retire

4. **Gradual rollout with monitoring**: A/B test, staged rollout, rollback capability

5. **Expert knowledge can become automated checks**: Turn repeated findings into skills

---

## Iteration 8: Scoped Learning

### Learning Scopes

```
WHO LEARNS FROM WHAT?

┌─────────────────────────────────────────────────────────────────────────────┐
│ INDIVIDUAL (Worker-level)                                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ What: Personal shortcuts, preferences, common patterns                      │
│ Storage: Worker's profile / personal config                                 │
│ Sharing: Not shared by default                                              │
│                                                                             │
│ Examples:                                                                   │
│ • "I always run /format after editing Python files"                        │
│ • "My preferred test command is 'pytest -x'"                               │
│ • "I like verbose output from security scans"                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ TEAM (Factory-level)                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ What: Team conventions, project-specific skills, shared workflows           │
│ Storage: Factory's .work/skills/                                            │
│ Sharing: Shared within factory, can be promoted to org                      │
│                                                                             │
│ Examples:                                                                   │
│ • "Our team's deploy process includes extra smoke tests"                   │
│ • "This project requires HIPAA compliance checks"                          │
│ • "We use a specific branching strategy"                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ ORGANIZATION                                                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ What: Company standards, cross-team skills, official procedures             │
│ Storage: Organization's .work/skills/                                       │
│ Sharing: Available to all factories, promoted from factory learnings        │
│                                                                             │
│ Examples:                                                                   │
│ • "Our security review process (mandatory for all)"                        │
│ • "Company-wide code style guide"                                          │
│ • "Standard deployment pipeline"                                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ ECOSYSTEM (Cross-organization)                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ What: Community skills, open-source best practices, industry standards      │
│ Storage: Public registry                                                    │
│ Sharing: Opt-in publishing, curated by registry maintainers                 │
│                                                                             │
│ Examples:                                                                   │
│ • "OWASP security checklist"                                               │
│ • "Kubernetes deployment best practices"                                    │
│ • "PCI compliance workflow"                                                 │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Knowledge Promotion

How does learning propagate upward?

```
KNOWLEDGE PROMOTION FLOW

Individual → Team → Org → Ecosystem

┌─────────────────────────────────────────────────────────────────────────────┐
│ PROMOTION TRIGGERS                                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Individual → Team:                                                          │
│ • Worker shares skill with team (explicit)                                  │
│ • Supervisor sees worker pattern, adopts for team                           │
│ • Multiple workers independently create similar skills                      │
│                                                                             │
│ Team → Org:                                                                 │
│ • Skill used successfully by 3+ factories                                   │
│ • CEO mandates skill as org standard                                        │
│ • Skill addresses org-wide need (compliance, etc.)                          │
│                                                                             │
│ Org → Ecosystem:                                                            │
│ • Organization opts to publish                                              │
│ • Skill is generalized (remove org-specific parts)                          │
│ • Registry maintainers accept and curate                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

PROMOTION PROCESS:

1. Candidate identified (manually or automatically)
2. Skill generalized if needed
3. Approval obtained (team lead, CEO, or registry)
4. Skill copied to higher scope
5. Original skill can reference promoted version
6. Metrics tracked at new scope
```

### Privacy & Isolation

```yaml
# .work/governance/learning-privacy.yaml

privacy:
  # What can be shared externally
  ecosystem_sharing:
    allowed:
      - skill_definitions  # The skill itself
      - aggregate_metrics  # Usage counts, success rates
      - anonymized_patterns  # Common sequences

    forbidden:
      - work_order_content  # Actual work being done
      - code  # Source code
      - identities  # Worker/agent identities
      - business_logic  # Proprietary processes

  # What stays within org
  org_internal:
    - detailed_metrics
    - worker_performance
    - project_names
    - customer_data

  # Opt-in for ecosystem contribution
  ecosystem_contribution:
    enabled: true
    require_review: true  # Human reviews before publishing
    anonymize: true
```

---

## Iteration 8 Key Insights

1. **Four scopes**: Individual, Team, Organization, Ecosystem

2. **Promotion flows upward**: Good patterns bubble up through scopes

3. **Privacy is enforced**: Clear boundaries on what can be shared

4. **Generalization required for promotion**: Remove org-specific details

5. **Opt-in for ecosystem**: Organizations choose what to publish

---

## Iteration 9: Extension Lifecycle & Trust

### Extension Lifecycle

```
EXTENSION LIFECYCLE STATES

┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│  DRAFT → TESTING → ACTIVE → DEPRECATED → RETIRED                           │
│    │        │        │          │           │                               │
│    │        │        │          │           │                               │
│    ▼        ▼        ▼          ▼           ▼                               │
│  Local   Staged   Production  Warning    Removed                            │
│  dev     rollout  available   period     from                               │
│  only    (10%)    to all      (migrate)  registry                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

LIFECYCLE TRANSITIONS:

DRAFT → TESTING:
  Trigger: Author marks ready
  Requirements: Manifest valid, tests pass, permissions declared
  Approval: None (author decision)

TESTING → ACTIVE:
  Trigger: Testing period complete (7 days) OR manual promotion
  Requirements:
    - Success rate > 95%
    - No security issues found
    - Positive feedback from testers
  Approval: Supervisor (factory) or CEO (org)

ACTIVE → DEPRECATED:
  Trigger: Better alternative exists, or critical flaw found
  Requirements:
    - Replacement identified (if applicable)
    - Migration path documented
    - Warning period set (default: 30 days)
  Approval: CEO or Human

DEPRECATED → RETIRED:
  Trigger: Warning period expired
  Requirements:
    - All users migrated
    - No active invocations in last 7 days
  Approval: Automatic
```

### Trust Model

```
EXTENSION TRUST LEVELS

┌─────────────────────────────────────────────────────────────────────────────┐
│ LEVEL 0: UNTRUSTED (default for unknown sources)                             │
├─────────────────────────────────────────────────────────────────────────────┤
│ Permissions: None                                                           │
│ Execution: Sandboxed, no network, no file writes                           │
│ Installation: Requires human approval                                       │
│                                                                             │
│ How to elevate: Security review + human approval                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ LEVEL 1: SAFE (verified safe extensions)                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│ Permissions: read:code, read:config, write:reports                          │
│ Execution: Sandboxed with limited file access                               │
│ Installation: Supervisor approval                                           │
│                                                                             │
│ How to elevate: Track record + elevated permission request                  │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ LEVEL 2: STANDARD (common working extensions)                                │
├─────────────────────────────────────────────────────────────────────────────┤
│ Permissions: Level 1 + invoke:git, invoke:network, read:work_orders         │
│ Execution: Container with network access                                    │
│ Installation: CEO approval                                                  │
│                                                                             │
│ How to elevate: Extended track record + security audit                      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ LEVEL 3: ELEVATED (powerful extensions)                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│ Permissions: Level 2 + write:code, block:merge, invoke:llm                  │
│ Execution: Full container access                                            │
│ Installation: Human approval                                                │
│                                                                             │
│ How to elevate: Full security audit + ongoing monitoring                    │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ LEVEL 4: PRIVILEGED (system-level extensions)                                │
├─────────────────────────────────────────────────────────────────────────────┤
│ Permissions: All (including invoke:shell, read:secrets)                     │
│ Execution: Host access                                                      │
│ Installation: Human approval + security contract                            │
│                                                                             │
│ Reserved for: Core infrastructure, critical integrations                    │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Security Verification

```python
# vermas/extensions/security.py

class ExtensionSecurityChecker:
    """Verify extension security before installation."""

    def verify(self, extension: Extension) -> SecurityVerification:
        results = SecurityVerification()

        # 1. Static analysis of extension code
        results.static_analysis = self.static_analyze(extension.code_path)

        # 2. Permission analysis
        results.permission_analysis = self.analyze_permissions(
            requested=extension.manifest.permissions,
            code_uses=self.detect_permission_usage(extension.code_path)
        )

        # 3. Dependency check
        results.dependency_check = self.check_dependencies(
            extension.dependencies
        )

        # 4. Sandbox test
        results.sandbox_test = self.run_in_sandbox(extension)

        # 5. Behavior analysis
        results.behavior_analysis = self.analyze_behavior(
            extension, test_scenarios=self.standard_scenarios
        )

        # Calculate trust score
        results.trust_score = self.calculate_trust_score(results)
        results.recommended_level = self.recommend_trust_level(results)

        return results

    def recommend_trust_level(self, results: SecurityVerification) -> int:
        if results.trust_score < 0.5:
            return 0  # Untrusted
        elif results.trust_score < 0.7:
            return 1  # Safe
        elif results.trust_score < 0.85:
            return 2  # Standard
        elif results.trust_score < 0.95:
            return 3  # Elevated
        else:
            return 4  # Privileged (still needs human review)
```

### Extension Audit Trail

```yaml
# Every extension action is logged

extension.installed:
  extension_id: security-scanner
  version: 1.2.3
  installed_by: ceo
  scope: org
  trust_level: 2
  permissions_granted: [read:code, read:config, write:reports, invoke:git]
  approval_chain: [supervisor-alpha, ceo]

extension.invoked:
  extension_id: security-scanner
  tool: scan_code
  invoker: worker-1
  work_order: wo-abc123
  duration_ms: 2340
  result: success

extension.permission_used:
  extension_id: security-scanner
  permission: invoke:git
  action: "git diff HEAD~1"
  context: "Scanning changed files"

extension.deprecated:
  extension_id: old-linter
  deprecated_by: ceo
  reason: "Replaced by new-linter with better performance"
  replacement: new-linter
  warning_period_days: 30
  migration_guide: "https://..."
```

---

## Iteration 9 Key Insights

1. **Five lifecycle states**: Draft → Testing → Active → Deprecated → Retired

2. **Five trust levels**: Untrusted → Safe → Standard → Elevated → Privileged

3. **Security verification is multi-layered**: Static analysis, sandboxing, behavior analysis

4. **Approval scales with risk**: More permissions = more approval required

5. **Full audit trail**: Every install, invoke, and permission use is logged

---

## Iteration 10: Synthesis

### Priority + Extensibility: The Complete Model

```
THE UNIFIED MODEL

┌─────────────────────────────────────────────────────────────────────────────┐
│                           PRIORITY SYSTEM                                    │
│               "What needs to be done and when"                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Dynamic Priority Score = Base + Urgency Factors + Importance Factors      │
│                                                                             │
│   Urgency: Deadline, Blockers, External Pressure, Age, Escalation           │
│   Importance: Objectives, Business Value, Risk, CEO Directive               │
│                                                                             │
│   Classes: CRITICAL → HIGH → MEDIUM → LOW → BACKLOG                         │
│                                                                             │
│   Decay: Old work rises, ancient work reviewed                              │
│                                                                             │
│   Controls: Budget (max % per class), Justification, Expiration             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                    "Who can do this work?"
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         EXTENSIBILITY SYSTEM                                 │
│              "Bringing in the right capabilities"                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Extension Types:                                                          │
│   ├── Plugins (tools from external services)                                │
│   ├── Skills (reusable workflows and procedures)                            │
│   ├── Experts (specialized agent consultants)                               │
│   └── Templates (packaged best practices)                                   │
│                                                                             │
│   Integration Patterns:                                                     │
│   ├── Consultation (on-demand advice)                                       │
│   ├── Review Gate (mandatory checkpoint)                                    │
│   ├── Pair Work (collaborative session)                                     │
│   └── Delegation (handoff to specialist)                                    │
│                                                                             │
│   Trust Levels: Untrusted → Safe → Standard → Elevated → Privileged         │
│                                                                             │
│   Lifecycle: Draft → Testing → Active → Deprecated → Retired                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                    "How do we get better?"
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          LEARNING SYSTEM                                     │
│              "Improving from experience"                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   Learning Sources: Usage → Outcomes → Feedback                             │
│                                                                             │
│   Learning Scopes:                                                          │
│   ├── Individual (personal shortcuts)                                       │
│   ├── Team (factory-specific skills)                                        │
│   ├── Organization (company standards)                                      │
│   └── Ecosystem (community knowledge)                                       │
│                                                                             │
│   Knowledge Flow: Individual → Team → Org → Ecosystem                       │
│                                                                             │
│   Promotion: Pattern detection → Proposal → Approval → Rollout              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### File Structure

```
.work/
├── priority/
│   ├── config.yaml           # Priority calculation weights
│   ├── classes.yaml          # Priority class definitions
│   ├── decay-rules.yaml      # Staleness and decay rules
│   └── budgets.yaml          # Priority class budgets
│
├── extensions/
│   ├── installed/            # Installed extensions
│   │   ├── security-scanner/
│   │   │   ├── manifest.yaml
│   │   │   └── config.yaml
│   │   └── code-quality/
│   │       └── ...
│   ├── registry.yaml         # Registry configuration
│   └── trust-levels.yaml     # Trust level definitions
│
├── skills/
│   ├── atomic/               # Single-operation skills
│   │   ├── lint.yaml
│   │   └── test.yaml
│   ├── composite/            # Multi-step workflows
│   │   ├── deploy-staging.yaml
│   │   └── release.yaml
│   ├── knowledge/            # Domain knowledge prompts
│   │   └── api-standards.yaml
│   └── registry.yaml         # Skill registry
│
├── experts/
│   ├── security-expert/
│   │   ├── profile.yaml
│   │   └── system-prompt.md
│   ├── performance-expert/
│   │   └── ...
│   ├── budget.yaml           # Expert time budgets
│   └── reputation.yaml       # Reputation scores
│
└── learning/
    ├── patterns.jsonl        # Detected patterns
    ├── proposals.jsonl       # Improvement proposals
    ├── privacy.yaml          # Sharing rules
    └── promotions.jsonl      # Promotion history
```

### Key Events

```yaml
# Priority events
priority.calculated:
  work_order_id: string
  score: float
  class: string
  factors: dict

priority.class_changed:
  work_order_id: string
  from_class: string
  to_class: string
  reason: string

priority.exhaustion_detected:
  critical_percent: float
  recommendation: string

# Extension events
extension.installed:
  extension_id: string
  version: string
  scope: string
  trust_level: int

extension.invoked:
  extension_id: string
  tool: string
  invoker: string
  result: string

extension.promoted:
  extension_id: string
  from_scope: string
  to_scope: string

# Learning events
learning.pattern_detected:
  pattern_type: string
  confidence: float
  evidence: list

learning.proposal_created:
  proposal_id: string
  type: string
  description: string

learning.knowledge_promoted:
  skill_id: string
  from_scope: string
  to_scope: string
```

---

## Summary: What We Designed

### Priority System (Iterations 1-3)
1. **Dynamic priority scoring** combining urgency and importance factors
2. **Five priority classes** with SLAs and behaviors
3. **Decay mechanisms** to prevent eternal deferral
4. **Exhaustion detection** and reset protocols
5. **Explainable priority** with factor breakdowns

### Extensibility System (Iterations 4-6)
1. **Four extension types**: Plugins, Skills, Experts, Templates
2. **Scoped installation** with approval chains
3. **Permission model** with trust levels
4. **Expert integration patterns**: Consultation, Gate, Pair, Delegation
5. **Skill registry** with matching and gap analysis

### Learning System (Iterations 7-9)
1. **Multi-source learning**: Usage, outcomes, feedback
2. **Four scopes**: Individual, Team, Org, Ecosystem
3. **Knowledge promotion** with privacy controls
4. **Extension lifecycle** from Draft to Retired
5. **Trust verification** with security analysis

---

## Approval Status

| Section | Status |
|---------|--------|
| Iteration 1: Priority Fundamentals | Pending Review |
| Iteration 2: Priority Signals | Pending Review |
| Iteration 3: Priority Decay | Pending Review |
| Iteration 4: Extensibility Model | Pending Review |
| Iteration 5: Expert Integration | Pending Review |
| Iteration 6: Skills Registry | Pending Review |
| Iteration 7: Learning from Extensions | Pending Review |
| Iteration 8: Scoped Learning | Pending Review |
| Iteration 9: Extension Lifecycle | Pending Review |
| Iteration 10: Synthesis | Pending Review |

