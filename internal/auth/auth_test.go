package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordOne(t *testing.T) {
	password := "Test123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("error hashing password| pass: %s, error: %v", password, err)
	}
	err = CheckPasswordHash(hash, password)
	if err != nil {
		t.Errorf("password doesn't match hash: %v", err)
	}
}

func TestCheckPasswordEmpty(t *testing.T) {
	password := " "
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("error hashing password| pass: %s, error: %v", password, err)
	}
	err = CheckPasswordHash(hash, password)
	if err != nil {
		t.Errorf("password doesn't match hash: %v", err)
	}
}

func TestCheckPasswordSpaces(t *testing.T) {
	password := "this pass has spaces"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("error hashing password| pass: %s, error: %v", password, err)
	}
	err = CheckPasswordHash(hash, password)
	if err != nil {
		t.Errorf("password doesn't match hash: %v", err)
	}
}

func TestCheckJWTBase(t *testing.T) {
	newID, err := uuid.NewUUID()

	if err != nil {
		t.Errorf("error generating UUID: %v", err)
	}
	duration := 5 * time.Minute

	tokenSecret := "testme"
	tokenString, err := MakeJWT(newID, tokenSecret, duration)
	if err != nil {
		t.Errorf("error making JWT: %v", err)
	}

	_, err = ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Errorf("error validating JWT: %v", err)
	}
}

func TestCheckJWTBadSecret(t *testing.T) {
	newID, err := uuid.NewUUID()

	if err != nil {
		t.Errorf("error generating UUID: %v", err)
	}
	duration := 5 * time.Minute

	tokenSecret := "testme"
	tokenString, err := MakeJWT(newID, tokenSecret, duration)
	if err != nil {
		t.Errorf("error making JWT: %v", err)
	}
	wrongSecret := "notThis"
	_, err = ValidateJWT(tokenString, wrongSecret)
	if err == nil {
		t.Errorf("validated with incorrect secret: %v", err)
	}
}

func TestCheckJWTWaiting(t *testing.T) {
	newID, err := uuid.NewUUID()

	if err != nil {
		t.Errorf("error generating UUID: %v", err)
	}
	duration := 5 * time.Microsecond

	tokenSecret := "testme"
	tokenString, err := MakeJWT(newID, tokenSecret, duration)
	if err != nil {
		t.Errorf("error making JWT: %v", err)
	}

	time.Sleep(1 * time.Second)
	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Errorf("validated after time expired: %v", err)
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
