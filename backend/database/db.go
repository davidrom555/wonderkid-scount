package database

import (
	"fmt"
	"log"
	"time"

	"fifa-scout/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "root:@tcp(localhost:3306)/fifa_scouting?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[DB] Cannot connect to MariaDB: %v", err)
	}

	DB.AutoMigrate(&models.Player{}, &models.WatchlistItem{})

	var count int64
	DB.Model(&models.Player{}).Count(&count)
	if int(count) != len(seedPlayers) {
		fmt.Printf("[DB] Re-syncing: %d → %d players\n", count, len(seedPlayers))
		DB.Exec("DELETE FROM players")
		for _, p := range seedPlayers {
			DB.Create(&p)
		}
	}

	fmt.Println("[DB] MariaDB connected successfully")
}

var seedPlayers = []models.Player{
	{ID: 1,  Name: "Jude Bellingham", Position: "CAM", Age: 20, Overall: 88, Potential: 94, MarketValue: 180000000, Club: "Real Madrid",  Nationality: "England",     SofifaID: 252371, PAC: 82, SHO: 82, PAS: 87, DRI: 90, DEF: 68, PHY: 80, CreatedAt: time.Now()},
	{ID: 2,  Name: "E. Haaland",      Position: "ST",  Age: 23, Overall: 91, Potential: 95, MarketValue: 200000000, Club: "Man City",     Nationality: "Norway",      SofifaID: 239085, PAC: 89, SHO: 94, PAS: 66, DRI: 82, DEF: 44, PHY: 88, CreatedAt: time.Now()},
	{ID: 3,  Name: "Gavi",            Position: "CM",  Age: 19, Overall: 85, Potential: 92, MarketValue: 90000000,  Club: "Barcelona",    Nationality: "Spain",       SofifaID: 264240, PAC: 71, SHO: 74, PAS: 87, DRI: 88, DEF: 76, PHY: 72, CreatedAt: time.Now()},
	{ID: 4,  Name: "Pedri",           Position: "CAM", Age: 21, Overall: 86, Potential: 93, MarketValue: 120000000, Club: "Barcelona",    Nationality: "Spain",       SofifaID: 251854, PAC: 73, SHO: 76, PAS: 88, DRI: 91, DEF: 71, PHY: 71, CreatedAt: time.Now()},
	{ID: 5,  Name: "J. Musiala",      Position: "CAM", Age: 20, Overall: 86, Potential: 93, MarketValue: 130000000, Club: "Bayern",       Nationality: "Germany",     SofifaID: 256790, PAC: 83, SHO: 78, PAS: 82, DRI: 90, DEF: 62, PHY: 68, CreatedAt: time.Now()},
	{ID: 6,  Name: "B. Saka",         Position: "RW",  Age: 22, Overall: 87, Potential: 91, MarketValue: 150000000, Club: "Arsenal",      Nationality: "England",     SofifaID: 246669, PAC: 90, SHO: 83, PAS: 84, DRI: 88, DEF: 59, PHY: 71, CreatedAt: time.Now()},
	{ID: 7,  Name: "V. Koopmeiners",  Position: "CM",  Age: 25, Overall: 84, Potential: 88, MarketValue: 70000000,  Club: "Juventus",     Nationality: "Netherlands", SofifaID: 240679, PAC: 74, SHO: 80, PAS: 84, DRI: 82, DEF: 74, PHY: 78, CreatedAt: time.Now()},
	{ID: 8,  Name: "M. Akanji",       Position: "CB",  Age: 28, Overall: 87, Potential: 89, MarketValue: 65000000,  Club: "Man City",     Nationality: "Switzerland", SofifaID: 229237, PAC: 78, SHO: 49, PAS: 71, DRI: 68, DEF: 87, PHY: 84, CreatedAt: time.Now()},
	{ID: 9,  Name: "R. Dias",         Position: "CB",  Age: 26, Overall: 88, Potential: 91, MarketValue: 85000000,  Club: "Man City",     Nationality: "Portugal",    SofifaID: 239818, PAC: 77, SHO: 46, PAS: 72, DRI: 70, DEF: 89, PHY: 85, CreatedAt: time.Now()},
	{ID: 10, Name: "L. Yamal",        Position: "RW",  Age: 17, Overall: 82, Potential: 96, MarketValue: 180000000, Club: "Barcelona",    Nationality: "Spain",       SofifaID: 277643, PAC: 92, SHO: 77, PAS: 80, DRI: 89, DEF: 33, PHY: 58, CreatedAt: time.Now()},
	{ID: 11, Name: "W. Foden",        Position: "CAM", Age: 23, Overall: 87, Potential: 92, MarketValue: 140000000, Club: "Man City",     Nationality: "England",     SofifaID: 209658, PAC: 81, SHO: 83, PAS: 87, DRI: 89, DEF: 58, PHY: 69, CreatedAt: time.Now()},
	{ID: 12, Name: "K. Mbappé",       Position: "ST",  Age: 25, Overall: 92, Potential: 94, MarketValue: 220000000, Club: "Real Madrid",  Nationality: "France",      SofifaID: 231747, PAC: 97, SHO: 91, PAS: 82, DRI: 93, DEF: 45, PHY: 78, CreatedAt: time.Now()},
	{ID: 13, Name: "V. Junior",       Position: "LW",  Age: 23, Overall: 89, Potential: 92, MarketValue: 200000000, Club: "Real Madrid",  Nationality: "Brazil",      SofifaID: 238794, PAC: 95, SHO: 84, PAS: 79, DRI: 93, DEF: 32, PHY: 70, CreatedAt: time.Now()},
	{ID: 14, Name: "J. Sancho",       Position: "LW",  Age: 24, Overall: 84, Potential: 87, MarketValue: 40000000,  Club: "Man United",   Nationality: "England",     SofifaID: 233049, PAC: 83, SHO: 78, PAS: 80, DRI: 86, DEF: 44, PHY: 63, CreatedAt: time.Now()},
	{ID: 15, Name: "C. Gakpo",        Position: "LW",  Age: 24, Overall: 86, Potential: 91, MarketValue: 75000000,  Club: "Liverpool",    Nationality: "Netherlands", SofifaID: 242516, PAC: 84, SHO: 82, PAS: 79, DRI: 84, DEF: 54, PHY: 72, CreatedAt: time.Now()},
	{ID: 16, Name: "F. De Jong",      Position: "CDM", Age: 26, Overall: 86, Potential: 90, MarketValue: 80000000,  Club: "Barcelona",    Nationality: "Netherlands", SofifaID: 228702, PAC: 74, SHO: 71, PAS: 87, DRI: 86, DEF: 78, PHY: 74, CreatedAt: time.Now()},
	{ID: 17, Name: "A. Guler",        Position: "CAM", Age: 19, Overall: 79, Potential: 88, MarketValue: 45000000,  Club: "Real Madrid",  Nationality: "Turkey",      SofifaID: 272123, PAC: 77, SHO: 80, PAS: 83, DRI: 88, DEF: 52, PHY: 62, CreatedAt: time.Now()},
	{ID: 18, Name: "P. Cubarsi",      Position: "CB",  Age: 17, Overall: 79, Potential: 90, MarketValue: 35000000,  Club: "Barcelona",    Nationality: "Spain",       SofifaID: 278854, PAC: 72, SHO: 40, PAS: 70, DRI: 70, DEF: 84, PHY: 78, CreatedAt: time.Now()},
	{ID: 19, Name: "D. Olmo",         Position: "CAM", Age: 26, Overall: 85, Potential: 88, MarketValue: 80000000,  Club: "Barcelona",    Nationality: "Spain",       SofifaID: 222665, PAC: 78, SHO: 79, PAS: 85, DRI: 87, DEF: 64, PHY: 70, CreatedAt: time.Now()},
	{ID: 20, Name: "R. Hojlund",      Position: "ST",  Age: 21, Overall: 80, Potential: 89, MarketValue: 55000000,  Club: "Man United",   Nationality: "Denmark",     SofifaID: 265243, PAC: 82, SHO: 80, PAS: 64, DRI: 78, DEF: 36, PHY: 78, CreatedAt: time.Now()},
	{ID: 21, Name: "D. Frimpong",     Position: "RB",  Age: 23, Overall: 82, Potential: 87, MarketValue: 45000000,  Club: "Leverkusen",   Nationality: "Netherlands", SofifaID: 245359, PAC: 90, SHO: 68, PAS: 72, DRI: 82, DEF: 75, PHY: 72, CreatedAt: time.Now()},
	{ID: 22, Name: "W. Zaire-Emery",  Position: "CM",  Age: 19, Overall: 80, Potential: 91, MarketValue: 55000000,  Club: "PSG",          Nationality: "France",      SofifaID: 271920, PAC: 76, SHO: 72, PAS: 82, DRI: 84, DEF: 73, PHY: 74, CreatedAt: time.Now()},
	{ID: 23, Name: "D. Doue",         Position: "LW",  Age: 19, Overall: 80, Potential: 90, MarketValue: 55000000,  Club: "PSG",          Nationality: "France",      SofifaID: 272538, PAC: 88, SHO: 76, PAS: 75, DRI: 85, DEF: 38, PHY: 64, CreatedAt: time.Now()},
	{ID: 24, Name: "J. Duran",        Position: "ST",  Age: 20, Overall: 77, Potential: 87, MarketValue: 35000000,  Club: "Aston Villa",  Nationality: "Colombia",    SofifaID: 265262, PAC: 81, SHO: 77, PAS: 58, DRI: 73, DEF: 28, PHY: 80, CreatedAt: time.Now()},
	{ID: 25, Name: "O. Nkunku",       Position: "CAM", Age: 26, Overall: 84, Potential: 87, MarketValue: 60000000,  Club: "Chelsea",      Nationality: "France",      SofifaID: 244985, PAC: 82, SHO: 83, PAS: 77, DRI: 84, DEF: 44, PHY: 72, CreatedAt: time.Now()},
	{ID: 26, Name: "K. Adeyemi",      Position: "LW",  Age: 23, Overall: 82, Potential: 87, MarketValue: 45000000,  Club: "Dortmund",     Nationality: "Germany",     SofifaID: 256952, PAC: 93, SHO: 76, PAS: 68, DRI: 82, DEF: 40, PHY: 65, CreatedAt: time.Now()},
	{ID: 27, Name: "W. Pacho",        Position: "CB",  Age: 23, Overall: 80, Potential: 87, MarketValue: 35000000,  Club: "PSG",          Nationality: "Ecuador",     SofifaID: 263882, PAC: 78, SHO: 38, PAS: 68, DRI: 66, DEF: 84, PHY: 82, CreatedAt: time.Now()},
	{ID: 28, Name: "T. Fofana",       Position: "CDM", Age: 23, Overall: 82, Potential: 86, MarketValue: 50000000,  Club: "Chelsea",      Nationality: "France",      SofifaID: 239003, PAC: 76, SHO: 60, PAS: 75, DRI: 76, DEF: 82, PHY: 80, CreatedAt: time.Now()},
	{ID: 29, Name: "M. Balogun",      Position: "ST",  Age: 23, Overall: 79, Potential: 86, MarketValue: 30000000,  Club: "Monaco",       Nationality: "USA",         SofifaID: 257244, PAC: 83, SHO: 77, PAS: 64, DRI: 78, DEF: 30, PHY: 68, CreatedAt: time.Now()},
	{ID: 30, Name: "Lautaro M.",      Position: "ST",  Age: 26, Overall: 89, Potential: 91, MarketValue: 110000000, Club: "Inter",        Nationality: "Argentina",   SofifaID: 237692, PAC: 78, SHO: 89, PAS: 74, DRI: 84, DEF: 44, PHY: 82, CreatedAt: time.Now()},
	{ID: 31, Name: "M. Ødegaard",     Position: "CAM", Age: 26, Overall: 88, Potential: 90, MarketValue: 100000000, Club: "Arsenal",      Nationality: "Norway",      SofifaID: 230621, PAC: 78, SHO: 80, PAS: 90, DRI: 88, DEF: 64, PHY: 68, CreatedAt: time.Now()},
	{ID: 32, Name: "B. Fernandes",    Position: "CAM", Age: 29, Overall: 87, Potential: 88, MarketValue: 70000000,  Club: "Man United",   Nationality: "Portugal",    SofifaID: 212831, PAC: 75, SHO: 82, PAS: 88, DRI: 84, DEF: 62, PHY: 70, CreatedAt: time.Now()},
	{ID: 33, Name: "V. Osimhen",      Position: "ST",  Age: 25, Overall: 88, Potential: 90, MarketValue: 90000000,  Club: "Galatasaray",  Nationality: "Nigeria",     SofifaID: 235212, PAC: 91, SHO: 87, PAS: 68, DRI: 82, DEF: 36, PHY: 82, CreatedAt: time.Now()},
	{ID: 34, Name: "J. Doku",         Position: "LW",  Age: 22, Overall: 83, Potential: 88, MarketValue: 55000000,  Club: "Man City",     Nationality: "Belgium",     SofifaID: 258517, PAC: 96, SHO: 72, PAS: 74, DRI: 88, DEF: 36, PHY: 60, CreatedAt: time.Now()},
	{ID: 35, Name: "M. Salah",        Position: "RW",  Age: 32, Overall: 90, Potential: 90, MarketValue: 50000000,  Club: "Liverpool",    Nationality: "Egypt",       SofifaID: 209331, PAC: 88, SHO: 88, PAS: 82, DRI: 88, DEF: 44, PHY: 76, CreatedAt: time.Now()},
	{ID: 36, Name: "H. Kane",         Position: "ST",  Age: 30, Overall: 90, Potential: 90, MarketValue: 80000000,  Club: "Bayern",       Nationality: "England",     SofifaID: 202126, PAC: 70, SHO: 91, PAS: 84, DRI: 80, DEF: 48, PHY: 83, CreatedAt: time.Now()},
	{ID: 37, Name: "M. Ugarte",       Position: "CDM", Age: 23, Overall: 80, Potential: 87, MarketValue: 40000000,  Club: "Man United",   Nationality: "Uruguay",     SofifaID: 263745, PAC: 72, SHO: 60, PAS: 76, DRI: 74, DEF: 82, PHY: 82, CreatedAt: time.Now()},
	{ID: 38, Name: "Alisson",         Position: "GK",  Age: 31, Overall: 90, Potential: 90, MarketValue: 50000000,  Club: "Liverpool",    Nationality: "Brazil",      SofifaID: 211110, PAC: 87, SHO: 14, PAS: 50, DRI: 13, DEF: 18, PHY: 88, CreatedAt: time.Now()},
	{ID: 39, Name: "E. Martínez",     Position: "GK",  Age: 32, Overall: 89, Potential: 89, MarketValue: 45000000,  Club: "Aston Villa",  Nationality: "Argentina",   SofifaID: 204931, PAC: 84, SHO: 13, PAS: 48, DRI: 12, DEF: 17, PHY: 88, CreatedAt: time.Now()},
}
