# Challenge 30

**Break an MD4 keyed MAC using length extension**

## Challenge description

> Break an MD4 keyed MAC using length extension
>
> Second verse, same as the first, but use MD4 instead of SHA-1. 
>
> Having done this attack once against SHA-1, the MD4 variant should take much less time; mostly just the time you'll spend Googling for an implementation of MD4.

## Solution

As the challenge description states, this challenge is very similar to the previous.

MD4 is a very similar hash algorithm to SHA-1. MD4 computes a 16-byte / 128-bit digest (rather than 20 bytes), represented by 4 32-bit registers (rather than 5). The block size is the same as SHA-1 (64 bytes / 512 bits). Like SHA-1, the heart of the MD4 algorithm is the block processing function, which I once again use the Go standard library code for.

My custom SHA-1 and MD4 classes both implement the `hash.Hash` interface from the standard library, making it easy to switch between them. Moreover, I implemented my `hashx.HashUtils` from the previous challenge for MD4. Since my `ExtendMac` takes in a `hashx.HashUtils` as input, all I had to do for this challenge was pass in the MD4 `HashUtils` implementation instead of the SHA-1 implementation.

