-- +goose Up
-- users テーブルから workstation_id カラムを削除するのだ。
-- これで、関連する外部キー制約も自動的に削除されるのだ。
ALTER TABLE public.users DROP COLUMN workstation_id;

-- +goose Down
-- もしロールバック（down）したら、カラムと外部キー制約を元通り復元するのだ。
ALTER TABLE public.users ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
