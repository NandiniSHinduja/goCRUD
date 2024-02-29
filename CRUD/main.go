package main

import (
	"encoding/json"
	"errors"
	"github.com/bmdavis419/fiber-mongo-example/common"
	"github.com/bmdavis419/fiber-mongo-example/router"
	"os"

	"github.com/bmdavis419/fiber-mongo-example/models"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	err := run()

	if err != nil {
		panic(err)
	}
}

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

func run() error {
	// init env
	err := common.LoadEnv()
	if err != nil {
		return err
	}

	// init db
	err = common.InitDB()
	if err != nil {
		return err
	}

	err = InitDB2()
	if err != nil {
		return err
	}

	// defer closing db
	defer common.CloseDB()

	// create app
	app := fiber.New()

	// add basic middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// add routes
	router.AddBookGroup(app)
	AddtitlesGroup(app)

	// router.AddtitlesGroup(app)

	// start server
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}
	app.Listen(":" + port)

	return nil
}

var redisClient *redis.Client

func AddtitlesGroup(app *fiber.App) {
	bookGroup := app.Group("/titles")

	bookGroup.Get("/", getstitles)
	bookGroup.Get("/:id", getsTitle)
	// bookGroup.Post("/", createBook)
	// bookGroup.Put("/:id", updateBook)
	// bookGroup.Delete("/:id", deleteBook)
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

func getsTitle(c *fiber.Ctx) error {
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
