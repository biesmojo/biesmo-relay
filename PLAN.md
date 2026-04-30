# Relay Take-Home Project Plan

## Architecture Decision: CRM-Centered Hybrid

Chosen approach: **CRM-centered with internal event bus**
- CRM tables (customers, sessions, tickets, messages) are the system of record
- Events are how state changes propagate and trigger rules
- Rules sit next to the CRM, not above it
- This balances Lovina's CRM strength with SIGAP's event-driven patterns

## Scope Cuts (to fit time budget)

### Explicitly Cutting:
- Rule-builder UI → rules-as-code (YAML/JSON) instead
- Real auth → hardcoded admin env var
- Multi-tenant → single-tenant only
- Production-grade vector search → keyword search for KB
- Mobile/native app → web only
- Real telephony → web chat only

### Prioritization:
**Must have**: Phase 1 (CRM core) + Phase 2 (rule engine) + Phase 3 (AI agent)
**Nice to have**: Phase 4 (multi-channel delivery)

## Phase Breakdown

### Phase 1: CRM Core
- [ ] Postgres schema: customers, sessions, tickets, messages, kb_articles
- [ ] Go API: CRUD endpoints
- [ ] Web chat UI: message input, transcript display
- [ ] Auto-ticket rule: sentiment Negative → create ticket
- [ ] Admin UI: session list, click into transcript

### Phase 2: Event + Rule Layer
- [ ] POST /events endpoint
- [ ] Rules table + evaluator
- [ ] 2 working rules (auto-ticket-on-complaint, auto-notify-on-overdue)
- [ ] Idempotency + retry logic
- [ ] Audit trail (rule_firings table)

### Phase 3: AI Agent (CENTERPIECE)
- [ ] Tool schemas: lookup_customer, search_kb, create_ticket, escalate_to_human
- [ ] Multi-turn conversation loop
- [ ] Session artifacts: summary, sentiment, predicted CSAT
- [ ] Fallback + handoff policies
- [ ] Agent evaluation setup

### Phase 4: Multi-channel Delivery
- [ ] Telegram bot OR Resend email (real network call)
- [ ] Delivery pipeline integration
- [ ] Event explorer UI

## AI Workflow

Primary tool: **Claude Code** (or CLI mode)
Strategy:
1. Bootstrap with planning agent → get structure fast
2. Implement phase-by-phase with focused prompts
3. Review AI output before accepting
4. Run verification tests

## Verification Strategy

- Backend: Go tests + curl manual verification
- Frontend: browser manual walkthrough
- AI Agent: scripted conversation tests with expected outcomes
- End-to-end: demo walkthrough video
