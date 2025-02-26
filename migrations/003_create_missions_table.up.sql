CREATE TABLE IF NOT EXISTS missions (
  id SERIAL PRIMARY KEY,
  spy_id INTEGER UNIQUE REFERENCES spies(id),
  -- target_1_id INTEGER UNIQUE REFERENCES targets(id),
  -- target_2_id INTEGER UNIQUE REFERENCES targets(id),
  -- target_3_id INTEGER UNIQUE REFERENCES targets(id),
  completed BOOLEAN DEFAULT FALSE
)