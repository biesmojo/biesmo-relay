# Relay Architecture Decision (298 words)

**Architecture: Hybrid CRM + Event Bus**

CRM tables (customers, sessions, tickets, messages) as system of record — they're queried constantly by admin UI and agent. Events provide loose coupling between modules: chat → session event → rules → actions → Telegram delivery. Events enable idempotency, audit, retry — no hand-rolled queues.

**Reasoning:** Pure event-sourced too heavy for take-home (projections, replay). CRM-only can't handle multi-source (webchat webhook), retries, audit. Hybrid gets both: fast reads from CRM, reliable delivery from events.

**Cuts for scope:**
- Multi-channel: Only Telegram (real API), no WhatsApp/Email sandbox setup
- Rule-builder UI → rules as JSON in DB (add UI next week)
- KB search: ILIKE keyword (no pgvector embeddings — real search needs setup)
- Auth: Mock Bearer token (no JWT/Supabase Auth)
- CSAT: Predict only (real survey send would need Resend + templates)

**Real vs Fake:**
- **Real:** Telegram delivery (Bot API HTTP), AI agent (Claude tools), Postgres events/rules (idempotent)
- **Fake:** create_ticket check (stubbed 24h logic), no real webhook channel, session summary generation (Claude prompt)

**One more week:**
- pgvector + pgai embeddings for semantic KB search
- Rule-builder UI with JSON editor/preview
- WhatsApp Business sandbox + channel adapter
- Real CSAT survey via Resend email templates + link tracking
- pgmq for async rule processing (scale)
- Supabase Edge Functions for serverless deploy

Tradeoffs explicit, production-ready core (agent loop, tool calling, retry logic). Slice demonstrates full flow: chat → tools → session.close → event → rule → Telegram delivery → admin dashboard.

