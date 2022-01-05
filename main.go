package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gyujae/starcraft_scrapper/scrapper"
)

func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders: "Access-Control-Allow-Origin",
	}))

	app.Get("/maps", func(c *fiber.Ctx) error {
		maps := scrapper.MapScraper()
		return c.JSON(maps)
	})

	app.Get("/asl_maps", func(c *fiber.Ctx) error {
		maps := scrapper.ASLMapScrapper()
		return c.JSON(maps)
	})

	log.Fatal(app.Listen(":4000"))

}
