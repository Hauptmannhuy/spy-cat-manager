CREATE TABLE IF NOT EXISTS missions (
  id SERIAL PRIMARY KEY,
  spy_id INTEGER REFERENCES spies(id),
  completed BOOLEAN DEFAULT FALSE
);
  
  
  CREATE UNIQUE INDEX unique_spy ON missions(spy_id) WHERE spy_id IS NOT NULL;