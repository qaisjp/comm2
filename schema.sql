--
-- PostgreSQL database dump
--

-- Dumped from database version 11.3
-- Dumped by pg_dump version 11.3

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

SET default_with_oids = false;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: mta
--

CREATE TABLE public.accounts (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    username character varying(254) NOT NULL,
    password character(60) NOT NULL,
    email character varying(254) NOT NULL,
    activated boolean DEFAULT false NOT NULL
);


ALTER TABLE public.accounts OWNER TO mta;

--
-- Name: accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: mta
--

CREATE SEQUENCE public.accounts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.accounts_id_seq OWNER TO mta;

--
-- Name: accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mta
--

ALTER SEQUENCE public.accounts_id_seq OWNED BY public.accounts.id;


--
-- Name: resource_ratings; Type: TABLE; Schema: public; Owner: mta
--

CREATE TABLE public.resource_ratings (
    account integer NOT NULL,
    resource integer NOT NULL,
    positive boolean NOT NULL
);


ALTER TABLE public.resource_ratings OWNER TO mta;

--
-- Name: resources; Type: TABLE; Schema: public; Owner: mta
--

CREATE TABLE public.resources (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    creator integer NOT NULL,
    name character varying(254) NOT NULL,
    title text NOT NULL,
    description text NOT NULL
);


ALTER TABLE public.resources OWNER TO mta;

--
-- Name: resources_id_seq; Type: SEQUENCE; Schema: public; Owner: mta
--

CREATE SEQUENCE public.resources_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.resources_id_seq OWNER TO mta;

--
-- Name: resources_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mta
--

ALTER SEQUENCE public.resources_id_seq OWNED BY public.resources.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: mta
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO mta;

--
-- Name: accounts id; Type: DEFAULT; Schema: public; Owner: mta
--

ALTER TABLE ONLY public.accounts ALTER COLUMN id SET DEFAULT nextval('public.accounts_id_seq'::regclass);


--
-- Name: resources id; Type: DEFAULT; Schema: public; Owner: mta
--

ALTER TABLE ONLY public.resources ALTER COLUMN id SET DEFAULT nextval('public.resources_id_seq'::regclass);


--
-- Name: accounts accounts_email_key; Type: CONSTRAINT; Schema: public; Owner: mta
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_email_key UNIQUE (email);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: mta
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: accounts accounts_username_key; Type: CONSTRAINT; Schema: public; Owner: mta
--

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_username_key UNIQUE (username);


--
-- Name: resource_ratings resource_ratings_pkey; Type: CONSTRAINT; Schema: public; Owner: mta
--

ALTER TABLE ONLY public.resource_ratings
    ADD CONSTRAINT resource_ratings_pkey PRIMARY KEY (account, resource);


--
-- Name: resources resources_name_key; Type: CONSTRAINT; Schema: public; Owner: mta
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_name_key UNIQUE (name);


--
-- Name: resources resources_pkey; Type: CONSTRAINT; Schema: public; Owner: mta
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: mta
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: resource_ratings resource_ratings_account_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mta
--

ALTER TABLE ONLY public.resource_ratings
    ADD CONSTRAINT resource_ratings_account_fkey FOREIGN KEY (account) REFERENCES public.accounts(id);


--
-- Name: resource_ratings resource_ratings_resource_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mta
--

ALTER TABLE ONLY public.resource_ratings
    ADD CONSTRAINT resource_ratings_resource_fkey FOREIGN KEY (resource) REFERENCES public.resources(id);


--
-- Name: resources resources_creator_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mta
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_creator_fkey FOREIGN KEY (creator) REFERENCES public.accounts(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 11.3
-- Dumped by pg_dump version 11.3

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
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: mta
--

COPY public.schema_migrations (version, dirty) FROM stdin;
20180909013807	f
\.


--
-- PostgreSQL database dump complete
--

