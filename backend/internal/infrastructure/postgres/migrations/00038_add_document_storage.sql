-- +goose Up
-- Evrak dosyaları MinIO/S3 nesne deposunda tutulur; bu kolonlar nesne anahtarını
-- ve bütünlük doğrulaması için sha256 sağlama toplamını saklar.
ALTER TABLE documents ADD COLUMN IF NOT EXISTS storage_bucket VARCHAR(255);
ALTER TABLE documents ADD COLUMN IF NOT EXISTS storage_key    TEXT;
ALTER TABLE documents ADD COLUMN IF NOT EXISTS checksum       VARCHAR(128);

CREATE INDEX IF NOT EXISTS idx_documents_storage_key ON documents(storage_key);

-- +goose Down
DROP INDEX IF EXISTS idx_documents_storage_key;
ALTER TABLE documents DROP COLUMN IF EXISTS checksum;
ALTER TABLE documents DROP COLUMN IF EXISTS storage_key;
ALTER TABLE documents DROP COLUMN IF EXISTS storage_bucket;
