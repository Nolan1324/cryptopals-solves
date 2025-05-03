package main

import "testing"

func TestNotAdmin(t *testing.T) {
	app := makeApplication()
	data, err := app.CreateDataEncrypted("hello")
	if err != nil {
		t.Fatal(err)
	}
	isAdmin, err := app.IsAdmin(data)
	if err != nil {
		t.Fatal(err)
	}
	if isAdmin {
		t.Error("Normal user marked as admin")
	}
}

func TestPreventInjection(t *testing.T) {
	app := makeApplication()
	data, err := app.CreateDataEncrypted("hello;admin=true")
	if err != nil {
		t.Fatal(err)
	}
	isAdmin, err := app.IsAdmin(data)
	if err != nil {
		t.Fatal(err)
	}
	if isAdmin {
		t.Error("Normal user could inject admin")
	}
}

func TestParseAdmin(t *testing.T) {
	isAdmin, err := isAdminDecrypted("comment1=hello;admin=true;comment2=world")
	if err != nil {
		t.Fatal(err)
	}
	if !isAdmin {
		t.Error("Admin string not detected as admin")
	}
}
