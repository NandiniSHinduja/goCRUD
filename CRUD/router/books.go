package router

import (
	"encoding/json"
	"os"

	// "errors"
	// "fmt"
	// "strconv"

	// "golang-fiber-crud/common"
	"golang-fiber-crud/models"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// Global Redis client instance (replace with your connection details)
var redisClient *redis.Client

func init() {
	uri2 := os.Getenv("REDIS_URI")
	redisClient = redis.NewClient(&redis.Options{
		Addr: uri2, // Replace with your Redis server address and port
	})
}

func AddBooksGroup(app *fiber.App) {
	bookGroup := app.Group("/books")

	bookGroup.Get("/", getBooks)
	bookGroup.Get("/:id", getBook)
	bookGroup.Post("/", createBook)
	bookGroup.Put("/:id", updateBook)
	bookGroup.Delete("/:id", deleteBook)
}

func getsBooks(c *fiber.Ctx) error {
	// Get all book keys from a Redis Set named "books"
	keys, err := redisClient.SMembers(c.Context(), "books").Result()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Create an empty slice to store books
	books := make([]models.Book, 0)

	// Loop through each key and get the corresponding book data using HGETALL
	for _, key := range keys {
		val, err := redisClient.HGetAll(c.Context(), key).Result()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// Convert the map[string]string to a Book struct
		book := models.Book{}
		if err := json.Unmarshal([]byte(val["title"]), &book.Title); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if err := json.Unmarshal([]byte(val["author"]), &book.Author); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if err := json.Unmarshal([]byte(val["year"]), &book.Year); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		books = append(books, book)
	}

	return c.Status(200).JSON(fiber.Map{"data": books})
}

func getsBook(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "id is required"})
	}

	// Use HGETALL to retrieve all fields for the book with the given ID
	val, err := redisClient.HGetAll(c.Context(), "book:"+id).Result()
	if err != nil {
		if err == redis.Nil {
			return c.Status(404).JSON(fiber.Map{"error": "book not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Convert the map[string]string to a Book struct
	book := models.Book{}
	if err := json.Unmarshal([]byte(val["title"]), &book.Title); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := json.Unmarshal([]byte(val["author"]), &book.Author); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := json.Unmarshal([]byte(val["year"]), &book.Year); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"data": book})
}

// type createDTO struct {
// 	Title
