package router

import (
	"encoding/json"
	"os"

	"errors"
	// "fmt"
	// "strconv"

	// "golang-fiber-crud/common"
	"golang-fiber-crud/models"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// Global Redis client instance (replace with your connection details)
var redisClient *redis.Client

func InitDB2() error {
	uri2 := os.Getenv("REDIS_URI")
	if uri2 == "" {
		return errors.New("you must set your 'REDIS_URI' environmental variable")

	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: uri2, // Replace with your Redis server address and port
	})
	return nil

}

func AddtitlesGroup(app *fiber.App) {
	bookGroup := app.Group("/titles")

	bookGroup.Get("/", getBooks)
	bookGroup.Get("/:id", getBook)
	bookGroup.Post("/", createBook)
	bookGroup.Put("/:id", updateBook)
	bookGroup.Delete("/:id", deleteBook)
}

func getstitles(c *fiber.Ctx) error {
	// Get all book keys from a Redis Set named "titles"
	keys, err := redisClient.SMembers(c.Context(), "titles").Result()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Create an empty slice to store titles
	titles := make([]models.Book, 0)

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

		titles = append(titles, book)
	}

	return c.Status(200).JSON(fiber.Map{"data": titles})
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
