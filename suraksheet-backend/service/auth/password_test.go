package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password %v", err)
	}
	if hash == "" {
		t.Error("expected hash to be not empty")
	}
	if hash == "password" {
		t.Error("expected hash to be different from password")
	}
}

func TestComparePassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password %v", err)
	}
	res := ComparePassword(hash, "password")
	if res != true {
		t.Errorf("expected match to be true but got %v", res)
	}
	res = ComparePassword(hash, "passwd")
	if res == true {
		t.Errorf("expected match to be false but got %v", res)
	}
}
