# Challenge 32

**Break HMAC-SHA1 with a slightly less artificial timing leak**

## Challenge description

> Reduce the sleep in your "insecure_compare" until your previous solution breaks. (Try 5ms to start.)
> 
> Now break it again.

## Solution

When the sleep is reduced to about 5 milliseconds, the solution from challenge 31 does not always produce the correct answer by just choosing the guess corresponding to the longest request. This is because the timing analysis is more susceptible to random noise; just 5 ms of variation in the computation could skew the results.

The simplest way to handle variation in statistics is to take multiple samples and then take the mean. So, for each guess $g \in [0, 256)$ for $c_i$, we make $n$ requests to the server and take the average response time. We then choose the $g$ that corresponds to the longest average time. This improves the quality of the results significantly, even when $n = 10$.

## Future ideas

There are certainly more complex statistical analysis we could explore. For instance, rather than just considering the sample mean, we could also consider the sample variance. We could use this to quantify how confident we are that the largest sample mean significantly differs from the means for the other guesses. If the confidence is below some threshold, we could take more samples until we reach the desired confidence.
