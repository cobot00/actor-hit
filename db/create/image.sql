CREATE TABLE image (
  id SERIAL,
  name character varying(40) NOT NULL,
  type character varying(10) NOT NULL,
  path character varying(100) NOT NULL,
  voice_actor character varying(40) NOT NULL DEFAULT '',
  created_at timestamp with time zone DEFAULT now(),
  updated_at timestamp with time zone DEFAULT now(),
  PRIMARY KEY (id)
);
