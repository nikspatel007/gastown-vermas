# VerMAS Evaluation

> How to evaluate if the system is working correctly

## Event-Driven Evaluation

All evaluation in VerMAS is derived from the **event log**. Every metric can be computed by querying `events.jsonl`. See [EVENTS.md](./EVENTS.md) for the event sourcing model.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      EVENT-DRIVEN EVALUATION                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   events.jsonl                                                              │
│        │                                                                    │
│        ├─→ Correctness metrics  (completion rates, assignment compliance)  │
│        │                                                                    │
│        ├─→ Reliability metrics  (uptime, recovery times)                   │
│        │                                                                    │
│        ├─→ Efficiency metrics   (throughput, latency)                      │
│        │                                                                    │
│        └─→ Verification metrics (accuracy, false positive rate)            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Evaluation Dimensions

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          EVALUATION FRAMEWORK                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │  1. CORRECTNESS                                                     │  │
│   │     Does the system do what it's supposed to?                       │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │  2. RELIABILITY                                                     │  │
│   │     Does it keep running without intervention?                      │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │  3. EFFICIENCY                                                      │  │
│   │     How quickly does work get done?                                 │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │  4. VERIFICATION QUALITY                                            │  │
│   │     Does VerMAS catch real issues? (VerMAS-specific)               │  │
│   └─────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 1. Correctness Metrics

### Work Completion Rate

**What it measures:** Do tasks get completed?

| Metric | Formula | Target |
|--------|---------|--------|
| Completion rate | completed / (completed + failed + abandoned) | >90% |
| First-pass rate | completed_first_try / completed | >70% |
| Rework rate | rework_requests / completed | <30% |

**How to measure:**
- Count work order status transitions
- Track REWORK_REQUEST messages
- Count abandoned work orders (killed workers without completion)

### Assignment Principle Compliance

**What it measures:** Do agents execute immediately when assigned?

| Metric | Formula | Target |
|--------|---------|--------|
| Assignment response time | time(assignment_created → work_started) | <30 sec |
| Assignment violations | agents_that_waited_for_confirmation | 0 |

**How to measure:**
- Timestamp assignment file creation
- Timestamp first agent action after startup
- Review session logs for confirmation prompts

### Message Delivery

**What it measures:** Do messages reach their destinations?

| Metric | Formula | Target |
|--------|---------|--------|
| Delivery rate | delivered / sent | >99% |
| Processing rate | processed / delivered | >99% |
| Lost messages | sent - delivered | 0 |

**How to measure:**
- Compare sent timestamps with read timestamps
- Check for unprocessed messages in inboxes
- Audit message JSONL for orphans

---

## 2. Reliability Metrics

### Uptime

**What it measures:** Is the system running?

| Metric | Formula | Target |
|--------|---------|--------|
| Operations uptime | running_time / total_time | >99% |
| Supervisor uptime | running_time / total_time | >99% |
| QA uptime | running_time / total_time | >99% |

**How to measure:**
- Track session start/stop times
- Count Operations restarts
- Monitor watchdog chain activity

### Recovery

**What it measures:** Does the system recover from failures?

| Metric | Formula | Target |
|--------|---------|--------|
| Auto-recovery rate | auto_recovered / failures | >95% |
| Mean time to recovery | avg(failure → recovered) | <5 min |
| Manual interventions | human_restarts / total_restarts | <5% |

**How to measure:**
- Track restart events
- Timestamp failure detection → recovery
- Count human interventions

### Work Continuity

**What it measures:** Does work survive failures?

| Metric | Formula | Target |
|--------|---------|--------|
| Assignment persistence | assignments_surviving_crash / assignments_at_crash | 100% |
| Workspace recovery | workspaces_with_recoverable_work / killed_workers | >80% |
| Handoff success | handoffs_continued / handoffs_created | >90% |

**How to measure:**
- Kill sessions, verify assignments persist
- Check workspace state after worker kill
- Track handoff message → next session action

---

## 3. Efficiency Metrics

### Throughput

**What it measures:** How much work gets done?

| Metric | Formula | Target |
|--------|---------|--------|
| Work orders/hour | completed_work_orders / elapsed_hours | Baseline + 20% |
| Parallel utilization | avg_active_workers / max_slots | >60% |
| Queue wait time | avg(queued → started) | <10 min |

**How to measure:**
- Count work order completions over time
- Sample worker slot utilization
- Timestamp queue events

### Resource Usage

**What it measures:** How efficiently are resources used?

| Metric | Formula | Target |
|--------|---------|--------|
| Slot utilization | active_slots / total_slots | >60% |
| Idle time | worker_idle_time / worker_total_time | <20% |
| Claude calls/work order | llm_invocations / completed_work_orders | <50 |

**How to measure:**
- Track slot allocation/release
- Measure time between worker actions
- Count Claude CLI invocations per session

### Pipeline Latency

**What it measures:** How long does the full cycle take?

| Metric | Formula | Target |
|--------|---------|--------|
| Time to merge | avg(work_order_created → merged) | Baseline |
| Verification time | avg(READY_FOR_QA → verdict) | <5 min |
| Review cycle time | avg(REWORK_REQUEST → re-submit) | <30 min |

**How to measure:**
- Timestamp work order lifecycle events
- Track verification process duration
- Measure rework turnaround

---

## 4. Verification Quality (VerMAS-specific)

### Detection Rate

**What it measures:** Does verification catch real issues?

| Metric | Formula | Target |
|--------|---------|--------|
| True positive rate | issues_caught / actual_issues | >80% |
| False positive rate | false_failures / total_failures | <10% |
| False negative rate | missed_issues / actual_issues | <20% |

**How to measure:**
- Human review of verification verdicts
- Track issues found post-merge
- Compare verification outcome with human judgment

### Criterion Quality

**What it measures:** Are the generated criteria meaningful?

| Metric | Formula | Target |
|--------|---------|--------|
| Testable criteria | shell_testable / total_criteria | >80% |
| Redundant criteria | duplicates / total_criteria | <10% |
| Coverage | requirements_with_criteria / total_requirements | >90% |

**How to measure:**
- Review Strategist output
- Check for duplicate test commands
- Map criteria back to requirements

### Adversarial Review Quality

**What it measures:** Does the Advocate/Critic/Judge process work?

| Metric | Formula | Target |
|--------|---------|--------|
| Argument quality | arguments_with_evidence / total_arguments | >90% |
| Judge reasoning | verdicts_with_reasoning / total_verdicts | 100% |
| Flip rate | judge_disagreed_with_auditor / total_reviews | 5-20% |

**How to measure:**
- Review Advocate/Critic output for evidence citations
- Check Judge output for reasoning
- Compare Auditor assessment with final verdict

---

## Test Scenarios

### Scenario 1: Happy Path

**Setup:** Simple work order, no issues

**Expected:**
1. Work order created → assigned within 30 sec
2. Worker executes immediately (Assignment Principle)
3. Work completed → WORKER_DONE sent
4. Supervisor forwards → READY_FOR_QA
5. QA runs tests → PASS
6. Verification runs → PASS
7. Merge completes → MERGED sent

**Success criteria:**
- End-to-end in <15 min (excluding actual coding time)
- No manual intervention
- All messages delivered

### Scenario 2: Verification Failure

**Setup:** Work order with implementation bug

**Expected:**
1. Worker completes work
2. Tests pass (bug not caught by tests)
3. Verification runs → FAIL (Critic finds issue)
4. REWORK_REQUEST sent
5. Worker fixes issue
6. Second attempt → PASS
7. Merge completes

**Success criteria:**
- Bug caught by verification
- Rework request contains actionable feedback
- Second attempt succeeds

### Scenario 3: Worker Failure

**Setup:** Worker gets stuck

**Expected:**
1. Worker goes idle >5 min
2. Supervisor sends NUDGE
3. Worker still idle >15 min
4. Supervisor kills session, releases slot
5. Work remains in workspace
6. New worker can continue
7. Work eventually completes

**Success criteria:**
- Detection within patrol interval
- Slot released back to pool
- Work recoverable

### Scenario 4: Watchdog Chain

**Setup:** Kill Supervisor process

**Expected:**
1. Supervisor dies
2. Operations detects (within 60 sec)
3. Operations restarts Supervisor
4. Supervisor resumes patrol
5. No workers lost

**Success criteria:**
- Detection within 2 patrol intervals
- Auto-restart succeeds
- No human intervention

### Scenario 5: Handoff Continuity

**Setup:** Session ends mid-work

**Expected:**
1. Agent creates handoff mail
2. Attaches handoff to assignment for next session
3. Session ends
4. New session starts
5. New session finds assignment with handoff
6. Work continues from context

**Success criteria:**
- Handoff persists across sessions
- Next session executes handoff (Assignment Principle)
- Context sufficient to continue

---

## Evaluation Methods

### Automated Monitoring

**What:** Continuous metrics collection

**Implementation:**
- Instrument message sending/receiving
- Track work order state transitions
- Log session start/stop events
- Record Claude CLI calls

**Output:**
- Dashboard with real-time metrics
- Alerts for anomalies
- Historical trends

### Manual Review

**What:** Human assessment of quality

**When:**
- Weekly sample review
- After significant changes
- When metrics indicate problems

**What to review:**
- Verification reasoning quality
- Agent decision appropriateness
- Message content clarity
- Prompt effectiveness

### Chaos Testing

**What:** Deliberately break things

**Tests:**
- Kill random sessions
- Corrupt message files
- Create conflicting work orders
- Flood with work

**Purpose:**
- Verify recovery mechanisms
- Find failure modes
- Test watchdog chain

---

## Benchmarks

### Baseline Establishment

Before optimization, establish baselines:

| Metric | How to establish |
|--------|------------------|
| Completion time | Run 10 similar work orders, average |
| Verification accuracy | Human review 20 verdicts |
| Recovery time | Kill 5 sessions, measure recovery |
| Message latency | Timestamp 50 messages end-to-end |

### Comparison Points

Compare against:
- Manual development (same tasks, human only)
- Single-agent Claude Code (no VerMAS)
- Previous system version

### Improvement Targets

Set targets based on:
- Current baseline + X%
- Industry benchmarks
- Theoretical limits

---

## Reporting

### Daily Report

```
VERMAS DAILY REPORT - {date}

WORK
- Work orders completed: X
- Work orders failed: Y
- Completion rate: Z%

RELIABILITY
- Uptime: X%
- Restarts: Y
- Manual interventions: Z

VERIFICATION
- Reviews run: X
- Pass rate: Y%
- Flip rate: Z%

ISSUES
- [List any anomalies]
```

### Weekly Summary

```
VERMAS WEEKLY SUMMARY - {week}

TRENDS
- Completion rate: X% (↑/↓ from last week)
- Throughput: Y work orders (↑/↓ from last week)
- Verification accuracy: Z% (based on manual review)

TOP ISSUES
1. [Issue description]
2. [Issue description]

ACTIONS TAKEN
- [What was fixed]

NEXT WEEK
- [Planned improvements]
```

---

## Computing Metrics from Events

All metrics derive from the event log. Here's how to compute them using the CLI or programmatically.

### CLI Commands

```bash
# Completion rate from events
wo eval completion --since=7d

# Assignment compliance
wo eval assignment-compliance --since=1d

# Throughput
wo eval throughput --since=24h

# Verification accuracy (requires human labels)
wo eval verify-accuracy --since=30d

# Full evaluation report
wo eval report --since=7d --output=report.json
```

### Programmatic Computation

```python
from vermas.eval import EventMetrics
from datetime import timedelta

metrics = EventMetrics(".work/events.jsonl")

# Completion rate
created = metrics.count("work_order.created", since=timedelta(days=7))
closed = metrics.count("work_order.status_changed",
                       filter={"to_status": "closed"},
                       since=timedelta(days=7))
completion_rate = closed / created if created > 0 else 0

# Assignment compliance (response time < 30s)
assignment_checks = metrics.get_events("assignment.checked", since=timedelta(days=1))
fast_responses = [e for e in assignment_checks if e.data["response_ms"] < 30000]
assignment_compliance = len(fast_responses) / len(assignment_checks)

# Average time to merge
def time_to_merge(wo_id: str) -> timedelta:
    created = metrics.get_event("work_order.created", filter={"work_order_id": wo_id})
    merged = metrics.get_event("work_order.status_changed",
                               filter={"work_order_id": wo_id, "to_status": "merged"})
    return merged.timestamp - created.timestamp

# Verification pass rate
verdicts = metrics.get_events("verify.verdict", since=timedelta(days=7))
passes = [v for v in verdicts if v.data["verdict"] == "PASS"]
pass_rate = len(passes) / len(verdicts)
```

### Event Queries for Each Metric

| Metric | Event Query |
|--------|-------------|
| Completion rate | `work_order.created` vs `work_order.status_changed(to=closed)` |
| Assignment compliance | `assignment.checked` where `response_ms < 30000` |
| Message delivery | `mail.sent` vs `mail.delivered` |
| Recovery time | `agent.stopped` to `agent.started` duration |
| Verification accuracy | `verify.verdict` cross-referenced with human labels |
| Throughput | `work_order.status_changed(to=closed)` per hour |
| Idle time | `agent.working` to `agent.idle` gaps |

### Continuous Monitoring

Set up a metrics daemon that tails the event feed:

```python
async def metrics_daemon():
    """Continuously compute and expose metrics."""
    metrics = RealTimeMetrics(".work/feed.jsonl")

    async for event in metrics.watch():
        # Update running totals
        if event.event_type == "work_order.created":
            metrics.increment("work_orders.created")
        elif event.event_type == "work_order.status_changed":
            if event.data["to_status"] == "closed":
                metrics.increment("work_orders.closed")

        # Emit metrics event every minute
        if metrics.should_snapshot():
            emit_event(Event(
                event_type="metrics.snapshot",
                data=metrics.snapshot()
            ))
```

---

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [OPERATIONS.md](./OPERATIONS.md) - Monitoring and maintenance
- [AGENTS.md](./AGENTS.md) - Agent roles
- [HOOKS.md](./HOOKS.md) - Claude Code integration and git worktrees
- [WORKFLOWS.md](./WORKFLOWS.md) - Process system
- [MESSAGING.md](./MESSAGING.md) - Communication patterns
- [EVENTS.md](./EVENTS.md) - Event sourcing and change feeds
- [VERIFICATION.md](./VERIFICATION.md) - VerMAS verification pipeline
- [SCHEMAS.md](./SCHEMAS.md) - Data specifications
- [CLI.md](./CLI.md) - Evaluation commands
