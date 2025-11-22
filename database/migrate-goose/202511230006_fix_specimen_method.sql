-- +goose Up
-- ============================================================
-- Specimen Methods: ID (Integer/Seq) -> UUID (Text)
-- ============================================================

-- PK削除
ALTER TABLE specimen_methods DROP CONSTRAINT IF EXISTS specimen_methods_pkey CASCADE;

-- IDカラムの型変更とデフォルト値設定
ALTER TABLE specimen_methods ALTER COLUMN specimen_methods_id DROP DEFAULT;
ALTER TABLE specimen_methods ALTER COLUMN specimen_methods_id TYPE text USING gen_random_uuid()::text;
ALTER TABLE specimen_methods ALTER COLUMN specimen_methods_id SET DEFAULT gen_random_uuid()::text;

-- PK再設定
ALTER TABLE specimen_methods ADD CONSTRAINT specimen_methods_pkey PRIMARY KEY (specimen_methods_id);

-- +goose Down
-- This migration is destructive.
