# Challenge 15

**The CBC padding oracle**

## Solution

### The oracle

The application provides an oracle of the form

$$
O(c, v) = \mathrm{isPadValid}(D_\mathrm{CBC}(c, v, k))
$$

where $\mathrm{isPadValid}$ returns true iff the final plaintext block has valid PKCS7 padding. Valid PKCS7 padding is defined in my write-up for challenge 15.

In my solution I only ever input a single block to the oracle,  so for simplicity we can write it as

$$
\begin{align*}
O(c, v) &= \mathrm{isPadValid}(D_\mathrm{CBC\_block}(c, v, k)) \\
&= \mathrm{isPadValid}(D_\mathrm{AES}(c, k) \oplus v)
\end{align*}
$$

### The attack

The application gives us a ciphertext and IV encrypted with AES CBC, and our goal is to determine the plaintext. We can perform this attack on each block individually, so let $c$ be the current ciphertext block, $v$ be the previous ciphertext block (or the IV if $c$ is the first block), and $m$ be the plaintext block corresponding to $c$. So we have the relationship $c = E_\mathrm{CBC\_block}(c, v, k) = E_\mathrm{AES}(m \oplus v, k)$. Given $c,v$ we want to find $m$.

#### Learning the final byte of $m$ ($m_{s-1}$)

Let's first just look at how we can determine the final byte of $m$ using the oracle $O$.

Recall that if we replace $v$ with $v'$, then $c$ now decrypts to $m' := D_\mathrm{CBC\_block}(c, v', k) = D_\mathrm{AES}(c, k) \oplus v' = m \oplus v \oplus v'$. So $O(c, v')$ returns true iff $m' = m \oplus v \oplus v'$ has valid PKCS7 padding. Thus, by varying the value of $v'_{s-1}$ and checking $O(c, v')$, we can learn information about $m'_{s-1}$, and by extension $m_{s-1}$.

For each possible byte $0 \leq b < 256$, let $v' = (0)_{\times s-1} \mid\mid b$ and run $O(c,v')$. Iff $O$ returns true, then we know that $m_{s-1}' = m_{s-1} \oplus v_{s-1} \oplus b$ is a valid padding byte! So $1 \leq m_{s-1}' \leq s$.

However, we would like to pinpoint an exact value for $m_{s-1}'$; for instance, $m_{s-1}'=1$. To do this, whenever we find a $b$ such that $O(c, v')$ returns true, we can let $v'' = (0)_{\times s-2} \mid\mid 1 \mid\mid b$ and run $O(c,v'')$. Compared to $v'$, all we did was change the last $0$ to a $1$, so we have that $m''_{s-2} \neq m'_{s-2}$ but still $m''_{s-1} = m'_{s-1}$. If this $O(c,v'')$ _also_ returns true, then we know that $m'_{s-1}$ is a valid padding byte even if we change the value of the byte proceeding it. Therefore, it must be that $m'_{s-1} = 1$, because if it was any value $>1$, $m'_{s-2}$ or $m''_{s-2}$ would have differed from $m'_{s-1}$ and caused invalid padding.

Then we finally learn that $m_{s-1} = m_{s-1}' \oplus v_{s-1} \oplus b = 1 \oplus v_{s-1} \oplus b$.

#### Learning byte $m_i$

We will start with $i = s-2$ and iterate to $i = 0$. By induction, assume that we have already solved for bytes $m_{i+1:s}$. We will now solve for $m_i$, which is actually a bit simpler than solving for the last byte.

**Idea:** Since $m_i$ is the $(s-i)$-th byte from the right of $m$, we will craft $v'$ so that bytes $m_{i+1:s}'$ equal $s-i$ and then search for a $v'_{i}$ that makes $m'_{i} = s-1$, causing $m'$ to have valid padding (since the last $s-i$ bytes would equal $s-i$, in this case).

For each $0 \leq b < 256$, let $v' = (0)_{\times i} \mid\mid b \mid\mid (v_{i+1:s} \oplus m_{i+1:s} \oplus (s-i)_{\times s-i-1})$ and run $O(c, v')$. For each $i+1 \leq j < s$, we have that $m'_j = m_j \oplus v_j \oplus v'_j = m_j \oplus v_j \oplus (v_j \oplus m_j \oplus (s-i)) = s-i$. Thus, $O(c, v')$ returns true iff $m'$ has valid padding iff $m_i' = s-i$. So when $O(c, v')$ does return true (which will happen _exactly once_), we learn that $m_i = m_i' \oplus v_i \oplus b = (s-i) \oplus v_i \oplus b$.

By repeating the above process, we learn all the bytes of $m$.


<!-- Likewise, if we only change byte $i$ of $v$ to $v'_i = v_i \oplus m_i \oplus b$, then $D_\mathrm{CBC\_block}(c, v', k)_i = b$.  -->

<!-- Recall from Challenge 15 "CBC bitflipping attacks" that if we set $v' = v \oplus m \oplus m'$, then $D_\mathrm{CBC\_block}(c, v', k) = m'$. Likewise, if we only change byte $i$ of $v$ to $v'_i = v_i \oplus m_i \oplus b$, then $D_\mathrm{CBC\_block}(c, v', k)_i = b$.  -->

<!-- Focusing on the last byte $i = s-1$, if we happen to chose a $b$ such that $D_\mathrm{CBC\_block}(c, v', k)_i = 1$ (the byte with _value_ 1) -->