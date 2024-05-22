package primality

import (
	"fmt"
	"math"
	"slices"
)

// fastPower computes b^e fast
func fastPower(b, e int) int {
	r := 1
	for e > 0 {
		if e%2 == 1 {
			r *= b
		}
		b *= b
		e >>= 1
	}
	return r
}

// fastPowerMod computes b^e mod m fast
func fastPowerMod(b, e, m int) int {
	r := 1
	if e < 0 {
		e = -e
	}
	for e > 0 {
		if e%2 == 1 {
			r = (r * b) % m
		}
		b = (b * b) % m
		e >>= 1
	}
	return r
}

// Check n != a^b for a,b > 1, returning true if it is otherwise false if not.
func basePowerCheck(n int) bool {
	var bMax = int(math.Log2(float64(n))) + 1
	for b := 2; b <= bMax; b++ {
		var aMin = 2
		var aMax = int(math.Pow(float64(2), float64(64)/float64(b))) - 1
		if fastPower(b, aMin) == n {
			return true
		}
		if fastPower(b, aMax) == n {
			return true
		}
		// binary search for a base a for exponent b
		for aMax-aMin != 1 {
			var aMid = (aMin + aMax) / 2
			var nMid = fastPower(aMid, b)
			if nMid == n {
				return true
			} else if nMid < n {
				aMin = aMid
			} else { // nMid > n
				aMax = aMid
			}
		}
	}
	return false
}

// orderMod finds the order of coprimes a and r such that a^i % r = 1,
// where i is the order found
func orderMod(a, r int) int {
	// Check if a and r are coprime
	if gcd(a, r) != 1 {
		return 0
	}
	// Find the order i such that a^i % r = 1
	i := 0
	for {
		fmt.Println(i)
		if fastPowerMod(a, i, r) == 1 {
			return i
		}
		i++
	}
}

// Find the smallest order such that Ord_r(n) > (log_2(n))^2
func findCorrectOrder(n int) int {
	fn := float64(n)
	lower := int(math.Floor(math.Log2(fn)*math.Log2(fn) + 0.5))
	upperR := max(3, int(math.Floor(math.Pow(math.Log2(fn), 5)+0.5)))
	for r := 1; r < upperR; r++ {
		if gcd(r, n) != 1 {
			continue
		}
		k := 2
		for {
			if fastPowerMod(r, k, n) == 1 {
				break
			}
			k++
		}
		if r > lower && r != n {
			return r - 1
		}
	}
	return 0
}

// gcd finds the greatest common divisor of a and b
func gcd(a, b int) int {
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}
	return a
}

// gcdChecker returns false when the gcd of all values [2, r] and n aren't
// coprime otherwise it returns true - meaning the n is a composite number
func gcdChecker(n, r int) bool {
	// Check all the gcd values in the interval [2, Ord_r(n)]
	for i := r; i > 1; i-- {
		d := gcd(i, n)
		// Check if i and n are coprime
		if d > 1 && d < n {
			return true // n is a composite number
		}
	}
	return false
}

// eulerTotient is an implementation of Euler's Totient function, which counts
// the number of relatively prime numbers to n, in the interval [0, sqrt(n)]
func eulerTotient(n int) int {
	res := 0
	// Check up to the square root of n for numbers coprime with n
	for i := 1; i <= n; i++ {
		g := gcd(i, n)
		if g == 1 {
			res++
		}
	}
	return res
}

// polynomialMod does a term-wise reduction on the coefficients of a polymoial p
// using modulus m where p is the coefficients of some (x+a)^e in ascending order
// of the powers of x from x = 0
func polynomialMod(p []int, m int) []int {
	r := make([]int, len(p))
	for i, x := range p {
		mod := x % m
		if mod < 0 {
			mod += x
		}
		if mod != 0 {
			r[i] = mod
		}
	}
	stripTrailingZeros(&r)
	return r
}

// polynomialExpansion finds the coefficients of the polynomial expansion of
// (x+a)^e and returns the coefficients in ascending powers of x from x = 0
func polynomialExpansion(e, a int) []int {
	c := make([]int, e+1)
	c[0], c[e] = 1, 1
	for i := 0; i < int(e/2); i++ {
		x := c[i] * (e - i) / (i + 1)
		c[i+1], c[e-1-i] = x, x
	}
	if a != 0 {
		for i := range c {
			c[i] *= fastPower(a, e-i)
		}
	}
	return c
}

func stripTrailingZeros(p *[]int) {
	for i := len(*p) - 1; i >= 0; i-- {
		if (*p)[i] != 0 {
			*p = (*p)[:i+1]
			return
		}
	}
	*p = nil
}

func degree(p []int) int {
	for i := len(p) - 1; i >= 0; i-- {
		if i != 0 {
			return i
		}
	}
	return -1
}

func polynomialMultiplyScalar(p []int, n int) []int {
	res := make([]int, len(p))
	for i := range p {
		res[i] = p[i] * n
	}
	return res
}

func polynomialRemainder(p1, p2 []int) []int {
	dp1 := degree(p1)
	dp2 := degree(p2)
	if dp2 < 0 {
		return nil
	}
	if dp1 < dp2 {
		return p1
	}
	q := make([]int, dp1)
	for dp1 >= dp2 {
		d := make([]int, dp1+1)
		copy(d[dp1-dp2:], p2)
		q[dp1-dp2] = p1[dp1] / d[len(d)-1]
		d = polynomialMultiplyScalar(d, q[dp1-dp2])
		p1 = polynomialSubtraction(p1, d)
		dp1 = degree(p1)
	}
	return p1
}

// polynomialModRemainder finds the remainder of polynomials p1/p2 and does a
// term-wise reduction modulo m on the result, returning a slice of coefficients
// for a polynomial in ascending order of x.
func polynomialModRemainder(p1, p2 []int, m int) []int {
	return polynomialMod(polynomialRemainder(p1, p2), m)
}

// polynomialSubtraction subtracts p1 from p2 in a term-wise fashion.
// The function orders the polynomials longest first prior to subtraction.
func polynomialSubtraction(p1, p2 []int) []int {
	res := make([]int, len(p1))
	i := 0
	for len(p2) > 0 && len(p1) > 0 {
		res[i] = p1[0] - p2[0]
		p1 = p1[1:]
		p2 = p2[1:]
		i++
	}
	stripTrailingZeros(&res)
	if len(p1) > 0 {
		res = append(res, p1...)
	}
	if len(p2) > 0 {
		for _, x := range p2 {
			res = append(res, -x)
		}
	}
	stripTrailingZeros(&res)
	return res
}

// AKS is an implementation of the AKS deterministic primality test.
// Step 5 takes up the majority of the time and as such results in a slow test.
func AKS(n int) bool {
	// Initial check that n > 1 and not even (expect for 2 itself)
	if n < 2 || (n > 2 && n&1 == 0) {
		return false
	}

	// Step 1
	composite := basePowerCheck(n)
	if composite {
		return false
	}

	// Step 2
	r := findCorrectOrder(n)

	// Step 3
	composite = gcdChecker(n, r)
	if composite {
		return false
	}

	// Step 4
	if n <= 5690034 && n <= r {
		return true
	}

	// Step 5
	maxA := int(
		math.Floor(
			(math.Log2(float64(n)) * math.Sqrt(float64(eulerTotient(r)))) + 0.5,
		),
	)
	xna := polynomialExpansion(n, 0)
	xr1 := polynomialExpansion(r, 0)
	var xa, remA, remB []int // stop throwing away each loop
	for a := 1; a <= maxA; a++ {
		xna[0] += a
		xa = polynomialExpansion(n, a)
		remA = polynomialModRemainder(xa, xr1, n)
		remB = polynomialRemainder(xna, xr1)
		// fmt.Println(remA, remB)
		if slices.Equal(remA, remB) {
			return false
		}
	}

	// Step 6
	return true
}
