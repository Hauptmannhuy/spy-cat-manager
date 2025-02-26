CREATE TABLE IF NOT EXISTS spies (
    id SERIAL PRIMARY KEY,               
    name TEXT UNIQUE NOT NULL,  
    breed TEXT NOT NULL,  
    salary INT,
    experience INT
);