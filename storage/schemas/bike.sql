CREATE TABLE bikes (
    id serial   NOT NULL,
    public_id character varying(20) UNIQUE NOT NULL,
    latitude real NOT NULL,
    longitude real NOT NULL,
    status integer  NOT NULL DEFAULT 1,
    CONSTRAINT bikes_pkey PRIMARY KEY (id)
)
With(OIDS=FALSE);