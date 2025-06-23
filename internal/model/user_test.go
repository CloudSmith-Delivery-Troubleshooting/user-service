package model

import (
	"encoding/json"
	"testing"
)

// TestUserJSONSerialization ensures that User structs serialize to JSON correctly
func TestUserJSONSerialization(t *testing.T) {
	user := User{
		Email: "test@example.com",
		Name:  "Test User",
		Age:   30,
	}

	// Test marshaling
	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	// Verify JSON structure
	expected := `{"email":"test@example.com","name":"Test User","age":30}`
	if string(data) != expected {
		t.Errorf("JSON serialization incorrect\nExpected: %s\nGot: %s", expected, string(data))
	}

	// Test unmarshaling
	var unmarshaled User
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal user: %v", err)
	}

	// Verify fields
	if unmarshaled.Email != user.Email {
		t.Errorf("Email field mismatch: expected %s, got %s", user.Email, unmarshaled.Email)
	}
	if unmarshaled.Name != user.Name {
		t.Errorf("Name field mismatch: expected %s, got %s", user.Name, unmarshaled.Name)
	}
	if unmarshaled.Age != user.Age {
		t.Errorf("Age field mismatch: expected %d, got %d", user.Age, unmarshaled.Age)
	}
}
