# Challenge 18

**Implement CTR, the stream cipher mode**

## Challenge description

> The string:
> 
> L77na/nrFsKvynd6HzOoG7GHTLXsTVu9qvY/2syLXzhPweyyMTJULu/6/kXX0KSvoOLSFQ==
> ... decrypts to something approximating English in CTR mode, which is an AES block cipher mode that turns AES into a stream cipher, with the > following parameters:
> 
>       key=YELLOW SUBMARINE
>       nonce=0
>       format=64 bit unsigned little endian nonce,
>              64 bit little endian block count (byte count / 16)
> CTR mode is very simple.
> 
> Instead of encrypting the plaintext, CTR mode encrypts a running counter, producing a 16 byte block of keystream, which is XOR'd against the > plaintext.
> 
> For instance, for the first 16 bytes of a message with these parameters:
> 
> keystream = AES("YELLOW SUBMARINE",
>                 "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00")
> ... for the next 16 bytes:
> 
> keystream = AES("YELLOW SUBMARINE",
>                 "\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00")
> ... and then:
> 
> keystream = AES("YELLOW SUBMARINE",
>                 "\x00\x00\x00\x00\x00\x00\x00\x00\x02\x00\x00\x00\x00\x00\x00\x00")
> CTR mode does not require padding; when you run out of plaintext, you just stop XOR'ing keystream and stop generating keystream.
> 
> Decryption is identical to encryption. Generate the same keystream, XOR, and recover the plaintext.
> 
> Decrypt the string at the top of this function, then use your CTR function to encrypt and decrypt other things.

## Solution

### Theory

CTR mode takes a pretty different approach than ECB and CBC mode. Rather than using the AES key to encrypt the plaintext, it used it to generate an infinite length "keystream" by encrypting a sequence of increasing numbers. We then just XOR the keystream with the plaintext to preform encryption/decryption. In previous challenges, we saw that encrypting by XORing with a _repeated_ key is very insecure. However, since the keystream is infinite, we do not need to repeat it at all! This still is not enough to ensure security though; as we will see in later challenges, if we use the same keystream to encrypt multiple plaintexts, that is just as bad as repeating the key when encrypting one plaintext. We thus also introduce a **nonce** value that is used together with AES key to generate the keystream. **We must use a different nonce for each plaintext we encrypt to ensure that the keystreams are different.** Like the IV in CBC mode, the nonce used to encrypt a ciphertext can be stored publicly, since you still need both the nonce **and** the AES key to generate the keystream and decrypt the ciphertext.

In general, a stream cipher provides a way to turn a finite length key into an infinite length keystream. The finite length key in CTR is the 16-byte AES key combined with the 64-bit nonce. The infinite length keystream is used for XOR encryption/decryption. As mentioned earlier, we keep the AES key the same for every encryption but change the nonce every time. Thus, we keep the AES key private but publicly store the nonce with the ciphertext.

### Implementation

#### Converting the counter to bytes

To convert the counter to bytes, I use `binary.LittleEndian.PutUint64`. The challenge description does not quite explain this, but the 8 counter bytes are to be placed in the **latter** half of the 16-byte block. Thus, the first half of the block is always just all 0's.

#### Encrypting the keystream bytes

I regenerate the keystream every time `Encrypt`/`Decrypt` is called.

#### Encrypt vs decrypt

Encrypt and decrypt work exactly the same in CTR, thanks to the self-inverting property of XOR. Thus, my `Decrypt` function just calls `Encrypt`.

#### Interface

This time I decided to implement the cipher as a datatype exposing encrypt/decrypt functions so that the internal block cipher object could be reused. A user can create a CTR cipher with `cipher, err := cipherx.NewAesCtr(key)` and then call the `Encrypt` and `Decrypt` functions on it.

### Commentary

#### CTR vs PRNG

I find it interesting that the CTR cipher does not use the decrypt function of AES at all. Instead, it essentially uses the AES encrypt function to act as a psuedo-random number generator (PRNG), with the AES key combined with the nonce essentially acting as the seed. Even more interesting, the confidentiality property of AES means that even if someone could obtain part of the keystream, they could not easily reverse engineer the key nor nonce, nor could they predict the later values in the keystream. These properties are unlike traditional PRNGs; for instance, we will see in the next few challenges that with the Mersenne Twister PRNG, if you know the first 624 numbers returned by the PRNG, you can predict all of the later numbers returned.
