# Challenge 31

**Implement and break HMAC-SHA1 with an artificial timing leak**

## Challenge description

> The psuedocode on Wikipedia should be enough. HMAC is very easy.
> 
> Using the web framework of your choosing (Sinatra, web.py, whatever), write a tiny application that has a URL that takes a "file" argument and a "signature" argument, like so:
> 
> http://localhost:9000/test?file=foo&signature=46b4ec586117154dacd49d664e5d63fdc88efb51
> Have the server generate an HMAC key, and then verify that the "signature" on incoming requests is valid for "file", using the "==" operator to compare the valid MAC for a file with the "signature" parameter (in other words, verify the HMAC the way any normal programmer would verify it).
> 
> Write a function, call it "insecure_compare", that implements the == operation by doing byte-at-a-time comparisons with early exit (ie, return false at the first non-matching byte).
> 
> In the loop for "insecure_compare", add a 50ms sleep (sleep 50ms after each byte).
> 
> Use your "insecure_compare" function to verify the HMACs on incoming requests, and test that the whole contraption works. Return a 500 if the MAC is invalid, and a 200 if it's OK.
> 
> Using the timing leak in this application, write a program that discovers the valid MAC for any file.

## Background

### HMAC

In the previous challenge, we saw how the naive implementation of MAC is vulnerable to length-extension attacks. HMAC on the other hand is a much more well-thought-out protocol. The basic idea of HMAC is to feed the output of naive MAC $H(k \mid\mid m)$ into another naive MAC, giving the form $H(k \mid\mid H(k \mid\mid m))$. Recall that length-extending a hash expression of the form $c := H(k \mid\mid m)$ requires knowing both $m$ and $c$. In HMAC, the outer hash function hides the value of the inner $H(k \mid\mid m)$ expression, preventing us from length-extending the inner hash. Moreover, not knowing $H(k \mid\mid m)$ also prevents us from length-extending the outer hash.

The HMAC full protocol also pads (XORs) the inner and outer keys with different magic numbers, giving

```math
\text{HMAC}(k, m) = H((k \oplus \text{opad}) \mid\mid H((k \oplus \text{ipad} \mid\mid m))
```

### Timing attack

Since HMAC itself is secure, this challenge explores a different type of attack called a **timing attack**. A timing attack is where analyzing how long a computation takes can reveal specific information about that computation. In this case, we are interested in a computation that compares two HMAC signatures byte-by-byte.

## Challenge setup

In this challenge, we have an HTTP server that holds a secret key $k$ for signing HMAC signatures. The user can provide a plaintext $m$ and signature $c'$ to the server. The server will then sign the message with its key to get $c := \text{HMAC}(k, m)$. Finally, it compares $c$ to $c'$ and returns the result to the user.

Since HMAC is secure, we cannot feasibly craft a $c'$ for our $m$ such that $c = c'$ without having knowledge of $k$. Thus, we will instead try to learn what value of $c$ the server computed. 

The server in this challenge uses an insecure computation to compare $c = c'$. It compares the signatures byte-at-a-time, and when it finds the first byte that differs, it returns `false` immediately. This means that, generally, the later that the first differing byte occurs, the longer the function will take to compute. In this challenge, we exaggerate this issue by having the function sleep for 50 milliseconds between each byte comparison. Thus, by inspecting how long the function took to run, we can make reasonable guesses about which byte the signatures first differ at.

## The web application

The Go standard library contains HTTP server functionality, so implementing a HTTP server is fairly straightforward.

For convenience, I launch the HTTP server in the same process that I run the client in, via a goroutine. In a real-world scenario, the HTTP server would likley be running in a process on a another computer over the network, which would introduce a plethora of time attacking challenges such as network latency, but it seems this challenge just wants us to complete basic a proof-of-concept.

Like the challenge explains, we implement an endpoint `test` that takes query parameters `file` and `signature`. It first decodes `signature` from hexadecimal to raw bytes. Our server uses SHA-1 HMAC, so the signature is 20 bytes long. Then, it computes the signature of `file` using its secret key. Finally, it compares these two signatures using the insecure compare function.

The insecure compare function is implemented as follows:

```go
// insecureCompare compares two byte arrays in an insecure manner.
// compareDuration how long each byte comparison takes.
func insecureCompare(buf1 []byte, buf2 []byte, compareDuration time.Duration) bool {
	if len(buf1) != len(buf2) {
		return false
	}
	for i := range buf1 {
		if buf1[i] != buf2[i] {
			return false
		}
		time.Sleep(compareDuration)
	}
	return true
}
```

## The attack

Suppose we know (with high confidence) the first $i$ bytes $c_{0:i}$ of the true signature of $m$. To learn the next byte $c_i$, we enumerate all possible bytes $g \in [0, 256)$ and construct our signature as $c' := c_{0:i} \mid\mid g \mid\mid (0)^{\times (20 - i - 1)}$ (the 0 padding is just to make the HMAC length 20).

For each $g$, we will pass $c'$ and $m$ to the server. If $g \neq c_i$, then since $c_{0:i} = c_{0:i}'$, the server makes $i$ byte comparisons before returning `false`. If $g = c_i$, then the server makes at least $i+1$ byte comparisons (its "at least" since it might so happen that $c'_{i+1} = 0$ and so on). Therefore, the $g$ corresponding to the request that took the longest to process on the server is most likley to be the correct byte $c_i$.

We repeat this process until we have guesses for all bytes in $c$. We can then pass $m$ and our guess for $c$ to the server to check our answer; it will return `true` if we are right.

## Secure HMAC comparison

This challenge begs the question: how can you implement secure  comparison of HMACs (or byte arrays in general)? The Go standard library implements a secure byte array comparison function in its cryptography package. `crypto/subtle` has the function `ConstantTimeCompare`:

```go
// ConstantTimeCompare returns 1 if the two slices, x and y, have equal contents
// and 0 otherwise. The time taken is a function of the length of the slices and
// is independent of the contents. If the lengths of x and y do not match it
// returns 0 immediately.
func ConstantTimeCompare(x, y []byte) int {
	if len(x) != len(y) {
		return 0
	}

	var v byte

	for i := 0; i < len(x); i++ {
		v |= x[i] ^ y[i]
	}

	return ConstantTimeByteEq(v, 0)
}
```

The key idea appears to be to remove the early return. This means that no matter if the arrays differ at the first byte or the final byte, the function will always compare all the bytes, theoretically taking the same amount of time.

I am not entirely certain what specific advantage the bitwise operations provide. My guess is that bitwise operations reveal less information than something like 

```go
if x[i] != y[i] {
    equal = false
}
```

because the branching could reveal some timing information.

<!-- We discussed in early challenges how $x = y$ iff $x \oplus y = 0$. -->

