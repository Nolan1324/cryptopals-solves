# Challenge 8

**Detect AES in ECB mode**

## Overview

This challenge provides a list of ciphertexts and asks to detect which one was encrypted with AES in ECB mode.

## Solution

In AES ECB mode, each block is encrypted by the same key and algorithm. Thus, if two blocks are equal in the plaintext, then they are also equal in the ciphertext. Thus, if the ciphertext has duplicate, its likely it was encrypted with AES ECB. In this challenge, checking for duplicate blocks allows us to find which ciphertext was encrypted with AES ECB.

## Commentary

Notably, there are many real scenarios where the plaintext could have duplicate blocks. For instance, if the plaintext represents a page of memory, there could be a region of many "0" bytes, resulting in multiple all-zero blocks.
