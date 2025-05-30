package env
import (
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

func Init() {
	// Load environment variables from a .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Overload(".env"); err != nil {
			panic("Error loading .env file")
		}
	}
}
// GetEnv retrieves the value of the environment variable named by the key.
func GetEnv(key string,fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
// GetEnvAsInt retrieves the value of the environment variable named by the key and converts it to an integer.
func GetEnvAsInt(key string,fallback int) (int) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	
	return intValue
}

