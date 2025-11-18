package main

import (
	"context"
	"cryptopals/internal/apps/mitm"
	"cryptopals/internal/dh"
	"fmt"
	"math/big"
	"sync"
)

type attackerChannels = mitm.AttackerChannels[Message]

func intercept(
	c attackerChannels,
	makeG func(p *big.Int) *big.Int,
) (dh.DiffeHellman, EncryptedMessage, EncryptedMessage) {
	var p *big.Int
	var g *big.Int
	var encA, encB EncryptedMessage
	c.AttackerLoop(context.Background(),
		func(msg Message) Message {
			switch msg := msg.(type) {
			case NegotiateGroup:
				p = msg.P
				g = makeG(msg.P)
				return NegotiateGroup{G: g, P: msg.P}
			case SharePublicKey:
				return msg
			case EncryptedMessage:
				encA = msg
				return msg
			default:
				panic("unknown message type")
			}
		}, func(msg Message) Message {
			switch msg := msg.(type) {
			case NegotiateGroup:
				p = msg.P
				g = makeG(msg.P)
				return NegotiateGroup{G: g, P: msg.P}
			case SharePublicKey:
				return msg
			case EncryptedMessage:
				encB = msg
				return msg
			default:
				panic("unknown message type")
			}
		})
	return dh.MakeDiffeHellman(p, g), encA, encB
}

func attack1(c attackerChannels) ([][]byte, [][]byte) {
	_, encA, encB := intercept(c, func(p *big.Int) *big.Int { return big.NewInt(1) })
	sharedKey := big.NewInt(1)
	plaintextA := DecryptMessage(sharedKey, encA)
	plaintextB := DecryptMessage(sharedKey, encB)
	return [][]byte{plaintextA}, [][]byte{plaintextB}
}

func attack2(c attackerChannels) ([][]byte, [][]byte) {
	_, encA, encB := intercept(c, func(p *big.Int) *big.Int { return p })
	sharedKey := big.NewInt(0)
	plaintextA := DecryptMessage(sharedKey, encA)
	plaintextB := DecryptMessage(sharedKey, encB)
	return [][]byte{plaintextA}, [][]byte{plaintextB}
}

func attack3(c attackerChannels) ([][]byte, [][]byte) {
	minus1 := func(p *big.Int) *big.Int {
		return new(big.Int).Sub(p, big.NewInt(1)) // g = p-1
	}

	group, encA, encB := intercept(c, minus1)
	sharedKeyCandidates := []*big.Int{big.NewInt(1), minus1(group.P())}
	var plaintextACandidates, plaintextBCandidates [][]byte
	for _, sharedKey := range sharedKeyCandidates {
		plaintextACandidates = append(plaintextACandidates, DecryptMessage(sharedKey, encA))
		plaintextBCandidates = append(plaintextBCandidates, DecryptMessage(sharedKey, encB))
	}
	return plaintextACandidates, plaintextBCandidates
}

func demoAttack(attack func(c attackerChannels) ([][]byte, [][]byte)) {
	sim := MakeSimulation()

	messageToSend := []byte("hello world")

	var messageReturned []byte
	var plaintextACandidates, plaintextBCandidates [][]byte

	var wg sync.WaitGroup

	wg.Go(func() { messageReturned = sim.RunClientA(messageToSend) })
	wg.Go(sim.RunClientB)
	wg.Go(func() { plaintextACandidates, plaintextBCandidates = attack(sim.AttackerChannels()) })
	wg.Wait()

	fmt.Println("--- Final results ---")
	fmt.Printf("Message to send was '%s'\n", messageToSend)
	fmt.Printf("Client A got '%s' back from client B\n", messageReturned)
	fmt.Printf("Attacker got from client A:\n")
	for _, plaintext := range plaintextACandidates {
		fmt.Printf("%s\n", plaintext)
	}
	fmt.Printf("Attacker got from client B:\n")
	for _, plaintext := range plaintextBCandidates {
		fmt.Printf("%s\n", plaintext)
	}
}

func main() {
	demoAttack(attack1)
	demoAttack(attack2)
	demoAttack(attack3)
}
