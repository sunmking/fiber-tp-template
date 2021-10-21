package main

import (
	"expvar"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	expvarmw "github.com/gofiber/fiber/v2/middleware/expvar"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/template/html"
)

// Field names should start with an uppercase letter
type Person struct {
	Name string `json:"name" xml:"name" form:"name"`
	Pass string `json:"pass" xml:"pass" form:"pass"`
}

var count = expvar.NewInt("count")

func main() {
	engine := html.New("./views", ".html")
	// Pass engine to Fiber's Views Engine
	app := fiber.New(fiber.Config{
		Views: engine,
		// Views Layout is the global layout for all template render until override on Render function.
		ViewsLayout: "layouts/main",
	})
	app.Use(func(c *fiber.Ctx) error {
		// Set some security headers:
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Download-Options", "noopen")
		c.Set("Strict-Transport-Security", "max-age=5184000")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-DNS-Prefetch-Control", "off")

		// Go to next middleware:
		return c.Next()
	})

	file, _ := os.OpenFile("./runtime/log/"+time.Now().Format("2006-01-01")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	defer file.Close()

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Output: file,
		Format: "${pid} - [${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Use(expvarmw.New())
	// Or extend your config for customization
	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     1 * time.Hour,
		KeyGenerator:   utils.UUID,
	}))
	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://gofiber.io, https://gofiber.net",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	// Provide a custom compression level
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

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

	log.Fatal(app.Listen(":3000"))
}
