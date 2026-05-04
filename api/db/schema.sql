-- Relay CRM Schema
-- Run with psql on Supabase

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. customers
CREATE TABLE IF NOT EXISTS customers (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  phone VARCHAR(20),
  email VARCHAR(255) UNIQUE,
  channel VARCHAR(20) DEFAULT 'web' CHECK (channel IN ('web', 'telegram', 'email')),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. sessions
CREATE TABLE IF NOT EXISTS sessions (
  id SERIAL PRIMARY KEY,
  customer_id INTEGER REFERENCES customers(id),
  status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'closed', 'escalated')),
  channel VARCHAR(20) DEFAULT 'web',
  sentiment VARCHAR(20) CHECK (sentiment IN ('Positif', 'Netral', 'Negatif', 'Sangat Negatif')),
  summary TEXT,
  csat_score INTEGER CHECK (csat_score BETWEEN 1 AND 5),
  csat_rationale TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  closed_at TIMESTAMP
);

-- 3. messages
CREATE TABLE IF NOT EXISTS messages (
  id SERIAL PRIMARY KEY,
  session_id INTEGER REFERENCES sessions(id),
  role VARCHAR(20) CHECK (role IN ('user', 'assistant', 'system')),
  content TEXT NOT NULL,
  kb_article_id INTEGER REFERENCES kb_articles(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. tickets
CREATE TABLE IF NOT EXISTS tickets (
  id SERIAL PRIMARY KEY,
  customer_id INTEGER REFERENCES customers(id),
  session_id INTEGER REFERENCES sessions(id),
  category VARCHAR(100) NOT NULL,
  priority VARCHAR(20) DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high', 'urgent')),
  status VARCHAR(20) DEFAULT 'open' CHECK (status IN ('open', 'in_progress', 'resolved')),
  summary TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  sla_due_at TIMESTAMP,
  resolved_at TIMESTAMP
);

-- 5. kb_articles
CREATE TABLE IF NOT EXISTS kb_articles (
  id SERIAL PRIMARY KEY,
  title VARCHAR(500) NOT NULL,
  content TEXT NOT NULL,
  tags TEXT[],
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 6. events
CREATE TABLE IF NOT EXISTS events (
  id SERIAL PRIMARY KEY,
  type VARCHAR(100) NOT NULL,
  source VARCHAR(100),
  payload JSONB,
  idempotency_key VARCHAR(255) UNIQUE,
  processed_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 7. rules
CREATE TABLE IF NOT EXISTS rules (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  event_type VARCHAR(100),
  match_conditions JSONB,
  action_type VARCHAR(100) CHECK (action_type IN ('create_ticket', 'send_message', 'update_customer', 'webhook')),
  action_payload JSONB,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 8. deliveries
CREATE TABLE IF NOT EXISTS deliveries (
  id SERIAL PRIMARY KEY,
  event_id INTEGER REFERENCES events(id),
  rule_id INTEGER REFERENCES rules(id),
  channel VARCHAR(20),
  recipient VARCHAR(255),
  payload JSONB,
  status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'delivered', 'failed')),
  attempt_count INTEGER DEFAULT 0,
  last_error TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  sent_at TIMESTAMP
);

-- 9. rule_audit_log
CREATE TABLE IF NOT EXISTS rule_audit_log (
  id SERIAL PRIMARY KEY,
  rule_id INTEGER REFERENCES rules(id),
  event_id INTEGER REFERENCES events(id),
  action_type VARCHAR(100),
  outcome VARCHAR(50),
  error_msg TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_sessions_customer ON sessions(customer_id);
CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions(status);
CREATE INDEX IF NOT EXISTS idx_messages_session ON messages(session_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_customer ON tickets(customer_id);
CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);
CREATE INDEX IF NOT EXISTS idx_events_idempotency ON events(idempotency_key);
CREATE INDEX IF NOT EXISTS idx_deliveries_status ON deliveries(status);
CREATE INDEX IF NOT EXISTS idx_rule_audit_rule ON rule_audit_log(rule_id);

