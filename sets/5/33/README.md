# Challenge 33

**Implement Diffie-Hellman**

## Challenge description

> For one of the most important algorithms in cryptography this exercise couldn't be a whole lot easier.> 
> Set a variable "p" to 37 and "g" to 5. This algorithm is so easy I'm not even going to explain it. Just do what I do.> 
> Generate "a", a random number mod 37. Now generate "A", which is "g" raised to the "a" power mode 37 --- A = (g**a) % p.> 
> Do the same for "b" and "B".> 
> "A" and "B" are public keys. Generate a session key with them; set "s" to "B" raised to the "a" power mod 37 --- s = (B**a) % p.> 
> Do the same with A**b, check that you come up with the same "s".> 
> To turn "s" into a key, you can just hash it to create 128 bits of key material (or SHA256 it to create a key for encrypting and a key for a MAC).> 
> Ok, that was fun, now repeat the exercise with bignums like in the real world. Here are parameters NIST likes:> 
> ```
> p:
> ffffffffffffffffc90fdaa22168c234c4c6628b80dc1cd129024
> e088a67cc74020bbea63b139b22514a08798e3404ddef9519b3cd
> 3a431b302b0a6df25f14374fe1356d6d51c245e485b576625e7ec
> 6f44c42e9a637ed6b0bff5cb6f406b7edee386bfb5a899fa5ae9f
> 24117c4b1fe649286651ece45b3dc2007cb8a163bf0598da48361
> c55d39a69163fa8fd24cf5f83655d23dca3ad961c62f356208552
> bb9ed529077096966d670c354e4abc9804f1746c08ca237327fff
> fffffffffffff
>  
> g: 2
> ```
> This is very easy to do in Python or Ruby or other high-level languages that auto-promote fixnums to bignums, but it isn't "hard" anywhere.> 
> Note that you'll need to write your own modexp (this is blackboard math, don't freak out), because you'll blow out your bignum library raising "a" to the 1024-bit-numberth power. You can find modexp routines on Rosetta Code for most languages.> 

## Background: modular arithmetic

This document assumes some background on modular arithmetic, but a refresher is provided here.

Something to remember when working with modular arithmetic is that there are many ways to think and write about it, but the mechanics and result are often the same.

### Modular equivalence

Let $a, b$ be integers and $n$ be a positive integer. The notation $a \equiv b \pmod{n}$ reads as "$a$ and $b$ are equivalent mod $n$" and $n$ is referred to as the modulus.

We define modular equivalence as follows:

> $a \equiv b \pmod{n}$ iff there exists some integer $k$ such that $a = b + kn$.

Notice that $a = b + kn$ iff $a - b = kn$ iff $n \mid (a - b)$ ("$n$ divides $a-b$"). This gives us an equivalent definition.

> $a \equiv b \pmod{n}$ iff $n \mid (a - b)$.

Notice critically that $kn \equiv 0 \pmod{n}$ for any integer $k$. In other words, every multiple of $n$ is equivalent to $0$.

Modular arithmetic respects addition and multiplication. Formally, this means that if $a \equiv b \pmod{n}$ and $c \equiv d \pmod{n}$, then

- $a + c \equiv b + d \pmod{n}$
- $ac \equiv bd \pmod{n}$

In other words, in modular expressions involving addition and multiplication, individual terms can be freely exchanged with equivalent terms (under that same modulus).

Since exponentiation is just repeated multiplication, this also applies to it. Formally, if $a \equiv b \pmod{n}$ and $e$ is some integer, then

- $a^e \equiv b^e \pmod{n}$

Notice that we are always using the same modulus ($n$) in these laws. There are some scenarios where you can change between different modulo, but for this document we always use the same modulus.

### Remainders

A bit confusingly, there is a related, but different, notation that uses the word $\bmod$.

We define $a \bmod n = r$ to mean "$r$ the remainder of $a$ when divided by $n$." We enforce that remainders are never negative. Thus, $0 \leq r < a$.

In programming, we often write $a \bmod n = r$ as $a ~ \% ~ n = r$.

Notice that this gives us an alternative definition for modular equivalence:

- $a \equiv b \pmod{n}$ iff $a \bmod n = b \bmod n$

Notice that the "mod" notation is being used here in two different ways.

In the first expression, the $\pmod{n}$ is in parentheses and applies to the entire $\equiv$ operator, representing modular equivalence.

In the second expression, the $\mod{n}$ has no parentheses, and is an operator between $a$ and $n$ meaning "the remainder of $a$ when divided by $n$" (likewise for $b \bmod n$).

It should be clear that this means our laws about modular addition, multiplication, and exponentiation work with remainders as well. Namely, if $a \equiv b \pmod{n}$, $c \equiv d \pmod{n}$, and $e$ is some integer, then

- $a + c \bmod{n} = b + d \bmod{n}$
- $ac \bmod{n} = bd \bmod{n}$
- $a^e \bmod{n} = b^e \bmod{n}$

All we did was swap between modular equivalence and remainder definitions here, nothing special.

### Equivalence class notation

There is another notation we can use for modular arithmetic. Let $[a]_n$ be the set of all integers that are equivalent to $n$ modulo $a$. We call $[a]_n$ the "equivalence class" of $a$ modulo $n$. In other words, $[a]_n = \{x \in \mathbb{Z} \mid x \equiv a \pmod{n}\}$, or $[a]_n = \{a + nk \mid k \in \mathbb{Z}\}$ (think about why these definitions are the same). For example, $[1]_3 = {\ldots, -2, 1, 4, 7, \ldots}$. 

As an exercise, think about why $[1]_3 = [4]_3$. More generally, think about (or prove!) why the following statements are all equivalent:

- $x \equiv y \pmod{n}$
- $x \in [y]_n$
- $[x]_n = [y]_n$

Since $[x]_n = [y]_n$ iff $x \equiv y \pmod{n}$, our addition/multiplication/exponentiation laws from before apply to equivalence classes as well. Namely,

- $[a]_n + [b]_n = [a+b]_n$
- $[a]_n [b]_n = [ab]_n$
- $([a]_n)^d = [a^d]_n$

### "Lazy" notation

Sometimes, if we are always working under the same modulus $n$, we omit the $n$ subscript when writing equivalence classes, writing just $[a]$. Moreover, we sometimes get even more lazy and omit the parenthesis, writing just $a$ to represent the equivalence class of $a$ under $\mod n$.

When doing this, we will shorthand $[a]_n = [b]_n$ or $a \equiv b \pmod{n}$ with simply $a = b$.

Thus, writing something like the following would be reasonable:

> Working under $\mod n$, we have that $(n-1) + 1 = n = 0$.

## The protocol

In Diffe-Hellmen, we work under mod $p$ where $p$ is prime. We pick some generator $g$ with $2 \leq g < p$ (we omit discussion on how to pick a secure generator).

Client A generates a random integer $a \in [0, p)$. Client B generates $b \in [0, p)$. These are their private keys, which they will not share.

Client A then generates their public key $A := g^a \bmod{p}$. Likewise, client B generates $B := g^b \bmod{p}$. They then exchange their public keys.

Client A now has $a$ and $B$. They use these to compute:
```math
K_a := B^a = (g^b)^a = g^{ab}
```

(notice that I'm being lazy and not writing $\mod{p}$ everywhere. Since we are working under $\mod{p}$ in this whole discussion, we will let it be implicit).

Client B has $b$ and $A$, which they use to compute:

```math
K_b := A^b = (g^a)^b = g^{ab}
```

Notice that $K_a = K_b$, which we call the shared key $K$. Clients A and B can now use $K$ to encrypt and decrypt messages.

### Security

The basic idea (this is an over-simplification) of the security of this protocol is as follows. An attacker in the middle would only learn $A := g^a \bmod{p}$ and $B := g^b \bmod{p}$. There is no obvious way to use these to compute $K := g^{ab} \bmod{p}$. It would seem that you would somehow have to derive $a$ from $A$ and/or $b$ from $B$. In regular, real-numbered arithmetic, we could just compute $\log_g{A}$ to get $a$ easily. However, remember that we are working with integers under $\mod{p}$, which are "discrete" values (as opposed to real numbers, which are "continuous" values). Computing $a$ from $A := g^a \bmod{p}$ in this context is known as the "discrete logarithm problem" (DLP), which currently has no known efficient solution (outside of quantum computing, at least).
