# Challenge 7

**AES in ECB mode**

## Overview

This challenge asks to implement the AES ECB mode. In this mode, each block of the plaintext is independently encrypted with AES. The block size is typically 16 bytes.

## Solution

The Go standard library package `crypto/aes` provides `aes.NewCipher` to create an AES cipher. This cipher can be used to encrypt/decrypt individual blocks. I use this to encrypt/decrypt each block independently to implement AES ECB mode.

