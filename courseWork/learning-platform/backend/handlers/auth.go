package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"platform/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	database := c.MustGet("db").(*sql.DB)
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Проверяем, существует ли пользователь
	var existingUserID int
	err := database.QueryRow("SELECT id FROM users WHERE username = ?", user.Username).Scan(&existingUserID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь уже существует"})
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хеширования пароля"})
		return
	}

	// Сохранение в БД
	result, err := database.Exec("INSERT INTO users (username, password) VALUES (?, ?)",
		user.Username, string(hashedPassword))
	if err != nil {
		log.Printf("Ошибка при создании пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания пользователя"})
		return
	}

	userID, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{
		"message": "Пользователь создан",
		"id":      userID,
	})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	database := c.MustGet("db").(*sql.DB)
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	var user models.User
	err := database.QueryRow("SELECT id, username, password FROM users WHERE username = ?", req.Username).Scan(
		&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		}
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверные учетные данные"})
		return
	}

	// Генерируем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    tokenString,
		"user_id":  user.ID,
		"username": user.Username,
	})
}
