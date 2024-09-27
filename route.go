package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// User model represents a user in the system
type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
}

// Post model represents a post by a user
type Post struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	UserID  int       `json:"user_id"`
	Created time.Time `json:"created"`
}

var users = []User{}
var posts = []Post{}

// Dummy user for authentication simulation
var dummyUser = User{
	ID:       1,
	Username: "admin",
	Email:    "admin@example.com",
	Password: "password123",
	Created:  time.Now(),
}

// Middleware for basic authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok || username != dummyUser.Username || password != dummyUser.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Function to check if a user exists
func userExists(username string) bool {
	for _, user := range users {
		if user.Username == username {
			return true
		}
	}
	return false
}

// Validate user input
func validateUserInput(user User) bool {
	return user.Username != "" && user.Email != "" && user.Password != ""
}

// Logging middleware to log requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		log.Printf("%s %s %s in %v", c.Request.Method, c.Request.URL.Path, c.ClientIP(), latency)
	}
}

// Main function that sets up the Gin server
func hew() {
	router := gin.Default()

	// Use middleware for logging and authentication
	router.Use(LoggerMiddleware())
	auth := router.Group("/", AuthMiddleware())

	// User Routes
	router.GET("/users", getUsers)
	router.POST("/users", createUser)
	router.PUT("/users/:id", updateUser)
	router.DELETE("/users/:id", deleteUser)

	// Post Routes
	auth.GET("/posts", getPosts)
	auth.POST("/posts", createPost)
	auth.PUT("/posts/:id", updatePost)
	auth.DELETE("/posts/:id", deletePost)

	// Start the server
	router.Run(":8080")
}

// Get all users
func getUsers(c *gin.Context) {
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No users found"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Create a new user
func createUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !validateUserInput(newUser) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user input"})
		return
	}
	if userExists(newUser.Username) {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}
	newUser.ID = len(users) + 1
	newUser.Created = time.Now()
	users = append(users, newUser)
	c.JSON(http.StatusCreated, newUser)
}

// Update an existing user
func updateUser(c *gin.Context) {
	id := c.Param("id")
	for i, user := range users {
		if fmt.Sprintf("%d", user.ID) == id {
			var updatedUser User
			if err := c.BindJSON(&updatedUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			users[i].Username = updatedUser.Username
			users[i].Email = updatedUser.Email
			users[i].Password = updatedUser.Password
			c.JSON(http.StatusOK, users[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

// Delete an existing user
func deleteUser(c *gin.Context) {
	id := c.Param("id")
	for i, user := range users {
		if fmt.Sprintf("%d", user.ID) == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

// Get all posts
func getPosts(c *gin.Context) {
	if len(posts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No posts found"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// Create a new post
func createPost(c *gin.Context) {
	var newPost Post
	if err := c.BindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newPost.ID = len(posts) + 1
	newPost.Created = time.Now()
	posts = append(posts, newPost)
	c.JSON(http.StatusCreated, newPost)
}

// Update an existing post
func updatePost(c *gin.Context) {
	id := c.Param("id")
	for i, post := range posts {
		if fmt.Sprintf("%d", post.ID) == id {
			var updatedPost Post
			if err := c.BindJSON(&updatedPost); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			posts[i].Title = updatedPost.Title
			posts[i].Content = updatedPost.Content
			c.JSON(http.StatusOK, posts[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
}

// Delete an existing post
func deletePost(c *gin.Context) {
	id := c.Param("id")
	for i, post := range posts {
		if fmt.Sprintf("%d", post.ID) == id {
			posts = append(posts[:i], posts[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
}

// Helper function to validate post input
func validatePostInput(post Post) bool {
	return post.Title != "" && post.Content != ""
}

// Helper function to find a post by ID
func findPostByID(id int) *Post {
	for _, post := range posts {
		if post.ID == id {
			return &post
		}
	}
	return nil
}

// Logging function for different levels
func logRequest(level string, message string) {
	switch level {
	case "INFO":
		log.Printf("[INFO] %s", message)
	case "ERROR":
		log.Printf("[ERROR] %s", message)
	default:
		log.Printf("[DEBUG] %s", message)
	}
}

// Function to simulate complex business logic
func complexBusinessLogic(data string) string {
	// Simulate heavy computation or logic
	time.Sleep(2 * time.Second)
	return fmt.Sprintf("Processed: %s", data)
}

// Error handler for JSON binding errors
func handleBindError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

// Mock function for database transaction simulation
func simulateTransaction() error {
	// Simulate DB transaction
	time.Sleep(1 * time.Second)
	// Simulate a successful transaction
	return nil
}

// Health check route
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
