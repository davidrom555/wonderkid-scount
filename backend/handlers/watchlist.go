package handlers

import (
	"fifa-scout/database"
	"fifa-scout/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetWatchlist(c *fiber.Ctx) error {
	var items []models.WatchlistItem
	database.DB.Preload("Player").Find(&items)
	var players []models.Player
	for _, item := range items {
		players = append(players, item.Player)
	}
	if players == nil {
		players = []models.Player{}
	}
	return c.JSON(fiber.Map{"players": players, "count": len(players)})
}

func AddToWatchlist(c *fiber.Ctx) error {
	id := c.Params("id")

	var player models.Player
	if err := database.DB.First(&player, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Player not found"})
	}

	var existing models.WatchlistItem
	err := database.DB.Where("player_id = ?", id).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		database.DB.Create(&models.WatchlistItem{PlayerID: player.ID})
	}

	return c.JSON(fiber.Map{"success": true})
}

func RemoveFromWatchlist(c *fiber.Ctx) error {
	id := c.Params("id")
	database.DB.Where("player_id = ?", id).Delete(&models.WatchlistItem{})
	return c.JSON(fiber.Map{"success": true})
}
