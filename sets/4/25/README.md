# Challenge 25

**Break "random access read/write" AES CTR**

## Challenge Description

> Back to CTR. Encrypt the recovered plaintext from this file (the ECB exercise) under CTR with a random key (for this exercise the key should be unknown to you, but hold on to it).
> 
> Now, write the code that allows you to "seek" into the ciphertext, decrypt, and re-encrypt with different plaintext. Expose this as a function, like, "edit(ciphertext, key, offset, newtext)".
> 
> Imagine the "edit" function was exposed to attackers by means of an API call that didn't reveal the key or the original plaintext; the attacker has the ciphertext and controls the offset and "new text".
> 
> Recover the original plaintext.

## The application

The application I implemented exposes the following functions

```go
// EditByte changes the encrypted byte at offset to newByte
func (a Application) EditByte(offset int, newByte byte)

// ReadCiphertext reads the entire current ciphertext
func (a Application) ReadCiphertext() []byte
```

This could be thought of encrypted random-access memory that provides write access but not (decrypted) read access.

## Solution

Similar to previous the CTR challenges, this attack exploits XOR keystream reuse.

Let $K$ be the keystream generated from CTR, $m$ be the original plaintext, and $c := m \oplus K$ be the original ciphertext.

First we call `ReadCiphertext` so that we can read out the original ciphertext $c$.

If we call `EditByte` on each byte to write a new message $m'$, we have replaced the ciphertext $c$ with $c' = m' \oplus K$. We can call `ReadCiphertext` to read out $c'$.

Then we can discover the keystream since $K = m' \oplus c'$ and $m'$ and $c'$ are both known to us. To further simplify this, we could just write the message $m' = 0$ so that $K = 0 \oplus c' = c'$.

Now that we have the keystream, we decrypt the original ciphertext with simply $m = c \oplus K$.
