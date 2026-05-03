package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"math/rand"
)
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type User struct {
	Username string
	Password string
}


var (
	albums = []Album{
		{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: "2", Title: "Time Out", Artist: "Dave Brubeck", Price: 37.99},
		{ID: "3", Title: "Flying Beagle", Artist: "Himiko Kikuchi", Price: 69.99},
	}

	// Registered users
	users = map[string]string{
		"admin": "admin123",
		"user1": "pass1",
		"user2": "pass2",
	}

	// token -> username
	tokens   = map[string]string{}
	tokensMu sync.RWMutex

	albumsMu sync.RWMutex
)


const tokenChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func generateToken() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = tokenChars[rand.Intn(len(tokenChars))]
	}
	return string(b)
}

func getUserFromToken(c *gin.Context) (string, bool) {
	auth := c.GetHeader("Authorization")
	if len(auth) < 8 || auth[:7] != "Bearer " {
		return "", false
	}
	token := auth[7:]
	tokensMu.RLock()
	username, ok := tokens[token]
	tokensMu.RUnlock()
	return username, ok
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, ok := getUserFromToken(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: invalid or missing token",
			})
			return
		}
		c.Set("username", username)
		c.Next()
	}
}


func loginHandler(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Basic auth credentials required"})
		return
	}

	expectedPassword, exists := users[username]
	if !exists || expectedPassword != password {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token := generateToken()
	tokensMu.Lock()
	tokens[token] = username
	tokensMu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Hi " + username + ", welcome to the Store System",
		"token":   token,
	})
}

func logoutHandler(c *gin.Context) {
	username := c.GetString("username")

	auth := c.GetHeader("Authorization")
	token := auth[7:]

	tokensMu.Lock()
	delete(tokens, token)
	tokensMu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Bye " + username + ", your token has been revoked",
	})
}

func getAlbumsHandler(c *gin.Context) {
	albumsMu.RLock()
	defer albumsMu.RUnlock()
	c.JSON(http.StatusOK, albums)
}

func getAlbumByIDHandler(c *gin.Context) {
	id := c.Param("id")

	albumsMu.RLock()
	defer albumsMu.RUnlock()

	for _, a := range albums {
		if a.ID == id {
			c.JSON(http.StatusOK, a)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Album with ID " + id + " not found"})
}

func createAlbumHandler(c *gin.Context) {
	var newAlbum Album
	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if newAlbum.ID == "" || newAlbum.Title == "" || newAlbum.Artist == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields id, title, and artist are required"})
		return
	}

	albumsMu.Lock()
	defer albumsMu.Unlock()

	for _, a := range albums {
		if a.ID == newAlbum.ID {
			c.JSON(http.StatusConflict, gin.H{"error": "Album with ID " + newAlbum.ID + " already exists"})
			return
		}
	}

	albums = append(albums, newAlbum)
	c.JSON(http.StatusCreated, newAlbum)
}

func statusHandler(c *gin.Context) {
	username := c.GetString("username")
	c.JSON(http.StatusOK, gin.H{
		"message": "Hi " + username + ", the Vinyl Store API is Up and Running",
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	})
}

func main() {
	rand.Seed(time.Now().UnixNano())

	r := gin.Default()
	r.GET("/login", loginHandler)

	auth := r.Group("/")
	auth.Use(authMiddleware())
	{
		auth.GET("/logout", logoutHandler)
		auth.GET("/albums", getAlbumsHandler)
		auth.GET("/albums/:id", getAlbumByIDHandler)
		auth.POST("/createAlbum", createAlbumHandler)
		auth.GET("/status", statusHandler)
	}

	log.Println("Vinyl Store API starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}