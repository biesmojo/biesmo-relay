# Relay — Take-Home

> An AI-native CRM × event-driven workflow platform — think contact-center CRM + AI agents + n8n-style automations in one box. *Relay* events into actions, customer messages into AI replies, AI sessions into human handoffs. This take-home is a small vertical slice of it.

---

## 1. The role

This is an **AI Manager** role, not a software engineering role.

We do not value you typing code. We value you **steering AI to ship code**. The hire is for someone whose daily output is "the AI shipped this" — not "I shipped this with AI helping a bit". If you instinctively reach for the keyboard before reaching for the agent, this isn't the role.

What we are testing in this take-home:

1. **Your AI workflow.** How do you brief, plan, dispatch, review, and merge work that an AI does for you? Show us your slash commands, planning agents, subagent dispatches, eval loops, prompt iterations.
2. **Your taste in scope.** The brief is intentionally larger than the time budget. What you direct the AI to skip — and why — is graded harder than what you ship.
3. **Your judgment on AI products.** Tool definitions, agent prompts, fallback policies, hallucination bounds. Anyone can wrap a `chat()` call. We want someone who can ship an agent that is reliable in production.
4. **Your verification instinct.** When the AI says "done", how do you know? This is unhinted on purpose — you'll figure out the answer in the work.

What we are explicitly **not** testing:

- Whether you can hand-write Go or React. Use the AI for that.
- Whether you remember syntax, library APIs, or framework conventions. The AI does.
- Whether you write "clean code" in your own hand. We grade the result, not the keystrokes.

A rough sniff test: at the end of this take-home, if your `AI_LOG.md` could plausibly read "I directed Claude Code / Codex / Cursor through the entire build via planning agents and reviewed the output", that's the right shape. If it reads "I wrote most of it and used Copilot for autocomplete", you're applying for the wrong role.

---

## 2. The product: Relay

Relay is the working name for an internal product that fuses three things into one platform:

- A **CRM** — customer records, sessions, tickets, knowledge base. Inspired by [Lovina](https://www.ailovina.com), Metatech's AI-native contact-center product.
- **AI agents** — tool-using LLM agents that handle inbound conversations, retrieve knowledge, create tickets, and hand off to humans when stuck.
- A **workflow engine** — events come in, rules match, actions fire across channels (WhatsApp / push / email / SMS / webhook). Think n8n / Zapier / Make.com, but native to the CRM and the agents. Inspired by **SIGAP**, an event-driven communication platform commissioned by Badan Gizi Nasional (Indonesia's national nutrition agency).

Most products in this space pick one center of gravity. **A CRM with bolted-on AI** is what most contact-center incumbents are doing. **An automation tool with bolted-on chat** is what most workflow vendors are doing. Relay tries to put all three in the middle, and how you architect that is the open design question.

You don't need to build all of any one of these. You need to build a **vertical slice** that exercises every layer: event in → rule fires → AI agent acts → message goes out → ticket gets created → operator sees it on a screen.

---

## 3. Ground rules

- **Target effort:** 5–8 hours of focused work, spread over up to 1 week. We don't enforce a timer.
- **AI tools — required, not optional.** Claude Code / Codex / Cursor / Aider — pick your daily driver. Manual coding is *discouraged*. If you find yourself typing implementation by hand, stop and ask why.
- **Stack — pinned.** **Backend in Go**, **frontend in React** (Next.js, Remix, Vite + React, TanStack Start — your call within the React family). This matches our production stack; the AI you use needs to be productive in it. Don't bikeshed the stack choice — pick a React framework in 5 minutes and move on.
- **LLM:** Any provider. We use Claude (Sonnet 4.6 / 4.7) in production, slight preference for that. OpenAI, Gemini, etc. all fine.
- **Libraries / SaaS:** Use whatever you want. Supabase, Neon, Upstash, Vercel, Fly, Railway, Resend, Twilio, WhatsApp Cloud API sandbox, Telegram bot — all fair game. Free tiers are fine.
- **Language:** Customer-facing strings (notifications, AI replies, error messages a user sees) must be in **Bahasa Indonesia**. Everything else (code, comments, design doc, AI log, READMEs) is in **English**. This matches the production context.
- **Authentication:** A single hardcoded admin account is fine. Don't burn time on real auth.
- **Deployment:** Either a live deployed URL OR a 5-minute Loom of you walking through the happy paths. We need to *see* it work.

---

## 4. Greenfield — you direct the AI to bootstrap

There is no starter scaffold. The repo starts empty. Direct the AI to set up the Go backend, the React frontend, the schema, local dev, deploy.

This is intentional. The brief in Section 5 is **larger than the time budget**. We want to see how you delegate the bootstrap to an AI agent without rabbit-holing.

A few hints so you don't waste time:

- **Postgres** is the safe DB choice — Supabase / Neon / Render / local Docker, all work. `sqlc` or `gorm` for the Go side, both fine.
- **Go BE, React FE — typically split** into two folders or two services. A monorepo with `/api` (Go) and `/web` (React) is the obvious shape. Don't reinvent.
- **Don't have the AI hand-roll auth, queues, or ORMs.** Direct it to pull the boring library.
- **Don't pre-build for 100k events/sec.** Postgres + a worker goroutine + `pgmq` / `River` / cron is more than enough.

We are NOT grading the cleanliness of your initial commit. We are grading the slice you ship by the end, and the quality of the AI direction that produced it. Get past `hello world` fast.

---

## 5. What to build — 4 phases

Each phase is a **vertical slice**. Ship depth, not breadth. A working Phase 1+2+3 beats a shallow attempt at all four.

### Phase 1 — CRM core

The Lovina-flavored half. Enough CRM primitives that an event has somewhere to land.

- **Customer**, **Session**, **Ticket**, **Message** models — direct the AI to design the schema. Sensible columns for a CRM, no over-engineering.
- One inbound channel that creates Sessions. **Web chat is fine** (a `/chat` page where a user types, messages persist, agent replies — see Phase 3).
- **Auto-ticket from session** rule: when a session ends with sentiment Negative or contains escalation keywords, a Ticket is created with `customer_id`, `category`, `priority`, `summary` filled in. Dedup: same customer + same category in 24h appends to existing ticket instead of creating a new one.
- **KB retrieval into agent context:** when the agent answers, retrieved KB article content is in the prompt and the article ID is logged on the message.
- A bare admin UI: list of sessions, click into a session to see the transcript and any linked ticket.

### Phase 2 — Event + rule layer (the SIGAP marriage)

The SIGAP-flavored half. Make the rule engine real.

- `POST /events` accepts a JSON event with `type`, `source`, `payload`, idempotency key. Persist to `events` table.
- A **rule evaluator**: rules live in the database (table) or in YAML / JSON files (your call). Each rule has a `match` (event type + payload conditions) and an `action` (one of: `create_ticket`, `send_message`, `update_customer`, `webhook`). Bonus if you build a tiny rule-builder UI; not required if rules-as-code is cleaner for you.
- Ship at least **2 working rules**:
  1. `auto-ticket-on-complaint`: incoming web-chat session with sentiment ≤ Negative → create a ticket.
  2. `auto-notify-on-overdue-ticket`: ticket past SLA → send a message to the customer (in Bahasa Indonesia) on the channel they originally used.
- **Idempotency:** sending the same event twice (same idempotency key) MUST NOT trigger the action twice.
- **Retry:** action failures retry with exponential backoff up to 3 times. Persist delivery state (`pending` → `sent` → `failed`) on the `deliveries` row so failures are visible.
- **Audit trail:** every rule firing writes a row that names the rule, the event id, the action, and the outcome.

### Phase 3 — AI agent (the Lovina-flavored core)

The "AI product" muscle. **This phase is the centerpiece of the rubric.**

Build a tool-using agent that handles inbound web-chat conversations. The agent:

- Has **at least these tools**, defined as proper LLM tool schemas:
  - `lookup_customer({phone | email})` — returns customer record from `customers`
  - `search_kb({query})` — returns top-k matching KB article excerpts (vector or keyword, your call)
  - `create_ticket({customer_id, category, priority, summary})` — creates a ticket; returns id
  - `escalate_to_human({reason, urgency})` — flags the session for human pickup; closes the agent loop
- Drives a real multi-turn loop: it should call tools when it needs data, not hallucinate customer fields.
- Produces, at session end, four artifacts attached to the session row:
  1. **Transcript** — full message log (already in DB)
  2. **Summary** — 2–4 sentences in Bahasa Indonesia
  3. **Sentiment** — one of `Positif | Netral | Negatif | Sangat Negatif`
  4. **Predicted CSAT** — integer 1–5 with a one-sentence rationale
- Has a **fallback policy** when intent is unclear (asks 1 clarifying question, then escalates) and a **handoff path** when the user explicitly asks for a human (calls `escalate_to_human` and closes).
- **Agent evaluations:** how do you know the agent works? Direct the AI to set up whatever you need to be able to answer that question with confidence — the form is your call. We will look at how you verify your agent's behavior across the fixtures you choose.

### Phase 4 — Multi-channel delivery + observability

Make at least one channel besides web chat work for real-ish.

- Pick one: **WhatsApp Cloud API sandbox**, **Telegram bot**, or **email via Resend free tier**. Real network calls. (Twilio SMS works too if you have credits.)
- Outbound delivery flows through the same `deliveries` pipeline as in Phase 2: the rule engine fires, writes a delivery row, the channel adapter picks it up, sends, updates status.
- **Delivery status webhook (or polling)** updates the row to `delivered` / `read` / `failed`.
- A tiny **event explorer** page (`/admin/events`) showing recent events, which rules matched, which actions fired, status. This is the SIGAP-D.1 outcome — let us debug your system.

---

## 6. The marriage decision

Lovina centers the CRM. SIGAP centers the event bus.

You pick the architecture. There are reasonable choices on both sides:

- **CRM-centered:** Sessions and tickets are the heart; events are how Sessions get rules applied to them. Rule engine sits *next to* the CRM, not above.
- **Event-centered:** Everything (a chat message, a ticket update, a CSAT response) is an event. CRM tables are downstream projections. Rules subscribe to event streams.
- **Hybrid:** CRM is the system of record; an internal event bus connects modules so rules can react. (This is roughly what production Lovina is moving toward.)

There is no wrong answer. **In `DESIGN.md`, write 200–400 words on which you picked, what you cut to make the time budget, and what you'd do differently with another week.** We grade your reasoning, not your choice.

---

## 7. AI-tool expectations

We expect you to use AI tooling daily. We're not testing whether you can avoid AI; we're testing how *good* you are with it.

You will keep an `AI_LOG.md` that includes:

- **Which tools** you used (Claude Code, Codex, Cursor, Aider, etc.) and roughly when in the project.
- **Your slash commands / subagents / skills** — paste them in. If you wrote a custom skill or `CLAUDE.md`, include it.
- **3 specific moments** with timestamps and what happened:
  1. **A win** — where AI saved you serious time. Show the prompt and the result.
  2. **A miss** — where AI misled you, hallucinated, or churned. What did you do?
  3. **A move** — where you outgrew the default workflow (e.g. you stopped accepting one-shot completions and started running a planning agent first; you started using subagents for research while you wrote tests).
- **Final agent prompts** (the system prompts you ended up shipping for Phase 3). These are the most important artifact in this file. We will read them carefully.

`AI_LOG.md` is graded. A thin or vague log is a worse signal than no log at all — it suggests you don't reflect on your tools.

We are also looking at your **commit cadence and messages**. AI-assisted work tends to leave a recognizable shape: many small commits, sometimes co-authored. That's fine. We're not penalizing that — we're confirming it.

---

## 8. What you submit

Four things:

1. **GitHub repo** (private) shared with `@listiarso` (or the GitHub handle in the email reply, if I asked you to add a different one). Main branch is what I'll grade.
2. **Live demo** — deployed URL (Vercel/Fly/Render/Railway — your call) **OR** a 5-minute Loom walking through: web chat → AI agent reply → KB cited → ticket created → outbound delivery via the non-web channel → event explorer view. The Loom must be ≤7 minutes.
3. **`DESIGN.md`** in the repo root — the marriage decision (Section 6), what you cut, what's faked vs real, what you'd do next.
4. **`AI_LOG.md`** in the repo root — Section 7 contents.

`README.md` should be the entrypoint: how to run it locally in <5 minutes, link to the demo, link to `DESIGN.md` and `AI_LOG.md`.

---

## 9. Evaluation rubric

| Weight | Criterion | What we look for |
|---|---|---|
| 30% | **AI direction quality** | `AI_LOG.md` shows planning, dispatch, review. Prompts are sharp. You used the right tool for the right step. Commit cadence and AI-co-author trail align with the workflow you describe. |
| 25% | **AI-product engineering** | Tool definitions are tight, agent prompts are well-constructed, fallback/handoff logic is sound, hallucinations are bounded. The Phase 3 agent feels like a product, not a demo. |
| 25% | **Does it work end-to-end** | The slice runs. Web chat → agent → KB → ticket → outbound channel → event explorer. Verified live in your demo. |
| 15% | **Scope & judgment** | The marriage decision (Section 6) is reasoned. What you cut is justified. What you shipped reflects taste, not panic. |
| 5% | **Communication** | `DESIGN.md`, `AI_LOG.md`, demo, README — clear, brief, no fluff. |

There are also a small number of **unannounced sub-criteria** we look for. They're not hidden in a tricky way — they're things a senior AI Manager would naturally produce as part of the workflow, that a junior wouldn't think to ask the AI for. We trust you to figure out what those are.

---

## 10. Out of scope (please don't build these)

These are explicit cuts. They're real Lovina/SIGAP features but they would blow your time budget:

- **Telephony / SIP / VoIP** — no real voice calls. Web chat is the only channel that needs an interactive surface.
- **Real CSAT survey send** — the agent's *prediction* is in scope (Phase 3). Sending a real survey link to a customer is not.
- **Billing / RAB / payment integration** — no SAKTI, no Stripe, no invoicing.
- **RBAC beyond user/admin** — one admin login is enough. Don't build 30 resource × 4 permission matrices.
- **Multi-tenant isolation** — single-tenant is fine. Don't build company-scoping.
- **PDP-compliant data residency** — don't worry about Indonesia data sovereignty rules. Use whatever cloud is fastest.
- **SAML/OAuth/SSO** — out.
- **Outbound campaign predictive dialing, A/B test framework, agent QA scorecards, bulk audio evaluation** — out.
- **WhatsApp Business API production approval** — sandbox only.
- **Production-grade observability** — a Postgres `events` table with a tiny UI is the bar.
- **A mobile app, a native app, an Electron app** — web only.

If you're unsure whether something is in scope, default to "no". Ship the slice.

---

## 11. FAQ

**Q: The stack — Go and React. Any sub-pinning?**
A: No. Pick the Go web framework you'd direct an AI agent through fastest (`net/http` + `chi`, Gin, Echo, Fiber — all fine). Pick the React framework you trust an AI to scaffold cleanly (Next.js is the safe default; Vite + React is fine; Remix / TanStack Start work too). Don't use Svelte / Vue / Angular — wrong stack.

**Q: Can I use Supabase / Neon / Postgres on Render?**
A: Yes. Use whatever Postgres is fastest for you.

**Q: Can I skip auth entirely?**
A: One hardcoded admin (env var) is fine. We won't test auth.

**Q: Which LLM should I use?**
A: Any. Claude Sonnet 4.6 or 4.7 is what we use. Anthropic, OpenAI, Gemini all fine.

**Q: Can I vibe-code the whole thing with Claude Code / Cursor / Codex?**
A: That is literally what we're testing. Show your direction in `AI_LOG.md`. We'll be able to tell whether you understood what shipped or just accepted whatever the agent produced.

**Q: Is the rule-builder UI required?**
A: No. Rules-as-code (YAML / JSON / Go files) is acceptable. A small UI is a nice-to-have if you have time after Phase 4.

**Q: How real does WhatsApp delivery have to be?**
A: WhatsApp Cloud API sandbox sending to your own test number counts. Telegram bot to a test channel counts. Resend email to your own inbox counts. SMS via Twilio trial counts. The point is a real network call leaves your code.

**Q: What if I can't finish all 4 phases?**
A: Common. Ship Phases 1+2+3 deeply and skip Phase 4 if you must. Note what you cut in `DESIGN.md`.

**Q: What if the spec is ambiguous?**
A: Make the call, direct the AI accordingly, document it in `DESIGN.md`. Over-clarifying is a soft signal that you can't make calls.

**Q: Can I use a UI library / component kit?**
A: Yes — shadcn/ui, Mantine, MUI, whatever. Don't burn time on visual polish; we don't grade pixels.

**Q: I write a lot of code by hand because I'm faster than the AI for [X]. Is that OK?**
A: Treat that as a yellow flag for whether this role is right for you. The hire is for someone whose AI workflow is faster than their hand-coding. If that's not yet you, this isn't the right time to apply.

---

## 12. Submission

This take-home is being announced openly in my mentee group. There is no recruiter funnel, no formal pipeline, no take-home invite email.

If you want to apply:

1. Build the slice. Make your repo private.
2. Email **gogo@metatech.id** with: a one-paragraph intro about you, the repo URL, the demo URL or Loom link, and the GitHub handle to add as a collaborator.
3. **Do not expect a reply.** I might respond in a day, in a month, or never. I read everything but I don't promise feedback. If you need a fast yes/no for life planning, this isn't that.

If the spec is contradictory in a way you can't paper over, or an external API is down, email me. Asking for *direction* on judgment calls is a soft negative; asking for *help* on real blockers is fine.

Have fun. The brief is real — we're actually building this.
