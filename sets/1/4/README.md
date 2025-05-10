# Challenge 4

**Detect single-character XOR**

## Challenge Description

> **Detect single-character XOR**
>  
> One of the 60-character strings in this file has been encrypted by single-character XOR.
>
> Find it.
>
> (Your code from #3 should help.)

## Solution

For each string, we can compute the score of the best plaintext guess as in Challenge 3. We can then pick best guess out of these 60 strings with the highest score as our final guess.
