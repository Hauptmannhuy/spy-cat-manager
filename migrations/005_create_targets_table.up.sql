CREATE TABLE IF NOT EXISTS targets(
  id SERIAL PRIMARY KEY,
  name TEXT,
  mission_id INT REFERENCES missions(id),
  country TEXT,
  completed boolean
)