package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tsivinsky/fileasy/internal/db"
	"github.com/tsivinsky/fileasy/internal/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Static("/", "./static", fiber.Static{})

	app.Get("/api/files", router.HandleListAllFiles)
	app.Get("/api/:name", router.HandleFindFileByName)

	app.Post("/api/upload", router.HandleUploadFile)

	log.Fatal(app.Listen(":5000"))
}
