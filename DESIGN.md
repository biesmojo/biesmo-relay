# DESIGN.md - Relay Architecture Decision

## Architecture Choice: CRM-Centered Hybrid (400 words)

We chose **CRM-centered hybrid** architecture:
- **CRM tables as system of record** (customers, sessions, tickets, messages) - Lovina style.
- **Internal event bus** for rules/AI/workflow - SIGAP style.
- Rules evaluate on event insert, fire actions synchronously for simplicity (worker goroutine for retries).

**Why not pure event-sourced?** CRM queries dominate (admin UI lists sessions/tickets). Projections would add latency/complexity for take-home.

**Why not CRM-only?** Events needed for multi-source (webchat, WhatsApp webhook), idempotency, audit.

**Implementation**:
- Event POST persists, checks idempotency (processed_keys table).
- Rule evaluator: DB query enabled rules by event_type, JSONB conditions match (e.g. payload.sentiment in rule.conditions.sentiment).
- Actions: create_ticket impl, send_message stub (Resend Phase4).
- AI agent: Anthropic tools in /api/chat loop.

**Cuts for time budget**:
- Phase4 light (Resend email only, no WhatsApp sandbox).
- Keyword KB search (no vector/pgvector).
- Rules DB only (no YAML loader).
- No real auth (hardcoded).
- No full CRUD/admin UI beyond sessions list.
- No agent evals framework (manual curl tests).

**Faked**:
- Delivery retries (single attempt).
- Multi-turn chat poll (no WS).

**With another week**:
- pgvector + embeddings for KB.
- Worker queue (pgmq).
- Rule-builder UI.
- WhatsApp sandbox.
- Auth (NextAuth).
- Tests (95% coverage).
- Deploy Railway Postgres + full stack.

This slice demonstrates end-to-end: chat → agent (tools/KB) → session summary/sentiment → event → rule → ticket/email.

Tradeoffs explicit, production-ready where mattered (agent prompts, tools).

