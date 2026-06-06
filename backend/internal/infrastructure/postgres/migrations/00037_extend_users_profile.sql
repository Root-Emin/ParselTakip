-- +goose Up
-- Kullanıcı tablosunu detaylı kimlik/iletişim/firma ve KVKK alanları ile genişlet.
-- national_id (TCKN) ve hassas alanlar uygulama katmanında AES-256-GCM ile şifrelenir;
-- national_id_hash, şifreli kolonda eşitlik araması için HMAC kör-indeksidir (KVKK).
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone                VARCHAR(30);
ALTER TABLE users ADD COLUMN IF NOT EXISTS national_id          TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS national_id_hash     VARCHAR(64);
ALTER TABLE users ADD COLUMN IF NOT EXISTS user_type            VARCHAR(30) NOT NULL DEFAULT 'citizen';
ALTER TABLE users ADD COLUMN IF NOT EXISTS address              TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS city                 VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS district             VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS company_id           UUID REFERENCES contractor_companies(id) ON DELETE SET NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS kvkk_consent_at      TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS kvkk_consent_version VARCHAR(20);
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_at        TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS failed_login_count   INT NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS locked_until         TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_users_national_id_hash ON users(national_id_hash);
CREATE INDEX IF NOT EXISTS idx_users_user_type ON users(user_type);
CREATE INDEX IF NOT EXISTS idx_users_company ON users(company_id);

-- +goose Down
DROP INDEX IF EXISTS idx_users_company;
DROP INDEX IF EXISTS idx_users_user_type;
DROP INDEX IF EXISTS idx_users_national_id_hash;
ALTER TABLE users DROP COLUMN IF EXISTS locked_until;
ALTER TABLE users DROP COLUMN IF EXISTS failed_login_count;
ALTER TABLE users DROP COLUMN IF EXISTS last_login_at;
ALTER TABLE users DROP COLUMN IF EXISTS kvkk_consent_version;
ALTER TABLE users DROP COLUMN IF EXISTS kvkk_consent_at;
ALTER TABLE users DROP COLUMN IF EXISTS company_id;
ALTER TABLE users DROP COLUMN IF EXISTS district;
ALTER TABLE users DROP COLUMN IF EXISTS city;
ALTER TABLE users DROP COLUMN IF EXISTS address;
ALTER TABLE users DROP COLUMN IF EXISTS user_type;
ALTER TABLE users DROP COLUMN IF EXISTS national_id_hash;
ALTER TABLE users DROP COLUMN IF EXISTS national_id;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
