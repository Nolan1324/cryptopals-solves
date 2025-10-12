package mitm

import "context"

type Simulation[T any] struct {
	outgoingA chan T
	incomingA chan T
	outgoingB chan T
	incomingB chan T
}

type AttackerChannels[T any] struct {
	OutgoingA <-chan T
	IncomingA chan<- T
	OutgoingB <-chan T
	IncomingB chan<- T
}

// MakeSimulation makes a man-in-the-middle simulation using Go channels.
//
// If cap is set to 0, the channels block until the other client reads it, which could cause deadlock.
// To avoid deadlock, either:
//   - Ensure the clients order their sends / reads in a way that prevents wait-for cycles.
//   - Set the capacity as need (for instance, cap=1, if each client only sends one message at a time).
//   - Have clients send messages in another goroutine to model async communication.
func MakeSimulation[T any](cap int) Simulation[T] {
	return Simulation[T]{outgoingA: make(chan T, cap), incomingA: make(chan T, cap), outgoingB: make(chan T, cap), incomingB: make(chan T, cap)}
}

// ClientAChannels are the channels that client A communicates with client B through
func (s Simulation[T]) ClientAChannels() (chan<- T, <-chan T) {
	return s.outgoingA, s.incomingA
}

// ClientBChannels are the channels that client B communicates with client A through
func (s Simulation[T]) ClientBChannels() (chan<- T, <-chan T) {
	return s.outgoingB, s.incomingB
}

// AttackerChannels are the channels that the man-in-the-middle uses to intercept communication between the two clients
func (s Simulation[T]) AttackerChannels() AttackerChannels[T] {
	return AttackerChannels[T]{OutgoingA: s.outgoingA, IncomingA: s.incomingA, OutgoingB: s.outgoingB, IncomingB: s.incomingB}
}

// AttackerLoop runs a loop for the attacker to intercept and modify messages between A and B and B and A.
// The handler functions are never called concurrently.
// The loop completes when both clients A and B have closed their outgoing channels,
// or if the provided context is done.
func (c AttackerChannels[T]) AttackerLoop(ctx context.Context, handleAToB, handleBToA func(T) T) {
	defer close(c.IncomingA)
	defer close(c.IncomingB)

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.OutgoingA:
			if !ok {
				c.OutgoingA = nil
				break
			}
			c.IncomingB <- handleAToB(msg)
		case msg, ok := <-c.OutgoingB:
			if !ok {
				c.OutgoingB = nil
				break
			}
			c.IncomingA <- handleBToA(msg)
		}

		if c.OutgoingA == nil && c.OutgoingB == nil {
			break
		}
	}
}

// Passthrough runs the man-in-the-middle to simply passthrough messages between the clients without altering them.
func (c AttackerChannels[T]) Passthrough(ctx context.Context) {
	identity := func(msg T) T { return msg }
	c.AttackerLoop(ctx, identity, identity)
}
