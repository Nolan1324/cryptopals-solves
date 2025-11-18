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

type NegotiateGroup struct {
	P *big.Int
	G *big.Int
}

type SharePublicKey struct {
	PublicKey dh.PublicKey
}

// EncryptedMessage is messaged encrypted with AES-CBC using the,
// previously estabilished shared key from the key exchange
type EncryptedMessage struct {
	Ciphertext []byte
	Iv         []byte
}

func (NegotiateGroup) message()   {}
func (SharePublicKey) message()   {}
func (EncryptedMessage) message() {}

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
	log.Println("Client A sending group negotiation")
	proposedDh := dh.MakeNistDiffeHellman()
	send <- NegotiateGroup{P: proposedDh.P(), G: proposedDh.G()}
	log.Println("Client A reading group negotiation response")
	negotiateResponse, ok := (<-recv).(NegotiateGroup)
	if !ok {
		panic("client B set invalid response type")
	}
	dhClient := dh.MakeClientWithRandomKey(dh.MakeDiffeHellman(negotiateResponse.P, negotiateResponse.G))
	log.Println("Client A sharing public key")
	send <- SharePublicKey{dhClient.PublicKey()}
	log.Println("Client A reading client B's public key")
	shareKeyResponse, ok := (<-recv).(SharePublicKey)
	if !ok {
		panic("client B set invalid response type")
	}
	return dhClient.SharedKey(shareKeyResponse.PublicKey)
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
	log.Println("Client B reading group negotiation")
	negotiateGroup, ok := (<-recv).(NegotiateGroup)
	if !ok {
		panic("client B set invalid response type")
	}
	group := dh.MakeDiffeHellman(negotiateGroup.P, negotiateGroup.G)
	log.Println("Client B sending group negotiation response")
	send <- NegotiateGroup{P: group.P(), G: group.G()}
	log.Println("Client B reading client A's public key")
	shareKeyResponse, ok := (<-recv).(SharePublicKey)
	if !ok {
		panic("client A set invalid response type")
	}
	log.Println("Client B sharing public key")
	dhClient := dh.MakeClientWithRandomKey(group)
	send <- SharePublicKey{dhClient.PublicKey()}
	return dhClient.SharedKey(shareKeyResponse.PublicKey)
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
	plaintext, _ = cipherx.RemovePkcs7Padding(plaintext) // ignore error
	return plaintext
}
