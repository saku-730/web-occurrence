-- +goose Up
-- ワークステーションとユーザー紐付けテーブルの全データを削除し、
-- workstationテーブルのシーケンス（workstation_id）をリセットするのだ
-- CASCADEオプションで、workstation_userテーブルも同時にクリアされるのだ
TRUNCATE TABLE workstation RESTART IDENTITY CASCADE;

-- +goose Down
-- データ削除のDownマイグレーションは、通常は行わないのだ。
-- データを復元するための適切なバックアップ戦略を別に立ててほしいのだ。
SELECT 1;
