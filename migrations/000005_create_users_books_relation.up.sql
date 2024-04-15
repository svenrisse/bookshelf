CREATE TABLE IF NOT EXISTS usersBooksRelation (
  id bigserial PRIMARY KEY,
  bookId bigserial,
  userId bigserial, 
  reviewId bigserial,
  read BOOLEAN NOT NULL,
  rating FLOAT(2),
  reviewBody text,
  added_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  date_read timestamp(0) with time zone,
  reviewed_at timestamp(0) with time zone,
  version integer NOT NULL DEFAULT 1,

  FOREIGN KEY (bookId) REFERENCES books(id),
  FOREIGN KEY (userId) REFERENCES users(id),
  UNIQUE (bookId, userId)
);
