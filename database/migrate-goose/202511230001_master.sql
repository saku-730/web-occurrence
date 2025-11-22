-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- 1. 既存データを全削除（スキーマは維持）
-- CASCADE を付けることで、外部キーで繋がっているテーブルも全部まとめて消すのだ
TRUNCATE TABLE 
    users,
    projects,
    project_members,
    occurrence,
    observations,
    observation_methods,
    specimen,
    make_specimen,
    specimen_methods,
    identifications,
    places,
    place_names_json,
    classification_json,
    attachments,
    attachment_goup,
    wiki_pages,
    workstation,
    workstation_user,
    languages,
    file_types,
    file_extensions,
    user_roles
RESTART IDENTITY CASCADE;

-- 2. マスターデータの投入

-- 言語 (Languages)
INSERT INTO public.languages (language_id, language_short, language_common) VALUES
(1, 'ja', '日本語'),
(2, 'en', 'English');

-- ユーザーロール (User Roles)
INSERT INTO public.user_roles (role_id, role_name) VALUES
(1, 'administrator'),
(2, 'editor'),
(3, 'viewer');

-- ファイル種別 (File Types)
INSERT INTO public.file_types (file_type_id, type_name) VALUES
(1, 'Image'),
(2, 'Audio'),
(3, 'Video'),
(4, 'Document');

-- ファイル拡張子 (File Extensions)
INSERT INTO public.file_extensions (extension_text, file_type_id) VALUES
('jpg', 1),
('jpeg', 1),
('png', 1),
('gif', 1),
('tiff',1),
('webp',1),
('mp3', 2),
('wav', 2),
('mp4', 3),
('mov', 3),
('pdf', 4),
('txt', 4),
('csv', 4);

-- デフォルトのワークステーション (Workstation)
-- ※これがないとユーザー登録時の紐付けで困るかもしれないので1つ作っておくのだ
INSERT INTO public.workstation (workstation_id, workstation_name) VALUES
(1, 'Default Workstation');


-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

-- ロールバック時は、今回追加したマスターデータを削除するのだ
-- (ユーザーが追加したデータも消えてしまうけど、Seedのリセットという意味ではこれでOKなのだ)
TRUNCATE TABLE 
    languages,
    file_types,
    file_extensions,
    user_roles,
    workstation
RESTART IDENTITY CASCADE;
