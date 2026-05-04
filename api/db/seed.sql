-- Relay Seed Data

-- Sample customers
INSERT INTO customers (name, phone, email, channel) VALUES
('Budi Santoso', '+6281234567890', 'budi@example.com', 'web'),
('Ani Wijaya', '+6282345678901', 'ani@example.com', 'telegram'),
('Dewi Lestari', '+6283456789012', 'dewi@example.com', 'email')
ON CONFLICT (email) DO NOTHING;

-- KB Articles - LubriMax lubricant product (Bahasa Indonesia)
INSERT INTO kb_articles (title, content, tags) VALUES
('Cara Mengganti Oli Mesin LubriMax', 'Langkah-langkah mengganti oli mesin dengan LubriMax:
1. Panaskan mesin 5 menit
2. Matikan mesin, buka baut pembuangan
3. Bersihkan filter oli lama
4. Pasang filter baru LubriMax
5. Tuang 4L LubriMax SAE 10W-40
6. Tutup baut, hidupkan mesin 1 menit
7. Cek kebocoran', ARRAY['oli', 'ganti oli', 'mesin']),
('Manfaat LubriMax SAE 10W-40', 'LubriMax SAE 10W-40:
- Sintetik full untuk mesin modern
- Perlindungan 10.000 km
- Mengurangi gesekan 30%
- Tahan suhu ekstrem -30°C to 150°C
- Cocok motor & mobil', ARRAY['manfaat', 'sintetik', 'perlindungan']),
('Frekuensi Ganti Oli LubriMax', 'Ganti oli LubriMax setiap:
- Motor: 2.000 - 3.000 km
- Mobil: 5.000 - 10.000 km
- Atau setiap 6 bulan
Gunakan saringan ori LubriMax untuk hasil optimal.', ARRAY['frekuensi', 'jadwal', 'maintenance']),
('Gejala Oli Perlu Diganti', 'Tanda oli sudah harus diganti:
- Mesin berisik/keras
- Asap knalpot biru
- Konsumsi BBM naik
- Oli hitam kental/bau gosong
- Lampu indikator oli menyala
Jangan tunggu, ganti segera dengan LubriMax!', ARRAY['gejala', 'tanda', 'warning']),
('Pilih Viskositas LubriMax', 'Pilih oli LubriMax berdasarkan:
- 10W-40: Cuaca tropis umum
- 20W-50: Mesin tua/panas berat
- 5W-30: Mesin modern EFI
- 15W-40: Diesel truk
Lihat manual kendaraan Anda.', ARRAY['viskositas', 'pilih oli', 'rekomendasi'])
ON CONFLICT DO NOTHING;

-- Rules
INSERT INTO rules (name, event_type, match_conditions, action_type, action_payload, is_active) VALUES
('auto-ticket-on-complaint', 'session.closed', '{"sentiment": ["Negatif", "Sangat Negatif"]}', 'create_ticket', '{"priority": "urgent", "category": "keluhan"}', true),
('auto-notify-on-overdue-ticket', 'ticket.overdue', '{}', 'send_message', '{"template": "Tiket Anda sudah lewat SLA. Tim kami sedang proses."}', true)
ON CONFLICT (name) DO NOTHING;

