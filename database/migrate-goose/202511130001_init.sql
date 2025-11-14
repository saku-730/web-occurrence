-- +goose Up
--
-- Name: postgis; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA public;


--
-- Name: attachment_goup; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.attachment_goup (
    occurrence_id integer NOT NULL,
    attachment_id integer NOT NULL,
    priority integer
);


--
-- Name: attachments; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.attachments (
    attachment_id integer NOT NULL,
    file_path text NOT NULL,
    extension_id integer,
    user_id integer
);


--
-- Name: attachments_attachment_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.attachments_attachment_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: attachments_attachment_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.attachments_attachment_id_seq OWNED BY public.attachments.attachment_id;


--
-- Name: change_logs; Type: TABLE; Schema: public; Owner: admin
--

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


--
-- Name: change_logs_log_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.change_logs_log_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: change_logs_log_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.change_logs_log_id_seq OWNED BY public.change_logs.log_id;


--
-- Name: classification_json; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.classification_json (
    classification_id integer NOT NULL,
    class_classification jsonb
);


--
-- Name: classification_json_classification_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.classification_json_classification_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: classification_json_classification_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.classification_json_classification_id_seq OWNED BY public.classification_json.classification_id;


--
-- Name: file_extensions; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.file_extensions (
    extension_id integer NOT NULL,
    extension_text character varying(255),
    file_type_id integer
);


--
-- Name: file_extensions_extension_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.file_extensions_extension_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: file_extensions_extension_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.file_extensions_extension_id_seq OWNED BY public.file_extensions.extension_id;


--
-- Name: file_types; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.file_types (
    file_type_id integer NOT NULL,
    type_name text
);


--
-- Name: file_types_file_type_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.file_types_file_type_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: file_types_file_type_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.file_types_file_type_id_seq OWNED BY public.file_types.file_type_id;


--
-- Name: identifications; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.identifications (
    identification_id integer NOT NULL,
    user_id integer,
    occurrence_id integer,
    source_info text,
    identificated_at timestamp with time zone DEFAULT now(),
    timezone smallint NOT NULL
);


--
-- Name: identifications_identification_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.identifications_identification_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: identifications_identification_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.identifications_identification_id_seq OWNED BY public.identifications.identification_id;


--
-- Name: languages; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.languages (
    language_id integer NOT NULL,
    language_short text,
    language_common text
);


--
-- Name: language_language_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.language_language_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: language_language_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.language_language_id_seq OWNED BY public.languages.language_id;


--
-- Name: make_specimen; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.make_specimen (
    make_specimen_id integer NOT NULL,
    occurrence_id integer,
    user_id integer,
    specimen_id integer,
    date date,
    specimen_method_id integer,
    created_at timestamp with time zone DEFAULT now(),
    timezone smallint NOT NULL
);


--
-- Name: make_specimen_make_specimen_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.make_specimen_make_specimen_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: make_specimen_make_specimen_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.make_specimen_make_specimen_id_seq OWNED BY public.make_specimen.make_specimen_id;


--
-- Name: observation_methods; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.observation_methods (
    observation_method_id integer NOT NULL,
    method_common_name text,
    pageid integer
);


--
-- Name: observation_methods_observation_method_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.observation_methods_observation_method_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: observation_methods_observation_method_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.observation_methods_observation_method_id_seq OWNED BY public.observation_methods.observation_method_id;


--
-- Name: observations; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.observations (
    observations_id integer NOT NULL,
    user_id integer,
    occurrence_id integer,
    observation_method_id integer,
    behavior text,
    observed_at timestamp with time zone DEFAULT now(),
    timezone smallint NOT NULL
);


--
-- Name: observations_observations_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.observations_observations_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: observations_observations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.observations_observations_id_seq OWNED BY public.observations.observations_id;


--
-- Name: occurrence; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.occurrence (
    occurrence_id integer NOT NULL,
    project_id integer,
    user_id integer,
    individual_id integer,
    lifestage text,
    sex text,
    classification_id integer,
    place_id integer,
    attachment_group_id integer,
    body_length numeric,
    language_id integer,
    note text,
    created_at timestamp with time zone DEFAULT now(),
    timezone smallint NOT NULL
);


--
-- Name: occurrence_occurrence_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.occurrence_occurrence_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: occurrence_occurrence_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.occurrence_occurrence_id_seq OWNED BY public.occurrence.occurrence_id;


--
-- Name: place_names_json; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.place_names_json (
    place_name_id integer NOT NULL,
    class_place_name jsonb
);


--
-- Name: place_names_json_place_name_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.place_names_json_place_name_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: place_names_json_place_name_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.place_names_json_place_name_id_seq OWNED BY public.place_names_json.place_name_id;


--
-- Name: places; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.places (
    place_id integer NOT NULL,
    coordinates public.geography(Point,4326),
    place_name_id integer,
    accuracy numeric
);


--
-- Name: places_place_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.places_place_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: places_place_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.places_place_id_seq OWNED BY public.places.place_id;


--
-- Name: project_members; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.project_members (
    project_member_id integer NOT NULL,
    project_id integer,
    user_id integer,
    join_day date,
    finish_day date
);


--
-- Name: project_members_project_member_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.project_members_project_member_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: project_members_project_member_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.project_members_project_member_id_seq OWNED BY public.project_members.project_member_id;


--
-- Name: projects; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.projects (
    project_id integer NOT NULL,
    project_name text NOT NULL,
    disscription text,
    start_day date,
    finished_day date,
    updated_day date,
    note text
);


--
-- Name: projects_project_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.projects_project_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: projects_project_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.projects_project_id_seq OWNED BY public.projects.project_id;


--
-- Name: specimen; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.specimen (
    specimen_id integer NOT NULL,
    occurrence_id integer,
    specimen_method_id integer,
    institution_id text,
    collection_id text
);


--
-- Name: specimen_methods; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.specimen_methods (
    specimen_methods_id integer NOT NULL,
    method_common_name text,
    page_id integer
);


--
-- Name: specimen_methods_specimen_methods_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.specimen_methods_specimen_methods_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: specimen_methods_specimen_methods_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.specimen_methods_specimen_methods_id_seq OWNED BY public.specimen_methods.specimen_methods_id;


--
-- Name: specimen_specimen_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.specimen_specimen_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: specimen_specimen_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.specimen_specimen_id_seq OWNED BY public.specimen.specimen_id;


--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.user_roles (
    role_id integer NOT NULL,
    role_name text NOT NULL
);


--
-- Name: user_roles_role_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.user_roles_role_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_roles_role_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.user_roles_role_id_seq OWNED BY public.user_roles.role_id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.users (
    user_id integer NOT NULL,
    user_name character varying(255) NOT NULL,
    display_name character varying(255) NOT NULL,
    mail_address character varying(255),
    password text,
    created_at timestamp with time zone DEFAULT now(),
    timezone smallint NOT NULL
);


--
-- Name: users_userid_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.users_userid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_userid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.users_userid_seq OWNED BY public.users.user_id;


--
-- Name: wiki_pages; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.wiki_pages (
    page_id integer NOT NULL,
    title text,
    user_id integer,
    created_date timestamp without time zone DEFAULT now(),
    updated_date timestamp without time zone,
    content_path text
);


--
-- Name: wiki_pages_page_id_seq; Type: SEQUENCE; Schema: public; Owner: admin
--

CREATE SEQUENCE public.wiki_pages_page_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: wiki_pages_page_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: admin
--

ALTER SEQUENCE public.wiki_pages_page_id_seq OWNED BY public.wiki_pages.page_id;


--
-- Name: workstation; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.workstation (
    workstation_id integer NOT NULL,
    workstation_name text NOT NULL
);


--
-- Name: workstation_user; Type: TABLE; Schema: public; Owner: admin
--

CREATE TABLE public.workstation_user (
    workstation_id integer NOT NULL,
    user_id integer NOT NULL,
    role_id integer NOT NULL
);


--
-- Name: attachments attachment_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.attachments ALTER COLUMN attachment_id SET DEFAULT nextval('public.attachments_attachment_id_seq'::regclass);


--
-- Name: change_logs log_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.change_logs ALTER COLUMN log_id SET DEFAULT nextval('public.change_logs_log_id_seq'::regclass);


--
-- Name: classification_json classification_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.classification_json ALTER COLUMN classification_id SET DEFAULT nextval('public.classification_json_classification_id_seq'::regclass);


--
-- Name: file_extensions extension_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.file_extensions ALTER COLUMN extension_id SET DEFAULT nextval('public.file_extensions_extension_id_seq'::regclass);


--
-- Name: file_types file_type_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.file_types ALTER COLUMN file_type_id SET DEFAULT nextval('public.file_types_file_type_id_seq'::regclass);


--
-- Name: identifications identification_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.identifications ALTER COLUMN identification_id SET DEFAULT nextval('public.identifications_identification_id_seq'::regclass);


--
-- Name: languages language_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.languages ALTER COLUMN language_id SET DEFAULT nextval('public.language_language_id_seq'::regclass);


--
-- Name: make_specimen make_specimen_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.make_specimen ALTER COLUMN make_specimen_id SET DEFAULT nextval('public.make_specimen_make_specimen_id_seq'::regclass);


--
-- Name: observation_methods observation_method_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.observation_methods ALTER COLUMN observation_method_id SET DEFAULT nextval('public.observation_methods_observation_method_id_seq'::regclass);


--
-- Name: observations observations_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.observations ALTER COLUMN observations_id SET DEFAULT nextval('public.observations_observations_id_seq'::regclass);


--
-- Name: occurrence occurrence_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.occurrence ALTER COLUMN occurrence_id SET DEFAULT nextval('public.occurrence_occurrence_id_seq'::regclass);


--
-- Name: place_names_json place_name_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.place_names_json ALTER COLUMN place_name_id SET DEFAULT nextval('public.place_names_json_place_name_id_seq'::regclass);


--
-- Name: places place_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.places ALTER COLUMN place_id SET DEFAULT nextval('public.places_place_id_seq'::regclass);


--
-- Name: project_members project_member_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.project_members ALTER COLUMN project_member_id SET DEFAULT nextval('public.project_members_project_member_id_seq'::regclass);


--
-- Name: projects project_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.projects ALTER COLUMN project_id SET DEFAULT nextval('public.projects_project_id_seq'::regclass);


--
-- Name: specimen specimen_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.specimen ALTER COLUMN specimen_id SET DEFAULT nextval('public.specimen_specimen_id_seq'::regclass);


--
-- Name: specimen_methods specimen_methods_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.specimen_methods ALTER COLUMN specimen_methods_id SET DEFAULT nextval('public.specimen_methods_specimen_methods_id_seq'::regclass);


--
-- Name: user_roles role_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.user_roles ALTER COLUMN role_id SET DEFAULT nextval('public.user_roles_role_id_seq'::regclass);


--
-- Name: users user_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users ALTER COLUMN user_id SET DEFAULT nextval('public.users_userid_seq'::regclass);


--
-- Name: wiki_pages page_id; Type: DEFAULT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.wiki_pages ALTER COLUMN page_id SET DEFAULT nextval('public.wiki_pages_page_id_seq'::regclass);


--
-- Name: attachment_goup attachment_goup_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.attachment_goup
    ADD CONSTRAINT attachment_goup_pkey PRIMARY KEY (occurrence_id, attachment_id);


--
-- Name: attachments attachments_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.attachments
    ADD CONSTRAINT attachments_pkey PRIMARY KEY (attachment_id);


--
-- Name: change_logs change_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.change_logs
    ADD CONSTRAINT change_logs_pkey PRIMARY KEY (log_id);


--
-- Name: classification_json classification_json_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.classification_json
    ADD CONSTRAINT classification_json_pkey PRIMARY KEY (classification_id);


--
-- Name: file_extensions file_extensions_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.file_extensions
    ADD CONSTRAINT file_extensions_pkey PRIMARY KEY (extension_id);


--
-- Name: file_types file_types_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.file_types
    ADD CONSTRAINT file_types_pkey PRIMARY KEY (file_type_id);


--
-- Name: identifications identifications_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.identifications
    ADD CONSTRAINT identifications_pkey PRIMARY KEY (identification_id);


--
-- Name: languages language_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.languages
    ADD CONSTRAINT language_pkey PRIMARY KEY (language_id);


--
-- Name: make_specimen make_specimen_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.make_specimen
    ADD CONSTRAINT make_specimen_pkey PRIMARY KEY (make_specimen_id);


--
-- Name: observation_methods observation_methods_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.observation_methods
    ADD CONSTRAINT observation_methods_pkey PRIMARY KEY (observation_method_id);


--
-- Name: observations observations_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.observations
    ADD CONSTRAINT observations_pkey PRIMARY KEY (observations_id);


--
-- Name: occurrence occurrence_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.occurrence
    ADD CONSTRAINT occurrence_pkey PRIMARY KEY (occurrence_id);


--
-- Name: place_names_json place_names_json_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.place_names_json
    ADD CONSTRAINT place_names_json_pkey PRIMARY KEY (place_name_id);


--
-- Name: places places_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.places
    ADD CONSTRAINT places_pkey PRIMARY KEY (place_id);


--
-- Name: project_members project_members_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.project_members
    ADD CONSTRAINT project_members_pkey PRIMARY KEY (project_member_id);


--
-- Name: projects projects_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.projects
    ADD CONSTRAINT projects_pkey PRIMARY KEY (project_id);


--
-- Name: specimen_methods specimen_methods_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.specimen_methods
    ADD CONSTRAINT specimen_methods_pkey PRIMARY KEY (specimen_methods_id);


--
-- Name: specimen specimen_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.specimen
    ADD CONSTRAINT specimen_pkey PRIMARY KEY (specimen_id);


--
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_pkey PRIMARY KEY (role_id);


--
-- Name: user_roles user_roles_role_name_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.user_roles
    ADD CONSTRAINT user_roles_role_name_key UNIQUE (role_name);


--
-- Name: users users_mail_address_key; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_mail_address_key UNIQUE (mail_address);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);


--
-- Name: wiki_pages wiki_pages_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.wiki_pages
    ADD CONSTRAINT wiki_pages_pkey PRIMARY KEY (page_id);


--
-- Name: workstation workstation_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.workstation
    ADD CONSTRAINT workstation_pkey PRIMARY KEY (workstation_id);


--
-- Name: workstation_user workstation_user_pkey; Type: CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.workstation_user
    ADD CONSTRAINT workstation_user_pkey PRIMARY KEY (user_id, workstation_id);


--
-- Name: attachment_goup attachment_goup_attachment_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.attachment_goup
    ADD CONSTRAINT attachment_goup_attachment_id_fkey FOREIGN KEY (attachment_id) REFERENCES public.attachments(attachment_id);


--
-- Name: attachment_goup attachment_goup_occurrence_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.attachment_goup
    ADD CONSTRAINT attachment_goup_occurrence_id_fkey FOREIGN KEY (occurrence_id) REFERENCES public.occurrence(occurrence_id);


--
-- Name: attachments attachments_extension_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.attachments
    ADD CONSTRAINT attachments_extension_id_fkey FOREIGN KEY (extension_id) REFERENCES public.file_extensions(extension_id);


--
-- Name: attachments attachments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.attachments
    ADD CONSTRAINT attachments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- Name: change_logs change_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.change_logs
    ADD CONSTRAINT change_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- Name: file_extensions file_extensions_file_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.file_extensions
    ADD CONSTRAINT file_extensions_file_type_id_fkey FOREIGN KEY (file_type_id) REFERENCES public.file_types(file_type_id);


--
-- Name: identifications identifications_occurrence_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.identifications
    ADD CONSTRAINT identifications_occurrence_id_fkey FOREIGN KEY (occurrence_id) REFERENCES public.occurrence(occurrence_id);


--
-- Name: identifications identifications_userid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.identifications
    ADD CONSTRAINT identifications_userid_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- Name: make_specimen make_specimen_occurrence_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.make_specimen
    ADD CONSTRAINT make_specimen_occurrence_id_fkey FOREIGN KEY (occurrence_id) REFERENCES public.occurrence(occurrence_id);


--
-- Name: make_specimen make_specimen_specimen_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.make_specimen
    ADD CONSTRAINT make_specimen_specimen_id_fkey FOREIGN KEY (specimen_id) REFERENCES public.specimen(specimen_id);


--
-- Name: make_specimen make_specimen_specimen_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.make_specimen
    ADD CONSTRAINT make_specimen_specimen_method_id_fkey FOREIGN KEY (specimen_method_id) REFERENCES public.specimen_methods(specimen_methods_id);


--
-- Name: make_specimen make_specimen_userid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.make_specimen
    ADD CONSTRAINT make_specimen_userid_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- Name: observation_methods observation_methods_pageid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.observation_methods
    ADD CONSTRAINT observation_methods_pageid_fkey FOREIGN KEY (pageid) REFERENCES public.wiki_pages(page_id);


--
-- Name: observations observations_observation_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.observations
    ADD CONSTRAINT observations_observation_method_id_fkey FOREIGN KEY (observation_method_id) REFERENCES public.observation_methods(observation_method_id);


--
-- Name: observations observations_occurrence_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.observations
    ADD CONSTRAINT observations_occurrence_id_fkey FOREIGN KEY (occurrence_id) REFERENCES public.occurrence(occurrence_id);


--
-- Name: observations observations_userid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.observations
    ADD CONSTRAINT observations_userid_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- Name: occurrence occurrence_classification_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.occurrence
    ADD CONSTRAINT occurrence_classification_id_fkey FOREIGN KEY (classification_id) REFERENCES public.classification_json(classification_id);


--
-- Name: occurrence occurrence_language_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.occurrence
    ADD CONSTRAINT occurrence_language_id_fkey FOREIGN KEY (language_id) REFERENCES public.languages(language_id);


--
-- Name: occurrence occurrence_place_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.occurrence
    ADD CONSTRAINT occurrence_place_id_fkey FOREIGN KEY (place_id) REFERENCES public.places(place_id);


--
-- Name: occurrence occurrence_project_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.occurrence
    ADD CONSTRAINT occurrence_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.projects(project_id);


--
-- Name: occurrence occurrence_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.occurrence
    ADD CONSTRAINT occurrence_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- Name: places places_place_name_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.places
    ADD CONSTRAINT places_place_name_id_fkey FOREIGN KEY (place_name_id) REFERENCES public.place_names_json(place_name_id);


--
-- Name: project_members project_members_project_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.project_members
    ADD CONSTRAINT project_members_project_id_fkey FOREIGN KEY (project_id) REFERENCES public.projects(project_id);


--
-- Name: project_members project_members_userid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.project_members
    ADD CONSTRAINT project_members_userid_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- Name: workstation_user role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.workstation_user
    ADD CONSTRAINT role_id_fkey FOREIGN KEY (role_id) REFERENCES public.user_roles(role_id) NOT VALID;


--
-- Name: specimen_methods specimen_methods_page_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.specimen_methods
    ADD CONSTRAINT specimen_methods_page_id_fkey FOREIGN KEY (page_id) REFERENCES public.wiki_pages(page_id);


--
-- Name: specimen specimen_occurrence_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.specimen
    ADD CONSTRAINT specimen_occurrence_id_fkey FOREIGN KEY (occurrence_id) REFERENCES public.occurrence(occurrence_id);


--
-- Name: specimen specimen_specimen_method_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.specimen
    ADD CONSTRAINT specimen_specimen_method_id_fkey FOREIGN KEY (specimen_method_id) REFERENCES public.specimen_methods(specimen_methods_id);


--
-- Name: workstation_user user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.workstation_user
    ADD CONSTRAINT user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id) NOT VALID;


--
-- Name: wiki_pages wiki_pages_userid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.wiki_pages
    ADD CONSTRAINT wiki_pages_userid_fkey FOREIGN KEY (user_id) REFERENCES public.users(user_id);


--
-- Name: workstation_user workstation_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: admin
--

ALTER TABLE ONLY public.workstation_user
    ADD CONSTRAINT workstation_id_fkey FOREIGN KEY (workstation_id) REFERENCES public.workstation(workstation_id) NOT VALID;


-- +goose Down
-- スキーマ全体を削除するDownマイグレーション（DROP TABLEなど）は、
-- 依存関係（FK）の逆順に書く必要があって複雑だから、ここでは省略するのだ。
-- もし必要なら、テーブルをすべてDROPする処理をここに書くのだ。
-- 例:
-- DROP TABLE public.workstation_user;
-- DROP TABLE public.workstation;
-- ... (すべてのテーブル、シーケンス、拡張機能を逆順に) ...
-- DROP EXTENSION IF EXISTS postgis;
