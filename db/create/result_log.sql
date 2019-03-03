CREATE TABLE result_log (
  id SERIAL,
  ip character varying(64) NOT NULL,
  hit smallint NOT NULL,
  created_at timestamp with time zone DEFAULT now(),
  updated_at timestamp with time zone DEFAULT now(),
  PRIMARY KEY (id)
);
