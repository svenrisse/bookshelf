CREATE TABLE IF NOT EXISTS reviews (
  id bigserial PRIMARY KEY,
  bookId bigserial,
  userId bigserial,
  rating FLOAT(2),
  body text,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  version integer NOT NULL DEFAULT 1,

  FOREIGN KEY (bookId) REFERENCES books(id),
  FOREIGN KEY (userId) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS usersBooksRelation (
  bookId bigserial,
  usersId bigserial, 
  reviewId bigserial,
  read BOOLEAN NOT NULL,
  added_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  date_read timestamp(0) with time zone,

  FOREIGN KEY (bookId) REFERENCES books(id),
  FOREIGN KEY (usersId) REFERENCES users(id),
  FOREIGN KEY (reviewId) REFERENCES reviews(id),
  UNIQUE (bookId, usersId)
);
