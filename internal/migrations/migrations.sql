DROP TABLE IF EXISTS users, films, actors, castfilms;
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY, 
  email VARCHAR (150) UNIQUE NOT NULL,
  password VARCHAR (200) NOT NULL,
  is_admin boolean DEFAULT FALSE);


CREATE TABLE films (
  id SERIAL PRIMARY KEY,
  name VARCHAR(150) NOT NULL,
  description VARCHAR(1000) NOT NULL,
  release_date DATE NOT NULL,
  rating INT NOT NULL
);

CREATE TABLE actors (
  id SERIAL PRIMARY KEY,
  name VARCHAR(150) NOT NULL,
  sex VARCHAR(10) NOT NULL,
  birth_date DATE NOT NULL
);

CREATE TABLE castfilms (
  actor_id INT NOT NULL REFERENCES actors ON DELETE CASCADE,
  film_id INT NOT NULL REFERENCES films ON DELETE CASCADE,
  UNIQUE(actor_id, film_id)
);

INSERT INTO users VALUES (0, 'admin@admin.com', 'test', true); 