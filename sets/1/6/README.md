# Challenge 6

**Break repeating-key XOR**

## Overview

The challenge asks us to break the repeating-key XOR we implemented in the previous challenge.

## Solution

### Notation

In this writeup and in other writeups I use the following notation. Let $m$ be the plaintext/message, $c$ be the ciphertext, and $k$ be the key. Each of these are bytes sequences. If $x$ is a byte sequence, let $x_i$ denote the $i$-th byte (0-indexed) and $x_{i:j}$ denote bytes in the range $[i, j)$.

In this challenge, the plaintext $m$ was encrypted with a key $k$ of unknown length $s^\star$ so that $c_i = m_i \oplus k_{i \bmod s^\star}$.

### Hamming distance

First the challenge asks us to implement the Hamming distance between two byte/bit sequences. Hamming distance $h(a,b)$ is defined as the number of bits that differ between $a$ and $b$. Another way to think of this is the number of 1's in the string $a \oplus b$, since XORing two bits returns 1 if and only if the bits differ. So if we let $S(x)$ denote the number of 1 bits in $x$, then $h(a, b) = S(a \oplus b)$.

### Cracking the key size

The challenge outlines the following method for guessing the key size. For a key size guess $s$, compute $\frac{1}{s} \cdot h(c_{0:s}, c_{s:2s})$. Pick a few of the key sizes that produced the lowest hamming distance (after dividing by $s$ to average it out), and then try those.

The challenge does not really explain _why_ the average hamming distance between blocks provides a good indicator of key size correctness, so I wanted to explore that a bit more. Let $s^\star$ be the true key size and $s$ be our current guess. Then we compute $h(c_{0:s}, c_{s:2s}) = S(c_{0:s} \oplus c_{s:2s}) = \sum_{i=0}^{s-1} S(c_i \oplus c_{i+s})$. Zooming in a bit, the summand is $S(c_i \oplus c_{i+s}) = S((m_i \oplus k_{i \bmod s^\star}) \oplus (m_{i+s} \oplus k_{i+s \bmod s^\star})) = S((m_i \oplus m_{i+s}) \oplus (k_{i \bmod s^\star} \oplus k_{i+s \bmod s^\star}))$.

The term $k_{i \bmod s^\star} \oplus k_{i+s \bmod s^\star}$ is of interest because there are two cases that it could fall under. If $i \equiv i+s \pmod{s^\star}$, then these are the exact same key bytes and XOR to $0$! This happens iff $0 \equiv s \pmod{s^\star}$, or in other words if the correct key size divides our key size guess ($s^\star | s$). \
Otherwise, these are two _different_ key bytes. Then assuming that each key bit is i.i.d, the number of differing bits between these two key bytes follows a binomial distribution $S(k_{i \bmod s^\star} \oplus k_{i+s \bmod s^\star}) \sim \mathrm{Bin}(8, 1/2)$. So $\mathbb{E}(S(k_{i \bmod s^\star} \oplus k_{i+s \bmod s^\star})) = 8 \cdot \frac{1}{2} = 4$.

Now lets zoom out a bit. In the case where $s^\star | s$, we get that $S((m_i \oplus m_{i+s}) \oplus (k_{i \bmod s^\star} \oplus k_{i+s \bmod s^\star})) = S(m_i \oplus m_{i+s})$. A key insight now is that in an English plaintext, the expected value of $S(m_i \oplus m_{i+s})$ is likely lower than if the two bytes were random. As we saw before, if $x,y$ are random bytes, then $S(x \oplus y) \sim \mathrm{Bin}(8, 1/2)$ so $\mathbb{E}(S(x \oplus y)) = 4$. However, $m_i$ and $m_{i+s}$ are very likely to be lowercase ASCII letters, which typically have many identical bits. Namely, if $x,y$ are both random lowercase ASCII letters, $\mathbb{E}(S(x \oplus y))$ is about $2.47$. More specifically, the probability distribution looks like

| $S(x \oplus y)$ | 0    | 1    | 2    | 3    | 4    | 5    | 6 | 7 | 8 |
|--------------------------|------|------|------|------|------|------|------|------|------|
| $p$                     | 0.04 | 0.16 | 0.31 | 0.30 | 0.16 | 0.03 | 0 | 0 | 0 |

which is quite different from $\mathrm{Bin}(8, 1/2)$.

In the other case where $s^\star \nmid s$, recall that $S(k_{i \bmod s^\star} \oplus k_{i+s \bmod s^\star}) \sim \mathrm{Bin}(8, 1/2)$. However, if we XOR $k_{i \bmod s^\star} \oplus k_{i+s \bmod s^\star}$ by $m_i \oplus m_{i+s}$, the probability distribution does not actually change. This is because XORing it by a fixed byte bijectively maps each possible byte to another unique byte, so all it really does is permute the sample space which does not impact the probability distribution (with respect to $k$). So $S((m_i \oplus m_{i+s}) \oplus (k_{i \bmod s^\star} \oplus k_{i+s \bmod s^\star})) \sim \mathrm{Bin}(8, 1/2)$. Thus $\mathbb{E}(S((m_i \oplus m_{i+s}) \oplus (k_{i \bmod s^\star} \oplus k_{i+s \bmod s^\star})) = 4$.

We can now see that the expected value of the average hamming distance is significantly lower when $s^\star | s$. Notably, even if we pick an $s$ such that $s^\star | s$ but $s > s^\star$, the guess is still technically "correct" since we could just view the key as repeating itself $s / s^\star$ times. Therefore, a low average hamming distance is a good (but not perfect) indicator that our guess may be correct.

### Cracking the plaintext

For each key size guess $s$, we can view $c$ as a sequence of blocks each of size $s$, where each block was XORed against $k$. If we look at the byte 0 of each block, each of these bytes were XORed against the same byte $k_0$. Thus, we can crack the concatenated string $c_0 | c_k | c_{2k} | \ldots$ as single-byte XOR using our function from the previous challenge, yielding us a guess for $k_0$. We can repeat this process for each byte offset: for each $0 \leq i < s$, solve $c_{i} | c_{k+i} | c_{2k+i} | \ldots$ as-single byte XOR to get a guess for $k_i$. In the end, we get a guess for $k$ and can decrypt the ciphertext. We can compute a score for each guess using the same histogram metric from the previous challenges, and then pick the guess with the best score.
