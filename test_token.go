package main
import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)
func main() {
	godotenv.Load()
	fmt.Println("TOKEN IS:", os.Getenv("AI_INTERNAL_TOKEN"))
}
