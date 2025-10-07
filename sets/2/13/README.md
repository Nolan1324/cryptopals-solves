# Challenge 12

**ECB cut-and-paste**

## Challenge description

Write a k=v parsing routine, as if for a structured cookie. The routine should take:

```
foo=bar&baz=qux&zap=zazzle
```

... and produce:

```json
{
  foo: 'bar',
  baz: 'qux',
  zap: 'zazzle'
}
```

(you know, the object; I don't care if you convert it to JSON).

Now write a function that encodes a user profile in that format, given an email address. You should have something like:

```
profile_for("foo@bar.com")
```

... and it should produce:

```json
{
  email: 'foo@bar.com',
  uid: 10,
  role: 'user'
}
```

... encoded as:

```
email=foo@bar.com&uid=10&role=user
```

Your "profile_for" function should not allow encoding metacharacters (& and =). Eat them, quote them, whatever you want to do, but don't let people set their email address to "foo@bar.com&role=admin".

Now, two more easy functions. Generate a random AES key, then:

Encrypt the encoded user profile under the key; "provide" that to the "attacker".
Decrypt the encoded user profile and parse it.
Using only the user input to profile_for() (as an oracle to generate "valid" ciphertexts) and the ciphertexts themselves, make a role=admin profile.

## Context: message signing

This is the first challenge where encryption is used to (attempt to) achieve **integrity**, rather than **confidentiality**.

**Confidentiality** is when the defender wants to keep the plaintext confidential; the attacker wants to crack the plaintext from the ciphertext, as we did in previous challenges.

**Integrity** is when a receiver wants to verify that a message came untampered from a certain sender. In this challenge, the sender and receiver are both the "application." When the user requests a profile with `profile_for`, the application creates a profile string and encrypts it with its secret AES key $k$ to produce a ciphertext, which it returns to the user. **This is not done to hide the plaintext profile from the user.** Rather, the application is essentially "signing" the profile with $k$ to show that it authored it. In the future, the user may pass the ciphertext (the signed profile) back to the application; for instance, to request some administrator action. The application will decrypt the profile with $k$ and then check if `role=admin` is set in the decrypted profile. **In order for the ciphertext to decrypt to a valid profile, the profile must have been originally encrypted it with $k$, which (ideally) only the application is capable of doing.** The goal of the attacker is then to fabricate a ciphertext that when decrypted with $k$, produces a valid profile with `role=admin` set.

<!-- From the application's point of view, this ciphertext would -->

<!-- This would allow the attacker to pretend as though they previously obtained a signed profile with `role=admin` from the application and thus preform . -->

## The application

Unfortunately, the challenge was a bit unclear about the specifics of how the application should be implemented. Does `uid` always equal `10`, or does it monotonically increase, or is it random? Does the URL encoder escape all characters, or just `&` and `=`? Does the URL parser throw an error if it detects excess parameters other than `email`, `uid`, and `role`? Somewhat arbitrarily, I decided that `uid` is always 10, the URL only escapes `&` and `=`, and excess parameters are not allowed. I found afterwords that other people used different assumptions, but I do not think that the exact assumptions meaningfully impacted the solution or lesson of the challenge.

I wrote a few test cases to ensure that the application I implemented behaved as expected, such as forbidding `&` and `=` characters in emails.

## Solution

The key insight of this challenge is that since ECB encrypts each block independently and deterministically, we can stitch together ciphertext blocks returned from the application by `profile_for` to construct a new valid ciphertext. The challenge then becomes to determine what strings we should call `profile_for` with to get the necessary ciphertext blocks to construct a valid "admin" ciphertext.

Since my application does not allow excess parameters in the profile, we must construct a ciphertext that decrypts to `email=<some_email>&uid=10&role=admin[PKCS7 PAD]`. We can get the ciphertext for `email=<some_email>&uid=10&role=user[PKCS7 PAD]` by just calling `profile_for(<some_email>)`. This contains everything we need except for `admin` at the end.

We could get the ciphertext for some string containing `admin` by including `admin` in our email, but we cannot directly get the ciphertext for `role=admin` in this way since `profile_for` forbids the `=` character in the email. Thus, to get the `role=` portion, we will instead call `profile_for(<some_email>)` in such a way so that `role=` in `email=<some_email>&uid=10&role=user[PKCS7 PAD]` lies at the end of some block. That way, we can just exclude the ciphertext block for `user[PKCS7 PAD]` when we pick and stitch the ciphertext blocks later.

It happens that if we call `profile_for` with a 13-byte email, like `profile_for("mmail@foo.bar")`, we get the ciphertext for the following (new lines denote 16-byte block boundaries):

```
email=mmail@foo.
bar&uid=10&role=
user[PKCS7 PAD ]
```

We can then take just the first two blocks, giving us the ciphertext for `email=mmail@foo.bar&uid=10&role=`.

Now we just need the ciphertext for the block `admin[PKCS7 PAD]` (the PKCS7 padding is needed to make the string valid). If we call `profile_for("mmail@foo.admin\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B")`, we get ciphertext for

```
email=mmail@foo. 
admin[PKCS7 pad]
&uid=10&role=use
r[PKCS7 pad    ]
```

We can then take just the second block and append it to the two ciphertext blocks we crafted earlier, giving a final valid ciphertext for `email=mmail@foo.bar&uid=10&role=admin[PKCS7 PAD]`.