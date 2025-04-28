package main

import (
	"testing"
)

func TestProfileFor(t *testing.T) {
	app := makeApplication()
	profile, err := app.ProfileFor("foo@bar.com")
	if err != nil {
		t.Fatal(err)
	}
	if profile != "email=foo@bar.com&uid=10&role=user" {
		t.Error("incorrect profile")
	}
}

func TestProfileForbiddenChars(t *testing.T) {
	app := makeApplication()
	_, err := app.ProfileFor("foo&role=admin")
	if err == nil {
		t.Error("error was not thrown by ProfileFor on forbidden characters")
	}
}

func TestRejectNormalUser(t *testing.T) {
	app := makeApplication()
	profile, err := app.EncryptedProfileFor("foo@bar.com")
	if err != nil {
		t.Fatal(err)
	}
	isAdmin, err := app.IsAdmin(profile)
	if err != nil {
		t.Fatal(err)
	}
	if isAdmin {
		t.Error("normal user is recognized as admin")
	}
}
