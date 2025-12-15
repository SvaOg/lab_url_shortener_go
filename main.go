package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/jackc/pgx/v5/stdlib" // Драйвер для Postgres
)

// Структура для JSON запроса
type CreateRequest struct {
	URL string `json:"url"`
}

// Структура для JSON ответа
type UrlInfo struct {
	URL  string `json:"url"`
	Slug string `json:"slug"`
}

var db *sql.DB

func main() {
	// 1. Подключение к БД
	var err error
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		// Дефолт для локального запуска без докера
		dbURL = "postgres://postgres:postgres@localhost:6432/postgres"
	}

	db, err = sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	// 2. Создание таблицы (аналог lifespan из Python)
	// В Go миграции обычно делают отдельно, но для простоты создадим тут
	query := `
	CREATE TABLE IF NOT EXISTS short_urls (
		slug VARCHAR(10) PRIMARY KEY,
		long_url TEXT NOT NULL
	);`
	if _, err := db.Exec(query); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// 3. Настройка роутера (Chi)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Настройка CORS (чтобы фронтенд мог стучаться)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
	}))

	// 4. Ручки (Endpoints)
	r.Post("/", createShortURL)
	r.Get("/{slug}", redirectURL)

	log.Println("Server starting on :8000")
	http.ListenAndServe(":8000", r)
}

// Handler: Создание ссылки
func createShortURL(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	var slug string
	// Попытки генерации (Retry logic)
	for i := 0; i < 5; i++ {
		slug = generateSlug(6)

		// Пытаемся вставить. Если slug занят, Postgres вернет ошибку.
		_, err := db.Exec("INSERT INTO short_urls (slug, long_url) VALUES ($1, $2)", slug, req.URL)
		if err == nil {
			// Успех
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(UrlInfo{URL: req.URL, Slug: slug})
			return
		}
		// Если ошибка, цикл пойдет дальше
	}

	http.Error(w, "Could not generate unique slug", http.StatusInternalServerError)
}

// Handler: Редирект
func redirectURL(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var longURL string
	err := db.QueryRow("SELECT long_url FROM short_urls WHERE slug = $1", slug).Scan(&longURL)

	if err == sql.ErrNoRows {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 307 Temporary Redirect
	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}

// Генератор случайных строк (безопасный)
func generateSlug(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	ret := make([]byte, n)
	for i := range n {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}
	return string(ret)
}
