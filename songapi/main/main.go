package main

import (
	"database/sql"
	"fmt"
	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"

	"github.com/joho/godotenv"
	"github.com/pressly/goose"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// @title Song API
// @version 1.0
// @description API for managing songs with verses
// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

var db *sql.DB

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	// Подключаемся к базе данных
	ping := dbConnect()
	log.Printf("[INFO] Database connection status: %s", ping)

	// Создаем маршрутизатор
	router := gin.Default()

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Роуты
	router.DELETE("/deleteSong", deleteSong)
	router.PATCH("/updateSong", updateSong)
	router.POST("/addSong", addSong)
	router.GET("/getSong", getSong)
	router.GET("/gettextWithPagination", gettextWithPagination)

	// Запускаем сервер
	err := router.Run("localhost:8081")
	if err != nil {
		log.Fatalf("[ERROR] Failed to start server: %v", err)
	}
}

func dbConnect() error {
	var err error

	// Получаем параметры для подключения к БД из .env
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return err
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"))

	// Открываем соединение с базой данных
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Выполняем миграции с помощью goose
	migrationDir := "./db/migrations"
	if err := goose.Up(db, migrationDir); err != nil {
		log.Fatalf("[ERROR] Could not run migrations: %v", err)
		return err
	}

	// Проверяем соединение с базой
	if err = db.Ping(); err != nil {
		log.Printf("[ERROR] Database ping failed: %v", err)
		return err
	}

	log.Println("[INFO] Database connection established successfully")
	return nil
}

// @Summary Get song with pagination
// @Description Returns songs with pagination based on text search
// @Param text query string false "Text to search for"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {array} openapi.SongDetail
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /gettextWithPagination [get]
func gettextWithPagination(c *gin.Context) {

	page := c.DefaultQuery("page", os.Getenv("PAGE_DEF"))
	limit := c.DefaultQuery("limit", os.Getenv("LIMIT_DEF"))

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}

	id := c.DefaultQuery("id", "")
	releaseDate := c.DefaultQuery("releaseDate", "")
	link := c.DefaultQuery("link", "")
	text := c.DefaultQuery("text", "")

	query := `SELECT sd.release_date, v.text, sd.link 
              FROM songdetail sd 
              JOIN verse v ON sd.id = v.song_id 
              WHERE 1=1`

	var args []interface{}
	argIndex := 1

	if id != "" {
		query += " AND v.song_id = $" + strconv.Itoa(argIndex)
		args = append(args, id)
		argIndex++
	}

	if releaseDate != "" {
		query += " AND sd.release_date ILIKE $" + strconv.Itoa(argIndex)
		args = append(args, "%"+releaseDate+"%")
		argIndex++
	}

	if text != "" {
		query += " AND v.text ILIKE $" + strconv.Itoa(argIndex)
		args = append(args, "%"+text+"%")
		argIndex++
	}

	if link != "" {
		query += " AND sd.link ILIKE $" + strconv.Itoa(argIndex)
		args = append(args, "%"+link+"%")
		argIndex++
	}

	query += " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, limitInt, (pageInt-1)*limitInt)

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var songs []openapi.SongDetail

	for rows.Next() {
		var song openapi.SongDetail
		if err := rows.Scan(&song.ReleaseDate, &song.Text, &song.Link); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		songs = append(songs, song)
	}

	c.IndentedJSON(http.StatusOK, songs)
}

// @Summary Delete song
// @Description Deletes a song based on query parameters
// @Param id query string false "Song ID"
// @Param releaseDate query string false "Release Date"
// @Param text query string false "Song text"
// @Param link query string false "Song link"
// @Success 200 {object} map[string]string "message: Song deleted successfully"
// @Failure 500 {object} map[string]string "error"
// @Router /deleteSong [delete]
func deleteSong(c *gin.Context) {
	id := c.Query("id")
	releaseDate := c.Query("releaseDate")
	text := c.Query("text")
	link := c.Query("link")

	query := `DELETE FROM songdetail WHERE `
	var args []interface{}
	conditions := []string{}

	if id != "" {
		conditions = append(conditions, " id = $1")
		args = append(args, id)
	}
	if releaseDate != "" {
		conditions = append(conditions, " release_date ILIKE $"+strconv.Itoa(len(args)+1))
		args = append(args, "%"+releaseDate+"%")
	}
	if text != "" {
		conditions = append(conditions, " text ILIKE $"+strconv.Itoa(len(args)+1))
		args = append(args, "%"+text+"%")
	}
	if link != "" {
		conditions = append(conditions, " link ILIKE $"+strconv.Itoa(len(args)+1))
		args = append(args, "%"+link+"%")
	}

	if len(conditions) > 0 {
		query += strings.Join(conditions, " AND ")
	}

	_, err := db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}

// @Summary Update song
// @Description Updates song details and verses
// @Param id query string true "Song ID"
// @Param updateReleaseDate query string false "New Release Date"
// @Param updateLink query string false "New Song Link"
// @Param verses body []map[string]interface{} true "Verses to update"
// @Success 200 {object} map[string]string "message: Song updated successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "error"
// @Router /updateSong [patch]
func updateSong(c *gin.Context) {
	id := c.Query("id")
	updateReleaseDate := c.Query("updateReleaseDate")
	updateLink := c.Query("updateLink")

	var verses []struct {
		VerseNumber int    `json:"verseNumber"`
		Text        string `json:"text"`
	}
	if err := c.ShouldBindJSON(&verses); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var args []interface{}
	var conditions []string

	if updateReleaseDate != "" {
		conditions = append(conditions, " release_date = $"+strconv.Itoa(len(args)+1))
		args = append(args, updateReleaseDate)
	}
	if updateLink != "" {
		conditions = append(conditions, " link = $"+strconv.Itoa(len(args)+1))
		args = append(args, updateLink)
	}
	if len(conditions) > 0 {
		query := "UPDATE songdetail SET " + strings.Join(conditions, ", ") + " WHERE id = $" + strconv.Itoa(len(args)+1)
		args = append(args, id)

		_, err := db.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	for _, verse := range verses {
		verseQuery := `UPDATE verse SET text = $1 WHERE song_id = $2 AND verse_number = $3`
		_, err := db.Exec(verseQuery, verse.Text, id, verse.VerseNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

// @Summary Add new song
// @Description Adds a new song with verses
// @Param songData body map[string]interface{} true "Song data"
// @Success 200 {object} map[string]string "message: Song added successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "error"
// @Router /addSong [post]
func addSong(c *gin.Context) {
	var songData struct {
		ReleaseDate string `json:"releaseDate"`
		Link        string `json:"link"`
		Verses      []struct {
			VerseNumber int    `json:"verseNumber"`
			Text        string `json:"text"`
		} `json:"verses"`
	}

	if err := c.ShouldBindJSON(&songData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := `INSERT INTO songdetail (release_date, link) VALUES ($1, $2) RETURNING id`
	var songID int
	err := db.QueryRow(query, songData.ReleaseDate, songData.Link).Scan(&songID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, verse := range songData.Verses {
		verseQuery := `INSERT INTO verse (song_id, verse_number, text) VALUES ($1, $2, $3)`
		_, err := db.Exec(verseQuery, songID, verse.VerseNumber, verse.Text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song added successfully"})
}

// @Summary Get song by ID
// @Description Returns song details and verses based on song ID
// @Param id query string true "Song ID"
// @Success 200 {object} openapi.SongDetail
// @Failure 404 {object} map[string]string "Song not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /getSong [get]
func getSong(c *gin.Context) {
	songID := c.DefaultQuery("id", "")

	if songID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Song ID is required"})
		return
	}

	var song openapi.SongDetail
	query := `SELECT sd.release_date, v.text, sd.link 
              FROM songdetail sd
              JOIN verse v ON sd.id = v.song_id 
              WHERE sd.id = $1`

	row := db.QueryRow(query, songID)
	err := row.Scan(&song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, song)
}
