package main

import (
	"bytes"
	"context"
	"sync"
	"testing"
)

func TestSimDefault(t *testing.T) {
	sim := MakeSimulation()

	messageToSend := []byte("hello world")

	var messageReturned []byte

	var wg sync.WaitGroup
	wg.Go(func() { messageReturned = sim.RunClientA(messageToSend) })
	wg.Go(sim.RunClientB)
	wg.Go(func() { sim.AttackerChannels().Passthrough(context.Background()) })
	wg.Wait()

	if !bytes.Equal(messageToSend, messageReturned) {
		t.Errorf("Expected message %v returned, but got message %v", messageToSend, messageReturned)
	}
}
