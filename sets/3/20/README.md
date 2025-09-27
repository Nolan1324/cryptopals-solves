# Challenge 20

**Break fixed-nonce CTR statistically**

## Challenge description

> In this file find a similar set of Base64'd plaintext. Do with them exactly what you did with the first, but solve the problem differently.
> 
> Instead of making spot guesses at to known plaintext, treat the collection of ciphertexts the same way you would repeating-key XOR.
> 
> Obviously, CTR encryption appears different from repeated-key XOR, but with a fixed nonce they are effectively the same thing.
> 
> To exploit this: take your collection of ciphertexts and truncate them to a common length (the length of the smallest ciphertext will work).
> 
> Solve the resulting concatenation of ciphertexts as if for repeating-key XOR, with a key size of the length of the ciphertext you XOR'd.

## Solution

Thankfully, the code we implemented to solve repeating key XOR works pretty directly for this challenge. Let $n$ be the number of ciphertexts, $c^j$ be ciphertext $j$, $m^j$ be message $j$, and $K$ be the keystream generated from AES key $k$ and nonce $0$ (could write this as $(k, 0) \mapsto K$).

It is clear that for each $i$, for each $j$ we have $c_i^j = m_i^j \oplus K_i$. So $c_i^0 \mid\mid c_i^1 \mid\mid \ldots \mid\mid c_i^n = (m_i^0 \mid\mid m_i^1 \mid\mid \ldots \mid\mid m_i^n) \oplus K_i$. This is just single-byte XOR, so we can solve this by trying all guesses for $K_i$ and choosing the one with the best histogram score, as we did in the past (also like before, we assume the plaintexts are English messages).

### Caveats

There are a few caveats with this approach in practice. First of all, the ciphertexts are not all equal length. So for larger $i$, there might not exist a byte $i$ in every ciphertext. This gives us less bytes to build the observed histogram from and thus makes accurately scoring the guesses harder. Thus, I was unable to automatically solve for the last few characters of the longest ciphertext.

Another caveat is that when $i=0$, it turns out that the histogram looks for $m_i^0 \mid\mid m_i^1 \mid\mid \ldots \mid\mid m_i^n$ looks quite different than usual. This is because the first character of each plaintext is likely a capital letter. However, in the expected histogram that we compare against, capital letters are much rarer than lower case letters. Therefore, the correct guess here gets a comparatively low score, resulting in an incorrect solution. The simplest way to solve this would just find the best guess when $i=0$ by manual inspection (it is only 256 guesses, after all). However, I still wanted to automate it, so I generated a new expected histogram **from only the first character of each sentence in the original corpus from `nltk`**, and then used this expected histogram when $i=0$. This ended up finding the correct solution.

