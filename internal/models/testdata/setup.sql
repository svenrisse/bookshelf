CREATE TABLE IF NOT EXISTS books (
    id bigserial PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    author text NOT NULL,
    year integer NOT NULL,
    pages integer NOT NULL,
    genres text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);

ALTER TABLE books ADD CONSTRAINT books_page_check CHECK (pages >= 0);
ALTER TABLE books ADD CONSTRAINT books_year_check CHECK (year BETWEEN 1 and date_part('year', now()));
ALTER TABLE books ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 10);

CREATE INDEX IF NOT EXISTS books_title_idx ON books USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS books_genres_idx ON books USING GIN (genres);

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1
);

INSERT INTO users (name, email, password_hash, activated) VALUES('Alice Jones','alice@example.com','$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',true);
