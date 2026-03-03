package main

import (
    "fmt"
    "net/http"
    "sync"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type Pass struct {
    ID           int    `json:"id"`
    EmployeeName string `json:"employeeName"`
    Department   string `json:"department"`
    ValidUntil   string `json:"validUntil"`
    UserID       int    `json:"userId"`      // кто создал
    CreatorName  string `json:"creatorName"` // заполняется при ответе
}

var (
    usersMu     sync.RWMutex
    users       = make(map[int]User)
    passesMu    sync.RWMutex
    passes      = make(map[int]Pass)
    nextUserID  = 1
    nextPassID  = 1
)

var jwtSecret = []byte("super-secret-key-change-in-production")

type Claims struct {
    UserID int `json:"userId"`
    jwt.RegisteredClaims
}

func generateToken(userID int) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

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

func main() {
    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    r.POST("/register", register)
    r.POST("/login", login)

    auth := r.Group("/")
    auth.Use(authMiddleware())
    {
        auth.GET("/passes", getPasses)
        auth.GET("/passes/:id", getPass)
        auth.POST("/passes", createPass)
        auth.PUT("/passes/:id", updatePass)
        auth.DELETE("/passes/:id", deletePass)
    }

    r.Run(":8080")
}

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

    usersMu.Lock()
    defer usersMu.Unlock()

    for _, u := range users {
        if u.Username == req.Username {
            c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
            return
        }
    }

    hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    user := User{
        ID:       nextUserID,
        Username: req.Username,
        Password: string(hashed),
    }
    users[nextUserID] = user
    nextUserID++

    token, _ := generateToken(user.ID)
    c.JSON(http.StatusCreated, gin.H{"token": token})
}

func login(c *gin.Context) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    usersMu.RLock()
    defer usersMu.RUnlock()

    var found User
    for _, u := range users {
        if u.Username == req.Username {
            found = u
            break
        }
    }
    if found.ID == 0 {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    token, _ := generateToken(found.ID)
    c.JSON(http.StatusOK, gin.H{"token": token})
}

// getPasses возвращает все пропуски (без фильтрации по userId)
func getPasses(c *gin.Context) {
    passesMu.RLock()
    defer passesMu.RUnlock()

    list := make([]Pass, 0, len(passes))
    for _, p := range passes {
        p.CreatorName = getUserName(p.UserID)
        list = append(list, p)
    }
    c.JSON(http.StatusOK, list)
}

func getPass(c *gin.Context) {
    id := 0
    if _, err := fmt.Sscanf(c.Param("id"), "%d", &id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    passesMu.RLock()
    defer passesMu.RUnlock()

    p, ok := passes[id]
    if !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "pass not found"})
        return
    }
    p.CreatorName = getUserName(p.UserID)
    c.JSON(http.StatusOK, p)
}

func createPass(c *gin.Context) {
    userID := c.GetInt("userId")
    var p Pass
    if err := c.ShouldBindJSON(&p); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
        return
    }
    if p.EmployeeName == "" || p.Department == "" || p.ValidUntil == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "all fields are required"})
        return
    }

    passesMu.Lock()
    defer passesMu.Unlock()
    p.ID = nextPassID
    p.UserID = userID
    p.CreatorName = getUserName(userID)
    passes[nextPassID] = p
    nextPassID++
    c.JSON(http.StatusCreated, p)
}

func updatePass(c *gin.Context) {
    //userID := c.GetInt("userId")
    id := 0
    if _, err := fmt.Sscanf(c.Param("id"), "%d", &id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var updated Pass
    if err := c.ShouldBindJSON(&updated); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
        return
    }

    passesMu.Lock()
    defer passesMu.Unlock()
    p, ok := passes[id]
    if !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "pass not found"})
        return
    }
    // При обновлении сохраняем исходного создателя
    updated.ID = id
    updated.UserID = p.UserID
    updated.CreatorName = getUserName(p.UserID)
    passes[id] = updated
    c.JSON(http.StatusOK, updated)
}

func deletePass(c *gin.Context) {
    id := 0
    if _, err := fmt.Sscanf(c.Param("id"), "%d", &id); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    passesMu.Lock()
    defer passesMu.Unlock()
    if _, ok := passes[id]; !ok {
        c.JSON(http.StatusNotFound, gin.H{"error": "pass not found"})
        return
    }
    delete(passes, id)
    c.Status(http.StatusNoContent)
}

// вспомогательная функция для получения имени пользователя по ID
func getUserName(userID int) string {
    usersMu.RLock()
    defer usersMu.RUnlock()
    if user, ok := users[userID]; ok {
        return user.Username
    }
    return "Неизвестный пользователь"
}