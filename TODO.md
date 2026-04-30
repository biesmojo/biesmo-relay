# Relay Take-Home TODO

## Phase 0: Setup & Fixes
- [ ] Fix migrations/001_initial_schema.sql (JSONB seeds)
- [ ] docker-compose.yml + .env.example (Postgres)
- [ ] execute docker-compose up (verify DB/seeds)
- [x] Add Gorilla Mux + basic customers handler / main mux update ✅
- [ ] Implement Phase 1 handlers

## Phase 1: CRM Core
- [ ] CRUD handlers: customers, sessions, messages, tickets, KB
- [ ] Web chat UI (/chat)
- [ ] Admin UI (/admin/sessions)
- [ ] Auto-ticket stub

## Phase 2: Event + Rule Layer
- [ ] /api/events + evaluator
- [ ] 2 rules working
- [ ] Idempotency/retry/audit
- [ ] Admin events explorer

## Phase 3: AI Agent
- [ ] Tools + multi-turn /api/chat
- [ ] Session artifacts
- [ ] Fallback/handoff
- [ ] Agent evals/tests

## Phase 4: Multi-channel
- [ ] Resend email delivery
- [ ] Event explorer UI

## Docs & Ship
- [x] DESIGN.md ✅
- [x] AI_LOG.md ✅
- [x] README.md ✅
- [ ] Deploy/Loom

