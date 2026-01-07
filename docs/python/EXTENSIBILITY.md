# VerMAS Extensibility Design

> Bringing in external expertise through plugins, skills, experts, and templates

---

## Initial Implementation Scope (Approved)

Start with a constrained, safe model. Expand later based on experience.

### Core Principles

1. **Anthropic-approved only**: Plugins must come from Anthropic's curated registry
2. **Ephemeral by default**: Install for the work, remove when done
3. **Autonomous at lower scopes**: Worker/Team can self-serve
4. **Human gates at higher scopes**: Department/Company requires human approval

### Installation Scope Matrix

```
PLUGIN INSTALLATION AUTHORITY

┌──────────────┬─────────────────┬─────────────────┬─────────────────────────┐
│ Scope        │ Who Approves    │ Lifetime        │ Source Restriction      │
├──────────────┼─────────────────┼─────────────────┼─────────────────────────┤
│ Worker       │ Worker (self)   │ Until work done │ Anthropic registry only │
│ Team/Factory │ Supervisor      │ Until work done │ Anthropic registry only │
│ Department   │ Human required  │ Persistent      │ Any approved source     │
│ Company/Org  │ Human required  │ Persistent      │ Any approved source     │
└──────────────┴─────────────────┴─────────────────┴─────────────────────────┘
```

### Ephemeral Plugin Lifecycle

```
WORK-SCOPED PLUGIN LIFECYCLE

   Work Order Created
          │
          ▼
   ┌─────────────────┐
   │ Analyze work    │  ← What capabilities needed?
   │ requirements    │
   └────────┬────────┘
            │
            ▼
   ┌─────────────────┐
   │ Search Anthropic│  ← Only approved plugins
   │ registry        │
   └────────┬────────┘
            │
            ▼
   ┌─────────────────┐
   │ Install plugin  │  ← Scoped to worker/team
   │ (ephemeral)     │
   └────────┬────────┘
            │
            ▼
   ┌─────────────────┐
   │ Use plugin      │  ← During work execution
   │ tools           │
   └────────┬────────┘
            │
            ▼
   ┌─────────────────┐
   │ Work completed  │
   └────────┬────────┘
            │
            ▼
   ┌─────────────────┐
   │ Uninstall       │  ← Automatic cleanup
   │ plugin          │
   └─────────────────┘
```

### Why This Model

**Trust**: Anthropic-approved plugins are vetted for safety and quality.

**Simplicity**: Workers don't need to evaluate arbitrary extensions.

**Clean state**: Ephemeral installation prevents plugin sprawl.

**Escalation path**: If a plugin is needed frequently, promote to Department/Company level with human approval.

### Example Flow

```
Worker receives work order: "Add OAuth to user service"

1. Worker analyzes: needs OAuth expertise, security review
2. Worker searches Anthropic registry: finds "oauth-helper" plugin
3. Worker installs oauth-helper (worker scope, ephemeral)
4. Worker uses oauth-helper tools during implementation
5. Work order completed
6. oauth-helper automatically uninstalled

Later: Team notices oauth-helper used 10 times this month
→ Supervisor requests human approval for Team-level persistent install
→ Human approves
→ Plugin now available to all team workers without per-work install
```

---

## Future Expansion (Not Initial Scope)

The iterations below explore the full extensibility model for future implementation:
- Skills (workflow automation)
- Experts (specialized agents)
- Templates (best practice packages)
- Learning and knowledge promotion

---

## The Questions

1. **How do we bring in external expertise?** - Consultants, specialists, plugins
2. **What's the extension model?** - Tools, commands, sub-agents, skills
3. **How do extensions integrate?** - Discovery, installation, invocation
4. **Who learns from extensions?** - Individual, team, department, ecosystem
5. **How do we trust extensions?** - Verification, sandboxing, audit
6. **How do extensions evolve?** - Lifecycle, deprecation, replacement

---

## Iteration 1: Extensibility Model

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

## Iteration 1 Key Insights

1. **Four extension types**: Plugins (tools), Skills (workflows), Experts (agents), Templates (packages)

2. **Scoped installation**: Worker, Factory, Organization levels

3. **Permission model**: Read/Write/Invoke categories with approval levels

4. **Discovery via registry**: Official, private, git, local sources

5. **Claude Code parallel**: MCP Servers → Plugins, Skills → Skills, Sub-agents → Experts

---

## Iteration 2: Expert/Consultant Integration

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

## Iteration 2 Key Insights

1. **Four integration patterns**: Consultation, Review Gate, Pair Work, Delegation

2. **Budget controls**: Monthly limits, rate limiting, priority access

3. **Reputation is multi-dimensional**: Accuracy, speed, quality, satisfaction, false positives

4. **Experts are agents with specialized prompts**: Configuration defines expertise

5. **Request flow includes availability check**: Queue when busy, reject when over budget

---

## Iteration 3: Skills & Capabilities Registry

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

## Iteration 3 Key Insights

1. **Three skill types**: Atomic (single ops), Composite (workflows), Knowledge (prompts)

2. **Registry tracks all skills**: With metrics, scope, and implementation

3. **Skill matching is multi-source**: Explicit tags, type inference, content analysis

4. **Gap analysis guides investment**: Shows where to add experts or training

5. **Skills have scopes**: Worker, Factory, Organization

---

## Iteration 4: Learning from Extensions

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

## Iteration 4 Key Insights

1. **Learn from usage, outcomes, and feedback**: Three data sources

2. **Pattern detection drives improvement**: Identify common sequences, failures, gaps

3. **Proposals require approval**: Different levels for create/modify/retire

4. **Gradual rollout with monitoring**: A/B test, staged rollout, rollback capability

5. **Expert knowledge can become automated checks**: Turn repeated findings into skills

---

## Iteration 5: Scoped Learning

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

## Iteration 5 Key Insights

1. **Four scopes**: Individual, Team, Organization, Ecosystem

2. **Promotion flows upward**: Good patterns bubble up through scopes

3. **Privacy is enforced**: Clear boundaries on what can be shared

4. **Generalization required for promotion**: Remove org-specific details

5. **Opt-in for ecosystem**: Organizations choose what to publish

---

## Iteration 6: Extension Lifecycle & Trust

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

## Iteration 6 Key Insights

1. **Five lifecycle states**: Draft → Testing → Active → Deprecated → Retired

2. **Five trust levels**: Untrusted → Safe → Standard → Elevated → Privileged

3. **Security verification is multi-layered**: Static analysis, sandboxing, behavior analysis

4. **Approval scales with risk**: More permissions = more approval required

5. **Full audit trail**: Every install, invoke, and permission use is logged

---

## Summary: The Extensibility Model

### Extension Types
1. **Plugins**: External tools (like MCP Servers) providing specialized capabilities
2. **Skills**: Named, reusable workflows and procedures
3. **Experts**: Specialized agent configurations for domain expertise
4. **Templates**: Pre-packaged best practices and configurations

### Integration Patterns
1. **Consultation**: On-demand advice from experts
2. **Review Gate**: Mandatory checkpoints requiring expert approval
3. **Pair Work**: Collaborative sessions with expert guidance
4. **Delegation**: Handoff to specialists for complex work

### Trust & Security
- **Five trust levels**: Untrusted → Safe → Standard → Elevated → Privileged
- **Permission model**: Read/Write/Invoke categories with approval chains
- **Lifecycle management**: Draft → Testing → Active → Deprecated → Retired
- **Full audit trail**: Every action logged for compliance

### Learning System
- **Four scopes**: Individual → Team → Organization → Ecosystem
- **Knowledge promotion**: Good patterns bubble up through scopes
- **Privacy controls**: Clear boundaries on what can be shared
- **Pattern detection**: Turn repeated expert findings into automated skills

### File Structure

```
.work/
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

---

## Approval Status

| Section | Status |
|---------|--------|
| Iteration 1: Extensibility Model | Pending Review |
| Iteration 2: Expert Integration | Pending Review |
| Iteration 3: Skills Registry | Pending Review |
| Iteration 4: Learning from Extensions | Pending Review |
| Iteration 5: Scoped Learning | Pending Review |
| Iteration 6: Extension Lifecycle | Pending Review |

