# Relay

Relay is a small vertical slice of an AI-native CRM, event-driven workflow engine, and AI agent system.

It demonstrates:

- CRM primitives: customers, sessions, tickets, and messages
- An event ingestion and rule engine
- A tool-using AI agent for inbound web chat
- Outbound delivery through a non-web channel
- An admin UI for reviewing sessions, tickets, and events

## Demo

- Live demo: <YOUR_DEPLOYED_URL>
- Loom walkthrough: <YOUR_LOOM_URL>

## Architecture

- Backend: Go
- Frontend: React
- Database: PostgreSQL
- AI provider: <YOUR_LLM_PROVIDER>
- Non-web channel: <WhatsApp / Telegram / Email>

See [DESIGN.md](./DESIGN.md) for the architecture decision, tradeoffs, and scope cuts.

## Features

### CRM Core
- Customer records
- Chat sessions
- Tickets
- Message history
- Session transcript view in admin UI

### Event + Rule Engine
- `POST /events` endpoint
- Persistent event storage
- Rule matching and action execution
- Idempotency protection
- Retry with exponential backoff
- Delivery state tracking
- Audit trail for rule firing

### AI Agent
- Inbound web chat handling
- Tool calling for:
  - customer lookup
  - KB search
  - ticket creation
  - human escalation
- Session summary in Bahasa Indonesia
- Sentiment classification
- Predicted CSAT scoring

### Admin UI
- Session list
- Session detail page
- Linked ticket visibility
- Event explorer

## Project Structure

```txt
/api    # Go backend
/web    # React frontend