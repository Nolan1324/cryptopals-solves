package main

import (
	"cryptopals/internal/apps/mitm"
	"cryptopals/internal/cipherx"
	"cryptopals/internal/dh"
	"cryptopals/internal/randx"
	"log"
	"math/big"
)

type Message interface {
	message()
}

type KeyExchangeRequest struct {
	P         *big.Int
	G         *big.Int
	PublicKey *big.Int
}

type KeyExchangeResponse struct {
	PublicKey *big.Int
}

type EncryptedMessage struct {
	Ciphertext []byte
	Iv         []byte
}

func (KeyExchangeRequest) message()  {}
func (KeyExchangeResponse) message() {}
func (EncryptedMessage) message()    {}

type Simulation struct {
	sim mitm.Simulation[Message]
}

func MakeSimulation() Simulation {
	return Simulation{sim: mitm.MakeSimulation[Message](0)}
}

func (s Simulation) RunClientA(messageToSend []byte) (messageReceived []byte) {
	send, recv := s.sim.ClientAChannels()
	defer close(send)

	sharedKey := runClientAKeyExchange(send, recv)
	log.Println("Client A sending encrypted message")
	send <- encryptMessage(sharedKey, messageToSend)
	log.Println("Client A reading encrypted response")
	encMsg, ok := (<-recv).(EncryptedMessage)
	if !ok {
		panic("client B sent invalid encrypted message type")
	}
	msg := decryptMessage(sharedKey, encMsg)
	// Client B should have sent back the same message
	log.Printf("Client A got message: %s\n", msg)
	return msg
}

func (s Simulation) RunClientB() {
	send, recv := s.sim.ClientBChannels()
	defer close(send)

	sharedKey := runClientBKeyExchange(send, recv)
	log.Println("Client B reading encrypted message")
	encMsg, ok := (<-recv).(EncryptedMessage)
	if !ok {
		panic("client B sent invalid encrypted message type")
	}
	msg := decryptMessage(sharedKey, encMsg)
	// Send the same message back to client A
	log.Println("Client B sending encrypted response")
	send <- encryptMessage(sharedKey, msg)
}

func (s Simulation) AttackerChannels() mitm.AttackerChannels[Message] {
	return s.sim.AttackerChannels()
}

func runClientAKeyExchange(send chan<- Message, recv <-chan Message) dh.SharedKey {
	log.Println("Client A sending key exchange request")
	dhClient := dh.MakeClientWithRandomKey(dh.MakeNistDiffeHellman())
	send <- KeyExchangeRequest{
		P:         dhClient.P(),
		G:         dhClient.G(),
		PublicKey: dhClient.PublicKey(),
	}
	log.Println("Client A reading key exchange response")
	response, ok := (<-recv).(KeyExchangeResponse)
	if !ok {
		panic("Client B set invalid response type")
	}
	return dhClient.SharedKey(response.PublicKey)
}

func runClientBKeyExchange(send chan<- Message, recv <-chan Message) dh.SharedKey {
	log.Println("Client B reading key exchange request")
	request, ok := (<-recv).(KeyExchangeRequest)
	if !ok {
		panic("Other party set invalid request type")
	}
	dhClient := dh.MakeClientWithRandomKey(dh.MakeDiffeHellman(request.P, request.G))
	sharedKey := dhClient.SharedKey(request.PublicKey)
	log.Println("Client B sending key exchange response")
	send <- KeyExchangeResponse{PublicKey: dhClient.PublicKey()}
	return sharedKey
}

func encryptMessage(sharedKey dh.SharedKey, plaintext []byte) EncryptedMessage {
	iv := randx.RandBytes(16)
	key := dh.ToAesKey(sharedKey)
	ciphertext, err := cipherx.EncryptAesCbc(cipherx.AddPkcs7Padding(plaintext, 16), key, iv)
	if err != nil {
		panic("failed to encrypt")
	}
	return EncryptedMessage{
		Ciphertext: ciphertext,
		Iv:         iv,
	}
}

func decryptMessage(sharedKey dh.SharedKey, message EncryptedMessage) []byte {
	key := dh.ToAesKey(sharedKey)
	plaintext, err := cipherx.DecryptAesCbc(message.Ciphertext, key, message.Iv)
	if err != nil {
		panic("failed to decrypt")
	}
	plaintext, err = cipherx.RemovePkcs7Padding(plaintext)
	if err != nil {
		panic("failed to remove PKCS7 padding")
	}
	return plaintext
}
