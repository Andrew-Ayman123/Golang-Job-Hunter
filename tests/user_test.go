package tests

import (
	"testing"
	"github.com/Andrew-Ayman123/GoProject/handler"
)

// TestUserLogin tests the user login functionality.
func TestAdd(t *testing.T) {
	result := handler.Add(2, 3)
	expected := 5

	if result != expected {
		t.Errorf("Add(2, 3) = %d; want %d", result, expected)
	}
}