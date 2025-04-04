-- +goose Up
CREATE TABLE users (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL, 
  updated_at TIMESTAMP NOT NULL,
  name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE feeds (
  id UUID PRIMARY KEY, 
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name VARCHAR(255) NOT NULL, 
  url VARCHAR(255) NOT NULL UNIQUE, 
  user_id UUID NOT NULL, 
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE feed_follows (
  id UUID PRIMARY KEY, 
  created_at TIMESTAMP NOT NULL, 
  updated_at TIMESTAMP NOT NULL, 
  feed_id UUID NOT NULL,
  user_id UUID NOT NULL,
  UNIQUE (feed_id, user_id),
  FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE, 
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


-- +goose Down 
DROP TABLE feed_follows;
DROP TABLE feeds;
DROP TABLE users;



