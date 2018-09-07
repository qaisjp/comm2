CREATE TABLE resources
(
    id serial NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    creator integer NOT NULL,
    name character varying(254) NOT NULL,
    title text NOT NULL,
    description text NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (name),
    FOREIGN KEY (creator)
        REFERENCES public.accounts (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)
