CREATE TABLE IF NOT EXISTS movies (
    id bigserial PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    author text NOT NULL,
    year integer NOT NULL,
    pages integer NOT NULL,
    genres text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);
