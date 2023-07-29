package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/tsivinsky/fileasy/internal/db"
	"github.com/tsivinsky/fileasy/internal/middleware"
	"github.com/tsivinsky/fileasy/internal/router"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(500).JSON(fiber.Map{
		"error": err.Error(),
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Use(cors.New())
	app.Use(recover.New())

	app.Static("/", "./static", fiber.Static{})

	app.Get("/api/auth/github", router.HandleGitHubLogin)
	app.Get("/api/auth/github/callback", router.HandleGitHubCallback)
	app.Post("/api/auth/refresh", router.HandleGetNewAccessToken)

	app.Get("/api/files", middleware.VerifyJWTToken, router.HandleListAllFiles)
	app.Get("/api/:name", middleware.VerifyJWTToken, router.HandleFindFileByName)

	app.Post("/api/upload", middleware.VerifyJWTToken, router.HandleUploadFile)

	log.Fatal(app.Listen(":5000"))
}
