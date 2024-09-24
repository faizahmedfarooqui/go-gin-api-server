package controllers

import (
	"database/sql"
	"net/http"

	"api-server/services"
	"api-server/validators"

	"github.com/gin-gonic/gin"
)

type ItemController struct {
	service *services.ItemService
}

func NewItemController(db *sql.DB) *ItemController {
	return &ItemController{
		service: services.NewItemService(db),
	}
}

func (ic *ItemController) GetItems(c *gin.Context) {
	items, err := ic.service.GetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (ic *ItemController) GetItem(c *gin.Context) {
	id := c.Param("id")
	item, err := ic.service.GetItemByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (ic *ItemController) CreateItem(c *gin.Context) {
	var input validators.CreateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := ic.service.CreateItem(input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}
