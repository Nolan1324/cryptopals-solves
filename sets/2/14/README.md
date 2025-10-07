# Challenge 14

**Byte-at-a-time ECB decryption (Harder)**

## Challenge Description

> Take your oracle function from #12. Now generate a random count of random bytes and prepend this string to every plaintext. You are now doing:
> 
> ```
> AES-128-ECB(random-prefix || attacker-controlled || target-bytes, random-key)
> ```
> 
> Same goal: decrypt the target-bytes.

## Overview

This challenge is very similar to challenge 12, except the oracle now prepends the plaintext with an unknown prefix. I approach this challenge by constructing a reduction from this challenge to challenge 12.

## Solution

### The oracle

The oracle gives us a function $O(m) = E_\text{ECB}(p \mid\mid m \mid\mid t, k)$ where $t$ is the target plaintext we seek to learn and $p$ is some unknown prefix.

### Determining the block size and ECB mode

In challenge 12, we already constructed algorithms to determine block size and ECB mode using an oracle of the form $O(m) = E_\text{ECB}(p \mid\mid m \mid\mid t, k)$ (even though we had $p = \varepsilon$ in challenge 12). Thus, we can just use these same algorithms again.

### Cracking the target text

If we can reduce the oracle $O(m) = E_\text{ECB}(p \mid\mid m \mid\mid t, k)$ to another oracle of the form $O'(m) = E_\text{ECB}(m \mid\mid t, k)$, then we can just invoke our algorithm from challenge 12 on $O'$.

#### The reduction

If we know the length of $p$ ($|p|$), then we can preform the reduction as follows. $\lfloor \frac{|p|}{s} \rfloor$ is the number of the block containing the final byte of $p$. $|p| \bmod s$ is how many bytes of $p$ are in the final block of $p$. In order to "ignore" $p$ in the oracle, we will prepend $s - (|p| \bmod s)$ padding bytes to $m$, essentially "filling" the remaining bytes in the final $p$ block with padding. Then, after encrypting with $c := O(m)$, we will discard all the prefix ciphertext blocks (the blocks up to and including block number $\lfloor \frac{|p|}{s} \rfloor$). Thus, we get the reduction

$$O'(m) := O(s - (|p| \bmod s) \text{ padding bytes} \mid\mid m)_{s \left( \left\lfloor \frac{|p|}{s} \right\rfloor + 1 \right) :}$$

which satisifies that $O'(m) = E_\text{ECB}(m \mid\mid t, k)$, as desired.

#### Determing the length of $p$

Now we just need to determine $|p|$ so that we can preform the reduction.

We first determine the ciphertext for a known plaintext block; namely, the block containing all 0's; call this the $z$ block. As done in previous algorithms, we set $m$ to many ($4s$) 0 bytes to force $c := O(m)$ to have at least two consecutive identical blocks. We detect these blocks in the ciphertext and conclude that they correspond to the two $z$ blocks in the plaintext*.

We will now initialize $m = \varepsilon$ and incrementally add a 0 byte to $m$ until we observe the ciphertext block for $z$ in $c := O(m)$ \*\*. Once we detect the $z$ block, this means that $m$ was long enough to fill the remaining bytes in the final block of $p$ (since $O$ directly appends $m$ to $p$) **and** fill an entire new block of size of $s$.

It is possible that the byte directly _after_ $m$ (the first byte of the target text $t$) also equaled the zero byte, which would have caused the new $z$ block to form earlier than expected. To resolve this, we include some other byte, (say, `'x'`) at the end of $m$ in every call to $O(m)$ to seperate $m$ from the first byte of $t$.

When we first detect the $z$ block, let $\ell$ be the length of $m$ (excluding the separation byte at the end) and let $b$ be the block number of the detected $z$ block. Then other than the zero bytes that formed the $z$ block, $\ell - s$ zero bytes joined with the end of the prefix to complete the final prefix block.
So the prefix length is $s \cdot b$ (the index where the prefix block stops and the zero block starts) minus the extra zero bytes in the prefix block, so $|p| = s \cdot b - (\ell - s)$.

<sub>*This assumes the original plaintext had no other consecutive duplicate blocks. If it did, we could have just injected more 0 bytes to create more consecutive identical blocks, and then detected this larger set of identical blocks in the ciphertext. Alternatively, we could have incrementally added a single 0 byte to $m$ and encrypted until we observed a new duplicate block appear in the ciphertext.</sub>

<sub>\*\*This assumes the original plaintext did not contain an all 0's block. If it did, we could have resolved it by just injecting a different repeated byte in the previous step (the choice constructing an all 0's block was arbitrary).
