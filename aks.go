package primality

import (
	"math"
	"slices"
)

// Check n != a^b for a,b > 1, returning true if it is otherwise false if not.
func basePowerCheck(n int) bool {
	for i := 2.0; i < math.Log2(float64(n)); i++ {
		a := math.Pow(float64(n), 1.0/float64(i))
		// Check if the rounded value a is equal to its integer value.
		if i == a {
			return true // n is a composite number.
		}
	}
	return false
}

// Find the order such that Ord_r(n) > (log_2(n))^2
func findCorrectOrder(n int) int {
	maxK := math.Floor(math.Pow(math.Log2(float64(n)), 2))
	maxR := max(3.0, math.Ceil(math.Pow(math.Log2(float64(n)), 5)))
	nextR := true
	i := 2.0
	for i = 2.0; nextR && i < maxR; i++ {
		nextR = false
		for j := 1.0; !nextR && j <= maxK; j++ {
			nextR = math.Mod(math.Pow(float64(n), j), i) == 1 ||
				math.Mod(math.Pow(float64(n), j), i) == 0
		}
	}
	return int(i - 1.0)
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

// eulerTotient is an implementation of Euler's Totient function
func eulerTotient(n int) int {
	res := n // Initialize result as n
	// Check up to the square root of n for factors of n
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			for n%p == 0 {
				n /= p // remove all factors p from n
			}
			res -= res / p // remove the current result's quotient of p
		}
	}
	if n > 1 {
		res -= res / n // remove the current result's quotient of n
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
		if mod != 0 {
			r[i] = mod
		}
	}
	return r
}

// polynomialExpansion finds the coefficients of the polynomial expansion of
// (x+a)^e and returns the coefficients in ascending powers of x from x = 0
func polynomialExpansion(e, a int) []int {
	c := make([][]int, e)
	for i := range c {
		c[i] = make([]int, i+2)
	}
	for n := 0; n < e; n++ {
		c[n][0] = 1
		c[n][len(c[n])-1] = 1
	}
	for n := 1; n < e; n++ {
		for k := 1; k < len(c[n]); k++ {
			if k > n {
				break
			}
			c[n][k] = (c[n-1][k-1] + c[n-1][k])
		}
	}
	for i, x := range c[e-1] {
		c[e-1][i] = x * int(math.Pow(float64(a), float64(e-i)))
	}
	return c[e-1]
}

// degree determines the degree of the polynomial p
//
//	Ref: https://rosettacode.org/wiki/Polynomial_long_division#Go
func degree(p []int) int {
	for d := len(p) - 1; d >= 0; d-- {
		if p[d] != 0 {
			return d
		}
	}
	return -1
}

// pld performs polynomial long division on the two polynomial coefficient slices
// provided. It expects the polynomials to be in ascending order of powers of x.
//
//	Ref: https://rosettacode.org/wiki/Polynomial_long_division#Go
func pld(nn, dd []int) (q, r []int, ok bool) {
	if degree(dd) < 0 {
		return
	}
	nn = append(r, nn...)
	if degree(nn) >= degree(dd) {
		q = make([]int, degree(nn)-degree(dd)+1)
		for degree(nn) >= degree(dd) {
			d := make([]int, degree(nn)+1)
			copy(d[degree(nn)-degree(dd):], dd)
			q[degree(nn)-degree(dd)] = nn[degree(nn)] / d[degree(d)]
			for i := range d {
				d[i] *= q[degree(nn)-degree(dd)]
				nn[i] -= d[i]
			}
		}
	}
	return q, nn, true
}

// polynomialModRemainder finds the remainder of polynomials p1/p2 and does a
// term-wise reduction modulo m on the result, returning a slice of coefficients
// for a polynomial in ascending order of x.
func polynomialModRemainder(p1, p2 []int, m int) []int {
	_, r, ok := pld(p1, p2)
	if !ok {
		return nil
	}
	return polynomialMod(r, m)
}

// polynomialSubtraction subtracts p1 from p2 in a term-wise fashion.
// The function orders the polynomials longest first prior to subtraction.
func polynomialSubtraction(p1, p2 []int) []int {
	longest, shortest := p1, p2
	if len(p2) > len(p2) {
		longest, shortest = p2, p1
	}
	res := make([]int, len(longest)+1)
	for i, x := range shortest {
		res[i] = longest[i] - x
	}
	copy(res[len(shortest):], longest[len(shortest)-1:])
	return res
}

// AKS is an implementation of the AKS deterministic primality test.
// Step 5 takes up the majority of the time and as such results in a slow test.
func AKS(n int) bool {
	// Initial check
	if n < 2 {
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
	for a := 1; a <= maxA; a++ {
		xa := polynomialExpansion(n, a)
		xna := polynomialExpansion(n, 1)
		xna[0] += a
		xr1 := polynomialExpansion(r, 1)
		xr1[0]--
		remA := polynomialModRemainder(xa, xr1, n)
		_, remB, ok := pld(xna, xr1)
		if !ok {
			panic("error dividing polynomials")
		}
		longest := make([]int, len(remA))
		if len(remB) > len(remA) {
			longest = make([]int, len(remB))
		}
		if slices.Equal(polynomialSubtraction(remA, remB), longest) {
			return false
		}
	}

	// Step 6
	return true
}
