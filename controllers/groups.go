package controllers

import (
	"../config"
	"../models"

	"strconv"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
)

func GroupsRoute(router router.Party) {
	// Route -> /api/groups/*
	groupsRoute := router.Party("/groups")

	groupsRoute.Get("/", handleGroupsGet)
	groupsRoute.Get("/{id:int}", handleGroupGet)
	groupsRoute.Post("/", handleGroupPost)
	groupsRoute.Delete("/{id:int}", handleGroupDelete)
}

func handleGroupsGet(ctx context.Context) {
	offset, err := ctx.URLParamInt("offset")
	if err != nil || offset < 0 {
		offset = 0
	}

	count, err := ctx.URLParamInt("count")
	if err != nil || count > 100 || count < 0 {
		count = 30
	}

	var groups []models.Group
	userId := models.GetCurrentUser().ID

	models.GetAllGroups(&groups, userId, count, offset)

	ctx.JSON(iris.Map{"result": "Groups received", "groups": groups})
}

func handleGroupGet(ctx context.Context) {
	groupId, err := strconv.ParseUint(ctx.Params().Get("id"), 10, 64)
	if err != nil || groupId < 1 {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Invalid ID"})
		return
	}

	var group models.Group
	userId := models.GetCurrentUser().ID

	models.GetGroup(&group, uint(groupId), userId, 40)

	if group.ID < 1 {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Group is not found"})
		return
	}

	ctx.JSON(iris.Map{"result": "Group received", "group": group})
}

func handleGroupPost(ctx context.Context) {
	name, description := ctx.PostValue("name"), ctx.PostValue("description")

	if name == "" {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Name is required"})
		return
	}

	if len(name) > 60 {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "[Name] - Max length is 60"})
		return
	}

	if len(description) > 300 {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "[Description] - Max length is 60"})
		return
	}

	group := models.Group{
		Name: name,
		Description: description,
		CreatorID: models.GetCurrentUser().ID,
	}

	isBlank := config.Db.NewRecord(group)

	if isBlank {
		config.Db.Create(&group)
		config.Db.Model(&group).Association("Users").Append(models.GetCurrentUser())
		group.Creator = models.GetCurrentUser()
		ctx.JSON(iris.Map{"result": "Group created", "group": group})
	} else {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Create error"})
	}
}

func handleGroupDelete(ctx context.Context) {
	groupId := ctx.Params().Get("id")

	var group models.Group
	config.Db.First(&group, groupId)

	if group.ID < 1 {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "Group does not exist"})
		return
	}

	currentUser := models.GetCurrentUser()

	if group.CreatorID != currentUser.ID && !currentUser.IsAdmin {
		ctx.StatusCode(iris.StatusUnprocessableEntity)
		ctx.JSON(iris.Map{"error": "No permission to delete"})
		return
	}

	config.Db.Delete(&group)
	ctx.JSON(iris.Map{"result": "Group deleted"})
}
