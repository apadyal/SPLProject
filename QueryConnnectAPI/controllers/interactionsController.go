package controllers

import (
	"QueryConnectAPI/models"
	"fmt"

	"github.com/kataras/iris/v12"
)

type InteractionsController struct{}

type SearchRequest struct {
	Tags string `json:"tags"`
}

func (c *InteractionsController) Get() []models.Interaction {
	interactions := models.GetAllInteractions()
	return interactions
}

func (c *InteractionsController) GetByType(ctx iris.Context) []models.Interaction {
	intrType := ctx.URLParam("intrType")
	items := models.GetInteractionsByType(intrType)
	return items
}

// "http://localhost:8080/interactions/owner?ownerid=123"
func (c *InteractionsController) GetInteractionsByOwner(ctx iris.Context) []models.Interaction {

	ownerId := ctx.URLParam("ownerid")

	items := models.GetInteractionsByOwnerID(ownerId)
	return items
}

// "http://localhost:8080/interactions/interaction?interactionobjectid=123"
func (c *InteractionsController) GetInteractionsByInteractionId(ctx iris.Context) []models.Interaction {

	interactionObjectId := ctx.URLParam("interactionobjectid")

	items := models.GetInteractionsByInteractionObjectID(interactionObjectId)
	return items
}

func (c *InteractionsController) Post(ctx iris.Context) {

	var interaction models.Interaction
	err := ctx.ReadJSON(&interaction)
	fmt.Printf(interaction.Description)
	if err != nil {
		fmt.Printf(err.Error())
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request payload"})
		return
	}
	// Add your logic here to save the interaction
	// For example, you could call a method on your model
	res := models.InsertInteraction(interaction)
	fmt.Printf("\n%s", res)
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(interaction)
}

func (c *InteractionsController) Edit(ctx iris.Context) {

	var interaction models.Interaction

	err := ctx.ReadJSON(&interaction)
	fmt.Printf(interaction.Description)
	if err != nil {
		fmt.Printf(err.Error())
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request payload"})
		return
	}
	// Add your logic here to save the interaction
	// For example, you could call a method on your model
	res := models.UpdateInteraction(interaction, false)
	fmt.Printf("\n%s", res)
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(interaction)
}

func (c *InteractionsController) Search(ctx iris.Context) []models.Interaction {
	tags := ctx.URLParam("tags")
	items := models.GetInteractionsByTags(tags)
	return items
}

func (c *InteractionsController) DeleteByID(ctx iris.Context) {
	id := ctx.URLParam("id")
	models.DeleteInteraction(id)
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Item deleted successfully"})
}
