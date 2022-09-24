create table users (
  "id" varchar(30) PRIMARY KEY,
  "user_id" varchar(255),
  "username" varchar(255),
  "name" varchar(255),
  created_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_id_uniq ON users(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_username_uniq ON users(username);

create table resi (
  "id" varchar(30) PRIMARY KEY,
  "user_id" varchar(255),
  "no_resi" varchar(100),
  "vendor" varchar(100),
  "history" TEXT DEFAULT '',
  "chat_id" varchar(100),
  "completed" boolean DEFAULT false,
  "name" varchar(100) DEFAULT '',
  "is_deleted" boolean DEFAULT false,
  created_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_resi_user_id_uniq ON resi(no_resi,user_id);

-- CREATE TABLE public.resi (
-- 	id varchar(30) NOT NULL,
-- 	user_id varchar(255) NULL,
-- 	no_resi varchar(100) NULL,
-- 	vendor varchar(100) NULL,
-- 	chat_id varchar(100) NULL,
-- 	created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
-- 	updated_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
-- 	history text NOT NULL DEFAULT ''::text,
-- 	completed bool NOT NULL DEFAULT false,
-- 	name varchar(255) NULL DEFAULT ''::character varying,
-- 	CONSTRAINT resi_pkey PRIMARY KEY (id)
-- );
-- CREATE UNIQUE INDEX idx_resi_user_id_uniq ON public.resi USING btree (user_id, no_resi);

-- CREATE TABLE public.users (
-- 	id int8 NOT NULL,
-- 	user_id varchar(255) NULL,
-- 	username varchar(255) NULL,
-- 	"name" varchar(255) NULL,
-- 	created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
-- 	updated_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
-- 	CONSTRAINT user_pkey PRIMARY KEY (id)
-- );
-- CREATE UNIQUE INDEX idx_user_id_uniq ON public.users USING btree (user_id);
-- CREATE UNIQUE INDEX idx_username_uniq ON public.users USING btree (username);