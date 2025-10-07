# Challenge 28

**Implement a SHA-1 keyed MAC**

## Challenge description

> Find a SHA-1 implementation in the language you code in.
> 
> Write a function to authenticate a message under a secret key by using a secret-prefix MAC, which is simply:
> 
> ```
> SHA1(key || message)
> ```
> Verify that you cannot tamper with the message without breaking the MAC you've produced, and that you can't produce a new MAC without knowing the secret key.

## Solution

### Hash algorithms

SHA-1 is a hash algorithm. Hash functions deterministically compute a fixed size "digest" / "checksum" of their inputs. A checksum essentially acts as a "fingerprint" of a message, as the same message will always hash to the same checksum. It is possible for two different messages to hash to the same digest (called a "hash collision") but this becomes increasingly unlikely the larger the digest is. For instance, SHA-1 uses 20-byte / 160-bit digests, so the probability of two hashes colliding is $2^{-160}$ which is incredibly small (however, there are more advanced attacks to find SHA-1 hash collisions, and hence SHA-1 is considered insecure).

Two key properties of a hash function are "avalanche effect" and non-invertible. The idea behind the avalanche effect is that if we change even one bit of the input to the hash function, many bits of the digest will change in an unpredictable (but deterministic) way. This makes the hash algorithm, in practice, "non-invertible" meaning that given the output digest, it is computationally intractable (but not impossible) to find a corresponding input.

### The SHA-1 algorithm

SHA-1 is hash algorithm that produces a 20-byte / 160-bit digests. It operates on 64-byte / 512-bit chunks of its input. To compute the digest, SHA-1 initializes the digest to fixed "magic" numbers. When it is fed a block, it preforms a new block operation to compute a new digest from the old digest. When all the chunks have been fed to the algorithm, you can simply read out the current digest as the final digest. A key detail here is that **the new digest is computed from the previous digest**, which will be critical in the next challenge.

Internally, the 160-bit digest is represented as 5 32-bit "registers" which are used during computation. To produce the final digest, we simply represent each register in big-endian and then concatenate them into a byte sequence. Conversely, given the digest as a byte sequence, you could easily split it back into registers.

Since SHA-1 operates on 512-bit blocks, if the input bit count is not divisible by 512, SHA-1 will pad the final block with MD padding. It does this by first appending a single `1` bit and then appending `0` bits until the message is exactly 8 bytes shy of the end of the block. Then, it writes the (unpadded) message length (as a 64-bit integer in big-endian) into the final those 8 bytes, completing the final block.

### SHA-1 implementation

The Go standard library has a SHA-1 implementation, but we will need to be able to edit the implementation for later challenges. So, I copied some of the code (with the appropriate license) into this project. The heart of the SHA-1 algorithm (and the most tedious part to code) is the `block` function, which consumes a 64-byte block to update 5 SHA-1 registers. Go actually has optimized implementations of `block` for different CPU architectures, but I just use the pure Go one.

### SHA-1 MAC implementation

As the challenge description explains, SHA-1 MAC can be implemented as simply

```
SHA1(key \mid\mid message)
```

(which is **_incredibly_** insecure, as we will see in the next challenge)

I make my MAC data structure take in an object implementing the `hash.Hash` interface as input, so that we can easily swap out the hash algorithm in a later challenge. So my MAC implementation looks like this:

```go
type Mac interface {
	Sign(message []byte) []byte
}

type mac struct {
	key  []byte
	hash hash.Hash
}

func NewMac(key []byte, hash hash.Hash) Mac {
	return &mac{key: key, hash: hash}
}

func NewSha1Mac(key []byte) Mac {
	return NewMac(key, sha1x.New())
}

func (m mac) Sign(message []byte) []byte {
	m.hash.Reset()
	m.hash.Write(append(bytes.Clone(m.key), message...))
	return m.hash.Sum(nil)
}
```
