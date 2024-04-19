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
    avatar text, 
    provider text NOT NULL
);

CREATE TABLE IF NOT EXISTS usersBooksRelation (
  id bigserial PRIMARY KEY,
  bookId bigserial,
  userId bigserial, 
  reviewId bigserial,
  read BOOLEAN NOT NULL,
  rating FLOAT(2),
  reviewBody text,
  added_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  read_at timestamp(0) with time zone,
  reviewed_at timestamp(0) with time zone,
  version int NOT NULL DEFAULT 1,

  FOREIGN KEY (bookId) REFERENCES books(id),
  FOREIGN KEY (userId) REFERENCES users(id),
  UNIQUE (bookId, userId)
);

INSERT INTO users (id, name, avatar, provider) VALUES (1, 'Alice Jones', 'avat', 'discord');
INSERT INTO books (id, title, author, year, pages, genres) VALUES (1, 'The Hobbit', 'JRR Tolkien', 1890, 320, ARRAY ['Fantasy', 'Childrens Literature']);
INSERT INTO books (id, title, author, year, pages, genres) VALUES (2, 'A Game Of Thrones', 'GRRM Martin', 1990, 700, ARRAY ['Fantasy', 'Epic']);

INSERT INTO usersBooksRelation (id, bookId, userId, read, rating, reviewBody, read_at, reviewed_at, version)
VALUES (14, 2, 1, true, 4.5, 'Very good book yes!', '2024-04-10 14:30:00', '2024-04-11 15:00:00', 1);
