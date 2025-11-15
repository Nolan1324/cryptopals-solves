package main

import (
	"cryptopals/internal/apps/mitm"
	"cryptopals/internal/cipherx"
	"cryptopals/internal/dh"
	"cryptopals/internal/randx"
	"log"
	"math/big"
)

// Message is an interface for a message that can between sent between clients A and B
type Message interface {
	// message is defined privately so that no external types can implement this interface
	message()
}

// KeyExchangeRequest is a request to start a key exchange,
// containing the Diffie-Hellman parameters and the requester's public key.
type KeyExchangeRequest struct {
	P         *big.Int
	G         *big.Int
	PublicKey dh.PublicKey
}

// KeyExchangeResponse is a response to a key exchange,
// containing the responder's public key.
type KeyExchangeResponse struct {
	PublicKey dh.PublicKey
}

// EncryptedMessage is messaged encrypted with AES-CBC using the,
// previously estabilished shared key from the key exchange
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

// MakeSimulation creates a new man-in-the-middle simulation for challenge 34
func MakeSimulation() Simulation {
	return Simulation{sim: mitm.MakeSimulation[Message](0)}
}

// AttackerChannels are the channels that the man-in-the-middle uses to intercept communication between the two clients
func (s Simulation) AttackerChannels() mitm.AttackerChannels[Message] {
	return s.sim.AttackerChannels()
}

// RunClientA runs the logic for client A talking to client B
func (s Simulation) RunClientA(messageToSend []byte) (messageReceived []byte) {
	send, recv := s.sim.ClientAChannels()
	defer close(send)

	sharedKey := runClientAKeyExchange(send, recv)
	return runClientAEncryptedConversation(send, recv, sharedKey, messageToSend)
}

// RunClientA runs the logic for client B talking to client A
func (s Simulation) RunClientB() {
	send, recv := s.sim.ClientBChannels()
	defer close(send)

	sharedKey := runClientBKeyExchange(send, recv)
	runClientBEncryptedConversation(send, recv, sharedKey)
}

func runClientAKeyExchange(send chan<- Message, recv <-chan Message) dh.SharedKey {
	log.Println("Client A sending key exchange request")
	dhClient := dh.MakeClientWithRandomKey(dh.MakeNistDiffeHellman())
	send <- KeyExchangeRequest{P: dhClient.P(), G: dhClient.G(), PublicKey: dhClient.PublicKey()}
	log.Println("Client A reading key exchange response")
	response, ok := (<-recv).(KeyExchangeResponse)
	if !ok {
		panic("client B set invalid response type")
	}
	return dhClient.SharedKey(response.PublicKey)
}

func runClientAEncryptedConversation(send chan<- Message, recv <-chan Message, sharedKey dh.SharedKey, messageToSend []byte) []byte {
	log.Println("Client A sending encrypted message")
	send <- EncryptMessage(sharedKey, messageToSend)
	log.Println("Client A reading encrypted response")
	encMsg, ok := (<-recv).(EncryptedMessage)
	if !ok {
		panic("client B sent invalid encrypted message type")
	}
	msg := DecryptMessage(sharedKey, encMsg)
	// Client B should have sent back the same message
	log.Printf("Client A got message: %s\n", msg)
	return messageToSend
}

func runClientBKeyExchange(send chan<- Message, recv <-chan Message) dh.SharedKey {
	log.Println("Client B reading key exchange request")
	request, ok := (<-recv).(KeyExchangeRequest)
	if !ok {
		panic("client A sent invalid request type")
	}
	dhClient := dh.MakeClientWithRandomKey(dh.MakeDiffeHellman(request.P, request.G))
	sharedKey := dhClient.SharedKey(request.PublicKey)
	log.Println("Client B sending key exchange response")
	send <- KeyExchangeResponse{PublicKey: dhClient.PublicKey()}
	return sharedKey
}

func runClientBEncryptedConversation(send chan<- Message, recv <-chan Message, sharedKey dh.SharedKey) {
	log.Println("Client B reading encrypted message")
	encMsg, ok := (<-recv).(EncryptedMessage)
	if !ok {
		panic("client A sent invalid encrypted message type")
	}
	msg := DecryptMessage(sharedKey, encMsg)
	// Send the same message back to client A
	log.Println("Client B sending encrypted response")
	send <- EncryptMessage(sharedKey, msg)
}

// EncryptMessage creates an encrypted message using the shared key
func EncryptMessage(sharedKey dh.SharedKey, plaintext []byte) EncryptedMessage {
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

// DecryptMessage decrypts a message using the shared key
func DecryptMessage(sharedKey dh.SharedKey, message EncryptedMessage) []byte {
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
