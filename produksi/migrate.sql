-- ============================================================
-- ChatGo v2 SaaS — Migration Script
-- Jalankan SEKALI saat upgrade ke versi SaaS (multi-tenant)
-- Semua ALTER TABLE pakai IF NOT EXISTS pattern manual
-- ============================================================

-- 1. Tambah user_id ke semua tabel utama (isolasi multi-user)
ALTER TABLE autoreplies ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE campaigns ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE scheduled ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE contacts ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE contact_groups ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE templates ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE tags ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE canned_responses ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE blacklist ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE ai_keys ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE devices ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE ai_trainings ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE drips ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE drip_steps ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL DEFAULT 0;

-- 2. Tabel session owner — restore userID setelah server restart
CREATE TABLE IF NOT EXISTS wa_session_owners (
    phone VARCHAR(64) PRIMARY KEY,
    user_id BIGINT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. Perbarui struktur autoreplies (jika belum ada)
ALTER TABLE autoreplies ADD COLUMN IF NOT EXISTS account_id VARCHAR(512) NOT NULL DEFAULT '';
ALTER TABLE autoreplies ADD COLUMN IF NOT EXISTS training_id BIGINT NOT NULL DEFAULT 0;

-- ============================================================
-- Verifikasi
SELECT 'migrate-v2-saas completed' AS status;
-- ============================================================
