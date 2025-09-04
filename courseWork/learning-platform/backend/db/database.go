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
		log.Fatal("Ошибка открытия базы данных:", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
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
            id INTEGER PRIMARY KEY,
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
            FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
            FOREIGN KEY (task_id) REFERENCES tasks (id) ON DELETE CASCADE
        );

        -- Создаем индекс для предотвращения дублирования
        CREATE UNIQUE INDEX IF NOT EXISTS idx_tasks_id ON tasks(id);
        CREATE INDEX IF NOT EXISTS idx_submissions_user_id ON submissions(user_id);
        CREATE INDEX IF NOT EXISTS idx_submissions_created_at ON submissions(created_at);
    `

	if _, err := db.Exec(createTables); err != nil {
		log.Fatal("Ошибка создания таблиц:", err)
	}

	// Проверяем и добавляем отсутствующие столбцы
	checkAndAddColumn(db, "submissions", "output", "TEXT")
	checkAndAddColumn(db, "submissions", "created_at", "DATETIME DEFAULT CURRENT_TIMESTAMP")
	checkAndAddColumn(db, "submissions", "test_results", "TEXT")

	log.Println("База данных успешно инициализирована")
	return db
}

// checkAndAddColumn проверяет существование столбца и добавляет его если нужно
func checkAndAddColumn(db *sql.DB, tableName, columnName, columnType string) {
	var exists int
	err := db.QueryRow(`
        SELECT COUNT(*) FROM pragma_table_info(?) WHERE name = ?
    `, tableName, columnName).Scan(&exists)

	if err != nil {
		log.Printf("Ошибка проверки столбца %s: %v", columnName, err)
		return
	}

	if exists == 0 {
		_, err := db.Exec(`ALTER TABLE ` + tableName + ` ADD COLUMN ` + columnName + ` ` + columnType)
		if err != nil {
			log.Printf("Ошибка добавления столбца %s: %v", columnName, err)
		} else {
			log.Printf("Добавлен столбец %s в таблицу %s", columnName, tableName)
		}
	}
}

// Task представляет структуру задания для синхронизации
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
}

// SyncTasks синхронизирует задания из переданного списка с БД
func SyncTasks(db *sql.DB, tasks []Task) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, task := range tasks {
		// Используем INSERT OR IGNORE чтобы избежать дублирования
		_, err := tx.Exec(`
            INSERT OR IGNORE INTO tasks (id, title, description, difficulty) 
            VALUES (?, ?, ?, ?)
        `, task.ID, task.Title, task.Description, task.Difficulty)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
