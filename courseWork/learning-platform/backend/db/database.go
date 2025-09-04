package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Инициализация БД
func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./platform.db")
	if err != nil {
		log.Fatal(err)
	}

	// Создание таблиц
	createTables := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS tasks (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            description TEXT NOT NULL,
            difficulty TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );

        CREATE TABLE IF NOT EXISTS submissions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER,
            task_id INTEGER,
            code TEXT NOT NULL,
            status TEXT NOT NULL,
            output TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users (id),
            FOREIGN KEY (task_id) REFERENCES tasks (id)
        );

        -- Добавляем тестовые задания
        INSERT OR IGNORE INTO tasks (title, description, difficulty) VALUES 
        ('Hello World', 'Напишите программу, которая выводит "Hello, World!"', 'easy'),
        ('Сумма двух чисел', 'Напишите программу, которая складывает два числа', 'easy'),
        ('Факториал', 'Напишите программу, которая вычисляет факториал числа', 'medium');
    `

	if _, err := db.Exec(createTables); err != nil {
		log.Fatal(err)
	}

	// Добавляем столбец output если его нет
	_, err = db.Exec(`ALTER TABLE submissions ADD COLUMN output TEXT`)
	if err != nil {
		// Игнорируем ошибку если столбец уже существует
		log.Println("Note: output column may already exist:", err)
	}

	return db
}
