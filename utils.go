package primality

import (
	"math/big"
	"slices"
)

// Return new slice containing only the elements in a that are not in b.
//
//	d = {x : x ∈ a and x ∉ b}
func difference(a, b []uint64) []uint64 {
	diff := make([]uint64, 0, len(a))
	for _, v := range a {
		if found := slices.Contains(b, v); !found {
			diff = append(diff, v)
		}
	}
	return diff
}

// Generate a bitmask for primes from lo->hi, the bitmask will index the primes
// at their relative bit indexes for quick comparison and evaluation of primes.
func primemask(lo, hi uint64) *big.Int {
	var low, high []uint64
	// Get primes 2->hi
	high = SieveOfEratosthenes(hi)
	// Only get lower limit primes if lo > 2
	if lo > 2 {
		low = SieveOfEratosthenes(lo)
	}
	// Use only difference between high and low slices
	diff := difference(high, low)
	// Create the mask by performing mask |= 1<<prime
	mask := new(big.Int)
	one := big.NewInt(1)
	for i := 0; i < len(diff); i++ {
		prime := new(big.Int)
		// Left shift
		prime = prime.Lsh(one, uint(diff[i]))
		mask = mask.Or(mask, prime)
	}
	return mask
}

// bigPowMod uses the binary exponentiation of the exponent provided to
// compute the base raised to the exponent mod the provided modulus.
func bigPowMod(b, e, m *big.Int) *big.Int {
	// Special case
	if m.Cmp(one) == 0 {
		return zero
	}
	// Copy arguments
	base := new(big.Int)
	base = base.Set(b)
	exp := new(big.Int)
	exp = exp.Set(e)
	mod := new(big.Int)
	mod = mod.Set(m)
	// Initialise remainder
	r := big.NewInt(int64(1))
	// Use the binary exponentiation of e to successively calculate b^e (mod m)
	base = base.Mod(base, mod)
	em := new(big.Int)
	for exp.Cmp(zero) > 0 {
		em = em.Mod(exp, two)
		if em.Cmp(one) == 0 {
			r = r.Mul(r, base).Mod(r, mod)
		}
		// Divide exponent by 2
		exp = exp.Rsh(exp, uint(1))
		// Repeatedly square b to achieve b=b^(2^i) (mod m)
		base = base.Mul(base, base).Mod(base, mod)
	}
	return r
}
