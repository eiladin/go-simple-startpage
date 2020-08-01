package network

import (
	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/gofiber/fiber"
)

// GetNetwork handles /api/network
func GetNetwork(c *fiber.Ctx) {
	db := database.DBConn
	var network Network
	db.Set("gorm:auto_preload", true).Find(&network)
	c.JSON(network)
}

// NewNetwork handles /api/network
func NewNetwork(c *fiber.Ctx) {
	var network Network
	err := c.BodyParser(&network)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return
	}
	type resp struct {
		id uint
	}
	db := database.DBConn
	db.Unscoped().Delete(Network{})
	db.Unscoped().Delete(Site{})
	db.Unscoped().Delete(Tag{})
	db.Unscoped().Delete(Link{})
	db.Create(&network)
	c.JSON(resp{id: network.ID})
}
