-- +goose Up
-- 既存のテーブルを削除して、UUID対応のスキーマで再作成するのだ
-- 開発環境用なのでデータをリセットする強硬策をとるのだ

DROP TABLE IF EXISTS attachment_group CASCADE;
DROP TABLE IF EXISTS attachments CASCADE;
DROP TABLE IF EXISTS observations CASCADE;
DROP TABLE IF EXISTS make_specimen CASCADE;
DROP TABLE IF EXISTS specimen CASCADE;
DROP TABLE IF EXISTS identifications CASCADE;
DROP TABLE IF EXISTS occurrence CASCADE;
DROP TABLE IF EXISTS places CASCADE;
DROP TABLE IF EXISTS place_names_json CASCADE;
DROP TABLE IF EXISTS classification_json CASCADE;
DROP TABLE IF EXISTS workstation_user CASCADE;
DROP TABLE IF EXISTS workstation CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- 1. Users
CREATE TABLE users (
    user_id text PRIMARY KEY, -- UUID
    user_name text NOT NULL,
    display_name text,
    mail_address text UNIQUE,
    password text,
    created_at timestamp with time zone DEFAULT now()
);

-- 2. Workstation
CREATE TABLE workstation (
    workstation_id text PRIMARY KEY, -- UUID
    workstation_name text NOT NULL
);

-- 3. Workstation User (Many-to-Many)
CREATE TABLE workstation_user (
    workstation_id text REFERENCES workstation(workstation_id) ON DELETE CASCADE,
    user_id text REFERENCES users(user_id) ON DELETE CASCADE,
    role_id integer DEFAULT 1,
    PRIMARY KEY (workstation_id, user_id)
);

-- 4. Classification
CREATE TABLE classification_json (
    classification_id text PRIMARY KEY,
    class_classification jsonb -- {"kingdom": "...", "species": "..."}
);

-- 5. Place Names
CREATE TABLE place_names_json (
    place_name_id text PRIMARY KEY,
    class_place_name jsonb -- {"ja": "...", "en": "..."}
);

-- 6. Places
CREATE TABLE places (
    place_id text PRIMARY KEY,
    place_name_id text REFERENCES place_names_json(place_name_id),
    coordinates jsonb, -- PostGISを使わない簡易版としてJSONBで保存
    accuracy numeric
);

-- 7. Occurrence (Main Table)
CREATE TABLE occurrence (
    occurrence_id text PRIMARY KEY,
    workstation_id text REFERENCES workstation(workstation_id),
    user_id text REFERENCES users(user_id), -- created_by
    project_id text,
    
    individual_id text,
    lifestage text,
    sex text,
    body_length numeric,
    note text,
    
    classification_id text REFERENCES classification_json(classification_id),
    place_id text REFERENCES places(place_id),
    language_id text,

    created_at timestamp with time zone,
    timezone text
);

-- 8. Identifications
CREATE TABLE identifications (
    identification_id text PRIMARY KEY,
    occurrence_id text REFERENCES occurrence(occurrence_id) ON DELETE CASCADE,
    user_id text REFERENCES users(user_id),
    source_info text,
    identificated_at timestamp with time zone
);

-- 9. Specimen
CREATE TABLE specimen (
    specimen_id text PRIMARY KEY,
    occurrence_id text REFERENCES occurrence(occurrence_id) ON DELETE CASCADE,
    institution_id text,
    collection_id text,
    specimen_method_id text
);

-- 10. Make Specimen (Specimen作成情報)
CREATE TABLE make_specimen (
    make_specimen_id text PRIMARY KEY,
    specimen_id text REFERENCES specimen(specimen_id) ON DELETE CASCADE,
    user_id text REFERENCES users(user_id),
    created_at timestamp with time zone
);

-- 11. Observations
CREATE TABLE observations (
    observation_id text PRIMARY KEY,
    occurrence_id text REFERENCES occurrence(occurrence_id) ON DELETE CASCADE,
    user_id text REFERENCES users(user_id),
    observation_method_id text,
    behavior text,
    observed_at timestamp with time zone
);

-- 12. Attachments (今回は簡易化のためテーブル定義のみ)
CREATE TABLE attachments (
    attachment_id text PRIMARY KEY,
    file_path text,
    user_id text REFERENCES users(user_id)
);

CREATE TABLE attachment_group (
    occurrence_id text REFERENCES occurrence(occurrence_id) ON DELETE CASCADE,
    attachment_id text REFERENCES attachments(attachment_id) ON DELETE CASCADE,
    priority integer,
    PRIMARY KEY (occurrence_id, attachment_id)
);

-- +goose Down
DROP TABLE IF EXISTS attachment_group CASCADE;
DROP TABLE IF EXISTS attachments CASCADE;
DROP TABLE IF EXISTS observations CASCADE;
DROP TABLE IF EXISTS make_specimen CASCADE;
DROP TABLE IF EXISTS specimen CASCADE;
DROP TABLE IF EXISTS identifications CASCADE;
DROP TABLE IF EXISTS occurrence CASCADE;
DROP TABLE IF EXISTS places CASCADE;
DROP TABLE IF EXISTS place_names_json CASCADE;
DROP TABLE IF EXISTS classification_json CASCADE;
DROP TABLE IF EXISTS workstation_user CASCADE;
DROP TABLE IF EXISTS workstation CASCADE;
DROP TABLE IF EXISTS users CASCADE;
