CREATE TABLE IF NOT EXISTS public.patients (
    id integer NOT NULL DEFAULT nextval('patients_id_seq'::regclass),
    first_name character varying(100) NOT NULL,
    last_name character varying(100) NOT NULL,
    gender character varying(10),
    contact_number character varying(15),
    address text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    date_of_birth character varying(10),
    CONSTRAINT patients_pkey PRIMARY KEY (id)
);
