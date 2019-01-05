package controllers

import (
	"../models"
	"../utils"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
)

func AuthRoute(router router.Party) {
	router.Get("/auth/tokens", handleTokens)
	router.Post("/auth", handleLogin)
}

func handleLogin(ctx context.Context) {
	login, password := ctx.PostValue("login"), ctx.PostValue("password")

	if login == "" || password == "" {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Login or password is empty"})
		return
	}

	user := models.GetUserByLogin(login)

	if user.ID <= 0 || !utils.CheckPasswordHash(password, user.Password) {
		ctx.StatusCode(422)
		ctx.JSON(iris.Map{"error": "Login or password is incorrect"})
		return
	}

	token, err := user.AddToken(ctx.RemoteAddr())

	if err == nil {
		ctx.JSON(iris.Map{"result": "Success", "token": token})
	} else {
		ctx.JSON(iris.Map{"error": "Error"})
	}
}

func handleTokens(ctx context.Context) {
	tokens, err := models.GetCurrentUser().GetTokens()

	if err == nil {
		ctx.JSON(iris.Map{"result": "Success", "tokens": tokens})
	} else {
		ctx.JSON(iris.Map{"error": "Error"})
	}
}
