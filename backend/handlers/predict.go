package handlers

import (
	"math"

	"fifa-scout/database"
	"fifa-scout/models"

	"github.com/gofiber/fiber/v2"
)

func PredictPlayer(c *fiber.Ctx) error {
	id := c.Params("id")

	var player models.Player
	if err := database.DB.First(&player, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Player not found"})
	}

	ageFactor := float64(player.Age-15) * 0.8
	if ageFactor <= 0 {
		ageFactor = 0.8
	}

	growthRate := float64(player.Potential-player.Overall) / ageFactor

	const peakAge = 26.0
	yearsLeft := peakAge - float64(player.Age)
	if yearsLeft < 0 {
		yearsLeft = 0
	}
	projectRating := int(math.Min(float64(player.Overall)+growthRate*yearsLeft, 99))

	roiScore := float64(player.Potential) / (float64(player.MarketValue) / 10_000_000)
	isUndervalued := roiScore > 1.2

	badge := "💎 PREMIUM ASSET"
	if roiScore > 1.2 {
		badge = "🔥 UNDERVALUED TARGET"
	} else if roiScore > 0.8 {
		badge = "📊 FAIR VALUE"
	}

	return c.JSON(models.PredictionResponse{
		GrowthRate:    math.Round(growthRate*100) / 100,
		ProjectRating: projectRating,
		IsUndervalued: isUndervalued,
		ROIScore:      math.Round(roiScore*100) / 100,
		Badge:         badge,
	})
}
