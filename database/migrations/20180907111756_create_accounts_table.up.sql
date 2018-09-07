CREATE TABLE accounts
(
    id serial NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    username character varying(254) NOT NULL,
    password character (60) NOT NULL,
    email character varying(254) NOT NULL,
    activated boolean DEFAULT false NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (username),
    UNIQUE (email)
)
