-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- 1. シーケンスを作成するのだ
CREATE SEQUENCE public.workstation_workstation_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- 2. シーケンスを workstation_id カラムに紐付けるのだ
-- (テーブルが消えたらシーケンスも消えるようにする設定)
ALTER SEQUENCE public.workstation_workstation_id_seq OWNED BY public.workstation.workstation_id;

-- 3. workstation_id のデフォルト値を、作成したシーケンスの次の値 (nextval) に設定するのだ
ALTER TABLE ONLY public.workstation ALTER COLUMN workstation_id SET DEFAULT nextval('public.workstation_workstation_id_seq'::regclass);

-- (オプション) もし既に手動でデータを入れてしまっている場合は、シーケンスの現在値を最大IDに合わせる必要があるのだ
-- データが空ならこの行はなくてもエラーにはならないけど、念のため入れておくと安全なのだ
SELECT setval('public.workstation_workstation_id_seq', COALESCE((SELECT MAX(workstation_id) FROM public.workstation), 1), false);


-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

-- 1. デフォルト値を削除するのだ
ALTER TABLE ONLY public.workstation ALTER COLUMN workstation_id DROP DEFAULT;

-- 2. シーケンスを削除するのだ
DROP SEQUENCE IF EXISTS public.workstation_workstation_id_seq;
