# Challenge 29

**Break a SHA-1 keyed MAC using length extension**

## Challenge description

> Secret-prefix SHA-1 MACs are trivially breakable.
> 
> The attack on secret-prefix SHA1 relies on the fact that you can take the ouput of SHA-1 and use it as a new starting point for SHA-1, thus taking an arbitrary SHA-1 hash and "feeding it more data".
> 
> Since the key precedes the data in secret-prefix, any additional data you feed the SHA-1 hash in this fashion will appear to have been hashed with the secret key.
> 
> To carry out the attack, you'll need to account for the fact that SHA-1 is "padded" with the bit-length of the message; your forged message will need to include that padding. We call this "glue padding". The final message you actually forge will be:
> 
> ```
> SHA1(key || original-message || glue-padding || new-message)
> ```
>
> (where the final padding on the whole constructed message is implied)
> 
> Note that to generate the glue padding, you'll need to know the original bit length of the message; the message itself is known to the attacker, but the secret key isn't, so you'll need to guess at it.
> 
> This sounds more complicated than it is in practice.
> 
> To implement the attack, first write the function that computes the MD padding of an arbitrary message and verify that you're generating the same padding that your SHA-1 implementation is using. This should take you 5-10 minutes.
> 
> Now, take the SHA-1 secret-prefix MAC of the message you want to forge --- this is just a SHA-1 hash --- and break it into 32 bit SHA-1 registers (SHA-1 calls them "a", "b", "c", &c).
> 
> Modify your SHA-1 implementation so that callers can pass in new values for "a", "b", "c" &c (they normally start at magic numbers). With the registers "fixated", hash the additional data you want to forge.
> 
> Using this attack, generate a secret-prefix MAC under a secret key (choose a random word from /usr/share/dict/words or something) of the string:
> 
> ```
> "comment1=cooking%20MCs;userdata=foo;comment2=%20like%20a%20pound%20of%20bacon"
> ```
>
> Forge a variant of this message that ends with ";admin=true".

## Solution

### Utility functions

I first extended my `sha1x` package with a few utility functions.

```go
type HashUtils interface {
	// FromRegisters reconstructs the hash state from the registers, assuming that the previously summed data was already padded.
	// h is the list of registers, and have exactly the same number of registers as the respective hash algorithm.
	// len is the length of the previously summed data and must be a multiple of 64.
	FromRegisters(h []uint32, len uint64) hash.Hash

	// DigestToRegisters splits the digest into 32-bit registers.
	// digest must be the correct length for the respective hash algorithm.
	DigestToRegisters(digest []byte) []uint32

	// Generates MD-padding for a message of length len for the respective hash algorithm.
	Padding(len uint64) []byte
}
```

Initially I just created these as standalone functions. However, the next challenge had me carry out the exact same attack for a different hash algorithm, which required these exact same utility functions. Creating an interface made it easy to swap between the different hash algorithms in the attack.

`FromRegisters` just constructs a `digest` object with the registers set to the provided values. It also sets the `len` value on the `digest` object, since the current message length will need to be known when `Sum` is called and the hash algorithm computes the padding.

`DigestToRegisters` simply splits the digest into registers, as described in the previous challenge.

`Padding` simply generates the MD-padding for a certain (unpadded) message length, as described in the previous challenge.

### Length extension attack

The challenge description explains the attack pretty well, so I will just add a bit of my own formalization and detail.

#### Function definitions

##### SHA-1

Let $B(d, b)$ be the SHA-1 block function that takes the current digest/registers $d$, updates it with block $b$, and outputs the new digest.

Let $S(d, m)$ be the summing function that takes the current digest $d$, updates it with the blocks $m = b_1 \mid\mid \ldots \mid\mid b_n$ ($m$ must be block-aligned), and outputs the new digest. We can define $S$ recursively as

```math
S(d, b_1 \mid\mid \ldots \mid\mid b_n) =
\begin{cases}
d & n=0 \\
S(B(d, b_1), b_2 \mid\mid \ldots \mid\mid b_n)  & n>0
\end{cases} 
```

One could easily prove the following property from this definition:

> **Lemma 1:**
> Let $d$ be some digest and $m_1$ and $m_2$ be two  block-aligned messages. Then 
> ```math
> S(d, m_1 \mid\mid m_2) = S(S(d, m_1), m_2)
> ```

This is the main property that the length extension attack uses.

Let $\mathfrak{v}$ be initial "magic" value of the digest that SHA-1 initializes to. Then we can define the full SHA-1 hash function $H$ as

```math
H(m) := S(\mathfrak{v}, m \mid\mid \mathrm{padding}(|m|))
```

where $\mathrm{padding}(|m|)$ is the MD padding for a message of length $|m|$.

##### MAC

We can define the SHA-1 MAC function as simply

```math
\mathrm{MAC}(m, k) := H(k \mid\mid m)
```

#### The attack

Let $m$ be the original message return by the application and let $c := \mathrm{MAC}(m, k)$ be the MAC computed for this message by the application. We can see that

```math
c = H(k \mid\mid m) = S(\mathfrak{v}, k \mid\mid m \mid\mid \mathrm{padding}(|k| + |m|))
```

As the attacker, we are provided $m$ and $c$. We are not directly provided the length of the key, $|k|$, but we can guess it through brute-force.

Let $g := \mathrm{padding}(|k| + |m|)$ which we have enough information to compute. The challenge calls this the "glue padding."

Let $e$ be the new data we want to extend the message with. We can update the current digest $c$ with this new data by calling $c' := S(c, e \mid\mid p)$, as long as we pad it with some $p$ to block-align it. We can see that

```math
\begin{align*}
c' &= S(c, e \mid\mid p) \\
&= S(S(\mathfrak{v}, k \mid\mid m \mid\mid g), \; e \mid\mid p) \\
&= S(\mathfrak{v}, k \mid\mid m \mid\mid g \mid\mid e \mid\mid p) & \text{Lemma 1}
\end{align*}
```

This looks very close to a MAC for the message $`m' := m \mid\mid g \mid\mid e`$. We just need to make sure $p$ is valid MD padding for $k \mid\mid m'$, so let $p = \text{padding}(|k| + |m'|)$. Then we finally get that

```math
\begin{align*}
c' &= S(\mathfrak{v}, k \mid\mid m' \mid\mid \text{padding}(|k| + |m'|)) \\
&= H(k \mid\mid m') \\
&= \text{MAC}(m', k) \\
\end{align*}
```

by the definition of $H$ and $\text{MAC}$

We summarize this final result as follows:


> Let $m$ be a message, $k$ be a MAC key, $c := \mathrm{MAC}(m, k)$, $e$ be the new data we want to extend the message with.
> 
> Let $g := \mathrm{padding}(|k| + |m|)$ and $`m' := m \mid\mid g \mid\mid e`$. Then
> 
> ```math
> \text{MAC}(m', k) = S(c, e \mid\mid \text{padding}(|k| + |m'|))
> ```

In other words, given a message $m$, its MAC, and the MAC key length, we fabricate a valid MAC for the message $m \mid\mid g \mid\mid e$.

> An aside: I like doing this formalization because it allows us to see concretely how all the pieces fall into place. It also helps us to convince ourselves that the properties of the attack are true, and it gives us a resource to look back to if we ever need to re-convince ourselves of these properties.

### Attack implementation

The attack implementation is quite succinct, so I include it here as well:

```go
func ExtendMac(hu hashx.HashUtils, mac []byte, keySize int, originalMessage []byte, newMessage []byte) ([]byte, []byte) {
	macStateLen := uint64(keySize + len(originalMessage))
	gluePadding := hu.Padding(macStateLen)

	h := hu.FromRegisters(hu.DigestToRegisters(mac), macStateLen+uint64(len(gluePadding)))
	h.Write(newMessage)
	extendedMac := h.Sum(nil) // nil here just means we append the digest to an empty slice

	fullMessage := slices.Concat(originalMessage, gluePadding, newMessage)

	return fullMessage, extendedMac
}
```

### Completing the challenge

Finally, the actual challenge provides us a SHA-1 MAC for the original message

```
comment1=cooking%20MCs;userdata=foo;comment2=%20like%20a%20pound%20of%20bacon
```

and we need to construct a new message ending with `";admin=true"` and its respective MAC, without knowing the secret key used to create the original MAC.

We perform the previously described length extension attack to extend the original MAC with the new message `";admin=true"`. This gives us a MAC for the message:

```
comment1=cooking%20MCs;userdata=foo;comment2=%20like%20a%20pound%20of%20bacon[glue padding];admin=true
```

> **An aside:** Notice that the glue padding is part of the message. In this case, it becomes part of the value for `comment2`. This could potentially make the message invalid to the application. For instance, if the length of `key \mid\mid originalMessage` happened to be 61, the padding would contain a `'='` character at the end, causing us to have a `=` without a matching `;`. In this challenge, we get lucky that the glue padding does not make the message invalid.

We then pass both the new message (including the glue padding in the middle) and the forged MAC to the application to gain administrator access.

## Commentary

I think this attack is an interesting example where a specific use of a cryptographic algorithm seems secure with an understanding of the basic properties of that algorithm, but it falls apart once you know more details about the algorithm. Namely, it seems like the 
"avalanche effect" and "non-invertible" properties of SHA-1 would make it so the digest $c := \text{SHA-1}(k \mid\mid m)$ can only be computed if you know $k$, and that if you are given just $m$ and $c$ there is no tractable way to reverse information about $k$ or produce your own hash $c' := \text{SHA-1}(k \mid\mid m')$. However, once you know that when SHA-1 consumes each message block, it computes the new digest from the previous digest, it becomes very clear that you can just use $c$ as the initial digest and keep hashing more message data to get $c'$.
