# Challenge 35

**Implement DH with negotiated groups, and break with malicious "g" parameters**

## Challenge description

> A->B \
> Send "p", "g" \
> B->A \
> Send ACK \
> A->B \
> Send "A" \
> B->A \
> Send "B" \
> A->B \
> Send AES-CBC(SHA1(s)[0:16], iv=random(16), msg) + iv \
> B->A \
> Send AES-CBC(SHA1(s)[0:16], iv=random(16), A's msg) + iv \
> Do the MITM attack again, but play with "g". What happens with:
> ```
> g = 1
> g = p
> g = p - 1
> ```
> Write attacks for each.

## The attacks

### Attack 1

If $g = 1$, then $K_a = 1^a = 1$ and $K_b = 1^b = 1$. Thus, a shared key $K = 1$ is established. The attacker can use this to decrypt the messages.

### Attack 2

If $g = p$, then $K_a = p^a = 0$ and $K_b = p^b = 0$ (under $\mod{p}$). Thus, a shared key $K = 0$ is established. The attacker can use this to decrypt the messages.

### Attack 3

If $g = p - 1$, then $K_a = (p-1)^a = (-1)^a$ (under $\mod{p}$). So $K_a = 1$ if $a$ is even or $K_a = -1 = p - 1$ if $a$ is odd. The same applies for $K_b$. Thus, a shared key _might_ not be established, for instance if $K_a = 1$ and $K_b = -1$, which the clients may notice when they try to decrypt each other's messages. Regardless, the attacker can still decrypt the intercepted messages by trying both $1$ and $p-1$ as the key and inspecting the output.

## Modelling the attacks

Once again, I used the man-in-the-middle simulation framework that I developed in the previous challenge to model this attack. The implementations of client A, client B, and the attacker are pretty similar to the previous challenge as well. Notably, the attacker now accepts a function `func (p *big.Int) *big.Int` to compute the malicious "g", which differs between each of the 3 attacks.
