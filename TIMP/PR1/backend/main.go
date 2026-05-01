package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


type User struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"` // хранится хеш
	CreatedAt time.Time `json:"-"`
}

type Pass struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	EmployeeName string    `gorm:"not null" json:"employeeName"`
	Department   string    `gorm:"not null" json:"department"`
	ValidUntil   string    `gorm:"not null" json:"validUntil"` // формат YYYY-MM-DD
	UserID       uint      `json:"userId"`                      // создатель
	User         User      `gorm:"foreignKey:UserID" json:"-"`  // для Preload
	CreatedAt    time.Time `json:"-"`
}


var db *gorm.DB
var jwtSecret = []byte("5144b6f10c66048d0859a84a60b9c3f9")


type Claims struct {
	UserID uint `json:"userId"`
	jwt.RegisteredClaims
}

// --- Middleware (аутентификация) ---
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		const prefix = "Bearer "
		if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}
		tokenStr := authHeader[len(prefix):]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("userId", claims.UserID)
		c.Next()
	}
}


// Регистрация
func register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}

	// Проверка существования пользователя
	var existing User
	if err := db.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	// Хеширование пароля
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := User{
		Username: req.Username,
		Password: string(hashed),
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"token": token})
}

// Вход
func login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Получить все пропуски (доступны всем авторизованным)
func getPasses(c *gin.Context) {
	var passes []Pass
	// Preload User, чтобы получить имя создателя
	if err := db.Preload("User").Find(&passes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch passes"})
		return
	}
	// Формируем ответ с именем создателя, защита от nil
	result := make([]gin.H, len(passes))
	for i, p := range passes {
		creatorName := "Неизвестный пользователь"
		if p.User.ID != 0 {
			creatorName = p.User.Username
		}
		result[i] = gin.H{
			"id":           p.ID,
			"employeeName": p.EmployeeName,
			"department":   p.Department,
			"validUntil":   p.ValidUntil,
			"creatorName":  creatorName,
		}
	}
	c.JSON(http.StatusOK, result)
}

// Получить один пропуск по ID
func getPass(c *gin.Context) {
	id := c.Param("id")
	var pass Pass
	if err := db.Preload("User").First(&pass, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pass not found"})
		return
	}
	creatorName := "Неизвестный пользователь"
	if pass.User.ID != 0 {
		creatorName = pass.User.Username
	}
	c.JSON(http.StatusOK, gin.H{
		"id":           pass.ID,
		"employeeName": pass.EmployeeName,
		"department":   pass.Department,
		"validUntil":   pass.ValidUntil,
		"creatorName":  creatorName,
	})
}

// Создать новый пропуск
func createPass(c *gin.Context) {
	userID := c.GetUint("userId")
	var req struct {
		EmployeeName string `json:"employeeName"`
		Department   string `json:"department"`
		ValidUntil   string `json:"validUntil"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	if req.EmployeeName == "" || req.Department == "" || req.ValidUntil == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all fields are required"})
		return
	}

	pass := Pass{
		EmployeeName: req.EmployeeName,
		Department:   req.Department,
		ValidUntil:   req.ValidUntil,
		UserID:       userID,
	}
	if err := db.Create(&pass).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create pass"})
		return
	}
	// Загружаем пользователя для ответа
	db.Preload("User").First(&pass, pass.ID)
	creatorName := "Неизвестный пользователь"
	if pass.User.ID != 0 {
		creatorName = pass.User.Username
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":           pass.ID,
		"employeeName": pass.EmployeeName,
		"department":   pass.Department,
		"validUntil":   pass.ValidUntil,
		"creatorName":  creatorName,
	})
}

// Обновить пропуск
func updatePass(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		EmployeeName string `json:"employeeName"`
		Department   string `json:"department"`
		ValidUntil   string `json:"validUntil"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	var pass Pass
	if err := db.First(&pass, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pass not found"})
		return
	}

	pass.EmployeeName = req.EmployeeName
	pass.Department = req.Department
	pass.ValidUntil = req.ValidUntil
	if err := db.Save(&pass).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update pass"})
		return
	}
	db.Preload("User").First(&pass, pass.ID)
	creatorName := "Неизвестный пользователь"
	if pass.User.ID != 0 {
		creatorName = pass.User.Username
	}
	c.JSON(http.StatusOK, gin.H{
		"id":           pass.ID,
		"employeeName": pass.EmployeeName,
		"department":   pass.Department,
		"validUntil":   pass.ValidUntil,
		"creatorName":  creatorName,
	})
}

// Удалить пропуск
func deletePass(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Pass{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pass not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

// --- Вспомогательная функция ---
func generateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// --- Точка входа ---
func main() {
	// Подключение к SQLite (Data Access Layer)
	var err error
	db, err = gorm.Open(sqlite.Open("passcontrol.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Автоматическое создание таблиц
	db.AutoMigrate(&User{}, &Pass{})

	r := gin.Default()

	// Настройка CORS (важно!)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Публичные маршруты
	r.POST("/register", register)
	r.POST("/login", login)

	// Защищённые маршруты
	auth := r.Group("/")
	auth.Use(authMiddleware())
	{
		auth.GET("/passes", getPasses)
		auth.GET("/passes/:id", getPass)
		auth.POST("/passes", createPass)
		auth.PUT("/passes/:id", updatePass)
		auth.DELETE("/passes/:id", deletePass)
	}

	r.Run(":8080") // сервер слушает порт 8080
}