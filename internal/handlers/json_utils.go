package handlers

import (
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func (h *UserHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// ValidateStruct validates any struct using validator tags
func (h *UserHandler)validateStruct(input interface{}) (map[string]string, error) {
    err := h.validator.Struct(input)
    if err == nil {
        return nil, nil
    }

    validationErrors := make(map[string]string)
    for _, fieldErr := range err.(validator.ValidationErrors) {
        validationErrors[fieldErr.Field()] = fmt.Sprintf("failed validation: %s", fieldErr.Tag())
    }
    return validationErrors, err
}