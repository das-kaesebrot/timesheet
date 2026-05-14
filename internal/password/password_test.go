package password

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("mySecurePassword123")
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned an empty hash")
	}
}

func TestPasswordMatches(t *testing.T) {
	password := "mySecurePassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if !PasswordMatches(password, hash) {
		t.Error("PasswordMatches returned false for the correct password")
	}
}

func TestPasswordMatchesWrongPassword(t *testing.T) {
	hash, err := HashPassword("correctPassword")
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if PasswordMatches("wrongPassword", hash) {
		t.Error("PasswordMatches returned true for an incorrect password")
	}
}

func TestPasswordMatchesEmptyPassword(t *testing.T) {
	hash, err := HashPassword("somePassword")
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if PasswordMatches("", hash) {
		t.Error("PasswordMatches returned true for an empty password")
	}
}

func TestHashPasswordEmptyString(t *testing.T) {
	hash, err := HashPassword("")
	if err != nil {
		t.Fatalf("HashPassword returned an error for empty string: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned an empty hash for empty input")
	}
}
