package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/risipatel9/go-mod/prisma/db"
)

type Todo struct {
	Id        int    `json:"id"`
	Completed bool   `json:"completed"`
	Title     string `json:"title"`
	Desc      string `json:"Desc"`
}

func main() {
	fmt.Println("hello world")
	app := fiber.New()
	todos := []Todo{}

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to load .env file")
	}

	PORT := os.Getenv("PORT")
	client := db.NewClient()

	if err := client.Prisma.Connect(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer client.Prisma.Disconnect()

	defer client.Disconnect()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "hello"})
	})

	app.Post("/api/todo", func(c *fiber.Ctx) error {
		todo := &Todo{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}
		fmt.Println(todo)

		if todo.Title == "" {
			return c.Status(400).JSON(fiber.Map{"Error": "Title is empty"})
		}
		newTodo, err := client.Todo.CreateOne(
			db.Todo.Title.Set(todo.Title),
			db.Todo.Desc.Set(todo.Desc),
			db.Todo.Completed.Set(false),
		).Exec(context.Background())

		fmt.Println(newTodo)
		fmt.Println(err)

		return c.Status(201).JSON(fiber.Map{"Success": "Added Todo Successfully"})
	})

	app.Get("/api/todo/", func(c *fiber.Ctx) error {
		todos, err := client.Todo.FindMany().Exec(context.Background())

		fmt.Println(todos, err)
		return c.Status(200).JSON(fiber.Map{"Todos": todos})
	})

	app.Patch("/api/todo/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		flag := false

		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos[i].Completed = true
				flag = true
			}
		}
		if !flag {
			return c.Status(404).JSON(fiber.Map{"Error": "There is no Todo associated with "})
		}

		return c.Status(201).JSON(fiber.Map{"Success": "Updated todo Successfully"})
	})

	app.Delete("/api/todo/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(201).JSON(fiber.Map{"Success": "Deleted todo Successfully"})
			}
		}
		return c.Status(404).JSON(fiber.Map{"Error": "There is no Todo associated with"})
	})
	log.Fatal((app.Listen(":" + PORT)))
}
