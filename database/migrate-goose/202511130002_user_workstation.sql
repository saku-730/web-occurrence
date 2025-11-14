-- +goose Up
--
-- workstation_id が足りないテーブルに追加 (FK設定も含む)
--
ALTER TABLE public.attachment_goup ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.attachments ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.change_logs ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.classification_json ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.identifications ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.make_specimen ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.observation_methods ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.observations ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.occurrence ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.place_names_json ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.places ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.project_members ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.projects ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.specimen ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.specimen_methods ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.users ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);
ALTER TABLE public.wiki_pages ADD COLUMN workstation_id INTEGER REFERENCES public.workstation(workstation_id);

--
-- user_id が足りないテーブルに追加 (FK設定も含む)
--
ALTER TABLE public.attachment_goup ADD COLUMN user_id INTEGER REFERENCES public.users(user_id);
ALTER TABLE public.classification_json ADD COLUMN user_id INTEGER REFERENCES public.users(user_id);
ALTER TABLE public.observation_methods ADD COLUMN user_id INTEGER REFERENCES public.users(user_id);
ALTER TABLE public.place_names_json ADD COLUMN user_id INTEGER REFERENCES public.users(user_id);
ALTER TABLE public.places ADD COLUMN user_id INTEGER REFERENCES public.users(user_id);
ALTER TABLE public.projects ADD COLUMN user_id INTEGER REFERENCES public.users(user_id);
ALTER TABLE public.specimen ADD COLUMN user_id INTEGER REFERENCES public.users(user_id);
ALTER TABLE public.specimen_methods ADD COLUMN user_id INTEGER REFERENCES public.users(user_id);

-- +goose Down
--
-- workstation_id を削除
--
ALTER TABLE public.attachment_goup DROP COLUMN workstation_id;
ALTER TABLE public.attachments DROP COLUMN workstation_id;
ALTER TABLE public.change_logs DROP COLUMN workstation_id;
ALTER TABLE public.classification_json DROP COLUMN workstation_id;
ALTER TABLE public.identifications DROP COLUMN workstation_id;
ALTER TABLE public.make_specimen DROP COLUMN workstation_id;
ALTER TABLE public.observation_methods DROP COLUMN workstation_id;
ALTER TABLE public.observations DROP COLUMN workstation_id;
ALTER TABLE public.occurrence DROP COLUMN workstation_id;
ALTER TABLE public.place_names_json DROP COLUMN workstation_id;
ALTER TABLE public.places DROP COLUMN workstation_id;
ALTER TABLE public.project_members DROP COLUMN workstation_id;
ALTER TABLE public.projects DROP COLUMN workstation_id;
ALTER TABLE public.specimen DROP COLUMN workstation_id;
ALTER TABLE public.specimen_methods DROP COLUMN workstation_id;
ALTER TABLE public.users DROP COLUMN workstation_id;
ALTER TABLE public.wiki_pages DROP COLUMN workstation_id;

--
-- user_id を削除
--
ALTER TABLE public.attachment_goup DROP COLUMN user_id;
ALTER TABLE public.classification_json DROP COLUMN user_id;
ALTER TABLE public.observation_methods DROP COLUMN user_id;
ALTER TABLE public.place_names_json DROP COLUMN user_id;
ALTER TABLE public.places DROP COLUMN user_id;
ALTER TABLE public.projects DROP COLUMN user_id;
ALTER TABLE public.specimen DROP COLUMN user_id;
ALTER TABLE public.specimen_methods DROP COLUMN user_id;
