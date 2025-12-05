-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE IF NOT EXISTS bioskop
(
    "ID" integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    "Nama" character varying(50) NOT NULL,
    "Lokasi" character varying(100) NOT NULL,
    "Rating" real,
    CONSTRAINT bioskop_pkey PRIMARY KEY ("ID")
)

-- +migrate StatementEnd