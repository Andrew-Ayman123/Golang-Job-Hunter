package handler
import (
	"net/http"
)

// HandleUserLogIn handles user login
// @Summary User Login
// @Description Authenticate user with username and password
// @Tags Authentication
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Router /user/login [post]
func HandleUserLogIn(w http.ResponseWriter, r *http.Request) {
	
	// strip header auth
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	authToken := authHeader[len("ApiKeyAuth "):]
	if authToken != "123" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	// If authenticated, proceed with the request
	w.Write([]byte("Hello World!"))
	w.WriteHeader(http.StatusOK)
	
}

// Hello returns a greeting for the named person.
func Add(a int, b int) int {
	return a + b
}