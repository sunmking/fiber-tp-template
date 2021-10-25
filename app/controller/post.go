package controller

import (
	"fiber-blog/app/model"
	"fiber-blog/app/service"
	"fiber-blog/config"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func HelloWorld(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]string{"a": "你好"})
}

func GetPost(ctx *fiber.Ctx) error {

	query := ctx.Query("id")
	if len(query) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "param invalid")
	}
	id, err := strconv.ParseInt(query, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "param invalid")
	}
	post, err := service.GetPostById(id)
	if err == config.ErrFoo {
		return ctx.SendStatus(fiber.StatusNoContent)
	}
	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(post)
}

func SavePost(ctx *fiber.Ctx) error {
	post := &model.Post{}
	err := ctx.BodyParser(post)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "param invalid")
	}
	err = service.SavePost(post)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return ctx.SendStatus(fiber.StatusOK)
}
