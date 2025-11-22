-- +goose Up
-- CASCADE: 紐付いている他のテーブルのデータも一緒に消す
-- RESTART IDENTITY: シーケンス番号を初期化する
TRUNCATE TABLE public.workstation RESTART IDENTITY CASCADE;
-- +goose Down
