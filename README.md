# aks

aks is a Golang library that provides an implementation of the AKS primality test
to deterministically determine whether an integer is prime or not. It works using
`math/big` so can work on arbitrarily large integers.

AKS is deterministic and thus provides a 100% guarantee on the primality of an
integer as opposed to other probabilistic methods such as the Miller-Rabin
primality test.
