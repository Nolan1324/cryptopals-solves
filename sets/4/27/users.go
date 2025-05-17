package main

import (
	"bytes"
	"cryptopals/internal/cipherx"
	"cryptopals/internal/randx"
	"fmt"
	"log"
)

type Users struct {
	key []byte
	sim Simulation
}

type AsciiError struct {
	Message []byte
}

func (e *AsciiError) Error() string {
	return fmt.Sprintf("Message %q contains high value ASCII character", e.Message)
}

func MakeUsers(sim Simulation) Users {
	return Users{key: randx.RandBytes(16), sim: sim}
}

func (u Users) RunSender() {
	sendMessage(u.sim.GetSenderChan(), []byte("Hello world, how are you doing in this moment???"), u.key)
}

func (u Users) RunReciever() {
	recvChan, errChan := u.sim.GetReceiverChans()
	recvMessages(recvChan, errChan, u.key)
}

func (u Users) KeyMatches(key []byte) bool {
	return bytes.Equal(u.key, key)
}

func sendMessage(senderChan chan<- []byte, message []byte, key []byte) {
	if len(message)%16 != 0 {
		panic("message must be block aligned")
	}
	ct, err := cipherx.EncryptAesCbc(message, key, key)
	if err != nil {
		panic(err)
	}
	senderChan <- ct
	close(senderChan)
}

func recvMessages(recvChan <-chan []byte, errChan chan<- error, key []byte) {
	sendErr := func(err error) {
		go func() {
			errChan <- err
		}()
	}

	handleMessage := func(ct []byte) error {
		pt, err := cipherx.DecryptAesCbc(ct, key, key)
		if err != nil {
			return err
		}
		if !isValidAscii(pt) {
			return &AsciiError{Message: pt}
		}
		log.Printf("Received: %q\n", pt)
		return nil
	}

	for ct := range recvChan {
		sendErr(handleMessage(ct))
	}

	close(errChan)
}

func isValidAscii(data []byte) bool {
	for _, c := range data {
		if c >= 128 {
			return false
		}
	}
	return true
}
