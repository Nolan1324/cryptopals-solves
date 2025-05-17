# Challenge 12

**Byte-at-a-time ECB decryption (Simple)**

## Challenge Description

> Copy your oracle function to a new function that encrypts buffers under ECB mode using a consistent but unknown key (for instance, assign a single random key, once, to a global variable).
> 
> Now take that same function and have it append to the plaintext, BEFORE ENCRYPTING, the following string:
> 
> ```
> Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkg
> aGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBq
> dXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUg
> YnkK
> ```
>
> Spoiler alert.
> Do not decode this string now. Don't do it.
> 
> Base64 decode the string before appending it. Do not base64 decode the string by hand; make your code do it. The point is that you don't know its contents.
> 
> What you have now is a function that produces:
> 
> `AES-128-ECB(your-string \mid\mid unknown-string, random-key)`
>
> It turns out: you can decrypt "unknown-string" with repeated calls to the oracle function!
> 
> Here's roughly how:
> 
> 1. Feed identical bytes of your-string to the function 1 at a time --- start with 1 byte ("A"), then "AA", then "AAA" and so on. Discover the block size of the cipher. You know it, but do this step anyway.
> 2. Detect that the function is using ECB. You already know, but do this step anyways.
> 3. Knowing the block size, craft an input block that is exactly 1 byte short (for instance, if the block size is 8 bytes, make "AAAAAAA"). Think about what the oracle function is going to put in that last byte position.
> 4. Make a dictionary of every possible last byte by feeding different strings to the oracle; for instance, "AAAAAAAA", "AAAAAAAB", "AAAAAAAC", remembering the first block of each invocation.
> 5. Match the output of the one-byte-short input to one of the entries in your dictionary. You've now discovered the first byte of unknown-string.
> 6. Repeat for the next byte.

## Solution

The challenge description already gives a general overview of the solution, so I will focus a bit more on the implementation details.

### The oracle

The oracle gives us a function $`O(m) = E_\text{ECB}(m \mid\mid t, k)`$ where $`t`$ is the target plaintext we seek to learn.

### Determining the block size $`s`$

We will devise a slightly stronger algorithm for determining the block size by assuming the oracle is of the form $`O(m) = E_\text{ECB}(p \mid\mid m \mid\mid t, k)`$ where $`p`$ is some unknown prefix. In this challenge, $`p = \varepsilon`$.

First set $`m = \varepsilon`$, call $`c := O(m)`$, and note down $`|c|`$ (the length in bytes of $`c`$). Then repeatedly call $`c' := O(m)`$ until $`c' > c`$. When this occurs, this means that we have extended the input to $`E_\text{ECB}`$ enough to create a new block. Therefore, $`s := c' - c`$ is the block size.

Notably, in real AES, the only valid block size is $`s = 16`$, which we take for granted in later challenges to streamline things a bit.

### Detecting ECB mode

We once again devise a slightly stronger algorithm where $`O(m) = E_\text{ECB}(p \mid\mid m \mid\mid t, k)`$. We can use the exact same algorithm as the previous challenge where we set $`m`$ to the same character repeated many ($`4s`$) times to force the creation of two identical plaintext blocks. If the ciphertext contains two duplicate ciphertext blocks, then we have detected ECB mode.

### Cracking the target text

For this algorithm, we will use the original oracle $`O(m) = E_\text{ECB}(m \mid\mid t, k)`$ that has no prefix.

We devise an inductive algorithm such that for any $`s - 1 \leq i < |c|`$, if we know the previous $`s-1`$ bytes $`t_{i-(s-1) ~:~ i}`$ of the target text, we can determine the next byte $`t_i`$.

#### Algorithm inductive step

**Dictionary construction:** First we perform the dictionary construction step. For each of the 255 possible guesses $`g`$ for the byte $`t_i`$, we construct the guess block $`m := t_{i-(s-1) ~:~ i} \mid\mid g`$ and then encrypt it with $`c := O(m)`$. At this step, we are only interested in what the block $`m`$ encrypts to, so we just take the first block $`c_{0:s}`$ of the resulting ciphertext and store the ciphertext-to-guess-byte mapping $`c_{0:s} \to g`$ in our dictionary.

**Ciphertext discovery:** We now need to discover the actual ciphertext for the full block $`t_{i-(s-1) ~:~ i} \mid\mid t_i`$ (aka, $`t_{i-(s-1) ~:~ i+1}`$). To do so, we need to set $`m`$ to a certain number of padding bytes so that $`t_{i-(s-1)}`$ lies at the start of some block in $`m \mid\mid t`$. Namely, if we set $`m`$ to $`j = s - 1 - (i \bmod s)`$ padding bytes, then the new index of $`t_{i-(s-1)}`$ in $`m \mid\mid t`$ becomes $`j + (i-(s-1)) = s - 1 - (i \bmod s) + (i-(s-1))= i - (i \bmod{s}) = i - (i - \lfloor \frac{i}{s} \rfloor s) = \lfloor \frac{i}{s} \rfloor s`$, which is at the start of block number $`b = \lfloor \frac{i}{s} \rfloor`$. \
Thus, to compute the ciphertext of block $`t_{i-(s-1) ~:~ i} \mid\mid t_i`$, we call $`c := O(m)`$ where $`m`$ equals $`j`$ padding bytes, and then grab the $`b`$-th ciphertext block $`c' = c_{s b ~:~ s (b+1)}`$. We then query our dictionary with $`c'`$ to find which guess byte $`g`$ results in $`E_\text{AES}(t_{i-(s-1) ~:~ i} \mid\mid t_i, k) = E_\text{AES}(t_{i-(s-1) ~:~ i} \mid\mid g, k)`$, and thus $`t_{i-(s-1) ~:~ i} \mid\mid t_i = t_{i-(s-1) ~:~ i} \mid\mid g`$, and thus $`t_i = g`$.

#### Handling the base cases ($`0 \leq i < s-1`$)

The above inductive step relies on us knowing the previous $`s-1`$ bytes of $`t`$. However, when $`0 \leq i < s-1`$, these are bytes are out of bounds on the left. To resolve this, we pretend as though $`t`$ starts with $`s-1`$ bytes of padding (say, the character $`a`$). Call this padded string $`t'`$, so $`i' = i + (s-1)`$ is the index of $`t_i`$ in $`t'`$.

Now, in the dictionary construction step, we can construct our guess block as $`m := t'_{i'-(s-1) ~:~ i'} \mid\mid g`$. Since we use $`t'`$ here, the indices are not out of bounds, and the guess block may start with some padding bytes.

In the ciphertext discovery step, we set $`m`$ to $`j`$ padding bytes as before. However, if $`0 \leq i < s-1`$, then we are inspecting the resulting ciphertext block $`c'`$ with block number 0, **for which the respective plaintext block starts with $`j`$ bytes of padding.** This is consistent with the dictionary construction step, in which the guess block also started with $`j`$ bytes of padding. Thus, this still gives us a correct ciphertext to query the dictionary with.

Notice that the padding bytes now do double duty in the ciphertext discovery step: they are used to shift the target text to align the guess block to the start of block number $`\lfloor \frac{i}{s} \rfloor`$, and they are used to make block 0 match the padded guest blocks when $`0 \leq i < s-1`$.

#### Full algorithm

We can merge the base case logic with the inductive step to arrive at a final, modified algorithm.

For each $`0 \leq i < \text{len}(t)`$, do the following:

**Dictionary construction:** Let $`t'`$ be $`t`$ prepended with $`s-1`$ bytes of padding. So $`i' = i + (s-1)`$ is the index of $`t_i`$ in $`t'`$. For each of the 255 possible guesses $`g`$ for the byte $`t_i`$, we construct the guess block $`m := t'_{i'-(s-1) ~:~ i'} \mid\mid g`$ and then encrypt it with $`c := O(m)`$. Take the first block $`c_{0:s}`$ of the resulting ciphertext and store the mapping $`c_{0:s} \to g`$ in our dictionary.

**Ciphertext discovery:** We now discover the ciphertext for block $`t'_{i'-(s-1) ~:~ i'} \mid\mid t_i`$. Set $`m`$ to $`j = s - 1 - (i \bmod s)`$ padding bytes so that the new index of $`t'_{i'-(s-1)} = t_{i-(s-1)}`$ in $`m \mid\mid t`$ becomes $`j + (i-(s-1)) = s - i (\bmod s) + (i-(s-1))= i - (i \bmod{s}) = i - (i - \lfloor \frac{i}{s} \rfloor s) = \lfloor \frac{i}{s} \rfloor s`$, which is at the start of block number $`b = \lfloor \frac{i}{s} \rfloor`$. \
**Notice that if $`0 \leq i < s-1`$, even though $`t_{i-(s-1)}`$ is out of bounds, $`t'_{i'-(s-1)} = t'_{i}`$ is still in bounds, and since we set $`m`$ to equal $`j`$ padding bytes, we in fact still get that $`(m \mid\mid t)_{i'-(s-1)} = t_{i'-(s-1)}`$ as desired.** \
To compute $`E_\text{AES}(t'_{i'-(s-1) ~:~ i'} \mid\mid t_i, k)`$, call $`c := O(m)`$, and then grab block $`c' = c_{s b ~:~ s (b+1)}`$. Then query our dictionary with $`c'`$ to find which guess byte $`g`$ results in $`t'_{i'-(s-1) ~:~ i'} \mid\mid t_i = t'_{i'-(s-1) ~:~ i'} \mid\mid g`$, and thus $`t_i = g`$.

### Detecting the ciphertext length (optional)

