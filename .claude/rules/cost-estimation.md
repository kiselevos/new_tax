# Cost Estimation Before Delegation

> Derived from COST_TRACKING.md and team agreement (2026-03-11).

## Rule

Before delegating any task to a subagent, you MUST show a cost estimate to the user.

## Format

```
→ agent-name (model) · ~NNK tokens · ~$X-Y
```

Example:
```
→ CTO (opus) · ~50K tokens · ~$4-5
→ frontend-developer (sonnet) · ~80K tokens · ~$1.5-2
→ Total: ~$5.5-7
```

## When to estimate

- **ALWAYS**: before launching any Agent tool call
- **Multiple agents**: show each + total
- **Background agents**: estimate before launch, not after

## How to estimate

| Model | Per 50K tokens (in+out) | Typical task |
|-------|------------------------|--------------|
| Opus | $3-5 | Strategy, triage, decomposition |
| Sonnet | $0.5-1.5 | Code, review, architecture |
| Haiku | $0.01-0.10 | Search, docs, validation |

Factors that increase cost:
- Many tool calls (file reads, edits) → more tokens
- Large codebase context → more input tokens
- Multiple iterations → multiply estimate

## Exceptions

No estimate needed for:
- Explore agents (always cheap, haiku-level)
- Direct tool calls (Read, Edit, Grep — not agents)
- Resuming an agent (cost already committed)

## Reference

Список агентов проекта: `.claude/agents/`
