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
('The Dark Knight', 'When the menace known as the Joker emerges from his mysterious past, he wreaks havoc and chaos on the people of Gotham.', '2008-07-18', 8),
('Зеленая миля', 'Тюремный охранник обнаруживает, что осужденный на смертную казнь обладает сверхъестественными способностями.', '1999-12-06', 9),
('Побег из Шоушенка', 'История заключенного, который, несмотря на все препятствия, находит способ сбежать из тюрьмы Шоушенк.', '1994-09-23', 9),
('Список Шиндлера', 'Реальная история оскароносного бизнесмена, спасшего более тысячи еврейских беженцев во время Холокоста.', '1993-12-15', 9),
('Властелин колец: Возвращение короля', 'Заключительная часть трилогии о противостоянии Сэма и Фродо Саурону.', '2003-12-17', 9),
('Титаник', 'Романтическая история на фоне трагедии гибели легендарного лайнера Титаник.', '1997-12-19', 8),
('Форрест Гамп', 'История необычного мужчины, который, несмотря на свои ограничения, оказывается вовлеченным в самые значимые события второй половины 20 века.', '1994-07-06', 8),
('Властелин колец: Братство кольца', 'Первая часть трилогии о борьбе против Саурона и поиске силы для уничтожения кольца.', '2001-12-19', 9),
('Крёстный отец', 'Эпическая история о семье мафии Корлеоне и её лидере Вито Корлеоне.', '1972-03-24', 9),
('Начало', 'Группа специалистов проникает в подсознание людей, чтобы внедрить идею.', '2010-07-16', 8),
('Матрица', 'Компьютерный хакер узнает правду о своей реальности и своей роли в войне против контролирующих её машин.', '1999-03-31', 8);

INSERT INTO actors (gender, birthdate, name) VALUES 
('Male', '1937-04-25', 'Al Pacino'),
('Male', '1955-05-06', 'Tom Hanks'),
('Female', '1975-11-04', 'Kate Winslet'),
('Male', '1963-03-18', 'Брюс Уиллис'),
('Female', '1983-12-12', 'Лорен Джерман'),
('Male', '1962-09-25', 'Майкл Мэдсен'),
('Male', '1961-06-25', 'Рики Джервэйс'),
('Male', '1962-09-25', 'Кристофер Рив'),
('Female', '1981-04-09', 'Джена Ушковиц'),
('Female', '1974-09-27', 'Кэри-Энн Мосс'),
('Female', '1978-05-10', 'Малин Акерман'),
('Male', '1964-09-02', 'Киану Ривз'),
('Male', '1967-11-02', 'Дэвид Швиммер'),
('Female', '1969-06-15', 'Кортни Кокс'),
('Male', '1963-07-24', 'Карлос Санта'),
('Male', '1954-09-07', 'Майкл Эмерсон'),
('Male', '1955-09-21', 'Билл Мюррей'),
('Female', '1975-07-26', 'Элизабет Бэнкс'),
('Female', '1982-11-28', 'Карен Гиллан'),
('Female', '1976-04-28', 'Элизабет Рем'),
('Female', '1971-07-24', 'Эмбер Тэмблин'),
('Male', '1956-12-03', 'Джулианн Мур');

INSERT INTO films_actors (film_id, actor_id) VALUES 
(1, 2),
(2, 1),
(3, 3),
(4, 1),
(4, 5),
(4, 6),
(5, 7),
(6, 3),
(7, 8),
(7, 1),
(7, 4),
(8, 9),
(9, 1),
(10, 7),
(11, 5),
(11, 1),
(11, 2);
