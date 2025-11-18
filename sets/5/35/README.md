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

If $g = 1$, then $A = 1^a = 1$ and $B = 1^b = 1$. So $K_a = B ^ a = 1^a = 1$ and $K_b = A ^ b = 1^b = 1$. Thus, a shared key $K = 1$ is established. The attacker can use this to decrypt the messages.

### Attack 2

If $g = 0$, then $A = 0^a = 0$ and $B = 0^b = 0$ (under $\mod{p}$). So $K_a = B ^ a = 0^a = 0$ and $K_b = A ^ b = 0^b = 0$. Thus, a shared key $K = 0$ is established. The attacker can use this to decrypt the messages.

### Attack 3

If $g = p - 1$, then $A = (p-1)^a = (-1)^a$ and $B = (p-1)^b = (-1)^b$ (under $\mod{p}$). So $K_a = B ^ a = (-1)^{ab}$ and $K_b = A ^ b = (-1)^{ab}$. So $K = 1$ if $ab$ is even and $K = -1 = p - 1$ if $ab$ is odd. Thus, the attacker can decrypt the intercepted messages by trying both $1$ and $p-1$ as the key and inspecting the output.

While just trying both possibilities for the key will work, we can also see which is more likely as an exercise. $a$ and $b$ are generated uniformly at random, so each is odd w.p. (with probability) $1/2$. $ab$ is odd iff $a$ and $b$ are both odd, so $ab$ is odd w.p. $1/2 + 1/2 = 1/4$. Thus, $K=1$ w.p. $0.75$ and $K = p - 1$ w.p. $0.25$.

## Modelling the attacks

Once again, I used the man-in-the-middle simulation framework that I developed in the previous challenge to model this attack. The implementations of client A, client B, and the attacker are pretty similar to the previous challenge as well. Notably, the attacker now accepts a function `func (p *big.Int) *big.Int` to compute the malicious "g", which differs between each of the 3 attacks.
