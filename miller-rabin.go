package primality

import (
	"math/big"
	"math/rand"
	"time"
)

var zero = big.NewInt(0)
var one = big.NewInt(1)
var two = big.NewInt(2)
var primeMask = primemask(2, 2<<16)

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
		// Check that primeMask&(1<<n) != 0 -> n == prime
		shift := big.NewInt(1)
		shift = shift.Lsh(shift, uint(n.Uint64()))
		prime := new(big.Int)
		prime = prime.And(primeMask, shift)
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
