CREATE TABLE IF NOT EXISTS public.appointments (
    id integer NOT NULL DEFAULT nextval('appointments_id_seq'::regclass),
    patient_id integer,
    doctor_id integer,
    appointment_date timestamp with time zone NOT NULL,
    notes text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    recep_id integer NOT NULL,
    CONSTRAINT appointments_pkey PRIMARY KEY (id),
    CONSTRAINT appointments_doctor_id_fkey FOREIGN KEY (doctor_id)
        REFERENCES public.users (id) ON UPDATE NO ACTION ON DELETE SET NULL,
    CONSTRAINT appointments_patient_id_fkey FOREIGN KEY (patient_id)
        REFERENCES public.patients (id) ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT appointments_recep_id_fkey FOREIGN KEY (recep_id)
        REFERENCES public.users (id) ON UPDATE NO ACTION ON DELETE SET NULL
);
