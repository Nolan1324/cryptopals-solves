# Challenge 3

## Challenge Description

> The hex encoded string:
> 
> `1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b336`
>
> ... has been XOR'd against a single character. Find the key, decrypt the message.
> 
> You can do this by hand. But don't: write code to do it for you.
> 
> How? Devise some method for "scoring" a piece of English > plaintext. Character frequency is a good metric. Evaluate each > output and choose the one with the best score.

## Overview

This challenges provides a string that has been XOR'd by a single byte. We have to find the "key" (the single byte) that was used to encrypt the message, and then decrypt it.

## Solution

### Brute-forcing the key

Thanks to self-inverting property of XOR, once we find the "key" that was used to encrypt the string, we can just XOR the string by the key again to decrypt it. Since the key is just a single byte, this means we can just exhaustively try decrypting using all 256 possible values for the key byte, and then see which decrypted string looks the most "correct."

With only 256 options, it is feasible to just look through all the decrypted strings and see which is correct. However, in later challenges, this becomes less feasible. Thus, it would be useful to have a "score" metric to judge how "correct" a candidate decrypted string looks.

### Score function

The score function ultimately depends on what we expect the plaintext to look like. In this challenge, we expect the plaintext to be an English message composed of ASCII characters. Thus, the challenge recommends using character frequency as a score metric. For instance, the plaintext is likely to have more occurences of the character "e" than the character "q".

This idea can be formalized as a "histogram" that tracks how often each character occurs in the string. If we have an "expected" histogram based on the English language and an "observed" histogram calculated from the plaintext guess, we can compute how similar these histograms are to score the guess.

#### Histogram representation

I played around with a few different ways to represent histograms. Initially, I only built histograms on alphabetic characters. However, this made it difficult to score (probably bad) guesses that were mostly composed of non-alphabetic characters. I tried to ignore guesses that had "non-word looking" character, but this inadvertly caused the correct guess to sometimes be ignored in later challenges, such as if it contained the `\n` character. Ultimately, I found the best solution was to include all 128 ASCII characters in the histogram. Guesses with non-ASCII bytes (value >= 128) would be ignored. This allowed me to nicely represent the fact that non-word characters are less common but still possible. Now, a histogram is formally a 128-dimensional vector.

#### Expected histogram computation

To compute the expected histogram, I used the `nltk` package in Python to download an English text corpus, and then simply computed the frequency of each ASCII character. This approach was able to calculate frequencies for non-word characters as well, such as `'` and `\n`. Naturally, many valid ASCII characters, like `\0`, had frequencies of 0.

We can compute the observed histogram of a plaintext guess in the same way (just counting characters). The sum of the elements of the observed histogram is the length of the plaintext.

#### Similarity score

To compare histograms, I first scale the expected histogram so that its elements sum to 1, so that each element represents a percentage. I then multiply the estimated histogram by the length of the plaintext. Now, the expected and observed histograms have the same sum.

I then used cosine similarity to compare the two scaled histogram vectors, which appeared to work fairly well.

#### Cosine similarity vs Chi-squared test

I believe that the more statistically sound comparison method would be the Chi-squared test. However, I encountered some issues with it for this problem. For instance, the formula for Chi-squared $\sum_{i=0}^n (O_i - E_i)^2 / E_i$ causes the p-value to explode to infinity if a character that was not expected to occur ($E_i = 0$) does occur ($O_i > 0$). However, in our problem, this very well could be possible; for instance, `\0` might occasionally occur in the real plaintext despite it never appearing in the corpus. Thus, we would not want this to ruin the p-value. Cosine similarity, on the other hand, would still somewhat penalize the score in this scenario, but if the rest of the plaintext looked good then the score could still be high. 

Thus, I decided to use cosine similarity for the time being, since it seemed to produce good results.
