CREATE TABLE authors
(
    id   BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name text   NOT NULL,
    bio  text
);

CREATE TABLE books (
                       book_id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
                       author_id integer NOT NULL,
                       isbn varchar(255) NOT NULL DEFAULT '' UNIQUE,
                       book_type ENUM('FICTION', 'NONFICTION') NOT NULL DEFAULT 'FICTION',
                       title text NOT NULL,
                       yr integer NOT NULL DEFAULT 2000,
                       available datetime NOT NULL DEFAULT NOW(),
                       tags text NOT NULL
    -- CONSTRAINT FOREIGN KEY (author_id) REFERENCES authors(author_id)
) ENGINE=InnoDB;

CREATE INDEX books_title_idx ON books(title(255), yr);

/*
CREATE FUNCTION say_hello(s text) RETURNS text
  DETERMINISTIC
  RETURN CONCAT('hello ', s);
*/

