package main
import (
	"fmt"
	"os"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"github.com/joho/godotenv"
)
func main() {
	godotenv.Load()
	dbURL := os.Getenv("DIRECT_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		dbURL = "postgresql://postgres.flfgtvtrxufyeqtwbolr:.IronHack1234@aws-1-eu-central-1.pooler.supabase.com:5432/postgres"
	}
	db.ConnectDB(dbURL)

	var maps []models.Map
	db.DB.Find(&maps)
	fmt.Printf("Total maps: %d\n", len(maps))
	for _, m := range maps {
		fmt.Printf("ID: %d | Name: %s | Layout: %s | Dictionary: %s\n", m.ID, m.Name, m.Layout, m.Dictionary)
	}
}
