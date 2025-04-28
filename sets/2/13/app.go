package main

import (
	"cryptopals/internal/cipherx"
	"cryptopals/internal/randx"
	"fmt"
	"net/url"
	"strings"
)

type Application struct {
	key []byte
}

func makeApplication() Application {
	return Application{key: randx.RandBytes(16)}
}

func (a Application) ProfileFor(email string) (string, error) {
	if strings.ContainsAny(email, "&=") {
		return "", fmt.Errorf("email contains illegal character")
	}

	// Would be much more secure to do Values.Encode here, but making this somewhat insecure
	//  for the sake of the challenge
	return fmt.Sprintf("email=%s&uid=10&role=user", email), nil

	// q := url.Values{}
	// q.Set("email", email)
	// q.Set("uid", "10")
	// q.Set("role", "user")
	// return q.Encode(), nil
}

func (a Application) EncryptedProfileFor(email string) ([]byte, error) {
	profileStr, err := a.ProfileFor(email)
	if err != nil {
		return nil, err
	}
	profile := cipherx.Pcks7Padding([]byte(profileStr), 16)
	ct, err := cipherx.EncryptAesEcb(profile, a.key)
	if err != nil {
		panic(err)
	}
	return ct, nil
}

func (a Application) IsAdmin(encryptedProfile []byte) (bool, error) {
	pt, err := cipherx.DecryptAesEcb(encryptedProfile, a.key)
	if err != nil {
		panic(err)
	}
	profile := cipherx.RemovePcks7Padding(pt)
	vals, err := url.ParseQuery(string(profile))
	if err != nil {
		return false, fmt.Errorf("parsing profile: %w", err)
	}
	if !vals.Has("email") {
		return false, fmt.Errorf("profile is missing email")
	}
	if !vals.Has("uid") {
		return false, fmt.Errorf("profile is missing uid")
	}
	return vals.Get("role") == "admin", nil
}
