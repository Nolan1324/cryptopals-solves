# Challenge 9

**Implement PKCS#7 padding**

## Challenge description

> A block cipher transforms a fixed-sized block (usually 8 or 16 bytes) of plaintext into ciphertext. But we almost never want to transform a single > block; we encrypt irregularly-sized messages.
> 
> One way we account for irregularly-sized messages is by padding, creating a plaintext that is an even multiple of the blocksize. The most popular > padding scheme is called PKCS#7.
> 
> So: pad any block to a specific block length, by appending the number of bytes of padding to the end of the block. For instance,
> 
> `"YELLOW SUBMARINE"`
> 
> ... padded to 20 bytes would be:
> 
> `"YELLOW SUBMARINE\x04\x04\x04\x04"`

## Solution

This challenge is fairly straightfoward. I implemented a `AddPcks7Padding` function to add the padding.

### Semantics of `AddPcks7Padding`

I debated a bit what the semantics of this function should be. Should it return a new byte buffer with the padded text? Or should it write its output to a buffer parameter? These options are easy to interpret, but both require unnecessarily copying the entire input, even though we are just adding a constant amount of padding.

To avoid copying the input, one option is to have `AddPcks7Padding` take in a pointer to the input slice (`*[]byte`), append to the slice, and then modify the input slice. Another option is to keep the input as type `[]byte`, append to the slice, and then return the new slice as a `[]byte`. The latter option has the exact same semantics as `append` in the standard library, so I chose that option.
Like `append`, the best way to use this function is like `buf = AddPcks7Padding(buf, 16)`. I gave `AddPcks7Padding` a docstring based on the existing one for `append` to communicate these semantics clearly.
