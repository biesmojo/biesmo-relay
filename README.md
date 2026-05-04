# Relay

AI-native CRM and event-driven workflow platform. Agents handle conversations with tools (KB lookup, ticket creation), rules fire actions (Telegram notifications), full admin dashboard for sessions/tickets/events.

**Demo:** <YOUR_URL>
**Loom:** <YOUR_LOOM_LINK>

## Tech Stack

| Component | Technology |
| --- | --- |
| Backend | Go + Chi + pgx + Supabase Postgres |
| Frontend | Next.js 14 App Router + Tailwind + TypeScript |
| AI | Anthropic Claude 3.5 Sonnet (tools) |
| Events | Postgres events + rules + deliveries |
| Channel | Telegram Bot API (outbound) |

## Quick Start (<5 min local)

```bash
git clone <repo>
cp .env.example .env # fill Supabase DATABASE_URL, Anthropic API key, Telegram BOT_TOKEN/CHAT_ID, ADMIN_TOKEN
cd api && go mod tidy && go run main.go # :8080
# (optional: psql DATABASE_URL -f api/db/schema.sql -f api/db/seed.sql)
cd ../web && npm i && npm run dev # :3000
```

## Happy Path Test

1. Backend running
2. Frontend /admin/events
3. POST /api/events (curl or Postman):
```json
{
  "type": "ticket.overdue",
  "source": "system",
  "idempotency_key": "test1",
  "payload": {"customer_id": 1, "channel": "telegram"}
}
```
4. Refresh /admin/events → see row status=sent/green
5. Telegram test chat receives notification
6. /chat → AI agent replies Bahasa, tools work

See DESIGN.md for architecture, AI_LOG.md for agent prompt/debug.

