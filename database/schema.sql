--
-- PostgreSQL database dump
--

\restrict hDZ1dpD8q0IlCBH3HwVEBZq4SqgpAtG8FrAxg5C4AhAbPYG1KkTGBJC8rAQEDi0

-- Dumped from database version 16.10 (Ubuntu 16.10-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.10 (Ubuntu 16.10-0ubuntu0.24.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: file_types; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.file_types (file_type_id, type_name) FROM stdin;
\.


--
-- Data for Name: file_extensions; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.file_extensions (extension_id, extension_text, file_type_id) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.users (user_id, user_name, display_name, mail_address, password, created_at, timezone) FROM stdin;
\.


--
-- Data for Name: workstation; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.workstation (workstation_id, workstation_name) FROM stdin;
\.


--
-- Data for Name: attachments; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.attachments (attachment_id, file_path, extension_id, user_id, workstation_id) FROM stdin;
\.


--
-- Data for Name: classification_json; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.classification_json (classification_id, class_classification, workstation_id, user_id) FROM stdin;
\.


--
-- Data for Name: languages; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.languages (language_id, language_short, language_common) FROM stdin;
\.


--
-- Data for Name: place_names_json; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.place_names_json (place_name_id, class_place_name, workstation_id, user_id) FROM stdin;
\.


--
-- Data for Name: places; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.places (place_id, coordinates, place_name_id, accuracy, workstation_id, user_id) FROM stdin;
\.


--
-- Data for Name: projects; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.projects (project_id, project_name, disscription, start_day, finished_day, updated_day, note, workstation_id, user_id) FROM stdin;
\.


--
-- Data for Name: occurrence; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.occurrence (occurrence_id, project_id, user_id, individual_id, lifestage, sex, classification_id, place_id, attachment_group_id, body_length, language_id, note, created_at, timezone, workstation_id) FROM stdin;
\.


--
-- Data for Name: attachment_goup; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.attachment_goup (occurrence_id, attachment_id, priority, workstation_id, user_id) FROM stdin;
\.


--
-- Data for Name: goose_db_version; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.goose_db_version (id, version_id, is_applied, tstamp) FROM stdin;
1	0	t	2025-11-13 17:36:22.100682
2	202511130001	t	2025-11-13 17:36:32.397668
3	202511130002	t	2025-11-13 17:48:34.745577
4	202511130003	t	2025-11-13 18:53:09.392494
5	202511130004	t	2025-11-13 19:03:09.523964
6	202511130005	t	2025-11-13 19:32:09.374519
\.


--
-- Data for Name: identifications; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.identifications (identification_id, user_id, occurrence_id, source_info, identificated_at, timezone, workstation_id) FROM stdin;
\.


--
-- Data for Name: wiki_pages; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.wiki_pages (page_id, title, user_id, created_date, updated_date, content_path, workstation_id) FROM stdin;
\.


--
-- Data for Name: specimen_methods; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.specimen_methods (specimen_methods_id, method_common_name, page_id, workstation_id, user_id) FROM stdin;
\.


--
-- Data for Name: specimen; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.specimen (specimen_id, occurrence_id, specimen_method_id, institution_id, collection_id, workstation_id, user_id) FROM stdin;
\.


--
-- Data for Name: make_specimen; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.make_specimen (make_specimen_id, occurrence_id, user_id, specimen_id, specimen_method_id, created_at, timezone, workstation_id) FROM stdin;
\.


--
-- Data for Name: observation_methods; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.observation_methods (observation_method_id, method_common_name, pageid, workstation_id, user_id) FROM stdin;
\.


--
-- Data for Name: observations; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.observations (observations_id, user_id, occurrence_id, observation_method_id, behavior, observed_at, timezone, workstation_id) FROM stdin;
\.


--
-- Data for Name: project_members; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.project_members (project_member_id, project_id, user_id, join_day, finish_day, workstation_id) FROM stdin;
\.


--
-- Data for Name: spatial_ref_sys; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.spatial_ref_sys (srid, auth_name, auth_srid, srtext, proj4text) FROM stdin;
\.


--
-- Data for Name: user_roles; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.user_roles (role_id, role_name) FROM stdin;
\.


--
-- Data for Name: workstation_user; Type: TABLE DATA; Schema: public; Owner: admin
--

COPY public.workstation_user (workstation_id, user_id, role_id) FROM stdin;
\.


--
-- Name: attachments_attachment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.attachments_attachment_id_seq', 1, false);


--
-- Name: classification_json_classification_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.classification_json_classification_id_seq', 1, false);


--
-- Name: file_extensions_extension_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.file_extensions_extension_id_seq', 1, false);


--
-- Name: file_types_file_type_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.file_types_file_type_id_seq', 1, false);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.goose_db_version_id_seq', 6, true);


--
-- Name: identifications_identification_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.identifications_identification_id_seq', 1, false);


--
-- Name: language_language_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.language_language_id_seq', 1, false);


--
-- Name: make_specimen_make_specimen_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.make_specimen_make_specimen_id_seq', 1, false);


--
-- Name: observation_methods_observation_method_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.observation_methods_observation_method_id_seq', 1, false);


--
-- Name: observations_observations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.observations_observations_id_seq', 1, false);


--
-- Name: occurrence_occurrence_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.occurrence_occurrence_id_seq', 1, false);


--
-- Name: place_names_json_place_name_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.place_names_json_place_name_id_seq', 1, false);


--
-- Name: places_place_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.places_place_id_seq', 1, false);


--
-- Name: project_members_project_member_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.project_members_project_member_id_seq', 1, false);


--
-- Name: projects_project_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.projects_project_id_seq', 1, false);


--
-- Name: specimen_methods_specimen_methods_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.specimen_methods_specimen_methods_id_seq', 1, false);


--
-- Name: specimen_specimen_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.specimen_specimen_id_seq', 1, false);


--
-- Name: user_roles_role_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.user_roles_role_id_seq', 1, false);


--
-- Name: users_userid_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.users_userid_seq', 1, false);


--
-- Name: wiki_pages_page_id_seq; Type: SEQUENCE SET; Schema: public; Owner: admin
--

SELECT pg_catalog.setval('public.wiki_pages_page_id_seq', 1, false);


--
-- PostgreSQL database dump complete
--

\unrestrict hDZ1dpD8q0IlCBH3HwVEBZq4SqgpAtG8FrAxg5C4AhAbPYG1KkTGBJC8rAQEDi0

