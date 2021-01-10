CREATE TABLE "domain" (
    id bigint NOT NULL,
    hostname character varying(255) NOT NULL,
    filename character varying(255),
    def_time timestamp with time zone DEFAULT now(),
    mod_time timestamp with time zone DEFAULT now()
);

CREATE SEQUENCE domain_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE domain_id_seq OWNED BY "domain".id;
ALTER TABLE ONLY "domain" ALTER COLUMN id SET DEFAULT nextval('domain_id_seq'::regclass);
ALTER TABLE ONLY "domain" ADD CONSTRAINT domain_pkey PRIMARY KEY (id);
CREATE INDEX domain_domain_index ON "domain"(hostname);

CREATE OR REPLACE FUNCTION update_mod_time_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.mod_time = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_domain_mod_time BEFORE UPDATE
    ON "domain" FOR EACH ROW EXECUTE PROCEDURE
    update_mod_time_column();
