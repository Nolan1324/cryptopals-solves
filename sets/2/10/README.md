# Challenge 10

**Implement CBC mode**

## Challenge Description

> CBC mode is a block cipher mode that allows us to encrypt > irregularly-sized messages, despite the fact that a block cipher > natively only transforms individual blocks.
> 
> In CBC mode, each ciphertext block is added to the next plaintext > block before the next call to the cipher core.
> 
> The first plaintext block, which has no associated previous > ciphertext block, is added to a "fake 0th ciphertext block" > called the initialization vector, or IV.
> 
> Implement CBC mode by hand by taking the ECB function you wrote > earlier, making it encrypt instead of decrypt (verify this by > decrypting whatever you encrypt to test), and using your XOR > function from the previous exercise to combine them.
> 
> The file here is intelligible (somewhat) when CBC decrypted > against "YELLOW SUBMARINE" with an IV of all ASCII 0 > (\x00\x00\x00 &c)

## Solution

This challenge is fairly straightfoward to implement. I introduce a bit of notation here for CBC for later challenges write ups.

Let $m$ be the current plaintext block, $c$ be the current ciphertext block, and $v$ be the previous ciphertext block (or the IV). Then the encrypt/decrypt functions for this block are:

$E_\mathrm{CBC\_block}(m, v, k) := E_\mathrm{AES}(m \oplus v, k)$

$D_\mathrm{CBC\_block}(c, v, k) := D_\mathrm{AES}(c, k) \oplus v$