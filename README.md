# Relay — AI-Native CRM (Vertical Slice)

A minimal vertical slice of an AI-native CRM × event-driven workflow platform.

This project demonstrates:

* Web chat → AI agent reply (with tools)
* KB retrieval into responses
* Auto ticket creation via rules
* Event-driven rule engine with idempotency + retries
* Admin visibility (sessions, tickets, events)

---

## 🚀 Demo

* Live URL: <YOUR_URL>
* Loom walkthrough (≤7 min): <YOUR_LOOM_LINK>

---

## 🧱 Tech Stack

* Backend: Go (chi + PostgreSQL)
* Frontend: React (Next.js)
* DB: Postgres (Supabase)
* LLM: Claude Sonnet (tool-using agent)

---

## ⚡ Quick Start (local, <5 min)

### 1. Clone repo

```
git clone <repo-url>
cd relay
```

### 2. Setup env

```
cp .env.example .env
```

Fill:

* DATABASE_URL
* LLM_API_KEY

---

### 3. Run backend

```
cd api
go run main.go
```

---

### 4. Run frontend

```
cd web
npm install
npm run dev
```

---

Open:

* Chat: http://localhost:3000/chat
* Admin: http://localhost:3000/admin

---

## 🔁 Happy Path (what to test)

1. User sends message in `/chat`
2. Agent:

   * retrieves KB
   * responds in Bahasa Indonesia
3. Negative sentiment → ticket created
4. Rule triggers outbound message
5. Check `/admin/events` for rule execution

---

## 📂 Docs

* Architecture & decisions → `DESIGN.md`
* AI workflow & prompts → `AI_LOG.md`

---

## ⚠️ Notes

* Auth is mocked (single admin)
* Some channels simulated except one real integration (email/Telegram)
* Focus is on agent + rule engine reliability, not UI polish
