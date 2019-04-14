-- Drop table

-- DROP TABLE public.migrations

CREATE  TABLE public.migrations (
	id serial NOT NULL,
	filename varchar(1024) NOT NULL,
	applied bool NOT NULL DEFAULT false,
	CONSTRAINT migrations_pk PRIMARY KEY (id),
	CONSTRAINT migrations_un UNIQUE (filename)
);