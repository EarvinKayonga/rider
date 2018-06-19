
CREATE TABLE trips (
    id serial   NOT NULL,
    started_at date NOT NULL DEFAULT CURRENT_DATE,
    ended_at date,
    public_id character varying(20) UNIQUE NOT NULL,
    bike_id character varying(20)   NOT NULL,
    status integer  NOT NULL,
    CONSTRAINT trips_pkey PRIMARY KEY (id)
)
With(OIDS=FALSE);


CREATE TABLE locations (
    id serial   NOT NULL,
    latitude real NOT NULL,
    longitude real NOT NULL,
    trip_id character varying(20) NOT NULL,
    CONSTRAINT locations_pkey PRIMARY KEY (id) ,
    FOREIGN KEY (trip_id) REFERENCES trips(public_id)
)
With(OIDS=FALSE);