CREATE USER api_user WITH PASSWORD '123456';
CREATE DATABASE reports OWNER api_user;

CREATE TABLE public.time_session_partial
(
    time_session_partial_id uuid NOT NULL,
    time_session_name character varying(128) COLLATE pg_catalog."default" NOT NULL,
    time_session_partial_start timestamp with time zone NOT NULL,
    time_session_partial_end timestamp with time zone,
    time_session_id uuid NOT NULL,
    time_session_completed boolean,
    user_id uuid NOT NULL,
    CONSTRAINT time_session_partial_pkey PRIMARY KEY (time_session_partial_id)
)

ALTER TABLE public.time_session_partial
    OWNER to api_user;
