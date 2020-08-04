package network

import (
	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
	"github.com/gofiber/fiber"
)

// Handler handles Network commands
type Handler struct {
	NetworkService interfaces.NetworkService
}

// GetNetwork handles /api/network
func (h Handler) GetNetwork(c *fiber.Ctx) {
	var net interfaces.Network
	h.NetworkService.FindNetwork(&net)
	c.Status(fiber.StatusOK).JSON(net)
}

// NewNetwork handles /api/network
func (h Handler) NewNetwork(c *fiber.Ctx) {
	var net interfaces.Network
	err := c.BodyParser(&net)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return
	}

	h.NetworkService.CreateNetwork(&net)

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id": net.ID,
	})
}
