-- Users table stores users
CREATE TABLE IF NOT EXISTS users (
  id            VARCHAR(255)  NOT NULL,
  first_name    VARCHAR(255),
  last_name     VARCHAR(255),
  email         VARCHAR(255),
  created_at    TIMESTAMPTZ   NOT NULL  DEFAULT now(),

  CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

-- User actions table stores every action that a user takes in the bot
CREATE TABLE IF NOT EXISTS user_actions (
  id            SERIAL,
  user_id       VARCHAR(255)  NOT NULL,
  url           TEXT,
  button        VARCHAR(255),
  quick_reply   VARCHAR(255),
  message       TEXT,
  payload       TEXT,
  created_at    TIMESTAMPTZ   NOT NULL  DEFAULT now(),
  CONSTRAINT user_actions_pkey          PRIMARY KEY (id),
  CONSTRAINT user_actions_user_id_fkey  FOREIGN KEY (user_id) REFERENCES users(id) NOT DEFERRABLE
);