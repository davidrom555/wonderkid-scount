package models

type WatchlistItem struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	PlayerID uint   `json:"playerId" gorm:"uniqueIndex;not null"`
	Player   Player `json:"player" gorm:"foreignKey:PlayerID"`
}
