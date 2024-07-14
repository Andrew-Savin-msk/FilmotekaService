CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR NOT NULL UNIQUE,
    encrypted_password VARCHAR NOT NULL,
    is_admin BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS films (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    description VARCHAR(1000),
    release_date DATE NOT NULL,
    assesment INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS actors (
    id BIGSERIAL PRIMARY KEY,
    gender VARCHAR(20),
    birthdate DATE,
    name VARCHAR(150)
);

CREATE TABLE IF NOT EXISTS films_actors (
    film_id BIGINT NOT NULL,
    actor_id BIGINT NOT NULL,
    PRIMARY KEY (film_id, actor_id),
    FOREIGN KEY (film_id) REFERENCES films (id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES actors (id) ON DELETE CASCADE
);

INSERT INTO users (email, encrypted_password, is_admin) VALUES ('admin@mail.ru', '$2a$04$IpAACgn5/jGuN3ZzZCUJqu727qIC3CtYQYE1iH3BcZdOTX7gvnG.O', true);

INSERT INTO films (name, description, release_date, assesment) VALUES 
('The Shawshank Redemption', 'Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.', '1994-09-22', 9),
('The Godfather', 'The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.', '1972-03-24', 9),
('The Dark Knight', 'When the menace known as the Joker emerges from his mysterious past, he wreaks havoc and chaos on the people of Gotham.', '2008-07-18', 8);

INSERT INTO actors (gender, birthdate, name) VALUES 
('Male', '1937-04-25', 'Al Pacino'),
('Male', '1955-05-06', 'Tom Hanks'),
('Female', '1975-11-04', 'Kate Winslet');

INSERT INTO films_actors (film_id, actor_id) VALUES 
(1, 2),
(2, 1),
(3, 3);
