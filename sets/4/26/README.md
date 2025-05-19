# Challenge 26

**CTR bitflipping**

## Challenge description

> There are people in the world that believe that CTR resists bit flipping attacks of the kind to which CBC mode is susceptible.
> 
> Re-implement the CBC bitflipping exercise from earlier to use CTR mode instead of CBC mode. Inject an "admin=true" token.

## The application

The application here is very similar to Challenge 16 "CBC bitflipping attacks", except it now signs the data with CTR mode instead of CBC mode. It also generates a random CTR nonce with each call to `CreateDataEncrypted`.

The application exposes the following functions to the user

```go
// CreateDataEncrypted creates a string containing "userdata=<userData>;" signed with AES CTR.
// It returns the AES CTR ciphertext along with the nonce
func (a Application) CreateDataEncrypted(userData string) ([]byte, cipherx.Count)

// IsAdmin takes in encryptedData signed with AES CTR along with the nonce, and checks if the decrypted
// data contains "admin=true;"
func (a Application) IsAdmin(encryptedData []byte, nonce cipherx.Count) (bool, error)
```

## Solution

CTR bit-flipping is actually even easier than CBC bitflipping. Let's input, say, 12 bytes of known plaintext $m$ into `CreateDataEncrypted`; for instance, `m = "aaaaaaaaaaaa"`. Then we get the ciphertext $c$ for

```
comment1=cooking%20MCs;userdata=aaaaaaaaaaaa;comment2=%20like%20a%20pound%20of%20bacon
```

along with the nonce $u$.

We can slice out the ciphertext $c$ for the plaintext we inputted (bytes $`[32, 32+12)`$). We can then find the keystream bytes $`[32, 32+12)`$ (when using nonce $u$) by simply computing $K = m \oplus c$. Then if we want to generate new ciphertext bytes $c'$ that decrypt to a new message $m'$, we just construct $c' := m' \oplus K$ which gives that $`D_\text{CTR}(c') = c' \oplus K = (m' \oplus K) \oplus K = m'`$ as desired. If we let $m'$ equal `a;admin=true"` and then construct $c'$ accordingly, the ciphertext now decrypts to

```
comment1=cooking%20MCs;userdata=a;admin=true;comment2=%20like%20a%20pound%20of%20bacon
```

Now if we pass $c'$ and $u$ back to `IsAdmin`, we get admin access.
