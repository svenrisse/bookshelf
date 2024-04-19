CREATE TABLE IF NOT EXISTS usersBooksRelation (
  id bigserial PRIMARY KEY,
  bookId bigserial,
  userId bigserial, 
  read BOOLEAN NOT NULL,
  rating REAL,
  reviewBody text,
  added_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  read_at timestamp(0) with time zone,
  reviewed_at timestamp(0) with time zone,
  version integer NOT NULL DEFAULT 1,

  FOREIGN KEY (bookId) REFERENCES books(id),
  FOREIGN KEY (userId) REFERENCES users(id),
  UNIQUE (bookId, userId)
);
