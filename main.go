package main

import (
	"blogs/model"
	"blogs/schema"
	"html"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type blogBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {
	model.CreateTable()
	app := fiber.New()
	app.Use(recover.New())

	app.Get("/blogs", func(c *fiber.Ctx) error {
		db := model.ConnectDB()
		defer model.CloseDB(db)

		blogs := []schema.Blog{}
		db.Raw("SELECT * FROM blogs").Scan(&blogs)

		var status int
		if len(blogs) > 0 {
			status = fiber.StatusOK
		} else {
			status = fiber.StatusNoContent
		}

		return c.Status(status).JSON(&fiber.Map{
			"success": true,
			"blogs":   blogs,
		})
	})

	app.Post("/blog", func(c *fiber.Ctx) error {
		db := model.ConnectDB()
		defer model.CloseDB(db)

		body := blogBody{}
		if err := c.BodyParser(&body); err != nil {
			panic(err)
		}

		body.Title = html.EscapeString(body.Title)
		body.Description = html.EscapeString(body.Description)

		db.Exec("INSERT INTO blogs (title, description) VALUES (?, ?)", body.Title, body.Description)

		return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
			"success": true,
		})
	})

	app.Put("/blog/:id", func(c *fiber.Ctx) error {
		db := model.ConnectDB()
		defer model.CloseDB(db)

		id, err := c.ParamsInt("id")
		if err != nil {
			panic(err)
		}

		var body blogBody
		if err := c.BodyParser(&body); err != nil {
			panic(err)
		}

		body.Title = html.EscapeString(body.Title)
		body.Description = html.EscapeString(body.Description)

		db.Exec("UPDATE blogs SET title=?, description=? WHERE id=?", body.Title, body.Description, id)

		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
		})
	})

	app.Delete("/blog/:id", func(c *fiber.Ctx) error {
		db := model.ConnectDB()
		defer model.CloseDB(db)

		id, err := c.ParamsInt("id")
		if err != nil {
			panic(err)
		}

		db.Exec("DELETE FROM blogs WHERE id=?", id)

		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
		})
	})

	app.Listen(":8000")
}
