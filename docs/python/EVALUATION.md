# VerMAS Evaluation

> How to evaluate if the system is working correctly

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
- Count bead status transitions
- Track REWORK_REQUEST messages
- Count abandoned beads (killed polecats without completion)

### GUPP Compliance

**What it measures:** Do agents execute immediately when hooked?

| Metric | Formula | Target |
|--------|---------|--------|
| Hook response time | time(hook_created → work_started) | <30 sec |
| GUPP violations | agents_that_waited_for_confirmation | 0 |

**How to measure:**
- Timestamp hook file creation
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
| Deacon uptime | running_time / total_time | >99% |
| Witness uptime | running_time / total_time | >99% |
| Refinery uptime | running_time / total_time | >99% |

**How to measure:**
- Track session start/stop times
- Count Deacon restarts
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
| Hook persistence | hooks_surviving_crash / hooks_at_crash | 100% |
| Sandbox recovery | sandboxes_with_recoverable_work / killed_polecats | >80% |
| Handoff success | handoffs_continued / handoffs_created | >90% |

**How to measure:**
- Kill sessions, verify hooks persist
- Check sandbox state after polecat kill
- Track handoff message → next session action

---

## 3. Efficiency Metrics

### Throughput

**What it measures:** How much work gets done?

| Metric | Formula | Target |
|--------|---------|--------|
| Beads/hour | completed_beads / elapsed_hours | Baseline + 20% |
| Parallel utilization | avg_active_polecats / max_slots | >60% |
| Queue wait time | avg(queued → started) | <10 min |

**How to measure:**
- Count bead completions over time
- Sample polecat slot utilization
- Timestamp queue events

### Resource Usage

**What it measures:** How efficiently are resources used?

| Metric | Formula | Target |
|--------|---------|--------|
| Slot utilization | active_slots / total_slots | >60% |
| Idle time | polecat_idle_time / polecat_total_time | <20% |
| Claude calls/bead | llm_invocations / completed_beads | <50 |

**How to measure:**
- Track slot allocation/release
- Measure time between polecat actions
- Count Claude CLI invocations per session

### Pipeline Latency

**What it measures:** How long does the full cycle take?

| Metric | Formula | Target |
|--------|---------|--------|
| Time to merge | avg(bead_created → merged) | Baseline |
| Verification time | avg(MERGE_READY → verdict) | <5 min |
| Review cycle time | avg(REWORK_REQUEST → re-submit) | <30 min |

**How to measure:**
- Timestamp bead lifecycle events
- Track verification molecule duration
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

**Setup:** Simple bead, no issues

**Expected:**
1. Bead created → hooked within 30 sec
2. Polecat executes immediately (GUPP)
3. Work completed → POLECAT_DONE sent
4. Witness forwards → MERGE_READY
5. Refinery runs tests → PASS
6. Verification runs → PASS
7. Merge completes → MERGED sent

**Success criteria:**
- End-to-end in <15 min (excluding actual coding time)
- No manual intervention
- All messages delivered

### Scenario 2: Verification Failure

**Setup:** Bead with implementation bug

**Expected:**
1. Polecat completes work
2. Tests pass (bug not caught by tests)
3. Verification runs → FAIL (Critic finds issue)
4. REWORK_REQUEST sent
5. Polecat fixes issue
6. Second attempt → PASS
7. Merge completes

**Success criteria:**
- Bug caught by verification
- Rework request contains actionable feedback
- Second attempt succeeds

### Scenario 3: Polecat Failure

**Setup:** Polecat gets stuck

**Expected:**
1. Polecat goes idle >5 min
2. Witness sends NUDGE
3. Polecat still idle >15 min
4. Witness kills session, releases slot
5. Work remains in sandbox
6. New polecat can continue
7. Work eventually completes

**Success criteria:**
- Detection within patrol interval
- Slot released back to pool
- Work recoverable

### Scenario 4: Watchdog Chain

**Setup:** Kill Witness process

**Expected:**
1. Witness dies
2. Deacon detects (within 60 sec)
3. Deacon restarts Witness
4. Witness resumes patrol
5. No polecats lost

**Success criteria:**
- Detection within 2 patrol intervals
- Auto-restart succeeds
- No human intervention

### Scenario 5: Handoff Continuity

**Setup:** Session ends mid-work

**Expected:**
1. Agent creates handoff mail
2. Hooks handoff for next session
3. Session ends
4. New session starts
5. New session finds hooked handoff
6. Work continues from context

**Success criteria:**
- Handoff persists across sessions
- Next session executes handoff (GUPP)
- Context sufficient to continue

---

## Evaluation Methods

### Automated Monitoring

**What:** Continuous metrics collection

**Implementation:**
- Instrument message sending/receiving
- Track bead state transitions
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
- Create conflicting beads
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
| Completion time | Run 10 similar beads, average |
| Verification accuracy | Human review 20 verdicts |
| Recovery time | Kill 5 sessions, measure recovery |
| Message latency | Timestamp 50 messages end-to-end |

### Comparison Points

Compare against:
- Manual development (same tasks, human only)
- Single-agent Claude Code (no Gas Town)
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
- Beads completed: X
- Beads failed: Y
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
- Throughput: Y beads (↑/↓ from last week)
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

## See Also

- [ARCHITECTURE.md](./ARCHITECTURE.md) - System architecture
- [AGENTS.md](./AGENTS.md) - Agent roles
- [WORKFLOWS.md](./WORKFLOWS.md) - Molecule system
- [MESSAGING.md](./MESSAGING.md) - Communication patterns
