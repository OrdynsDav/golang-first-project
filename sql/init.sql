DROP TABLE IF EXISTS tasks;

CREATE TABLE
    tasks (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description TEXT,
        completed BOOLEAN DEFAULT false,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

INSERT INTO
    tasks (title, description, completed)
VALUES
    ('Изучить Go', 'Пройти базовый курс', true),
    (
        'Написать REST API',
        'Посмотреть видео и написать самому',
        false
    ),
    (
        'Зарелизить приложение',
        'Развернуть приложение на сервере',
        false
    );