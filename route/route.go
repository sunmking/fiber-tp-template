package route

import (
	"expvar"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var count = expvar.NewInt("count")

func New(app *fiber.App) *fiber.App {
	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("1st route!")
		return c.Next()
	})

	app.Get("*", func(c *fiber.Ctx) error {
		fmt.Println("2nd route!")
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("3rd route!")
		count.Add(1)
		//panic("I'm an error")
		return c.SendString(fmt.Sprintf("hello expvar count %d", count.Value()))
	})

	app.Get("/teapot", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusTeapot).SendString("🍵 short and stout 🍵")
	})

	app.Get("/v1/some/resource/name\\:customVerb", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Community")
	})

	v1 := app.Group("/v1")

	v1.Get("/a2", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index/index", fiber.Map{
			"Title": "Hello, World!",
		})
	})
	return app
}
