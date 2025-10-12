package main

import (
	"context"
	"cryptopals/internal/apps/mitm"
	"fmt"
	"math/big"
	"sync"
)

// attack performs the man-in-the-middle attacker for challenge 34
// and returns the decrypted plaintext sent by clients A and B, respectively.
func attack(c mitm.AttackerChannels[Message]) ([]byte, []byte) {
	var p *big.Int
	var plaintextA, plaintextB []byte
	c.AttackerLoop(context.Background(),
		func(msg Message) Message {
			switch msg := msg.(type) {
			case KeyExchangeRequest:
				p = msg.P
				return KeyExchangeRequest{G: msg.G, P: msg.P, PublicKey: msg.P}
			case KeyExchangeResponse:
				panic("unexpected message from client A: KeyExchangeResponse")
			case EncryptedMessage:
				plaintextA = DecryptMessage(big.NewInt(0), msg)
				return msg
			default:
				panic("unknown message type")
			}
		}, func(msg Message) Message {
			switch msg := msg.(type) {
			case KeyExchangeRequest:
				panic("unexpected message from client B: KeyExchangeRequest")
			case KeyExchangeResponse:
				if p == nil {
					panic("client A never sent key exchange request before client B sent key exchange response")
				}
				return KeyExchangeResponse{PublicKey: p}
			case EncryptedMessage:
				plaintextB = DecryptMessage(big.NewInt(0), msg)
				return msg
			default:
				panic("unknown message type")
			}
		})
	return plaintextA, plaintextB
}

func main() {
	sim := MakeSimulation()

	messageToSend := []byte("hello world")

	var messageReturned, plaintextA, plaintextB []byte

	var wg sync.WaitGroup

	wg.Go(func() { messageReturned = sim.RunClientA(messageToSend) })
	wg.Go(sim.RunClientB)
	wg.Go(func() { plaintextA, plaintextB = attack(sim.AttackerChannels()) })
	wg.Wait()

	fmt.Println("--- Final results ---")
	fmt.Printf("Message to send was '%s'\n", messageToSend)
	fmt.Printf("Client A got '%s' back from client B\n", messageReturned)
	fmt.Printf("Attacker got '%s' from client A\n", plaintextA)
	fmt.Printf("Attacker got '%s' from client B\n", plaintextB)
}
