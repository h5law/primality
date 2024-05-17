# primality

Primality is a Golang library for determining whether any given integer $i\in
\mathbb{Z}+,~i\ge1$. The library provides two implementations of algorithms for
determining the primality of an arbitrarily sized integer, using the `math/big`
library.

The two methods provided are the Miller-Rabin and AKS primality test the former
being probabilistic and the latter deterministic, meaning that the Miller-Rabin
primality test doesn't guarantee 100% accuracy in certain situations. Whereas
the AKS primality test is alwasy 100% accurate - but a lot slower on larger
integers.
