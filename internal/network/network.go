package network

import (
	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/gofiber/fiber"
)

// GetNetwork handles /api/network
func GetNetwork(c *fiber.Ctx) {
	db := database.DBConn
	var net Network
	db.Preload("Sites.Tags").Preload("Sites").Preload("Links").Find(&net)
	c.Status(fiber.StatusOK).JSON(net)
}

// NewNetwork handles /api/network
func NewNetwork(c *fiber.Ctx) {
	var net Network
	err := c.BodyParser(&net)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return
	}

	db := database.DBConn
	db.Unscoped().Where("1 = 1").Delete(&Tag{})
	db.Unscoped().Where("1 = 1").Delete(&Site{})
	db.Unscoped().Where("1 = 1").Delete(&Link{})
	db.Unscoped().Where("1 = 1").Delete(&Network{})
	db.Create(&net)

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id": net.ID,
	})
}
