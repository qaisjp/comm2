CREATE TABLE resource_ratings
(
    account integer NOT NULL,
    resource integer NOT NULL,
    positive boolean NOT NULL,
    PRIMARY KEY (account, resource),
    FOREIGN KEY (resource)
        REFERENCES resources (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    FOREIGN KEY (account)
        REFERENCES accounts (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)