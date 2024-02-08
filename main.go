package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// Task model
type Task struct {
	gorm.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

func main() {
	// Initialize SQLite database using gorm
	initDB()

	// Set the Gin mode to release
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	router := gin.New()

	// Define CRUD endpoints
	router.POST("/tasks", createTask)
	router.GET("/tasks/:id", getTask)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteTask)
	router.GET("/tasks", listTasks)

	// Run the server on port 8080
	router.Run(":8080")
}

// Initialize SQLite database using gorm
func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("./tasks.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto Migrate the schema
	db.AutoMigrate(&Task{})
}

// Create a new task
func createTask(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&task)

	c.JSON(http.StatusCreated, task)
}

// Retrieve a task
func getTask(c *gin.Context) {
	var task Task
	id := c.Param("id")

	if err := db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// Update a task
func updateTask(c *gin.Context) {
	var task Task
	id := c.Param("id")

	if err := db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Updates(&task)

	c.JSON(http.StatusOK, task)
}

// Delete a task
func deleteTask(c *gin.Context) {
	var task Task
	id := c.Param("id")

	if err := db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	db.Delete(&task)

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// List all tasks
func listTasks(c *gin.Context) {
	var tasks []Task
	db.Find(&tasks)

	c.JSON(http.StatusOK, tasks)
}
