create table users (
  "id" bigserial PRIMARY KEY,
  "user_id" varchar(255),
  "username" varchar(255),
  "name" varchar(255),
  created_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
)
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_id_uniq ON users(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_username_uniq ON users(username);

create table resi (
  "id" bigserial PRIMARY KEY,
  "user_id" varchar(255),
  "no_resi" varchar(100),
  "vendor" varchar(100),
  "history" TEXT DEFAULT '',
  "chat_id" varchar(100),
  "completed" boolean DEFAULT false,
  "name" varchar(100) DEFAULT '',
  created_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ  DEFAULT CURRENT_TIMESTAMP
)
CREATE UNIQUE INDEX IF NOT EXISTS idx_resi_user_id_uniq ON resi(no_resi,user_id);
