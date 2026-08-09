package utils

import (
	"testing"
)

func TestHashedPassword(t *testing.T) {
	password := "securepassword123"
	hashed, err := HashedPassword(password)
	if err != nil {
		t.Fatalf("HashedPassword failed: %v", err)
	}
	if hashed == "" {
		t.Fatal("expected non-empty hash")
	}
	if hashed == password {
		t.Fatal("hash should not equal plaintext password")
	}
}

func TestCheckPassword_Success(t *testing.T) {
	password := "securepassword123"
	hashed, err := HashedPassword(password)
	if err != nil {
		t.Fatalf("HashedPassword failed: %v", err)
	}
	if err := CheckPassword(password, hashed); err != nil {
		t.Fatalf("CheckPassword should succeed for correct password: %v", err)
	}
}

func TestCheckPassword_Failure(t *testing.T) {
	password := "securepassword123"
	hashed, err := HashedPassword(password)
	if err != nil {
		t.Fatalf("HashedPassword failed: %v", err)
	}
	if err := CheckPassword("wrongpassword", hashed); err == nil {
		t.Fatal("CheckPassword should fail for incorrect password")
	}
}

func TestGenerateRandomID(t *testing.T) {
	id, err := GenerateRandomID(8)
	if err != nil {
		t.Fatalf("GenerateRandomID failed: %v", err)
	}
	if len(id) != 16 { // 8 bytes = 16 hex chars
		t.Fatalf("expected length 16, got %d", len(id))
	}
}

func TestGenerateRandomID_InvalidLen(t *testing.T) {
	_, err := GenerateRandomID(0)
	if err == nil {
		t.Fatal("expected error for byteLen <= 0")
	}
}

func TestMustRandomID(t *testing.T) {
	id := MustRandomID()
	if id == "" {
		t.Fatal("MustRandomID returned empty string")
	}
	if len(id) != 16 {
		t.Fatalf("expected length 16, got %d", len(id))
	}
}
