-- +goose Up
-- +goose StatementBegin
CREATE TABLE songdetail (
                            id SERIAL PRIMARY KEY,               -- Уникальный идентификатор
                            group_name VARCHAR(255),             -- Название группы
                            song_name VARCHAR(255),              -- Название песни
                            release_date DATE NOT NULL,          -- Дата выпуска песни
                            link VARCHAR(255)                    -- Ссылка на ресурс (например, на YouTube)
);
CREATE TABLE verse (
                       id SERIAL PRIMARY KEY,                 -- Уникальный идентификатор куплета
                       song_id INT NOT NULL,                  -- Внешний ключ на songdetail
                       verse_number INT NOT NULL,             -- Номер куплета
                       text TEXT NOT NULL,                    -- Текст куплета
                       FOREIGN KEY (song_id) REFERENCES songdetail(id) ON DELETE CASCADE, -- Связь с songdetail
                       UNIQUE (song_id, verse_number)         -- Уникальное ограничение для комбинации song_id и verse_number
);
-- +goose StatementEnd
-- +goose Down
DROP TABLE songdetail;
DROP TABLE verse;