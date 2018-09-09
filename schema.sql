--
-- PostgreSQL database dump
--

-- Dumped from database version 10.1
-- Dumped by pg_dump version 10.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: mtasa_hub
--

CREATE TABLE accounts (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    username character varying(254) NOT NULL,
    password character(60) NOT NULL,
    email character varying(254) NOT NULL,
    activated boolean DEFAULT false NOT NULL
);


ALTER TABLE accounts OWNER TO mtasa_hub;

--
-- Name: accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: mtasa_hub
--

CREATE SEQUENCE accounts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE accounts_id_seq OWNER TO mtasa_hub;

--
-- Name: accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mtasa_hub
--

ALTER SEQUENCE accounts_id_seq OWNED BY accounts.id;


--
-- Name: resources; Type: TABLE; Schema: public; Owner: mtasa_hub
--

CREATE TABLE resources (
    id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    creator integer NOT NULL,
    name character varying(254) NOT NULL,
    title text NOT NULL,
    description text NOT NULL
);


ALTER TABLE resources OWNER TO mtasa_hub;

--
-- Name: resources_id_seq; Type: SEQUENCE; Schema: public; Owner: mtasa_hub
--

CREATE SEQUENCE resources_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE resources_id_seq OWNER TO mtasa_hub;

--
-- Name: resources_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: mtasa_hub
--

ALTER SEQUENCE resources_id_seq OWNED BY resources.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: mtasa_hub
--

CREATE TABLE schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE schema_migrations OWNER TO mtasa_hub;

--
-- Name: accounts id; Type: DEFAULT; Schema: public; Owner: mtasa_hub
--

ALTER TABLE ONLY accounts ALTER COLUMN id SET DEFAULT nextval('accounts_id_seq'::regclass);


--
-- Name: resources id; Type: DEFAULT; Schema: public; Owner: mtasa_hub
--

ALTER TABLE ONLY resources ALTER COLUMN id SET DEFAULT nextval('resources_id_seq'::regclass);


--
-- Name: accounts accounts_email_key; Type: CONSTRAINT; Schema: public; Owner: mtasa_hub
--

ALTER TABLE ONLY accounts
    ADD CONSTRAINT accounts_email_key UNIQUE (email);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: mtasa_hub
--

ALTER TABLE ONLY accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: accounts accounts_username_key; Type: CONSTRAINT; Schema: public; Owner: mtasa_hub
--

ALTER TABLE ONLY accounts
    ADD CONSTRAINT accounts_username_key UNIQUE (username);


--
-- Name: resources resources_name_key; Type: CONSTRAINT; Schema: public; Owner: mtasa_hub
--

ALTER TABLE ONLY resources
    ADD CONSTRAINT resources_name_key UNIQUE (name);


--
-- Name: resources resources_pkey; Type: CONSTRAINT; Schema: public; Owner: mtasa_hub
--

ALTER TABLE ONLY resources
    ADD CONSTRAINT resources_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: mtasa_hub
--

ALTER TABLE ONLY schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: resources resources_creator_fkey; Type: FK CONSTRAINT; Schema: public; Owner: mtasa_hub
--

ALTER TABLE ONLY resources
    ADD CONSTRAINT resources_creator_fkey FOREIGN KEY (creator) REFERENCES accounts(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

