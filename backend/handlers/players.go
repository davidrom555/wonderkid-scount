package handlers

import (
	"fmt"
	"fifa-scout/database"
	"fifa-scout/models"

	"github.com/gofiber/fiber/v2"
)

func GetPlayers(c *fiber.Ctx) error {
	var players []models.Player
	database.DB.Find(&players)
	fmt.Printf("[API] Returning %d players from DB\n", len(players))
	return c.JSON(players)
}

func FilterPlayers(c *fiber.Ctx) error {
	minPotential := c.QueryInt("minPotential", 0)
	position := c.Query("position", "")
	search := c.Query("search", "")

	query := database.DB.Model(&models.Player{})
	if minPotential > 0 {
		query = query.Where("potential >= ?", minPotential)
	}
	if position != "" {
		query = query.Where("position = ?", position)
	}
	if search != "" {
		query = query.Where("name LIKE ? OR club LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	var players []models.Player
	query.Find(&players)
	return c.JSON(players)
}
