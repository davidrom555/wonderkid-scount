package models

import "time"

type Player struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Position    string    `json:"position" gorm:"size:10;not null"`
	Age         int       `json:"age" gorm:"not null"`
	Overall     int       `json:"overall" gorm:"not null"`
	Potential   int       `json:"potential" gorm:"not null"`
	MarketValue int64     `json:"marketValue" gorm:"not null"`
	Club        string    `json:"club" gorm:"size:100"`
	Nationality string    `json:"nationality" gorm:"size:100"`
	SofifaID    int       `json:"sofifaId" gorm:"not null"`
	PAC         int       `json:"pac" gorm:"not null"`
	SHO         int       `json:"sho" gorm:"not null"`
	PAS         int       `json:"pas" gorm:"not null"`
	DRI         int       `json:"dri" gorm:"not null"`
	DEF         int       `json:"def" gorm:"not null"`
	PHY         int       `json:"phy" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt"`
}

type PredictionResponse struct {
	GrowthRate    float64 `json:"growthRate"`
	ProjectRating int     `json:"projectRating"`
	IsUndervalued bool    `json:"isUndervalued"`
	ROIScore      float64 `json:"roiScore"`
	Badge         string  `json:"badge"`
}
