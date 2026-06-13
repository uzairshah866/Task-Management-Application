package auth

import (
	"testing"

	"github.com/user/taskapp/internal/models"
)

func TestHashPassword(t *testing.T) {
	password := "test123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" {
		t.Error("Hash is empty")
	}

	if hash == password {
		t.Error("Hash should not equal plain password")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "test123"
	hash, _ := HashPassword(password)

	if !VerifyPassword(hash, password) {
		t.Error("Password verification failed for correct password")
	}

	if VerifyPassword(hash, "wrongpassword") {
		t.Error("Password verification should fail for incorrect password")
	}
}

func TestGenerateToken(t *testing.T) {
	user := &models.User{
		ID:    "test-user-id",
		Email: "test@example.com",
	}

	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Token is empty")
	}

	claims, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("Failed to verify token: %v", err)
	}

	if claims["user_id"] != user.ID {
		t.Error("Token claims do not match user ID")
	}

	if claims["email"] != user.Email {
		t.Error("Token claims do not match email")
	}
}
