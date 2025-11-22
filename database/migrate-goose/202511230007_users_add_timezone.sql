-- +goose Up
-- usersテーブルに timezone カラムを追加するのだ
-- デフォルト値は 'Asia/Tokyo' (またはUTCなど) にしておくと安全なのだ
ALTER TABLE users ADD COLUMN timezone text DEFAULT 'Asia/Tokyo';

-- +goose Down
ALTER TABLE users DROP COLUMN timezone;
