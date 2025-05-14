# Challenge 11

**An ECB/CBC detection oracle**

## Challenge Description

> Now that you have ECB and CBC working:
> 
> Write a function to generate a random AES key; that's just 16 random bytes.
> 
> Write a function that encrypts data under an unknown key --- that is, a function that generates a random key and encrypts under it.
> 
> The function should look like:
> 
> ```
> encryption_oracle(your-input)
> => [MEANINGLESS JIBBER JABBER]
> ```
> 
> Under the hood, have the function append 5-10 bytes (count chosen randomly) before the plaintext and 5-10 bytes after the plaintext.
> 
> Now, have the function choose to encrypt under ECB 1/2 the time, and under CBC the other half (just use random IVs each time for CBC). Use rand(2) to decide which to use.
> 
> Detect the block cipher mode the function is using each time. You should end up with a piece of code that, pointed at a block box that might be encrypting ECB or CBC, tells you which one is happening.

## Solution

### Implementing the oracle

This is the first challenge that asks us to implemenet the "application" that we will be attacking. To do this, I created a struct in the `main` package (usually called `Application`) to encapsult the hidden application data (in this case, the application key and the encryption mode). The application can be created with a function like `makeApplication`.

### Detecting the cipher mode

The oracle gives us a function of the form $O(m) = E(m_p \mid\mid m \mid\mid m_s, k)$ where $m_p$ and $m_s$ are the random prefix and suffix and $E$ is either ECB or CBC.

If we set $m$ to one character repeated many times — say, the character `'a'` repeated 64 times — we can force the plaintext passed to $E$ to have at least two blocks that are all a's. Note that since we do not know the size of the prefix and suffix, some of these 'a's might join the last prefix block and/or the first suffix block. However, at maximum 15*2=30 of the a's may join blocks. So as long as we input at least 30+32=52 a's, we can still create at least two full blocks of all a's. I reuse this is a technique in a few of the later challenges as well; if the bytes that we inject may join with the prefix/suffix, just inject a lot of bytes.

If the oracle encrypts with ECB, the two all a's blocks will be encrypted to the same ciphertext, which we can detect. However, if encrypting with CBC, it is extremly unlikely that any two blocks will be encrypted to the same ciphertext, since each block is XORed against an effectively random block before encrypting it. Therefore, detecting the prescence of two equal ciphertext blocks gives us a good way to determine if ECB is being used.
