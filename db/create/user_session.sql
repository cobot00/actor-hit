CREATE TABLE user_session(
  id SERIAL,
  turn smallint NOT NULL DEFAULT 0,
  hit smallint NOT NULL DEFAULT 0,
  answer character varying(40) NOT NULL DEFAULT '',
  answer_character_names character varying(250) NOT NULL DEFAULT '',
  result_summary character varying(250) NOT NULL DEFAULT '',
  expired_at timestamp with time zone NOT NULL,
  created_at timestamp with time zone DEFAULT now(),
  updated_at timestamp with time zone DEFAULT now(),
  PRIMARY KEY (id)
);
