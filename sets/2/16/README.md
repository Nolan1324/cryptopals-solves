# Challenge 15

**CBC bitflipping attacks**

## Challenge description

> Generate a random AES key.
> 
> Combine your padding code and CBC code to write two functions.
> 
> The first function should take an arbitrary input string, prepend the string:
> 
> ```
> "comment1=cooking%20MCs;userdata="
> ```
> 
> .. and append the string:
> 
> ```
> ";comment2=%20like%20a%20pound%20of%20bacon"
> ```
> 
> The function should quote out the ";" and "=" characters.
> 
> The function should then pad out the input to the 16-byte AES block length and encrypt it under the random AES key.
> 
> The second function should decrypt the string and look for the characters ";admin=true;" (or, equivalently, decrypt, split the string on ";", convert each resulting string into 2-tuples, and look for the "admin" tuple).
> 
> Return true or false based on whether the string exists.
> 
> If you've written the first function properly, it should not be possible to provide user input to it that will generate the string the second function is looking for. We'll have to break the crypto to do that.
> 
> Instead, modify the ciphertext (without knowledge of the AES key) to accomplish this.
> 
> You're relying on the fact that in CBC mode, a 1-bit error in a ciphertext block:
> - Completely scrambles the block the error occurs in
> - Produces the identical 1-bit error(/edit) in the next ciphertext block.

## Overview

Similar to challenge 12 "ECB cut-and-paste", in this challenge encryption is used for integrity rather than confidentiality. The application exposes a function `createData(data)` that returns the string `comment1=cooking%20MCs;userdata={data};comment2=%20like%20a%20pound%20of%20bacon` encrypted with AES CBC using private key $k$. The user can then later pass the ciphertext back to the application's `isAdmin` function; the app will decrypt the ciphertext with $k$ and check if `admin=true` is set. `createData` forbids the `;` and `=` characters, so the user cannot call `createData(";admin=true")` to create encrypted data with `admin=true`. Thus, the attacker once again needs to fabricate a ciphertext that decrypts to a string with `admin=true` set. However, the server is now using CBC mode, so we cannot just cut and paste blocks like in challenge 12.

## The application

The `createData` function URL encodes the input. Thus, the user cannot directly inject `;` or `=` into the data to create encrypted data with `admin=true`.

## Solution

### CBC bitflipping

Recall the encrypt/decrypt functions for a CBC block. Let $m$ be the current plaintext block, $c$ be the current ciphertext block, and $v$ be the previous ciphertext block (or the IV). Then the encrypt/decrypt functions for this block are:

$E_\mathrm{CBC\_block}(m, v, k) = E_\mathrm{AES}(m \oplus v, k)$

$D_\mathrm{CBC\_block}(c, v, k) = D_\mathrm{AES}(c, k) \oplus v$

Assume that the application has already performed the encryption, so we have $c$ and $v$. Also assume we know $m$. Focusing on the decrypt function, we can that by tampering with the ciphertext block $v$, we gain complete control over what $c$ decrypts to. Namely, if we want $c$ to decrypt to $m'$, we just change $v$ to $v' := v \oplus m \oplus m'$, which results in

$$
\begin{align*}
D_\mathrm{CBC\_block}(c, v', k) &= D_\mathrm{AES}(c, k) \oplus v' \\
&= (m \oplus v) \oplus (v \oplus m \oplus m') \\
&= m'
\end{align*}
$$

However, this also means that when $v'$ is decrypted in the previous iteration with $D_\mathrm{CBC\_block}(v', w, k)$ where $w$ is the ciphertext block (or IV) before $v/v'$, the result will be completely scrambled, even if we only changed one bit in $v$ to get $v'$.

Therefore, we gain the ability to completely control what one block decrypts to at the cost of completely scrambling the previous block. This is what CBC bitflipping refers to.

### The attack

I am going to cheat with my own notation a bit here and have $m_i$/$c_i$ represent **block** $i$ rather than byte $i$.

Based on the CBC bitflipping method, it seems like we may want to inject two blocks into the plaintext (injecting only one could work as well but I find two a bit easier to reason about). If we input a two-block-long string into `createData`, the application returns ciphertext $c$ for:

```
comment1=cooking
%20MCs;userdata=
[DATA BLOCK 2  ]
[DATA BLOCK 3  ]
;comment2=%20lik
e%20a%20pound%20
of%20bacon[PCKS]
```

Incidentally, our input is already aligned to block boundaries without adding any padding. Let $m_2, m_3$ be the two blocks we inputted. It turns out the content of these blocks does not matter too much, so I just input randomly generated letters. We want to avoid inputting any special characters because `createData` will URL encode them and thus change the length of the plaintext.

Now, CBC bitflipping allows us to alter what $c_3$ decrypts to by tampering with $c_2$, at the cost of scrambling what $c_2$ decrypts to. If we preform CBC bitflipping to make $c_3$ decrypt to `;admin=true;a=aa`, $c$ now decrypts to

```
comment1=cooking
%20MCs;userdata=
[SCRAMBLED DATA]
;admin=true;a=aa
;comment2=%20lik
e%20a%20pound%20
of%20bacon[PCKS]
```

This gives `admin=true` in the string, as desired. It also causes `userdata` to equal the scrambled block, but thats fine as long as the scrambled data is valid.

Now, to actually perform CBC bit-flipping, let $m_3'$ equal `;admin=true;a=aa`. Then let $c_2' = c_2 \oplus m_3 \oplus m_3'$. Then we make the tampered ciphertext $c' = c_{0:2} || c_2' || c_3 || c_{4:}$, which decrypts to $m' = m_{0:2} || \text{scrambled} || m_3' || m_{4:}$ as desired.

#### Errors caused by the scrambled block

When my application decrypts the encrypted string to check for `admin=true`, it also checks if the string is malformed. For instance, the string `userdata=a;=a=a` is malformed because not every `=` is paired with a `;`. When we caused $m_2$ to become scrambled, we might have incidentally made the decrypted string malformed. If this occurs, I simply retry the attack with different random bytes for $m_3$ and $m_4$ until it succeeds.
