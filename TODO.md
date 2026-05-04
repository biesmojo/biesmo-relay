# Relay Development TODO

## Current Progress
- [x] Landing page implementation complete.

## Approved Plan Steps (Phases)

### Prep (0)
- [x] Step 0.1: Create .env.example with Supabase/Claude/admin vars
- [x] Step 0.2: Update api/ deps to chi/pgx + Anthropic SDK (run cd api && go mod tidy manually)
- [x] Step 0.3: Update web/ deps to add @anthropic-ai/sdk (npm i)
- [x] Step 0.4: Schema already complete in 001_initial_schema.sql (all tables/seeds/indexes/Bahasa KB done!)

### Phase 1: CRM Core
- [ ] Step 1.1: Add models: session.go, ticket.go, message.go, kb_article.go
- [ ] Step 1.2: Add handlers: sessions.go, tickets.go, messages.go with CRUD
- [ ] Step 1.3: Update main.go routes for new endpoints
- [ ] Step 1.4: Frontend admin UI: /admin/customers, /admin/sessions, /admin/tickets (list views)

### Phase 2: Event + Rule Layer
- [ ] Step 2.1: Add models: event.go, rule.go
- [ ] Step 2.2: handlers/events.go: POST /events with idempotency, rule evaluator
- [ ] Step 2.3: handlers/rules.go: CRUD rules (JSONB conditions)
- [ ] Step 2.4: Admin UI event/rule explorer

### Phase 3: AI Agent
- [ ] Step 3.1: handlers/chat.go: POST /chat with Claude API (tools: lookup_customer, search_kb, create_ticket)
- [ ] Step 3.2: Update web/src/app/chat/page.tsx: integrate /api/chat, multi-turn, Bahasa UI
- [ ] Step 3.3: Session artifacts (summary/sentiment)

### Phase 4: Multi-channel
- [ ] Step 4.1: handlers/delivery.go: Resend email stub→real
- [ ] Step 4.2: Rule actions include send_email

### Final
- [ ] Update README.md/DESIGN.md with details
- [ ] Test E2E: npm run dev & go run, curl chat→event→rule
- [ ] attempt_completion

**Next: Prep steps. Run after each: api/ go run main.go test endpoints; web/ npm run dev check UI.**
