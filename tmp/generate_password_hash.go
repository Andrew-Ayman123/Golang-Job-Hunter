package main




import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// This is a simple Go program that hashes a password using bcrypt.
// to ingest a demo admin in the database.
func main2() {
	password := "12345678"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}
	fmt.Println("Hashed Password:", string(hashedPassword))
}