package main

import (
	"cryptopals/internal/cipherx"
	"errors"
	"fmt"
	"os"
	"sync"
)

// Runs an attacker that just passes the message unchanged
func runNoopAttacker(attackerChans AttackerChans) {
	ct := <-attackerChans.InputChan
	attackerChans.OutputChan <- ct
	err := <-attackerChans.ErrorChan
	if err != nil {
		panic(err)
	}
	close(attackerChans.OutputChan)
}

func runAttacker(attackerChans AttackerChans) []byte {
	ct := <-attackerChans.InputChan
	if len(ct) < 48 {
		panic("attacker needs ciphertext to be at least three blocks")
	}

	// Tamper with the ciphertext
	c1 := ct[0:16]
	c2 := ct[16:32]
	c3 := ct[32:48]
	for i := range c2 {
		c2[i] = 0
	}
	copy(c3, c1)

	attackerChans.OutputChan <- ct

	err := <-attackerChans.ErrorChan

	var key []byte

	var asciiErr *AsciiError
	if errors.As(err, &asciiErr) {
		fmt.Println("Attacker successfully caused ASCII error")
		pt := asciiErr.Message
		key = make([]byte, 16)
		cipherx.XorBytes(key, pt[0:16], pt[32:48])
	} else if err != nil {
		panic(err)
	} else {
		panic("attacker did not cause any error")
	}

	close(attackerChans.OutputChan)

	return key
}

func main() {
	sim := MakeSimulation()
	users := MakeUsers(sim)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		users.RunSender()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		users.RunReciever()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// runNoopAttacker(sim.GetAttackerChans())
		key := runAttacker(sim.GetAttackerChans())
		fmt.Printf("Attacker found key: %v\n", key)
		if users.KeyMatches(key) {
			fmt.Printf("Key found by attacker is correct!\n")
		} else {
			fmt.Printf("Key found by attacker is incorrect\n")
			os.Exit(1)
		}
	}()

	wg.Wait()
}
