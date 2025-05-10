# Challenge 5

**Implement repeating-key XOR**

## Solution

To implement this, we simply just repeatedly XOR the key with the plaintext. For example, if the plaintext is `Hello world!` and the key is `ICE`, we XOR the following two byte sequences:

`Hello world!` \
`ICEICEICEICE`
