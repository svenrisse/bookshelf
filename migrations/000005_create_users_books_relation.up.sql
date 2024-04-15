CREATE TABLE IF NOT EXISTS usersBooksRelation (
  bookId bigserial,
  usersId bigserial, 
  read BOOLEAN NOT NULL,
  added_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  date_read timestamp(0) with time zone,

  FOREIGN KEY (bookId) REFERENCES books(id),
  FOREIGN KEY (usersId) REFERENCES users(id),
  UNIQUE (bookId, usersId)
);
