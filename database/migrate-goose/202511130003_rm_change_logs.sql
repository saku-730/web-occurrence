-- +goose Up
-- change_logs テーブルを削除するのだ。
-- これで、関連するシーケンス、PK、FK（change_logs_user_id_fkey）も自動的に削除されるのだ。
DROP TABLE public.change_logs;

-- +goose Down
-- もしロールバック（down）したら、テーブルと関連設定をすべて元通り再作成するのだ。

-- 1. テーブルの作成
CREATE TABLE public.change_logs (
    log_id integer NOT NULL,
    type text,
    changed_id integer,
    before_value text,
    after_value text,
    user_id integer,
    date timestamp without time zone DEFAULT now(),
    "row" text
);

-- 2. シーケンスの作成と所有権の設定
CREATE SEQUENCE public.change_logs_log_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.change_logs_log_id_seq OWNED BY public.change_logs.log_id;

-- 3. デフォルト値（シーケンス）の設定
ALTER TABLE ONLY public.change_logs ALTER COLUMN log_id SET DEFAULT nextval('public.change_logs_log_id_seq'::regclass);

-- 4. 主キー（PK）制約の復元
ALTER TABLE ONLY public.change_logs
    ADD CONSTRAINT change_logs_pkey PRIMARY KEY (log_id);

-- 5. 外部キー（FK）制約の復元
ALTER TABLE ONLY public.change_logs
    ADD CONSTRAINT change_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);
