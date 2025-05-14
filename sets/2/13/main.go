package main

import (
	"cryptopals/internal/slicex"
	"fmt"
)

func craftAdminProfile(app Application) []byte {
	const bs = 16

	// Goal is to produce the following string in the first two blocks
	// email=mmail@foo.
	// bar&uid=10&role=
	ct, err := app.EncryptedProfileFor("mmail@foo.bar")
	if err != nil {
		panic(err)
	}
	// email=mmail@foo.bar&uid=10&role=
	craftedProfile := ct[0 : 2*bs]

	// Goal is to produce the following string in the first blocks
	// email=mmail@foo.
	// admin[PKCS7 PAD]
	// bs - len("admin") = 11 for PKCS7 padding
	ct, err = app.EncryptedProfileFor(
		"mmail@foo." + string(append([]byte("admin"), slicex.Repeat(byte(11), 11)...)),
	)
	if err != nil {
		panic(err)
	}
	// email=mmail@foo.bar&uid=10&role=admin[PKCS7 PAD]
	craftedProfile = append(craftedProfile, ct[bs:2*bs]...)

	return craftedProfile
}

func main() {
	app := makeApplication()
	craftedProfile := craftAdminProfile(app)
	isAdmin, err := app.IsAdmin(craftedProfile)
	if err != nil {
		panic(err)
	}
	if isAdmin {
		fmt.Println("Authenticated!")
	} else {
		fmt.Println("Profile is not admin")
	}
}
