# Challenge 24

**Create the MT19937 stream cipher and break it**

## Challenge description

> You can create a trivial stream cipher out of any PRNG; use it to generate a sequence of 8 bit outputs and call those outputs a keystream. XOR each byte of plaintext with each successive byte of keystream.
> 
> Write the function that does this for MT19937 using a 16-bit seed. Verify that you can encrypt and decrypt properly. This code should look similar to your CTR code.
> 
> Use your function to encrypt a known plaintext (say, 14 consecutive 'A' characters) prefixed by a random number of random characters.
> 
> From the ciphertext, recover the "key" (the 16 bit seed).
> 
> Use the same idea to generate a random "password reset token" using MT19937 seeded from the current time.
> 
> Write a function to check if any given password token is actually the product of an MT19937 PRNG seeded with the current time.

## The cipher

To generate the keystream, we construct an MT19937 with a 16-bit seed (with the seed acting as the "key") and then generate 8-bit outputs. It was unclear to me if we should truncate each 32-bit output to an 8-bit output, or if we should split each 32-bit output into four 8-bit outputs. I chose the former.

The keystream is then XOR'ed with the plaintext/ciphertext to encrypt/decrypt it, just like in CTR.

## Part 1

The challenge first asks us to "encrypt a known plaintext (say, 14 consecutive 'A' characters) prefixed by a random number of random characters" and then use the ciphertext to figure out the 16-bit seed.

So the "application" exposes a function $E(m) := (p \mid\mid m) \oplus K$, where $p$ is the random padding and $K$ is the keystream generated from the seed $k$.

Let $m = ('A')^{\times 14}$ as suggested and call $c := E(m)$. We can discover the padding length $|p|$ via $|p| = |c| - |m|$. This tells us that $c_{|p|:|p|+|m|} = m \oplus K_{|p|:|p|+|m|}$, and thus we can compute $K_{|p|:|p|+|m|} = m \oplus c_{|p|:|p|+|m|}$ since we know $c$ and $m$.

We know just need to figure out which seed generates the keystream bytes $K_{|p|:|p|+|m|}$. Since the seed is only 16-bits, there are only $2^{16} = 65536$ possibilities, so we can feasibly just try each seed. For each seed, construct a MT19937, skip the first $|p|$ outputs, then read the next $|m|$ outputs (each truncated to 8-bits) and check if they equal $K_{|p|:|p|+|m|}$. If they are equal, we have found the seed (with high probability).

## Part 2

I was a bit confused what this part was saying when I first read it, but I think the idea is to imagine we are generating a "password reset token" for some user by seeding an MT19937 with the current time and then pulling, say, sixteen 8-bit outputs.

If we want to check if a 16-byte token was generated using this method, we can just seed our own MT19937 with the current time, generate a token, and then compare them. However, the token probably was generated a little bit before the current time, so we could just try timestamp seeds for the past, say, 500 seconds and check if any of them generate the token. Note that if none of these generate the token, it is still possible that the token was generated from an MT19937 seeded sometime before 500 seconds ago.
