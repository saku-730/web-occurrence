-- +goose Up
-- make_specimen テーブルから date カラムを削除するのだ。
ALTER TABLE public.make_specimen DROP COLUMN date;

-- +goose Down
-- もしロールバック（down）したら、date カラムを元の型 (date) で復元するのだ。
ALTER TABLE public.make_specimen ADD COLUMN date date;
