DROP TABLE IF EXISTS users, actors;

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY, 
  email VARCHAR (150) UNIQUE NOT NULL,
  password VARCHAR (200) NOT NULL,
  is_admin boolean DEFAULT FALSE);

  CREATE TABLE actors (
  id SERIAL PRIMARY KEY,
  name VARCHAR(150) NOT NULL,
  sex VARCHAR(10) NOT NULL,
  birth_date DATE NOT NULL
);

INSERT INTO users VALUES (0, 'admin@admin.com', 'test', true);