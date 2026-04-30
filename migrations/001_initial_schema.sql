-- Relay Take-Home Schema
-- CRM-centered hybrid architecture

-- Customers table
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- KB Articles table
CREATE TABLE IF NOT EXISTS kb_articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    tags VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sessions table (chat conversations)
CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id),
    status VARCHAR(50) DEFAULT 'active', -- active, ended
    summary TEXT,
    sentiment VARCHAR(50), -- Positif, Netral, Negatif, Sangat Negatif
    predicted_csat INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    session_id INTEGER REFERENCES sessions(id),
    sender VARCHAR(50) NOT NULL, -- customer, agent, bot
    content TEXT NOT NULL,
    kb_article_id INTEGER REFERENCES kb_articles(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tickets table
CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id),
    session_id INTEGER REFERENCES sessions(id),
    category VARCHAR(100) NOT NULL,
    priority VARCHAR(50) DEFAULT 'medium', -- low, medium, high, critical
    summary TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'open', -- open, in_progress, resolved, closed
    sla_deadline TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Events table (for event-driven architecture)
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    type VARCHAR(100) NOT NULL,
    source VARCHAR(100),
    payload JSONB,
    idempotency_key VARCHAR(255) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Rules table (configurable rules)
CREATE TABLE IF NOT EXISTS rules (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    event_type VARCHAR(100),
    conditions JSONB, -- conditions to match
    action_type VARCHAR(100) NOT NULL, -- create_ticket, send_message, update_customer, webhook
    action_config JSONB,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Rule firings audit trail
CREATE TABLE IF NOT EXISTS rule_firings (
    id SERIAL PRIMARY KEY,
    rule_id INTEGER REFERENCES rules(id),
    rule_name VARCHAR(100),
    event_id INTEGER REFERENCES events(id),
    action VARCHAR(100),
    outcome VARCHAR(50), -- success, failed, skipped
    error_msg TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Deliveries table (message delivery tracking)
CREATE TABLE IF NOT EXISTS deliveries (
    id SERIAL PRIMARY KEY,
    event_id INTEGER REFERENCES events(id),
    rule_firing_id INTEGER REFERENCES rule_firings(id),
    channel VARCHAR(50), -- whatsapp, telegram, email, sms, web
    recipient VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pending', -- pending, sent, delivered, read, failed
    attempts INTEGER DEFAULT 0,
    last_error TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Processed idempotency keys for deduplication
CREATE TABLE IF NOT EXISTS processed_keys (
    key VARCHAR(255) PRIMARY KEY,
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_sessions_customer ON sessions(customer_id);
CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions(status);
CREATE INDEX IF NOT EXISTS idx_messages_session ON messages(session_id);
CREATE INDEX IF NOT EXISTS idx_tickets_customer ON tickets(customer_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);
CREATE INDEX IF NOT EXISTS idx_rule_firings_rule ON rule_firings(rule_id);
CREATE INDEX IF NOT EXISTS idx_rule_firings_event ON rule_firings(event_id);
CREATE INDEX IF NOT EXISTS idx_deliveries_status ON deliveries(status);

-- Seed KB articles with Bahasa Indonesia content
INSERT INTO kb_articles (title, content, tags) VALUES
('Cara Reset Password', 'Untuk mereset password Anda, silakan ikuti langkah berikut:\n1. Buka halaman login\n2. Klik "Lupa Password"\n3. Masukkan email Anda\n4. Cek inbox untuk link reset\n5. Ikuti instruksi untuk password baru', 'password, akun, reset'),
('Jam Operasional Layanan', 'Kami beroperasi 24 jam sehari, 7 hari seminggu.\nUntuk konsultasi produk, hubungi:\n- Telepon: 1500-XXX\n- WhatsApp: +62-XXX-XXXX-XXXX\n- Email: support@relay.id', 'operasional, jam, kontak'),
('Kebijakan Pengembalian', 'Produk dapat dikembalikan dalam waktu 14 hari sejak pembelian.\nSyarat:\n- Produk belum digunakan\n- Packaging masih lengkap\n- Siapkan bukti purchase', 'pengembalian, refund, garansi'),
('Cara Menghubungi Customer Service', 'Tim customer service kami siap membantu 24/7.\nMetode kontak:\n1. Live chat di website\n2. WhatsApp: +62-XXX-XXXX-XXXX\n3. Email: cs@relay.id\n4. Telepon: 1500-XXX', 'kontak, cs, customer service')
ON CONFLICT DO NOTHING;

-- Seed sample customer
INSERT INTO customers (name, email, phone) VALUES
('Budi Santoso', 'budi@example.com', '+6281234567890'),
('Ani Wijaya', 'ani@example.com', '+6282345678901'),
('Dewi Lestari', 'dewi@example.com', '+6283456789012')
ON CONFLICT DO NOTHING;

-- Seed sample rules
INSERT INTO rules (name, description, event_type, conditions, action_type, action_config, enabled) VALUES
('auto-ticket-on-complaint', 'Create ticket when session has negative sentiment', 'session.ended', '{"sentiment": ["Negatif", "Sangat Negatif"]}'::jsonb, 'create_ticket', '{"priority": "high", "category": "complaint"}'::jsonb, true),
('auto-notify-on-overdue-ticket', 'Notify customer when ticket is past SLA', 'ticket.overdue', '{}'::jsonb, 'send_message', '{"message": "Tiket Anda sudah melampaui batas waktu. Tim kami sedang memproses."}'::jsonb, true)
ON CONFLICT DO NOTHING;
