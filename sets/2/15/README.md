# Challenge 15

**PKCS#7 padding validation**

## Challenge Description

> Write a function that takes a plaintext, determines if it has valid PKCS#7 padding, and strips the padding off.
>
> The string:
>
> ```
> "ICE ICE BABY\x04\x04\x04\x04"
> ```
> 
> ... has valid padding, and produces the result "ICE ICE BABY".
> 
> The string:
> 
> ```
> "ICE ICE BABY\x05\x05\x05\x05"
> ```
> 
> ... does not have valid padding, nor does:
> 
> ```
> "ICE ICE BABY\x01\x02\x03\x04"
> ```
> 
> If you are writing in a language with exceptions, like Python or > Ruby, make your function throw an exception on bad padding.
> 
> Crypto nerds know where we're going with this. Bear with us.

## Solution

This challenge is fairly straightforward, though we have to be careful of a few edge cases, depending on the implementation. My implementation first reads the padding length `padLen := buf[len(buf)-1]` from the final byte, and then checks that the bytes in `buf[len(buf)-padLen:]` all equal `padLen`. However, we first need to check that `len(buf) > 0` and `padLen >= len(buf)` to avoid out-of-bounds panics. More subtely, we also need to check that `padLen != 0`, because currently this algorithm would return "valid" if `padLen == 0` even though `\x00` is **not** valid PCKS7 padding. Unit test cases are helpful in verifying these edge cases.

### Returning errors

The convention in Go for a function returning an recoverable error (rather than an unrecoverable error, for we typically just `panic`) is to return an object of type `error`, which equals `nil` iff no error occured. If the function also needs to return a value of type `valType` when there is no error, the convention is to return a tuple of type `(valType, error)`. The caller then calls `val, err := func()` and is responsible for checking `err != nil` before trying to use `val`. My remove PCKS7 padding function has signature `RemovePcks7Padding([]byte) ([]byte, error)`. If the padding is valid, it returns a slice of the input that excludes the padding, along with a `nil` error. Otherwise, it returns a non-`nil` error.
