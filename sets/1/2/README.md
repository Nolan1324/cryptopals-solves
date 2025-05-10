# Challenge 2

## Challenge description

> Write a function that takes two equal-length buffers and produces their XOR combination.
> 
> If your function works properly, then when you feed it the string:
> 
> `1c0111001f010100061a024b53535009181c`
>
> ... after hex decoding, and when XOR'd against:
> 
> `686974207468652062756c6c277320657965`
>
> ... should produce:
> 
> `746865206b696420646f6e277420706c6179`

## Overview

This challenge asks to XOR a two byte sequences together.

## XOR Basics

XOR ($\oplus$) is an operation between two bits. It returns $0$ if the two bits are the same and $1$ if the two bits are different. This gives the truth table:

| $A$ | $B$ | $A \oplus B$ |
|---|---|-------|
| $0$ | $0$ |   $0$   |
| $0$ | $1$ |   $1$   |
| $1$ | $0$ |   $1$   |
| $1$ | $1$ |   $0$   |

To XOR two sequences of bits, we just perform bit-wise XOR. For instance, XORing two bytes may look like $\verb|00000101| \oplus \verb|00001001| = \verb|00001100|$. Likewise, to XOR two sequences of bytes, we also just perform bit-wise XOR on the bytes.

XOR shares a few nice properties with traditional addition. Let $a, b, c$ be bit (or byte) sequences. Then

- $(a \oplus b) \oplus c = a \oplus (b \oplus c)$ (associativity)
- $a \oplus b = b \oplus a$ (commutativity)
- $a \oplus 0 = 0 \oplus a = a$ ($0$ acts like an identity)

However, the most interesting and useful property of XOR, which traditional addition does **not** share, is

$$
a \oplus a = 0
$$

Combining this with the other properties, this means that XORing by the same value $a$ twice is equivalent to never XORing at all. This holds true even if we XOR by other values in the middle, thanks to associativity and commutativity. For instance, $x \oplus (y \oplus a) \oplus (z \oplus a) = x \oplus y \oplus z$.

We see in later challenges that this allows us to view XOR as a "self-inverting" operator; if we XOR something by `a`, we can undo/invert this operation by simplying XORing it by `a` again.

## Implementation

In most programming languages, including Go, `a ^ b` computes the bitwise-XOR of two scalar values (bytes, ints, etc) `a` and `b`. To compute the XOR of two byte sequences, we can just apply this operating to each pair of bytes.

I implemented `func XorBytes(dst []byte, src1 []byte, src2 []byte)` in the `cipherx` package which XORs `src1` with `src2` and stores the result in `dst`. Originally I had `XorBytes` return the result as a newly allocated `[]byte`, but when completing later challenges, I found that sometimes I wanted to write the result to a buffer I had already allocated. Similar functions in the Go standard library followed this pattern of making the destination buffer the first function parameter.

### Panicking

In Go, calling `panic` produces a stack trace and exits. It is well-suited for when an unexpected/unrecoverable error occurs. For instance, if an API function is called with malformed arguments, it may be appropriate for the callee to panic. Then the developer calling the API can inspect the stack trace to figure out why the arguments they were passing were malformed and resolve the issue in their calling code.

Following this mentality, I decided to make `XorBytes` panic if `src1` and `src2` differ in length or if `dst` is shorter than `src1`/`src2`. I found this is consistent with how the Go standard library implements this function.
