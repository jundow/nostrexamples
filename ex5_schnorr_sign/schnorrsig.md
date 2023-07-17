
## Algebraic representation of Schnorr sign and verification in BIP0340

This short note aims to learn how Schnorr signature is represented as a set of formalisms.

The formalisms will lack of programming-level details however, I believe this will still helpful for people to understand how the signing process works and revisit the process when coding the details.

This article hide the generation process of random numbers so called nonce to simplify the explanation here. I believe this can be explained separately.


### Notation

- $\bold{G}$ : Bold - An element of an elliptic curve (Here, it's Secp256k1.)
- $\bold{G}.x$ or $\bold{G}.y$ are x or y coordinate of the elliptic curve $\bold{G}$. Each number is integer.
- $d$ : Italic - An integer
- $a[32]$ : With [num] - Byte array with a fixed length
- $m[]$ : With[] - Byte array with an arbitral length
- $d \cdot \bold{G} = \bold{P} $ - Scalar multiplication of an elliptic curve element returns an element of the elliptic curve.
- $int(a[32]) = b$ - Function $int()$ returns the corresponding integer of a fixed byte array.
- $ Rand_{BIP0340}(a[32], d, \bold{P},m[])$ - The function returns the random number from $a[32]$; auxiliary random data, d; integer, P; a point in elliptic curve and m; message to sign.

### Signature

#### Input

- sk[32] - as a fixed length byte array of secret key
- message[] - as message to sign
- a[32] - Auxiliary generated 32 bytes of random data  

#### Constant

- $\bold{G}$ Generator of Secp256k1 elliptic curve
- $p_r$ = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F; Prime number assigned to Secp256k1 elliptic curve
- $n$ = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141; Order of the generator $\bold{G}$ in scalar multiplication

#### Signing Process

- Let $d' = int(sk[32])$
- Let $\bold{P}=d'\cdot\bold{G} $
- if $\bold{P}.y$ is even, let $d=d'$, otherwise ,$d=n-d'$
- $k'= int(Rand_{BIP0340}(a[32], d, \bold{P}, message )) \mod n$
- $\bold{R}=k'\bold{G}$
- if $\bold{R}.y$ is even, let $k=k'$, otherwise ,$k=n-k'$
- Convert necessary information to hash by letting,
   $e$ = int( sha256( 
    sha256( byte("BIP0340/challenge") ) || sha256( byte("BIP0340/challenge" )) ||
    byte($\bold{R}.x$) || byte($\bold{P}.x$) || m ) ) $\mod n$ 
- Return byte($\bold{R}.x$) || byte ($k+ed$) as 64 bytes of byte array as sig[64]

#### Verification

#### Input

- pk[32] - as the public key derived from the signer's secret key
- m[] - the same message used to generate the signature
- sig[64] - the signature

#### Verification

- Calculate the point of elliptic curve using $int(pk[32])$
    - Let $\bold{P}_v$ be the point of the public key in the elliptic curve for verification
    - $\bold{P}_v.x = int(pk[32])$
    - Calculate $\bold{P}.y$ from the equation of the elliptic curve $\bold{P}_v.y^2 = \bold{P}_v.x^2 + 7 \mod p_r$
    - The equation has the two possible, even and odd. Take the even one for $\bold{P}_v$.
-  $r[32] = sig[0:31]$ and $s = sig[32:63]$ 
- Re-calculate hash by letting,
   $e$ = int( sha256( 
    sha256( byte("BIP0340/challenge") ) || sha256( byte("BIP0340/challenge" )) ||
    byte( byte($r$) || byte($\bold{P}_v.x$) || m ) ) $\mod n$ 
- Let $\bold{R}_v = s\cdot\bold{G} - e\cdot\bold{P}_v$ 
- Compare $\bold{R}_v.x$ and $r$. If they are the same, verification is succeeded otherwise failed. 
  - By calculating $\bold{R}_v$, 
    $\bold{R}_v = s\cdot\bold{G} - e\cdot\bold{P}_v = (k+ed)\cdot\bold{G} - e\cdot\bold{P}_v = k \cdot \bold{G}$
    - Since $k = \pm k'$ , + for $k'\cdot \bold{G}$ is even and - for odd. So, $d\cdot\bold{G}$ is always even . This means   $e\cdot(d\cdot\bold{G}) - e\cdot\bold{P}_v = \bold{O}$, because $\bold{P}_v.y$ was selected for $\bold{P}_v$ to be even. 
  - If the sign is valid, the x coordinates of both "$r = k'\cdot\bold{G}$" and "$\bold{R}_v=k\cdot\bold{G}$" is equal.

### Generating Random Number

(Under construction)


