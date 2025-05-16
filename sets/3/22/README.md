# Challenge 22

**Crack an MT19937 seed**

## Challenge description

> Make sure your MT19937 accepts an integer seed value. Test it (verify that you're getting the same sequence of outputs given a seed).
> 
> Write a routine that performs the following operation:
> 
> - Wait a random number of seconds between, I don't know, 40 and 1000.
> - Seeds the RNG with the current Unix timestamp
> - Waits a random number of seconds again.
> - Returns the first 32 bit output of the RNG.
> - You get the idea. Go get coffee while it runs. Or just simulate the passage of time, although you're missing some of the fun of this exercise if you do that.
> 
> From the 32 bit RNG output, discover the seed.

## Solution

In this challenge, we are given the first 32-bit output from a MT19937 PRNG and want to determine what seed was used. 

To verify if a seed guess is correct, we can just seed our own MT19937 instance with the guess, generate its first output, and compare it to the output we were given. This works because MT19937 generates the sequence deterministicly from the seed. _Technically_, it is possible that two different seeds could produce the same 32-bit output, but that is incredibly unlikely (probability $1/2^{32}$).

We know that the PRNG was seeded with the current Unix timestamp (in seconds) sometime in the last few minutes. To discover the seed, we can start with the current Unix timestamp as our guess, and then decrement it by 1 until we find a seed that produces the same output. Since the PRNG was seeded in the last few minutes, we will only need to try about 500 seeds.