# primality

Primality is a Golang library for determining whether any given integer $i\in
\mathbb{Z}^+,~i\ge1$. The library provides two implementations of algorithms for
determining the primality of an arbitrarily sized integer, using the `math/big`
library.

The two methods provided are the Miller-Rabin and AKS primality test the former
being probabilistic and the latter deterministic, meaning that the Miller-Rabin
primality test doesn't guarantee 100% accuracy in certain situations. Whereas
the AKS primality test is always 100% accurate - but a lot slower on larger
integers.

## Features

### Miller-Rabin:

- Highly accurate probabilistic primality test
  - 25 rounds and force usage of base 2 recommended
- Arbitrarily large integer support (`big.Int`)

### AKS

- Deterministic primality test
  - Slow on larger integers
- `int` support only

## TODOs

- Fix false positives with AKS method
- Add unit tests
- Improve speed of the AKS method
  - Specifically step 5
- Make the AKS method work on arbitrarily sized integers
  - use `big.Int`s over `int`s
