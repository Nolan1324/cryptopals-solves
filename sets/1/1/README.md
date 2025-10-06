u# Challenge 1

This challenge is fairly straightforward, so I will use this space to talk about the codebase setup.

## Programming language: Go

I implemented both the challenges and my solutions in Go. I chose Go as it provides static typing (as opposed to say, Python) while still being fairly convenient and quick to implement with. Additionally, I also wanted to leverage its package system to create some structure to the codebase.

## Codebase structure

### Go module background

This entire repository represents a single Go **module**, defined by the `go.mod` file.

In Go, a module can be split into many **packages**. Each directory can contain at most one package. Each `.go` file in the directory starts with the line `package <packageName>` where `<packageName>` should be the same across all the files.

`main` is a special package name that indicates that the package builds as an executable. One file in the package should define `func main()` as the entrypoint for the executable. A module can have multiple `main` packages, they just need to be placed in seperate directories.

By default, packages in this module could be imported by other modules. However, in this project, I just want to use the packages I write internally without exporting them. Go lets you do this by placing the packages in an directory named `internal`.

When importing a package into another package or module, only symbols starting with a capital letter will be accessible from that package. For example, suppose package `a` implements functions `Foo()`, `bar()`, and defines a struct `type Data struct { Bizz int; bazz int }`. Then package `b` can access `a.Foo()` but not `a.bar()`. `b` can also access the type `Data`, but can only access its member `Bizz`, not `bazz`.

### My module

I structured this codebase to have one `main` package for each cryptopals challenge. Namely, the `sets` directory contains a numbered directory for each challenge set, which contains a numbered directory for each challenge. Each challenge directory contains a `main` package that can be ran with `go run .`, demo-ing the solution to the challenge.

I also have a few internal packages for implementing common logic. Most notably, `cipherx` (name includes an `x` to avoid confusion with the standard library `crypto/cipher` package) contains commonly re-used cipher functiosn (AES, etc) and `crack` contains functions to execute common attacks.

For example, set 1 challenge 1 is implemented in `sets/1/1/*.go` as a `main` package. It imports the internal `enc` package I implemented in `internal/enc/*.go` which exports the functions `HexToBase64`, `HexDecode`, and `Base64Encode`.

## Test cases

In Go, test cases for a package are implemented within one or many `*test.go` files. I implemented test cases for important functions in my internal packages. These tests were especially helpful when I changed my internal functions for later challenges, to ensure they did not break.

Some challenges involved implementing the "application" to be attacked. In these challenges, I would sometimes implement test cases for the application within that challenge's `main` package. These were helpful in ensuring I implemented the application correctly before attacking it.

## Representing bytes

Crytography code often needs to operate on a sequence of bytes. For instance, in typical encryption, the "plaintext" to encrypt is usually represented as a sequence of bytes, as is the "ciphertext" that it gets encrypted to. Moreover, the "key" used in the encryption is likely also a sequence of bytes.

Go has a few facilities for representing byte sequences: the primitive `[]byte` (just a slice of `byte` values) and the data type `bytes.Buffer`.

It turns out that `bytes.Buffer` is best for, well, buffers. `bytes.Buffer` primarily implements `Read` and `Write` for reading and writing to and form it. Code involving IO protocols can benefit from buffers to efficiently read and write encoded data.

The more primitive `[]byte` turns out to be fairly well suited for crytographic operations.

`[]byte` is a slice of bytes. Slices in Go are an interesting hybrid between an array view and a vector data structure. They act as array views since they can be constructed from subsets of arrays. For instance, `var x [10]byte` defines a fixed-sized array of bytes. `y := data[0:5]` constructs a slice of type `[]byte`, providing a view into that array. Since slices act like views, they essentially consist of a pointer `ptr` to the start of the data, and the length `len` of the data they view.

You can also construct slices directly with `make([]byte, len, cap)` which constructs both the slice and the underlying data array. This is what I commonly use to allocate byte arrays.

Slices store a third attribute `cap` which is the total capacity of the underlying data array that they point to. So, suppose `y := make([]byte, 5, 10)` is a byte slice with attributes `{ptr: 0xA0, len: 5, cap: 10}`. We can append a value to the underlying data array with `y = append(y, 'a')`, returning a new slice with properties `{ptr: 0xA0, len: 6, cap: 10}`. If we try to `append` when the slice is at is capacity, `append` will allocate a **new** data array with larger capacity, copy the old data to that array, and then return a slice pointing to the **new** array. For instance, if `y` has attributes `{ptr: 0xA0, len: 10, cap: 10}`, then `y = append(y, 'a')` could return `y = {ptr: 0xB0, len: 11, cap: 20}` where `0xB0` is the address of a newly allocated data array (and the old array at `0xA0` will evenetually be garbage collected). This is the vector-like behavior of a slice.

You can create a byte slice from a string with `[]byte(theString)` and create a string from a byte slice with `string(theBytes)`.

### Considerations with slices

I have found that you should exercise caution in mixing the "array view" and "vector" properties of a slice. For instance, suppose we have `x := []byte("hello world")`. The `x` slice may look something like `{ptr: 0xA0, len: 11, cap: 11}`. Then we take the slice `y := data[0:5]`.. The `y` slice object looks something like `{ptr: 0xA0, len: 5, cap: 11}`. Notice critically that the capacity is still `11`. So what happens if we do `y = append(y, 'a')`? Well, we get the slice `{ptr: 0xA0, len: 6, cap: 11}`. **Notably, the original data at `x` is now modified.** So, not only does `y[5]` equal `'a'`, but so does `x[5]` now. I think its up for debate if this is inherently "bad" behavior, but either way it requires some caution. Generally, I just try to not mix these operations too much.

## Challenge 1

Challenge 1 asks to just convert a hex (base 16) encoded string to a base 64 encoded string. I first convert hex to raw bytes, and then to base 64. In the cryptopals challenges, encrypted data is typically provided in base 16 or 64, but all the algorithms operate on the raw bytes.