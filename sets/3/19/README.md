# Challenge 19

**Break fixed-nonce CTR mode using substitutions**

## Challenge description

> Take your CTR encrypt/decrypt function and fix its nonce value to 0. Generate a random AES key.
> 
> In successive encryptions (not in one big running CTR stream), encrypt each line of the base64 decodes of the following, producing multiple independent ciphertexts:
> 
> [bunch of strings]
> 
> Because the CTR nonce wasn't randomized for each encryption, each ciphertext has been encrypted against the same keystream. This is very bad.
> 
> Understanding that, like most stream ciphers (including RC4, and obviously any block cipher run in CTR mode), the actual "encryption" of a byte of data boils down to a single XOR operation, it should be plain that:
> 
> CIPHERTEXT-BYTE XOR PLAINTEXT-BYTE = KEYSTREAM-BYTE
>
> And since the keystream is the same for every ciphertext:
> 
> CIPHERTEXT-BYTE XOR KEYSTREAM-BYTE = PLAINTEXT-BYTE (ie, "you don't say!")
>
> Attack this cryptosystem piecemeal: guess letters, use expected English language frequence to validate guesses, catch common English trigrams, and so on.
> 
> Don't overthink it.
>
> Points for automating this, but part of the reason I'm having you > do this is that I think this approach is suboptimal.

## ?

This challenge asks to solve the exact same problem as challenge 20. However, seemingly this challenge wants you to solve it in a different way than challenge 19. However, I could not quite understand what they wanted you to do for this challenge; all of my ideas seemed to fall under the solution method for challenge 20. Thus, I just moved on to challenge 20.
