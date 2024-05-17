package primality

import (
	"math/big"
	"math/rand"
	"slices"
	"time"
)

// SieveOfEratosthenes is an implementation of the ancient sieve of Eratosthenes
// to generate a list of primes up to the provided integer.
func SieveOfEratosthenes(n uint64) []uint64 {
	if n < 2 {
		return nil // 1 is not prime
	}
	// Initialise a slice to store the primes we generate
	primes := make([]uint64, 0, n)
	// Initialise a boolean slice with all false values
	b := make([]bool, n+1)
	b[0], b[1] = true, true
	// Set all multiples of numbers up to sqrt(n) to true
	for i := uint64(2); i*i <= n; i++ {
		if b[i] {
			continue
		}
		for j := i * i; j <= n; j += i {
			b[j] = true
		}
	}
	// Primes are remaining false values indexes
	for i, v := range b {
		if !v {
			primes = append(primes, uint64(i))
		}
	}
	return primes
}

var zero = big.NewInt(0)
var one = big.NewInt(1)
var two = big.NewInt(2)

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
	low := make([]uint64, lo)
	high := make([]uint64, hi)
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

// MillerRabin is an implementation of the Miller-Rabin primality test, it is
// a probabilistic method - and as such using 25 repetitions/rounds is
// recommended as it ensures a higher probability of accuracy of the result.
// The test is 100% accurate up to 2^64 and some values higher - forcing the
// use of base 2 when randomly generating bases for the
func MillerRabin(n *big.Int, reps int, force2 bool) bool {
	// n must and odd integer > 2
	one := big.NewInt(1)
	even := new(big.Int)
	even = even.And(n, one)
	lt := n.Cmp(two)
	if lt < 0 || (lt != 0 && even.Cmp(zero) == 0) {
		return false
	}
	// Get prime bitmask for primes between 2-2^16
	if n.Cmp(big.NewInt(2<<16)) <= 0 {
		primeBitmask := primemask(2, 2<<16)
		// Check that primeBitmask&(1<<n) != 0 -> n == prime
		shift := big.NewInt(1)
		shift = shift.Lsh(shift, uint(n.Uint64()))
		prime := new(big.Int)
		prime = prime.And(primeBitmask, shift)
		return prime.Cmp(zero) != 0
	}
	// Create offset values
	nm1 := big.NewInt(-1)
	nm1 = nm1.Add(n, nm1)
	nm4 := big.NewInt(-4)
	nm4 = nm4.Add(n, nm4)
	// s > 0 and d > 0 such that n-1 = 2^s(d)
	s := nm1.TrailingZeroBits() // get highest power of 2 from n-1
	d := new(big.Int)
	d = d.Rsh(nm1, s) // get odd multiplier such that (n-1)/2^s=d
	// Get random source
	src := rand.NewSource(time.Now().UTC().UnixNano())
	rand := rand.New(src)
	y := new(big.Int)
	a := new(big.Int)
	for i := 0; i < reps; i++ {
		// n is always a probable prime to base 1 and n-1
		a = a.Rand(rand, nm4).Add(a, big.NewInt(int64(2))) // random(2, n-2)
		if force2 && i == reps-1 {
			a = two // force base of 2 on final repetition if required
		}
		x := bigPowMod(a, d, n)
		for j := uint(0); j < s; j++ {
			y = bigPowMod(x, two, n)
			if y.Cmp(one) == 0 && x.Cmp(one) != 0 && x.Cmp(nm1) != 0 {
				return false // nontrivial square root of 1 (mod n)
			}
			x = x.Set(y)
		}
		if y.Cmp(one) != 0 {
			return false
		}
	}
	return true
}
