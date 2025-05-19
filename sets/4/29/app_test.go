package main

import "testing"

func TestSignature(t *testing.T) {
	app := MakeApplication()
	data, signature := app.CreateSignedData()
	if len(signature) != 20 {
		t.Fatalf("expected signature length 20, actual length %v", len(signature))
	}
	isAdmin, err := app.IsAdmin(data, signature)
	if err != nil {
		t.Fatalf("error checking if admin: %v", err)
	}
	if isAdmin {
		t.Errorf("user incorrectly recognized as admin")
	}
}

func TestTamperedData(t *testing.T) {
	app := MakeApplication()
	data, signature := app.CreateSignedData()
	// Tamper with the data
	data[len(data)-1] ^= 1
	_, err := app.IsAdmin(data, signature)
	if err == nil {
		t.Fatalf("no error returned")
	}
}

func TestBadSignature(t *testing.T) {
	app := MakeApplication()
	data, signature := app.CreateSignedData()
	signature[len(signature)-1] ^= 1
	_, err := app.IsAdmin(data, signature)
	if err == nil {
		t.Fatalf("no error returned")
	}
}
