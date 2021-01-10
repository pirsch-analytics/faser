CREATE TABLE "domains" (
    id bigint NOT NULL,
    domain character varying(255) NOT NULL,
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE domains_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE domains_id_seq OWNED BY "domains".id;
ALTER TABLE ONLY "domains" ALTER COLUMN id SET DEFAULT nextval('domains_id_seq'::regclass);
ALTER TABLE ONLY "domains" ADD CONSTRAINT domains_pkey PRIMARY KEY (id);
CREATE INDEX domains_domain_index ON "domains"(domain);

CREATE OR REPLACE FUNCTION update_mod_time_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.mod_time = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_domains_mod_time BEFORE UPDATE
    ON "domains" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();
