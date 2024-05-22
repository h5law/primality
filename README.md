# primality

Primality is a Golang library for determining whether any given integer
$i\in\mathbb{Z}^+,~i\ge1$. The library provides two implementations of algorithms
for determining the primality of an arbitrarily sized integer, using the `math/big`
library.

The two methods provided are the Miller-Rabin and AKS primality test the former
being probabilistic and the latter deterministic, meaning that the Miller-Rabin
primality test doesn't guarantee 100% accuracy in certain situations. Whereas
the AKS primality test is always 100% accurate - but a lot slower on larger
integers.

## Features

### Miller-Rabin:

- Highly accurate probabilistic primality test
  - 25 rounds and force usage of base 2 recommended for near 100% accuracy
- Arbitrarily large integer support (`big.Int`)
- $\mathcal{O}(r\cdot s)$ time complexity (assuming `big.Int` operations are
  $\mathcal{O}(1)$ where $r$ is the number of repetitions and $s$ the number of
  trailing zeros on $n$, the number being tested - otherwise it is related to
  the operations of `big.Int` integers and $n$ itself.
  - As a prime must be odd $s\le7$ meaning the time complexity in its worst case
    is $\mathcal{O}(175)=\mathcal{O}(1)$ with 25 repetitions.

### AKS

- Deterministic primality test
  - Slow overall - but guarantees 100% a valid outcome
- `int` support only

## TODOs

- Improve speed of the AKS method
  - Specifically step 5 but overall it is slow
- Make the AKS method work on arbitrarily sized integers
  - use `big.Int`s over `int`s
- Determine the true time and space complexities for both methods
