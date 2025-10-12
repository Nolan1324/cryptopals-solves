package mitm

import (
	"context"
	"sync"
	"testing"
)

func TestPassthrough(t *testing.T) {
	sim := MakeSimulation[int](1)

	msg1 := 5
	msg2 := 7

	client := func(clientName string, msgToSend, expectedMsg int, send chan<- int, recv <-chan int) {
		defer close(send)
		t.Logf("Client %v started", clientName)
		send <- msgToSend
		t.Logf("Client %v sending message", clientName)
		msgRecieved := <-recv
		t.Logf("Client %v receiving message", clientName)
		if msgRecieved != expectedMsg {
			t.Errorf("Client %v expected message %v, got %v", clientName, expectedMsg, msgRecieved)
		}
	}

	var wg sync.WaitGroup
	wg.Go(func() {
		send, recv := sim.ClientAChannels()
		client("A", msg1, msg2, send, recv)
	})
	wg.Go(func() {
		send, recv := sim.ClientBChannels()
		client("B", msg2, msg1, send, recv)
	})
	wg.Go(func() {
		sim.AttackerChannels().Passthrough(context.Background())
	})
	wg.Wait()
}
