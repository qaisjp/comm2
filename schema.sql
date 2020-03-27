--
-- PostgreSQL database dump
--

-- Dumped from database version 12.1
-- Dumped by pg_dump version 12.1

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: resource_collaborators; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.resource_collaborators (
    resource_id integer NOT NULL,
    user_id integer NOT NULL,
    accepted boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: resource_comments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.resource_comments (
    id integer NOT NULL,
    resource_id integer NOT NULL,
    author_id integer,
    message text NOT NULL,
    deleted boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    edited_at timestamp without time zone DEFAULT now() NOT NULL,
    deleted_at timestamp without time zone DEFAULT now() NOT NULL,
    parent_comment_id integer
);


--
-- Name: resource_comments_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.resource_comments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: resource_comments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.resource_comments_id_seq OWNED BY public.resource_comments.id;


--
-- Name: resource_followings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.resource_followings (
    resource_id integer NOT NULL,
    user_id integer NOT NULL,
    type integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone NOT NULL
);


--
-- Name: resource_media; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.resource_media (
    id integer NOT NULL,
    resource_id integer NOT NULL,
    title text NOT NULL,
    description text DEFAULT ''::text NOT NULL,
    filename text NOT NULL,
    author_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: resource_media_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.resource_media_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: resource_media_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.resource_media_id_seq OWNED BY public.resource_media.id;


--
-- Name: resource_packages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.resource_packages (
    id integer NOT NULL,
    resource_id integer NOT NULL,
    filename text NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    description text DEFAULT ''::text NOT NULL,
    version character varying(10) NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    author_id integer NOT NULL,
    draft boolean DEFAULT true NOT NULL
);


--
-- Name: resource_packages_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.resource_packages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: resource_packages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.resource_packages_id_seq OWNED BY public.resource_packages.id;


--
-- Name: resource_votes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.resource_votes (
    user_id integer NOT NULL,
    resource_id integer NOT NULL,
    positive boolean NOT NULL
);


--
-- Name: resources; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.resources (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    author_id integer NOT NULL,
    name character varying(254) NOT NULL,
    title text NOT NULL,
    description text NOT NULL,
    status text DEFAULT 'public'::text NOT NULL
);


--
-- Name: resources_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.resources_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: resources_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.resources_id_seq OWNED BY public.resources.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: user_bans; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_bans (
    id integer NOT NULL,
    author_id integer NOT NULL,
    target_user_id integer NOT NULL,
    reason text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    expires_at timestamp without time zone DEFAULT now(),
    active boolean DEFAULT true NOT NULL
);


--
-- Name: user_bans_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_bans_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_bans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_bans_id_seq OWNED BY public.user_bans.id;


--
-- Name: user_followings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_followings (
    source_user_id integer NOT NULL,
    target_user_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: user_profile; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_profile (
    user_id integer NOT NULL,
    location character varying(255),
    organisation character varying(255),
    website character varying(255),
    bio character varying(255)
);


--
-- Name: user_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_tokens (
    id character(20) NOT NULL,
    user_id integer NOT NULL,
    type integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: COLUMN user_tokens.type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.user_tokens.type IS 'activation token? password reset token?';


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    username character varying(254) NOT NULL,
    password character(60) NOT NULL,
    email character varying(254) NOT NULL,
    is_activated boolean DEFAULT false NOT NULL,
    level integer DEFAULT 0 NOT NULL,
    is_banned boolean DEFAULT false NOT NULL
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: resource_comments id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_comments ALTER COLUMN id SET DEFAULT nextval('public.resource_comments_id_seq'::regclass);


--
-- Name: resource_media id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_media ALTER COLUMN id SET DEFAULT nextval('public.resource_media_id_seq'::regclass);


--
-- Name: resource_packages id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_packages ALTER COLUMN id SET DEFAULT nextval('public.resource_packages_id_seq'::regclass);


--
-- Name: resources id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resources ALTER COLUMN id SET DEFAULT nextval('public.resources_id_seq'::regclass);


--
-- Name: user_bans id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_bans ALTER COLUMN id SET DEFAULT nextval('public.user_bans_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: resource_collaborators resource_collaborators_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_collaborators
    ADD CONSTRAINT resource_collaborators_pk PRIMARY KEY (resource_id, user_id);


--
-- Name: resource_comments resource_comments_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_comments
    ADD CONSTRAINT resource_comments_pk PRIMARY KEY (id);


--
-- Name: resource_followings resource_followings_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_followings
    ADD CONSTRAINT resource_followings_pk PRIMARY KEY (resource_id, user_id);


--
-- Name: resource_media resource_media_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_media
    ADD CONSTRAINT resource_media_pk PRIMARY KEY (id);


--
-- Name: resource_packages resource_packages_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_packages
    ADD CONSTRAINT resource_packages_pk PRIMARY KEY (id);


--
-- Name: resource_votes resource_votes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_votes
    ADD CONSTRAINT resource_votes_pkey PRIMARY KEY (user_id, resource_id);


--
-- Name: resources resources_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: user_bans user_bans_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_bans
    ADD CONSTRAINT user_bans_pk PRIMARY KEY (id);


--
-- Name: user_followings user_followings_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_followings
    ADD CONSTRAINT user_followings_pk PRIMARY KEY (source_user_id, target_user_id);


--
-- Name: user_profile user_profile_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_profile
    ADD CONSTRAINT user_profile_pk PRIMARY KEY (user_id);


--
-- Name: user_tokens user_tokens_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_tokens
    ADD CONSTRAINT user_tokens_pk PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: resources_name_author_id_uindex; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX resources_name_author_id_uindex ON public.resources USING btree (name, author_id);


--
-- Name: resource_collaborators resource_collaborators_resources_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_collaborators
    ADD CONSTRAINT resource_collaborators_resources_id_fk FOREIGN KEY (resource_id) REFERENCES public.resources(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resource_collaborators resource_collaborators_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_collaborators
    ADD CONSTRAINT resource_collaborators_users_id_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resource_comments resource_comments_resource_comments_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_comments
    ADD CONSTRAINT resource_comments_resource_comments_id_fk FOREIGN KEY (parent_comment_id) REFERENCES public.resource_comments(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resource_comments resource_comments_resources_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_comments
    ADD CONSTRAINT resource_comments_resources_id_fk FOREIGN KEY (resource_id) REFERENCES public.resources(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resource_comments resource_comments_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_comments
    ADD CONSTRAINT resource_comments_users_id_fk FOREIGN KEY (author_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resource_followings resource_followings_resources_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_followings
    ADD CONSTRAINT resource_followings_resources_id_fk FOREIGN KEY (resource_id) REFERENCES public.resources(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resource_followings resource_followings_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_followings
    ADD CONSTRAINT resource_followings_users_id_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resource_media resource_media_resources_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_media
    ADD CONSTRAINT resource_media_resources_id_fk FOREIGN KEY (resource_id) REFERENCES public.resources(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resource_media resource_media_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_media
    ADD CONSTRAINT resource_media_users_id_fk FOREIGN KEY (author_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resource_packages resource_packages_resources_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_packages
    ADD CONSTRAINT resource_packages_resources_id_fk FOREIGN KEY (resource_id) REFERENCES public.resources(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resource_packages resource_packages_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_packages
    ADD CONSTRAINT resource_packages_users_id_fk FOREIGN KEY (author_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resource_votes resource_votes_resource_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_votes
    ADD CONSTRAINT resource_votes_resource_fkey FOREIGN KEY (resource_id) REFERENCES public.resources(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resource_votes resource_votes_user_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resource_votes
    ADD CONSTRAINT resource_votes_user_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: resources resources_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_users_id_fk FOREIGN KEY (author_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_bans user_bans_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_bans
    ADD CONSTRAINT user_bans_users_id_fk FOREIGN KEY (author_id) REFERENCES public.users(id) ON UPDATE RESTRICT ON DELETE RESTRICT;


--
-- Name: user_bans user_bans_users_id_fk_2; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_bans
    ADD CONSTRAINT user_bans_users_id_fk_2 FOREIGN KEY (target_user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_followings user_followings_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_followings
    ADD CONSTRAINT user_followings_users_id_fk FOREIGN KEY (target_user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_followings user_followings_users_id_fk_2; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_followings
    ADD CONSTRAINT user_followings_users_id_fk_2 FOREIGN KEY (source_user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_profile user_profile_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_profile
    ADD CONSTRAINT user_profile_users_id_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: user_tokens user_tokens_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_tokens
    ADD CONSTRAINT user_tokens_users_id_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 12.1
-- Dumped by pg_dump version 12.1

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
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.schema_migrations (version, dirty) FROM stdin;
20200114225241	f
\.


--
-- PostgreSQL database dump complete
--
