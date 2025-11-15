# Challenge 34

**Implement a MITM key-fixing attack on Diffie-Hellman with parameter injection**

## Challenge description

> Use the code you just worked out to build a protocol and an "echo" bot. You don't actually have to do the network part of this if you don't want; just simulate that. The protocol is:
> 
> A->B \
> Send "p", "g", "A" \
> B->A \
> Send "B"  \
> A->B \
> Send AES-CBC(SHA1(s)[0:16], iv=random(16), msg) + iv \
> B->A \
> Send AES-CBC(SHA1(s)[0:16], iv=random(16), A's msg) + iv \
> (In other words, derive an AES key from DH with SHA1, use it in both directions, and do CBC with random IVs appended or prepended to the message).
> 
> Now implement the following MITM attack:
> 
> A->M \
> Send "p", "g", "A" \
> M->B \
> Send "p", "g", "p" \
> B->M \
> Send "B" \
> M->A \
> Send "p" \
> A->M \
> Send AES-CBC(SHA1(s)[0:16], iv=random(16), msg) + iv \
> M->B \
> Relay that to B \
> B->M \
> Send AES-CBC(SHA1(s)[0:16], iv=random(16), A's msg) + iv \
> M->A \
> Relay that to A \
> M should be able to decrypt the messages. "A" and "B" in the protocol --- the public keys, over the wire --- have been swapped out with "p". Do the DH math on this quickly to see what that does to the predictability of the key.
> 
> Decrypt the messages from M's vantage point as they go by.
> 
> Note that you don't actually have to inject bogus parameters to make this attack work; you could just generate Ma, MA, Mb, and MB as valid DH parameters to do a generic MITM attack. But do the parameter injection attack; it's going to come up again.

## The attack

First I will explain why the attack works since it is fairly short. Then in the next section I will go on a spiel about a man-in-the-middle simulation framework I implemented for this challenge using Go channels.

This attack replaces the exchanged public keys with $A := p$ and $B := p$. Thus, clients A and B compute the new shared key as:

- $K_a = B^a \bmod{p} = p^a \bmod{p}$
- $K_b = A^b \bmod{p} = p^b \bmod{p}$

But, recall that $p \equiv 0 \pmod{p}$. Also recall that modular arithmetic respects exponentiation. Thus,

- $K_a = 0^a \bmod{p} = 0$
- $K_b = 0^b \bmod{p} = 0$

This gives that $K_a = K_b$, so a shared key was still established despite our intervention. However, we now know that the shared key must be $0$, allowing us to decrypt the messages in transit.

Note that we could have replaced $A$ and $B$ with any integer equivalent to $0$ under $\mod p$; $0, -p, p, 2p, 30p$, etc., would have all worked just fine. Perhaps we picked $p$ instead of $0$ to be a bit more sneaky.

## Man-in-the-middle simulation framework in Go with channels

In a previous challenge, I created a basic framework to model a man-in-the-middle attack using Go channels. In this challenge, I create a more general framework for this challenge and future challenges.

The goal of the simulation is to model communication between two clients A and B where a man-in-the-middle may be intercepting and tampering with the messages. We define the following channels to model this

```go
type Simulation[T any] struct {
	outgoingA chan T
	incomingA chan T
	outgoingB chan T
	incomingB chan T
}
```

`T` is the message type. The simulation is created with `MakeSimulation` by specifying `T` and the capacity of the channels

```go
func MakeSimulation[T any](cap int) Simulation[T]
```

Client A sends on `outgoingA` and reads on `incomingA`

```go
func (s Simulation[T]) ClientAChannels() (chan<- T, <-chan T) {
	return s.outgoingA, s.incomingA
}
```

Likewise, client B sends on `outgoingB` and `incomingB`.

```go
func (s Simulation[T]) ClientBChannels() (chan<- T, <-chan T) {
	return s.outgoingB, s.incomingB
}
```

If there were no man-in-the-middle, we would simply let `outgoingA == incomingB` and `outgoingB == incomingA`.

If there is a man-in-the-middle (called the "attacker") we provide them the following channels.

```go
type AttackerChannels[T any] struct {
	OutgoingA <-chan T
	IncomingA chan<- T
	OutgoingB <-chan T
	IncomingB chan<- T
}

func (s Simulation[T]) AttackerChannels() AttackerChannels[T] {
	return AttackerChannels[T]{OutgoingA: s.outgoingA, IncomingA: s.incomingA, OutgoingB: s.outgoingB, IncomingB: s.incomingB}
}
```

For type safety, we cast the channels to send-only or read-only as appropriate (this by no means provides actual security, it just prevents us from accidentally using the channels incorrectly).

The attacker code would typically read messages from `OutgoingA` modify them, and then send them to `IncomingA`. At the same time, it would read messages from `OutgoingB` to `IncomingB`. We can create a helper function to implement this common pattern:

```go
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
```

We use `select` here to read from whichever outgoing channel has a message ready, or handle cancellation.

Note that `AttackerLoop` does not handle the channels in parallel. If we wanted to handle them in parallel, we could use separate goroutines like so:

```go
func (c AttackerChannels[T]) AttackerLoopConcurrent(handleAToB, handleBToA func(T) T) {
	defer close(c.IncomingA)
	defer close(c.IncomingB)

	var wg sync.WaitGroup
	wg.Go(func() {
		for msg := range c.OutgoingA {
			c.IncomingB <- handleAToB(msg)
		}
	})
	wg.Go(func() {
		for msg := range c.OutgoingB {
			c.IncomingA <- handleAToB(msg)
		}
	})
	wg.Wait()
}
```

We can use `AttackerLoop` to easily implement an attacker that just lets messages pass through without modifying them, as so

```go
func (c AttackerChannels[T]) Passthrough(ctx context.Context) {
	identity := func(msg T) T { return msg }
	c.AttackerLoop(ctx, identity, identity)
}
```

### Deadlock

When creating and using the simulation, deadlock may occur if not done properly. Suppose we create it with `MakeSimulation[int](0)` so that the channels have capacity `0`, and suppose we implement the clients as follows:

```go
func clientA(sim Simulation[int]) {
	send, recv := sim.ClientAChannels()
	send <- 43
	x := <-recv
}

func clientA(sim Simulation[int]) {
	send, recv := sim.ClientBChannels()
	send <- 57
	x := <-recv
}
```

Since the channels have capacity 0, sending will block until the other side reads the value. Thus, the following interleaving will deadlock:

```
A: send <- 43 (waiting on A to read)
B: send <- 57 (waiting on B to read)
```

There are multiple ways to solve this. One is to order the sends and reads to prevent wait-for cycles, like so

```go
func clientA(sim Simulation[int]) {
	send, recv := sim.ClientAChannels()
	send <- 43
	x := <-recv
}

func clientA(sim Simulation[int]) {
	send, recv := sim.ClientBChannels()
	x := <-recv
	send <- 57
}
```
 
Another method is to increase the channel capacity so that sending does not block. In our example, since we only put at most one message on each channel, increasing the capacity to `1` would avoid deadlock. However, this solution might not be safe in more complex scenarios, and it would be bad practice to try solving more complex deadlocks by just bumping up the capacity arbitrarily.

Finally, we could also solve the deadlock by wrapping each send in a Goroutine.

```go
func clientA(sim Simulation[int]) {
	send, recv := sim.ClientAChannels()
	go func() {
		send <- 43
	}()
	x := <-recv
}

func clientA(sim Simulation[int]) {
	send, recv := sim.ClientBChannels()
	go func() {
		send <- 57
	}()
	x := <-recv
}
```

This effectively allows us to model asynchronous communication, rather than synchronous communication. However, in more complex scenarios, this "fire-and-forget" pattern could prove problematic, as Goroutines that were launched in the past may cause unexpected results in the future.

### Modelling this challenge

The man-in-the-middle simulation framework allows us to model this challenge.

#### The message type

We will define the `Message` type. There are 3 types of messages clients may exchange: a key exchange request, a key exchange response, and a message encrypted with the key. Since there is a limited set of possible messages, we could model `Message` as a **sum type** (aka tagged union, discriminated union) with the variant types being `KeyExchangeRequest`, `KeyExchangeResponse`, and `EncryptedMessage`.

Go does not directly support sum types (at least not to the degree that Rust or a functional language does), but we can loosely approximate them using an interface. We define the `Message` type as an `interface` with one package-private function `message()`. Then we define three `struct` types `KeyExchangeRequest`, `KeyExchangeResponse`, `EncryptedMessage` that each implement `message()` as a no-op; they only implement it so that the Go compiler considers them as implementations of the `Message` interface. We made `message()` package-private so that external packages cannot implement more variants of `Message`.

```go
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
```

Given some `msg Message`, Go lets you cast it as so

```go
request, ok := msg.(KeyExchangeRequest)
if !ok {
	// msg is not of type KeyExchangeRequest, so request is zero-initialized
	return
}
// use request
```

Go also lets you switch on its concrete type, safely casting it in each case

```go
switch msg := msg.(type) {
	case KeyExchangeRequest:
		// msg is now type KeyExchangeRequest
	case KeyExchangeResponse:
		// msg is now type KeyExchangeResponse
	case EncryptedMessage:
		// msg is now type EncryptedMessage
	default:
		panic("unreachable")
}
```

If you try to cast it to a type that does not implement `Message`, then you will get a compilation error.

#### The clients

Now that we have the simulation and message types set up, we can implement the clients fairly directly.

```go
func (s Simulation) RunClientA(messageToSend []byte) (messageReceived []byte) {
	send, recv := s.sim.ClientAChannels()
	defer close(send)

	sharedKey := runClientAKeyExchange(send, recv)
	return runClientAEncryptedConversation(send, recv, sharedKey, messageToSend)
}

func runClientAKeyExchange(send chan<- Message, recv <-chan Message) dh.SharedKey {
	dhClient := dh.MakeClientWithRandomKey(dh.MakeNistDiffeHellman())
	send <- KeyExchangeRequest{P: dhClient.P(), G: dhClient.G(), PublicKey: dhClient.PublicKey()}
	response, ok := (<-recv).(KeyExchangeResponse)
	if !ok {
		panic("client B set invalid response type")
	}
	return dhClient.SharedKey(response.PublicKey)
}

func runClientAEncryptedConversation(send chan<- Message, recv <-chan Message, sharedKey dh.SharedKey, messageToSend []byte) []byte {
	send <- EncryptMessage(sharedKey, messageToSend)
	encMsg, ok := (<-recv).(EncryptedMessage)
	if !ok {
		panic("client B sent invalid encrypted message type")
	}
	msg := DecryptMessage(sharedKey, encMsg)
	return messageToSend
}
```

```go
func (s Simulation) RunClientB() {
	send, recv := s.sim.ClientBChannels()
	defer close(send)

	sharedKey := runClientBKeyExchange(send, recv)
	runClientBEncryptedConversation(send, recv, sharedKey)
}

func runClientBKeyExchange(send chan<- Message, recv <-chan Message) dh.SharedKey {
	request, ok := (<-recv).(KeyExchangeRequest)
	if !ok {
		panic("client A sent invalid request type")
	}
	dhClient := dh.MakeClientWithRandomKey(dh.MakeDiffeHellman(request.P, request.G))
	sharedKey := dhClient.SharedKey(request.PublicKey)
	send <- KeyExchangeResponse{PublicKey: dhClient.PublicKey()}
	return sharedKey
}

func runClientBEncryptedConversation(send chan<- Message, recv <-chan Message, sharedKey dh.SharedKey) {
	encMsg, ok := (<-recv).(EncryptedMessage)
	if !ok {
		panic("client A sent invalid encrypted message type")
	}
	msg := DecryptMessage(sharedKey, encMsg)
	send <- EncryptMessage(sharedKey, msg)
}
```

#### The attacker

We can use the `AttackerLoop` helper function from earlier to easily implement the attacker.

```go
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
```

As explained in the previous section, when the attacker changes the public keys to $p$, the shared key becomes $0$, so the attacker can decrypt the messages with $0$ as the shared key.
