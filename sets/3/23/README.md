# Challenge 23

**Clone an MT19937 RNG from its output**

## Challenge description

> The internal state of MT19937 consists of 624 32 bit integers.
> 
> For each batch of 624 outputs, MT permutes that internal state. By permuting state regularly, MT19937 achieves a period of 2**19937, which is Big.
> 
> Each time MT19937 is tapped, an element of its internal state is subjected to a tempering function that diffuses bits through the result.
> 
> The tempering function is invertible; you can write an "untemper" function that takes an MT19937 output and transforms it back into the corresponding element of the MT19937 state array.
> 
> To invert the temper transform, apply the inverse of each of the operations in the temper transform in reverse order. There are two kinds of operations in the temper transform each applied twice; one is an XOR against a right-shifted value, and the other is an XOR against a left-shifted value AND'd with a magic number. So you'll need code to invert the "right" and the "left" operation.
> 
> Once you have "untemper" working, create a new MT19937 generator, tap it for 624 outputs, untemper each of them to recreate the state of the generator, and splice that state into a new instance of the MT19937 generator.
> 
> The new "spliced" generator should predict the values of the original.

## Solution

### Inverting the tempering function

Let $x$ be a 32-bit value generated from MT19937 before tempering it. Let $b,c$ be 32-bit integers, and let $1 \leq u,s,t,l \leq 32$ be integers.
The tempering function involves the following sequences of operations.

$$
\begin{align*}y&= x\oplus (x\gg u)\\
y &= y \oplus ((y\ll s) ~ \And ~ b)\\
y &= y \oplus ((y\ll t) ~ \And ~ c)\\
&= y\oplus (y\gg l)
\end{align*}
$$

To invert these operations, there are two different types of operations we need to invert: $y = x \oplus ((x \ll s) ~ \And ~ m)$ and $y = x \oplus (x \gg s)$.

#### Inverting $y = x \oplus ((x \ll s) ~ \And ~ m)$

We want to invert the function

$$
y = x \oplus ((x \ll s) ~ \And ~ m)
$$

where $x$ is a 32-bit input, $y$ is the output, $1 \leq s \leq 32$ is the shift amount, and $m$ is a 32-bit mask.

It is a bit easier to analyze the function if we define it bitwise. Let $0 \leq i < 32$ and let $x_i$ denote the $i$-th of bit $x$. Then

$$
y_i = 
\begin{cases}
x_i \oplus (x_{i-s} ~ \And ~ m_i) & i \geq s \\
x_i & \text{otherwise}
\end{cases}
$$

We can then write the inverse as a recursive function

$$
x_i = 
\begin{cases}
y_i \oplus (x_{i-s} ~ \And ~ m_i) & i \geq s \\
y_i & \text{otherwise}
\end{cases}
$$

This is a recursive function because we can compute $x_i$ from $x_{i-s}$, which we compute from $x_{i-2s}$, etc until we reach some base case $x_j$ where $j < s$.

We can also use this recursive function to compute $x$ from $y$ in an iterative fashion. Initialize $x_i = y_i$ for all $0 \leq i < s$. Then for $i = s$ to $31$, compute $x_i = y_i \oplus (x_{i-s} ~ \And ~ m_i)$. (Technically, this is dynamic programming!). Now we have a method to compute $x$ from $y$!

#### Inverting $y = x \oplus (x \gg s)$

We want to invert the function

$$
y = x \oplus (x \gg s)
$$

where $x$ is a 32-bit input, $y$ is the output and $1 \leq s \leq 32$ is the shift amount.

As before, define the function bitwise

$$
y_i = 
\begin{cases}
x_i \oplus x_{i+s} & i < 32 - s \\
x_i & \text{otherwise}
\end{cases}
$$

We can then write the inverse as a recursive function

$$
x_i = 
\begin{cases}
y_i \oplus x_{i+s} & i < 32 - s \\
y_i & \text{otherwise}
\end{cases}
$$

Like before, we use this recursive function to compute $x$ from $y$ in an iterative fashion. However, with this function, we must start at the highest bit of $x$ and iterate to the lowest bit. Initialize $x_i = y_i$ for all $32 - s \leq i < 32$. Then for $i = 31 - s$ to $0$, compute $x_i = y_i \oplus x_{i+s}$.

### Reconstructing the state

Now that we have an untemper function, we can reconstruct the state of a MT19937 instance given its first $n = 624$ 32-bit outputs. We are given the outputs $z_n, \ldots, z_{2n-1}$ (since the first value MT19937 outputs is $z_n$, not $z_0$). Applying the untemper function to each one gives $x_n, \ldots, x_{2n-1}$, which is exactly equal to the state vector at the time right after $z_{2n-1}$ was outputted. Thus, we can initialize our own MT19937 instance with this state vector, and this instance will act as an exact "clone" of the original instance.

